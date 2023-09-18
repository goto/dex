// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	gcs "github.com/goto/dex/internal/server/gcs"
	mock "github.com/stretchr/testify/mock"

	models "github.com/goto/dex/generated/models"
)

// BlobStorageClient is an autogenerated mock type for the BlobStorageClient type
type BlobStorageClient struct {
	mock.Mock
}

// ListDlqMetadata provides a mock function with given fields: bucketInfo
func (_m *BlobStorageClient) ListDlqMetadata(bucketInfo gcs.BucketInfo) ([]models.DlqMetadata, error) {
	ret := _m.Called(bucketInfo)

	var r0 []models.DlqMetadata
	var r1 error
	if rf, ok := ret.Get(0).(func(gcs.BucketInfo) ([]models.DlqMetadata, error)); ok {
		return rf(bucketInfo)
	}
	if rf, ok := ret.Get(0).(func(gcs.BucketInfo) []models.DlqMetadata); ok {
		r0 = rf(bucketInfo)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.DlqMetadata)
		}
	}

	if rf, ok := ret.Get(1).(func(gcs.BucketInfo) error); ok {
		r1 = rf(bucketInfo)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewBlobStorageClient creates a new instance of BlobStorageClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBlobStorageClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *BlobStorageClient {
	mock := &BlobStorageClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
