package helm

import (
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
)

type HelmClient struct {
	Url          string
	RepoName     string
	ChartName    string
	ReleaseName  string
	Namespace    string
	Args         map[string]string
	Settings     *cli.EnvSettings
	ActionConfig *action.Configuration
}
type HelmConfig struct {
	Repo       string
	Name       string
	Url        string
	Args       []string
	Cmd        *cobra.Command
	ValuesYaml string
	ValueOpts  *values.Options
}
