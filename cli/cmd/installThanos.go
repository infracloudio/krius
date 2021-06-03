package cmd

import (
	"fmt"
	"log"

	"github.com/infracloudio/krius/pkg/helm"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/cli/values"
)

var thanosCmd = &cobra.Command{
	Use:   "thanos",
	Short: "Install thanos component",
	Run:   configureThanos,
}

func init() {
	installCmd.AddCommand(thanosCmd)
	addInstallFlags(thanosCmd)
	addThanosFlags(thanosCmd)
}

func configureThanos(cmd *cobra.Command, args []string) {
	var chartConfiguration *helm.HelmConfig

	// Grabbing the configureAs from the flag
	configureAs, err := cmd.Flags().GetString("configure-as")
	if err != nil {
		log.Fatalf("Could not fetch configure-as from the flags, Please try again")
	}

	if configureAs == "sidecar" {
		release, err := cmd.Flags().GetString("release")
		if release == "" {
			log.Fatal("Please provide the prometheus release name to inject the sidecar to")
		} else if err != nil {
			log.Fatalf("Could not fetch release name from the flag: %v", err)
		}
		chartConfiguration = &helm.HelmConfig{
			Repo: "prometheus-community",
			Name: "kube-prometheus-stack",
			Url:  "https://prometheus-community.github.io/helm-charts",
			Cmd:  cmd,
		}
		installThanosSidecar(chartConfiguration)
	} else if configureAs == "receiver" {
		chartConfiguration = &helm.HelmConfig{
			Repo: "bitnami",
			Name: "thanos",
			Url:  "https://charts.bitnami.com/bitnami",
		}
		installThanosReceiver(chartConfiguration, cmd)
	} else {
		log.Fatal("Invalid Option for configure-as flag, Please set either of receiver/sidecar")
	}
}

func installThanosSidecar(chartConfig *helm.HelmConfig) {
	//Creating Helm Client object out from the chart cofiguration
	helmClient, err := createHelmClientObject(chartConfig)
	if err != nil {
		log.Fatalf("Error while creating helm client object: %v", err)
	}

	//Fetching release name from the release flag
	release, err := chartConfig.Cmd.Flags().GetString("release")
	if err != nil {
		log.Fatalf("Could not grab release name from the flag: %v",err)
	}
	// Listing releases
	results, err := helmClient.ListDeployedReleases()
	if err != nil {
		log.Fatalf("Error while fetching existing releases: %v", err)
	}
	// Checking if the release name provided by the user exists
	exists := false
	for _, v := range results {
		if v.Name == release {
			exists = true
		}
	}
	if !exists {
		log.Fatalf("Invalid release name %s, Please enter a valid release name", release)
	}
	// creating a values map to upgrade the stack with
	Values := createSidecarValuesMap()
	result, err := helmClient.UpgradeChart(Values)
	if err != nil {
		log.Fatalf("Could not upgrade The Observability Stack %v", err)
	} else {
		fmt.Println(*result)
	}
}

func installThanosReceiver(chartConfig *helm.HelmConfig, cmd *cobra.Command) {
	fmt.Printf("Need to implement the receiver component")
}

func createSidecarValuesMap() *values.Options {
	var valueOpts values.Options
	valueOpts.Values = []string{fmt.Sprintf("prometheus.prometheusSpec.thanos.image=%s", "thanosio/thanos:v0.21.0-rc.0"),
		fmt.Sprintf("prometheus.prometheusSpec.thanos.sha=%s", "dbf064aadd18cc9e545c678f08800b01a921cf6817f4f02d5e2f14f221bee17c"),
		fmt.Sprintf("prometheus.thanosService.enabled=%s", "true"),
		fmt.Sprintf("prometheus.prometheusSpec.thanos.objectStorageConfig.name=%s", "krius-sidecar-config"),
		fmt.Sprintf("prometheus.prometheusSpec.thanos.objectStorageConfig.key=%s", "sidecar.yml")}
	return &valueOpts
}
