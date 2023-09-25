package alert

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/goto/dex/generated/client/operations"
	"github.com/goto/dex/generated/models"
	"github.com/goto/dex/internal/server/reqctx"
	"github.com/goto/dex/internal/server/utils"
)

type Handler struct {
	subscriptionService *SubscriptionService
	alertService        *Service
}

func NewHandler(subscriptionService *SubscriptionService, alertService *Service) *Handler {
	return &Handler{
		subscriptionService: subscriptionService,
		alertService:        alertService,
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
		"subscription": MapToSubscription(subscription),
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
		"subscriptions": MapToSubscriptionList(subscriptions),
	})
}

func (h *Handler) createSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqCtx := reqctx.From(ctx)

	userEmail := reqCtx.UserEmail
	if userEmail == "" {
		utils.WriteErrMsg(w, http.StatusUnauthorized, "identity headers are required")
		return
	}

	var requestPayload models.SubscriptionForm
	if err := utils.ReadJSON(r, &requestPayload); err != nil {
		utils.WriteErr(w, err)
		return
	}
	if err := requestPayload.Validate(nil); err != nil {
		utils.WriteErrMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	form := SubscriptionForm{
		UserID:             userEmail,
		ChannelCriticality: ChannelCriticality(*requestPayload.ChannelCriticality),
		AlertSeverity:      AlertSeverity(*requestPayload.AlertSeverity),
		ProjectID:          *requestPayload.ProjectID,
		GroupID:            *requestPayload.GroupID,
		ResourceID:         *requestPayload.ResourceID,
		ResourceType:       *requestPayload.ResourceType,
	}
	subscriptionID, err := h.subscriptionService.CreateSubscription(ctx, form)
	if err != nil {
		if errors.Is(err, ErrNoShieldGroup) || errors.Is(err, ErrNoShieldOrg) || errors.Is(err, ErrNoShieldProject) {
			utils.WriteErrMsg(w, http.StatusNotFound, err.Error())
		} else if errors.Is(err, ErrNoShieldParentSlackReceiver) ||
			errors.Is(err, ErrInvalidShieldParentSlackReceiver) ||
			errors.Is(err, ErrNoShieldSirenNamespace) ||
			errors.Is(err, ErrInvalidShieldSirenNamespace) ||
			errors.Is(err, ErrNoSirenReceiver) {
			utils.WriteErrMsg(w, http.StatusUnprocessableEntity, err.Error())
		} else {
			utils.WriteErr(w, err)
		}

		return
	}

	subscription, err := h.subscriptionService.FindSubscription(ctx, subscriptionID)
	if err != nil {
		utils.WriteErr(w, fmt.Errorf("error finding subscription: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"subscription": MapToSubscription(subscription),
	})
}

func (h *Handler) updateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqCtx := reqctx.From(ctx)

	userEmail := reqCtx.UserEmail
	if userEmail == "" {
		utils.WriteErrMsg(w, http.StatusUnauthorized, "identity headers are required")
		return
	}

	subscriptionIDStr := chi.URLParam(r, "subscription_id")
	subscriptionID, err := strconv.Atoi(subscriptionIDStr)
	if err != nil {
		utils.WriteErrMsg(w, http.StatusBadRequest, "subscription identifier has to be a number")
		return
	}

	var requestPayload models.SubscriptionForm
	if err := utils.ReadJSON(r, &requestPayload); err != nil {
		utils.WriteErr(w, err)
		return
	}
	if err := requestPayload.Validate(nil); err != nil {
		utils.WriteErrMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	form := SubscriptionForm{
		UserID:             userEmail,
		ChannelCriticality: ChannelCriticality(*requestPayload.ChannelCriticality),
		AlertSeverity:      AlertSeverity(*requestPayload.AlertSeverity),
		ProjectID:          *requestPayload.ProjectID,
		GroupID:            *requestPayload.GroupID,
		ResourceID:         *requestPayload.ResourceID,
		ResourceType:       *requestPayload.ResourceType,
	}

	if err := h.subscriptionService.UpdateSubscription(ctx, subscriptionID, form); err != nil {
		if errors.Is(err, ErrNoShieldGroup) || errors.Is(err, ErrNoShieldOrg) || errors.Is(err, ErrNoShieldProject) {
			utils.WriteErrMsg(w, http.StatusNotFound, err.Error())
		} else if errors.Is(err, ErrNoShieldParentSlackReceiver) ||
			errors.Is(err, ErrInvalidShieldParentSlackReceiver) ||
			errors.Is(err, ErrNoShieldSirenNamespace) ||
			errors.Is(err, ErrInvalidShieldSirenNamespace) ||
			errors.Is(err, ErrNoSirenReceiver) {
			utils.WriteErrMsg(w, http.StatusUnprocessableEntity, err.Error())
		} else {
			utils.WriteErr(w, err)
		}
		return
	}

	subscription, err := h.subscriptionService.FindSubscription(ctx, subscriptionID)
	if err != nil {
		utils.WriteErr(w, fmt.Errorf("error finding subscription: %w", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"subscription": MapToSubscription(subscription),
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

func (h *Handler) getAlertChannels(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	groupID := chi.URLParam(r, "group_id")

	alertChannels, err := h.subscriptionService.GetAlertChannels(ctx, groupID)
	if err != nil {
		if errors.Is(err, ErrNoShieldGroup) {
			utils.WriteErrMsg(w, http.StatusNotFound, err.Error())
		} else {
			utils.WriteErr(w, err)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"alert_channels": alertChannels,
	})
}

func (h *Handler) setAlertChannels(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqCtx := reqctx.From(ctx)
	userEmail := reqCtx.UserEmail
	if userEmail == "" {
		utils.WriteErrMsg(w, http.StatusUnauthorized, "identity headers are required")
		return
	}

	groupID := chi.URLParam(r, "group_id")

	var body operations.SetGroupAlertChannelsBody
	if err := utils.ReadJSON(r, &body); err != nil {
		utils.WriteErr(w, err)
		return
	}
	if err := body.Validate(nil); err != nil {
		utils.WriteErrMsg(w, http.StatusBadRequest, err.Error())
		return
	}

	var alertChannelsForms []AlertChannelForm
	for _, ac := range body.AlertChannels {
		alertChannelsForms = append(alertChannelsForms, AlertChannelForm{
			ChannelCriticality:  ChannelCriticality(*ac.ChannelCriticality),
			ChannelName:         ac.ChannelName,
			ChannelType:         string(*ac.ChannelType),
			PagerdutyServiceKey: ac.PagerdutyServiceKey,
		})
	}

	alertChannels, err := h.subscriptionService.SetAlertChannels(ctx, userEmail, groupID, alertChannelsForms)
	if err != nil {
		if errors.Is(err, ErrNoShieldGroup) || errors.Is(err, ErrNoShieldOrg) {
			utils.WriteErrMsg(w, http.StatusNotFound, err.Error())
		} else if errors.Is(err, ErrNoShieldParentSlackReceiver) ||
			errors.Is(err, ErrInvalidShieldParentSlackReceiver) ||
			errors.Is(err, ErrNoShieldSirenNamespace) ||
			errors.Is(err, ErrInvalidShieldSirenNamespace) {
			utils.WriteErrMsg(w, http.StatusUnprocessableEntity, err.Error())
		} else {
			utils.WriteErr(w, err)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"alert_channels": alertChannels,
	})
}

func (h *Handler) getAlerts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Println("here")
	projectSlug := chi.URLParam(r, "project_slug")
	resourceUrn := chi.URLParam(r, "resource_urn")

	alerts, err := h.alertService.ListAlerts(ctx, projectSlug, resourceUrn)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK,
		utils.ListResponse[Alert]{Items: alerts})
}

func (h *Handler) getAlertPolicy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resourceType := chi.URLParam(r, "resource_type")
	projectSlug := chi.URLParam(r, "project_slug")
	resourceUrn := chi.URLParam(r, "resource_urn")

	policy, err := h.alertService.GetAlertPolicy(ctx, projectSlug, resourceUrn, resourceType)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}

	var suppliedAlertVariableNames = []string{"name", "team", "entity"}

	policy.Rules = RemoveSuppliedVariablesFromRules(policy.Rules, suppliedAlertVariableNames)

	utils.WriteJSON(w, http.StatusOK, policy)
}
