package optimus

import (
	"errors"
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

func (h *Handler) findJob(w http.ResponseWriter, r *http.Request) {
	jobName := chi.URLParam(r, "job_name")
	projectSlug := chi.URLParam(r, "project_slug")

	jobSpecResp, err := h.service.FindJobSpec(r.Context(), jobName, projectSlug)

	if err != nil {
		if errors.Is(err, ErrOptimusHostNotFound) {
			utils.WriteErrMsg(w, http.StatusNotFound, ErrOptimusHostNotFound.Error())
			return
		} else if errors.Is(err, ErrOptimusHostInvalid) {
			utils.WriteErrMsg(w, http.StatusUnprocessableEntity, ErrOptimusHostInvalid.Error())
			return
		}
		utils.WriteErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, jobSpecResp)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	projectSlug := chi.URLParam(r, "project_slug")

	listResp, err := h.service.ListJobs(r.Context(), projectSlug)
	if err != nil {
		if errors.Is(err, ErrOptimusHostNotFound) {
			utils.WriteErrMsg(w, http.StatusNotFound, ErrOptimusHostNotFound.Error())
			return
		} else if errors.Is(err, ErrOptimusHostInvalid) {
			utils.WriteErrMsg(w, http.StatusUnprocessableEntity, ErrOptimusHostInvalid.Error())
			return
		}
		utils.WriteErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, listResp)
}
