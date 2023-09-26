package entropy

type KubeJobConfig struct {
	Stopped    bool               `json:"stopped,omitempty"`
	Replicas   int32              `json:"replicas"`
	Namespace  string             `json:"namespace,omitempty"`
	Name       string             `json:"name,omitempty"`
	Containers []KubeJobContainer `json:"containers,omitempty"`
	JobLabels  map[string]string  `json:"job_labels,omitempty"`
	Volumes    []KubeJobVolume    `json:"volumes,omitempty"`
}

type KubeJobVolume struct {
	Name string
	Kind string // secret or config-map. secret is for gcs/bq credential
}

type KubeJobContainer struct {
	Name              string             `json:"name"`
	Image             string             `json:"image"`
	ImagePullPolicy   string             `json:"image_pull_policy,omitempty"`
	Command           []string           `json:"command,omitempty"`
	SecretsVolumes    []KubeJobSecret    `json:"secrets_volumes,omitempty"`
	ConfigMapsVolumes []KubeJobConfigMap `json:"config_maps_volumes,omitempty"`
	Limits            *UsageSpec         `json:"limits,omitempty"`
	Requests          *UsageSpec         `json:"requests,omitempty"`
	EnvConfigMaps     []string           `json:"env_config_maps,omitempty"`
	EnvVariables      map[string]string  `json:"env_variables,omitempty"`
}

type KubeJobSecret struct {
	Name  string `json:"name"`
	Mount string `json:"mount"`
}

type KubeJobConfigMap struct {
	Name  string `json:"name"`
	Mount string `json:"mount"`
}
