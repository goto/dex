package gcs

import (
	"context"

	"cloud.google.com/go/storage"
)

// StorageClient This is used in service
type StorageClient interface {
	ListTopicDates(bucketInfo BucketInfo) (map[string]map[string]int64, error)
}

// WrapperClient This is used to abstract actual gcs client
type WrapperClient interface {
	Objects(ctx context.Context, bucket string, query *storage.Query) WrapperIterator
}

type WrapperIterator interface {
	Next() (*storage.ObjectAttrs, error)
}

type Client struct {
	storageClient WrapperClient
}

type BucketInfo struct {
	BucketName string
	Prefix     string
	Delim      string
}

type SClient struct {
	gcsClient *storage.Client
}

func (c SClient) Objects(ctx context.Context, bucket string, query *storage.Query) WrapperIterator {
	return c.gcsClient.Bucket(bucket).Objects(ctx, query)
}
