package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure the give component",
	Run:   configureCluster,
}

func init() {
	RootCmd.AddCommand(configureCmd)
}

func configureCluster(cmd *cobra.Command, args []string) {
	fmt.Println("Need to implement configuration")
}
