package optimus

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	optimusv1beta1 "github.com/goto/optimus/protos/gotocompany/optimus/core/v1beta1"

	"github.com/goto/dex/internal/server/utils"
)

func handleGetOptimus(client optimusv1beta1.JobSpecificationServiceClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobName := chi.URLParam(r, "job_name")
		projectName := chi.URLParam(r, "project_name")

		res, err := getOptimus(r.Context(), client, jobName, projectName)
		if err != nil {
			utils.WriteErr(w, err)
			return
		}
		utils.WriteJSON(w, http.StatusOK, res)
	}
}
