package cmd

import (
	"fmt"
	"log"

	"github.com/infracloudio/krius/pkg/helm"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/cli/values"
)

const (
	storageType = "type"
	accessKey   = "access_key"
	secretKey   = "secret_key"
	endpoint    = "endpoint"
	bucket      = "bucket"
	configFile  = "config-file"
	release     = "release"
	namespace   = "namespace"
)

var objStoreCmd = &cobra.Command{
	Use:   "objstore [Name]",
	Short: "Configure object storage",
	Run:   configureObjStore,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		configFileFlag, _ := cmd.Flags().GetString(configFile)
		if configFileFlag == "" {
			if storeTypeFlag, _ := cmd.Flags().GetString(storageType); storeTypeFlag == "" {
				return fmt.Errorf("Flag missing - please set %s", storageType)
			}
			if accessKeyFlag, _ := cmd.Flags().GetString(accessKey); accessKeyFlag == "" {
				return fmt.Errorf("Flag missing - please set %s", accessKey)
			}
			if secretKeyFlag, _ := cmd.Flags().GetString(secretKey); secretKeyFlag == "" {
				return fmt.Errorf("Flag missing - please set %s", secretKey)
			}
			if endpointFlag, _ := cmd.Flags().GetString(endpoint); endpointFlag == "" {
				return fmt.Errorf("Flag missing - please set %s", endpoint)
			}
			if bucketFlag, _ := cmd.Flags().GetString(bucket); bucketFlag == "" {
				return fmt.Errorf("Flag missing - please set %s", bucket)
			}
			if releaseFlag, _ := cmd.Flags().GetString(release); releaseFlag == "" {
				return fmt.Errorf("Flag missing - please set %s", release)
			}
			if namespaceFlag, _ := cmd.Flags().GetString(namespace); namespaceFlag == "" {
				return fmt.Errorf("Flag missing - please set %s", namespace)
			}
		}
		return nil
	},
}

func init() {
	configureCmd.AddCommand(objStoreCmd)
	addConfigureObjStoreFlags(objStoreCmd)
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

	results, err := helmClient.ListDeployedReleases()
	if err != nil {
		log.Fatal(err)
		return
	}
	exists := false
	for _, v := range results {
		if v.Name == helmClient.ReleaseName {
			exists = true
		}
	}
	if exists {
		Values := &values.Options{}
		configPath, _ := cmd.Flags().GetString(configFile)
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
	storageType, _ := cmd.Flags().GetString(storageType)
	accessKey, _ := cmd.Flags().GetString(accessKey)
	secretKey, _ := cmd.Flags().GetString(secretKey)
	endpoint, _ := cmd.Flags().GetString(endpoint)
	bucket, _ := cmd.Flags().GetString(bucket)
	valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("objstoreConfig.type=%s", storageType))
	valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("objstoreConfig.config.bucket=%s", bucket))
	valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("objstoreConfig.config.endpoint=%s", endpoint))
	valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("objstoreConfig.config.access_key=%s", accessKey))
	valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("objstoreConfig.config.secret_key=%s", secretKey))
	valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("bucketweb.enabled=%s", "true"))
	return valueOpts
}
