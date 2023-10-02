package dlq

import (
	"fmt"

	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/goto/dex/entropy"
	"github.com/goto/dex/generated/models"
	"github.com/goto/dex/internal/server/utils"
)

// from app config
// from firehose
// hardcoded
// user

func mapToEntropyResource(job *models.DlqJob) (*entropyv1beta1.Resource, error) {
	cfgStruct, err := makeConfigStruct(job)
	if err != nil {
		return nil, err
	}

	return &entropyv1beta1.Resource{
		Urn:     job.Urn,
		Kind:    entropy.ResourceKindJob,
		Name:    buildEntropyResourceName(job),
		Project: job.Project,
		Labels:  buildResourceLabels(job),
		Spec: &entropyv1beta1.ResourceSpec{
			Configs: cfgStruct,
			Dependencies: []*entropyv1beta1.ResourceDependency{
				{
					Key:   "kube_cluster",
					Value: "", // from firehose
				},
			},
		},
	}, nil
}

func makeConfigStruct(job *models.DlqJob) (*structpb.Value, error) {
	return utils.GoValToProtoStruct(entropy.JobConfig{
		Stopped:   job.Stopped,
		Replicas:  1,
		Namespace: "", // same with firehose deployment namespace
		Name:      buildJobName(job),
		Containers: []entropy.JobContainer{
			{
				Name:            "dlq-job",
				Image:           "asia.gcr.io/systems-0001/dlq-processor:0.1.0", // from app config
				ImagePullPolicy: "IfNotPresent",
				SecretsVolumes: []entropy.JobSecret{ // confirm with streaming for required secrets
					{
						Name:  "firehose-bigquery-sink-credential", // from firehose config
						Mount: "/etc/secret/gcp",                   // from firehose config
					},
				},
				Limits: &entropy.UsageSpec{
					CPU:    "0.5", // user
					Memory: "2gb", // user
				},
				Requests: &entropy.UsageSpec{
					CPU:    "0.5", // user
					Memory: "2gb", // user
				},
				EnvConfigMaps: []string{},
				EnvVariables: map[string]string{
					// all firehose ENV_VARS +
					"DLQ_BATCH_SIZE":     "100",                // user
					"DLQ_NUM_THREADS":    "1",                  // user
					"DLQ_ERROR_TYPES":    "DEFAULT_ERROR",      // user
					"DLQ_INPUT_DATE":     "2023-04-10",         // user (internally created)
					"DLQ_TOPIC_NAME":     "gofood-booking-log", // user
					"METRIC_STATSD_TAGS": "TBA",
				},
			},
			{
				Name:            "telegraf",
				Image:           "telegraf:1.18.0-alpine",
				ImagePullPolicy: "IfNotPresent",
				ConfigMapsVolumes: []entropy.JobConfigMap{ // confirm with streaming for required secrets
					{
						Name:  "dlq-processor-telegraf",
						Mount: "/etc/telegraf",
					},
				},
				EnvConfigMaps: []string{},
				EnvVariables: map[string]string{
					// To be updated by streaming
					"APP_NAME":        "", //
					"PROMETHEUS_HOST": "", // from app config
				},
			},
		},
		JobLabels: map[string]string{
			"firehose_deployment": "", // firehose deployment
			"topic":               "",
			"date":                "",
		},
		Volumes: []entropy.JobVolume{
			{
				Name: "firehose-bigquery-sink-credential", // make sure it is similar to how we mount it
				Kind: "secret",
			},
			{
				Name: "dlq-processor-telegraf", // make sure it is similar to how we mount it
				Kind: "configMap",
			},
		},
	})
}

func buildResourceLabels(job *models.DlqJob) map[string]string {
	return map[string]string{
		"firehose": job.ResourceID,
		"date":     job.Date,
		"topic":    job.Topic,
		"job_type": "dlq",
	}
}

func buildEntropyResourceName(job *models.DlqJob) string {
	return fmt.Sprintf(
		"%s-%s-%s-%s",
		job.ResourceID,   // firehose urn
		job.ResourceType, // firehose / dagger
		job.Topic,        //
		job.Date,         //
	)
}

func buildJobName(job *models.DlqJob) string {
	randomKey := "91238192"
	return fmt.Sprintf(
		"%s-%s-%s",
		job.Date,
		job.Topic,
		randomKey,
	)
}
