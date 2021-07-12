package cmd

import (
	"fmt"

	"github.com/infracloudio/krius/pkg/helm"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var thanosCmd = &cobra.Command{
	Use:   "thanos",
	Short: "Install thanos component",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("missing argument to configure thanos as: sidecar, receiver\n\n For example: krius install thanos sidecar")
		}
		return nil
	},
	Run: thanosInstall,
}

func init() {
	installCmd.AddCommand(thanosCmd)
	addInstallFlags(thanosCmd)
}

func thanosInstall(cmd *cobra.Command, args []string) {
	helmConfiguration := &helm.HelmConfig{
		Repo: "bitnami",
		Name: "thanos",
		Url:  "https://charts.bitnami.com/bitnami",
	}
	fmt.Printf("Need to implement thanos %s, %s and %s", helmConfiguration.Name, helmConfiguration.Repo, helmConfiguration.Url)
}
