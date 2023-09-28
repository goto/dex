package optimus_test

import (
	"context"
	"testing"

	optimusv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/optimus/core/v1beta1"
	shieldv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/shield/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/goto/dex/internal/server/v1/optimus"
	"github.com/goto/dex/mocks"
	"github.com/goto/dex/pkg/errors"
)

func TestServiceFindJobSpec(t *testing.T) {
	jobName := "sample-optimus-job-name"
	projectSlug := "g-pilotdata-gl"
	hostname := "optimus.staging.golabs.io:80"

	t.Run("should return job spec using job name and project name from argument", func(t *testing.T) {

		projectRes := &shieldv1beta1.GetProjectResponse{
			Project: &shieldv1beta1.Project{
				Slug: "test-project",
				Metadata: newStruct(t, map[string]interface{}{
					"optimus_host": hostname,
				}),
			},
		}

		shieldClient := new(mocks.ShieldServiceClient)
		shieldClient.On("GetProject", context.TODO(), &shieldv1beta1.GetProjectRequest{
			Id: projectSlug,
		}).Return(projectRes, nil)

		jobSpecRes := &optimusv1beta1.JobSpecificationResponse{
			ProjectName:   "test-project",
			NamespaceName: "test-namespcace",
			Job: &optimusv1beta1.JobSpecification{
				Version:  1,
				Name:     jobName,
				Owner:    "goto",
				TaskName: "sample-task-name",
			},
		}

		optimusClient := new(optimus.OptimusClientMock)
		service := optimus.NewService(shieldClient, optimusClient)
		defer shieldClient.AssertExpectations(t)

		client := mocks.NewJobSpecificationServiceClient(t)
		optimusClient.On("BuildOptimusClient", hostname).Return(client, nil)
		defer optimusClient.AssertExpectations(t)

		client.On("GetJobSpecifications", context.TODO(), &optimusv1beta1.GetJobSpecificationsRequest{
			ProjectName: projectSlug,
			JobName:     jobName,
		}).Return(&optimusv1beta1.GetJobSpecificationsResponse{
			JobSpecificationResponses: []*optimusv1beta1.JobSpecificationResponse{
				jobSpecRes,
			},
		}, nil)
		defer client.AssertExpectations(t)
		job, err := service.FindJobSpec(context.TODO(), jobName, projectSlug)
		assert.NoError(t, err)
		assert.Equal(t, jobSpecRes, job)
	})

	t.Run("should return error, if project not found", func(t *testing.T) {
		expectedErr := status.Error(codes.NotFound, "Not found")

		shieldClient := new(mocks.ShieldServiceClient)
		shieldClient.On("GetProject", context.TODO(), &shieldv1beta1.GetProjectRequest{
			Id: projectSlug,
		}).Return(nil, expectedErr)

		optimusClient := new(optimus.OptimusClientMock)
		service := optimus.NewService(shieldClient, optimusClient)
		defer shieldClient.AssertExpectations(t)
		defer optimusClient.AssertExpectations(t)

		_, err := service.FindJobSpec(context.TODO(), jobName, projectSlug)
		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("should return error if metadata doesn't contain optimus hostname", func(t *testing.T) {
		projectRes := &shieldv1beta1.GetProjectResponse{
			Project: &shieldv1beta1.Project{
				Slug:     "test-project",
				Metadata: newStruct(t, map[string]interface{}{}),
			},
		}

		shieldClient := new(mocks.ShieldServiceClient)
		shieldClient.On("GetProject", context.TODO(), &shieldv1beta1.GetProjectRequest{
			Id: projectSlug,
		}).Return(projectRes, nil)

		optimusClient := new(optimus.OptimusClientMock)
		service := optimus.NewService(shieldClient, optimusClient)
		defer shieldClient.AssertExpectations(t)
		defer optimusClient.AssertExpectations(t)

		_, err := service.FindJobSpec(context.TODO(), jobName, projectSlug)
		assert.ErrorIs(t, err, optimus.ErrOptimusHostNotFound)
	})

	t.Run("should return error, if creation of optimus client fails to create", func(t *testing.T) {
		expectedErr := status.Error(codes.Internal, "Internal")

		projectRes := &shieldv1beta1.GetProjectResponse{
			Project: &shieldv1beta1.Project{
				Slug: "test-project",
				Metadata: newStruct(t, map[string]interface{}{
					"optimus_host": hostname,
				}),
			},
		}

		shieldClient := new(mocks.ShieldServiceClient)
		shieldClient.On("GetProject", context.TODO(), &shieldv1beta1.GetProjectRequest{
			Id: projectSlug,
		}).Return(projectRes, nil)

		optimusClient := new(optimus.OptimusClientMock)
		service := optimus.NewService(shieldClient, optimusClient)
		defer shieldClient.AssertExpectations(t)

		_ = mocks.NewJobSpecificationServiceClient(t)
		optimusClient.On("BuildOptimusClient", hostname).Return(nil, expectedErr)
		defer optimusClient.AssertExpectations(t)

		_, err := service.FindJobSpec(context.TODO(), jobName, projectSlug)
		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("should return not found, if job could not be found", func(t *testing.T) {
		projectRes := &shieldv1beta1.GetProjectResponse{
			Project: &shieldv1beta1.Project{
				Slug: "test-project",
				Metadata: newStruct(t, map[string]interface{}{
					"optimus_host": hostname,
				}),
			},
		}

		shieldClient := new(mocks.ShieldServiceClient)
		shieldClient.On("GetProject", context.TODO(), &shieldv1beta1.GetProjectRequest{
			Id: projectSlug,
		}).Return(projectRes, nil)

		optimusClient := new(optimus.OptimusClientMock)
		service := optimus.NewService(shieldClient, optimusClient)
		defer shieldClient.AssertExpectations(t)

		client := mocks.NewJobSpecificationServiceClient(t)
		optimusClient.On("BuildOptimusClient", hostname).Return(client, nil)
		defer optimusClient.AssertExpectations(t)

		client.On("GetJobSpecifications", context.TODO(), &optimusv1beta1.GetJobSpecificationsRequest{
			ProjectName: projectSlug,
			JobName:     jobName,
		}).Return(&optimusv1beta1.GetJobSpecificationsResponse{}, nil)
		defer client.AssertExpectations(t)

		_, err := service.FindJobSpec(context.TODO(), jobName, projectSlug)
		assert.ErrorIs(t, err, errors.ErrNotFound)
	})
}

func TestServiceListJobs(t *testing.T) {
	projectName := "test-project"
	hostname := "optimus.staging.golabs.io:80"

	t.Run("should return list of jobs using project in argument", func(t *testing.T) {
		jobSpecRes := &optimusv1beta1.JobSpecificationResponse{
			ProjectName:   projectName,
			NamespaceName: "test-namespcace",
			Job: &optimusv1beta1.JobSpecification{
				Version:  1,
				Name:     "sample-job",
				Owner:    "goto",
				TaskName: "sample-task-name",
			},
		}

		expectedResp := []*optimusv1beta1.JobSpecificationResponse{jobSpecRes}

		projectRes := &shieldv1beta1.GetProjectResponse{
			Project: &shieldv1beta1.Project{
				Slug: "test-project",
				Metadata: newStruct(t, map[string]interface{}{
					"optimus_host": hostname,
				}),
			},
		}

		shieldClient := new(mocks.ShieldServiceClient)
		shieldClient.On("GetProject", context.TODO(), &shieldv1beta1.GetProjectRequest{
			Id: projectName,
		}).Return(projectRes, nil)

		optimusClient := new(optimus.OptimusClientMock)
		service := optimus.NewService(shieldClient, optimusClient)
		defer shieldClient.AssertExpectations(t)

		client := mocks.NewJobSpecificationServiceClient(t)
		optimusClient.On("BuildOptimusClient", hostname).Return(client, nil)
		defer optimusClient.AssertExpectations(t)

		client.On("GetJobSpecifications", context.TODO(), &optimusv1beta1.GetJobSpecificationsRequest{
			ProjectName: projectName,
		}).Return(&optimusv1beta1.GetJobSpecificationsResponse{
			JobSpecificationResponses: []*optimusv1beta1.JobSpecificationResponse{
				jobSpecRes,
			},
		}, nil)
		defer client.AssertExpectations(t)
		resp, err := service.ListJobs(context.TODO(), projectName)
		assert.Equal(t, expectedResp, resp)
		assert.NoError(t, err)
	})

	t.Run("should return error, if project not found", func(t *testing.T) {
		expectedErr := status.Error(codes.NotFound, "Not found")

		shieldClient := new(mocks.ShieldServiceClient)
		shieldClient.On("GetProject", context.TODO(), &shieldv1beta1.GetProjectRequest{
			Id: projectName,
		}).Return(nil, expectedErr)

		optimusClient := new(optimus.OptimusClientMock)
		service := optimus.NewService(shieldClient, optimusClient)
		defer shieldClient.AssertExpectations(t)
		defer optimusClient.AssertExpectations(t)

		_, err := service.ListJobs(context.TODO(), projectName)
		assert.ErrorIs(t, err, expectedErr)
	})

	t.Run("should return error if metadata doesn't contain optimus hostname", func(t *testing.T) {
		projectRes := &shieldv1beta1.GetProjectResponse{
			Project: &shieldv1beta1.Project{
				Slug:     "test-project",
				Metadata: newStruct(t, map[string]interface{}{}),
			},
		}

		shieldClient := new(mocks.ShieldServiceClient)
		shieldClient.On("GetProject", context.TODO(), &shieldv1beta1.GetProjectRequest{
			Id: projectName,
		}).Return(projectRes, nil)

		optimusClient := new(optimus.OptimusClientMock)
		service := optimus.NewService(shieldClient, optimusClient)
		defer shieldClient.AssertExpectations(t)
		defer optimusClient.AssertExpectations(t)

		_, err := service.ListJobs(context.TODO(), projectName)
		assert.ErrorIs(t, err, optimus.ErrOptimusHostNotFound)
	})

	t.Run("should return error, if creation of optimus client fails to create", func(t *testing.T) {
		expectedErr := status.Error(codes.Internal, "Internal")

		projectRes := &shieldv1beta1.GetProjectResponse{
			Project: &shieldv1beta1.Project{
				Slug: "test-project",
				Metadata: newStruct(t, map[string]interface{}{
					"optimus_host": hostname,
				}),
			},
		}

		shieldClient := new(mocks.ShieldServiceClient)
		shieldClient.On("GetProject", context.TODO(), &shieldv1beta1.GetProjectRequest{
			Id: projectName,
		}).Return(projectRes, nil)

		optimusClient := new(optimus.OptimusClientMock)
		service := optimus.NewService(shieldClient, optimusClient)
		defer shieldClient.AssertExpectations(t)

		_ = mocks.NewJobSpecificationServiceClient(t)
		optimusClient.On("BuildOptimusClient", hostname).Return(nil, expectedErr)
		defer optimusClient.AssertExpectations(t)

		_, err := service.ListJobs(context.TODO(), projectName)
		assert.ErrorIs(t, err, expectedErr)
	})

}

func newStruct(t *testing.T, d map[string]interface{}) *structpb.Struct {
	t.Helper()

	strct, err := structpb.NewStruct(d)
	require.NoError(t, err)
	return strct
}
