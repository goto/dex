package dlq

import (
	"context"
	"testing"

	entropyv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/entropy/v1beta1/entropyv1beta1grpc"
	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	"github.com/goto/dex/internal/server/gcs"
)

type entropyService struct {
	mock.Mock
}

func (e *entropyService) ListResources(ctx context.Context, in *entropyv1beta1.ListResourcesRequest, opts ...grpc.CallOption) (*entropyv1beta1.ListResourcesResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (e *entropyService) GetResource(ctx context.Context, in *entropyv1beta1.GetResourceRequest, opts ...grpc.CallOption) (*entropyv1beta1.GetResourceResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (e *entropyService) CreateResource(ctx context.Context, in *entropyv1beta1.CreateResourceRequest, opts ...grpc.CallOption) (*entropyv1beta1.CreateResourceResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (e *entropyService) UpdateResource(ctx context.Context, in *entropyv1beta1.UpdateResourceRequest, opts ...grpc.CallOption) (*entropyv1beta1.UpdateResourceResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (e *entropyService) DeleteResource(ctx context.Context, in *entropyv1beta1.DeleteResourceRequest, opts ...grpc.CallOption) (*entropyv1beta1.DeleteResourceResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (e *entropyService) ApplyAction(ctx context.Context, in *entropyv1beta1.ApplyActionRequest, opts ...grpc.CallOption) (*entropyv1beta1.ApplyActionResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (e *entropyService) GetLog(ctx context.Context, in *entropyv1beta1.GetLogRequest, opts ...grpc.CallOption) (entropyv1beta1rpc.ResourceService_GetLogClient, error) {
	// TODO implement me
	panic("implement me")
}

func (e *entropyService) GetResourceRevisions(ctx context.Context, in *entropyv1beta1.GetResourceRevisionsRequest, opts ...grpc.CallOption) (*entropyv1beta1.GetResourceRevisionsResponse, error) {
	// TODO implement me
	panic("implement me")
}

type gcsClient struct {
	mock.Mock
}

func (g *gcsClient) ListTopicDates(bucketInfo gcs.BucketInfo) (map[string]map[string]int64, error) {
	// TODO implement me
	panic("implement me")
}

func TestListTopicDates(t *testing.T) {
	_ = Handler{service: NewService(&entropyService{}, &gcsClient{})}
}
