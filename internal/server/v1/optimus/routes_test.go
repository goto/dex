package optimus_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	optimusv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/optimus/core/v1beta1"
	shieldv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/shield/v1beta1"
	"github.com/go-chi/chi/v5"
	"github.com/goto/dex/internal/server/v1/optimus"
	"github.com/goto/dex/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRoutesFindJobSpec(t *testing.T) {
	jobName := "sample-optimus-job-name"
	hostname := "optimus.staging.golabs.io:80"
	projectSlug := "sample-project"
	method := http.MethodGet
	path := fmt.Sprintf("/projects/%s/jobs/%s", projectSlug, jobName)

	t.Run("should return 200 with job spec", func(t *testing.T) {
		jobSpecRes := &optimusv1beta1.JobSpecificationResponse{
			ProjectName:   projectSlug,
			NamespaceName: "test-namespcace",
			Job: &optimusv1beta1.JobSpecification{
				Version:  1,
				Name:     "sample-job",
				Owner:    "goto",
				TaskName: "sample-task-name",
			},
		}

		projectRes := &shieldv1beta1.GetProjectResponse{
			Project: &shieldv1beta1.Project{
				Slug: "test-project",
				Metadata: newStruct(t, map[string]interface{}{
					"optimus_host": hostname,
				}),
			},
		}

		shieldClient := new(mocks.ShieldServiceClient)
		shieldClient.On("GetProject", mock.Anything, &shieldv1beta1.GetProjectRequest{
			Id: projectSlug,
		}).Return(projectRes, nil)

		optimusClient := new(optimus.OptimusClientMock)
		defer shieldClient.AssertExpectations(t)

		client := mocks.NewJobSpecificationServiceClient(t)
		optimusClient.On("BuildOptimusClient", hostname).Return(client, nil)
		defer optimusClient.AssertExpectations(t)

		client.On("GetJobSpecifications", mock.Anything, &optimusv1beta1.GetJobSpecificationsRequest{
			ProjectName: projectSlug,
			JobName:     jobName,
		}).Return(&optimusv1beta1.GetJobSpecificationsResponse{
			JobSpecificationResponses: []*optimusv1beta1.JobSpecificationResponse{
				jobSpecRes,
			},
		}, nil)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := chi.NewRouter()
		optimus.Routes(shieldClient, optimusClient)(router)
		router.ServeHTTP(response, request)
		// assert
		assert.Equal(t, http.StatusOK, response.Code)
		resultJSON := response.Body.Bytes()
		expectedJSON, err := json.Marshal(jobSpecRes)
		require.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), string(resultJSON))
	})

	t.Run("should return 404 if job could not be found", func(t *testing.T) {
		projectRes := &shieldv1beta1.GetProjectResponse{
			Project: &shieldv1beta1.Project{
				Slug: "test-project",
				Metadata: newStruct(t, map[string]interface{}{
					"optimus_host": hostname,
				}),
			},
		}

		shieldClient := new(mocks.ShieldServiceClient)
		shieldClient.On("GetProject", mock.Anything, &shieldv1beta1.GetProjectRequest{
			Id: projectSlug,
		}).Return(projectRes, nil)

		optimusClient := new(optimus.OptimusClientMock)
		defer shieldClient.AssertExpectations(t)

		client := mocks.NewJobSpecificationServiceClient(t)
		optimusClient.On("BuildOptimusClient", hostname).Return(client, nil)
		defer optimusClient.AssertExpectations(t)

		client.On("GetJobSpecifications", mock.Anything, &optimusv1beta1.GetJobSpecificationsRequest{
			ProjectName: projectSlug,
			JobName:     jobName,
		}).Return(&optimusv1beta1.GetJobSpecificationsResponse{}, nil)
		defer client.AssertExpectations(t)
		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := chi.NewRouter()
		optimus.Routes(shieldClient, optimusClient)(router)
		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("should return 500 for internal error", func(t *testing.T) {
		clientError := status.Error(codes.Internal, "Internal")

		projectRes := &shieldv1beta1.GetProjectResponse{
			Project: &shieldv1beta1.Project{
				Slug: "test-project",
				Metadata: newStruct(t, map[string]interface{}{
					"optimus_host": hostname,
				}),
			},
		}

		shieldClient := new(mocks.ShieldServiceClient)
		shieldClient.On("GetProject", mock.Anything, &shieldv1beta1.GetProjectRequest{
			Id: projectSlug,
		}).Return(projectRes, nil)

		optimusClient := new(optimus.OptimusClientMock)
		defer shieldClient.AssertExpectations(t)

		client := mocks.NewJobSpecificationServiceClient(t)
		optimusClient.On("BuildOptimusClient", hostname).Return(client, nil)
		defer optimusClient.AssertExpectations(t)
		client.On("GetJobSpecifications", mock.Anything, &optimusv1beta1.GetJobSpecificationsRequest{
			ProjectName: projectSlug,
			JobName:     jobName,
		}).Return(nil, clientError)
		defer client.AssertExpectations(t)
		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := chi.NewRouter()
		optimus.Routes(shieldClient, optimusClient)(router)
		router.ServeHTTP(response, request)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestRoutesListJobs(t *testing.T) {
	projectSlug := "sample-project"
	hostname := "optimus.staging.golabs.io:80"
	method := http.MethodGet
	path := fmt.Sprintf("/projects/%s/jobs", projectSlug)

	t.Run("should return 200 with jobs list", func(t *testing.T) {
		jobSpec1 := &optimusv1beta1.JobSpecificationResponse{
			ProjectName:   "test-project",
			NamespaceName: "test-namespcace",
			Job: &optimusv1beta1.JobSpecification{
				Version:  1,
				Name:     "sample-job1",
				Owner:    "goto",
				TaskName: "sample-task-name",
			},
		}

		jobSpec2 := &optimusv1beta1.JobSpecificationResponse{
			ProjectName: "test-project",
			Job: &optimusv1beta1.JobSpecification{
				Version:  1,
				Name:     "sample-job2",
				Owner:    "goto",
				TaskName: "sample-task-name",
			},
		}

		expectedJobSpecRes := []*optimusv1beta1.JobSpecificationResponse{
			jobSpec1,
			jobSpec2,
		}

		projectRes := &shieldv1beta1.GetProjectResponse{
			Project: &shieldv1beta1.Project{
				Slug: "test-project",
				Metadata: newStruct(t, map[string]interface{}{
					"optimus_host": hostname,
				}),
			},
		}

		shieldClient := new(mocks.ShieldServiceClient)
		shieldClient.On("GetProject", mock.Anything, &shieldv1beta1.GetProjectRequest{
			Id: projectSlug,
		}).Return(projectRes, nil)

		optimusClient := new(optimus.OptimusClientMock)
		defer shieldClient.AssertExpectations(t)

		client := mocks.NewJobSpecificationServiceClient(t)
		optimusClient.On("BuildOptimusClient", hostname).Return(client, nil)
		defer optimusClient.AssertExpectations(t)

		client.On("GetJobSpecifications", mock.Anything, &optimusv1beta1.GetJobSpecificationsRequest{
			ProjectName: projectSlug,
		}).Return(&optimusv1beta1.GetJobSpecificationsResponse{
			JobSpecificationResponses: expectedJobSpecRes,
		}, nil)
		defer client.AssertExpectations(t)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)

		router := chi.NewRouter()
		optimus.Routes(shieldClient, optimusClient)(router)
		router.ServeHTTP(response, request)

		// assert
		assert.Equal(t, http.StatusOK, response.Code)
		resultJSON := response.Body.Bytes()
		expectedJSON, err := json.Marshal(expectedJobSpecRes)
		require.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), string(resultJSON))
	})

	t.Run("should return 500 for internal error", func(t *testing.T) {
		clientError := status.Error(codes.Internal, "Internal")

		projectRes := &shieldv1beta1.GetProjectResponse{
			Project: &shieldv1beta1.Project{
				Slug: "test-project",
				Metadata: newStruct(t, map[string]interface{}{
					"optimus_host": hostname,
				}),
			},
		}

		shieldClient := new(mocks.ShieldServiceClient)
		shieldClient.On("GetProject", mock.Anything, &shieldv1beta1.GetProjectRequest{
			Id: projectSlug,
		}).Return(projectRes, nil)

		optimusClient := new(optimus.OptimusClientMock)
		defer shieldClient.AssertExpectations(t)

		client := mocks.NewJobSpecificationServiceClient(t)
		optimusClient.On("BuildOptimusClient", hostname).Return(client, nil)
		defer optimusClient.AssertExpectations(t)

		client.On("GetJobSpecifications", mock.Anything, &optimusv1beta1.GetJobSpecificationsRequest{
			ProjectName: projectSlug,
		}).Return(nil, clientError)
		defer client.AssertExpectations(t)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(method, path, nil)
		router := chi.NewRouter()
		optimus.Routes(shieldClient, optimusClient)(router)
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}
