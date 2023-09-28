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
		return nil, args[1].(error)
	}
	return args.Get(0).(optimusv1beta1grpc.JobSpecificationServiceClient), args.Error(1)
}
