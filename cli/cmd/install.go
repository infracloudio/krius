package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/infracloudio/krius/pkg/helm"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/cli"
)

// Struct for the Chart Configuration
type helmConfig struct {
	repo string
	name string
	url  string
	args []string
	cmd  *cobra.Command
}

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
func addAndInstallChart(config *helmConfig) {
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
	result, err := helmClient.InstallOrUpgradeChart()
	if err != nil {
		fmt.Printf("could not install The Observability Stack %s", err)
	}
	fmt.Println(*result)
}

func createHelmClientObject(helmConfig *helmConfig) (*helm.HelmClient, error) {
	var releaseName string
	if len(helmConfig.args) > 0 {
		releaseName = helmConfig.args[0]
	}
	namespace, err := helmConfig.cmd.Flags().GetString("namespace")
	if err != nil {
		namespace = "default"
	}
	os.Setenv("HELM_NAMESPACE", namespace)
	settings = cli.New()

	client, err := helm.InitializeHelmAction(settings)
	if err != nil {
		log.Fatal(err)
	}
	helmClient := helm.HelmClient{
		RepoName:    helmConfig.repo,
		Url:         helmConfig.url,
		ReleaseName: releaseName,
		Namespace:   namespace,
		ChartName:   helmConfig.name,
		Client:      client,
		Settings:    settings,
	}
	return &helmClient, err
}
