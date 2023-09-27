package dlq

import (
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
	conf := &entropy.Config{}
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

type dlqJobReqBody struct {
	ErrorTypes string `json:"error_types,omitempty"`
	BatchSize  int64  `json:"batch_size,omitempty"`
	BlobBatch  int64  `json:"blob_batch,omitempty"`
	NumThreads int64  `json:"num_threads,omitempty"`
	Topic      string `json:"topic,omitempty"`
}

func (h *Handler) createDlqJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var def dlqJobReqBody

	if err := utils.ReadJSON(r, &def); err != nil {
		utils.WriteErr(w, err)
		return
	}

	dlq_job, err := h.service.mapDlqJob(def, ctx)
	if err != nil {
		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"error": err,
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"dlq_job": dlq_job,
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
