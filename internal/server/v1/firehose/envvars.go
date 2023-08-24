package firehose

import (
	"context"
	"fmt"
	"strings"

	"github.com/goto/dex/generated/models"
	"github.com/goto/dex/odin"
)

func (api *firehoseAPI) buildEnvVars(ctx context.Context, firehose *models.Firehose, userID string, fetchStencilURL bool) error {
	streamURN := buildStreamURN(*firehose.Configs.StreamName, firehose.Project)

	if firehose.Configs.EnvVars[configSourceKafkaBrokers] == "" {
		sourceKafkaBroker, err := odin.GetOdinStream(ctx, api.OdinAddr, streamURN)
		if err != nil {
			return fmt.Errorf("error getting odin stream: %w", err)
		}
		firehose.Configs.EnvVars[configSourceKafkaBrokers] = sourceKafkaBroker
	}

	firehose.Configs.EnvVars[configStencilEnable] = trueString
	if fetchStencilURL && firehose.Configs.EnvVars[configStencilURL] == "" {
		stencilUrls, err := api.getStencilURLs(
			ctx,
			userID,
			firehose.Configs.EnvVars[configSourceKafkaTopic],
			streamURN,
			firehose.Project,
			firehose.Configs.EnvVars[configProtoClassName],
		)
		if err != nil {
			return fmt.Errorf("error getting stencil url: %w", err)
		}

		firehose.Configs.EnvVars[configStencilURL] = stencilUrls
	}

	sinkType := firehose.Configs.EnvVars[configSinkType]
	firehose.Configs.EnvVars = buildEnvVarsBySink(sinkType, firehose.Configs.EnvVars, firehose.Configs)

	return nil
}

func buildEnvVarsBySink(sinkType string, envVars map[string]string, cfg *models.FirehoseConfig) map[string]string {
	// BQ or GCS sink
	if sinkType == bigquerySinkType || sinkType == blobSinkType {
		envVars[configJavaOptions] = "-Xmx1550m -Xms1550m"
		envVars[configDLQRetryFailAfterMaxAttempEnable] = trueString
		envVars[configDLQRetryMaxAttempts] = "5"
		envVars[configDLQSinkEnable] = trueString
		envVars[configDLQWriterType] = "BLOB_STORAGE"
		envVars[configErrorTypesForDLQ] = "DESERIALIZATION_ERROR,INVALID_MESSAGE_ERROR,UNKNOWN_FIELDS_ERROR,SINK_4XX_ERROR,SINK_5XX_ERROR,SINK_UNKNOWN_ERROR,DEFAULT_ERROR"
		envVars[configErrorTypesForRetry] = "SINK_5XX_ERROR,SINK_UNKNOWN_ERROR,DEFAULT_ERROR"
		envVars[configRetryExponentialBackoffInitialMS] = "100"
		envVars[configRetryMaxAttempts] = "10"
		envVars[configSourceKafkaPollTimeoutMS] = "60000"
	}

	// BQ only
	if sinkType == bigquerySinkType {
		envVars[configSourceKafkaConsumerMode] = "async"
		envVars[configSinkPoolNumThreads] = "8"
		envVars[configSinkPoolQueuePollTimeoutMS] = "100"
		envVars[configBigqueryTablePartitioningEnable] = trueString
		envVars[configBigqueryRowInsertIDEnable] = trueString
		envVars[configBigqueryClientReadTimeoutMS] = "-1"
		envVars[configBigqueryClientConnectTimeoutMS] = "-1"
		envVars[configBigqueryMetadataNamespace] = "__kafka_metadata"
		defaultIfEmpty(envVars, configBigqueryTableName, func() string {
			t := envVars[configSourceKafkaTopic]
			t = strings.ReplaceAll(t, ".", "_")
			t = strings.ReplaceAll(t, "-", "_")
			return t
		})
		defaultIfEmpty(envVars, configBigqueryDatasetName, func() string {
			return *cfg.StreamName
		})
	}

	return envVars
}

func defaultIfEmpty(vars map[string]string, key string, defaultVal func() string) {
	val, exists := vars[key]
	if !exists || val == "" {
		vars[key] = defaultVal()
	}
}
