package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

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
	addInstallFlags(thanosCmd)
}

func thanosInstall(cmd *cobra.Command, args []string) {

	helm.HelmRepoAdd(THANOS_CHART_REPO, THANOS_CHART_URL)
	if strings.ToLower(args[0]) == "sidecar" {
		releasename := args[1]
		cmds := []string{"upgrade", releasename, THANOS_CHART_REPO + "/" + THANOS_CHART}
		cmds = append(cmds, "--set", "prometheus.thanos.create=true")
		namespace, err := cmd.Flags().GetString("namespace")
		if namespace == "" {
			namespace = "default"
		}
		cmds = append(cmds, "--namespace", namespace)

		install := exec.Command("helm", cmds...)
		fmt.Printf("Installing Thanos %s", args[0])
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
