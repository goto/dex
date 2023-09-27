package optimus

import (
	"context"

	optimusv1beta1grpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/optimus/core/v1beta1/corev1beta1grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OptimusClient struct {
}

type OptimusClientBuilder interface {
	BuildOptimusClient(ctx context.Context, hostname string) (optimusv1beta1grpc.JobSpecificationServiceClient, error)
}

func (client *OptimusClient) BuildOptimusClient(ctx context.Context, hostname string) (optimusv1beta1grpc.JobSpecificationServiceClient, error) {
	optimusConn, err := grpc.Dial(hostname, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return optimusv1beta1grpc.NewJobSpecificationServiceClient(optimusConn), nil
}
