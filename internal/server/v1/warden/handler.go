package warden

import (
	"errors"
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

func (h *handler) teamList(w http.ResponseWriter, r *http.Request) {
	reqCtx := reqctx.From(r.Context())
	const errEmailMissedInHeader = "user email not in header"

	if reqCtx.UserEmail == "" {
		utils.WriteErrMsg(w, http.StatusUnauthorized, errEmailMissedInHeader)
		return
	}

	teamListResp, err := h.service.TeamList(r.Context(), reqCtx.UserEmail)
	if err != nil {
		if errors.Is(err, ErrEmailNotOnWarden) {
			utils.WriteErrMsg(w, http.StatusNotFound, ErrEmailNotOnWarden.Error())
			return
		}
		utils.WriteErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, teamListResp)
}

func (h *handler) updateGroupMetadata(w http.ResponseWriter, r *http.Request) {
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

	resShield, err := h.service.UpdateGroupMetadata(r.Context(), groupID, body.WardenTeamID)
	if err != nil {
		utils.WriteErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, resShield)
}
