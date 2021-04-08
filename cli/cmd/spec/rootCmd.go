package spec

import (
	"fmt"

	"github.com/infracloudio/krius/cli/cmd"
	"github.com/spf13/cobra"
)

var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "Profile to be created",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("spec called")
	},
}

func init() {
	cmd.RootCmd.AddCommand(specCmd)
}
