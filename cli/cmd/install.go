package cmd

import (
	"fmt"
	"log"

	"github.com/infracloudio/krius/pkg/helm"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the given component",
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	RootCmd.AddCommand(installCmd)
}

// Adds/Updates the repo and Installs the chart
func addAndInstallChart(helmClient *helm.HelmClient) {
	err := helmClient.AddRepo()
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
	result, err := helmClient.InstallOrUpgradeChart()
	if err != nil {
		fmt.Printf("could not install The Observability Stack %s", err)
	}
	fmt.Println(*result)
}
