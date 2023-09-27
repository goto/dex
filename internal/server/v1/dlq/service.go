package dlq

import (
	"context"

	entropyv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/entropy/v1beta1/entropyv1beta1grpc"
	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"

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

func (s *Service) mapDlqJob(def dlqJobReqBody, ctx context.Context) (*entropyv1beta1.CreateResourceResponse, error) {
	res := &entropyv1beta1.Resource{
		Kind: "kube",
	}
	rpcReq := &entropyv1beta1.CreateResourceRequest{Resource: res}
	rpcResp, err := s.client.CreateResource(ctx, rpcReq)
	return rpcResp, err
}
