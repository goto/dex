package dlq

import (
	"fmt"

	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/go-openapi/strfmt"
	"github.com/goto/dex/entropy"
	"github.com/goto/dex/generated/models"
	"github.com/goto/dex/internal/server/utils"
)

// from app config
// from firehose
// hardcoded
// user
type DlqJob struct {

	// batch size
	BatchSize int64 `json:"batch_size,omitempty"`

	// blob batch
	BlobBatch int64 `json:"blob_batch,omitempty"`

	// created at
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"created_at,omitempty"`

	// created by
	CreatedBy string `json:"created_by,omitempty"`

	// date
	// Example: 2012-10-30
	Date string `json:"date,omitempty"`

	// List of firehose error types, comma separated
	ErrorTypes string `json:"error_types,omitempty"`

	// num threads
	NumThreads int64 `json:"num_threads,omitempty"`

	// Shield's project slug
	Project string `json:"project,omitempty"`

	// replicas
	Replicas int64 `json:"replicas,omitempty"`

	// resource id
	ResourceID string `json:"resource_id,omitempty"`

	// resource type
	// Enum: [firehose]
	ResourceType string `json:"resource_type,omitempty"`

	// status
	// Enum: [pending error running stopped]
	Status string `json:"status,omitempty"`

	// stopped
	Stopped bool `json:"stopped,omitempty"`

	// topic
	Topic string `json:"topic,omitempty"`

	// updated at
	// Format: date-time
	UpdatedAt strfmt.DateTime `json:"updated_at,omitempty"`

	// updated by
	UpdatedBy string `json:"updated_by,omitempty"`

	// urn
	Urn string `json:"urn,omitempty"`

	//firehose
	EnvVars []string `json:"env_vars,omitempty"`

	KubeCluster string `json:"kube_cluster,omitempty"`

	FirehoseDeployment string `json:"firehose_deployment,omitempty"`

	ContainerImage string `json:"container_image,omitempty"`

	PrometheusHost string `json:"prometheus_host,omitempty"`
}

func enrichDlqJob(job *DlqJob, f models.Firehose, cfg *DlqJobConfig) error {
	var env_vars []string
	for key := range f.Configs.EnvVars {
		env_vars = append(env_vars, key)
	}

	job = &DlqJob{
		FirehoseDeployment: f.Configs.DeploymentID,
		KubeCluster:        *f.Configs.KubeCluster,
		EnvVars:            env_vars,
		ContainerImage:     cfg.DlqJobImage,
		PrometheusHost:     cfg.PrometheusHost,
	}
	return nil
}

// DlqJob param here is expected to have been enriched with firehose config
func mapToEntropyResource(job DlqJob) (*entropyv1beta1.Resource, error) {
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
					Value: job.KubeCluster, // from firehose configs.kube_cluster
				},
			},
		},
	}, nil
}

// TODO:
func mapToDlqJob(r *entropyv1beta1.Resource) (*DlqJob, error) {
	var resourceID, resourceType, topic, date string

	_, err := fmt.Sscanf(r.Name, "%s-%s-%s-%s", &resourceID, &resourceType, &topic, &date)
	if err != nil {
		return nil, err
	}
	job := DlqJob{
		BatchSize:          0, // Need to find the value
		BlobBatch:          0, // Need to find the value
		CreatedAt:          strfmt.DateTime(r.CreatedAt.AsTime()),
		CreatedBy:          "", // Need to find the value
		Date:               date,
		ErrorTypes:         "", // Need to find the value
		NumThreads:         0,  // Need to find the value
		Project:            r.Project,
		Replicas:           1, // Don't know how to access spec.Config.Replicas
		ResourceID:         resourceID,
		ResourceType:       resourceType,
		Status:             "",    // Need to find the value
		Stopped:            false, // Need to find the value
		Topic:              topic,
		UpdatedAt:          strfmt.DateTime{}, // need to find the value
		Urn:                r.Urn,
		KubeCluster:        "", // Don't know how to access r.Spec.Dependencies.Value,
		FirehoseDeployment: "", // Don't know how to access r.Spec.Configs.Namespace,
	}

	return &job, nil
}

func makeConfigStruct(job DlqJob) (*structpb.Value, error) {
	return utils.GoValToProtoStruct(entropy.JobConfig{
		Replicas:  1,
		Namespace: job.FirehoseDeployment, // same with firehose deployment namespace
		Name:      "",
		Containers: []entropy.JobContainer{
			{
				Name:            "dlq-job",
				Image:           job.ContainerImage, // from app config
				ImagePullPolicy: "IfNotPresent",
				SecretsVolumes: []entropy.JobSecret{ // confirm with streaming for required secrets
					{
						Name:  "firehose-bigquery-sink-credential", // hard coded, try look in entropy firehose?
						Mount: "/etc/secret/gcp",                   // from firehose configs.DLQ_GCS_CREDENTIAL_PATH
					},
				},
				Limits: entropy.UsageSpec{
					CPU:    "0.5", // user
					Memory: "2gb", // user
				},
				Requests: entropy.UsageSpec{
					CPU:    "0.5", // user
					Memory: "2gb", // user
				},
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
				Name:  "telegraf",
				Image: "telegraf:1.18.0-alpine",
				ConfigMapsVolumes: []entropy.JobConfigMap{ // confirm with streaming for required secrets
					{
						Name:  "dlq-processor-telegraf",
						Mount: "/etc/telegraf",
					},
				},
				EnvVariables: map[string]string{
					// To be updated by streaming
					"APP_NAME":        "",                 //
					"PROMETHEUS_HOST": job.PrometheusHost, // from app config
					"DEPLOYMENT_NAME": "deployment-name",
					"TEAM":            "de",
					"TOPIC":           "test-topic",
					"environment":     "production",
					"organization":    "de",
					"projectID":       "projectID",
				},
				Command: []string{
					"/bin/bash",
				},
				Args: []string{
					"-c",
					"telegraf & while [ ! -f /shared/job-finished ]; do sleep 5; done; sleep 20 && exit 0",
				},
				Limits: entropy.UsageSpec{
					CPU:    "100m",  // user
					Memory: "300Mi", // user
				},
				Requests: entropy.UsageSpec{
					CPU:    "100m",  // user
					Memory: "300Mi", // user
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

func buildResourceLabels(job DlqJob) map[string]string {
	return map[string]string{
		"firehose": job.ResourceID,
		"date":     job.Date,
		"topic":    job.Topic,
		"job_type": "dlq",
	}
}

func buildEntropyResourceName(job DlqJob) string {
	return fmt.Sprintf(
		"%s-%s-%s-%s",
		job.ResourceID,   // firehose urn
		job.ResourceType, // firehose / dagger
		job.Topic,        //
		job.Date,         //
	)
}

func buildJobName(job DlqJob) string {
	randomKey := "91238192"
	return fmt.Sprintf(
		"%s-%s-%s",
		job.Date,
		job.Topic,
		randomKey,
	)
}
