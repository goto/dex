package dlq

import (
	"context"

	entropyv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/entropy/v1beta1/entropyv1beta1grpc"
	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"

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
	cfg       *DlqJobConfig
}

func NewService(client entropyv1beta1rpc.ResourceServiceClient, gcsClient gcs.BlobStorageClient, cfg *DlqJobConfig) *Service {
	return &Service{
		client:    client,
		gcsClient: gcsClient,
		cfg:       cfg,
	}
}

func (s *Service) getDlqJob(ctx context.Context, jobURN string) (*models.DlqJob, error) {
	res, err := s.client.GetResource(ctx, &entropyv1beta1.GetResourceRequest{Urn: jobURN})
	if err != nil {
		return nil, err
	}

	dlqJob, err := MapToDlqJob(res.GetResource())
	if err != nil {
		return nil, err
	}
	return dlqJob, nil
}
