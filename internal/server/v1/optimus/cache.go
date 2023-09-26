package optimus

import (
	"sync"

	optimusv1beta1grpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/optimus/core/v1beta1/corev1beta1grpc"
)

type Cache struct {
	mu   sync.RWMutex
	data map[string]optimusv1beta1grpc.JobSpecificationServiceClient
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]optimusv1beta1grpc.JobSpecificationServiceClient),
	}
}
