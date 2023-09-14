package gcs_test

import (
	"fmt"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/api/iterator"

	"github.com/goto/dex/internal/server/gcs"
	"github.com/goto/dex/mocks"
)

func TestListTopicDates(t *testing.T) {
	mt := &mocks.ObjectIterator{}
	mt.On("Next").Return(&storage.ObjectAttrs{Name: "prefix/test-topic1/2023-08-26/file1", Size: 123}, nil).Once()
	mt.On("Next").Return(&storage.ObjectAttrs{Name: "prefix/test-topic1/2023-08-26/file2", Size: 456}, nil).Once()
	mt.On("Next").Return(&storage.ObjectAttrs{Name: "prefix/test-topic1/2023-08-27/file3", Size: 789}, nil).Once()
	mt.On("Next").Return(&storage.ObjectAttrs{Name: "prefix/test-topic1/2023-08-27/file4", Size: 101}, nil).Once()
	mt.On("Next").Return(&storage.ObjectAttrs{Name: "prefix/test-topic2/2023-08-28/file5", Size: 707}, nil).Once()
	mt.On("Next").Return(&storage.ObjectAttrs{Name: "prefix/test-topic2/2023-08-28/file6", Size: 989}, nil).Once()
	mt.On("Next").Return(nil, iterator.Done).Once()
	mc := &mocks.BlobObjectClient{}
	mc.On("Objects", mock.Anything, mock.Anything, mock.Anything).Return(mt).Once()
	client := gcs.Client{StorageClient: mc}
	topicDates, err := client.ListTopicDates(gcs.BucketInfo{
		BucketName: "test-bucket",
		Prefix:     "prefix",
		Delim:      "",
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(topicDates))
	assert.Equal(t, 2, len(topicDates["test-topic1"]))
	assert.Equal(t, 1, len(topicDates["test-topic2"]))
	assert.Equal(t, int64(579), topicDates["test-topic1"]["2023-08-26"])
	assert.Equal(t, int64(890), topicDates["test-topic1"]["2023-08-27"])
	assert.Equal(t, int64(1696), topicDates["test-topic2"]["2023-08-28"])
}

func TestErrorOnListTopic(t *testing.T) {
	mc := &mocks.BlobObjectClient{}
	mt := &mocks.ObjectIterator{}
	mt.On("Next").Return(nil, fmt.Errorf("test-error")).Once()
	mc.On("Objects", mock.Anything, mock.Anything, mock.Anything).Return(mt).Once()
	client := gcs.Client{StorageClient: mc}
	topicDates, err := client.ListTopicDates(gcs.BucketInfo{
		BucketName: "test-bucket",
		Prefix:     "prefix",
		Delim:      "",
	})
	assert.Nil(t, topicDates)
	assert.Error(t, err)
	assert.Equal(t, "Bucket(\"test-bucket\").Objects(): test-error", err.Error())
}
