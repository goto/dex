package alert_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	sirenv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/siren/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goto/dex/internal/server/v1/alert"
	"github.com/goto/dex/mocks"
	"github.com/goto/dex/pkg/errors"
	"github.com/goto/dex/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlerFindSubscription(t *testing.T) {
	t.Run("should return subscription on success", func(t *testing.T) {
		subscriptionID := 110
		subscription := &sirenv1beta1.Subscription{
			Id:        uint64(subscriptionID),
			Urn:       "sample-http-call-urn",
			Namespace: 1,
			Receivers: []*sirenv1beta1.ReceiverMetadata{
				{Id: 32},
			},
		}

		shieldClient := new(mocks.ShieldServiceClient)
		sirenClient := new(mocks.SirenServiceClient)
		sirenClient.On("GetSubscription", mock.Anything, &sirenv1beta1.GetSubscriptionRequest{Id: subscription.Id}).
			Return(&sirenv1beta1.GetSubscriptionResponse{
				Subscription: subscription,
			}, nil)
		defer sirenClient.AssertExpectations(t)

		service := alert.NewSubscriptionService(sirenClient)
		handler := alert.NewHandler(service, shieldClient)

		response := httptest.NewRecorder()
		request := httptest.NewRequest("", "/", nil)
		request = tests.SetChiParams(request, map[string]string{
			"subscription_id": fmt.Sprintf("%d", subscriptionID),
		})
		handler.FindSubscription(response, request)

		// assert status
		assert.Equal(t, http.StatusOK, response.Result().StatusCode)

		// assert response
		resultJSON := response.Body.Bytes()
		expectedJSON, err := json.Marshal(map[string]interface{}{
			"subscription": subscription,
		})
		require.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), string(resultJSON))
	})

	t.Run("should return 404 if id is not found", func(t *testing.T) {
		subscriptionID := 102
		expectedError := status.Error(codes.NotFound, "not found")

		shieldClient := new(mocks.ShieldServiceClient)
		sirenClient := new(mocks.SirenServiceClient)
		sirenClient.On("GetSubscription", mock.Anything, &sirenv1beta1.GetSubscriptionRequest{Id: uint64(subscriptionID)}).
			Return(nil, expectedError)
		defer sirenClient.AssertExpectations(t)

		service := alert.NewSubscriptionService(sirenClient)
		handler := alert.NewHandler(service, shieldClient)

		response := httptest.NewRecorder()
		request := httptest.NewRequest("", "/", nil)
		request = tests.SetChiParams(request, map[string]string{
			"subscription_id": fmt.Sprintf("%d", subscriptionID),
		})
		handler.FindSubscription(response, request)

		// assert status
		assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)

		// assert body
		resultJSON := response.Body.Bytes()
		expectedJSON, err := json.Marshal(errors.Error{
			Status:  http.StatusNotFound,
			Message: alert.ErrSubscriptionNotFound.Error(),
		})
		require.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), string(resultJSON))
	})
}
