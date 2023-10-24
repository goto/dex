package dlq_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/goto/dex/entropy"
	"github.com/goto/dex/generated/models"
	"github.com/goto/dex/internal/server/utils"
	"github.com/goto/dex/internal/server/v1/dlq"
	"github.com/goto/dex/mocks"
	"github.com/goto/dex/pkg/test"
)

func TestServiceListDLQJob(t *testing.T) {
	var (
		project      = "test-project-1"
		resourceID   = "test-resource-id"
		resourceType = "test-resource-type"
		kubeCluster  = "test-kube-cluster"
		date         = "2022-10-21"
		topic        = "test-topic"
		group        = "test-group"
	)

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
				"resource_id":             resourceID,
				"resource_type":           resourceType,
				"date":                    date,
				"topic":                   topic,
				"job_type":                "dlq",
				"group":                   group,
				"prometheus_host":         "test.prom.com",
				"replicas":                "1",
				"batch_size":              "3",
				"num_threads":             "5",
				"error_types":             "SCHEMA_ERROR",
				"dlq_gcs_credential_path": "/etc/test/123",
				"container_image":         "test-image-dlq:1.0.0",
				"namespace":               "dlq-namespace",
			},
			Spec: &entropyv1beta1.ResourceSpec{
				Dependencies: []*entropyv1beta1.ResourceDependency{
					{
						Key:   "kube_cluster",
						Value: kubeCluster, // from firehose configs.kube_cluster
					},
				},
			},
			CreatedBy: "user@test.com",
			UpdatedBy: "user@test.com",
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
				"resource_id":             resourceID,
				"resource_type":           resourceType,
				"date":                    date,
				"topic":                   topic,
				"job_type":                "dlq",
				"group":                   group,
				"prometheus_host":         "test2.prom.com",
				"replicas":                "12",
				"batch_size":              "4",
				"num_threads":             "12",
				"error_types":             "TEST_ERROR",
				"dlq_gcs_credential_path": "/etc/test/312",
				"container_image":         "test-image-dlq:2.0.0",
				"namespace":               "dlq-namespace-2",
			},
			Spec: &entropyv1beta1.ResourceSpec{
				Dependencies: []*entropyv1beta1.ResourceDependency{
					{
						Key:   "kube_cluster",
						Value: kubeCluster, // from firehose configs.kube_cluster
					},
				},
			},
			CreatedBy: "user@test.com",
			UpdatedBy: "user@test.com",
		},
	}

	expectedDlqJob := []models.DlqJob{
		{
			// from input
			BatchSize:    3,
			ResourceID:   "test-resource-id",
			ResourceType: "test-resource-type",
			Topic:        "test-topic",
			Name:         "test1-firehose-test-topic-2022-10-21",
			NumThreads:   5,
			Date:         "2022-10-21",
			ErrorTypes:   "SCHEMA_ERROR",

			// firehose resource
			ContainerImage:       "test-image-dlq:1.0.0",
			DlqGcsCredentialPath: "/etc/test/123",
			Group:                "test-group", //
			KubeCluster:          kubeCluster,
			Namespace:            "dlq-namespace",
			Project:              project,
			PrometheusHost:       "test.prom.com",

			// hardcoded
			Replicas: 1,

			// job resource
			Urn:       "test-urn-1",
			Status:    "STATUS_UNSPECIFIED",
			CreatedAt: strfmt.DateTime(dummyEntropyResources[0].CreatedAt.AsTime()),
			CreatedBy: "user@test.com",
			UpdatedAt: strfmt.DateTime(dummyEntropyResources[0].UpdatedAt.AsTime()),
			UpdatedBy: "user@test.com",
		},
		{
			// from input
			BatchSize:    4,
			ResourceID:   "test-resource-id",
			ResourceType: "test-resource-type",
			Topic:        "test-topic",
			Name:         "test2-firehose-test-topic-2022-10-21",
			NumThreads:   12,
			Date:         "2022-10-21",
			ErrorTypes:   "TEST_ERROR",

			// firehose resource
			ContainerImage:       "test-image-dlq:2.0.0",
			DlqGcsCredentialPath: "/etc/test/312",
			Group:                "test-group", //
			KubeCluster:          kubeCluster,
			Namespace:            "dlq-namespace-2",
			Project:              project,
			PrometheusHost:       "test2.prom.com",

			// hardcoded
			Replicas: 12,

			// job resource
			Urn:       "test-urn-2",
			Status:    "STATUS_UNSPECIFIED",
			CreatedAt: strfmt.DateTime(dummyEntropyResources[0].CreatedAt.AsTime()),
			CreatedBy: "user@test.com",
			UpdatedAt: strfmt.DateTime(dummyEntropyResources[0].UpdatedAt.AsTime()),
			UpdatedBy: "user@test.com",
		},
	}

	t.Run("should return dlqjob list", func(t *testing.T) {
		ctx := context.TODO()

		expectedRPCResp := &entropyv1beta1.ListResourcesResponse{
			Resources: dummyEntropyResources,
		}

		labelFilter := map[string]string{
			"resource_id":   resourceID,
			"resource_type": resourceType,
			"date":          date,
		}
		entropyClient := new(mocks.ResourceServiceClient)
		entropyClient.On(
			"ListResources", ctx, &entropyv1beta1.ListResourcesRequest{
				Kind: entropy.ResourceKindJob, Labels: labelFilter,
			},
		).Return(expectedRPCResp, nil)
		defer entropyClient.AssertExpectations(t)

		service := dlq.NewService(entropyClient, nil, dlq.DlqJobConfig{})

		dlqJob, err := service.ListDlqJob(ctx, labelFilter)
		assert.NoError(t, err)
		assert.Equal(t, expectedDlqJob, dlqJob)
	})
}

// nolint
func skipTestServiceCreateDLQJob(t *testing.T) {
	t.Run("should return ErrFirehoseNotFound if resource cannot be found in entropy", func(t *testing.T) {
		// inputs
		ctx := context.TODO()
		dlqJob := models.DlqJob{
			ResourceID:   "test-resource-id",
			ResourceType: "firehose",
		}
		expectedErr := status.Error(codes.NotFound, "Not found")

		// set conditions
		entropyClient := new(mocks.ResourceServiceClient)
		entropyClient.On(
			"GetResource", ctx, &entropyv1beta1.GetResourceRequest{Urn: dlqJob.ResourceID},
		).Return(nil, expectedErr)
		defer entropyClient.AssertExpectations(t)
		service := dlq.NewService(entropyClient, nil, dlq.DlqJobConfig{})

		_, err := service.CreateDLQJob(ctx, "", dlqJob)
		assert.ErrorIs(t, err, dlq.ErrFirehoseNotFound)
	})

	t.Run("should return error when there is an error getting firehose in entropy", func(t *testing.T) {
		// inputs
		ctx := context.TODO()
		dlqJob := models.DlqJob{
			ResourceID:   "test-resource-id",
			ResourceType: "firehose",
		}
		expectedErr := status.Error(codes.Internal, "Any Error")

		// set conditions
		entropyClient := new(mocks.ResourceServiceClient)
		entropyClient.On(
			"GetResource", ctx, &entropyv1beta1.GetResourceRequest{Urn: dlqJob.ResourceID},
		).Return(nil, expectedErr)
		defer entropyClient.AssertExpectations(t)
		service := dlq.NewService(entropyClient, nil, dlq.DlqJobConfig{})

		_, err := service.CreateDLQJob(ctx, "", dlqJob)
		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("should create a entropy resource with job kind", func(t *testing.T) {
		// inputs
		ctx := context.TODO()
		namespace := "test-namespace"
		kubeCluster := "test-kube-cluster"
		userEmail := "test@example.com"
		config := dlq.DlqJobConfig{
			PrometheusHost: "http://sample-prom-host",
			DlqJobImage:    "test-image",
		}

		dlqJob := models.DlqJob{
			BatchSize:    int64(5),
			Date:         "2012-10-30",
			ErrorTypes:   "DESERILIAZATION_ERROR",
			Group:        "",
			NumThreads:   2,
			ResourceID:   "test-resource-id",
			ResourceType: "firehose",
			Topic:        "test-create-topic",
		}

		// setup firehose resource
		firehoseEnvVars := map[string]string{
			"SINK_TYPE":                               "bigquery",
			"DLQ_FINISHED_STATUS_FILE":                "/shared/job-finished",
			"DLQ_GCS_BUCKET_NAME":                     "g-pilotdata-gl-dlq",
			"DLQ_GCS_CREDENTIAL_PATH":                 "/etc/secret/gcp/token",
			"DLQ_GCS_GOOGLE_CLOUD_PROJECT_ID":         "pilotdata-integration",
			"JAVA_TOOL_OPTIONS":                       "-javaagent:jolokia-jvm-agent.jar=port=8778,host=localhost",
			"_JAVA_OPTIONS":                           "-Xmx1800m -Xms1800m",
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
		firehoseConfig, err := utils.GoValToProtoStruct(entropy.FirehoseConfig{
			EnvVariables: firehoseEnvVars,
		})
		require.NoError(t, err)
		firehoseResource := &entropyv1beta1.Resource{
			Name: "test-firehose-name",
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
				Output: structpb.NewStructValue(test.NewStruct(t, map[string]interface{}{
					"namespace": namespace,
				})),
			},
		}

		updatedEnvVars := map[string]string{
			"DLQ_BATCH_SIZE":                          "5",
			"DLQ_NUM_THREADS":                         "2",
			"DLQ_ERROR_TYPES":                         dlqJob.ErrorTypes,
			"DLQ_INPUT_DATE":                          dlqJob.Date,
			"DLQ_TOPIC_NAME":                          dlqJob.Topic,
			"METRIC_STATSD_TAGS":                      "a=b", // TBA
			"DLQ_PREFIX_DIR":                          "test-firehose",
			"DLQ_FINISHED_STATUS_FILE":                "/shared/job-finished",
			"SINK_TYPE":                               firehoseEnvVars["SINK_TYPE"],
			"DLQ_GCS_BUCKET_NAME":                     firehoseEnvVars["DLQ_GCS_BUCKET_NAME"],
			"DLQ_GCS_CREDENTIAL_PATH":                 firehoseEnvVars["DLQ_GCS_CREDENTIAL_PATH"],
			"DLQ_GCS_GOOGLE_CLOUD_PROJECT_ID":         firehoseEnvVars["DLQ_GCS_GOOGLE_CLOUD_PROJECT_ID"],
			"JAVA_TOOL_OPTIONS":                       firehoseEnvVars["JAVA_TOOL_OPTIONS"],
			"_JAVA_OPTIONS":                           firehoseEnvVars["_JAVA_OPTIONS"],
			"INPUT_SCHEMA_PROTO_CLASS":                firehoseEnvVars["INPUT_SCHEMA_PROTO_CLASS"],
			"SCHEMA_REGISTRY_STENCIL_ENABLE":          firehoseEnvVars["SCHEMA_REGISTRY_STENCIL_ENABLE"],
			"SCHEMA_REGISTRY_STENCIL_URLS":            firehoseEnvVars["SCHEMA_REGISTRY_STENCIL_URLS"],
			"SINK_BIGQUERY_ADD_METADATA_ENABLED":      firehoseEnvVars["SINK_BIGQUERY_ADD_METADATA_ENABLED"],
			"SINK_BIGQUERY_CLIENT_CONNECT_TIMEOUT_MS": firehoseEnvVars["SINK_BIGQUERY_CLIENT_CONNECT_TIMEOUT_MS"],
			"SINK_BIGQUERY_CLIENT_READ_TIMEOUT_MS":    firehoseEnvVars["SINK_BIGQUERY_CLIENT_READ_TIMEOUT_MS"],
			"SINK_BIGQUERY_CREDENTIAL_PATH":           firehoseEnvVars["SINK_BIGQUERY_CREDENTIAL_PATH"],
			"SINK_BIGQUERY_DATASET_LABELS":            firehoseEnvVars["SINK_BIGQUERY_DATASET_LABELS"],
			"SINK_BIGQUERY_DATASET_LOCATION":          firehoseEnvVars["SINK_BIGQUERY_DATASET_LOCATION"],
			"SINK_BIGQUERY_DATASET_NAME":              firehoseEnvVars["SINK_BIGQUERY_DATASET_NAME"],
			"SINK_BIGQUERY_GOOGLE_CLOUD_PROJECT_ID":   firehoseEnvVars["SINK_BIGQUERY_GOOGLE_CLOUD_PROJECT_ID"],
			"SINK_BIGQUERY_ROW_INSERT_ID_ENABLE":      firehoseEnvVars["SINK_BIGQUERY_ROW_INSERT_ID_ENABLE"],
			"SINK_BIGQUERY_STORAGE_API_ENABLE":        firehoseEnvVars["SINK_BIGQUERY_STORAGE_API_ENABLE"],
			"SINK_BIGQUERY_TABLE_LABELS":              firehoseEnvVars["SINK_BIGQUERY_TABLE_LABELS"],
			"SINK_BIGQUERY_TABLE_NAME":                firehoseEnvVars["SINK_BIGQUERY_TABLE_NAME"],
			"SINK_BIGQUERY_TABLE_PARTITION_EXPIRY_MS": firehoseEnvVars["SINK_BIGQUERY_TABLE_PARTITION_EXPIRY_MS"],
			"SINK_BIGQUERY_TABLE_PARTITION_KEY":       firehoseEnvVars["SINK_BIGQUERY_TABLE_PARTITION_KEY"],
			"SINK_BIGQUERY_TABLE_PARTITIONING_ENABLE": firehoseEnvVars["SINK_BIGQUERY_TABLE_PARTITIONING_ENABLE"],
		}

		// setup expected request to entropy
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
							Mount: updatedEnvVars["DLQ_GCS_CREDENTIAL_PATH"],
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
					EnvVariables: updatedEnvVars,
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
						"TEAM":            dlqJob.Group,
						"TOPIC":           dlqJob.Topic,
						"environment":     "production", // TBA
						"organization":    "de",         // TBA
						"projectID":       dlqJob.Project,
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
				"firehose": dlqJob.ResourceID,
				"topic":    dlqJob.Topic,
				"date":     dlqJob.Date,
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
		expectedJobRequestToEntropy := &entropyv1beta1.Resource{
			Kind:    entropy.ResourceKindJob,
			Name:    "test-firehose-name-firehose-test-create-topic-2012-10-30",
			Project: firehoseResource.Project,
			Labels: map[string]string{
				"resource_id":             dlqJob.ResourceID,
				"resource_type":           dlqJob.ResourceType,
				"date":                    dlqJob.Date,
				"topic":                   dlqJob.Topic,
				"job_type":                "dlq",
				"group":                   dlqJob.Group,
				"prometheus_host":         config.PrometheusHost,
				"batch_size":              "5",
				"num_threads":             "2",
				"error_types":             "DESERILIAZATION_ERROR",
				"dlq_gcs_credential_path": updatedEnvVars["DLQ_GCS_CREDENTIAL_PATH"],
				"container_image":         config.DlqJobImage,
				"namespace":               namespace,
				"replicas":                "1",
			},
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
		entropyCtx := metadata.AppendToOutgoingContext(ctx, "user-id", userEmail)

		entropyClient := new(mocks.ResourceServiceClient)
		entropyClient.On(
			"GetResource", ctx, &entropyv1beta1.GetResourceRequest{Urn: dlqJob.ResourceID},
		).Return(&entropyv1beta1.GetResourceResponse{
			Resource: firehoseResource,
		}, nil)
		entropyClient.On("CreateResource", entropyCtx, &entropyv1beta1.CreateResourceRequest{
			Resource: expectedJobRequestToEntropy,
		}).Return(&entropyv1beta1.CreateResourceResponse{
			Resource: &entropyv1beta1.Resource{
				Urn:       "test-urn",
				CreatedBy: "test-created-by",
				UpdatedBy: "test-updated-by",
				CreatedAt: timestamppb.New(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt: timestamppb.New(time.Date(2022, 2, 2, 1, 1, 1, 1, time.UTC)),
				Kind:      expectedJobRequestToEntropy.Kind,
				Name:      expectedJobRequestToEntropy.Name,
				Project:   expectedJobRequestToEntropy.Project,
				Labels:    expectedJobRequestToEntropy.GetLabels(),
				Spec:      expectedJobRequestToEntropy.GetSpec(),
			},
		}, nil)
		defer entropyClient.AssertExpectations(t)
		service := dlq.NewService(entropyClient, nil, config)

		result, err := service.CreateDLQJob(ctx, userEmail, dlqJob)
		assert.NoError(t, err)

		// assertions
		expectedDlqJob := models.DlqJob{
			Replicas:             1,
			Urn:                  "test-urn",
			Status:               "STATUS_UNSPECIFIED",
			Name:                 expectedJobRequestToEntropy.Name,
			NumThreads:           dlqJob.NumThreads,
			Date:                 dlqJob.Date,
			ErrorTypes:           dlqJob.ErrorTypes,
			BatchSize:            dlqJob.BatchSize,
			ResourceID:           dlqJob.ResourceID,
			ResourceType:         dlqJob.ResourceType,
			Topic:                dlqJob.Topic,
			ContainerImage:       config.DlqJobImage,
			DlqGcsCredentialPath: updatedEnvVars["DLQ_GCS_CREDENTIAL_PATH"],
			EnvVars:              updatedEnvVars,
			KubeCluster:          kubeCluster,
			Namespace:            namespace,
			Project:              firehoseResource.Project,
			PrometheusHost:       config.PrometheusHost,
			CreatedAt:            strfmt.DateTime(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			UpdatedAt:            strfmt.DateTime(time.Date(2022, 2, 2, 1, 1, 1, 1, time.UTC)),
			CreatedBy:            "test-created-by",
			UpdatedBy:            "test-updated-by",
		}
		assert.Equal(t, expectedDlqJob, result)
	})
}

func TestServiceGetDlqJob(t *testing.T) {
	sampleJobURN := "test-dlq-job-urn"
	config := dlq.DlqJobConfig{
		PrometheusHost: "http://sample-prom-host",
		DlqJobImage:    "test-image",
	}
	t.Run("should return dlq job details", func(t *testing.T) {
		entropyClient := new(mocks.ResourceServiceClient)
		service := dlq.NewService(entropyClient, nil, config)

		envVars := map[string]string{
			"DLQ_BATCH_SIZE":                          "34",
			"DLQ_ERROR_TYPES":                         "DEFAULT_ERROR",
			"DLQ_FINISHED_STATUS_FILE":                "/shared/job-finished",
			"DLQ_GCS_BUCKET_NAME":                     "test_bucket_name",
			"DLQ_GCS_CREDENTIAL_PATH":                 "/path/to/token",
			"DLQ_GCS_GOOGLE_CLOUD_PROJECT_ID":         "pilot-project",
			"DLQ_INPUT_DATE":                          "2023-04-10",
			"DLQ_NUM_THREADS":                         "10",
			"DLQ_PREFIX_DIR":                          "test-firehose",
			"DLQ_TOPIC_NAME":                          "test-topic",
			"INPUT_SCHEMA_PROTO_CLASS":                "test-proto",
			"JAVA_TOOL_OPTIONS":                       "test-tool-options",
			"METRIC_STATSD_TAGS":                      "a=b",
			"SCHEMA_REGISTRY_STENCIL_ENABLE":          "true",
			"SCHEMA_REGISTRY_STENCIL_URLS":            "http://random-test-url",
			"SINK_BIGQUERY_ADD_METADATA_ENABLED":      "true",
			"SINK_BIGQUERY_CLIENT_CONNECT_TIMEOUT_MS": "-1",
			"SINK_BIGQUERY_CLIENT_READ_TIMEOUT_MS":    "-1",
			"SINK_BIGQUERY_CREDENTIAL_PATH":           "/path/to/token",
			"SINK_BIGQUERY_DATASET_LABELS":            "shangchi=legend,lord=voldemort",
			"SINK_BIGQUERY_DATASET_LOCATION":          "US",
			"SINK_BIGQUERY_DATASET_NAME":              "bq_test",
			"SINK_BIGQUERY_GOOGLE_CLOUD_PROJECT_ID":   "pilot-project",
			"SINK_BIGQUERY_ROW_INSERT_ID_ENABLE":      "false",
			"SINK_BIGQUERY_STORAGE_API_ENABLE":        "true",
			"SINK_BIGQUERY_TABLE_LABELS":              "hello=world,john=doe",
			"SINK_BIGQUERY_TABLE_NAME":                "bq_dlq_test1",
			"SINK_BIGQUERY_TABLE_PARTITIONING_ENABLE": "true",
			"SINK_BIGQUERY_TABLE_PARTITION_EXPIRY_MS": "2629800000",
			"SINK_BIGQUERY_TABLE_PARTITION_KEY":       "event_timestamp",
			"SINK_TYPE":                               "bigquery",
			"_JAVA_OPTIONS":                           "-Xmx1800m -Xms1800m",
		}

		jobContainers := []entropy.JobContainer{
			{
				Name:            "dlq-job",
				Image:           "test/image/path",
				ImagePullPolicy: "Always",
				Limits: entropy.UsageSpec{
					CPU:    "1",
					Memory: "2000Mi",
				},
				Requests: entropy.UsageSpec{
					CPU:    "1",
					Memory: "2000Mi",
				},
				EnvVariables: envVars,
				SecretsVolumes: []entropy.JobSecret{
					{
						Name:  "firehose-bigquery-sink-credential",
						Mount: "/etc/secret/gcp",
					},
				},
			},
			{
				Args: []string{
					"-c",
					"test args random args",
				},
				Command: []string{
					"/bin/bash",
				},
				ConfigMapsVolumes: []entropy.JobConfigMap{
					{
						Name:  "dlq-processor-telegraf",
						Mount: "etc/telegraf",
					},
				},
				EnvVariables: map[string]string{
					"APP_NAME":        "testing",
					"DEPLOYMENT_NAME": "deployment-name",
					"PROMETHEUS_URL":  "test/prometheus/url",
					"TEAM":            "test-team",
					"TOPIC":           "test-topic",
					"environment":     "production",
					"organization":    "test-team", "projectID": "projectID",
				},
				Image: "telegraf:1.18.0-alpine",
				Limits: entropy.UsageSpec{
					CPU:    "100m",
					Memory: "300Mi",
				},
				Name: "telegraf",
				Requests: entropy.UsageSpec{
					CPU:    "100m",
					Memory: "300Mi",
				},
			},
		}

		jobConfig, err := utils.GoValToProtoStruct(entropy.JobConfig{
			Replicas:   1,
			Namespace:  "default",
			Containers: jobContainers,
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

		output, err := structpb.NewValue(map[string]any{
			"jobName":   "test-dlq-job",
			"namespace": "default",
			"pods":      nil,
		})
		require.NoError(t, err)

		resourceResponse := &entropyv1beta1.GetResourceResponse{
			Resource: &entropyv1beta1.Resource{
				Urn:     "test-dlq-job-urn",
				Kind:    "job",
				Name:    "gofood-test-dlq-4",
				Project: "test-project",
				Labels: map[string]string{
					"team":            "test-team",
					"prometheus_host": "http://sample-prom-host",
					"topic":           "test-topic",
					"resource_type":   "dlq",
					"resource_id":     "test_resource_id",
					"date":            "2023-10-22",
				},
				CreatedAt: &timestamppb.Timestamp{Seconds: 1697176965, Nanos: 496343000},
				UpdatedAt: &timestamppb.Timestamp{Seconds: 1697215353, Nanos: 779372000},
				CreatedBy: "test@user.com",
				UpdatedBy: "test@user.com",
				Spec: &entropyv1beta1.ResourceSpec{
					Configs: jobConfig,
					Dependencies: []*entropyv1beta1.ResourceDependency{
						{
							Key:   "kube_cluster",
							Value: "test-kube-cluster",
						},
					},
				},
				State: &entropyv1beta1.ResourceState{
					Status: entropyv1beta1.ResourceState_STATUS_COMPLETED,
					Output: output,
				},
			},
		}

		expectedDlqJob := &models.DlqJob{
			BatchSize:            34,
			ContainerImage:       "test/image/path",
			CreatedAt:            strfmt.DateTime(resourceResponse.Resource.CreatedAt.AsTime()),
			UpdatedAt:            strfmt.DateTime(resourceResponse.Resource.UpdatedAt.AsTime()),
			CreatedBy:            "test@user.com",
			DlqGcsCredentialPath: "/path/to/token",
			EnvVars:              envVars,
			ErrorTypes:           "DEFAULT_ERROR",
			KubeCluster:          "test-kube-cluster",
			Name:                 "gofood-test-dlq-4",
			Namespace:            "default",
			NumThreads:           10,
			Project:              "test-project",
			PrometheusHost:       "http://sample-prom-host",
			Replicas:             1,
			ResourceType:         "dlq",
			Status:               "\x04",
			UpdatedBy:            "test@user.com",
			Urn:                  "test-dlq-job-urn",
			ResourceID:           "test_resource_id",
			Date:                 "2023-10-22",
			Topic:                "test-topic",
		}

		entropyClient.On("GetResource", context.TODO(), &entropyv1beta1.GetResourceRequest{Urn: sampleJobURN}).Return(&entropyv1beta1.GetResourceResponse{
			Resource: resourceResponse.GetResource(),
		}, nil)

		dlqJob, err := service.GetDlqJob(context.TODO(), sampleJobURN)
		require.NoError(t, err)
		assert.Equal(t, expectedDlqJob, dlqJob)
	})

	t.Run("should return ErrJobNotFound if resource cannot be found in entropy", func(t *testing.T) {
		ctx := context.TODO()
		dlqJob := models.DlqJob{
			Urn:          "test-job-urn",
			ResourceType: "dlq",
		}
		config := dlq.DlqJobConfig{
			PrometheusHost: "http://sample-prom-host",
			DlqJobImage:    "test-image",
		}
		expectedErr := status.Error(codes.NotFound, "Not found")

		entropyClient := new(mocks.ResourceServiceClient)
		entropyClient.On(
			"GetResource", ctx, &entropyv1beta1.GetResourceRequest{Urn: dlqJob.Urn},
		).Return(nil, expectedErr)
		defer entropyClient.AssertExpectations(t)
		service := dlq.NewService(entropyClient, nil, config)

		_, err := service.GetDlqJob(ctx, "test-job-urn")
		assert.ErrorIs(t, err, dlq.ErrJobNotFound)
	})

	t.Run("should return error when there is an error getting job in entropy", func(t *testing.T) {
		// inputs
		ctx := context.TODO()
		dlqJob := models.DlqJob{
			Urn:          "test-job-urn",
			ResourceType: "dlq",
		}
		config := dlq.DlqJobConfig{
			PrometheusHost: "http://sample-prom-host",
			DlqJobImage:    "test-image",
		}
		expectedErr := status.Error(codes.Internal, "Any Error")

		entropyClient := new(mocks.ResourceServiceClient)
		entropyClient.On(
			"GetResource", ctx, &entropyv1beta1.GetResourceRequest{Urn: dlqJob.Urn},
		).Return(nil, expectedErr)
		defer entropyClient.AssertExpectations(t)
		service := dlq.NewService(entropyClient, nil, config)

		_, err := service.GetDlqJob(ctx, "test-job-urn")
		assert.ErrorIs(t, err, expectedErr)
	})
}
