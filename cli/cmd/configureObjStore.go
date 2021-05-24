package cmd

import (
	"fmt"
	"log"
	"strings"

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
			errors := []string{}
			if storeTypeFlag, _ := cmd.Flags().GetString(storageType); storeTypeFlag == "" {
				errors = append(errors, storageType)
			}
			if accessKeyFlag, _ := cmd.Flags().GetString(accessKey); accessKeyFlag == "" {
				errors = append(errors, accessKey)
			}
			if secretKeyFlag, _ := cmd.Flags().GetString(secretKey); secretKeyFlag == "" {
				errors = append(errors, secretKey)
			}
			if endpointFlag, _ := cmd.Flags().GetString(endpoint); endpointFlag == "" {
				errors = append(errors, endpoint)
			}
			if bucketFlag, _ := cmd.Flags().GetString(bucket); bucketFlag == "" {
				errors = append(errors, bucket)
			}
			if releaseFlag, _ := cmd.Flags().GetString(release); releaseFlag == "" {
				errors = append(errors, release)
			}
			if namespaceFlag, _ := cmd.Flags().GetString(namespace); namespaceFlag == "" {
				errors = append(errors, namespace)
			}
			if len(errors) > 0 {
				formattedError := strings.Join(errors, ",")
				return fmt.Errorf("Flag missing - please set %s", formattedError)
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
		log.Fatalf("Error creating helm client: %v", err)
		return
	}

	results, err := helmClient.ListDeployedReleases()
	if err != nil {
		log.Fatalf("helm list error: %v", err)
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
			log.Println(*result)
		}
	} else {
		fmt.Print("Release name does not exist")
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
