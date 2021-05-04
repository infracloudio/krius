package cmd

import (
	"github.com/spf13/cobra"
)

func addInstallFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("namespace", "n", "default", "namespace in which the chart need to be installed")
}
