package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [Name]",
	Short: "Deletes the installed stack",
	Run:   helmUninstall,
}

func init() {
	RootCmd.AddCommand(uninstallCmd)
}

func helmUninstall(cmd *cobra.Command, args []string) {
	releasename := ""
	if len(args) > 0 {
		releasename = args[0]
	}
	cmds := []string{"delete", releasename}
	namespace, err := cmd.Flags().GetString("namespace")
	if namespace == "" {
		namespace = "default"
	}
	cmds = append(cmds, "--create-namespace", "--namespace", namespace)
	uninstall := exec.Command("helm", cmds...)

	fmt.Println("Deleting the Prometheus stack")
	out, err := uninstall.CombinedOutput()
	if err != nil {
		fmt.Printf("could not uninstall The Observability Stack: %s \nOutput: %s", err.Error(), string(out))
	}
}
