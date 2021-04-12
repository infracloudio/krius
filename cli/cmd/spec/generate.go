package spec

import (
	"fmt"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a profile based on questions asked to user",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Need to implement generate profile on multicluster")
	},
}

func init() {
	specCmd.AddCommand(generateCmd)
}
