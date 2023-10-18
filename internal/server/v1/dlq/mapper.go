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

const (
	kubeClusterDependenciesKey = "kube_cluster"
	//nolint
	dlqCredentialsSecretName = "firehose-bigquery-sink-credential"
	dlqTelegrafConfigName    = "dlq-processor-telegraf"
)

func EnrichDlqJob(job *models.DlqJob, res *entropyv1beta1.Resource, cfg *DlqJobConfig) error {
	var kubeCluster string
	for _, dep := range res.Spec.GetDependencies() {
		if dep.GetKey() == kubeClusterDependenciesKey {
			kubeCluster = dep.GetValue()
		}
	}
	var modConf entropy.FirehoseConfig
	if err := utils.ProtoStructToGoVal(res.Spec.GetConfigs(), &modConf); err != nil {
		return err
	}

	outputMap := res.GetState().GetOutput().GetStructValue().AsMap()
	namespaceAny, exists := outputMap["namespace"]
	if !exists {
		return ErrFirehoseNamespaceNotFound
	}
	namespace, ok := namespaceAny.(string)
	if !ok {
		return ErrFirehoseNamespaceInvalid
	}

	envs := modConf.EnvVariables
	job.ResourceID = res.GetUrn()
	job.Namespace = namespace
	job.KubeCluster = kubeCluster
	job.ContainerImage = cfg.DlqJobImage
	job.PrometheusHost = cfg.PrometheusHost
	job.EnvVars = map[string]string{
		"DLQ_BATCH_SIZE":                          fmt.Sprintf("%d", job.BatchSize),
		"DLQ_NUM_THREADS":                         fmt.Sprintf("%d", job.NumThreads),
		"DLQ_ERROR_TYPES":                         job.ErrorTypes,
		"DLQ_INPUT_DATE":                          job.Date,
		"DLQ_TOPIC_NAME":                          job.Topic,
		"METRIC_STATSD_TAGS":                      "a=b", // TBA
		"SINK_TYPE":                               envs["SINK_TYPE"],
		"DLQ_PREFIX_DIR":                          "test-firehose",
		"DLQ_FINISHED_STATUS_FILE":                "/shared/job-finished",
		"DLQ_GCS_BUCKET_NAME":                     envs["DLQ_GCS_BUCKET_NAME"],
		"DLQ_GCS_CREDENTIAL_PATH":                 envs["DLQ_GCS_CREDENTIAL_PATH"],
		"DLQ_GCS_GOOGLE_CLOUD_PROJECT_ID":         envs["DLQ_GCS_GOOGLE_CLOUD_PROJECT_ID"],
		"JAVA_TOOL_OPTIONS":                       envs["JAVA_TOOL_OPTIONS"],
		"_JAVA_OPTIONS":                           envs["_JAVA_OPTIONS"],
		"INPUT_SCHEMA_PROTO_CLASS":                envs["INPUT_SCHEMA_PROTO_CLASS"],
		"SCHEMA_REGISTRY_STENCIL_ENABLE":          envs["SCHEMA_REGISTRY_STENCIL_ENABLE"],
		"SCHEMA_REGISTRY_STENCIL_URLS":            envs["SCHEMA_REGISTRY_STENCIL_URLS"],
		"SINK_BIGQUERY_ADD_METADATA_ENABLED":      envs["SINK_BIGQUERY_ADD_METADATA_ENABLED"],
		"SINK_BIGQUERY_CLIENT_CONNECT_TIMEOUT_MS": envs["SINK_BIGQUERY_CLIENT_CONNECT_TIMEOUT_MS"],
		"SINK_BIGQUERY_CLIENT_READ_TIMEOUT_MS":    envs["SINK_BIGQUERY_CLIENT_READ_TIMEOUT_MS"],
		"SINK_BIGQUERY_CREDENTIAL_PATH":           envs["SINK_BIGQUERY_CREDENTIAL_PATH"],
		"SINK_BIGQUERY_DATASET_LABELS":            envs["SINK_BIGQUERY_DATASET_LABELS"],
		"SINK_BIGQUERY_DATASET_LOCATION":          envs["SINK_BIGQUERY_DATASET_LOCATION"],
		"SINK_BIGQUERY_DATASET_NAME":              envs["SINK_BIGQUERY_DATASET_NAME"],
		"SINK_BIGQUERY_GOOGLE_CLOUD_PROJECT_ID":   envs["SINK_BIGQUERY_GOOGLE_CLOUD_PROJECT_ID"],
		"SINK_BIGQUERY_ROW_INSERT_ID_ENABLE":      envs["SINK_BIGQUERY_ROW_INSERT_ID_ENABLE"],
		"SINK_BIGQUERY_STORAGE_API_ENABLE":        envs["SINK_BIGQUERY_STORAGE_API_ENABLE"],
		"SINK_BIGQUERY_TABLE_LABELS":              envs["SINK_BIGQUERY_TABLE_LABELS"],
		"SINK_BIGQUERY_TABLE_NAME":                envs["SINK_BIGQUERY_TABLE_NAME"],
		"SINK_BIGQUERY_TABLE_PARTITION_EXPIRY_MS": envs["SINK_BIGQUERY_TABLE_PARTITION_EXPIRY_MS"],
		"SINK_BIGQUERY_TABLE_PARTITION_KEY":       envs["SINK_BIGQUERY_TABLE_PARTITION_KEY"],
		"SINK_BIGQUERY_TABLE_PARTITIONING_ENABLE": envs["SINK_BIGQUERY_TABLE_PARTITIONING_ENABLE"],
	}
	job.DlqGcsCredentialPath = modConf.EnvVariables["DLQ_GCS_CREDENTIAL_PATH"]

	return nil
}

// DlqJob param here is expected to have been enriched with firehose config
func MapToEntropyResource(job models.DlqJob) (*entropyv1beta1.Resource, error) {
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

func makeConfigStruct(job models.DlqJob) (*structpb.Value, error) {
	return utils.GoValToProtoStruct(entropy.JobConfig{
		Replicas:  int32(job.Replicas),
		Namespace: job.Namespace,
		Containers: []entropy.JobContainer{
			{
				Name:            "dlq-job",
				Image:           job.ContainerImage,
				ImagePullPolicy: "Always",
				SecretsVolumes: []entropy.JobSecret{
					{
						Name:  dlqCredentialsSecretName,
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
				EnvVariables: job.EnvVars,
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
				Name: dlqCredentialsSecretName,
				Kind: "secret",
			},
			{
				Name: dlqTelegrafConfigName,
				Kind: "configMap",
			},
		},
	})
}

func MapToDlqJob(r *entropyv1beta1.Resource) (*models.DlqJob, error) {
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
	batchSize, err := strconv.ParseInt(envVars["DLQ_BATCH_SIZE"], 10, 64)
	if err != nil {
		return nil, err
	}
	numThreads, err := strconv.ParseInt(envVars["DLQ_NUM_THREADS"], 10, 64)
	if err != nil {
		return nil, err
	}
	errorTypes := envVars["DLQ_ERROR_TYPES"]

	job := models.DlqJob{
		Urn:                  r.Urn,
		Name:                 r.Name,
		ResourceID:           labels["resource_id"],
		ResourceType:         labels["resource_type"],
		Date:                 labels["date"],
		Topic:                labels["topic"],
		PrometheusHost:       labels["prometheus_host"],
		Namespace:            modConf.Namespace,
		ContainerImage:       modConf.Containers[0].Image,
		ErrorTypes:           errorTypes,
		BatchSize:            batchSize,
		NumThreads:           numThreads,
		Replicas:             int64(modConf.Replicas),
		KubeCluster:          kubeCluster,
		Project:              r.Project,
		CreatedBy:            r.CreatedBy,
		UpdatedBy:            r.UpdatedBy,
		Status:               string(r.GetState().Status),
		CreatedAt:            strfmt.DateTime(r.CreatedAt.AsTime()),
		UpdatedAt:            strfmt.DateTime(r.UpdatedAt.AsTime()),
		EnvVars:              envVars,
		DlqGcsCredentialPath: envVars["DLQ_GCS_CREDENTIAL_PATH"],
	}

	return &job, nil
}

func buildResourceLabels(job models.DlqJob) map[string]string {
	return map[string]string{
		"resource_id": job.ResourceID,
		"type":        job.ResourceType,
		"date":        job.Date,
		"topic":       job.Topic,
		"job_type":    "dlq",
	}
}

func buildEntropyResourceName(job models.DlqJob) string {
	return fmt.Sprintf(
		"%s-%s-%s-%s",
		job.ResourceID,   // firehose urn
		job.ResourceType, // firehose / dagger
		job.Topic,        //
		job.Date,         //
	)
}
