package cmd

import (
	"fmt"
	"log"

	"github.com/infracloudio/krius/pkg/helm"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/cli/values"
)

var objStoreCmd = &cobra.Command{
	Use:   "objstore [Name]",
	Short: "Configure object storage",
	Run:   configureObjStore,
}

func init() {
	configureCmd.AddCommand(objStoreCmd)
	addConfigureObjStoreFlags(objStoreCmd)
}

func debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	log.Output(2, fmt.Sprintf(format, v...))
}

func configureObjStore(cmd *cobra.Command, args []string) {
	chartConfiguration := &helm.HelmConfig{
		Repo: "bitnami",
		Name: "thanos",
		Url:  "https://charts.bitnami.com/bitnami",
		Args: args,
		Cmd:  cmd,
	}

	helmClient, err := createHelmClientObject(chartConfiguration)
	if err != nil {
		log.Fatal(err)
		return
	}

	release, err := cmd.Flags().GetString("release")
	if err != nil {
		log.Fatal(err)
		return
	}
	results, err := helmClient.ListDeployedReleases()
	if err != nil {
		log.Fatal(err)
		return
	}
	exists := false
	for _, v := range results {
		if v.Name == release {
			exists = true
		}
	}
	if exists {
		configPath := getVarFromCmd(cmd, "config-file", "bucket.yaml")
		Values := &values.Options{}
		if configPath != "" {
			Values.ValueFiles = []string{configPath}
		} else {
			Values = createObjStoreValuesMap(cmd, Values)
		}
		result, err := helmClient.UpgradeChart(Values)
		if err != nil {
			fmt.Printf("could not upgrade The Observability Stack %s", err)
		} else {
			fmt.Println(*result)
		}
	}
}

func createObjStoreValuesMap(cmd *cobra.Command, valueOpts *values.Options) *values.Options {
	storageType := getVarFromCmd(cmd, "type", "")
	accessKey := getVarFromCmd(cmd, "access_key", "")
	secretKey := getVarFromCmd(cmd, "secret_key", "")
	endpoint := getVarFromCmd(cmd, "endpoint", "")
	bucket := getVarFromCmd(cmd, "bucket", "")
	valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("objstoreConfig.type=%s", storageType))
	valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("objstoreConfig.config.bucket=%s", bucket))
	valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("objstoreConfig.config.endpoint=%s", endpoint))
	valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("objstoreConfig.config.access_key=%s", accessKey))
	valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("objstoreConfig.config.secret_key=%s", secretKey))
	valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("bucketweb.enabled=%s", "true"))
	return valueOpts
}
