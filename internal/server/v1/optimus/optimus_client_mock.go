package optimus

import (
	"context"

	optimusv1beta1grpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/optimus/core/v1beta1/corev1beta1grpc"
	"github.com/stretchr/testify/mock"
)

type OptimusClientMock struct {
	mock.Mock
}

func (mock *OptimusClientMock) BuildOptimusClient(ctx context.Context, hostname string) (optimusv1beta1grpc.JobSpecificationServiceClient, error) {
	args := mock.Called(ctx, hostname)
	return args.Get(0).(optimusv1beta1grpc.JobSpecificationServiceClient), args.Error(1)
}
