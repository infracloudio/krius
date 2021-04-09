package helm

import (
	"os/exec"

	"github.com/spf13/cobra"
)

func HelmRepoAdd(name, url string) {
	exec.Command("helm", "repo", "add", name, url)
}

func AddInstallFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("namespace", "n", "default", "namespace in which the chart need to be installed")
}
