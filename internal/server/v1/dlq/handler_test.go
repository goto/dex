package dlq_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/goto/dex/entropy"
	"github.com/goto/dex/generated/models"
	"github.com/goto/dex/internal/server/gcs"
	"github.com/goto/dex/internal/server/reqctx"
	"github.com/goto/dex/internal/server/utils"
	"github.com/goto/dex/internal/server/v1/dlq"
	"github.com/goto/dex/internal/server/v1/firehose"
	"github.com/goto/dex/mocks"
)

const (
	emailHeaderKey = "X-Auth-Email"
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
	handler := dlq.NewHandler(dlq.NewService(eService, gClient, dlq.DlqJobConfig{}))
	httpWriter := &testHTTPWriter{}
	httpRequest := &http.Request{}
	config := &entropy.FirehoseConfig{
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
		InitContainer: entropy.FirehoseInitContainer{},
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
	topicDates := []models.DlqMetadata{
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
	gClient.On("ListDlqMetadata", gcs.BucketInfo{
		BucketName: "test-bucket",
		Prefix:     "test-prefix",
		Delim:      "",
	}).Return(topicDates, nil)
	handler.ListFirehoseDLQ(httpWriter, httpRequest)
	expectedMap := make(map[string][]models.DlqMetadata)
	err := json.Unmarshal([]byte(httpWriter.messages[0]), &expectedMap)
	require.NoError(t, err)
	assert.Equal(t, topicDates, expectedMap["dlq_list"])
}

func TestErrorFromGCSClient(t *testing.T) {
	eService := &mocks.ResourceServiceClient{}
	gClient := &mocks.BlobStorageClient{}
	handler := dlq.NewHandler(dlq.NewService(eService, gClient, dlq.DlqJobConfig{}))
	httpWriter := &testHTTPWriter{}
	httpRequest := &http.Request{}
	config := &entropy.FirehoseConfig{
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
		InitContainer: entropy.FirehoseInitContainer{},
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
	gClient.On("ListDlqMetadata", gcs.BucketInfo{
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
	handler := dlq.NewHandler(dlq.NewService(eService, gClient, dlq.DlqJobConfig{}))
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

func TestCreateDlqJob(t *testing.T) {
	var (
		method        = http.MethodPost
		path          = fmt.Sprintf("/jobs")
		resource_id   = "test-resource-id"
		resource_type = "test-resource-type"
		error_types   = "DESERIALIZATION_ERROR"
		date          = "21-10-2022"
		batch_size    = 0
		num_threads   = 0
		topic         = "test-topic"
		jsonPayload   = fmt.Sprintf(`{
			"resource_id": "%s",
			"resource_type": "%s",
			"error_types": "%s",
			"batch_size": %d,
			"num_threads": %d,
			"topic": "%s",
			"date": "%s"
		}`, resource_id, resource_type, error_types, batch_size, num_threads, topic, date)
	)

	t.Run("Should return resource urn", func(t *testing.T) {
		// initt input
		namespace := "test-namespace"
		kubeCluster := "test-kube-cluster"
		userEmail := "test@example.com"
		config := dlq.DlqJobConfig{
			PrometheusHost: "http://sample-prom-host",
			DlqJobImage:    "test-image",
		}
		envVars := map[string]string{
			"SINK_TYPE":                               "bigquery",
			"DLQ_BATCH_SIZE":                          "34",
			"DLQ_NUM_THREADS":                         "10",
			"DLQ_PREFIX_DIR":                          "test-firehose",
			"DLQ_FINISHED_STATUS_FILE":                "/shared/job-finished",
			"DLQ_GCS_BUCKET_NAME":                     "g-pilotdata-gl-dlq",
			"DLQ_ERROR_TYPES":                         "DEFAULT_ERROR",
			"DLQ_GCS_CREDENTIAL_PATH":                 "/etc/secret/gcp/token",
			"DLQ_GCS_GOOGLE_CLOUD_PROJECT_ID":         "pilotdata-integration",
			"DLQ_INPUT_DATE":                          "2023-04-10",
			"JAVA_TOOL_OPTIONS":                       "-javaagent:jolokia-jvm-agent.jar=port=8778,host=localhost",
			"_JAVA_OPTIONS":                           "-Xmx1800m -Xms1800m",
			"DLQ_TOPIC_NAME":                          "gofood-booking-log",
			"INPUT_SCHEMA_PROTO_CLASS":                "gojek.esb.booking.GoFoodBookingLogMessage",
			"METRIC_STATSD_TAGS":                      "a=b",
			"SCHEMA_REGISTRY_STENCIL_ENABLE":          "true",
			"SCHEMA_REGISTRY_STENCIL_URLS":            "http://p-godata-systems-stencil-v1beta1-ingress.golabs.io/v1beta1/namespaces/gojek/schemas/esb-log-entities",
			"SINK_BIGQUERY_ADD_METADATA_ENABLED":      "true",
			"SINK_BIGQUERY_CLIENT_CONNECT_TIMEOUT_MS": "-1",
			"SINK_BIGQUERY_CLIENT_READ_TIMEOUT_MS":    "-1",
			"SINK_BIGQUERY_CREDENTIAL_PATH":           "/etc/secret/gcp/token",
			"SINK_BIGQUERY_DATASET_LABELS":            "shangchi=legend,lord=voldemort",
			"SINK_BIGQUERY_DATASET_LOCATION":          "US",
			"SINK_BIGQUERY_DATASET_NAME":              "bq_test",
			"SINK_BIGQUERY_GOOGLE_CLOUD_PROJECT_ID":   "pilotdata-integration",
			"SINK_BIGQUERY_ROW_INSERT_ID_ENABLE":      "false",
			"SINK_BIGQUERY_STORAGE_API_ENABLE":        "true",
			"SINK_BIGQUERY_TABLE_LABELS":              "hello=world,john=doe",
			"SINK_BIGQUERY_TABLE_NAME":                "bq_dlq_test1",
			"SINK_BIGQUERY_TABLE_PARTITION_EXPIRY_MS": "2629800000",
			"SINK_BIGQUERY_TABLE_PARTITION_KEY":       "event_timestamp",
			"SINK_BIGQUERY_TABLE_PARTITIONING_ENABLE": "true",
		}

		outputStruct, err := structpb.NewStruct(map[string]interface{}{
			"namespace": namespace,
		})
		require.NoError(t, err)

		firehoseConfig, err := utils.GoValToProtoStruct(entropy.FirehoseConfig{
			EnvVariables: envVars,
		})
		require.NoError(t, err)

		firehoseResource := &entropyv1beta1.Resource{
			Spec: &entropyv1beta1.ResourceSpec{
				Dependencies: []*entropyv1beta1.ResourceDependency{
					{
						Key:   "kube_cluster",
						Value: kubeCluster,
					},
				},
				Configs: firehoseConfig,
			},
			State: &entropyv1beta1.ResourceState{
				Output: structpb.NewStructValue(outputStruct),
			},
		}

		jobResource := &entropyv1beta1.Resource{
			Urn: "test-urn",
			State: &entropyv1beta1.ResourceState{
				Output: structpb.NewStructValue(outputStruct),
			},
		}

		expectedEnvVars := map[string]string{
			"DLQ_BATCH_SIZE":                          fmt.Sprintf("%d", batch_size),
			"DLQ_NUM_THREADS":                         fmt.Sprintf("%d", num_threads),
			"DLQ_ERROR_TYPES":                         error_types,
			"DLQ_INPUT_DATE":                          date,
			"DLQ_TOPIC_NAME":                          topic,
			"METRIC_STATSD_TAGS":                      "a=b", // TBA
			"SINK_TYPE":                               envVars["SINK_TYPE"],
			"DLQ_PREFIX_DIR":                          "test-firehose",
			"DLQ_FINISHED_STATUS_FILE":                "/shared/job-finished",
			"DLQ_GCS_BUCKET_NAME":                     envVars["DLQ_GCS_BUCKET_NAME"],
			"DLQ_GCS_CREDENTIAL_PATH":                 envVars["DLQ_GCS_CREDENTIAL_PATH"],
			"DLQ_GCS_GOOGLE_CLOUD_PROJECT_ID":         envVars["DLQ_GCS_GOOGLE_CLOUD_PROJECT_ID"],
			"JAVA_TOOL_OPTIONS":                       envVars["JAVA_TOOL_OPTIONS"],
			"_JAVA_OPTIONS":                           envVars["_JAVA_OPTIONS"],
			"INPUT_SCHEMA_PROTO_CLASS":                envVars["INPUT_SCHEMA_PROTO_CLASS"],
			"SCHEMA_REGISTRY_STENCIL_ENABLE":          envVars["SCHEMA_REGISTRY_STENCIL_ENABLE"],
			"SCHEMA_REGISTRY_STENCIL_URLS":            envVars["SCHEMA_REGISTRY_STENCIL_URLS"],
			"SINK_BIGQUERY_ADD_METADATA_ENABLED":      envVars["SINK_BIGQUERY_ADD_METADATA_ENABLED"],
			"SINK_BIGQUERY_CLIENT_CONNECT_TIMEOUT_MS": envVars["SINK_BIGQUERY_CLIENT_CONNECT_TIMEOUT_MS"],
			"SINK_BIGQUERY_CLIENT_READ_TIMEOUT_MS":    envVars["SINK_BIGQUERY_CLIENT_READ_TIMEOUT_MS"],
			"SINK_BIGQUERY_CREDENTIAL_PATH":           envVars["SINK_BIGQUERY_CREDENTIAL_PATH"],
			"SINK_BIGQUERY_DATASET_LABELS":            envVars["SINK_BIGQUERY_DATASET_LABELS"],
			"SINK_BIGQUERY_DATASET_LOCATION":          envVars["SINK_BIGQUERY_DATASET_LOCATION"],
			"SINK_BIGQUERY_DATASET_NAME":              envVars["SINK_BIGQUERY_DATASET_NAME"],
			"SINK_BIGQUERY_GOOGLE_CLOUD_PROJECT_ID":   envVars["SINK_BIGQUERY_GOOGLE_CLOUD_PROJECT_ID"],
			"SINK_BIGQUERY_ROW_INSERT_ID_ENABLE":      envVars["SINK_BIGQUERY_ROW_INSERT_ID_ENABLE"],
			"SINK_BIGQUERY_STORAGE_API_ENABLE":        envVars["SINK_BIGQUERY_STORAGE_API_ENABLE"],
			"SINK_BIGQUERY_TABLE_LABELS":              envVars["SINK_BIGQUERY_TABLE_LABELS"],
			"SINK_BIGQUERY_TABLE_NAME":                envVars["SINK_BIGQUERY_TABLE_NAME"],
			"SINK_BIGQUERY_TABLE_PARTITION_EXPIRY_MS": envVars["SINK_BIGQUERY_TABLE_PARTITION_EXPIRY_MS"],
			"SINK_BIGQUERY_TABLE_PARTITION_KEY":       envVars["SINK_BIGQUERY_TABLE_PARTITION_KEY"],
			"SINK_BIGQUERY_TABLE_PARTITIONING_ENABLE": envVars["SINK_BIGQUERY_TABLE_PARTITIONING_ENABLE"],
		}

		jobConfig, err := utils.GoValToProtoStruct(entropy.JobConfig{
			Replicas:  1,
			Namespace: namespace,
			Containers: []entropy.JobContainer{
				{
					Name:            "dlq-job",
					Image:           config.DlqJobImage,
					ImagePullPolicy: "Always",
					SecretsVolumes: []entropy.JobSecret{
						{
							Name:  "firehose-bigquery-sink-credential",
							Mount: envVars["DLQ_GCS_CREDENTIAL_PATH"],
						},
					},
					Limits: entropy.UsageSpec{
						CPU:    "0.5", // user
						Memory: "2gb", // user
					},
					Requests: entropy.UsageSpec{
						CPU:    "0.5", // user
						Memory: "2gb", // user
					},
					EnvVariables: expectedEnvVars,
				},
				{
					Name:  "telegraf",
					Image: "telegraf:1.18.0-alpine",
					ConfigMapsVolumes: []entropy.JobConfigMap{
						{
							Name:  "dlq-processor-telegraf",
							Mount: "/etc/telegraf",
						},
					},
					EnvVariables: map[string]string{
						// To be updated by streaming
						"APP_NAME":        "", // TBA
						"PROMETHEUS_HOST": config.PrometheusHost,
						"DEPLOYMENT_NAME": "deployment-name",
						"TEAM":            "",
						"TOPIC":           topic,
						"environment":     "production", // TBA
						"organization":    "de",         // TBA
						"projectID":       "",
					},
					Command: []string{
						"/bin/bash",
					},
					Args: []string{
						"-c",
						"telegraf & while [ ! -f /shared/job-finished ]; do sleep 5; done; sleep 20 && exit 0",
					},
					Limits: entropy.UsageSpec{
						CPU:    "100m",  // user
						Memory: "300Mi", // user
					},
					Requests: entropy.UsageSpec{
						CPU:    "100m",  // user
						Memory: "300Mi", // user
					},
				},
			},
			JobLabels: map[string]string{
				"firehose": resource_id,
				"topic":    topic,
				"date":     date,
			},
			Volumes: []entropy.JobVolume{
				{
					Name: "firehose-bigquery-sink-credential",
					Kind: "secret",
				},
				{
					Name: "dlq-processor-telegraf",
					Kind: "configMap",
				},
			},
		})
		require.NoError(t, err)

		newJobResourcePayload := &entropyv1beta1.Resource{
			Urn:  "",
			Kind: entropy.ResourceKindJob,
			Name: fmt.Sprintf(
				"%s-%s-%s-%s",
				resource_id,   // firehose urn
				resource_type, // firehose / dagger
				topic,         //
				date,          //
			),
			Project: firehoseResource.Project,
			Labels: map[string]string{
				"resource_id": resource_id,
				"type":        resource_type,
				"date":        date,
				"topic":       topic,
				"job_type":    "dlq",
			},
			CreatedBy: jobResource.CreatedBy,
			UpdatedBy: jobResource.UpdatedBy,
			Spec: &entropyv1beta1.ResourceSpec{
				Configs: jobConfig,
				Dependencies: []*entropyv1beta1.ResourceDependency{
					{
						Key:   "kube_cluster",
						Value: kubeCluster, // from firehose configs.kube_cluster
					},
				},
			},
		}

		// set conditions
		entropyClient := new(mocks.ResourceServiceClient)
		entropyClient.On(
			"GetResource", mock.Anything, &entropyv1beta1.GetResourceRequest{Urn: resource_id},
		).Return(&entropyv1beta1.GetResourceResponse{
			Resource: firehoseResource,
		}, nil)
		entropyClient.On("CreateResource", mock.Anything, &entropyv1beta1.CreateResourceRequest{
			Resource: newJobResourcePayload,
		}).Return(&entropyv1beta1.CreateResourceResponse{
			Resource: jobResource,
		}, nil)
		defer entropyClient.AssertExpectations(t)

		assert.NoError(t, err)
		requestBody := bytes.NewReader([]byte(jsonPayload))

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, requestBody)
		request.Header.Set(emailHeaderKey, userEmail)
		router := getRouter()
		dlq.Routes(entropyClient, nil, config)(router)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		resultJSON := response.Body.Bytes()
		expectedJSON, err := json.Marshal(map[string]interface{}{
			"dlq_list": "test-urn",
		})
		require.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), string(resultJSON))
	})
}

func getRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(reqctx.WithRequestCtx())

	return router
}
