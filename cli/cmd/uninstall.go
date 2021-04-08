package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Deletes the installed stack",
	Run:   helmUninstall,
}

func init() {
	RootCmd.AddCommand(uninstallCmd)
}

func helmUninstall(cmd *cobra.Command, args []string) {
	//TODO: need to remove hardcoded values
	uninstall := exec.Command("helm", "delete", "my-release")
	fmt.Println("Deleting the Prometheus stack")
	out, err := uninstall.CombinedOutput()
	if err != nil {
		fmt.Printf("could not uninstall The Observability Stack: %w \nOutput: %v", err, string(out))
	}
}
