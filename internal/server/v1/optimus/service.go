package optimus

import (
	"context"

	optimusv1beta1grpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/optimus/core/v1beta1/corev1beta1grpc"
	shieldv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/shield/v1beta1/shieldv1beta1grpc"
	optimusv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/optimus/core/v1beta1"
	shieldv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/shield/v1beta1"
	"github.com/goto/dex/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service struct {
	shieldClient shieldv1beta1rpc.ShieldServiceClient
}

func NewService(shieldClient shieldv1beta1rpc.ShieldServiceClient) *Service {
	return &Service{
		shieldClient: shieldClient,
	}
}

func (svc *Service) FindJobSpec(ctx context.Context, jobName, projectName string) (*optimusv1beta1.JobSpecificationResponse, error) {
	optimusCl, err := svc.createOptimusClient(ctx, projectName)
	if err != nil {
		return nil, err
	}

	res, err := optimusCl.GetJobSpecifications(ctx, &optimusv1beta1.GetJobSpecificationsRequest{
		ProjectName: projectName,
		JobName:     jobName,
	})
	if err != nil {
		return nil, err
	}

	list := res.JobSpecificationResponses
	if len(list) == 0 {
		return nil, errors.ErrNotFound
	}

	return list[0], nil
}

func (svc *Service) ListJobs(ctx context.Context, projectName string) ([]*optimusv1beta1.JobSpecificationResponse, error) {
	optimusCl, err := svc.createOptimusClient(ctx, projectName)
	if err != nil {
		return nil, err
	}

	res, err := optimusCl.GetJobSpecifications(ctx, &optimusv1beta1.GetJobSpecificationsRequest{
		ProjectName: projectName,
	})
	if err != nil {
		return nil, err
	}

	return res.JobSpecificationResponses, nil
}

func (svc *Service) createOptimusClient(ctx context.Context, projectName string) (optimusv1beta1grpc.JobSpecificationServiceClient, error) {
	prj, err := svc.shieldClient.GetProject(ctx, &shieldv1beta1.GetProjectRequest{
		Id: projectName,
	})
	if err != nil {
		return nil, err
	}

	metadata := prj.Project.Metadata.AsMap()
	optimusHost, exists := metadata[optimusHostKey]
	if !exists {
		return nil, ErrOptimusHostNotFound
	}

	optimusHostStr, isString := optimusHost.(string)
	if !isString {
		return nil, ErrOptimusHostNotFound
	}

	optimusConn, err := grpc.Dial(optimusHostStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return optimusv1beta1grpc.NewJobSpecificationServiceClient(optimusConn), nil
}
