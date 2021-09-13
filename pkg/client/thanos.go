package client

import (
	"errors"
	"fmt"
	"log"

	"github.com/infracloudio/krius/pkg/helm"
	k "github.com/infracloudio/krius/pkg/kubeClient"
	"gopkg.in/yaml.v2"
)

func NewThanosClient(thanosCluster *Cluster) (Client, error) {
	thanosConfig, err := GetConfig(thanosCluster.Data, "thanos")
	if err != nil {
		log.Printf("Error getting config %s", err)
		return nil, err
	}
	spec, _ := yaml.Marshal(thanosConfig)
	var thanos Thanos
	err = yaml.Unmarshal(spec, &thanos)
	if err != nil {
		log.Printf("Error unmarshaling %s", err)
		return nil, err
	}
	return &thanos, nil
}

func (t *Thanos) PreflightChecks(clusterConfig *Config, clusterName string) ([]string, error) {
	thanosErrs := []string{}

	if clusterConfig.Order == 2 {
		receiver := Receiver{}
		if (t.Receiver) == receiver {
			e := fmt.Sprintf("cluster.%s: %s,", clusterName, "Receiver not set")
			thanosErrs = append(thanosErrs, e)
		}
	}
	kubeClient, err := k.GetKubeClient(t.Namespace, clusterName)
	if err != nil {
		return nil, err
	}
	err = kubeClient.CreateNSIfNotExist()
	if err != nil {
		e := fmt.Sprintf("cluster.%s: %s,", clusterName, err)
		thanosErrs = append(thanosErrs, e)
		return thanosErrs, nil // don't create secret, if error in creating namespace
	}

	found := false
	for _, v := range clusterConfig.ObjStoreConfigslist {
		if v.Name == t.ObjStoreConfig {
			found = true
			secretSpec, err := createSecretforObjStore(v.Type, v.Config)
			if err != nil {
				return nil, err
			}
			err = kubeClient.CreateSecret(secretSpec, t.ObjStoreConfig)

			if err != nil {
				e := fmt.Sprintf("cluster.%s: %s,", t.Name, err)
				thanosErrs = append(thanosErrs, e)
			}
			break
		}
	}
	if !found {
		e := fmt.Sprintf("cluster.%s: Bucket config doesn't exist,", clusterName)
		thanosErrs = append(thanosErrs, e)
	}
	return thanosErrs, nil
}

func (t *Thanos) InstallClient(clusterName string, targets []string) (string, error) {

	chartConfiguration := &helm.Config{
		Repo: "bitnami",
		Name: "thanos",
		URL:  "https://charts.bitnami.com/bitnami",
	}

	helmClient, err := createHelmClientObject(clusterName, t.Namespace, chartConfiguration)
	if err != nil {
		return "", err
	}
	helmClient.ChartName = "thanos"
	helmClient.ReleaseName = t.Name
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
	t.Querier.Targets = targets
	values := t.createThanosValuesMap()

	_, err = helmClient.InstallChart(values)
	if err != nil {
		log.Printf("Error installing thanos: %s", err)
		return "", err
	}
	if t.Receiver.Name == "" { // sidecar mode
		return "", nil
	}
	receiveEndpoint := getReceiveEndpoint(clusterName, t.Namespace, t.Name)
	if len(receiveEndpoint) > 0 {
		return receiveEndpoint[0], nil
	}
	return "", nil
}

func (t *Thanos) UninstallClient(clusterName string) error {
	chartConfiguration := &helm.Config{
		Repo: "bitnami",
		Name: "thanos",
		URL:  "https://charts.bitnami.com/bitnami",
	}

	helmClient, err := createHelmClientObject(clusterName, t.Namespace, chartConfiguration)
	if err != nil {
		return err
	}

	helmClient.ChartName = "thanos"
	helmClient.ReleaseName = t.Name
	helmClient.Namespace = t.Namespace

	// thanos is already installed, check the release exist & mode, then uninstall the chart
	results, err := helmClient.ListDeployedReleases()
	if err != nil {
		return errors.New("helm list error")
	}
	exists := false
	for _, v := range results {
		if v.Name == helmClient.ReleaseName {
			exists = true
		}
	}

	if exists {
		_, err = helmClient.UninstallChart()
		if err != nil {
			log.Printf("Error uninstalling thanos: %s", err)
			return err
		}

	} else {
		errMsg := fmt.Sprintf("Release %s doesn't exist", helmClient.ReleaseName)
		return errors.New(errMsg)
	}

	return nil
}
