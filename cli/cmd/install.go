package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the given component",
	Args:  cobra.MinimumNArgs(1),
	Run:   helmInstall,
}

var thanosCmd = &cobra.Command{
	Use:   "thanos",
	Short: "Install thanos component",
	Args:  cobra.MinimumNArgs(1),
	Run:   thanosInstall,
}

func init() {
	RootCmd.AddCommand(installCmd)
	installCmd.AddCommand(thanosCmd)
}

func helmInstall(cmd *cobra.Command, args []string) {
	if strings.ToLower(args[0]) == "thanos" {
		return
	}
	//TODO: need to remove hardcoded values
	helmRepoAdd("bitnami", "https://charts.bitnami.com/bitnami")
	install := exec.Command("helm", "install", "my-release", args[0])
	fmt.Println("Installing the Prometheus stack")
	out, err := install.CombinedOutput()
	if err != nil {
		fmt.Printf("could not install The Observability Stack: %w \nOutput: %v", err, string(out))
	}
}

func helmRepoAdd(name, url string) {
	//TODO: below hardcoding need to be removed
	exec.Command("helm", "repo", "add", name, url)
}

func thanosInstall(cmd *cobra.Command, args []string) {
	//TODO: need to remove hardcoded values
	helmRepoAdd("bitnami", "https://charts.bitnami.com/bitnami")
	if strings.ToLower(args[0]) == "sidecar" {
		install := exec.Command("helm", "upgrade", "my-release", "--set",
			"prometheus.thanos.create=true", "bitnami/kube-prometheus")
		fmt.Println("Installing Thanos ", args[0])
		out, err := install.CombinedOutput()
		if err != nil {
			fmt.Printf("could not install Thanos : %w \nOutput: %v", args[0], string(out))
		}
	} else {
		if strings.ToLower(args[0]) == "receiver" {
			fmt.Println("Need to implement thanos receiver")
		}
	}
}
