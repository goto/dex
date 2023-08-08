package optimus

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/goto/dex/internal/server/utils"
	"github.com/goto/dex/pkg/errors"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) findJob(w http.ResponseWriter, r *http.Request) {
	jobName := chi.URLParam(r, "job_name")
	projectName := chi.URLParam(r, "project_name")

	jobSpecResp, err := h.service.FindJobSpec(r.Context(), jobName, projectName)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, jobSpecResp)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	projectName := q.Get("project")

	if projectName == "" {
		utils.WriteErr(w, errors.ErrInvalid.WithMsgf("project query param is required"))
		return
	}

	listResp, err := h.service.ListJobs(r.Context(), projectName)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, listResp)
}
