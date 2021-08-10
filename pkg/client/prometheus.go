package client

import (
	"errors"
	"fmt"
	"log"

	"github.com/infracloudio/krius/pkg/helm"
	k "github.com/infracloudio/krius/pkg/kubeClient"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/cli/values"
)

type Objspec struct {
	ConfigType string          `yaml:"type"`
	Config     ObjBucketConfig `yaml:"config"`
}

type ObjBucketConfig struct {
	BucketName string `yaml:"bucket"`
	Endpoint   string `yaml:"endpoint"`
	AccessKey  string `yaml:"access_key"`
	SecretKey  string `yaml:"secret_key"`
	Insecure   bool   `yaml:"insecure"`
	Trace      Trace  `yaml:"trace"`
}

func NewPromClient(promCluster *Cluster) (Client, error) {
	promConfig, err := GetConfig(promCluster.Data, "prometheus")
	if err != nil {
		log.Printf("Error getting config %s", err)
		return nil, err
	}
	spec, _ := yaml.Marshal(promConfig)
	var Prom Prometheus
	err = yaml.Unmarshal(spec, &Prom)
	if err != nil {
		log.Printf("Error unmarshaling %s", err)
		return nil, err
	}
	return &Prom, nil
}

func (prom *Prometheus) PreflightChecks(clusterConfig *Config, clusterName string) ([]string, error) {
	if prom.Mode == "sidecar" && clusterConfig.Order == 0 {
		clusterConfig.Order = 1
	} else if clusterConfig.Order == 0 {
		clusterConfig.Order = 2
	}
	kubeClient, err := k.GetKubeClient(prom.Namespace, clusterName)
	if err != nil {
		return nil, err
	}
	promErrs := []string{}
	if prom.Install {
		err := kubeClient.CreateNSIfNotExist()
		if err != nil {
			e := fmt.Sprintf("cluster.%s: %s,", clusterName, err)
			promErrs = append(promErrs, e)
			return promErrs, nil // don't create secret, if error in creating namespace
		}
	} else {
		// if update namepsace should exist
		err := kubeClient.CheckNamespaceExist()
		if err != nil {
			e := fmt.Sprintf("cluster.%s: %s,", clusterName, err)
			promErrs = append(promErrs, e)
			return promErrs, nil // do not try to create secret, if no namespace
		}
	}
	found := false
	for _, v := range clusterConfig.ObjStoreConfigslist {
		if v.Name == prom.ObjStoreConfig {
			found = true
			secretSpec, err := createSecretforObjStore(v.Type, v.Config)
			if err != nil {
				return nil, err
			}
			err = kubeClient.CreateSecret(secretSpec, prom.ObjStoreConfig) // changes if secret name changed

			if err != nil {
				e := fmt.Sprintf("cluster.%s: %s,", prom.Name, err)
				promErrs = append(promErrs, e)
			}
			break
		}
	}
	if !found {
		e := fmt.Sprintf("cluster.%s: Bucket config doesn't exist,", clusterName)
		promErrs = append(promErrs, e)
	}
	return promErrs, nil
}

func (prom *Prometheus) InstallClient(clusterName string, receiveEndpoint []string) (string, error) {
	chartConfiguration := &helm.Config{
		Repo: "prometheus-community",
		Name: "kube-prometheus-stack",
		URL:  "https://prometheus-community.github.io/helm-charts",
	}
	helmClient, err := createHelmClientObject(clusterName, prom.Namespace, chartConfiguration)
	if err != nil {
		return "", err
	}
	helmClient.ReleaseName = prom.Name
	helmClient.Namespace = prom.Namespace
	if prom.Install {
		if prom.Mode == "sidecar" {
			Values := createSidecarValuesMap(prom.ObjStoreConfig)
			_, err = helmClient.InstallChart(Values)
			if err != nil {
				log.Printf("Error installing prometheus: %s", err)
				return "", err
			}
			target := GetPrometheusTargets(clusterName, prom.Namespace, prom.Name)
			return target[0], nil

		}
		Values := &values.Options{}
		if len(receiveEndpoint) > 0 && receiveEndpoint[0] != "" {
			Values = createPrometheusReceiverValues(receiveEndpoint[0])
		}
		_, err = helmClient.InstallChart(Values)
		if err != nil {
			log.Printf("Error installing prometheus: %s", err)
			return "", err
		}

	} else {
		// prometheus is already installed, check the release exist & mode, then upgrade the chart
		results, err := helmClient.ListDeployedReleases()
		if err != nil {
			return "", errors.New("helm list error")
		}
		exists := false
		for _, v := range results {
			if v.Name == helmClient.ReleaseName {
				exists = true
			}
		}
		if exists {
			if prom.Mode == "sidecar" {
				Values := createSidecarValuesMap(prom.ObjStoreConfig)
				_, err = helmClient.UpgradeChart(Values)
				if err != nil {
					return "", err
				}
				target := GetPrometheusTargets(clusterName, prom.Namespace, prom.Name)
				if len(target) > 0 {
					return target[0], nil
				}
				return "", errors.New("Error getting sidecar target info")
			}
			Values := createPrometheusReceiverValues(prom.ReceiveReference)
			_, err = helmClient.UpgradeChart(Values)
			if err != nil {
				log.Printf("Error installing prometheus: %s", err)
				return "", err
			}

		} else {
			errMsg := fmt.Sprintf("Release %s doesn't exist", helmClient.ReleaseName)
			return "", errors.New(errMsg)
		}
	}
	return "", err
}

func createSecretforObjStore(configType string, bucConfig BucketConfig) (map[string][]byte, error) {
	//create a secret for bucket config
	secretSpec := map[string][]byte{}
	var obj Objspec
	obj.ConfigType = configType
	obj.Config = ObjBucketConfig(bucConfig)
	objYaml, err := yaml.Marshal(obj)
	if err != nil {
		return nil, err
	}
	secretSpec["objstore.yml"] = objYaml
	return secretSpec, nil
}
