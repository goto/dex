package dlq

import (
	entropyv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/entropy/v1beta1/entropyv1beta1grpc"
	"github.com/go-chi/chi/v5"

	"github.com/goto/dex/internal/server/gcs"
)

func Routes(
	entropyClient entropyv1beta1rpc.ResourceServiceClient,
	gcsClient gcs.BlobStorageClient,
	cfg *DlqJobConfig,
) func(r chi.Router) {
	service := NewService(entropyClient, gcsClient, cfg)
	handler := NewHandler(service)

	return func(r chi.Router) {
		r.Get("/firehose/{firehose_urn}", handler.ListFirehoseDLQ)
		r.Get("/jobs", handler.listDlqJobs)
		r.Get("/jobs/{job_urn}", handler.GetDlqJob)
		r.Post("/jobs", handler.createDlqJob)
	}
}
