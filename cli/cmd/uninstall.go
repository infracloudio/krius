package cmd

import (
	"fmt"

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
	fmt.Println("Need to implement unistall")
}
