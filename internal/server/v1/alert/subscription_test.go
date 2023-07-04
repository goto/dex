package alert_test

import (
	"context"
	"errors"
	"testing"

	sirenv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/siren/v1beta1"
	"github.com/stretchr/testify/assert"

	"github.com/goto/dex/internal/server/v1/alert"
	"github.com/goto/dex/mocks"
)

func TestSubscriptionServiceFindSubscription(t *testing.T) {
	ctx := context.TODO()

	t.Run("should return subscription on success", func(t *testing.T) {
		subscriptionID := 102
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

	t.Run("should return if client return error", func(t *testing.T) {
		subscriptionID := 102
		expectedError := errors.New("sample-error")

		client := new(mocks.SirenServiceClient)
		client.On("GetSubscription", ctx, &sirenv1beta1.GetSubscriptionRequest{Id: uint64(subscriptionID)}).
			Return(nil, expectedError)
		defer client.AssertExpectations(t)

		service := alert.NewSubscriptionService(client)
		_, err := service.FindSubscription(ctx, subscriptionID)
		assert.ErrorIs(t, err, expectedError)
	})
}
