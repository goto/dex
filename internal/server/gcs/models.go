package gcs

import (
	"context"

	"cloud.google.com/go/storage"
)

// BlobStorageClient This is used in service
type BlobStorageClient interface {
	ListTopicDates(bucketInfo BucketInfo) ([]TopicMetaData, error)
}

type TopicMetaData struct {
	Topic       string `json:"topic"`
	Date        string `json:"date"`
	SizeInBytes int64  `json:"size_in_bytes"`
}

// BlobObjectClient This is used to abstract actual gcs client
type BlobObjectClient interface {
	Objects(ctx context.Context, bucket string, query *storage.Query) ObjectIterator
}

type ObjectIterator interface {
	Next() (*storage.ObjectAttrs, error)
}

type Client struct {
	StorageClient BlobObjectClient
}

type BucketInfo struct {
	BucketName string
	Prefix     string
	Delim      string
}

type SClient struct {
	gcsClient *storage.Client
}

func (c SClient) Objects(ctx context.Context, bucket string, query *storage.Query) ObjectIterator {
	return c.gcsClient.Bucket(bucket).Objects(ctx, query)
}
