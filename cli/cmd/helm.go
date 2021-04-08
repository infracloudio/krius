package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the given component",
	Args:  cobra.MinimumNArgs(1),
	Run:   helmInstall,
}

func init() {
	RootCmd.AddCommand(installCmd)
}

func helmInstall(cmd *cobra.Command, args []string) {
	//TODO: need to remove hardcoded values
	install := exec.Command("helm", "install", "my-release", args[0])
	fmt.Println("Installing the Prometheus stack")
	out, err := install.CombinedOutput()
	if err != nil {
		fmt.Printf("could not install The Observability Stack: %w \nOutput: %v", err, string(out))
	}
}

func helmRepoAdd() {
	//TODO: below hardcoding need to be removed
	exec.Command("helm", "repo", "add", "prometheus-community", "https://prometheus-community.github.io/helm-charts")
}
