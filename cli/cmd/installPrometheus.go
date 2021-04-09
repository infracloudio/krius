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
)

var settings *cli.EnvSettings

var prometheusCmd = &cobra.Command{
	Use:   "prometheus [Name]",
	Short: "Install prometheus stack",
	Run:   prometheusInstall,
}

func init() {
	installCmd.AddCommand(prometheusCmd)
	helm.AddInstallFlags(prometheusCmd)
}

func prometheusInstall(cmd *cobra.Command, args []string) {
	promRepo, ok := viper.Get("prometheus.repo").(string)
	if !ok {
		log.Fatalf("Invalid prometheus repo name")
	}

	promUrl, ok := viper.Get("prometheus.url").(string)
	if !ok {
		log.Fatalf("Invalid prometheus url")
	}

	promChart, ok := viper.Get("prometheus.chart").(string)
	if !ok {
		log.Fatalf("Invalid prometheus chart name")
	}

	helm.HelmRepoAdd(promRepo, promUrl)
	releaseName := "--generate-name"
	if len(args) > 0 {
		releaseName = args[0]
	}
	namespace, err := cmd.Flags().GetString("namespace")
	if err != nil {
		namespace = "default"
	}
	os.Setenv("HELM_NAMESPACE", namespace)
	settings = cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debug); err != nil {
		log.Fatal(err)
	}
	client := action.NewInstall(actionConfig)
	helmClient := helm.HelmClient{
		RepoName:    promRepo,
		ReleaseName: releaseName,
		Namespace:   namespace,
		ChartName:   promChart,
		Client:      client,
		Settings:    settings,
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
