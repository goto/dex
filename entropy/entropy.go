package entropy

import "time"

type Config struct {
	// Stopped flag when set forces the firehose to be stopped on next sync.
	Stopped bool `json:"stopped"`

	// StopTime can be set to schedule the firehose to be stopped at given time.
	StopTime *time.Time `json:"stop_time,omitempty"`

	// Replicas is the number of firehose instances to run.
	Replicas int `json:"replicas"`

	// Namespace is the target namespace where firehose should be deployed.
	// Inherits from driver config.
	Namespace string `json:"namespace,omitempty"`

	// DeploymentID will be used as the release-name for the deployment.
	// Must be shorter than 53 chars if set. If not set, one will be generated
	// automatically.
	DeploymentID string `json:"deployment_id,omitempty"`

	// EnvVariables contains all the firehose environment config values.
	EnvVariables map[string]string `json:"env_variables,omitempty"`

	// ResetOffset represents the value to which kafka consumer offset was set to
	ResetOffset string `json:"reset_offset,omitempty"`

	Limits        UsageSpec     `json:"limits,omitempty"`
	Requests      UsageSpec     `json:"requests,omitempty"`
	Telegraf      *Telegraf     `json:"telegraf,omitempty"`
	ChartValues   *ChartValues  `json:"chart_values,omitempty"`
	InitContainer InitContainer `json:"init_container,omitempty"`
}

type UsageSpec struct {
	CPU    string `json:"cpu,omitempty" validate:"required"`
	Memory string `json:"memory,omitempty" validate:"required"`
}

type InitContainer struct {
	Enabled bool `json:"enabled"`

	Args    []string `json:"args"`
	Command []string `json:"command"`

	Repository string `json:"repository"`
	ImageTag   string `json:"image_tag"`
	PullPolicy string `json:"pull_policy"`
}

type Telegraf struct {
	Enabled bool           `json:"enabled,omitempty"`
	Image   map[string]any `json:"image,omitempty"`
	Config  TelegrafConf   `json:"config,omitempty"`
}

type TelegrafConf struct {
	Output               map[string]any    `json:"output"`
	AdditionalGlobalTags map[string]string `json:"additional_global_tags"`
}

type ChartValues struct {
	ImageTag        string `json:"image_tag" validate:"required"`
	ChartVersion    string `json:"chart_version" validate:"required"`
	ImagePullPolicy string `json:"image_pull_policy" validate:"required"`
}

type ScaleParams struct {
	Replicas int `json:"replicas"`
}

type StartParams struct {
	StopTime *time.Time `json:"stop_time"`
}

type ResetParams struct {
	To string `json:"to"`
}
