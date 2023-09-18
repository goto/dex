package gcs

import (
	"context"

	"cloud.google.com/go/storage"

	"github.com/goto/dex/generated/models"
)

// BlobStorageClient This is used in service
type BlobStorageClient interface {
	ListDlqMetadata(bucketInfo BucketInfo) ([]models.DlqMetadata, error)
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
