package dlq_test

import (
	"context"
	"testing"

	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/goto/dex/entropy"
	"github.com/goto/dex/generated/models"
	"github.com/goto/dex/internal/server/utils"
	"github.com/goto/dex/internal/server/v1/dlq"
	"github.com/goto/dex/mocks"
)

func TestServiceGetDlqJob(t *testing.T) {
	sampleJobURN := "test-dlq-job-urn"
	config := dlq.DlqJobConfig{
		PrometheusHost: "http://sample-prom-host",
		DlqJobImage:    "test-image",
	}
	t.Run("should return dlq job details", func(t *testing.T) {
		entropyClient := new(mocks.ResourceServiceClient)
		service := dlq.NewService(entropyClient, nil, &config)

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
		service := dlq.NewService(entropyClient, nil, &config)

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
		service := dlq.NewService(entropyClient, nil, &config)

		_, err := service.GetDlqJob(ctx, "test-job-urn")
		assert.ErrorIs(t, err, expectedErr)
	})
}
