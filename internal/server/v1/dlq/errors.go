package dlq

import "errors"

var (
	ErrFirehoseNamespaceNotFound = errors.New("could not find firehose namespace from resource output")
	ErrFirehoseNamespaceInvalid  = errors.New("invalid firehose namespace from resource output")
	ErrJobNotFound               = errors.New("no job found for this URN")
)
