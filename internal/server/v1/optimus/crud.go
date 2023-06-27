package optimus

import (
	optimusv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/optimus/core/v1beta1/corev1beta1grpc"
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/goto/dex/internal/server/utils"
)

func handleGetOptimus(client optimusv1beta1rpc.JobSpecificationServiceClient) http.HandlerFunc {
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
