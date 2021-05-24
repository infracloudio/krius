package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var promConfigCmd = &cobra.Command{
	Use:   "prometheus [Name]",
	Short: "Configure prometheus",
	Run:   configurePrometheus,
}

func init() {
	configureCmd.AddCommand(promConfigCmd)
}

func configurePrometheus(cmd *cobra.Command, args []string) {
	fmt.Println("configure prometheus")
}
