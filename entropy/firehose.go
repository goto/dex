package entropy

import "time"

type FirehoseConfig struct {
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

	Limits        UsageSpec             `json:"limits,omitempty"`
	Requests      UsageSpec             `json:"requests,omitempty"`
	Telegraf      *FirehoseTelegraf     `json:"telegraf,omitempty"`
	ChartValues   *FirehoseChartValues  `json:"chart_values,omitempty"`
	InitContainer FirehoseInitContainer `json:"init_container,omitempty"`
}

type FirehoseInitContainer struct {
	Enabled bool `json:"enabled"`

	Args    []string `json:"args"`
	Command []string `json:"command"`

	Repository string `json:"repository"`
	ImageTag   string `json:"image_tag"`
	PullPolicy string `json:"pull_policy"`
}

type FirehoseTelegraf struct {
	Enabled bool                 `json:"enabled,omitempty"`
	Image   map[string]any       `json:"image,omitempty"`
	Config  FirehoseTelegrafConf `json:"config,omitempty"`
}

type FirehoseTelegrafConf struct {
	Output               map[string]any    `json:"output"`
	AdditionalGlobalTags map[string]string `json:"additional_global_tags"`
}

type FirehoseChartValues struct {
	ImageTag        string `json:"image_tag" validate:"required"`
	ChartVersion    string `json:"chart_version" validate:"required"`
	ImagePullPolicy string `json:"image_pull_policy" validate:"required"`
}

type FirehoseScaleParams struct {
	Replicas int `json:"replicas"`
}

type FirehoseStartParams struct {
	StopTime *time.Time `json:"stop_time"`
}

type FirehoseResetParams struct {
	To string `json:"to"`
}
