package spec

import (
	"errors"

	"github.com/infracloudio/krius/pkg/helm"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

func addSpecApplyFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("config-file", "c", "", "config file path")
	cmd.MarkFlagRequired("config-file")
}

func createHelmClientObject(context string) (*helm.HelmClient, error) {
	opt := &helm.KubeConfClientOptions{
		KubeContext: context,
	}

	action, err := helm.NewClientFromKubeConf(opt)
	if err != nil {
		return nil, err
	}
	helmClient := helm.HelmClient{
		ActionConfig: action,
	}
	return &helmClient, nil
}

type ClusterType interface {
	GetTypeName()
}

func (p Prometheus) GetTypeName() {
}
func (g Grafana) GetTypeName() {
}
func (t Thanos) GetTypeName() {
}

func (c Cluster) GetConfig(kind string, data Data) (ClusterType, error) {
	switch kind {
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
	return nil, errors.New("wrong cluster type")
}
