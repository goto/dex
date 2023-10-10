package dlq

import (
	"fmt"
	"strconv"
	"time"

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
	EnvVars map[string]string `json:"env_vars,omitempty"`

	KubeCluster string `json:"kube_cluster,omitempty"`

	FirehoseDeployment string `json:"firehose_deployment,omitempty"`

	ContainerImage string `json:"container_image,omitempty"`

	PrometheusHost string `json:"prometheus_host,omitempty"`

	DlqGcsCredentialPath string `json:"dlq_gcs_credential_path,omitempty"`
}

func enrichDlqJob(job *DlqJob, f models.Firehose, cfg *DlqJobConfig) error {

	job = &DlqJob{
		ResourceID:           f.Urn,
		ResourceType:         "firehose",                      // this can be firehose/dagger
		Date:                 time.Now().Format("2006-01-02"), // taking time when this is created still look for further changes
		FirehoseDeployment:   f.Configs.DeploymentID,
		KubeCluster:          *f.Configs.KubeCluster,
		ContainerImage:       cfg.DlqJobImage,
		PrometheusHost:       cfg.PrometheusHost,
		EnvVars:              f.Configs.EnvVars,
		DlqGcsCredentialPath: f.Configs.EnvVars["DLQ_GCS_CREDENTIAL_PATH"],
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
		BlobBatch:    0, // Need to find the value
		CreatedAt:    strfmt.DateTime(r.CreatedAt.AsTime()),
		CreatedBy:    "", // Need to find the value
		Project:      r.Project,
		ResourceID:   resourceID,
		ResourceType: resourceType,
		Status:       "",                // Need to find the value
		Stopped:      false,             // Need to find the value
		UpdatedAt:    strfmt.DateTime{}, // need to find the value
		Urn:          r.Urn,
	}

	mapDlqSpecAndLabels(&job, r.GetSpec())

	return &job, nil
}

func mapDlqSpecAndLabels(job *DlqJob, spec *entropyv1beta1.ResourceSpec) error {
	var modConf entropy.JobConfig
	if err := utils.ProtoStructToGoVal(spec.GetConfigs(), &modConf); err != nil {
		return err
	}

	var kubeCluster string
	for _, dep := range spec.GetDependencies() {
		if dep.GetKey() == "kube_cluster" {
			kubeCluster = dep.GetValue()
		}
	}

	job.BatchSize, _ = strconv.ParseInt(modConf.Containers[0].EnvVariables["DLQ_BATCH_SIZE"], 10, 64)
	job.ErrorTypes = modConf.Containers[0].EnvVariables["DLQ_ERROR_TYPES"]
	job.NumThreads, _ = strconv.ParseInt(modConf.Containers[0].EnvVariables["DLQ_NUM_THREADS"], 10, 64)
	job.Date = modConf.Containers[0].EnvVariables["DLQ_INPUT_DATE"]
	job.Topic = modConf.Containers[0].EnvVariables["DLQ_TOPIC_NAME"]
	job.KubeCluster = kubeCluster
	job.FirehoseDeployment = modConf.Namespace
	job.Replicas = int64(modConf.Replicas)

	return nil
}

func makeConfigStruct(job DlqJob) (*structpb.Value, error) {
	EnvVars := map[string]string{
		"DLQ_BATCH_SIZE":     strconv.FormatInt(job.BatchSize, 10),  // user
		"DLQ_NUM_THREADS":    strconv.FormatInt(job.NumThreads, 10), // user
		"DLQ_ERROR_TYPES":    job.ErrorTypes,                        // user
		"DLQ_INPUT_DATE":     job.Date,                              // user (internally created)
		"DLQ_TOPIC_NAME":     job.Topic,                             // user
		"METRIC_STATSD_TAGS": "TBA",
	}
	for key, value := range job.EnvVars {
		EnvVars[key] = value
	}
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
						Mount: job.DlqGcsCredentialPath,            // from firehose configs.DLQ_GCS_CREDENTIAL_PATH
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
				EnvVariables: EnvVars,
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
			"firehose_deployment": job.FirehoseDeployment, // firehose deployment
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
