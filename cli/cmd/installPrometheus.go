package cmd

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/infracloudio/krius/pkg/helm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	releasename := "--generate-name"
	if len(args) > 0 {
		releasename = args[0]
	}
	cmds := []string{"install", releasename, promRepo + "/" + promChart}
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
