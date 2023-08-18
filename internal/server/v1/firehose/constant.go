package firehose

import "time"

const (
	logSinkTTL = 24 * time.Hour
)

// Firehose Sink used/modified by dex
const (
	logSinkType      = "LOG"
	bigquerySinkType = "BIGQUERY"
)

// Firehose Configs used/modified by dex
// Refer https://goto.github.io/firehose/advance/generic/
const (
	configSourceKafkaTopic         = "SOURCE_KAFKA_TOPIC"
	configSinkType                 = "SINK_TYPE"
	configStreamName               = "STREAM_NAME"
	configProtoClassName           = "INPUT_SCHEMA_PROTO_CLASS"
	configSourceKafkaBrokers       = "SOURCE_KAFKA_BROKERS"
	configSourceKafkaConsumerGroup = "SOURCE_KAFKA_CONSUMER_GROUP_ID"
	configStencilURL               = "SCHEMA_REGISTRY_STENCIL_URLS"
	configStencilEnable            = "SCHEMA_REGISTRY_STENCIL_ENABLE"

	// BIGQUERY
	configBigqueryTableName   = "SINK_BIGQUERY_TABLE_NAME"
	configBigqueryDatasetName = "SINK_BIGQUERY_DATASET_NAME"
)
