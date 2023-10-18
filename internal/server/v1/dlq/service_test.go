package dlq_test

import (
	"context"
	"fmt"
	"testing"

	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"github.com/go-openapi/strfmt"
	"github.com/goto/dex/entropy"
	"github.com/goto/dex/generated/models"
	"github.com/goto/dex/internal/server/utils"
	"github.com/goto/dex/internal/server/v1/dlq"
	"github.com/goto/dex/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestServiceCreateDLQJob(t *testing.T) {
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

		dlqJob := models.DlqJob{
			BatchSize:  int64(5),
			Date:       "2012-10-30",
			ErrorTypes: "DESERILIAZATION_ERROR",
			// Group: "",
			NumThreads:   2,
			ResourceID:   "test-resource-id",
			ResourceType: "firehose",
			Topic:        "test-create-topic",
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
			"DLQ_BATCH_SIZE":                          fmt.Sprintf("%d", dlqJob.BatchSize),
			"DLQ_NUM_THREADS":                         fmt.Sprintf("%d", dlqJob.NumThreads),
			"DLQ_ERROR_TYPES":                         dlqJob.ErrorTypes,
			"DLQ_INPUT_DATE":                          dlqJob.Date,
			"DLQ_TOPIC_NAME":                          dlqJob.Topic,
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
			Replicas:  0,
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
		entropyCtx := metadata.AppendToOutgoingContext(ctx, "user-id", userEmail)
		newJobResourcePayload := &entropyv1beta1.Resource{
			Urn:  dlqJob.Urn,
			Kind: entropy.ResourceKindJob,
			Name: fmt.Sprintf(
				"%s-%s-%s-%s",
				dlqJob.ResourceID,   // firehose urn
				dlqJob.ResourceType, // firehose / dagger
				dlqJob.Topic,        //
				dlqJob.Date,         //
			),
			Project: firehoseResource.Project,
			Labels: map[string]string{
				"resource_id": dlqJob.ResourceID,
				"type":        dlqJob.ResourceType,
				"date":        dlqJob.Date,
				"topic":       dlqJob.Topic,
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
			"GetResource", ctx, &entropyv1beta1.GetResourceRequest{Urn: dlqJob.ResourceID},
		).Return(&entropyv1beta1.GetResourceResponse{
			Resource: firehoseResource,
		}, nil)
		entropyClient.On("CreateResource", entropyCtx, &entropyv1beta1.CreateResourceRequest{
			Resource: newJobResourcePayload,
		}).Return(&entropyv1beta1.CreateResourceResponse{
			Resource: jobResource,
		}, nil)
		defer entropyClient.AssertExpectations(t)
		service := dlq.NewService(entropyClient, nil, config)

		err = service.CreateDLQJob(ctx, userEmail, &dlqJob)

		// assertions
		expectedDlqJob := models.DlqJob{
			// from input
			BatchSize:    dlqJob.BatchSize,
			ResourceID:   dlqJob.ResourceID,
			ResourceType: dlqJob.ResourceType,
			Topic:        dlqJob.Topic,
			NumThreads:   dlqJob.NumThreads,
			Date:         dlqJob.Date,
			ErrorTypes:   dlqJob.ErrorTypes,

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

		assert.NoError(t, err)
		assert.Equal(t, expectedDlqJob, dlqJob)
	})
}
