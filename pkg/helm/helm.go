package helm

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"helm.sh/helm/v3/pkg/strvals"
)

// AddRepo adds repo with given name and url
func (client *Client) AddRepo() error {
	repoFile := client.Settings.RepositoryConfig

	//Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(repoFile, filepath.Ext(repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer func() {
			err = fileLock.Unlock()
			if err != nil {
				log.Println(err)
			}
		}()
	}
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadFile(repoFile)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		log.Fatal(err)
	}

	if f.Has(client.RepoName) {
		fmt.Printf("repository name (%s) already exists\n", client.RepoName)
		return nil
	}

	c := repo.Entry{
		Name: client.RepoName,
		URL:  client.URL,
	}

	r, err := repo.NewChartRepository(&c, getter.All(client.Settings))
	if err != nil {
		return err
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		err := errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", client.URL)
		return err
	}

	f.Update(&c)

	if err := f.WriteFile(repoFile, 0644); err != nil {
		return err
	}
	fmt.Printf("%q has been added to your repositories\n", client.RepoName)
	return nil
}

func (client *Client) UpdateRepo() error {
	repoFile := client.Settings.RepositoryConfig

	f, err := repo.LoadFile(repoFile)
	if os.IsNotExist(errors.Cause(err)) || len(f.Repositories) == 0 {
		return errors.New("no repositories found. You must add one before updating")
	}
	var repos []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(client.Settings))
		if err != nil {
			return err
		}
		repos = append(repos, r)
	}

	fmt.Printf("Hang tight while we grab the latest from your chart repositories...\n")
	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if _, err := re.DownloadIndexFile(); err != nil {
				fmt.Printf("...Unable to get an update from the %q chart repository (%s):\n\t%s\n", re.Config.Name, re.Config.URL, err)
			} else {
				fmt.Printf("...Successfully got an update from the %q chart repository\n", re.Config.Name)
			}
		}(re)
	}
	wg.Wait()
	fmt.Printf("Update Complete. ⎈ Happy Helming!⎈\n")
	return nil
}

func (c *Client) ListDeployedReleases() ([]*release.Release, error) {
	listClient := action.NewList(c.ActionConfig)
	return listClient.Run()
}

func (c *Client) InstallChart(valueOpts *values.Options) (*string, error) {
	client := action.NewInstall(c.ActionConfig)

	if client.Version == "" && client.Devel {
		client.Version = ">0.0.0-0"
	}

	if c.ReleaseName != "" {
		client.ReleaseName = c.ReleaseName
	}

	// Generate Random name for the release
	client.GenerateName = true
	client.ReleaseName, _, _ = client.NameAndChart([]string{c.ChartName})
	cp, err := client.ChartPathOptions.LocateChart(fmt.Sprintf("%s/%s", c.RepoName, c.ChartName), c.Settings)
	if err != nil {
		return nil, err
	}

	debug("CHART PATH: %s\n", cp)
	if valueOpts == nil {
		valueOpts = &values.Options{}
	}
	p := getter.All(c.Settings)
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
			if client.DependencyUpdate {
				man := &downloader.Manager{
					Out:              os.Stdout,
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
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

	client.Namespace = c.Settings.Namespace()
	release, err := client.Run(chartRequested, vals)
	if err != nil {
		return nil, err
	}
	return &release.Manifest, nil
}

func (c *Client) UpgradeChart(valueOpts *values.Options) (*string, error) {
	client := action.NewUpgrade(c.ActionConfig)

	if client.Version == "" && client.Devel {
		client.Version = ">0.0.0-0"
	}
	cp, err := client.ChartPathOptions.LocateChart(fmt.Sprintf("%s/%s", c.RepoName, c.ChartName), c.Settings)
	if err != nil {
		return nil, err
	}

	debug("CHART PATH: %s\n", cp)
	p := getter.All(c.Settings)
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

	if req := chartRequested.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			return nil, err
		}
	}
	release, err := client.Run(c.ReleaseName, chartRequested, vals)
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
	err := log.Output(2, fmt.Sprintf(format, v...))
	if err != nil {
		log.Printf("Error while logging: %v", err)
	}
}

func InitializeHelmAction(settings *cli.EnvSettings) (*action.Configuration, error) {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(),
		os.Getenv("HELM_DRIVER"), debug); err != nil {
		return nil, err
	}
	return actionConfig, nil
}
func NewClientFromKubeConf(options *KubeConfClientOptions, settings *cli.EnvSettings) (*action.Configuration, error) {
	if options.KubeContext != "" {
		settings.KubeContext = options.KubeContext
	}
	return InitializeHelmAction(settings)
}
