package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/infracloudio/krius/pkg/helm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var thanosCmd = &cobra.Command{
	Use:   "thanos",
	Short: "Install thanos component",
	Args:  cobra.MinimumNArgs(1),
	Run:   thanosInstall,
}

func init() {
	installCmd.AddCommand(thanosCmd)
	helm.AddInstallFlags(thanosCmd)
}

func thanosInstall(cmd *cobra.Command, args []string) {
	thanosRepo, ok := viper.Get("thanos.repo").(string)
	if !ok {
		log.Fatalf("Invalid thanos repo name")
	}

	thanosRepoUrl, ok := viper.Get("thanos.url").(string)
	if !ok {
		log.Fatalf("Invalid thanos url")
	}

	thanosChart, ok := viper.Get("thanos.chart").(string)
	if !ok {
		log.Fatalf("Invalid thanos chart name")
	}

	helm.HelmRepoAdd(thanosRepo, thanosRepoUrl)
	if strings.ToLower(args[0]) == "sidecar" {
		releasename := args[1]
		cmds := []string{"upgrade", releasename, thanosRepo + "/" + thanosChart}
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
