package spec

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uninstallSpecCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Deletes the entire stack across clusters",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Need to implement uninstall profile on multicluster")
	},
}

func init() {
	specCmd.AddCommand(uninstallSpecCmd)
}
