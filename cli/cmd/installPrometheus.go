package cmd

import (
	"github.com/infracloudio/krius/pkg/helm"
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

	chartConfiguration := &helm.Config{
		Repo: "prometheus-community",
		Name: "kube-prometheus-stack",
		URL:  "https://prometheus-community.github.io/helm-charts",
		Args: args,
		Cmd:  cmd,
	}

	addAndInstallChart(chartConfiguration)
}
