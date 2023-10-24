package dlq

import "errors"

var (
	ErrFirehoseNamespaceNotFound = errors.New("could not find firehose namespace from resource output")
	ErrFirehoseNamespaceInvalid  = errors.New("invalid firehose namespace from resource output")
	ErrFirehoseNotFound          = errors.New("firehose not found")
	ErrEmptyConfigImage          = errors.New("empty dlq job image")
	ErrEmptyConfigPrometheusHost = errors.New("empty prometheus host")
)
