package helm

import (
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
)

type Client struct {
	URL          string
	RepoName     string
	ChartName    string
	ReleaseName  string
	Namespace    string
	Args         map[string]string
	Settings     *cli.EnvSettings
	ActionConfig *action.Configuration
}
type Config struct {
	Repo       string
	Name       string
	URL        string
	Args       []string
	Cmd        *cobra.Command
	ValuesYaml string
	ValueOpts  *values.Options
}

type KubeConfClientOptions struct {
	*Options
	KubeContext string
	KubeConfig  []byte
}
type Options struct {
	Namespace        string
	RepositoryConfig string
	RepositoryCache  string
	Debug            bool
	Linting          bool
}
