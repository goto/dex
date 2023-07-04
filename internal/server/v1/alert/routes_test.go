package alert_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	sirenv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/siren/v1beta1"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goto/dex/internal/server/v1/alert"
	"github.com/goto/dex/mocks"
	"github.com/goto/dex/pkg/errors"
)

func TestRoutesFindSubscription(t *testing.T) {
	subscriptionID := 102
	path := "/102"
	method := http.MethodGet

	t.Run("should return subscription on success", func(t *testing.T) {
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

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := chi.NewRouter()
		alert.SubscriptionRoutes(sirenClient, shieldClient)(router)
		router.ServeHTTP(response, request)

		// assert status
		assert.Equal(t, http.StatusOK, response.Code)

		// assert response
		resultJSON := response.Body.Bytes()
		expectedJSON, err := json.Marshal(map[string]interface{}{
			"subscription": subscription,
		})
		require.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), string(resultJSON))
	})

	t.Run("should return 404 if id is not found", func(t *testing.T) {
		expectedError := status.Error(codes.NotFound, "not found")

		shieldClient := new(mocks.ShieldServiceClient)
		sirenClient := new(mocks.SirenServiceClient)
		sirenClient.On("GetSubscription", mock.Anything, &sirenv1beta1.GetSubscriptionRequest{Id: uint64(subscriptionID)}).
			Return(nil, expectedError)
		defer sirenClient.AssertExpectations(t)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := chi.NewRouter()
		alert.SubscriptionRoutes(sirenClient, shieldClient)(router)
		router.ServeHTTP(response, request)

		// assert
		expectedStatusCode := http.StatusNotFound
		assert.Equal(t, expectedStatusCode, response.Code)
		resultJSON := response.Body.Bytes()
		expectedJSON, err := json.Marshal(errors.Error{
			Status:  expectedStatusCode,
			Message: alert.ErrSubscriptionNotFound.Error(),
		})
		require.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), string(resultJSON))
	})

	t.Run("should return 500 for internal error", func(t *testing.T) {
		expectedError := status.Error(codes.Internal, "Internal")

		shieldClient := new(mocks.ShieldServiceClient)
		sirenClient := new(mocks.SirenServiceClient)
		sirenClient.On("GetSubscription", mock.Anything, &sirenv1beta1.GetSubscriptionRequest{Id: uint64(subscriptionID)}).
			Return(nil, expectedError)
		defer sirenClient.AssertExpectations(t)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := chi.NewRouter()
		alert.SubscriptionRoutes(sirenClient, shieldClient)(router)
		router.ServeHTTP(response, request)

		// assert
		expectedStatusCode := http.StatusInternalServerError
		assert.Equal(t, expectedStatusCode, response.Code)
	})
}

func TestRoutesGetSubscriptions(t *testing.T) {
	groupID := "sample-shield-group-id"
	resourceID := "sample-dagger-id-or-urn"
	resourceType := "dagger"
	method := http.MethodGet

	t.Run("should return subscription on success", func(t *testing.T) {
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

		shieldClient := new(mocks.ShieldServiceClient)
		sirenClient := new(mocks.SirenServiceClient)
		sirenClient.On("ListSubscriptions", mock.Anything, &sirenv1beta1.ListSubscriptionsRequest{Metadata: map[string]string{
			"group_id":      groupID,
			"resource_id":   resourceID,
			"resource_type": resourceType,
		}}).
			Return(&sirenv1beta1.ListSubscriptionsResponse{
				Subscriptions: subscriptions,
			}, nil)
		defer sirenClient.AssertExpectations(t)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(
			method,
			fmt.Sprintf("/?group_id=%s&resource_id=%s&resource_type=%s", groupID, resourceID, resourceType),
			nil,
		)
		router := chi.NewRouter()
		alert.SubscriptionRoutes(sirenClient, shieldClient)(router)
		router.ServeHTTP(response, request)

		// assert status
		assert.Equal(t, http.StatusOK, response.Code)

		// assert response
		resultJSON := response.Body.Bytes()
		expectedJSON, err := json.Marshal(map[string]interface{}{
			"subscriptions": subscriptions,
		})
		require.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), string(resultJSON))
	})

	t.Run("should return 400 if both group_id and resource_id is not passed", func(t *testing.T) {
		shieldClient := new(mocks.ShieldServiceClient)
		sirenClient := new(mocks.SirenServiceClient)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, "/", nil)
		router := chi.NewRouter()
		alert.SubscriptionRoutes(sirenClient, shieldClient)(router)
		router.ServeHTTP(response, request)

		// assert status
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("should return 500 on internal error", func(t *testing.T) {
		expectedError := status.Error(codes.Internal, "Internal")

		shieldClient := new(mocks.ShieldServiceClient)
		sirenClient := new(mocks.SirenServiceClient)
		sirenClient.On("ListSubscriptions", mock.Anything, &sirenv1beta1.ListSubscriptionsRequest{Metadata: map[string]string{
			"group_id":      groupID,
			"resource_id":   resourceID,
			"resource_type": resourceType,
		}}).
			Return(nil, expectedError)
		defer sirenClient.AssertExpectations(t)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(
			method,
			fmt.Sprintf("/?group_id=%s&resource_id=%s&resource_type=%s", groupID, resourceID, resourceType),
			nil,
		)
		router := chi.NewRouter()
		alert.SubscriptionRoutes(sirenClient, shieldClient)(router)
		router.ServeHTTP(response, request)

		// assert status
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestRoutesDeleteSubscription(t *testing.T) {
	subscriptionID := 202
	path := "/202"
	method := http.MethodDelete

	t.Run("should return 200 on success", func(t *testing.T) {
		shieldClient := new(mocks.ShieldServiceClient)
		sirenClient := new(mocks.SirenServiceClient)
		sirenClient.On("DeleteSubscription", mock.Anything, &sirenv1beta1.DeleteSubscriptionRequest{Id: uint64(subscriptionID)}).
			Return(nil, nil)
		defer sirenClient.AssertExpectations(t)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := chi.NewRouter()
		alert.SubscriptionRoutes(sirenClient, shieldClient)(router)
		router.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("should return 404 if id is not found", func(t *testing.T) {
		expectedError := status.Error(codes.NotFound, "not found")

		shieldClient := new(mocks.ShieldServiceClient)
		sirenClient := new(mocks.SirenServiceClient)
		sirenClient.On("DeleteSubscription", mock.Anything, &sirenv1beta1.DeleteSubscriptionRequest{Id: uint64(subscriptionID)}).
			Return(nil, expectedError)
		defer sirenClient.AssertExpectations(t)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := chi.NewRouter()
		alert.SubscriptionRoutes(sirenClient, shieldClient)(router)
		router.ServeHTTP(response, request)

		// assert
		expectedStatusCode := http.StatusNotFound
		assert.Equal(t, expectedStatusCode, response.Code)
		resultJSON := response.Body.Bytes()
		expectedJSON, err := json.Marshal(errors.Error{
			Status:  expectedStatusCode,
			Message: alert.ErrSubscriptionNotFound.Error(),
		})
		require.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), string(resultJSON))
	})

	t.Run("should return 500 for internal error", func(t *testing.T) {
		expectedError := status.Error(codes.Internal, "Internal")

		shieldClient := new(mocks.ShieldServiceClient)
		sirenClient := new(mocks.SirenServiceClient)
		sirenClient.On("DeleteSubscription", mock.Anything, &sirenv1beta1.DeleteSubscriptionRequest{Id: uint64(subscriptionID)}).
			Return(nil, expectedError)
		defer sirenClient.AssertExpectations(t)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := chi.NewRouter()
		alert.SubscriptionRoutes(sirenClient, shieldClient)(router)
		router.ServeHTTP(response, request)

		// assert
		expectedStatusCode := http.StatusInternalServerError
		assert.Equal(t, expectedStatusCode, response.Code)
	})
}
