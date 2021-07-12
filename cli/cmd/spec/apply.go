package spec

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sort"

	"github.com/infracloudio/krius/pkg/helm"
	spec "github.com/infracloudio/krius/pkg/specvalidate"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	configFile = "config-file"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies/Updates the give profie file",
	RunE:  applySpec,
}

func init() {
	specCmd.AddCommand(applyCmd)
	addSpecApplyFlags(applyCmd)
}

func applySpec(cmd *cobra.Command, args []string) (err error) {
	configFileFlag, _ := cmd.Flags().GetString(configFile)
	loader, ruleSchemaLoader, err := spec.GetLoaders(configFileFlag)
	if err != nil {
		return err
	}
	valid, errors := spec.ValidateYML(loader, ruleSchemaLoader)
	if !valid {
		log.Println(errors)
		return
	}
	log.Println("valid yaml")
	yamlFile, err := ioutil.ReadFile(configFileFlag)
	if err != nil {
		log.Fatalf("yamlFile.Get err #%v ", err)
	}
	var config Config
	yaml.Unmarshal(yamlFile, &config)
	preFlightErrors := []string{}
	// check for preflight errors for all the clusters
	for _, cluster := range config.Clusters {
		switch cluster.Type {
		case "prometheus":
			errs := cluster.preflightChecks(config.ObjStoreConfigslist, &config)
			if errs != nil {
				preFlightErrors = append(preFlightErrors, errs...)
			}
		case "thanos":
			errs := cluster.preflightChecks(config.ObjStoreConfigslist, &config)
			if errs != nil {
				preFlightErrors = append(preFlightErrors, errs...)
			}
		case "grafana":
		}
	}
	if len(preFlightErrors) > 0 {
		log.Printf("Preflight checks failed %s", preFlightErrors)
		return
	}
	// reorder clusters based on sidecar/reciever setup
	if config.Order == 1 {
		sort.Slice(config.Clusters, func(p, q int) bool {
			return config.Clusters[p].Type < config.Clusters[q].Type
		})
	} else if config.Order == 2 {
		sort.Slice(config.Clusters, func(p, q int) bool {
			return config.Clusters[p].Type > config.Clusters[q].Type
		})
	}
	log.Println("Preflight checks passed")
	var targets []string
	for _, cluster := range config.Clusters {
		switch cluster.Type {
		case "prometheus":
			target, err := cluster.installPrometheus()
			if err != nil {
				return err
			}
			targets = append(targets, target)
		case "thanos":
			cluster.installThanos(targets)
		case "grafana":
		}
	}
	return
}
func (cluster *Cluster) preflightChecks(objStores []ObjStoreConfig, c *Config) []string {
	if cluster.Type == "prometheus" {
		promSpec, err := cluster.GetConfig()
		if err != nil {
			log.Printf("Error getting config %s", err)
			return nil
		}
		spec, _ := yaml.Marshal(promSpec)
		var prom Prometheus
		yaml.Unmarshal(spec, &prom)
		if prom.Mode == "sidecar" && c.Order == 0 {
			c.Order = 1
		} else if c.Order == 0 {
			c.Order = 2
		}
		return prom.preCheckProm(objStores, cluster.Name)
	} else if cluster.Type == "thanos" {
		thanosSpec, err := cluster.GetConfig()
		if err != nil {
			log.Printf("Error getting config %s", err)
			return nil
		}
		spec, err := yaml.Marshal(thanosSpec)
		var thanos Thanos
		yaml.Unmarshal(spec, &thanos)
		return thanos.preCheckThanos(objStores, cluster.Name)
	}
	return nil

}

func (t *Thanos) preCheckThanos(objStores []ObjStoreConfig, clusterName string) []string {
	promErrs := []string{}

	found := false
	for _, v := range objStores {
		if v.Name == t.ObjStoreConfig {
			found = true
			err := createSecretforObjStore(clusterName, t.Namespace, v.Type, v.Name, v.Config)
			if err != nil {
				e := fmt.Sprintf("cluster.%s: %s,", t.Name, err)
				promErrs = append(promErrs, e)
			}
			break
		}
	}
	if !found {
		e := fmt.Sprintf("cluster.%s: Bucket config doesn't exist,", clusterName)
		promErrs = append(promErrs, e)
	}
	return promErrs
}

func (p *Prometheus) preCheckProm(objStores []ObjStoreConfig, clusterName string) []string {
	promErrs := []string{}
	if p.Install {
		err := CreateNameSpaceIfNotExist(clusterName, p.Namespace)
		if err != nil {
			e := fmt.Sprintf("cluster.%s: %s,", clusterName, err)
			promErrs = append(promErrs, e)
			return promErrs // do not try to create secret, if err in creating
		}
	} else {
		// if update namepsace should exist
		err := CheckNamespaceExist(clusterName, p.Namespace)
		if err != nil {
			e := fmt.Sprintf("cluster.%s: %s,", clusterName, err)
			promErrs = append(promErrs, e)
			return promErrs // do not try to create secret, if no namespace
		}
	}
	found := false
	for _, v := range objStores {
		if v.Name == p.ObjStoreConfig {
			found = true
			err := createSecretforObjStore(clusterName, p.Namespace, v.Type, v.Name, v.Config)
			if err != nil {
				e := fmt.Sprintf("cluster.%s: %s,", p.Name, err)
				promErrs = append(promErrs, e)
			}
			break
		}
	}
	if !found {
		e := fmt.Sprintf("cluster.%s: Bucket config doesn't exist,", clusterName)
		promErrs = append(promErrs, e)
	}
	return promErrs
}

func (cluster *Cluster) installPrometheus() (string, error) {
	promSpec, err := cluster.GetConfig()
	if err != nil {
		return "", errors.New("Error getting config")
	}
	spec, _ := yaml.Marshal(promSpec)
	var prom Prometheus
	yaml.Unmarshal(spec, &prom)

	chartConfiguration := &helm.HelmConfig{
		Repo: "prometheus-community",
		Name: "kube-prometheus-stack",
		Url:  "https://prometheus-community.github.io/helm-charts",
	}
	helmClient, err := createHelmClientObject(cluster.Name, prom.Namespace, chartConfiguration)
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
			} else {
				target := GetPrometheusTargets(cluster.Name, prom.Namespace, prom.Name)
				return target[0], nil
			}
		} // TODO mode is receiver

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
				} else {
					target := GetPrometheusTargets(cluster.Name, prom.Namespace, prom.Name)
					if len(target) > 0 {
						return target[0], nil
					}
					return "", errors.New("Error getting sidecar target info")

				}
			}
			// TODO mode is receiver
		} else {
			errMsg := fmt.Sprintf("Release %s doesn't exist", helmClient.ReleaseName)
			return "", errors.New(errMsg)
		}
	}
	return "", err
}

func (cluster *Cluster) installThanos(targets []string) error {
	thanosSpec, _ := cluster.GetConfig()
	spec, _ := yaml.Marshal(thanosSpec)
	var thanos Thanos
	yaml.Unmarshal(spec, &thanos)
	chartConfiguration := &helm.HelmConfig{
		Repo: "bitnami",
		Name: "thanos",
		Url:  "https://charts.bitnami.com/bitnami",
	}

	helmClient, err := createHelmClientObject(cluster.Name, thanos.Namespace, chartConfiguration)
	if err != nil {
		return err
	}
	helmClient.ChartName = "thanos"
	helmClient.ReleaseName = "thanos"
	var extraFlags []string
	if thanos.Querier.AutoDownsample {
		extraFlags = append(extraFlags, "--query.auto-downsampling")
	}
	if thanos.Querier.PartialResponse {
		extraFlags = append(extraFlags, "--query.partial-response")
	}
	thanos.Querier.ExtraFlags = extraFlags
	thanos.Querier.Targets = targets

	Values := createThanosValuesMap(thanos)
	_, err = helmClient.InstallChart(Values)
	log.Println("error installing Thanos", err)
	return err
}
