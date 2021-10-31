package client

import (
	"errors"

	"github.com/mitchellh/mapstructure"
)

type ClusterSpec interface {
	GetTypeName()
}

//any value which implements GetTypeName is also of type Clustter Spec interface
func (p Prometheus) GetTypeName() {
}
func (g Grafana) GetTypeName() {
}
func (t Thanos) GetTypeName() {
}

func getConfig(data map[string]interface{}, clusterType string) (ClusterSpec, error) {
	switch clusterType {
	case "prometheus":
		s := Prometheus{}
		err := mapstructure.Decode(data, &s)
		if err != nil {
			return nil, err
		}
		return s, nil
	case "thanos":
		t := Thanos{}
		err := mapstructure.Decode(data, &t)
		if err != nil {
			return nil, err
		}
		return t, nil
	case "grafana":
		g := Grafana{}
		err := mapstructure.Decode(data, &g)
		if err != nil {
			return nil, err
		}
		return g, nil
	}
	return nil, errors.New("cluster doesn't exist")
}
