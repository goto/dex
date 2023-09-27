package dlq

import (
	entropyv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/entropy/v1beta1/entropyv1beta1grpc"

	"github.com/goto/dex/internal/server/gcs"
)

type Service struct {
	client    entropyv1beta1rpc.ResourceServiceClient
	gcsClient gcs.BlobStorageClient
}

func NewService(client entropyv1beta1rpc.ResourceServiceClient, gcsClient gcs.BlobStorageClient) *Service {
	return &Service{
		client:    client,
		gcsClient: gcsClient,
	}
}
