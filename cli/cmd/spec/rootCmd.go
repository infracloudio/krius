package spec

import (
	"fmt"
	"log"

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

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies/Updates the give profie file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return manageApp(cmd)
	},
}

var uninstallSpecCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy the entire stack across clusters",
	RunE: func(cmd *cobra.Command, args []string) error {
		return manageApp(cmd)
	},
}

func init() {
	cmd.RootCmd.AddCommand(specCmd)
	specCmd.AddCommand(applyCmd)
	specCmd.AddCommand(uninstallSpecCmd)
	applyCmd.Flags().StringP("config-file", "c", "", "config file path")
	uninstallSpecCmd.Flags().StringP("config-file", "c", "", "config file path")
	applyCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "enable debug logs")

	err := applyCmd.MarkFlagRequired("config-file")
	if err != nil {
		log.Print(err)
	}
	err = uninstallSpecCmd.MarkFlagRequired("config-file")
	if err != nil {
		log.Print(err)
	}
}
