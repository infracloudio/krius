package cmd

import (
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

	chartConfiguration := &helmConfig{
		repo: "prometheus-community",
		name: "kube-prometheus-stack",
		url:  "https://prometheus-community.github.io/helm-charts",
		args: args,
		cmd:  cmd,
	}

	addAndInstallChart(chartConfiguration)
}
