package dlq

import (
	"context"
	"fmt"
	"time"

	entropyv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/entropy/v1beta1/entropyv1beta1grpc"
	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/goto/dex/entropy"
	"github.com/goto/dex/generated/models"
	"github.com/goto/dex/internal/server/gcs"
)

type DlqJobConfig struct {
	JobImage string
}

type Service struct {
	client    entropyv1beta1rpc.ResourceServiceClient
	gcsClient gcs.BlobStorageClient
	cfg       DlqJobConfig
}

func NewService(client entropyv1beta1rpc.ResourceServiceClient, gcsClient gcs.BlobStorageClient, cfg DlqJobConfig) *Service {
	return &Service{
		client:    client,
		gcsClient: gcsClient,
		cfg:       cfg,
	}
}

// TODO: replace *DlqJob with a generated models.DlqJob
func (s *Service) CreateDLQJob(ctx context.Context, userEmail string, dlqJob models.DlqJob) (models.DlqJob, error) {
	if s.cfg.JobImage == "" {
		return models.DlqJob{}, ErrEmptyConfigImage
	}

	timestamp := time.Now().Unix()
	dlqJob.Name = fmt.Sprintf("dlq-%s-%d", dlqJob.Date, timestamp)
	dlqJob.Replicas = 1
	dlqJob.ContainerImage = s.cfg.JobImage

	def, err := s.client.GetResource(ctx, &entropyv1beta1.GetResourceRequest{Urn: dlqJob.ResourceID})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			return models.DlqJob{}, ErrFirehoseNotFound
		}
		return models.DlqJob{}, fmt.Errorf("error getting firehose resource: %w", err)
	}
	// enrich DlqJob with firehose details
	if err := enrichDlqJob(&dlqJob, def.GetResource()); err != nil {
		return models.DlqJob{}, fmt.Errorf("error enriching dlq job: %w", err)
	}

	// map DlqJob to entropy resource -> return entropy.Resource (kind = job)
	resource, err := mapToEntropyResource(dlqJob)
	if err != nil {
		return models.DlqJob{}, fmt.Errorf("error mapping to entropy resource: %w", err)
	}
	// entropy create resource
	ctx = metadata.AppendToOutgoingContext(ctx, "user-id", userEmail)
	req := &entropyv1beta1.CreateResourceRequest{Resource: resource}
	res, err := s.client.CreateResource(ctx, req)
	if err != nil {
		return models.DlqJob{}, fmt.Errorf("error creating resource: %w", err)
	}

	updatedDlqJob, err := mapToDlqJob(res.GetResource())
	if err != nil {
		return models.DlqJob{}, fmt.Errorf("error mapping back to dlq job: %w", err)
	}

	return updatedDlqJob, nil
}

func (s *Service) ListDlqJob(ctx context.Context, labelFilter map[string]string) ([]models.DlqJob, error) {
	dlqJob := []models.DlqJob{}

	rpcReq := &entropyv1beta1.ListResourcesRequest{
		Kind:   entropy.ResourceKindJob,
		Labels: labelFilter,
	}

	rpcResp, err := s.client.ListResources(ctx, rpcReq)
	if err != nil {
		return nil, fmt.Errorf("error getting job resource list: %w", err)
	}
	for _, res := range rpcResp.GetResources() {
		def, err := mapToDlqJob(res)
		if err != nil {
			return nil, fmt.Errorf("error mapping to dlq job: %w", err)
		}
		dlqJob = append(dlqJob, def)
	}

	return dlqJob, nil
}

func (s *Service) GetDlqJob(ctx context.Context, jobURN string) (models.DlqJob, error) {
	res, err := s.client.GetResource(ctx, &entropyv1beta1.GetResourceRequest{Urn: jobURN})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			return models.DlqJob{}, ErrJobNotFound
		}
		return models.DlqJob{}, fmt.Errorf("error getting entropy resource: %w", err)
	}

	dlqJob, err := mapToDlqJob(res.GetResource())
	if err != nil {
		return models.DlqJob{}, fmt.Errorf("error mapping resource to dlq job: %w", err)
	}
	return dlqJob, nil
}
