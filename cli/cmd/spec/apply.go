package spec

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	client "github.com/infracloudio/krius/pkg/client"
	spec "github.com/infracloudio/krius/pkg/specvalidate"
)

const (
	configFile = "config-file"
)

func (r *AppRunner) applySpec(cmd *cobra.Command) (err error) {
	configFileFlag, _ := cmd.Flags().GetString(configFile)
	yamlFile, err := ioutil.ReadFile(configFileFlag)
	if err != nil {
		r.log.Error(err)
		return
	}
	r.status.Start("validating yaml")
	time.Sleep(1 * time.Second)
	loader, ruleSchemaLoader, err := spec.GetLoaders(configFileFlag)
	if err != nil {
		return err
	}
	valid, errors := spec.ValidateYML(loader, ruleSchemaLoader)
	if !valid {
		errs := []string{}
		for _, desc := range errors {
			errs = append(errs, desc.String())
		}
		r.status.Error("validating yaml: " + strings.Join(errs, ", "))
		return
	}
	r.status.Success()

	var config client.Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		r.log.Error(err)
		return
	}
	// check for preflight errors for all the clusters
	for _, cluster := range config.Clusters {
		r.status.Start(fmt.Sprintf("Preflight error checks in cluster %s", cluster.Name))
		switch cluster.Type {
		case "prometheus":
			pc, err := client.NewPromClient(&cluster)
			if err != nil {
				r.log.Error(err)
				return err
			}
			clusterErrors, err := pc.PreflightChecks(&config, cluster.Name)
			if err != nil {
				r.log.Error(err)
				return err
			}
			if len(clusterErrors) > 0 {
				r.status.Error(strings.Join(clusterErrors, ", "))
				return err
			}
			r.status.Success()
			r.status.Stop()

		case "thanos":
			tc, err := client.NewThanosClient(&cluster)
			if err != nil {
				r.log.Error(err)
				return err
			}
			clusterErrors, err := tc.PreflightChecks(&config, cluster.Name)
			if err != nil {
				r.log.Error(err)
				return err
			}
			if len(clusterErrors) > 0 {
				r.status.Error(strings.Join(clusterErrors, ", "))
				return err
			}
			r.status.Success()
			r.status.Stop()
		case "grafana":
		}
	}

	// reorder clusters based on sidecar/receiver setup
	if config.Order == 1 {
		sort.Slice(config.Clusters, func(p, q int) bool {
			return config.Clusters[p].Type < config.Clusters[q].Type
		})
	} else if config.Order == 2 {
		sort.Slice(config.Clusters, func(p, q int) bool {
			return config.Clusters[p].Type > config.Clusters[q].Type
		})
	}
	var targets []string
	var receiveEndpoints []string
	for _, cluster := range config.Clusters {
		m := fmt.Sprintf("🚀 Installing %s stack in cluster %s", cluster.Type, cluster.Name)
		if r.log.DebugLevel {
			r.log.Infof(m)
		} else {
			r.status.Start(m)
		}
		switch cluster.Type {
		case "prometheus":
			pc, err := client.NewPromClient(&cluster)
			if err != nil {
				r.status.Error()
				return err
			}
			target, err := pc.InstallClient(cluster.Name, receiveEndpoints, r.status.logger.DebugLevel)
			if err != nil {
				r.status.Error()
				return err
			}
			targets = append(targets, target+":10901")
			r.status.Success()
			r.status.Stop()
		case "thanos":
			tc, err := client.NewThanosClient(&cluster)
			if err != nil {
				r.status.Error()
				return err
			}
			endpoint, err := tc.InstallClient(cluster.Name, targets, r.status.logger.DebugLevel)
			if err != nil {
				r.status.Error()
				return err
			}
			receiveEndpoints = append(receiveEndpoints, endpoint)
			r.status.Success()
			r.status.Stop()
		case "grafana":
		}
	}
	return nil
}
