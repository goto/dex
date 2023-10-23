package dlq

import (
	"errors"
	"log"
	"net/http"

	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"github.com/go-chi/chi/v5"

	"github.com/goto/dex/entropy"
	"github.com/goto/dex/internal/server/gcs"
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

func (*Handler) listDlqJobs(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"dlq_jobs": []interface{}{},
	})
}

func (*Handler) createDlqJob(w http.ResponseWriter, _ *http.Request) {
	// transform request body into DlqJob (validation?)
	// call service.CreateDLQJob

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"dlq_job": nil,
	})
}

func (h *Handler) GetDlqJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	jobURN := h.jobURN(r)

	dlqJob, err := h.service.GetDlqJob(ctx, jobURN)
	if err != nil {
		if errors.Is(err, ErrJobNotFound) {
			utils.WriteErrMsg(w, http.StatusNotFound, ErrJobNotFound.Error())
			return
		}
		utils.WriteErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, dlqJob)
}

func (*Handler) firehoseURN(r *http.Request) string {
	return chi.URLParam(r, "firehose_urn")
}

func (*Handler) jobURN(r *http.Request) string {
	return chi.URLParam(r, "job_urn")
}
