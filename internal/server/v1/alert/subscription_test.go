package alert_test

import (
	"context"
	"testing"

	sirenv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/siren/v1beta1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goto/dex/internal/server/v1/alert"
	"github.com/goto/dex/mocks"
)

func TestSubscriptionServiceFindSubscription(t *testing.T) {
	ctx := context.TODO()
	subscriptionID := 105

	t.Run("should return subscription on success", func(t *testing.T) {
		subscription := &sirenv1beta1.Subscription{
			Id:        uint64(subscriptionID),
			Urn:       "sample-urn",
			Namespace: 1,
			Receivers: []*sirenv1beta1.ReceiverMetadata{
				{Id: 30},
			},
		}

		client := new(mocks.SirenServiceClient)
		client.On("GetSubscription", ctx, &sirenv1beta1.GetSubscriptionRequest{Id: subscription.Id}).
			Return(&sirenv1beta1.GetSubscriptionResponse{
				Subscription: subscription,
			}, nil)
		defer client.AssertExpectations(t)

		service := alert.NewSubscriptionService(client)
		result, err := service.FindSubscription(ctx, subscriptionID)
		assert.NoError(t, err)
		assert.Equal(t, subscription, result)
	})

	t.Run("should return not found error if optimus return NotFound code", func(t *testing.T) {
		grpcError := status.Error(codes.NotFound, "Not Found")

		client := new(mocks.SirenServiceClient)
		client.On("GetSubscription", ctx, &sirenv1beta1.GetSubscriptionRequest{Id: uint64(subscriptionID)}).
			Return(nil, grpcError)
		defer client.AssertExpectations(t)

		service := alert.NewSubscriptionService(client)
		_, err := service.FindSubscription(ctx, subscriptionID)
		assert.ErrorIs(t, err, alert.ErrSubscriptionNotFound)
	})

	t.Run("should return if client return error", func(t *testing.T) {
		expectedError := status.Error(codes.Internal, "Internal")

		client := new(mocks.SirenServiceClient)
		client.On("GetSubscription", ctx, &sirenv1beta1.GetSubscriptionRequest{Id: uint64(subscriptionID)}).
			Return(nil, expectedError)
		defer client.AssertExpectations(t)

		service := alert.NewSubscriptionService(client)
		_, err := service.FindSubscription(ctx, subscriptionID)
		assert.ErrorIs(t, err, expectedError)
	})
}

func TestSubscriptionServiceGetSubscriptions(t *testing.T) {
	ctx := context.TODO()

	t.Run("should return subscription on success", func(t *testing.T) {
		groupID := "19293012i31"
		resourceID := "sample-resource-id-or-urn"
		resourceType := "firehose"

		subscriptions := []*sirenv1beta1.Subscription{
			{
				Id:        1,
				Urn:       "sample-urn-1",
				Namespace: 1,
				Receivers: []*sirenv1beta1.ReceiverMetadata{
					{Id: 30},
				},
			},
			{
				Id:        2,
				Urn:       "sample-urn-2",
				Namespace: 2,
				Receivers: []*sirenv1beta1.ReceiverMetadata{
					{Id: 33},
				},
			},
		}

		client := new(mocks.SirenServiceClient)
		client.On("ListSubscriptions", ctx, &sirenv1beta1.ListSubscriptionsRequest{Metadata: map[string]string{
			"group_id":      groupID,
			"resource_id":   resourceID,
			"resource_type": resourceType,
		}}).
			Return(&sirenv1beta1.ListSubscriptionsResponse{
				Subscriptions: subscriptions,
			}, nil)
		defer client.AssertExpectations(t)

		service := alert.NewSubscriptionService(client)
		result, err := service.GetSubscriptions(ctx, groupID, resourceID, resourceType)
		assert.NoError(t, err)
		assert.Equal(t, subscriptions, result)
	})

	t.Run("should return if client return error", func(t *testing.T) {
		expectedError := status.Error(codes.Internal, "Internal")

		client := new(mocks.SirenServiceClient)
		client.On("ListSubscriptions", ctx, &sirenv1beta1.ListSubscriptionsRequest{Metadata: map[string]string{}}).
			Return(nil, expectedError)
		defer client.AssertExpectations(t)

		service := alert.NewSubscriptionService(client)
		_, err := service.GetSubscriptions(ctx, "", "", "")
		assert.ErrorIs(t, err, expectedError)
	})
}

func TestSubscriptionServiceDeleteSubscription(t *testing.T) {
	ctx := context.TODO()
	subscriptionID := 203

	t.Run("should not return error success", func(t *testing.T) {
		client := new(mocks.SirenServiceClient)
		client.On("DeleteSubscription", ctx, &sirenv1beta1.DeleteSubscriptionRequest{Id: uint64(subscriptionID)}).
			Return(nil, nil)
		defer client.AssertExpectations(t)

		service := alert.NewSubscriptionService(client)
		err := service.DeleteSubscription(ctx, subscriptionID)
		assert.NoError(t, err)
	})

	t.Run("should return not found error if optimus return NotFound code", func(t *testing.T) {
		expectedError := status.Error(codes.NotFound, "Not Found")

		client := new(mocks.SirenServiceClient)
		client.On("DeleteSubscription", ctx, &sirenv1beta1.DeleteSubscriptionRequest{Id: uint64(subscriptionID)}).
			Return(nil, expectedError)
		defer client.AssertExpectations(t)

		service := alert.NewSubscriptionService(client)
		err := service.DeleteSubscription(ctx, subscriptionID)
		assert.ErrorIs(t, err, alert.ErrSubscriptionNotFound)
	})

	t.Run("should return if client return error", func(t *testing.T) {
		expectedError := status.Error(codes.Internal, "Internal")

		client := new(mocks.SirenServiceClient)
		client.On("DeleteSubscription", ctx, &sirenv1beta1.DeleteSubscriptionRequest{Id: uint64(subscriptionID)}).
			Return(nil, expectedError)
		defer client.AssertExpectations(t)

		service := alert.NewSubscriptionService(client)
		err := service.DeleteSubscription(ctx, subscriptionID)
		assert.ErrorIs(t, err, expectedError)
	})
}
