package spec

import (
	"fmt"

	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "plan the specified profile",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Need to implement plan profile on multicluster")
	},
}

var preCheckCmd = &cobra.Command{
	Use:   "pre-check",
	Short: "pre-check the specified profile",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Need to implement pre-check profile on multicluster")
	},
}

func init() {
	specCmd.AddCommand(planCmd)
	specCmd.AddCommand(preCheckCmd)
}
