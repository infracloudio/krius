package spec

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listProfilesCmd = &cobra.Command{
	Use:   "list-profiles",
	Short: "Shows all the profile across multiple clusters",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Need to implement list-profiles on multicluster")
	},
}

func init() {
	specCmd.AddCommand(listProfilesCmd)
}
