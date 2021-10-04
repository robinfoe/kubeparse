package app

import "fmt"

type AppConfig struct {
	Namespace   string
	Name        string
	Replicas    string
	ConfigMaps  []string
	ApiURL      []string
	ServiceCall []*Service

	RefferedBy map[string]*AppConfig // key is based on namespace|name
	DBProp     []interface{}
}

func (a *AppConfig) GetKey() string {
	return fmt.Sprintf("%s|%s", a.Namespace, a.Name)
}
