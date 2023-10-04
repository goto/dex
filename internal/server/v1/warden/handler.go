package warden

import (
	"net/http"

	"github.com/goto/dex/internal/server/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) teamList(w http.ResponseWriter, r *http.Request) {
	teamListResp, err := h.service.TeamList(r.Context())
	if err != nil {
		utils.WriteErr(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, teamListResp)
}
