package firehose_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	shieldv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/shield/v1beta1"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/goto/dex/internal/server/reqctx"
	"github.com/goto/dex/internal/server/v1/firehose"
	"github.com/goto/dex/mocks"
)

func TestRoutesCreateFirehose(t *testing.T) {
	method := http.MethodPost
	path := "/"
	groupID := uuid.NewString()
	sampleFirehose := "sample-firehose"

	shieldGroup := &shieldv1beta1.Group{
		Slug: "test-group",
	}

	projectID := "5dab4194-9516-421a-aafe-72fd3d96ec56"
	shieldProject := &shieldv1beta1.Project{
		Id:    projectID,
		Slug:  projectID,
		OrgId: "orgID",
	}

	validJSONPayload := fmt.Sprintf(`
	{
		"configs": {
			"env_vars": {
				"SCHEMA_REGISTRY_STENCIL_URLS": "http://g-godata-systems-stencil-v1beta1-ingress.golabs.io/v1beta1/namespaces/gojek/schemas/esb-log-entities",
				"SOURCE_KAFKA_BROKERS": "10.84.52.20:6668,10.84.52.21:6668,10.84.52.22:6668",
				"SOURCE_KAFKA_CONSUMER_GROUP_ID": "g-pilotdata-gl-stewart-log-sink-firehose-0001",
				"SOURCE_KAFKA_TOPIC": "dagger-kafka-sink-es-enrichment-smoke-test-d5b2IHuEDC"
			},
			"kube_cluster": "orn:entropy:kubernetes:g-pilotdata-gl:sourcing",
			"replicas": 1,
			"stream_name": "mainstream"
		},
		"description": "stewart-log-sinkstewart-log-sinkstewart-log-sinkstewart-log-sinkstewart-log-sinkstewart-log-sinkstewart-log-sinkstewart-log-sinkstewart-log-sink",
		"group": "%s",
		"name": "%s",
		"project": "%s",
		"title": "%s"
	}`, groupID, sampleFirehose, projectID, sampleFirehose)

	t.Run("should return 200 with firehose", func(t *testing.T) {
		shieldClient := new(mocks.ShieldServiceClient)
		shieldClient.On("GetGroup", mock.Anything, &shieldv1beta1.GetGroupRequest{Id: groupID}).
			Return(&shieldv1beta1.GetGroupResponse{Group: shieldGroup}, nil)

		shieldClient.On("GetProject", mock.Anything, &shieldv1beta1.GetProjectRequest{Id: projectID}).
			Return(&shieldv1beta1.GetProjectResponse{Project: shieldProject}, nil)

		defer shieldClient.AssertExpectations(t)

		entropyClient := new(mocks.ResourceServiceClient)

		createResourceResp := &entropyv1beta1.CreateResourceResponse{
			Resource: &entropyv1beta1.Resource{
				Urn:       "test-urn",
				Kind:      "firehose",
				Name:      sampleFirehose,
				Project:   projectID,
				CreatedAt: timestamppb.Now(),
				UpdatedAt: timestamppb.Now(),
				CreatedBy: "ayushi.sharma@gojek.com",
				UpdatedBy: "ayushi.sharma@gojek.com",
				Labels: map[string]string{
					"title":       "",
					"group":       "",
					"description": "",
					"stream_name": "",
				},
				Spec: &entropyv1beta1.ResourceSpec{
					Dependencies: []*entropyv1beta1.ResourceDependency{
						{
							Key:   "kube_cluster",
							Value: "test-cluster",
						},
					},
					Configs: &structpb.Value{
						Kind: &structpb.Value_StructValue{
							StructValue: &structpb.Struct{
								Fields: map[string]*structpb.Value{
									"chart_values": {
										Kind: &structpb.Value_StructValue{
											StructValue: &structpb.Struct{
												Fields: map[string]*structpb.Value{
													"image_tag": {
														Kind: &structpb.Value_StringValue{
															StringValue: "test_image_tag",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				State: &entropyv1beta1.ResourceState{},
			},
		}

		entropyClient.On("CreateResource", mock.Anything, mock.AnythingOfType("*entropyv1beta1.CreateResourceRequest")).Return(createResourceResp, nil)
		defer shieldClient.AssertExpectations(t)
		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, strings.NewReader(validJSONPayload))
		request.Header.Set("Content-Type", "application/json")
		router := getRouter()

		firehose.Routes(entropyClient, shieldClient, nil, nil, "", "")(router)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("should return 400 when firehose name is missing", func(t *testing.T) {
		clientEntropy := new(mocks.ResourceServiceClient)
		clientShield := new(mocks.ShieldServiceClient)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := chi.NewRouter()
		firehose.Routes(clientEntropy, clientShield, nil, nil, "", "")(router)
		router.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusBadRequest, response.Code)

	})

	// t.Run("should return 500 for internal error", func(t *testing.T) {
	// 	clientError := status.Error(codes.Internal, "Internal")

	// 	clientEntropy := new(mocks.ResourceServiceClient)
	// 	clientEntropy.On("CreateResource", mock.Anything, mock.Anything).Return(nil, clientError)
	// 	defer clientEntropy.AssertExpectations(t)

	// 	clientShield := new(mocks.ShieldServiceClient)

	// 	response := httptest.NewRecorder()
	// 	request := httptest.NewRequest(method, path, nil)
	// 	router := chi.NewRouter()
	// 	firehose.Routes(clientEntropy, clientShield, nil, nil, "", "")(router)
	// 	router.ServeHTTP(response, request)

	// 	assert.Equal(t, http.StatusInternalServerError, response.Code)
	// })

}

// func TestRoutesListFirehose(t *testing.T) {
// 	method := http.MethodGet
// 	path := "/firehoses"

// 	t.Run("should return 200 with list of firehoses", func(t *testing.T) {
// 		firehoseRes := []*entropyv1beta1.Resource{
// 			{
// 				Name: "test-firehose",
// 				Kind: "firehose",
// 			},
// 		}

// 		clientEntropy := new(mocks.ResourceServiceClient)
// 		clientEntropy.On("ListResources", mock.Anything, &entropyv1beta1.ListResourcesRequest{
// 			Kind: "firehose",
// 		}).Return(&entropyv1beta1.ListResourcesResponse{
// 			Resources: firehoseRes,
// 		}, nil)
// 		defer clientEntropy.AssertExpectations(t)

// 		clientShield := new(mocks.ShieldServiceClient)
// 		clientCompass := new(mocks.CompassServiceClient)

// 		response := httptest.NewRecorder()
// 		request := httptest.NewRequest(method, path, nil)
// 		router := chi.NewRouter()
// 		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
// 		router.ServeHTTP(response, request)

// 		// assert
// 		assert.Equal(t, http.StatusOK, response.Code)
// 		resultJSON := response.Body.Bytes()
// 		expectedJSON, err := json.Marshal(firehoseRes)
// 		require.NoError(t, err)
// 		assert.JSONEq(t, string(expectedJSON), string(resultJSON))
// 	})

// 	t.Run("should return 500 for internal error", func(t *testing.T) {
// 		clientError := status.Error(codes.Internal, "Internal")

// 		clientEntropy := new(mocks.ResourceServiceClient)
// 		clientEntropy.On("ListResources", mock.Anything, mock.Anything).Return(nil, clientError)
// 		defer clientEntropy.AssertExpectations(t)

// 		clientShield := new(mocks.ShieldServiceClient)
// 		clientCompass := new(mocks.CompassServiceClient)

// 		response := httptest.NewRecorder()
// 		request := httptest.NewRequest(method, path, nil)
// 		router := chi.NewRouter()
// 		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
// 		router.ServeHTTP(response, request)

// 		assert.Equal(t, http.StatusInternalServerError, response.Code)
// 	})

// }

// func TestRoutesUpdateFirehose(t *testing.T) {
// 	method := http.MethodPut
// 	path := "/firehoses/{urn}"

// 	t.Run("should return 200 with updated firehose", func(t *testing.T) {
// 		firehoseRes := &entropyv1beta1.UpdateResourceRequest{
// 			Resource: &entropyv1beta1.Resource{
// 				Name: "test-firehose",
// 				Kind: "firehose",
// 			},
// 		}

// 		clientEntropy := new(mocks.ResourceServiceClient)
// 		clientEntropy.On("UpdateResource", mock.Anything, firehoseRes).Return(&entropyv1beta1.UpdateResourceResponse{
// 			Resource: firehoseRes.Resource,
// 		}, nil)
// 		defer clientEntropy.AssertExpectations(t)

// 		clientShield := new(mocks.ShieldServiceClient)
// 		clientCompass := new(mocks.CompassServiceClient)

// 		response := httptest.NewRecorder()
// 		request := httptest.NewRequest(method, path, nil)
// 		router := chi.NewRouter()
// 		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
// 		router.ServeHTTP(response, request)

// 		// assert
// 		assert.Equal(t, http.StatusOK, response.Code)
// 		resultJSON := response.Body.Bytes()
// 		expectedJSON, err := json.Marshal(firehoseRes)
// 		require.NoError(t, err)
// 		assert.JSONEq(t, string(expectedJSON), string(resultJSON))
// 	})

// 	t.Run("should return 400 when firehose name is missing", func(t *testing.T) {
// 		clientEntropy := new(mocks.ResourceServiceClient)
// 		clientShield := new(mocks.ShieldServiceClient)
// 		clientCompass := new(mocks.CompassServiceClient)

// 		response := httptest.NewRecorder()
// 		request := httptest.NewRequest(method, path, nil)
// 		router := chi.NewRouter()
// 		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
// 		router.ServeHTTP(response, request)

// 		// assert
// 		assert.Equal(t, http.StatusBadRequest, response.Code)
// 	})

// 	t.Run("should return 500 for internal error", func(t *testing.T) {
// 		clientError := status.Error(codes.Internal, "Internal")

// 		clientEntropy := new(mocks.ResourceServiceClient)
// 		clientEntropy.On("UpdateResource", mock.Anything, mock.Anything).Return(nil, clientError)
// 		defer clientEntropy.AssertExpectations(t)

// 		clientShield := new(mocks.ShieldServiceClient)
// 		clientCompass := new(mocks.CompassServiceClient)

// 		response := httptest.NewRecorder()
// 		request := httptest.NewRequest(method, path, nil)
// 		router := chi.NewRouter()
// 		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
// 		router.ServeHTTP(response, request)

// 		assert.Equal(t, http.StatusInternalServerError, response.Code)
// 	})

// }

// func TestRoutesGetFirehose(t *testing.T) {
// 	method := http.MethodGet
// 	path := "/firehoses/{urn}"

// 	t.Run("should return 200 with firehose", func(t *testing.T) {
// 		firehoseRes := &entropyv1beta1.GetResourceRequest{
// 			Urn: "test-firehose",
// 		}

// 		clientEntropy := new(mocks.ResourceServiceClient)
// 		clientEntropy.On("GetResource", mock.Anything, firehoseRes).Return(&entropyv1beta1.GetResourceResponse{
// 			Resource: &entropyv1beta1.Resource{
// 				Name: "test-firehose",
// 				Kind: "firehose",
// 			},
// 		}, nil)
// 		defer clientEntropy.AssertExpectations(t)

// 		clientShield := new(mocks.ShieldServiceClient)
// 		clientCompass := new(mocks.CompassServiceClient)

// 		response := httptest.NewRecorder()
// 		request := httptest.NewRequest(method, path, nil)
// 		request = request.WithContext(firehoseRes.Context())
// 		router := chi.NewRouter()
// 		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
// 		router.ServeHTTP(response, request)

// 		// assert
// 		assert.Equal(t, http.StatusOK, response.Code)
// 		resultJSON := response.Body.Bytes()
// 		expectedJSON, err := json.Marshal(firehoseRes)
// 		require.NoError(t, err)
// 		assert.JSONEq(t, string(expectedJSON), string(resultJSON))
// 	})

// 	t.Run("should return 400 when firehose urn is missing", func(t *testing.T) {
// 		clientEntropy := new(mocks.ResourceServiceClient)
// 		clientShield := new(mocks.ShieldServiceClient)
// 		clientCompass := new(mocks.CompassServiceClient)

// 		response := httptest.NewRecorder()
// 		request := httptest.NewRequest(method, path, nil)
// 		request = request.WithContext(contextWithUrn(""))
// 		router := chi.NewRouter()
// 		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
// 		router.ServeHTTP(response, request)

// 		// assert
// 		assert.Equal(t, http.StatusBadRequest, response.Code)
// 	})

// 	t.Run("should return 500 for internal error", func(t *testing.T) {
// 		clientError := status.Error(codes.Internal, "Internal")

// 		firehoseRes := &entropyv1beta1.GetResourceRequest{
// 			Urn: "test-firehose",
// 		}

// 		clientEntropy := new(mocks.ResourceServiceClient)
// 		clientEntropy.On("GetResource", mock.Anything, firehoseRes).Return(nil, clientError)
// 		defer clientEntropy.AssertExpectations(t)

// 		clientShield := new(mocks.ShieldServiceClient)
// 		clientCompass := new(mocks.CompassServiceClient)

// 		response := httptest.NewRecorder()
// 		request := httptest.NewRequest(method, path, nil)
// 		request = request.WithContext(firehoseRes.Context())
// 		router := chi.NewRouter()
// 		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
// 		router.ServeHTTP(response, request)

// 		assert.Equal(t, http.StatusInternalServerError, response.Code)
// 	})

// }

// func TestRoutesGetFirehoseHistory(t *testing.T) {
// 	method := http.MethodGet
// 	path := "/firehoses/{urn}/history"

// 	t.Run("should return 200 with firehose history", func(t *testing.T) {
// 		firehoseRes := &entropyv1beta1.GetResourceRevisionsRequest{
// 			Urn: "test-firehose",
// 		}

// 		clientEntropy := new(mocks.ResourceServiceClient)
// 		clientEntropy.On("GetResourceHistory", mock.Anything, firehoseRes).Return(&entropyv1beta1.GetResourceHistoryResponse{
// 			ResourceHistory: []*entropyv1beta1.ResourceHistory{
// 				{
// 					Resource: &entropyv1beta1.Resource{
// 						Name: "test-firehose",
// 						Kind: "firehose",
// 					},
// 				},
// 			},
// 		}, nil)
// 		defer clientEntropy.AssertExpectations(t)

// 		clientShield := new(mocks.ShieldServiceClient)
// 		clientCompass := new(mocks.CompassServiceClient)

// 		response := httptest.NewRecorder()
// 		request := httptest.NewRequest(method, path, nil)
// 		request = request.WithContext(firehoseRes.Context())
// 		router := chi.NewRouter()
// 		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
// 		router.ServeHTTP(response, request)

// 		// assert
// 		assert.Equal(t, http.StatusOK, response.Code)
// 		resultJSON := response.Body.Bytes()
// 		expectedJSON, err := json.Marshal(firehoseRes)
// 		require.NoError(t, err)
// 		assert.JSONEq(t, string(expectedJSON), string(resultJSON))
// 	})

// 	t.Run("should return 400 when firehose urn is missing", func(t *testing.T) {
// 		clientEntropy := new(mocks.ResourceServiceClient)
// 		clientShield := new(mocks.ShieldServiceClient)
// 		clientCompass := new(mocks.CompassServiceClient)

// 		response := httptest.NewRecorder()
// 		request := httptest.NewRequest(method, path, nil)
// 		request = request.WithContext(contextWithUrn(""))
// 		router := chi.NewRouter()
// 		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
// 		router.ServeHTTP(response, request)

// 		// assert
// 		assert.Equal(t, http.StatusBadRequest, response.Code)
// 	})

// 	t.Run("should return 500 for internal error", func(t *testing.T) {
// 		clientError := status.Error(codes.Internal, "Internal")

// 		firehoseRes := &entropyv1beta1.GetResourceRevisionsRequest{
// 			Urn: "test-firehose",
// 		}

// 		clientEntropy := new(mocks.ResourceServiceClient)
// 		clientEntropy.On("GetResourceHistory", mock.Anything, firehoseRes).Return(nil, clientError)
// 		defer clientEntropy.AssertExpectations(t)

// 		clientShield := new(mocks.ShieldServiceClient)
// 		clientCompass := new(mocks.CompassServiceClient)

// 		response := httptest.NewRecorder()
// 		request := httptest.NewRequest(method, path, nil)
// 		request = request.WithContext(firehoseRes.Context())
// 		router := chi.NewRouter()
// 		firehose.Routes(clientEntropy, clientShield, &alert.Service{}, clientCompass, "", "")(router)
// 		router.ServeHTTP(response, request)

// 		assert.Equal(t, http.StatusInternalServerError, response.Code)
// 	})

// }

func getRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(reqctx.WithRequestCtx())

	return router
}

func envVarsToStructValue(envVars map[string]string) *structpb.Struct {
	fields := make(map[string]*structpb.Value)
	for key, value := range envVars {
		fields[key] = &structpb.Value{
			Kind: &structpb.Value_StringValue{
				StringValue: value,
			},
		}
	}
	return &structpb.Struct{Fields: fields}
}
