package dlq

import (
	"errors"
	"log"
	"net/http"

	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"github.com/go-chi/chi/v5"

	"github.com/goto/dex/entropy"
	"github.com/goto/dex/generated/models"
	"github.com/goto/dex/internal/server/gcs"
	"github.com/goto/dex/internal/server/reqctx"
	"github.com/goto/dex/internal/server/utils"
	"github.com/goto/dex/internal/server/v1/firehose"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListFirehoseDLQ(w http.ResponseWriter, r *http.Request) {
	firehoseURN := h.firehoseURN(r)
	resp, err := h.service.client.GetResource(r.Context(), &entropyv1beta1.GetResourceRequest{Urn: firehoseURN})
	if err != nil {
		utils.WriteErr(w, err)
		log.Println(err)
		return
	}
	conf := &entropy.FirehoseConfig{}
	err = utils.ProtoStructToGoVal(resp.GetResource().GetSpec().GetConfigs(), conf)
	if err != nil {
		utils.WriteErr(w, err)
		log.Println(err)
		return
	}
	// check the variables for dlq related config.
	bucketName := conf.EnvVariables[firehose.ConfigDLQBucket]
	directoryPrefix := conf.EnvVariables[firehose.ConfigDLQDirectoryPrefix]
	topicDates, err := h.service.gcsClient.ListDlqMetadata(gcs.BucketInfo{
		BucketName: bucketName,
		Prefix:     directoryPrefix,
		Delim:      "",
	})
	if err != nil {
		utils.WriteErr(w, err)
		log.Println(err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"dlq_list": topicDates,
	})
}

func (h *Handler) listDlqJobs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	firehoseUrn := chi.URLParam(r, "firehoseURN")
	// fetch py resource (kind = job)
	// mapToDlqJob(entropyResource) -> DqlJob
	dlqJob, err := h.service.listDlqJob(ctx, firehoseUrn)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, dlqJob)
}

func (h *Handler) createDlqJob(w http.ResponseWriter, r *http.Request) {
	// transform request body into DlqJob (validation?)
	ctx := r.Context()
	reqCtx := reqctx.From(ctx)
	var dlqJob models.DlqJob

	if err := utils.ReadJSON(r, &dlqJob); err != nil {
		utils.WriteErr(w, err)
		return
	}

	// call service.CreateDLQJob
	err := h.service.CreateDLQJob(ctx, reqCtx.UserEmail, &dlqJob)
	if err != nil {
		if errors.Is(err, ErrFirehoseNotFound) {
			utils.WriteErrMsg(w, http.StatusNotFound, err.Error())
			return
		}
		utils.WriteErr(w, err)
		return
	}
	// return
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"dlq_list": dlqJob.Urn,
	})
}

func (h *Handler) getDlqJob(w http.ResponseWriter, r *http.Request) {
	// sample to get job urn from route params
	_ = h.jobURN(r)

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"dlq_job": nil,
	})
}

func (*Handler) firehoseURN(r *http.Request) string {
	return chi.URLParam(r, "firehose_urn")
}

func (*Handler) jobURN(r *http.Request) string {
	return chi.URLParam(r, "job_urn")
}
