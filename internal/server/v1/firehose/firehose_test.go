package firehose_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/goto/dex/internal/server/v1/alert"
	"github.com/goto/dex/internal/server/v1/firehose"
	"github.com/goto/dex/mocks"
)

func TestRoutesCreateFirehose(t *testing.T) {
	method := http.MethodPost
	path := "/firehoses"

	t.Run("should return 200 with firehose", func(t *testing.T) {

		firehoseRes := &entropyv1beta1.CreateResourceRequest{
			Resource: &entropyv1beta1.Resource{
				Name: "test-firehose",
				Kind: "firehose",
			},
		}

		clientEntropy := new(mocks.ResourceServiceClient)
		clientEntropy.On("CreateResource", mock.Anything, firehoseRes).Return(&entropyv1beta1.CreateResourceResponse{
			Resource: firehoseRes.Resource,
		}, nil)
		defer clientEntropy.AssertExpectations(t)

		clientShield := new(mocks.ShieldServiceClient)
		clientCompass := new(mocks.CompassServiceClient)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := chi.NewRouter()
		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
		router.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusOK, response.Code)
		resultJSON := response.Body.Bytes()
		expectedJSON, err := json.Marshal(firehoseRes)
		require.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), string(resultJSON))
	})

	t.Run("should return 400 when firehose name is missing", func(t *testing.T) {
		clientEntropy := new(mocks.ResourceServiceClient)
		clientShield := new(mocks.ShieldServiceClient)
		clientCompass := new(mocks.CompassServiceClient)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := chi.NewRouter()
		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
		router.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("should return 500 for internal error", func(t *testing.T) {
		clientError := status.Error(codes.Internal, "Internal")

		clientEntropy := new(mocks.ResourceServiceClient)
		clientEntropy.On("CreateResource", mock.Anything, mock.Anything).Return(nil, clientError)
		defer clientEntropy.AssertExpectations(t)

		clientShield := new(mocks.ShieldServiceClient)
		clientCompass := new(mocks.CompassServiceClient)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := chi.NewRouter()
		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

}

func TestRoutesListFirehose(t *testing.T) {
	method := http.MethodGet
	path := "/firehoses"

	t.Run("should return 200 with list of firehoses", func(t *testing.T) {
		firehoseRes := []*entropyv1beta1.Resource{
			{
				Name: "test-firehose",
				Kind: "firehose",
			},
		}

		clientEntropy := new(mocks.ResourceServiceClient)
		clientEntropy.On("ListResources", mock.Anything, &entropyv1beta1.ListResourcesRequest{
			Kind: "firehose",
		}).Return(&entropyv1beta1.ListResourcesResponse{
			Resources: firehoseRes,
		}, nil)
		defer clientEntropy.AssertExpectations(t)

		clientShield := new(mocks.ShieldServiceClient)
		clientCompass := new(mocks.CompassServiceClient)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := chi.NewRouter()
		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
		router.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusOK, response.Code)
		resultJSON := response.Body.Bytes()
		expectedJSON, err := json.Marshal(firehoseRes)
		require.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), string(resultJSON))
	})

	t.Run("should return 500 for internal error", func(t *testing.T) {
		clientError := status.Error(codes.Internal, "Internal")

		clientEntropy := new(mocks.ResourceServiceClient)
		clientEntropy.On("ListResources", mock.Anything, mock.Anything).Return(nil, clientError)
		defer clientEntropy.AssertExpectations(t)

		clientShield := new(mocks.ShieldServiceClient)
		clientCompass := new(mocks.CompassServiceClient)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := chi.NewRouter()
		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

}

func TestRoutesUpdateFirehose(t *testing.T) {
	method := http.MethodPut
	path := "/firehoses/{urn}"

	t.Run("should return 200 with updated firehose", func(t *testing.T) {
		firehoseRes := &entropyv1beta1.UpdateResourceRequest{
			Resource: &entropyv1beta1.Resource{
				Name: "test-firehose",
				Kind: "firehose",
			},
		}

		clientEntropy := new(mocks.ResourceServiceClient)
		clientEntropy.On("UpdateResource", mock.Anything, firehoseRes).Return(&entropyv1beta1.UpdateResourceResponse{
			Resource: firehoseRes.Resource,
		}, nil)
		defer clientEntropy.AssertExpectations(t)

		clientShield := new(mocks.ShieldServiceClient)
		clientCompass := new(mocks.CompassServiceClient)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := chi.NewRouter()
		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
		router.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusOK, response.Code)
		resultJSON := response.Body.Bytes()
		expectedJSON, err := json.Marshal(firehoseRes)
		require.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), string(resultJSON))
	})

	t.Run("should return 400 when firehose name is missing", func(t *testing.T) {
		clientEntropy := new(mocks.ResourceServiceClient)
		clientShield := new(mocks.ShieldServiceClient)
		clientCompass := new(mocks.CompassServiceClient)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := chi.NewRouter()
		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
		router.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("should return 500 for internal error", func(t *testing.T) {
		clientError := status.Error(codes.Internal, "Internal")

		clientEntropy := new(mocks.ResourceServiceClient)
		clientEntropy.On("UpdateResource", mock.Anything, mock.Anything).Return(nil, clientError)
		defer clientEntropy.AssertExpectations(t)

		clientShield := new(mocks.ShieldServiceClient)
		clientCompass := new(mocks.CompassServiceClient)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := chi.NewRouter()
		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

}

func TestRoutesGetFirehose(t *testing.T) {
	method := http.MethodGet
	path := "/firehoses/{urn}"

	t.Run("should return 200 with firehose", func(t *testing.T) {
		firehoseRes := &entropyv1beta1.GetResourceRequest{
			Urn: "test-firehose",
		}

		clientEntropy := new(mocks.ResourceServiceClient)
		clientEntropy.On("GetResource", mock.Anything, firehoseRes).Return(&entropyv1beta1.GetResourceResponse{
			Resource: &entropyv1beta1.Resource{
				Name: "test-firehose",
				Kind: "firehose",
			},
		}, nil)
		defer clientEntropy.AssertExpectations(t)

		clientShield := new(mocks.ShieldServiceClient)
		clientCompass := new(mocks.CompassServiceClient)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		request = request.WithContext(firehoseRes.Context())
		router := chi.NewRouter()
		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
		router.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusOK, response.Code)
		resultJSON := response.Body.Bytes()
		expectedJSON, err := json.Marshal(firehoseRes)
		require.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), string(resultJSON))
	})

	t.Run("should return 400 when firehose urn is missing", func(t *testing.T) {
		clientEntropy := new(mocks.ResourceServiceClient)
		clientShield := new(mocks.ShieldServiceClient)
		clientCompass := new(mocks.CompassServiceClient)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		request = request.WithContext(contextWithUrn(""))
		router := chi.NewRouter()
		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
		router.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("should return 500 for internal error", func(t *testing.T) {
		clientError := status.Error(codes.Internal, "Internal")

		firehoseRes := &entropyv1beta1.GetResourceRequest{
			Urn: "test-firehose",
		}

		clientEntropy := new(mocks.ResourceServiceClient)
		clientEntropy.On("GetResource", mock.Anything, firehoseRes).Return(nil, clientError)
		defer clientEntropy.AssertExpectations(t)

		clientShield := new(mocks.ShieldServiceClient)
		clientCompass := new(mocks.CompassServiceClient)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		request = request.WithContext(firehoseRes.Context())
		router := chi.NewRouter()
		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

}

func TestRoutesGetFirehoseHistory(t *testing.T) {
	method := http.MethodGet
	path := "/firehoses/{urn}/history"

	t.Run("should return 200 with firehose history", func(t *testing.T) {
		firehoseRes := &entropyv1beta1.GetResourceRevisionsRequest{
			Urn: "test-firehose",
		}

		clientEntropy := new(mocks.ResourceServiceClient)
		clientEntropy.On("GetResourceHistory", mock.Anything, firehoseRes).Return(&entropyv1beta1.GetResourceHistoryResponse{
			ResourceHistory: []*entropyv1beta1.ResourceHistory{
				{
					Resource: &entropyv1beta1.Resource{
						Name: "test-firehose",
						Kind: "firehose",
					},
				},
			},
		}, nil)
		defer clientEntropy.AssertExpectations(t)

		clientShield := new(mocks.ShieldServiceClient)
		clientCompass := new(mocks.CompassServiceClient)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		request = request.WithContext(firehoseRes.Context())
		router := chi.NewRouter()
		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
		router.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusOK, response.Code)
		resultJSON := response.Body.Bytes()
		expectedJSON, err := json.Marshal(firehoseRes)
		require.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), string(resultJSON))
	})

	t.Run("should return 400 when firehose urn is missing", func(t *testing.T) {
		clientEntropy := new(mocks.ResourceServiceClient)
		clientShield := new(mocks.ShieldServiceClient)
		clientCompass := new(mocks.CompassServiceClient)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		request = request.WithContext(contextWithUrn(""))
		router := chi.NewRouter()
		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
		router.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("should return 500 for internal error", func(t *testing.T) {
		clientError := status.Error(codes.Internal, "Internal")

		firehoseRes := &entropyv1beta1.GetResourceRevisionsRequest{
			Urn: "test-firehose",
		}

		clientEntropy := new(mocks.ResourceServiceClient)
		clientEntropy.On("GetResourceHistory", mock.Anything, firehoseRes).Return(nil, clientError)
		defer clientEntropy.AssertExpectations(t)

		clientShield := new(mocks.ShieldServiceClient)
		clientCompass := new(mocks.CompassServiceClient)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		request = request.WithContext(firehoseRes.Context())
		router := chi.NewRouter()
		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})

}
