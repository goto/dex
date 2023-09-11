package dlq

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	entropyv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/entropy/v1beta1/entropyv1beta1grpc"
	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	"github.com/goto/dex/entropy"
	"github.com/goto/dex/internal/server/gcs"
	"github.com/goto/dex/internal/server/utils"
	"github.com/goto/dex/internal/server/v1/firehose"
)

type mockEntropyService struct {
	mock.Mock
}

func (*mockEntropyService) ListResources(context.Context, *entropyv1beta1.ListResourcesRequest, ...grpc.CallOption) (*entropyv1beta1.ListResourcesResponse, error) {
	panic("")
}

func (e *mockEntropyService) GetResource(ctx context.Context, in *entropyv1beta1.GetResourceRequest, opts ...grpc.CallOption) (*entropyv1beta1.GetResourceResponse, error) {
	args := e.Called(ctx, in, opts)
	if args.Get(0) == nil {
		err, _ := args.Get(1).(error)
		return nil, err
	}
	response, _ := args.Get(0).(*entropyv1beta1.GetResourceResponse)
	return response, nil
}

func (*mockEntropyService) CreateResource(context.Context, *entropyv1beta1.CreateResourceRequest, ...grpc.CallOption) (*entropyv1beta1.CreateResourceResponse, error) {
	panic("")
}

func (*mockEntropyService) UpdateResource(context.Context, *entropyv1beta1.UpdateResourceRequest, ...grpc.CallOption) (*entropyv1beta1.UpdateResourceResponse, error) {
	panic("")
}

func (*mockEntropyService) DeleteResource(context.Context, *entropyv1beta1.DeleteResourceRequest, ...grpc.CallOption) (*entropyv1beta1.DeleteResourceResponse, error) {
	panic("")
}

func (*mockEntropyService) ApplyAction(context.Context, *entropyv1beta1.ApplyActionRequest, ...grpc.CallOption) (*entropyv1beta1.ApplyActionResponse, error) {
	panic("")
}

func (*mockEntropyService) GetLog(context.Context, *entropyv1beta1.GetLogRequest, ...grpc.CallOption) (entropyv1beta1rpc.ResourceService_GetLogClient, error) {
	panic("")
}

func (*mockEntropyService) GetResourceRevisions(context.Context, *entropyv1beta1.GetResourceRevisionsRequest, ...grpc.CallOption) (*entropyv1beta1.GetResourceRevisionsResponse, error) {
	panic("")
}

type mockGcsClient struct {
	mock.Mock
}

func (g *mockGcsClient) ListTopicDates(bucketInfo gcs.BucketInfo) (map[string]map[string]int64, error) {
	args := g.Called(bucketInfo)
	if args.Get(0) == nil {
		err, _ := args.Get(1).(error)
		return nil, err
	}
	topicDates, _ := args.Get(0).(map[string]map[string]int64)
	return topicDates, nil
}

type mockHTTPWriter struct {
	mock.Mock
	messages []string
}

func (*mockHTTPWriter) Header() http.Header {
	return http.Header{}
}

func (m *mockHTTPWriter) Write(bytes []byte) (int, error) {
	m.messages = append(m.messages, string(bytes[:]))
	return 0, nil
}

func (*mockHTTPWriter) WriteHeader(int) {
}

func TestListTopicDates(t *testing.T) {
	eService := &mockEntropyService{}
	config := &entropy.Config{
		Stopped:      false,
		StopTime:     nil,
		Replicas:     0,
		Namespace:    "",
		DeploymentID: "",
		EnvVariables: map[string]string{
			firehose.ConfigDLQBucket:          "test-bucket",
			firehose.ConfigDLQDirectoryPrefix: "test-prefix",
		},
		ResetOffset:   "",
		Limits:        entropy.UsageSpec{},
		Requests:      entropy.UsageSpec{},
		Telegraf:      nil,
		ChartValues:   nil,
		InitContainer: entropy.InitContainer{},
	}
	configProto, _ := utils.GoValToProtoStruct(config)
	eService.On(
		"GetResource",
		context.Background(),
		mock.Anything,
		[]grpc.CallOption(nil)).Return(
		&entropyv1beta1.GetResourceResponse{
			Resource: &entropyv1beta1.Resource{
				Urn:       "",
				Kind:      "",
				Name:      "",
				Project:   "",
				Labels:    nil,
				CreatedAt: nil,
				UpdatedAt: nil,
				Spec: &entropyv1beta1.ResourceSpec{
					Configs:      configProto,
					Dependencies: nil,
				},
				State:     nil,
				CreatedBy: "",
				UpdatedBy: "",
			},
		}, nil)
	gClient := &mockGcsClient{}
	topicDates := make(map[string]map[string]int64)
	topicDates["topic-1"] = make(map[string]int64)
	topicDates["topic-2"] = make(map[string]int64)
	topicDates["topic-1"]["2023-08-26"] = int64(1234)
	topicDates["topic-1"]["2023-12-10"] = int64(4321)
	topicDates["topic-2"]["2023-09-20"] = int64(99)
	gClient.On("ListTopicDates", gcs.BucketInfo{
		BucketName: "test-bucket",
		Prefix:     "test-prefix",
		Delim:      "",
	}).Return(topicDates, nil)
	handler := Handler{service: NewService(eService, gClient)}
	httpWriter := &mockHTTPWriter{}
	httpRequest := &http.Request{}
	handler.listFirehoseDLQ(httpWriter, httpRequest)
	expectedMap := make(map[string]map[string]map[string]int64)
	_ = json.Unmarshal([]byte(httpWriter.messages[0]), &expectedMap)
	assert.Equal(t, topicDates, expectedMap["dlq_list"])
}
