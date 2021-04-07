package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "krius",
	Short: "A tool to setup Prometheus, Thanos",
	Long:  `A tool to setup Prometheus, Thanos & friends across multiple clusters easily for scale .`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
