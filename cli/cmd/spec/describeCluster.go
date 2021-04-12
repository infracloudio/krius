package spec

import (
	"fmt"

	"github.com/spf13/cobra"
)

var describeClusterCmd = &cobra.Command{
	Use:   "describe-cluster",
	Short: "Describes the entire stack across multiple clusters and current state",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Need to implement describe-cluster on multicluster")
	},
}

func init() {
	specCmd.AddCommand(describeClusterCmd)
}
