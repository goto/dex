package dlq

import (
	"context"

	entropyv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/entropy/v1beta1/entropyv1beta1grpc"
	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"google.golang.org/grpc/metadata"

	"github.com/goto/dex/generated/models"
	"github.com/goto/dex/internal/server/gcs"
)

type DlqJobConfig struct {
	DlqJobImage    string `mapstructure:"dlq_job_image"`
	PrometheusHost string `mapstructure:"prometheus_host"`
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
func (s *Service) CreateDLQJob(ctx context.Context, userEmail string, dlqJob *models.DlqJob) error {
	// validate dlqJob for creation
	// fetch firehose details
	def, err := s.client.GetResource(ctx, &entropyv1beta1.GetResourceRequest{Urn: dlqJob.ResourceID})
	if err != nil {
		return ErrFirehoseNotFound
	}
	// enrich DlqJob with firehose details
	if err := enrichDlqJob(dlqJob, def.GetResource(), s.cfg); err != nil {
		return ErrFirehoseNotFound
	}

	// map DlqJob to entropy resource -> return entropy.Resource (kind = job)
	res, err := mapToEntropyResource(*dlqJob)
	if err != nil {
		return err
	}
	// entropy create resource
	entropyCtx := metadata.AppendToOutgoingContext(ctx, "user-id", userEmail)
	rpcReq := &entropyv1beta1.CreateResourceRequest{Resource: res}
	rpcResp, err := s.client.CreateResource(entropyCtx, rpcReq)
	dlqJob.Urn = rpcResp.Resource.Urn
	if err != nil {
		outErr := ErrInternal
		return outErr
	}

	return nil
}

func (s *Service) getDlqJob(ctx context.Context, firehoseUrn string) (*models.DlqJob, error) {
	dlqJob := &models.DlqJob{}

	def, err := s.client.GetResource(ctx, &entropyv1beta1.GetResourceRequest{Urn: firehoseUrn})
	if err != nil {
		return nil, ErrFirehoseNotFound
	}

	dlqJob, err = MapToDlqJob(def.GetResource())

	return dlqJob, nil
}
