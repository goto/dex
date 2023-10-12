package dlq

import (
	"context"

	entropyv1beta1rpc "buf.build/gen/go/gotocompany/proton/grpc/go/gotocompany/entropy/v1beta1/entropyv1beta1grpc"
	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"google.golang.org/grpc/metadata"

	"github.com/goto/dex/generated/models"
	"github.com/goto/dex/internal/server/gcs"
	"github.com/goto/dex/internal/server/reqctx"
)

type DlqJobConfig struct {
	DlqJobImage    string `mapstructure:"dlq_job_image"`
	PrometheusHost string `mapstructure:"prometheus_host"`
}

type Service struct {
	client    entropyv1beta1rpc.ResourceServiceClient
	gcsClient gcs.BlobStorageClient
	cfg       *DlqJobConfig
	Entropy   entropyv1beta1rpc.ResourceServiceClient
}

func NewService(client entropyv1beta1rpc.ResourceServiceClient, gcsClient gcs.BlobStorageClient, cfg *DlqJobConfig) *Service {
	return &Service{
		client:    client,
		gcsClient: gcsClient,
		cfg:       cfg,
	}
}

// TODO: replace *DlqJob with a generated models.DlqJob
func (s *Service) CreateDLQJob(ctx context.Context, dlqJob *models.DlqJob) (*entropyv1beta1.Resource, error) {
	// validate dlqJob for creation
	// fetch firehose details
	def, err := s.Entropy.GetResource(ctx, &entropyv1beta1.GetResourceRequest{Urn: dlqJob.ResourceID})
	if err != nil {
		return nil, ErrFirehoseNotFound
	}
	// enrich DlqJob with firehose details
	if err := enrichDlqJob(dlqJob, def.GetResource(), s.cfg); err != nil {
		return nil, ErrFirehoseNotFound
	}

	// map DlqJob to entropy resource -> return entropy.Resource (kind = job)
	res, err := mapToEntropyResource(*dlqJob)
	if err != nil {
		return nil, err
	}
	// entropy create resource
	reqCtx := reqctx.From(ctx)
	entropyCtx := metadata.AppendToOutgoingContext(ctx, "user-id", reqCtx.UserEmail)
	rpcReq := &entropyv1beta1.CreateResourceRequest{Resource: res}
	rpcResp, err := s.Entropy.CreateResource(entropyCtx, rpcReq)
	if err != nil {
		outErr := ErrInternal
		return nil, outErr
	}

	return rpcResp.Resource, nil
}
