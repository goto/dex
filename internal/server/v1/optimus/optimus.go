package optimus

import (
	"context"

	optimusv1beta1grpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/optimus/core/v1beta1/corev1beta1grpc"
	optimusv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/optimus/core/v1beta1"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goto/dex/pkg/errors"
)

func Routes(optimus optimusv1beta1grpc.JobSpecificationServiceClient) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/projects/{project_name}/jobs/{job_name}", handleGetOptimus(optimus))
	}
}

func getOptimus(ctx context.Context, optimusClient optimusv1beta1grpc.JobSpecificationServiceClient, jobName, projectName string) (*optimusv1beta1.GetJobSpecificationResponse, error) {
	res, err := optimusClient.GetJobSpecification(ctx, &optimusv1beta1.GetJobSpecificationRequest{
		ProjectName:   projectName,
		JobName:       jobName,
		NamespaceName: "smoke_test",
	})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return res, nil
}
