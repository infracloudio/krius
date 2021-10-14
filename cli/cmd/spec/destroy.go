package spec

import (
	"io/ioutil"
	"log"

	client "github.com/infracloudio/krius/pkg/client"
	spec "github.com/infracloudio/krius/pkg/specvalidate"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var uninstallSpecCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy Krius Managed stack across clusters defined in the config file",
	RunE:  uninstallSpec,
}

func init() {
	specCmd.AddCommand(uninstallSpecCmd)
	err := addSpecApplyFlags(uninstallSpecCmd)
	if err != nil {
		log.Printf("Error adding flags: %v", err)
	}
}

func uninstallSpec(cmd *cobra.Command, args []string) (err error) {
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

	for _, cluster := range config.Clusters {
		switch cluster.Type {
		case "prometheus":
			pc, err := client.NewPromClient(&cluster)
			if err != nil {
				return err
			}
			err = pc.UninstallClient(cluster.Name)
			if err != nil {
				log.Println(err)
			}
		case "thanos":
			tc, err := client.NewThanosClient(&cluster)
			if err != nil {
				return err
			}
			err = tc.UninstallClient(cluster.Name)
			if err != nil {
				log.Println(err)
			}
		case "grafana":
			log.Println("Grafana uninstall to be implemented")
		}
	}
	return nil
}
