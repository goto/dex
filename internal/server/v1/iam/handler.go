package iam

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/goto/dex/internal/server/reqctx"
	"github.com/goto/dex/internal/server/utils"
)

type handler struct {
	service *Service
}

func NewHandler(service *Service) *handler {
	return &handler{service: service}
}

func (h *handler) listUserWardenTeams(w http.ResponseWriter, r *http.Request) {
	reqCtx := reqctx.From(r.Context())
	const errEmailMissedInHeader = "user email not in header"

	if reqCtx.UserEmail == "" {
		utils.WriteErrMsg(w, http.StatusUnauthorized, errEmailMissedInHeader)
		return
	}

	teamListResp, err := h.service.UserWardenTeamList(r.Context(), reqCtx.UserEmail)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"teams": teamListResp,
	})
}

func (h *handler) linkGroupToWarden(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "group_id")

	var body struct {
		WardenTeamID string `json:"warden_team_id"`
	}
	if err := utils.ReadJSON(r, &body); err != nil {
		utils.WriteErr(w, err)
		return
	} else if body.WardenTeamID == "" {
		utils.WriteErrMsg(w, http.StatusBadRequest, "missing warden_team_id")
		return
	}

	resShield, err := h.service.LinkGroupToWarden(r.Context(), groupID, body.WardenTeamID)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resShield)
}
