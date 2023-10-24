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

// Some of firehose Configs used/modified in more than one place by dex
// Refer https://goto.github.io/firehose/advance/generic/
const (
	configSourceKafkaTopic   = "SOURCE_KAFKA_TOPIC"
	configSourceKafkaBrokers = "SOURCE_KAFKA_BROKERS"
	configSinkType           = "SINK_TYPE"
	configStreamName         = "STREAM_NAME"
	configStencilURL         = "SCHEMA_REGISTRY_STENCIL_URLS"
)

const (
	ConfigDLQBucket          = "DLQ_GCS_BUCKET_NAME"
	ConfigDLQDirectoryPrefix = "DLQ_GCS_DIRECTORY_PREFIX"
)
