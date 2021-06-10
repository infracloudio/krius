package spec

import (
	"fmt"
	"io/ioutil"
	"log"

	spec "github.com/infracloudio/krius/pkg/specvalidate"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	configFile = "config-file"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies/Updates the give profie file",
	Run:   applySpec,
}

func init() {
	specCmd.AddCommand(applyCmd)
	addSpecApplyFlags(applyCmd)
}

func applySpec(cmd *cobra.Command, args []string) {
	configFileFlag, _ := cmd.Flags().GetString(configFile)
	loader, ruleSchemaLoader, err := spec.GetLoaders(configFileFlag)
	if err != nil {
		log.Println(err)
		return
	}
	valid, errors := spec.ValidateYML(loader, ruleSchemaLoader)
	if !valid {
		log.Println(errors)
	}
	log.Println("valid yaml")

	yamlFile, err := ioutil.ReadFile(configFileFlag)
	if err != nil {
		log.Fatalf("yamlFile.Get err #%v ", err)
	}
	var config Config
	yaml.Unmarshal(yamlFile, &config)
	for _, cluster := range config.Clusters {
		c, err := cluster.GetConfig(cluster.Type, cluster.Data)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(c)

		//switch on kind and pass the prom/thanos/grafana to respective files
	}

	// This is for testing only to create a Helm client using kube context

	// helmClient, err := createHelmClientObject("kind-cluster1")
	// if err != nil {
	// 	return
	// }
	// results, err := helmClient.ListDeployedReleases()
	// if err != nil {
	// 	log.Fatalf("helm list error: %v", err)
	// 	return
	// }
	// for _, v := range results {
	// 	fmt.Println("v.Name", v.Name)
	// }
}
