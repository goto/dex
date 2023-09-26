package optimus

import (
	"context"

	optimusv1beta1grpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/optimus/core/v1beta1/corev1beta1grpc"
	shieldv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/shield/v1beta1/shieldv1beta1grpc"
	optimusv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/optimus/core/v1beta1"
	shieldv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/shield/v1beta1"

	"github.com/goto/dex/pkg/errors"
)

type Service struct {
	shieldClient shieldv1beta1rpc.ShieldServiceClient
	cache        *Cache
	builder      OptimusClientBuilder
}

func NewService(shieldClient shieldv1beta1rpc.ShieldServiceClient, builder OptimusClientBuilder) *Service {
	return &Service{
		shieldClient: shieldClient,
		cache:        NewCache(),
		builder:      builder,
	}
}

func (svc *Service) FindJobSpec(ctx context.Context, jobName, projectSlug string) (*optimusv1beta1.JobSpecificationResponse, error) {
	optimusCl, err := svc.getOptimusClient(ctx, projectSlug)
	if err != nil {
		return nil, err
	}

	res, err := optimusCl.GetJobSpecifications(ctx, &optimusv1beta1.GetJobSpecificationsRequest{
		ProjectName: projectSlug,
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

func (svc *Service) ListJobs(ctx context.Context, projectSlug string) ([]*optimusv1beta1.JobSpecificationResponse, error) {
	optimusCl, err := svc.getOptimusClient(ctx, projectSlug)
	if err != nil {
		return nil, err
	}

	res, err := optimusCl.GetJobSpecifications(ctx, &optimusv1beta1.GetJobSpecificationsRequest{
		ProjectName: projectSlug,
	})
	if err != nil {
		return nil, err
	}

	return res.JobSpecificationResponses, nil
}

func (svc *Service) getOptimusClient(ctx context.Context, projectSlug string) (optimusv1beta1grpc.JobSpecificationServiceClient, error) {
	// retrieve hostname from cache

	if cl, exists := svc.cache.data[projectSlug]; exists {

		return cl, nil
	} else {
		// retrieve hostname from shield
		prj, err := svc.shieldClient.GetProject(ctx, &shieldv1beta1.GetProjectRequest{
			Id: projectSlug,
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
			return nil, ErrOptimusHostNotString
		}

		cl, err := svc.builder.BuildOptimusClient(ctx, optimusHostStr)
		if err != nil {
			return nil, err
		}

		// store hostname in cache
		svc.cache.mu.Lock()
		svc.cache.data[projectSlug] = cl
		svc.cache.mu.Unlock()

		return cl, nil
	}
}
