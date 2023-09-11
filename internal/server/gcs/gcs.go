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
)

func NewClient(keyFilePath string) (*Client, error) {
	client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(keyFilePath))
	if err != nil {
		log.Printf("Failed to create GCSClient storageClient: %v\n", err)
		return nil, err
	}
	return &Client{storageClient: SClient{gcsClient: client}}, nil
}

var errWrongPath = errors.New("object is not in correct path")

func wrongPath(name string) error {
	return fmt.Errorf("%w: Path %s", errWrongPath, name)
}

func (client Client) ListTopicDates(bucketInfo BucketInfo) (map[string]map[string]int64, error) {
	bucket := bucketInfo.BucketName
	prefix := bucketInfo.Prefix
	delim := bucketInfo.Delim
	ctx := context.Background()
	// map(topic -> map(Date -> size))
	topicDateMap := make(map[string]map[string]int64)
	ctx, cancel := context.WithTimeout(ctx, time.Second*120)
	defer cancel()
	it := client.storageClient.Objects(ctx, bucket, &storage.Query{
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
		if len(splits) != 3 {
			return nil, wrongPath(attrs.Name)
		}
		topicName := splits[0]
		date := splits[1]
		if topicDateMap[topicName] == nil {
			topicDateMap[topicName] = make(map[string]int64)
		}
		topicDateMap[topicName][date] += attrs.Size
	}
	return topicDateMap, nil
}
