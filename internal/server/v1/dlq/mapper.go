package dlq

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

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

func enrichDlqJob(job *models.DlqJob, res *entropyv1beta1.Resource) error {
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

	// TODO: populate prometheus host using firehose's
	promRemoteWriteAny, exists := modConf.Telegraf.Config.Output["prometheus_remote_write"]
	if !exists {
		return errors.New("missing prometheus_remote_write")
	}
	promRemoteWriteMap, ok := promRemoteWriteAny.(map[string]any)
	if !ok {
		return errors.New("invalid prometheus_remote_write")
	}
	promURLAny, exists := promRemoteWriteMap["url"]
	if !exists {
		return errors.New("missing prometheus_remote_write.url")
	}
	promURLString, ok := promURLAny.(string)
	if !ok {
		return errors.New("invalid prometheus_remote_write.url")
	}

	envs := modConf.EnvVariables
	firehoseLabels := res.Labels
	sinkType := envs["SINK_TYPE"]
	dataType := envs["INPUT_SCHEMA_PROTO_CLASS"]
	groupID := firehoseLabels["group"]
	groupSlug := firehoseLabels["team"]
	dlqPrefixDirectory := strings.Replace(envs["DLQ_GCS_DIRECTORY_PREFIX"], "{{ .name }}", res.Name, 1)
	metricStatsDTag := fmt.Sprintf(
		"namespace=%s,app=%s,sink=%s,team=%s,proto=%s,firehose=%s",
		namespace,
		job.Name,
		sinkType,
		groupSlug,
		dataType,
		res.Urn,
	)

	job.PrometheusHost = promURLString
	job.Namespace = namespace
	job.Project = res.Project
	job.KubeCluster = kubeCluster
	job.Group = groupID
	job.Team = groupSlug
	job.DlqGcsCredentialPath = envs["DLQ_GCS_CREDENTIAL_PATH"]
	if job.DlqGcsCredentialPath == "" {
		return errors.New("missing DLQ_GCS_CREDENTIAL_PATH")
	}
	job.EnvVars = map[string]string{
		"DLQ_BATCH_SIZE":                          fmt.Sprintf("%d", job.BatchSize),
		"DLQ_NUM_THREADS":                         fmt.Sprintf("%d", job.NumThreads),
		"DLQ_ERROR_TYPES":                         job.ErrorTypes,
		"DLQ_INPUT_DATE":                          job.Date,
		"DLQ_TOPIC_NAME":                          job.Topic,
		"METRIC_STATSD_TAGS":                      metricStatsDTag,
		"SINK_TYPE":                               sinkType,
		"DLQ_PREFIX_DIR":                          dlqPrefixDirectory,
		"DLQ_FINISHED_STATUS_FILE":                "/shared/job-finished",
		"DLQ_GCS_BUCKET_NAME":                     envs["DLQ_GCS_BUCKET_NAME"],
		"DLQ_GCS_CREDENTIAL_PATH":                 job.DlqGcsCredentialPath,
		"DLQ_GCS_GOOGLE_CLOUD_PROJECT_ID":         envs["DLQ_GCS_GOOGLE_CLOUD_PROJECT_ID"],
		"JAVA_TOOL_OPTIONS":                       envs["JAVA_TOOL_OPTIONS"],
		"_JAVA_OPTIONS":                           envs["_JAVA_OPTIONS"],
		"INPUT_SCHEMA_PROTO_CLASS":                dataType,
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

	return nil
}

// DlqJob param here is expected to have been enriched with firehose config
func mapToEntropyResource(job models.DlqJob) (*entropyv1beta1.Resource, error) {
	cfgStruct, err := makeConfigStruct(job)
	if err != nil {
		return nil, err
	}

	return &entropyv1beta1.Resource{
		Urn:       job.Urn,
		Kind:      entropy.ResourceKindJob,
		Name:      job.Name,
		Project:   job.Project,
		Labels:    buildResourceLabels(job),
		CreatedBy: job.CreatedBy,
		UpdatedBy: job.UpdatedBy,
		Spec: &entropyv1beta1.ResourceSpec{
			Configs: cfgStruct,
			Dependencies: []*entropyv1beta1.ResourceDependency{
				{
					Key:   "kube_cluster",
					Value: job.KubeCluster,
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
					"APP_NAME":        job.ResourceID,
					"PROMETHEUS_HOST": job.PrometheusHost,
					"TEAM":            job.Team,
					"TOPIC":           job.Topic,
					"environment":     "production",
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
			"topic": job.Topic,
			"date":  job.Date,
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

func mapToDlqJob(r *entropyv1beta1.Resource) (models.DlqJob, error) {
	labels := r.Labels

	var envVars map[string]string
	if r.GetSpec().Configs != nil {
		var modConf entropy.JobConfig
		if err := utils.ProtoStructToGoVal(r.GetSpec().GetConfigs(), &modConf); err != nil {
			return models.DlqJob{}, fmt.Errorf("error parsing proto value: %w", err)
		}

		if len(modConf.Containers) > 0 {
			envVars = modConf.Containers[0].EnvVariables
		}
	}

	var kubeCluster string
	for _, dep := range r.Spec.GetDependencies() {
		if dep.GetKey() == "kube_cluster" {
			kubeCluster = dep.GetValue()
		}
	}

	batchSize, err := strconv.Atoi(labels["batch_size"])
	if err != nil {
		batchSize = 0
	}

	numThreads, err := strconv.Atoi(labels["num_threads"])
	if err != nil {
		numThreads = 0
	}

	replicas, err := strconv.Atoi(labels["replicas"])
	if err != nil {
		replicas = 0
	}

	job := models.DlqJob{
		Urn:                  r.Urn,
		Name:                 r.Name,
		ResourceID:           labels["resource_id"],
		ResourceType:         labels["resource_type"],
		Date:                 labels["date"],
		Topic:                labels["topic"],
		PrometheusHost:       labels["prometheus_host"],
		Group:                labels["group"],
		Namespace:            labels["namespace"],
		ContainerImage:       labels["container_image"],
		ErrorTypes:           labels["error_types"],
		BatchSize:            int64(batchSize),
		NumThreads:           int64(numThreads),
		Replicas:             int64(replicas),
		KubeCluster:          kubeCluster,
		Project:              r.Project,
		CreatedBy:            r.CreatedBy,
		UpdatedBy:            r.UpdatedBy,
		Status:               r.GetState().GetStatus().String(),
		CreatedAt:            strfmt.DateTime(r.CreatedAt.AsTime()),
		UpdatedAt:            strfmt.DateTime(r.UpdatedAt.AsTime()),
		EnvVars:              envVars,
		DlqGcsCredentialPath: labels["dlq_gcs_credential_path"],
	}

	return job, nil
}

func buildResourceLabels(job models.DlqJob) map[string]string {
	return map[string]string{
		"resource_id":             job.ResourceID,
		"resource_type":           job.ResourceType,
		"date":                    job.Date,
		"topic":                   job.Topic,
		"job_type":                "dlq",
		"group":                   job.Group,
		"prometheus_host":         job.PrometheusHost,
		"replicas":                fmt.Sprintf("%d", job.Replicas),
		"batch_size":              fmt.Sprintf("%d", job.BatchSize),
		"num_threads":             fmt.Sprintf("%d", job.NumThreads),
		"error_types":             job.ErrorTypes,
		"dlq_gcs_credential_path": job.DlqGcsCredentialPath,
		"container_image":         job.ContainerImage,
		"namespace":               job.Namespace,
	}
}
