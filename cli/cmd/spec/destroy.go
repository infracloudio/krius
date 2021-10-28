package spec

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	client "github.com/infracloudio/krius/pkg/client"
	spec "github.com/infracloudio/krius/pkg/specvalidate"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func (r *AppRunner) uninstallSpec(cmd *cobra.Command) (err error) {
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
	r.status.Stop()

	var config client.Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return err
	}

	for _, cluster := range config.Clusters {
		m := fmt.Sprintf("🧹 Uninstalling %s stack in cluster %s and its dependencies...", cluster.Type, cluster.Name)
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
			}
			err = pc.UninstallClient(cluster.Name)
			if err != nil {
				r.status.Error()
				r.log.Errorf(err.Error())
			} else {
				r.status.Success()
				r.status.Stop()
			}
		case "thanos":
			tc, err := client.NewThanosClient(&cluster)
			if err != nil {
				r.status.Error()
				return err
			}
			err = tc.UninstallClient(cluster.Name)
			if err != nil {
				r.status.Error()
				r.log.Errorf(err.Error())
			} else {
				r.status.Success()
				r.status.Stop()
			}
		case "grafana":
			log.Println("Grafana uninstall to be implemented")
		}
	}
	return nil
}
