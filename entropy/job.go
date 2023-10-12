package entropy

type JobConfig struct {
	Replicas   int32             `json:"replicas"`
	Namespace  string            `json:"namespace"`
	Name       string            `json:"name,omitempty"`
	Containers []JobContainer    `json:"containers,omitempty"`
	JobLabels  map[string]string `json:"job_labels,omitempty"`
	Volumes    []JobVolume       `json:"volumes,omitempty"`
	TTLSeconds *int32            `json:"ttl_seconds,omitempty"`
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
	Args              []string          `json:"args,omitempty"`
	SecretsVolumes    []JobSecret       `json:"secrets_volumes,omitempty"`
	ConfigMapsVolumes []JobConfigMap    `json:"config_maps_volumes,omitempty"`
	Limits            UsageSpec         `json:"limits,omitempty"`
	Requests          UsageSpec         `json:"requests,omitempty"`
	EnvConfigMaps     []string          `json:"env_config_maps,omitempty"`
	EnvVariables      map[string]string `json:"env_variables,omitempty"`
	PreStopCmd        []string          `json:"pre_stop_cmd,omitempty"`
	PostStartCmd      []string          `json:"post_start_cmd,omitempty"`
}

type JobSecret struct {
	Name  string `json:"name"`
	Mount string `json:"mount"`
}

type JobConfigMap struct {
	Name  string `json:"name"`
	Mount string `json:"mount"`
}
