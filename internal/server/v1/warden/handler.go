package warden

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/goto/dex/internal/server/utils"
)

type handler struct {
	service *Service
}

func NewHandler(service *Service) *handler {
	return &handler{service: service}
}

func (h *handler) teamList(w http.ResponseWriter, r *http.Request) {
	teamListResp, err := h.service.TeamList(r.Context())

	if errors.Is(err, ErrUserNotFound) {
		utils.WriteErrMsg(w, http.StatusUnauthorized, ErrUserNotFound.Error())
		return
	}
	if errors.Is(err, ErrTeamNotFound) {
		utils.WriteErrMsg(w, http.StatusNotFound, ErrTeamNotFound.Error())
		return
	}
	if err != nil {
		utils.WriteErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, teamListResp)
}

func (h *handler) updateGroupMetadata(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "group_id")

	var body struct {
		WardeTeamID string `json:"warden_team_id"`
	}
	if err := utils.ReadJSON(r, &body); err != nil {
		utils.WriteErr(w, err)
		return
	} else if body.WardeTeamID == "" {
		utils.WriteErrMsg(w, http.StatusBadRequest, "missing warden_team_id")
		return
	}

	resShield, err := h.service.UpdateGroupMetadata(r.Context(), groupID, body.WardeTeamID)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resShield)
}
