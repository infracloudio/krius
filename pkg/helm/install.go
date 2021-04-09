package helm

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/strvals"
)

func (c *HelmClient) InstallOrUpgradeChart() (*string, error) {
	if c.Client.Version == "" && c.Client.Devel {
		c.Client.Version = ">0.0.0-0"
	}

	c.Client.ReleaseName = c.ReleaseName
	cp, err := c.Client.ChartPathOptions.LocateChart(fmt.Sprintf("%s/%s", c.RepoName, c.ChartName), c.Settings)
	if err != nil {
		return nil, err
	}

	debug("CHART PATH: %s\n", cp)

	p := getter.All(c.Settings)
	valueOpts := &values.Options{}
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return nil, err
	}

	// Add args
	if err := strvals.ParseInto(c.Args["set"], vals); err != nil {
		m := errors.Wrap(err, "failed parsing --set data")
		return nil, m
	}

	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		return nil, err
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if c.Client.DependencyUpdate {
				man := &downloader.Manager{
					Out:              os.Stdout,
					ChartPath:        cp,
					Keyring:          c.Client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: c.Settings.RepositoryConfig,
					RepositoryCache:  c.Settings.RepositoryCache,
				}
				if err := man.Update(); err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}

	c.Client.Namespace = c.Settings.Namespace()
	release, err := c.Client.Run(chartRequested, vals)
	if err != nil {
		return nil, err
	}
	return &release.Manifest, nil
}

func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

func debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	log.Output(2, fmt.Sprintf(format, v...))
}
