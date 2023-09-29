package dlq

import (
	"fmt"

	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/goto/dex/entropy"
	"github.com/goto/dex/generated/models"
	"github.com/goto/dex/internal/server/utils"
)

func mapToEntropyResource(job *models.DlqJob) (*entropyv1beta1.Resource, error) {
	cfgStruct, err := makeConfigStruct(job)
	if err != nil {
		return nil, err
	}

	return &entropyv1beta1.Resource{
		Urn:     job.Urn,
		Kind:    entropy.ResourceKindJob,
		Name:    buildResourceName(job),
		Project: job.Project,
		Labels:  buildResourceLabels(job),
		Spec: &entropyv1beta1.ResourceSpec{
			Configs: cfgStruct,
		},
	}, nil
}

func makeConfigStruct(job *models.DlqJob) (*structpb.Value, error) {
	return utils.GoValToProtoStruct(entropy.JobConfig{
		Stopped:   job.Stopped,
		Replicas:  int32(job.Replicas),
		Namespace: "", // decide with streaming where to store namespace
		Name:      "", // decide with streaming job name
		Containers: []entropy.JobContainer{
			{
				Name:            "dlq",
				Image:           "test-image", // ask streaming for image name, and also if it can be modified via UI/client
				ImagePullPolicy: "Always",
				Command:         []string{""}, // confirm with streaming team for commands
				SecretsVolumes: []entropy.JobSecret{ // confirm with streaming for required secrets
					{
						Name:  "",
						Mount: "",
					},
				},
				ConfigMapsVolumes: []entropy.JobConfigMap{ // confirm with streaming for required secrets
					{
						Name:  "",
						Mount: "",
					},
				},
				Limits: &entropy.UsageSpec{
					CPU:    "",
					Memory: "",
				},
				Requests: &entropy.UsageSpec{
					CPU:    "",
					Memory: "",
				},
				EnvConfigMaps: []string{},
				EnvVariables:  map[string]string{},
			},
		},
		JobLabels: map[string]string{},
		Volumes: []entropy.JobVolume{
			{
				Name: "",
				Kind: "",
			},
		},
	})
}

func buildResourceLabels(job *models.DlqJob) map[string]string {
	return map[string]string{
		"batch_size":    fmt.Sprintf("%d", job.BatchSize),
		"blob_batch":    fmt.Sprintf("%d", job.BlobBatch),
		"num_threads":   job.ErrorTypes,
		"date":          job.Date,
		"resource_id":   job.ResourceID,
		"resource_type": job.ResourceType,
		"topic":         job.Topic,
	}
}

func buildResourceName(job *models.DlqJob) string {
	return fmt.Sprintf(
		"%s-%s-%s-%s",
		job.ResourceID,
		job.ResourceType,
		job.Topic,
		job.Date,
	)
}
