package optimus

import (
	shieldv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/shield/v1beta1/shieldv1beta1grpc"
	"github.com/go-chi/chi/v5"
)

func Routes(shieldClient shieldv1beta1rpc.ShieldServiceClient) func(r chi.Router) {
	service := NewService(shieldClient, &OptimusClient{})
	handler := NewHandler(service)

	return func(r chi.Router) {
		r.Get("/projects/{project_slug}/jobs/{job_name}", handler.findJob)
		r.Get("/projects/{project_slug}/jobs", handler.list)
	}
}
