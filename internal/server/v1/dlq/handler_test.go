package dlq_test

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/goto/dex/entropy"
	"github.com/goto/dex/generated/models"
	"github.com/goto/dex/internal/server/gcs"
	"github.com/goto/dex/internal/server/reqctx"
	"github.com/goto/dex/internal/server/utils"
	"github.com/goto/dex/internal/server/v1/dlq"
	"github.com/goto/dex/internal/server/v1/firehose"
	"github.com/goto/dex/mocks"
)

//go:embed fixtures/list_dlq_jobs.json
var listDlqJobsFixtureJSON []byte

// nolint
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

func TestListDlqJob(t *testing.T) {
	var (
		method       = http.MethodGet
		project      = "test-project-1"
		resourceID   = "test-resource-id"
		resourceType = "test-resource-type"
		kubeCluster  = "test-kube-cluster"
		date         = "2022-10-21"
		topic        = "test-topic"
		group        = "test-group"
	)

	t.Run("Should return error in firehose mapping", func(t *testing.T) {
		path := "/jobs?resource_id="
		expectedErr := status.Error(codes.Internal, "Not found")
		expectedLabels := map[string]string{}
		entropyClient := new(mocks.ResourceServiceClient)
		entropyClient.On(
			"ListResources", mock.Anything, &entropyv1beta1.ListResourcesRequest{
				Kind: entropy.ResourceKindJob, Labels: expectedLabels,
			},
		).Return(nil, expectedErr)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := getRouter()
		dlq.Routes(entropyClient, nil, dlq.DlqJobConfig{})(router)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

	t.Run("Should return list of dlqJobs", func(t *testing.T) {
		path := fmt.Sprintf("/jobs?resource_id=%s&resource_type=%s&date=%s", resourceID, resourceType, date)

		dummyEntropyResources := []*entropyv1beta1.Resource{
			{
				Urn:  "test-urn-1",
				Kind: entropy.ResourceKindJob,
				Name: fmt.Sprintf(
					"%s-%s-%s-%s",
					"test1",    // firehose title
					"firehose", // firehose / dagger
					topic,      //
					date,       //
				),
				Project: project,
				Labels: map[string]string{
					"resource_id":     resourceID,
					"resource_type":   resourceType,
					"date":            date,
					"topic":           topic,
					"job_type":        "dlq",
					"group":           group,
					"prometheus_host": "prom_host",
				},
				CreatedBy: "user@test.com",
				UpdatedBy: "user@test.com",
				CreatedAt: timestamppb.New(time.Date(2022, time.December, 10, 0, 0, 0, 0, time.UTC)),
				UpdatedAt: timestamppb.New(time.Date(2023, time.December, 10, 2, 0, 0, 0, time.UTC)),
				Spec: &entropyv1beta1.ResourceSpec{
					Dependencies: []*entropyv1beta1.ResourceDependency{
						{
							Key:   "kube_cluster",
							Value: kubeCluster, // from firehose configs.kube_cluster
						},
					},
				},
			},
			{
				Urn:  "test-urn-2",
				Kind: entropy.ResourceKindJob,
				Name: fmt.Sprintf(
					"%s-%s-%s-%s",
					"test2",    // firehose title
					"firehose", // firehose / dagger
					topic,      //
					date,       //
				),
				Project: project,
				Labels: map[string]string{
					"resource_id":     resourceID,
					"resource_type":   resourceType,
					"date":            date,
					"topic":           topic,
					"job_type":        "dlq",
					"group":           group,
					"prometheus_host": "prom_host2",
				},
				CreatedBy: "user@test.com",
				UpdatedBy: "user@test.com",
				CreatedAt: timestamppb.New(time.Date(2012, time.October, 10, 4, 0, 0, 0, time.UTC)),
				UpdatedAt: timestamppb.New(time.Date(2013, time.February, 12, 2, 4, 0, 0, time.UTC)),
				Spec: &entropyv1beta1.ResourceSpec{
					Dependencies: []*entropyv1beta1.ResourceDependency{
						{
							Key:   "kube_cluster",
							Value: kubeCluster, // from firehose configs.kube_cluster
						},
					},
				},
			},
		}

		expectedRPCResp := &entropyv1beta1.ListResourcesResponse{
			Resources: dummyEntropyResources,
		}

		expectedLabels := map[string]string{
			"resource_id":   "test-resource-id",
			"resource_type": "test-resource-type",
			"date":          "2022-10-21",
		}

		entropyClient := new(mocks.ResourceServiceClient)
		entropyClient.On(
			"ListResources", mock.Anything, &entropyv1beta1.ListResourcesRequest{
				Kind: entropy.ResourceKindJob, Labels: expectedLabels,
			},
		).Return(expectedRPCResp, nil)
		defer entropyClient.AssertExpectations(t)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := getRouter()
		dlq.Routes(entropyClient, nil, dlq.DlqJobConfig{})(router)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		resultJSON := response.Body.Bytes()

		expectedPayload := string(listDlqJobsFixtureJSON)
		assert.JSONEq(t, expectedPayload, string(resultJSON))
	})
}

// nolint
func skipTestCreateDlqJob(t *testing.T) {
	var (
		method       = http.MethodPost
		path         = "/jobs"
		resourceID   = "test-resource-id"
		resourceType = "test-resource-type"
		errorTypes   = "DESERIALIZATION_ERROR"
		date         = "2022-10-2"
		batchSize    = 0
		numThreads   = 0
		topic        = "test-topic"
		group        = ""
		jsonPayload  = fmt.Sprintf(`{
			"resource_id": "%s",
			"resource_type": "%s",
			"error_types": "%s",
			"batch_size": %d,
			"num_threads": %d,
			"topic": "%s",
			"date": "%s",
			"group": "%s"
		}`, resourceID, resourceType, errorTypes, batchSize, numThreads, topic, date, group)
	)

	t.Run("Should return error user header is required", func(t *testing.T) {
		requestBody := bytes.NewReader([]byte(jsonPayload))

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, requestBody)
		router := getRouter()
		dlq.Routes(nil, nil, dlq.DlqJobConfig{})(router)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Code)
	})

	t.Run("Should return error validation resource type not firehose", func(t *testing.T) {
		jsonBody := `{
			"resource_id": "test-resource-id",
			"resource_type": "dlq",
			"error_types": "test-error",
			"batch_size": 1,
			"num_threads": 2,
			"topic": "test-topic",
			"date": "2022-10-21",
			"group": "test-group"
		}`
		userEmail := "test@example.com"

		requestBody := bytes.NewReader([]byte(jsonBody))

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, requestBody)
		request.Header.Set(emailHeaderKey, userEmail)
		router := getRouter()
		dlq.Routes(nil, nil, dlq.DlqJobConfig{})(router)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Should return error validation missing required fields", func(t *testing.T) {
		jsonBody := `{
			"resource_id": "",
			"resource_type": "firehose",
			"error_types": "test-error",
			"batch_size": 1,
			"num_threads": 2,
			"topic": "test-topic",
			"date": "2022-10-21",
			"group": "test-group"
		}`
		userEmail := "test@example.com"

		requestBody := bytes.NewReader([]byte(jsonBody))

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, requestBody)
		request.Header.Set(emailHeaderKey, userEmail)
		router := getRouter()
		dlq.Routes(nil, nil, dlq.DlqJobConfig{})(router)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Should return error firehose not Found", func(t *testing.T) {
		expectedErr := status.Error(codes.NotFound, "Not found")
		entropyClient := new(mocks.ResourceServiceClient)
		entropyClient.On(
			"GetResource", mock.Anything, &entropyv1beta1.GetResourceRequest{Urn: resourceID},
		).Return(nil, expectedErr)
		defer entropyClient.AssertExpectations(t)

		requestBody := bytes.NewReader([]byte(jsonPayload))

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, requestBody)
		router := getRouter()
		dlq.Routes(entropyClient, nil, dlq.DlqJobConfig{})(router)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("Should return error in firehose mapping", func(t *testing.T) {
		expectedErr := status.Error(codes.Internal, "Not found")
		entropyClient := new(mocks.ResourceServiceClient)
		entropyClient.On(
			"GetResource", mock.Anything, &entropyv1beta1.GetResourceRequest{Urn: resourceID},
		).Return(nil, expectedErr)
		defer entropyClient.AssertExpectations(t)

		requestBody := bytes.NewReader([]byte(jsonPayload))

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, requestBody)
		router := getRouter()
		dlq.Routes(entropyClient, nil, dlq.DlqJobConfig{})(router)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

	t.Run("Should return resource urn", func(t *testing.T) {
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
			"DLQ_BATCH_SIZE":                          fmt.Sprintf("%d", batchSize),
			"DLQ_NUM_THREADS":                         fmt.Sprintf("%d", numThreads),
			"DLQ_ERROR_TYPES":                         errorTypes,
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
				"firehose": resourceID,
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
				firehoseResource.Name, // firehose urn
				"firehose",            // firehose / dagger
				topic,                 //
				date,                  //
			),
			Project: firehoseResource.Project,
			Labels: map[string]string{
				"resource_id":     resourceID,
				"resource_type":   resourceType,
				"date":            date,
				"topic":           topic,
				"job_type":        "dlq",
				"group":           group,
				"prometheus_host": config.PrometheusHost,
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
			"GetResource", mock.Anything, &entropyv1beta1.GetResourceRequest{Urn: resourceID},
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
		// assertions
		expectedDlqJob := models.DlqJob{
			// from input
			BatchSize:    int64(batchSize),
			ResourceID:   resourceID,
			ResourceType: resourceType,
			Topic:        topic,
			Name: fmt.Sprintf(
				"%s-%s-%s-%s",
				firehoseResource.Name, // firehose title
				"firehose",            // firehose / dagger
				topic,                 //
				date,                  //
			),

			NumThreads: int64(numThreads),
			Date:       date,
			ErrorTypes: errorTypes,

			// firehose resource
			ContainerImage:       config.DlqJobImage,
			DlqGcsCredentialPath: envVars["DLQ_GCS_CREDENTIAL_PATH"],
			EnvVars:              expectedEnvVars,
			Group:                "", //
			KubeCluster:          kubeCluster,
			Namespace:            namespace,
			Project:              firehoseResource.Project,
			PrometheusHost:       config.PrometheusHost,

			// hardcoded
			Replicas: 0,

			// job resource
			Urn:       jobResource.Urn,
			Status:    jobResource.GetState().GetStatus().String(),
			CreatedAt: strfmt.DateTime(jobResource.CreatedAt.AsTime()),
			CreatedBy: jobResource.CreatedBy,
			UpdatedAt: strfmt.DateTime(jobResource.UpdatedAt.AsTime()),
			UpdatedBy: jobResource.UpdatedBy,
		}
		assert.Equal(t, http.StatusOK, response.Code)
		resultJSON := response.Body.Bytes()
		expectedJSON, err := json.Marshal(map[string]interface{}{
			"dlq_job": expectedDlqJob,
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
