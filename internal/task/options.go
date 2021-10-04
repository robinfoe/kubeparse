package task

type PathOptions struct {
	DeploymentConfig string
	ConfigMap        string
	FileOutput       string
}

type ParseOptions struct {
	Path *PathOptions
}

//** ** ** ** ** ** ** ** ** **//
type DeploymentConfigList struct {
	APIVersion string             `yaml:"apiVersion"`
	Items      []DeploymentConfig `yaml:"items"`
}

type DeploymentConfig struct {
	ApiVersion string       `yaml:"apiVersion"`
	Kind       string       `yaml:"kind"`
	Metadata   KubeMetaData `yaml:"metadata"`
	Spec       KubeSpec     `yalm:"spec"`
}

type KubeMetaData struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type KubeSpec struct {
	Replicas string       `yaml:"replicas"`
	Template KubeTemplate `yaml:"template"`
}

type KubeTemplate struct {
	Spec struct {
		Containers []struct {
			Name string `yaml:"name"`
			Env  []struct {
				Name      string `yaml:"name"`
				Value     string `yaml:"value"`
				ValueFrom struct {
					ConfigMapKeyRef struct {
						Key  string `yaml:"key"`
						Name string `yaml:"name"`
					} `yaml:"configMapKeyRef"`
				} `yaml:"valueFrom"`
			} `yaml:"env"`

			EnvFrom []struct {
				ConfigMapRef struct {
					Name string `yaml:"name"`
				} `yaml:"configMapRef"`
			} `yaml:"envFrom"`
		} `yaml:"containers"`

		Volumes []struct {
			Name      string `yaml:"name"`
			ConfigMap struct {
				Name string `yaml:"name"`
			} `yaml:"configMap"`
		} `yaml:"volumes"`
	} `yaml:"spec"`
}

type ConfigMapList struct {
	APIVersion string      `yaml:"apiVersion"`
	Items      []ConfigMap `yaml:"items"`
}

type ConfigMap struct {
	ApiVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   KubeMetaData      `yaml:"metadata"`
	Data       map[string]string `yaml:"data"`
}
