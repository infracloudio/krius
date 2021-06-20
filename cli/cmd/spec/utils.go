package spec

import (
	"errors"
	"log"
	"os"

	kube "github.com/infracloudio/krius/pkg/kubeClient"
	"gopkg.in/yaml.v2"

	"github.com/infracloudio/krius/pkg/helm"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/cli"
)

type objspec struct {
	Type string          `yaml:"type"`
	Data ObjBucketConfig `yaml:"config"`
}

type ObjBucketConfig struct {
	BucketName string `yaml:"bucket"`
	Endpoint   string `yaml:"endpoint"`
	AccessKey  string `yaml:"access_key"`
	SecretKey  string `yaml:"secret_key"`
}

var settings *cli.EnvSettings

func addSpecApplyFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("config-file", "c", "", "config file path")
	cmd.MarkFlagRequired("config-file")
}

func createSecretforObjStore(clusterName, namespace, configType, secretName string, bucConfig BucketConfig) error {
	//create a secret for bucket config
	secretSpec := map[string][]byte{}
	bucket := ObjBucketConfig{bucConfig.BucketName, bucConfig.Endpoint, bucConfig.AccessKey, bucConfig.SecretKey}
	var obj objspec
	obj.Type = configType
	obj.Data = bucket
	objYaml, _ := yaml.Marshal(obj)
	secretSpec["sidecar"] = []byte(objYaml)
	kubeClient, err := GetKubeClient(namespace, clusterName)
	if err != nil {
		log.Printf("Error creating kube client %s", err)
		return err
	}
	err = kubeClient.CreateSecret(secretSpec, secretName)
	if err != nil {
		log.Printf("Error creating secret %s", err)
		return err
	}
	return nil
}

func CreateNameSpaceIfNotExist(clusterName, namespace string) error {
	kubeClient, err := GetKubeClient(namespace, clusterName)
	if err != nil {
		log.Print("Error creating kube client")
		return err
	}
	err = kubeClient.CreateNSIfNotExist()
	if err != nil {
		log.Print("Error creating namespace")
		return err
	}
	return nil
}
func CheckNamespaceExist(clusterName, namespace string) error {
	kubeClient, err := GetKubeClient(namespace, clusterName)
	if err != nil {
		log.Print("Error creating kube client")
		return err
	}
	err = kubeClient.CheckNamespaceExist()
	if err != nil {
		log.Print("Error getting namespace")
		return err
	}
	return nil
}

func createHelmClientObject(context, namespace string) (*helm.HelmClient, error) {
	opt := &helm.KubeConfClientOptions{
		KubeContext: context,
	}
	os.Setenv("HELM_NAMESPACE", namespace)
	settings = cli.New()
	action, err := helm.NewClientFromKubeConf(opt, settings)
	if err != nil {
		return nil, err
	}
	helmClient := helm.HelmClient{
		ActionConfig: action,
		Settings:     settings,
	}
	return &helmClient, nil
}

func GetKubeClient(namespace, context string) (*kube.KubeConfig, error) {
	kubeClient := kube.KubeConfig{
		Namespace: namespace,
		Context:   context,
	}
	err := kubeClient.InitClient()
	if err != nil {
		return nil, err
	}
	return &kubeClient, nil
}

type ClusterSpec interface {
	GetTypeName()
}

func (p Prometheus) GetTypeName() {
}
func (g Grafana) GetTypeName() {
}
func (t Thanos) GetTypeName() {
}

func (c *Cluster) GetConfig() (ClusterSpec, error) {
	switch c.Type {
	case "prometheus":
		s := Prometheus{}
		err := mapstructure.Decode(c.Data, &s)
		if err != nil {
			return nil, err
		}
		return s, nil
	case "thanos":
		t := Thanos{}
		err := mapstructure.Decode(c.Data, &t)
		if err != nil {
			return nil, err
		}
		return t, nil
	case "grafana":
		g := Grafana{}
		err := mapstructure.Decode(c.Data, &g)
		if err != nil {
			return nil, err
		}
		return g, nil
	}
	return nil, errors.New("wrong cluster type")
}
