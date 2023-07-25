package firehose

import (
	"context"
	"net/http"
	"time"

	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"github.com/go-chi/chi/v5"
	firehose "github.com/goto/entropy/modules/firehose"
	kafka "github.com/goto/entropy/pkg/kafka"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goto/dex/generated/models"
	"github.com/goto/dex/internal/server/reqctx"
	"github.com/goto/dex/internal/server/utils"
	"github.com/goto/dex/pkg/errors"
)

const (
	actionStop        = "stop"
	actionScale       = "scale"
	actionStart       = "start"
	actionUpgrade     = "upgrade"
	actionResetOffset = "reset"
)

//nolint:dupl
func (api *firehoseAPI) handleReset(w http.ResponseWriter, r *http.Request) {
	var reqBody kafka.ResetParams

	if err := utils.ReadJSON(r, &reqBody); err != nil {
		utils.WriteErr(w, err)
		return
	}

	urn := chi.URLParam(r, pathParamURN)
	existingFirehose, err := api.getFirehose(r.Context(), urn)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}

	updatedFirehose, err := api.executeAction(r.Context(), existingFirehose, actionResetOffset, reqBody)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, updatedFirehose)
}

//nolint:dupl
func (api *firehoseAPI) handleScale(w http.ResponseWriter, r *http.Request) {
	var reqBody firehose.ScaleParams

	if err := utils.ReadJSON(r, &reqBody); err != nil {
		utils.WriteErr(w, err)
		return
	}

	// Ensure that the URN refers to a valid firehose resource.
	urn := chi.URLParam(r, pathParamURN)
	existingFirehose, err := api.getFirehose(r.Context(), urn)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}

	updatedFirehose, err := api.executeAction(r.Context(), existingFirehose, actionScale, reqBody)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, updatedFirehose)
}

func (api *firehoseAPI) handleStart(w http.ResponseWriter, r *http.Request) {
	var reqBody firehose.StartParams

	// Ensure that the URN refers to a valid firehose resource.
	urn := chi.URLParam(r, pathParamURN)
	existingFirehose, err := api.getFirehose(r.Context(), urn)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}

	// for LOG sinkType, updating stop_time
	if existingFirehose.Configs.EnvVars[confSinkType] == logSinkType {
		t := time.Now().UTC().Add(logSinkTTL)
		reqBody.StopTime = &t
	}

	updatedFirehose, err := api.executeAction(r.Context(), existingFirehose, actionStart, reqBody)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, updatedFirehose)
}

func (api *firehoseAPI) handleStop(w http.ResponseWriter, r *http.Request) {
	var reqBody struct{}
	if err := utils.ReadJSON(r, &reqBody); err != nil {
		utils.WriteErr(w, err)
		return
	}

	// Ensure that the URN refers to a valid firehose resource.
	urn := chi.URLParam(r, pathParamURN)
	existingFirehose, err := api.getFirehose(r.Context(), urn)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}

	updatedFirehose, err := api.executeAction(r.Context(), existingFirehose, actionStop, reqBody)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}

	if err := api.stopAlerts(r.Context(), *updatedFirehose, projectSlugFromURN(urn)); err != nil {
		utils.WriteErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, updatedFirehose)
}

func (api *firehoseAPI) handleUpgrade(w http.ResponseWriter, r *http.Request) {
	var reqBody struct{}
	if err := utils.ReadJSON(r, &reqBody); err != nil {
		utils.WriteErr(w, err)
		return
	}

	// Ensure that the URN refers to a valid firehose resource.
	urn := chi.URLParam(r, pathParamURN)
	existingFirehose, err := api.getFirehose(r.Context(), urn)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}

	updatedFirehose, err := api.executeAction(r.Context(), existingFirehose, actionUpgrade, reqBody)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, updatedFirehose)
}

func (api *firehoseAPI) executeAction(ctx context.Context, existingFirehose *models.Firehose, actionType string, params any) (*models.Firehose, error) {
	reqCtx := reqctx.From(ctx)

	paramStruct, err := utils.GoValToProtoStruct(params)
	if err != nil {
		return nil, err
	}

	labels := cloneAndMergeMaps(existingFirehose.Labels, map[string]string{
		labelUpdatedBy: reqCtx.UserEmail,
	})

	rpcReq := &entropyv1beta1.ApplyActionRequest{
		Urn:    existingFirehose.Urn,
		Action: actionType,
		Params: paramStruct,
		Labels: labels,
	}

	rpcResp, err := api.Entropy.ApplyAction(ctx, rpcReq)
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.InvalidArgument {
			return nil, errors.ErrInvalid.WithMsgf(st.Message())
		} else if st.Code() == codes.NotFound {
			return nil, errFirehoseNotFound.WithMsgf(st.Message())
		}
		return nil, err
	}

	return mapEntropyResourceToFirehose(rpcResp.GetResource())
}
