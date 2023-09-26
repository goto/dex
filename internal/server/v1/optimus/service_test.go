package optimus_test

import (
	"context"
	"testing"

	optimusv1beta1grpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/optimus/core/v1beta1/corev1beta1grpc"
	optimusv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/optimus/core/v1beta1"
	shieldv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/shield/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/goto/dex/internal/server/v1/optimus"
	"github.com/goto/dex/mocks"
)

func TestServiceFindJobSpec(t *testing.T) {
	jobName := "sample-optimus-job-name"
	// projectName := "sample-project"
	projectSlug := "g-pilotdata-gl"

	t.Run("should return job spec using job name and project name from argument", func(t *testing.T) {

		hostname := "optimus.staging.golabs.io:80"

		projectRes := &shieldv1beta1.GetProjectResponse{

			Project: &shieldv1beta1.Project{
				Slug: "test-project",
				Metadata: newStruct(t, map[string]interface{}{
					"optimus_host": hostname,
				}),
			},
		}

		shielClient := new(mocks.ShieldServiceClient)
		shielClient.On("GetProject", context.TODO(), &shieldv1beta1.GetProjectRequest{
			Id: projectSlug,
		}).Return(projectRes, nil)

		jobSpecRes := &optimusv1beta1.JobSpecificationResponse{
			ProjectName:   "test-project",
			NamespaceName: "test-namespcace",
			Job: &optimusv1beta1.JobSpecification{
				Version:  1,
				Name:     "sample-job",
				Owner:    "goto",
				TaskName: "sample-task-name",
			},
		}

		shielClient.On("GetJobSpecifications", context.TODO(), &optimusv1beta1.GetJobSpecificationsRequest{
			ProjectName: projectSlug,
			JobName:     jobName,
		}).Return(&optimusv1beta1.GetJobSpecificationsResponse{
			JobSpecificationResponses: []*optimusv1beta1.JobSpecificationResponse{
				jobSpecRes,
			},
		}, nil)
		defer shielClient.AssertExpectations(t)

		optimusClient := new(optimus.OptimusClientMock)
		service := optimus.NewService(shielClient, optimusClient)

		client := optimusv1beta1grpc.NewJobSpecificationServiceClient(nil)
		optimusClient.On("BuildOptimusClient", context.Background(), hostname).Return(client, nil)
		job, err := service.FindJobSpec(context.TODO(), jobName, projectSlug)
		assert.NoError(t, err)
		assert.Equal(t, jobSpecRes, job)
	})

	// t.Run("should return not found, if job could not be found", func(t *testing.T) {
	// 	client := new(mocks.ShieldServiceClient)
	// 	client.On("GetJobSpecifications", context.TODO(), &optimusv1beta1.GetJobSpecificationsRequest{
	// 		ProjectName: projectName,
	// 		JobName:     jobName,
	// 	}).Return(&optimusv1beta1.GetJobSpecificationsResponse{}, nil)
	// 	defer client.AssertExpectations(t)

	// 	service := optimus.NewService(client, &optimus.OptimusClientMock{})

	// 	_, err := service.FindJobSpec(context.TODO(), jobName, projectName)
	// 	assert.ErrorIs(t, err, errors.ErrNotFound)
	// })

	// t.Run("should return error, if client fails", func(t *testing.T) {
	// 	expectedErr := status.Error(codes.Internal, "Internal")

	// 	client := new(mocks.ShieldServiceClient)
	// 	client.On("GetJobSpecifications", context.TODO(), &optimusv1beta1.GetJobSpecificationsRequest{
	// 		ProjectName: projectName,
	// 		JobName:     jobName,
	// 	}).Return(nil, expectedErr)
	// 	defer client.AssertExpectations(t)

	// 	service := optimus.NewService(client, &optimus.OptimusClientMock{})

	// 	_, err := service.FindJobSpec(context.TODO(), jobName, projectName)
	// 	assert.ErrorIs(t, err, expectedErr)
	// })
}

// func TestServiceListJobs(t *testing.T) {
// 	projectName := "test-project"
// 	t.Run("should return list of jobs using project in argument", func(t *testing.T) {
// 		jobSpecRes := &optimusv1beta1.JobSpecificationResponse{
// 			ProjectName:   projectName,
// 			NamespaceName: "test-namespcace",
// 			Job: &optimusv1beta1.JobSpecification{
// 				Version:  1,
// 				Name:     "sample-job",
// 				Owner:    "goto",
// 				TaskName: "sample-task-name",
// 			},
// 		}

// 		listJobsResp := &optimusv1beta1.GetJobSpecificationsResponse{
// 			JobSpecificationResponses: []*optimusv1beta1.JobSpecificationResponse{
// 				jobSpecRes,
// 			},
// 		}

// 		expectedResp := []*optimusv1beta1.JobSpecificationResponse{jobSpecRes}

// 		client := new(mocks.ShieldServiceClient)
// 		client.On("GetJobSpecifications", context.TODO(), &optimusv1beta1.GetJobSpecificationsRequest{
// 			ProjectName: projectName,
// 		}).Return(listJobsResp, nil)
// 		defer client.AssertExpectations(t)

// 		service := optimus.NewService(client, &optimus.OptimusClientMock{})

// 		resp, err := service.ListJobs(context.TODO(), projectName)
// 		assert.Equal(t, expectedResp, resp)
// 		assert.NoError(t, err)
// 	})

// 	t.Run("should return error if RPC request fails", func(t *testing.T) {
// 		expectedErr := status.Error(codes.Internal, "Internal")

// 		client := new(mocks.ShieldServiceClient)
// 		client.On("GetJobSpecifications", context.TODO(), &optimusv1beta1.GetJobSpecificationsRequest{
// 			ProjectName: projectName,
// 		}).Return(nil, expectedErr)
// 		defer client.AssertExpectations(t)

// 		service := optimus.NewService(client, &optimus.OptimusClientMock{})

// 		_, err := service.ListJobs(context.TODO(), projectName)
// 		assert.ErrorIs(t, err, expectedErr)
// 	})
// }

func newStruct(t *testing.T, d map[string]interface{}) *structpb.Struct {
	t.Helper()

	strct, err := structpb.NewStruct(d)
	require.NoError(t, err)
	return strct
}
