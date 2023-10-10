package dlq

import (
	"fmt"
	"strconv"

	entropyv1beta1 "buf.build/gen/go/gotocompany/proton/protocolbuffers/go/gotocompany/entropy/v1beta1"
	"github.com/go-openapi/strfmt"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/goto/dex/entropy"
	"github.com/goto/dex/generated/models"
	"github.com/goto/dex/internal/server/utils"
)

func enrichDlqJob(job *DlqJob, f models.Firehose, cfg *DlqJobConfig) {
	job.ResourceID = f.Urn
	job.Namespace = "" // TBA
	job.KubeCluster = *f.Configs.KubeCluster
	job.ContainerImage = cfg.DlqJobImage
	job.PrometheusHost = cfg.PrometheusHost
	job.EnvVars = f.Configs.EnvVars
	job.DlqGcsCredentialPath = f.Configs.EnvVars["DLQ_GCS_CREDENTIAL_PATH"]
}

// DlqJob param here is expected to have been enriched with firehose config
func mapToEntropyResource(job DlqJob) (*entropyv1beta1.Resource, error) {
	cfgStruct, err := makeConfigStruct(job)
	if err != nil {
		return nil, err
	}

	return &entropyv1beta1.Resource{
		Urn:       job.Urn,
		Kind:      entropy.ResourceKindJob,
		Name:      buildEntropyResourceName(job),
		Project:   job.Project,
		Labels:    buildResourceLabels(job),
		CreatedBy: job.CreatedBy,
		UpdatedBy: job.UpdatedBy,
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

func makeConfigStruct(job DlqJob) (*structpb.Value, error) {
	envVars := map[string]string{}
	for key, value := range job.EnvVars {
		envVars[key] = value
	}
	envVars["DLQ_BATCH_SIZE"] = fmt.Sprintf("%d", job.BatchSize)
	envVars["DLQ_NUM_THREADS"] = fmt.Sprintf("%d", job.NumThreads)
	envVars["DLQ_ERROR_TYPES"] = job.ErrorTypes
	envVars["DLQ_INPUT_DATE"] = job.Date
	envVars["DLQ_TOPIC_NAME"] = job.Topic
	envVars["METRIC_STATSD_TAGS"] = "" // TBA

	dlqCredentialSecretName := "firehose-bigquery-sink-credential"
	dlqTelegrafConfigName := "dlq-processor-telegraf"
	return utils.GoValToProtoStruct(entropy.JobConfig{
		Replicas:  int32(job.Replicas),
		Namespace: job.Namespace,
		Name:      "", // TBA
		Containers: []entropy.JobContainer{
			{
				Name:            "dlq-job",
				Image:           job.ContainerImage,
				ImagePullPolicy: "Always",
				SecretsVolumes: []entropy.JobSecret{
					{
						Name:  dlqCredentialSecretName,
						Mount: job.DlqGcsCredentialPath,
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
				EnvVariables: envVars,
			},
			{
				Name:  "telegraf",
				Image: "telegraf:1.18.0-alpine",
				ConfigMapsVolumes: []entropy.JobConfigMap{
					{
						Name:  dlqTelegrafConfigName,
						Mount: "/etc/telegraf",
					},
				},
				EnvVariables: map[string]string{
					// To be updated by streaming
					"APP_NAME":        "", // TBA
					"PROMETHEUS_HOST": job.PrometheusHost,
					"DEPLOYMENT_NAME": "deployment-name",
					"TEAM":            job.Group,
					"TOPIC":           job.Topic,
					"environment":     "production", // TBA
					"organization":    "de",         // TBA
					"projectID":       job.Project,
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
			"firehose": job.ResourceID,
			"topic":    job.Topic,
			"date":     job.Date,
		},
		Volumes: []entropy.JobVolume{
			{
				Name: dlqCredentialSecretName,
				Kind: "secret",
			},
			{
				Name: dlqTelegrafConfigName,
				Kind: "configMap",
			},
		},
	})
}

func mapToDlqJob(r *entropyv1beta1.Resource) (*DlqJob, error) {
	labels := r.Labels

	var modConf entropy.JobConfig
	if err := utils.ProtoStructToGoVal(r.Spec.GetConfigs(), &modConf); err != nil {
		return nil, err
	}

	var kubeCluster string
	for _, dep := range r.Spec.GetDependencies() {
		if dep.GetKey() == "kube_cluster" {
			kubeCluster = dep.GetValue()
		}
	}

	envVars := modConf.Containers[0].EnvVariables
	batchSize, _ := strconv.ParseInt(envVars["DLQ_BATCH_SIZE"], 10, 64)   // handler error properly
	numThreads, _ := strconv.ParseInt(envVars["DLQ_NUM_THREADS"], 10, 64) // handler error properly
	errorTypes := envVars["DLQ_ERROR_TYPES"]

	job := DlqJob{
		Urn:          r.Urn,
		ResourceID:   labels["resource_id"],
		ResourceType: labels["resource_type"],
		Date:         labels["date"],
		Topic:        labels["topic"],
		Namespace:    modConf.Namespace,
		ErrorTypes:   errorTypes,
		BatchSize:    batchSize,
		NumThreads:   numThreads,
		Replicas:     int64(modConf.Replicas),
		KubeCluster:  kubeCluster,
		Project:      r.Project,
		CreatedBy:    r.CreatedBy,
		UpdatedBy:    r.UpdatedBy,
		Status:       string(r.GetState().Status),
		CreatedAt:    strfmt.DateTime(r.CreatedAt.AsTime()),
		UpdatedAt:    strfmt.DateTime(r.UpdatedAt.AsTime()),
	}

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
	job.NumThreads, _ = strconv.ParseInt(modConf.Containers[0].EnvVariables["DLQ_NUM_THREADS"], 10, 64)
	job.ErrorTypes = modConf.Containers[0].EnvVariables["DLQ_ERROR_TYPES"]
	job.Replicas = int64(modConf.Replicas)
	job.KubeCluster = kubeCluster
	job.Namespace = modConf.Namespace

	return nil
}

func buildResourceLabels(job DlqJob) map[string]string {
	return map[string]string{
		"resource_id": job.ResourceID,
		"type":        job.ResourceType,
		"date":        job.Date,
		"topic":       job.Topic,
		"job_type":    "dlq",
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
