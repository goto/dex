package gcs

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/api/iterator"
)

type MockStorageClient struct {
	mock.Mock
}
type MockIterator struct {
	mock.Mock
}

func (m *MockIterator) Next() (*storage.ObjectAttrs, error) {
	args := m.Called()
	if args.Get(0) == nil {
		err, _ := args.Get(1).(error)
		return nil, err
	}
	attributes, _ := args.Get(0).(*storage.ObjectAttrs)
	return attributes, nil
}

var mockIterator = &MockIterator{}

func (*MockStorageClient) Objects(context.Context, string, *storage.Query) WrapperIterator {
	return mockIterator
}

func TestListTopicDates(t *testing.T) {
	mockIterator.On("Next").Return(&storage.ObjectAttrs{Name: "test-topic1/2023-08-26/file1", Size: 123}, nil).Once()
	mockIterator.On("Next").Return(&storage.ObjectAttrs{Name: "test-topic1/2023-08-26/file2", Size: 456}, nil).Once()
	mockIterator.On("Next").Return(&storage.ObjectAttrs{Name: "test-topic1/2023-08-27/file3", Size: 789}, nil).Once()
	mockIterator.On("Next").Return(&storage.ObjectAttrs{Name: "test-topic1/2023-08-27/file4", Size: 101}, nil).Once()
	mockIterator.On("Next").Return(&storage.ObjectAttrs{Name: "test-topic2/2023-08-28/file5", Size: 707}, nil).Once()
	mockIterator.On("Next").Return(&storage.ObjectAttrs{Name: "test-topic2/2023-08-28/file6", Size: 989}, nil).Once()
	mockIterator.On("Next").Return(nil, iterator.Done).Once()
	client := Client{storageClient: &MockStorageClient{}}
	topicDates, err := client.ListTopicDates(BucketInfo{
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
	mockIterator.On("Next").Return(nil, fmt.Errorf("test-error")).Once()
	client := Client{storageClient: &MockStorageClient{}}
	topicDates, err := client.ListTopicDates(BucketInfo{
		BucketName: "test-bucket",
		Prefix:     "prefix",
		Delim:      "",
	})
	assert.Nil(t, topicDates)
	assert.Error(t, err)
	assert.Equal(t, "Bucket(\"test-bucket\").Objects(): test-error", err.Error())
}

func TestErrorForWrongPath(t *testing.T) {
	mockIterator.On("Next").Return(&storage.ObjectAttrs{Name: "test-topic1/31", Size: 123}, nil).Once()
	client := Client{storageClient: &MockStorageClient{}}
	topicDates, err := client.ListTopicDates(BucketInfo{
		BucketName: "test-bucket",
		Prefix:     "prefix",
		Delim:      "",
	})
	assert.Nil(t, topicDates)
	assert.Error(t, err)
	assert.Equal(t, "Object is not in correct path, It should be topic/date/file-name\nPath: test-topic1/31", err.Error())
}
