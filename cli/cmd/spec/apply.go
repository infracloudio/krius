package spec

import (
	"fmt"
	"io/ioutil"
	"log"

	spec "github.com/infracloudio/krius/pkg/specvalidate"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/cli/values"
)

const (
	configFile = "config-file"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies/Updates the give profie file",
	Run:   applySpec,
}

func init() {
	specCmd.AddCommand(applyCmd)
	addSpecApplyFlags(applyCmd)
}

func applySpec(cmd *cobra.Command, args []string) {
	configFileFlag, _ := cmd.Flags().GetString(configFile)
	loader, ruleSchemaLoader, err := spec.GetLoaders(configFileFlag)
	if err != nil {
		log.Println(err)
		return
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
	objs := []ObjStoreConfigslist{}
	for _, object := range config.ObjStoreConfigslist {
		objs = append(objs, object)
	}
	preFlightErrors := []string{}
	// check for preflight errors for all the clusters
	for _, cluster := range config.Clusters {
		switch cluster.Type {
		case "prometheus":
			errs := cluster.preflightChecks(objs)
			if errs != nil {
				preFlightErrors = append(preFlightErrors, errs...)
			}
		case "thanos":
		case "grafana":
		}
	}
	if len(preFlightErrors) > 0 {
		log.Printf("Preflight checks failed %s", preFlightErrors)
		return
	}
	log.Println("Preflight checks passed")
	for _, cluster := range config.Clusters {
		switch cluster.Type {
		case "prometheus":
			cluster.installPrometheus()
		case "thanos":
		case "grafana":
		}
	}
}
func (cluster *Cluster) preflightChecks(objStores []ObjStoreConfigslist) []string {
	promSpec, err := cluster.GetConfig()
	if err != nil {
		log.Printf("Error getting config %s", err)
		return nil
	}
	spec, _ := yaml.Marshal(promSpec)
	var prom Prometheus
	yaml.Unmarshal(spec, &prom)
	return prom.preCheckProm(objStores, cluster.Name)

}
func (p *Prometheus) preCheckProm(objStores []ObjStoreConfigslist, clusterName string) []string {
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

func (cluster *Cluster) installPrometheus() {
	promSpec, err := cluster.GetConfig()
	if err != nil {
		log.Printf("Error getting config %s", err)
		return
	}
	spec, _ := yaml.Marshal(promSpec)
	var prom Prometheus
	yaml.Unmarshal(spec, &prom)
	helmClient, err := createHelmClientObject(cluster.Name, prom.Namespace)
	if err != nil {
		return
	}
	helmClient.ChartName = "kube-prometheus-stack"
	helmClient.RepoName = "prometheus-community"
	helmClient.Url = "https://prometheus-community.github.io/helm-charts"
	helmClient.ReleaseName = prom.Name
	helmClient.Namespace = prom.Namespace
	if prom.Install {
		Values := &values.Options{}
		if prom.Mode == "sidecar" {
			Values = createSidecarValuesMap(prom.ObjStoreConfig)
			_, err = helmClient.InstallChart(Values)
			if err != nil {
				log.Printf("Error installing prometheus: %s", err)
				return
			}
		} else {
			// TODO mode is receiver
		}
	} else {
		// prometheus is already installed, check the release exist & mode, then upgrade the chart
		results, err := helmClient.ListDeployedReleases()
		if err != nil {
			log.Fatalf("helm list error: %v", err)
			return
		}
		exists := false
		for _, v := range results {
			if v.Name == helmClient.ReleaseName {
				exists = true
			}
		}
		if exists {
			Values := &values.Options{}
			if prom.Mode == "sidecar" {
				Values = createSidecarValuesMap(prom.ObjStoreConfig)
				_, err = helmClient.UpgradeChart(Values)
				if err != nil {
					log.Printf("Error adding sidecar: %s", err)
					return
				}
			} else {
				// TODO mode is receiver
			}
		} else {
			log.Printf("Release %s doesn't exist", helmClient.ReleaseName)
			return
		}

	}
}

func createSidecarValuesMap(secretName string) *values.Options {
	var valueOpts values.Options
	valueOpts.Values = []string{fmt.Sprintf("prometheus.prometheusSpec.thanos.image=%s", "thanosio/thanos:v0.21.0-rc.0"),
		fmt.Sprintf("prometheus.prometheusSpec.thanos.sha=%s", "dbf064aadd18cc9e545c678f08800b01a921cf6817f4f02d5e2f14f221bee17c"),
		fmt.Sprintf("prometheus.thanosService.enabled=%s", "true"),
		fmt.Sprintf("prometheus.thanosServiceExternal.enabled=%s", "true"),
		fmt.Sprintf("prometheus.prometheusSpec.thanos.objectStorageConfig.name=%s", secretName),
		fmt.Sprintf("prometheus.prometheusSpec.thanos.objectStorageConfig.key=%s", "sidecar")}
	return &valueOpts
}
