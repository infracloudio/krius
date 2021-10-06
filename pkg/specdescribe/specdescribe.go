package specdescribe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	client "github.com/infracloudio/krius/pkg/client"
	spec "github.com/infracloudio/krius/pkg/specvalidate"
	"github.com/spf13/cobra"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
)

var (
	describeConfig client.Config
)

const (
	configFile = "config-file"
)

func DescribeClusterKrius(cmd *cobra.Command, args []string) (err error) {

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

	getSpecFile, err := ioutil.ReadFile(configFileFlag)
	if err != nil {
		log.Fatal("Parse Error: Unable to read the spec file ")
		return err
	}

	jsonSpecFile, err := yamlutil.ToJSON(getSpecFile)
	if err != nil {
		log.Fatal("Parse Error: Unable to conver the spec file to standard JSON format")
		return err
	}

	err = json.Unmarshal(jsonSpecFile, &describeConfig)
	if err != nil {
		log.Fatal("Parse Error: Unable to parse spec")
		return err
	}

	for _, each := range describeConfig.Clusters {
		fmt.Print("\n---------------------------------------------------------------------------")
		fmt.Print("\n Kubernetes Cluster Context: ", each.Name)
		fmt.Print("\n Krius Cluster")
		fmt.Print("\n - Name: ", each.Data["name"])
		fmt.Print("\n - Namespace: ", each.Data["namespace"])
		fmt.Print("\n - Type: ", each.Type)
		fmt.Print("\n - ObjectConfiguration Name: ", each.Data["objStoreConfig"])
		fmt.Print("\n---------------------------------------------------------------------------")

	}
	return err
}
