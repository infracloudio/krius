package spec

import (
	"log"

	"github.com/infracloudio/krius/pkg/specdescribe"
	"github.com/spf13/cobra"
)

var describeClusterCmd = &cobra.Command{
	Use:   "describe-cluster",
	Short: "Describes the entire stack across multiple clusters and current state",
	RunE:  specdescribe.DescribeClusterKrius,
}

func addDescribeConfigFileFlags(cmd *cobra.Command) error {
	cmd.Flags().StringP("config-file", "c", "", "config file path")
	err := cmd.MarkFlagRequired("config-file")
	return err
}

func init() {
	specCmd.AddCommand(describeClusterCmd)
	err := addDescribeConfigFileFlags(describeClusterCmd)
	if err != nil {
		log.Print("Error Adding Config Flag: ", err)
	}
}
