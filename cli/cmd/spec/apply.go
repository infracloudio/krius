package spec

import (
	"fmt"

	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies/Updates the give profie file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Need to implement apply profile on multicluster")
	},
}

func init() {
	specCmd.AddCommand(applyCmd)
}
