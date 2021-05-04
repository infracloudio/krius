package cmd

import (
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure the give component",
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	RootCmd.AddCommand(configureCmd)
}
