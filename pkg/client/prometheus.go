package client

import (
	"errors"
	"fmt"
	"log"

	"github.com/infracloudio/krius/pkg/helm"
	k "github.com/infracloudio/krius/pkg/kubeClient"
	"gopkg.in/yaml.v2"
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

	err = helmClient.AddRepo()
	if err != nil {
		log.Fatalf("helm add repo error: %v", err)
		return "", err
	}
	err = helmClient.UpdateRepo()
	if err != nil {
		log.Fatalf("helm update repo error: %v", err)
		return "", err
	}
	if prom.Install {
		if prom.Mode == "sidecar" {
			Values := prom.createPrometheusSidecarValues()
			_, err = helmClient.InstallChart(Values)
			if err != nil {
				log.Printf("Error installing prometheus: %s", err)
				return "", err
			}
			target := getPrometheusTargets(clusterName, prom.Namespace, prom.Name)
			if len(target) > 0 {
				return target[0], nil
			}
			return "", errors.New("error getting sidecar target info")

		}
		Values := prom.createPrometheusReceiverValues(receiveEndpoint)
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
				Values := prom.createPrometheusSidecarValues()
				_, err = helmClient.UpgradeChart(Values)
				if err != nil {
					return "", err
				}
				target := getPrometheusTargets(clusterName, prom.Namespace, prom.Name)
				if len(target) > 0 {
					return target[0], nil
				}
				return "", errors.New("error getting sidecar target info")
			}
			Values := prom.createPrometheusReceiverValues(receiveEndpoint)
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
