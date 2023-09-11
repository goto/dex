package dlq

import (
	"log"

	entropyv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/entropy/v1beta1/entropyv1beta1grpc"
	"github.com/go-chi/chi/v5"

	"github.com/goto/dex/internal/server/gcs"
)

func Routes(entropyClient entropyv1beta1rpc.ResourceServiceClient, gcsKeyFilePath string) func(r chi.Router) {
	gcsClient, err := gcs.NewClient(gcsKeyFilePath)
	if err != nil {
		log.Fatalf("Error while creating GCS %s\n", err.Error())
	}
	service := NewService(entropyClient, gcsClient)
	handler := NewHandler(service)

	return func(r chi.Router) {
		r.Get("/firehose/{firehose_urn}", handler.listFirehoseDLQ)
		r.Get("/jobs", handler.listDlqJobs)
		r.Get("/jobs/{job_urn}", handler.getDlqJob)
		r.Post("/jobs", handler.createDlqJob)
	}
}
