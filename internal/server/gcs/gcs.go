package gcs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/goto/dex/generated/models"
)

const clientTimeout = time.Second * 120

func NewClient(keyFilePath string) (BlobObjectClient, error) {
	client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(keyFilePath))
	if err != nil {
		log.Printf("Failed to create GCSClient storageClient: %v\n", err)
		return nil, err
	}
	return &SClient{gcsClient: client}, nil
}

func (client Client) ListDlqMetadata(bucketInfo BucketInfo) ([]models.DlqMetadata, error) {
	bucket := bucketInfo.BucketName
	prefix := bucketInfo.Prefix
	delim := bucketInfo.Delim
	ctx := context.Background()
	topicDateMap := make(map[string]map[string]int64)
	ctx, cancel := context.WithTimeout(ctx, clientTimeout)
	defer cancel()
	it := client.StorageClient.Objects(ctx, bucket, &storage.Query{
		Prefix:    prefix,
		Delimiter: delim,
	})
	for {
		attrs, err := it.Next()
		if errors.Is(iterator.Done, err) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Bucket(%q).Objects(): %w", bucket, err)
		}
		splits := strings.Split(attrs.Name, "/")
		if len(splits) != 4 {
			continue
		}
		// prefix/topic-name/date/object-name
		topicName := splits[1]
		date := splits[2]
		if topicDateMap[topicName] == nil {
			topicDateMap[topicName] = make(map[string]int64)
		}
		topicDateMap[topicName][date] += attrs.Size
	}
	var returnVal []models.DlqMetadata
	for topic, dates := range topicDateMap {
		for date, size := range dates {
			returnVal = append(returnVal, models.DlqMetadata{
				Topic:       topic,
				Date:        date,
				SizeInBytes: size,
			})
		}
	}
	return returnVal, nil
}
