package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/infracloudio/krius/pkg/helm"
	"github.com/spf13/cobra"
)

const (
	THANOS_CHART_REPO = "bitnami"
	THANOS_CHART      = "kube-prometheus"
	THANOS_CHART_URL  = "https://charts.bitnami.com/bitnami"
)

var thanosCmd = &cobra.Command{
	Use:   "thanos",
	Short: "Install thanos component",
	Args:  cobra.MinimumNArgs(1),
	Run:   thanosInstall,
}

func init() {
	installCmd.AddCommand(thanosCmd)
}

func thanosInstall(cmd *cobra.Command, args []string) {
	//TODO: need to remove hardcoded values
	helm.HelmRepoAdd(THANOS_CHART_REPO, THANOS_CHART_URL)
	if strings.ToLower(args[0]) == "sidecar" {
		install := exec.Command("helm", "upgrade", "my-release", "--set",
			"prometheus.thanos.create=true", THANOS_CHART_REPO+"/"+THANOS_CHART)
		fmt.Println("Installing Thanos ", args[0])
		out, err := install.CombinedOutput()
		if err != nil {
			fmt.Printf("could not install Thanos : %s \nOutput: %v", args[0], string(out))
		}
	} else {
		if strings.ToLower(args[0]) == "receiver" {
			fmt.Println("Need to implement thanos receiver")
		}
	}
}
