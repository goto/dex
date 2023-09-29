package entropy

type JobConfig struct {
	Stopped    bool              `json:"stopped,omitempty"`
	Replicas   int32             `json:"replicas"`
	Namespace  string            `json:"namespace,omitempty"`
	Name       string            `json:"name,omitempty"`
	Containers []JobContainer    `json:"containers,omitempty"`
	JobLabels  map[string]string `json:"job_labels,omitempty"`
	Volumes    []JobVolume       `json:"volumes,omitempty"`
}

type JobVolume struct {
	Name string
	Kind string // secret or config-map. secret is for gcs/bq credential
}

type JobContainer struct {
	Name              string            `json:"name"`
	Image             string            `json:"image"`
	ImagePullPolicy   string            `json:"image_pull_policy,omitempty"`
	Command           []string          `json:"command,omitempty"`
	SecretsVolumes    []JobSecret       `json:"secrets_volumes,omitempty"`
	ConfigMapsVolumes []JobConfigMap    `json:"config_maps_volumes,omitempty"`
	Limits            *UsageSpec        `json:"limits,omitempty"`
	Requests          *UsageSpec        `json:"requests,omitempty"`
	EnvConfigMaps     []string          `json:"env_config_maps,omitempty"`
	EnvVariables      map[string]string `json:"env_variables,omitempty"`
}

type JobSecret struct {
	Name  string `json:"name"`
	Mount string `json:"mount"`
}

type JobConfigMap struct {
	Name  string `json:"name"`
	Mount string `json:"mount"`
}
