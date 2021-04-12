package cmd

import (
	"fmt"
	"os/exec"

	"github.com/infracloudio/krius/pkg/helm"
	"github.com/spf13/cobra"
)

const (
	PROMETHEUS_CHART_REPO = "prometheus-community"
	PROMETHEUS_CHART      = "kube-prometheus-stack"
	PROMETHEUS_CHART_URL  = "https://prometheus-community.github.io/helm-charts"
)

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

	helm.HelmRepoAdd(PROMETHEUS_CHART_REPO, PROMETHEUS_CHART_URL)

	releasename := "--generate-name"
	if len(args) > 0 {
		releasename = args[0]
	}
	cmds := []string{"install", releasename, PROMETHEUS_CHART_REPO + "/" + PROMETHEUS_CHART}
	namespace, err := cmd.Flags().GetString("namespace")
	if namespace == "" {
		namespace = "default"
	}
	cmds = append(cmds, "--create-namespace", "--namespace", namespace)
	install := exec.Command("helm", cmds...)
	fmt.Println("Installing the Prometheus stack")
	out, err := install.CombinedOutput()
	if err != nil {
		fmt.Printf("could not install The Observability Stack: %s \nOutput: %v", err.Error(), string(out))
	}
}
