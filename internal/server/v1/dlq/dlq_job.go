package dlq

import "github.com/go-openapi/strfmt"

// from app config
// from firehose
// hardcoded
// user
type DlqJob struct {
	// batch size
	BatchSize int64 `json:"batch_size,omitempty"`

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

	Group string `json:"group,omitempty"`

	// status
	// Enum: [pending error running stopped]
	Status string `json:"status,omitempty"`

	// topic
	Topic string `json:"topic,omitempty"`

	// updated at
	// Format: date-time
	UpdatedAt strfmt.DateTime `json:"updated_at,omitempty"`

	// updated by
	UpdatedBy string `json:"updated_by,omitempty"`

	// urn
	Urn string `json:"urn,omitempty"`

	// firehose
	EnvVars map[string]string `json:"env_vars,omitempty"`

	KubeCluster string `json:"kube_cluster,omitempty"`

	Namespace string `json:"firehose_deployment,omitempty"`

	ContainerImage string `json:"container_image,omitempty"`

	PrometheusHost string `json:"prometheus_host,omitempty"`

	DlqGcsCredentialPath string `json:"dlq_gcs_credential_path,omitempty"`
}
