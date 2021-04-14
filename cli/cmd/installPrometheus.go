package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/infracloudio/krius/pkg/helm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"os/exec"

	"github.com/infracloudio/krius/pkg/helm"
	"github.com/spf13/cobra"
)

const (
	PROMETHEUS_CHART_REPO = "prometheus-community"
	PROMETHEUS_CHART      = "kube-prometheus-stack"
	PROMETHEUS_CHART_URL  = "https://prometheus-community.github.io/helm-charts"
)

var settings *cli.EnvSettings

var prometheusCmd = &cobra.Command{
	Use:   "prometheus [Name]",
	Short: "Install prometheus stack",
	Run:   prometheusInstall,
}

func init() {
	installCmd.AddCommand(prometheusCmd)
	addInstallFlags(prometheusCmd)
}

func prometheusInstall(cmd *cobra.Command, args []string) {
	releaseName := "--generate-name"
	if len(args) > 0 {
		releaseName = args[0]
	}
	namespace, err := cmd.Flags().GetString("namespace")
	if err != nil {
		namespace = "default"
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debug); err != nil {
		log.Fatal(err)
		return
	}
	client := action.NewInstall(actionConfig)
	helmClient := helm.HelmClient{
		RepoName:    promRepo,
		Url:         promUrl,
		ReleaseName: releaseName,
		Namespace:   namespace,
		ChartName:   promChart,
		Client:      client,
		Settings:    settings,
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

func debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	log.Output(2, fmt.Sprintf(format, v...))
}
