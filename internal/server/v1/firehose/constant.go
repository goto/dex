package firehose

import "time"

const (
	logSinkTTL = 24 * time.Hour
	trueString = "true"
)

// Firehose Sink used/modified by dex
const (
	logSinkType      = "LOG"
	bigquerySinkType = "BIGQUERY"
	blobSinkType     = "BLOB"
)

// Firehose Configs used/modified by dex
// Refer https://goto.github.io/firehose/advance/generic/
const (
	configSourceKafkaTopic           = "SOURCE_KAFKA_TOPIC"
	configSourceKafkaBrokers         = "SOURCE_KAFKA_BROKERS"
	configSourceKafkaConsumerGroup   = "SOURCE_KAFKA_CONSUMER_GROUP_ID"
	configSourceKafkaConsumerMode    = "SOURCE_KAFKA_CONSUMER_MODE"
	configSourceKafkaPollTimeoutMS   = "SOURCE_KAFKA_POLL_TIMEOUT_MS"
	configSinkType                   = "SINK_TYPE"
	configStreamName                 = "STREAM_NAME"
	configProtoClassName             = "INPUT_SCHEMA_PROTO_CLASS"
	configStencilURL                 = "SCHEMA_REGISTRY_STENCIL_URLS"
	configStencilEnable              = "SCHEMA_REGISTRY_STENCIL_ENABLE"
	configJavaOptions                = "_JAVA_OPTIONS"
	configSinkPoolNumThreads         = "SINK_POOL_NUM_THREADS"
	configSinkPoolQueuePollTimeoutMS = "SINK_POOL_QUEUE_POLL_TIMEOUT_MS"

	// DLQ
	configDLQRetryFailAfterMaxAttempEnable = "DLQ_RETRY_FAIL_AFTER_MAX_ATTEMPT_ENABLE"
	configDLQRetryMaxAttempts              = "DLQ_RETRY_MAX_ATTEMPTS"
	configDLQSinkEnable                    = "DLQ_SINK_ENABLE"
	configDLQWriterType                    = "DLQ_WRITER_TYPE"
	configErrorTypesForDLQ                 = "ERROR_TYPES_FOR_DLQ"

	// Retry
	configRetryExponentialBackoffInitialMS = "RETRY_EXPONENTIAL_BACKOFF_INITIAL_MS"
	configErrorTypesForRetry               = "ERROR_TYPES_FOR_RETRY"
	configRetryMaxAttempts                 = "RETRY_MAX_ATTEMPTS"

	// BIGQUERY
	configBigqueryTableName   = "SINK_BIGQUERY_TABLE_NAME"
	configBigqueryDatasetName = "SINK_BIGQUERY_DATASET_NAME"
)
