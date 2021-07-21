package spec

import (
	"fmt"
	"strings"

	"helm.sh/helm/v3/pkg/cli/values"
)

func createSidecarValuesMap(secretName string) *values.Options {
	var valueOpts values.Options
	valueOpts.Values = []string{fmt.Sprintf("prometheus.prometheusSpec.thanos.image=%s", "thanosio/thanos:v0.21.0-rc.0"),
		fmt.Sprintf("prometheus.prometheusSpec.thanos.sha=%s", "dbf064aadd18cc9e545c678f08800b01a921cf6817f4f02d5e2f14f221bee17c"),
		fmt.Sprintf("prometheus.thanosService.enabled=%s", "true"),
		fmt.Sprintf("prometheus.thanosServiceExternal.enabled=%s", "true"),
		fmt.Sprintf("prometheus.prometheusSpec.thanos.objectStorageConfig.name=%s", secretName),
		fmt.Sprintf("prometheus.prometheusSpec.thanos.objectStorageConfig.key=%s", "objstore.yml")}
	return &valueOpts
}

func createThanosValuesMap(thanos Thanos) *values.Options {
	var valueOpts values.Options
	targets := "{" + strings.Join(thanos.Querier.Targets, ",") + "}"
	extraFlags := "{"
	if thanos.Querier.AutoDownsample {
		extraFlags += "--query.auto-downsampling,"
	}
	if thanos.Querier.PartialResponse {
		extraFlags += "--query.partial-response"
	}

	extraFlags += "}"
	valueOpts.Values = []string{
		fmt.Sprintf("existingObjstoreSecret=%s", thanos.ObjStoreConfig),
		fmt.Sprintf("query.stores=%s", targets),
		fmt.Sprintf("storegateway.enabled=%s", "true"),
		fmt.Sprintf("query.extraFlags=%s", extraFlags),
		fmt.Sprintf("queryFrontend.enabled=%s", "true")}

	if thanos.Compactor.Name != "" {
		valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("compactor.enabled=%s", "true"))
		// valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("compactor.podLabels=%s", "{'key':'"+thanos.Compactor.Name+"'}"))
	}
	if thanos.Querierfe.Name != "" {
		valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("queryFrontend.enabled=%s", "true"))
		// valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("queryFrontend.podLabels=%s", "{'key':'"+thanos.Querier.Name+"'}"))
	}
	return &valueOpts
}
