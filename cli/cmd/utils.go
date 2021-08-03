package cmd

import (
	"github.com/spf13/cobra"
)

func addInstallFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("namespace", "n", "default", "namespace in which the chart need to be installed")
	cmd.Flags().StringP("release", "r", "default", "release name to be used for the specific install")
}

// unused for now
// func getVarFromCmd(cmd *cobra.Command, envVar, defaultValue string) string {
// 	envVar, err := cmd.Flags().GetString(envVar)
// 	if err != nil {
// 		envVar = defaultValue
// 	}
// 	return envVar
// }

func addConfigureObjStoreFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("config-file", "c", "", "config file path")
	cmd.Flags().StringP("namespace", "n", "", "namespace in which the chart has installed")
	cmd.Flags().StringP("release", "r", "", "release name of the chart")
	cmd.Flags().StringP("type", "t", "", "type of storage")
	cmd.Flags().StringP("bucket", "b", "", "bucket name")
	cmd.Flags().StringP("endpoint", "e", "", "bucket's endpoint")
	cmd.Flags().StringP("access_key", "a", "", "access key")
	cmd.Flags().StringP("secret_key", "s", "", "secret key")
}
