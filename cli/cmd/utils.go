package cmd

import (
	"log"
	"os"

	"github.com/infracloudio/krius/pkg/helm"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/cli"
)

// Struct for the Chart Configuration
type ChartConfig struct {
	CHART_REPO string
	CHART_NAME string
	CHART_URL  string
}

var settings *cli.EnvSettings

func addInstallFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("namespace", "n", "default", "namespace in which the chart need to be installed")
}

func createHelmClientObject(cmd *cobra.Command, args []string, chartConfig *ChartConfig) (*helm.HelmClient, error) {
	var releaseName string
	if len(args) > 0 {
		releaseName = args[0]
	}
	namespace, err := cmd.Flags().GetString("namespace")
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
		RepoName:    chartConfig.CHART_REPO,
		Url:         chartConfig.CHART_URL,
		ReleaseName: releaseName,
		Namespace:   namespace,
		ChartName:   chartConfig.CHART_NAME,
		Client:      client,
		Settings:    settings,
	}
	return &helmClient, err
}
