package optimus

import (
	optimusv1beta1grpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/optimus/core/v1beta1/corev1beta1grpc"
	"github.com/stretchr/testify/mock"
)

type OptimusClientMock struct {
	mock.Mock
}

type OptimusClientBuilderMock interface {
	BuildOptimusClient(hostname string) (optimusv1beta1grpc.JobSpecificationServiceClient, error)
}

func (mock *OptimusClientMock) BuildOptimusClient(hostname string) (optimusv1beta1grpc.JobSpecificationServiceClient, error) {
	args := mock.Called(hostname)
	if args[0] == nil {
		return nil, args.Error(1)
	}

	if client, ok := args[0].(optimusv1beta1grpc.JobSpecificationServiceClient); ok {
		return client, args.Error(1)
	}
	return nil, args.Error(1)
}
