package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
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

func configureObjStore(cmd *cobra.Command, args []string) {
	chartConfiguration := &ChartConfig{
		CHART_REPO: "bitnami",
		CHART_NAME: "thanos",
		CHART_URL:  "https://charts.bitnami.com/bitnami",
	}

	helmClient, err := createHelmUpgradeClientObject(cmd, args, chartConfiguration)
	if err != nil {
		log.Fatal(err)
	}
	result, err := helmClient.UpgradeChartValues()
	if err != nil {
		fmt.Printf("could not upgrade The Observability Stack %s", err)
	} else {
		fmt.Println(*result)
	}
}
