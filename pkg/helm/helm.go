package helm

import (
	"os/exec"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
)

type HelmClient struct {
	Url         string
	RepoName    string
	ChartName   string
	ReleaseName string
	Namespace   string
	Args        map[string]string
	Client      *action.Install
	Settings    *cli.EnvSettings
}

func HelmRepoAdd(name, url string) {
	exec.Command("helm", "repo", "add", name, url)
}

func AddInstallFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("namespace", "n", "default", "namespace in which the chart need to be installed")
}
