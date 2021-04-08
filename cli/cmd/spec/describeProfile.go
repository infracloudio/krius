package spec

import (
	"fmt"

	"github.com/spf13/cobra"
)

var describeProfileCmd = &cobra.Command{
	Use:   "describe-profile",
	Short: "Describes the profile across multiple clusters and current state",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Need to implement describe-profile on multicluster")
	},
}

func init() {
	specCmd.AddCommand(describeProfileCmd)
}
