package dlq

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/goto/dex/internal/server/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) listDlq(w http.ResponseWriter, r *http.Request) {
	// sample to get firehose urn from route params
	_ = h.firehoseURN(r)

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"dlq_list": []interface{}{},
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
