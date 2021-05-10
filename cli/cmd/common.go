package cmd

import (
	"log"
	"os"

	"github.com/infracloudio/krius/pkg/helm"
	"helm.sh/helm/v3/pkg/cli"
)

func createHelmClientObject(helmConfig *helm.HelmConfig) (*helm.HelmClient, error) {
	namespace, err := helmConfig.Cmd.Flags().GetString("namespace")
	if err != nil {
		namespace = "default"
	}
	releaseName, err := helmConfig.Cmd.Flags().GetString("release")
	if err != nil {
		releaseName = "my-release"
	}
	os.Setenv("HELM_NAMESPACE", namespace)
	settings = cli.New()

	action, err := helm.InitializeHelmAction(settings)
	if err != nil {
		log.Fatal(err)
	}
	helmClient := helm.HelmClient{
		RepoName:     helmConfig.Repo,
		Url:          helmConfig.Url,
		ReleaseName:  releaseName,
		Namespace:    namespace,
		ChartName:    helmConfig.Name,
		ActionConfig: action,
		Settings:     settings,
	}
	return &helmClient, err
}
