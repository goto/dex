package alert

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	shieldv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/shield/v1beta1/shieldv1beta1grpc"
	shieldv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/shield/v1beta1"
	"github.com/go-chi/chi/v5"

	"github.com/goto/dex/generated/models"
	"github.com/goto/dex/internal/server/reqctx"
	"github.com/goto/dex/internal/server/utils"
)

type Handler struct {
	subscriptionService *SubscriptionService
	shieldClient        shieldv1beta1rpc.ShieldServiceClient
}

func NewHandler(subscriptionService *SubscriptionService, shieldClient shieldv1beta1rpc.ShieldServiceClient) *Handler {
	return &Handler{
		subscriptionService: subscriptionService,
		shieldClient:        shieldClient,
	}
}

func (h *Handler) findSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	subscriptionIDStr := chi.URLParam(r, "subscription_id")
	subscriptionID, err := strconv.Atoi(subscriptionIDStr)
	if err != nil {
		utils.WriteErrMsg(w, http.StatusBadRequest, "subscription identifier has to be a number")
		return
	}

	subscription, err := h.subscriptionService.FindSubscription(ctx, subscriptionID)
	if err != nil {
		if errors.Is(err, ErrSubscriptionNotFound) {
			utils.WriteErrMsg(w, http.StatusNotFound, err.Error())
		} else {
			utils.WriteErr(w, err)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"subscription": subscription,
	})
}

func (h *Handler) getSubscriptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	groupID := strings.TrimSpace(r.URL.Query().Get("group_id"))
	resourceID := strings.TrimSpace(r.URL.Query().Get("resource_id"))
	resourceType := strings.TrimSpace(r.URL.Query().Get("resource_type"))

	if groupID == "" && resourceID == "" {
		utils.WriteErrMsg(w, http.StatusBadRequest, "requires either groupID or a combination of resource_id and resource_type")
		return
	}

	subscriptions, err := h.subscriptionService.GetSubscriptions(ctx, groupID, resourceID, resourceType)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"subscriptions": subscriptions,
	})
}

func (h *Handler) createSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqCtx := reqctx.From(ctx)

	var form models.SubscriptionForm
	if err := utils.ReadJSON(r, &form); err != nil {
		utils.WriteErr(w, err)
		return
	}
	if err := form.Validate(nil); err != nil {
		utils.WriteErr(w, err)
		return
	}

	channelName, err := h.getSlackChannelByCriticality(ctx, *form.GroupID, ChannelCriticality(*form.ChannelCriticality))
	if err != nil {
		utils.WriteErr(w, fmt.Errorf("error getting slack channel: %w", err))
		return
	}

	subscriptionID, err := h.subscriptionService.CreateSubscription(ctx, form, channelName, reqCtx.UserEmail)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}

	subscription, err := h.subscriptionService.FindSubscription(ctx, subscriptionID)
	if err != nil {
		utils.WriteErr(w, fmt.Errorf("error finding subscription: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"subscription": subscription,
	})
}

func (h *Handler) updateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqCtx := reqctx.From(ctx)

	subscriptionIDStr := chi.URLParam(r, "subscription_id")
	subscriptionID, err := strconv.Atoi(subscriptionIDStr)
	if err != nil {
		utils.WriteErrMsg(w, http.StatusBadRequest, "subscription identifier has to be a number")
		return
	}

	var form models.SubscriptionForm
	if err := utils.ReadJSON(r, &form); err != nil {
		utils.WriteErr(w, err)
		return
	}
	if err := form.Validate(nil); err != nil {
		utils.WriteErr(w, err)
		return
	}

	channelName, err := h.getSlackChannelByCriticality(ctx, *form.GroupID, ChannelCriticality(*form.ChannelCriticality))
	if err != nil {
		utils.WriteErr(w, fmt.Errorf("error getting slack channel: %w", err))
		return
	}

	if err := h.subscriptionService.UpdateSubscription(ctx, subscriptionID, form, channelName, reqCtx.UserEmail); err != nil {
		utils.WriteErr(w, err)
		return
	}

	subscription, err := h.subscriptionService.FindSubscription(ctx, subscriptionID)
	if err != nil {
		utils.WriteErr(w, fmt.Errorf("error finding subscription: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"subscription": subscription,
	})
}

func (h *Handler) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	subscriptionIDStr := chi.URLParam(r, "subscription_id")
	subscriptionID, err := strconv.Atoi(subscriptionIDStr)
	if err != nil {
		utils.WriteErrMsg(w, http.StatusBadRequest, "subscription identifier has to be a number")
		return
	}

	if err := h.subscriptionService.DeleteSubscription(ctx, subscriptionID); err != nil {
		if errors.Is(err, ErrSubscriptionNotFound) {
			utils.WriteErrMsg(w, http.StatusNotFound, err.Error())
		} else {
			utils.WriteErr(w, err)
		}

		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "subscription removed",
	})
}

func (h *Handler) getSlackChannelByCriticality(ctx context.Context, groupID string, criticality ChannelCriticality) (string, error) {
	resp, err := h.shieldClient.GetGroup(ctx, &shieldv1beta1.GetGroupRequest{
		Id: groupID,
	})
	if err != nil {
		return "", fmt.Errorf("error getting a group: %w", err)
	}

	group := resp.GetGroup()
	groupMetadata := group.GetMetadata().AsMap()

	// get slack metadata
	slack, exists := groupMetadata["slack"]
	if !exists {
		return "", ErrNoShieldSlackMetadata
	}
	slackMap, ok := slack.(map[string]interface{})
	if !ok {
		return "", ErrInvalidShieldSlackMetadata
	}

	// get channel name
	channelNameAny, exists := slackMap[string(criticality)]
	if !exists {
		return "", ErrNoShieldSlackChannel
	}
	channelName, ok := channelNameAny.(string)
	if !ok {
		return "", ErrInvalidSlackChannelFormat
	}

	return channelName, nil
}
