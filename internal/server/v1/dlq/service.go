package dlq

import (
	entropyv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/entropy/v1beta1/entropyv1beta1grpc"

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

// TODO: replace *DlqJob with a generated models.DlqJob
func (s *Service) CreateDLQJob(dlqJob *DlqJob) error {
	// validate dlqJob for creation

	// fetch firehose details
	// enrich DlqJob with firehose details
	// map DlqJob to entropy resource -> return entropy.Resource (kind = job)
	// entropy create resource
	return nil
}
