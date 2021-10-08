package client

import (
	"github.com/infracloudio/krius/pkg/helm"
)

func ChartStatusCheck(clusterName string, namespace string, chartName string) (status string, err error) {

	chartConfiguration := &helm.Config{
		Name: chartName,
	}

	helmClient, err := createHelmClientObject(clusterName, namespace, chartConfiguration)
	if err != nil {
		return "", err
	}
	status, err = helmClient.StatusHelmChart(chartName)
	if err != nil {
		return "", err
	}
	return status, err
}
