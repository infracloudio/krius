package cmd

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/infracloudio/krius/pkg/helm"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/cli"
)

type conf struct {
	ReleaseName string `yaml:"release"`
	Namespace   string `yaml:"namespace"`
}

func (c *conf) getConf(valuesYaml string) *conf {
	yamlFile, err := ioutil.ReadFile(valuesYaml)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}

func createHelmClientObject(helmConfig *helm.HelmConfig) (*helm.HelmClient, error) {
	var namespace string
	var releaseName string
	valuesYaml, _ := helmConfig.Cmd.Flags().GetString(configFile)
	if valuesYaml != "" {
		var c conf
		c.getConf(valuesYaml)
		namespace = c.Namespace
		releaseName = c.ReleaseName
	} else {
		namespace, _ = helmConfig.Cmd.Flags().GetString("namespace")
		releaseName, _ = helmConfig.Cmd.Flags().GetString("release")
	}
	if namespace == "" {
		return nil, errors.New("Please set Namespace")
	}
	if releaseName == "" {
		return nil, errors.New("Please set Release name")
	}
	os.Setenv("HELM_NAMESPACE", namespace)
	settings = cli.New()

	action, err := helm.InitializeHelmAction(settings)
	if err != nil {
		log.Fatal(err)
	}
	helmClient := helm.HelmClient{
		RepoName:     helmConfig.Repo,
		Url:          helmConfig.Url,
		ReleaseName:  releaseName,
		Namespace:    namespace,
		ChartName:    helmConfig.Name,
		ActionConfig: action,
		Settings:     settings,
	}
	return &helmClient, err
}
