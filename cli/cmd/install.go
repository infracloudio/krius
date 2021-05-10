package cmd

import (
	"fmt"
	"log"

	"github.com/infracloudio/krius/pkg/helm"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/cli"
)

// Struct for the Chart Configuration

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the given component",
	Args:  cobra.MinimumNArgs(1),
}

var settings *cli.EnvSettings

func init() {
	RootCmd.AddCommand(installCmd)
}

// Adds/Updates the repo and Installs the chart
func addAndInstallChart(config *helm.HelmConfig) {
	helmClient, err := createHelmClientObject(config)
	if err != nil {
		log.Fatal(err)
	}
	err = helmClient.AddRepo()
	if err != nil {
		log.Fatal(err)
		return
	}
	err = helmClient.UpdateRepo()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Installing the Prometheus stack")
	_, err = helmClient.InstallChart()
	if err != nil {
		fmt.Printf("could not install The Observability Stack %s", err)
	}
}
