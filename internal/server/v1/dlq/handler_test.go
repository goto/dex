package dlq_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/goto/dex/entropy"
	"github.com/goto/dex/internal/server/gcs"
	"github.com/goto/dex/internal/server/utils"
	"github.com/goto/dex/internal/server/v1/dlq"
	"github.com/goto/dex/internal/server/v1/firehose"
	"github.com/goto/dex/mocks"
)

type testHTTPWriter struct {
	messages []string
}

func (*testHTTPWriter) Header() http.Header {
	return http.Header{}
}

func (m *testHTTPWriter) Write(bytes []byte) (int, error) {
	m.messages = append(m.messages, string(bytes[:]))
	return 0, nil
}

func (*testHTTPWriter) WriteHeader(int) {
}

func TestListTopicDates(t *testing.T) {
	eService := &mocks.ResourceServiceClient{}
	gClient := &mocks.BlobStorageClient{}
	handler := dlq.NewHandler(dlq.NewService(eService, gClient))
	httpWriter := &testHTTPWriter{}
	httpRequest := &http.Request{}
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
		&entropyv1beta1.GetResourceRequest{Urn: ""}).Return(
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
	topicDates := []gcs.TopicMetaData{
		{
			Topic:       "topic-1",
			Date:        "2023-08-26",
			SizeInBytes: 1234,
		},
		{
			Topic:       "test-topic2",
			Date:        "2023-12-10",
			SizeInBytes: 4321,
		},
		{
			Topic:       "topic-2",
			Date:        "2023-09-20",
			SizeInBytes: 99,
		},
	}
	gClient.On("ListTopicDates", gcs.BucketInfo{
		BucketName: "test-bucket",
		Prefix:     "test-prefix",
		Delim:      "",
	}).Return(topicDates, nil)
	handler.ListFirehoseDLQ(httpWriter, httpRequest)
	expectedMap := make(map[string][]gcs.TopicMetaData)
	err := json.Unmarshal([]byte(httpWriter.messages[0]), &expectedMap)
	require.NoError(t, err)
	assert.Equal(t, topicDates, expectedMap["dlq_list"])
}

func TestErrorFromGCSClient(t *testing.T) {
	eService := &mocks.ResourceServiceClient{}
	gClient := &mocks.BlobStorageClient{}
	handler := dlq.NewHandler(dlq.NewService(eService, gClient))
	httpWriter := &testHTTPWriter{}
	httpRequest := &http.Request{}
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
		&entropyv1beta1.GetResourceRequest{Urn: ""}).Return(
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
	gClient.On("ListTopicDates", gcs.BucketInfo{
		BucketName: "test-bucket",
		Prefix:     "test-prefix",
		Delim:      "",
	}).Return(nil, fmt.Errorf("test-error"))
	handler.ListFirehoseDLQ(httpWriter, httpRequest)
	expectedMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(httpWriter.messages[0]), &expectedMap)
	require.NoError(t, err)
	assert.Equal(t, "test-error", expectedMap["cause"])
}

func TestErrorFromFirehoseResource(t *testing.T) {
	eService := &mocks.ResourceServiceClient{}
	gClient := &mocks.BlobStorageClient{}
	handler := dlq.NewHandler(dlq.NewService(eService, gClient))
	httpWriter := &testHTTPWriter{}
	httpRequest := &http.Request{}
	eService.On(
		"GetResource",
		context.Background(),
		mock.Anything).Return(nil, fmt.Errorf("test-error"))
	handler.ListFirehoseDLQ(httpWriter, httpRequest)
	expectedMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(httpWriter.messages[0]), &expectedMap)
	require.NoError(t, err)
	assert.Equal(t, "test-error", expectedMap["cause"])
}
