package spec

import "github.com/spf13/cobra"

func addSpecApplyFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("config-file", "c", "", "config file path")
	cmd.MarkFlagRequired("config-file")
}
