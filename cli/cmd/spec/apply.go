package spec

import (
	"fmt"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applys the specified profile",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Apply the spec")
	},
}

func init() {
	specCmd.AddCommand(applyCmd)
}
