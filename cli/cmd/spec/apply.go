package spec

import (
	"io/ioutil"
	"log"
	"sort"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	client "github.com/infracloudio/krius/pkg/client"
	spec "github.com/infracloudio/krius/pkg/specvalidate"
)

const (
	configFile = "config-file"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies/Updates the give profie file",
	RunE:  applySpec,
}

func init() {
	specCmd.AddCommand(applyCmd)
	err := addSpecApplyFlags(applyCmd)
	if err != nil {
		log.Printf("Error adding flags: %v", err)
	}
}

func addSpecApplyFlags(cmd *cobra.Command) error {
	cmd.Flags().StringP("config-file", "c", "", "config file path")
	err := cmd.MarkFlagRequired("config-file")
	return err
}

func applySpec(cmd *cobra.Command, args []string) (err error) {
	configFileFlag, _ := cmd.Flags().GetString(configFile)
	loader, ruleSchemaLoader, err := spec.GetLoaders(configFileFlag)
	if err != nil {
		return err
	}
	valid, errors := spec.ValidateYML(loader, ruleSchemaLoader)
	if !valid {
		log.Println(errors)
		return
	}
	log.Println("valid yaml")
	yamlFile, err := ioutil.ReadFile(configFileFlag)
	if err != nil {
		log.Fatalf("yamlFile.Get err #%v ", err)
	}
	var config client.Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return err
	}
	preFlightErrors := []string{}
	// check for preflight errors for all the clusters
	for _, cluster := range config.Clusters {
		switch cluster.Type {
		case "prometheus":
			pc, err := client.NewPromClient(&cluster)
			if err != nil {
				return err
			}
			clusterErrors, err := pc.PreflightChecks(&config, cluster.Name)
			if err != nil {
				return err
			}
			if clusterErrors != nil {
				preFlightErrors = append(preFlightErrors, clusterErrors...)
			}
		case "thanos":
			tc, err := client.NewThanosClient(&cluster)
			if err != nil {
				return err
			}
			clusterErrors, err := tc.PreflightChecks(&config, cluster.Name)
			if err != nil {
				return err
			}
			if clusterErrors != nil {
				preFlightErrors = append(preFlightErrors, clusterErrors...)
			}
		case "grafana":
		}
	}
	if len(preFlightErrors) > 0 {
		log.Printf("Preflight checks failed %s", preFlightErrors)
		return
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
	log.Println("Preflight checks passed")
	var targets []string
	var receiveEndpoints []string
	for _, cluster := range config.Clusters {
		switch cluster.Type {
		case "prometheus":
			pc, err := client.NewPromClient(&cluster)
			if err != nil {
				return err
			}
			target, err := pc.InstallClient(cluster.Name, receiveEndpoints)
			if err != nil {
				return err
			}
			targets = append(targets, target+":10901")
		case "thanos":
			tc, err := client.NewThanosClient(&cluster)
			if err != nil {
				return err
			}
			endpoint, err := tc.InstallClient(cluster.Name, targets)
			if err != nil {
				return err
			}
			receiveEndpoints = append(receiveEndpoints, endpoint)

		case "grafana":
		}
	}
	return nil
}
