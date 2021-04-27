package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

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

	chartConfiguration := &ChartConfig{
		CHART_REPO: "prometheus-community",
		CHART_NAME: "kube-prometheus-stack",
		CHART_URL:  "https://prometheus-community.github.io/helm-charts",
	}

	helmClient, err := createHelmClientObject(cmd, args, chartConfiguration)
	if err != nil {
		log.Fatal(err)
	}
	addAndInstallChart(helmClient)
}

func debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	log.Output(2, fmt.Sprintf(format, v...))
}
