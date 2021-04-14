package cmd

import (
	"fmt"

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
	fmt.Printf("Need to implement thanos %s, %s and %s", THANOS_CHART, THANOS_CHART_REPO, THANOS_CHART_URL)
}
