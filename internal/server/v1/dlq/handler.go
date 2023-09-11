package dlq

import (
	"context"
	"fmt"
	"net/http"

	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"github.com/go-chi/chi/v5"
	entropyFirehose "github.com/goto/entropy/modules/firehose"

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

func (h *Handler) listFirehoseDLQ(w http.ResponseWriter, r *http.Request) {
	firehoseURN := h.firehoseURN(r)
	resp, err := h.service.client.GetResource(context.Background(), &entropyv1beta1.GetResourceRequest{Urn: firehoseURN})
	if err != nil {
		utils.WriteErr(w, err)
		fmt.Println(err)
		return
	}
	conf := &entropyFirehose.Config{}
	err = utils.ProtoStructToGoVal(resp.GetResource().GetSpec().GetConfigs(), conf)
	if err != nil {
		utils.WriteErr(w, err)
		fmt.Println(err)
		return
	}
	// check the variables for dlq related config.
	bucketName := conf.EnvVariables[firehose.ConfigDLQBucket]
	directoryPrefix := conf.EnvVariables[firehose.ConfigDLQDirectoryPrefix]
	topicDates, err := h.service.gcsClient.ListTopicDates(gcs.BucketInfo{
		BucketName: bucketName,
		Prefix:     directoryPrefix,
		Delim:      "",
	})
	if err != nil {
		utils.WriteErr(w, err)
		fmt.Println(err)
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
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"dlq_job": nil,
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
