package client

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/infracloudio/krius/pkg/helm"
	kube "github.com/infracloudio/krius/pkg/kubeClient"
	"github.com/infracloudio/krius/pkg/utils"

	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
)

var settings *cli.EnvSettings

func createHelmClientObject(context, namespace string, debug bool, helmConfig *helm.Config) (*helm.Client, error) {
	opt := &helm.KubeConfClientOptions{
		KubeContext: context,
	}
	os.Setenv("HELM_NAMESPACE", namespace)
	helm.Settings.Debug = debug
	settings = cli.New()
	action, err := helm.NewClientFromKubeConf(opt, settings)
	if err != nil {
		return nil, err
	}
	helmClient := helm.Client{
		ActionConfig: action,
		Settings:     settings,
		RepoName:     helmConfig.Repo,
		URL:          helmConfig.URL,
		ChartName:    helmConfig.Name,
	}
	return &helmClient, nil
}
func GetKubeClient(namespace, context string) (*kube.KubeConfig, error) {
	kubeClient := kube.KubeConfig{
		Namespace: namespace,
		Context:   context,
	}
	err := kubeClient.InitClient()
	if err != nil {
		return nil, err
	}
	return &kubeClient, nil
}
func getPrometheusTargets(clusterName, namespace, promName string) []string {
	kubeClient, err := GetKubeClient(namespace, clusterName)
	if err != nil {
		return nil
	}
	m := make(map[string]string)
	m["app"] = "krius"
	m["prometheus"] = "sidecar"
	return kubeClient.GetServiceInfoByLabels(m)
}

func getReceiveEndpoint(clusterName, namespace, specName string) []string {
	kubeClient, err := GetKubeClient(namespace, clusterName)
	if err != nil {
		return nil
	}
	return kubeClient.GetServiceInfo(specName + "-receive")
}

func (p Prometheus) createPrometheusSidecarValues() *values.Options {
	var valueOpts values.Options
	valueOpts.Values = []string{
		fmt.Sprintf("commonLabels.replica=%s", p.Name),
		fmt.Sprintf("prometheus.prometheusSpec.thanos.image=%s", "quay.io/thanos/thanos:v0.28.1"),
		fmt.Sprintf("prometheus.thanosService.enabled=%s", "true"),
		fmt.Sprintf("prometheus.thanosServiceExternal.enabled=%s", "true"),
		fmt.Sprintf("prometheus.thanosServiceExternal.labels.app=%s", "krius"),
		fmt.Sprintf("prometheus.thanosServiceExternal.labels.prometheus=%s", "sidecar"),
		fmt.Sprintf("prometheus.prometheusSpec.thanos.objectStorageConfig.name=%s", p.ObjStoreConfig),
		fmt.Sprintf("prometheus.prometheusSpec.thanos.objectStorageConfig.key=%s", "objstore.yml")}
	return &valueOpts
}

func (p Prometheus) createPrometheusReceiverValues() *values.Options {
	var valueOpts values.Options

	valueOpts.Values = append(valueOpts.Values,
		fmt.Sprintf("commonLabels.replica=%s", p.Name))
	if len(p.RemoteWriteURL) > 0 && p.RemoteWriteURL[0] != "" {
		valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("prometheus.prometheusSpec.remoteWrite[0].url=http://%s:19291/api/v1/receive", p.RemoteWriteURL[0]))
	}
	return &valueOpts
}
func (thanos Thanos) createThanosValuesMap() (*values.Options, error) {
	var valueOpts values.Options
	hash := utils.RandStringRunes(10)
	// default
	valueOpts.Values = []string{
		fmt.Sprintf("existingObjstoreSecret=%s", thanos.ObjStoreConfig),
		fmt.Sprintf("storegateway.enabled=%s", "true")}

	// query config
	extraFlags := []string{}
	if thanos.Querier.AutoDownsample {
		extraFlags = append(extraFlags, "--query.auto-downsampling")
	}
	if thanos.Querier.DedupEnbaled {
		extraFlags = append(extraFlags, "--query.replica-label="+"app")
	}
	if thanos.Querier.PartialResponse {
		extraFlags = append(extraFlags, "--query.partial-response")
	}
	extraFlagsResult := "{" + strings.Join(extraFlags, ",") + "}"
	valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("query.extraFlags=%s", extraFlagsResult))
	// receive/sidecar config
	if thanos.Receiver.Name != "" {
		valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("receive.enabled=%s", "true"))
		valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("receive.service.type=%s", "LoadBalancer"))
	} else {
		targets := "{" + strings.Join(thanos.Querier.Targets, ",") + "}"
		valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("query.stores=%s", targets))
	}

	// compactor config
	if thanos.Compactor.Name != "" {
		valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("compactor.enabled=%s", "true"))
		extraFlagsCompactor := []string{}

		if !thanos.Compactor.Downsampling {
			extraFlagsCompactor = append(extraFlagsCompactor, "--downsampling.disable")
		}
		// prometheus instance replica labels
		if thanos.Compactor.Deduplication {
			extraFlagsCompactor = append(extraFlagsCompactor, "--deduplication.replica-label="+"app")
		}
		if thanos.Compactor.RetentionResolutionRaw != "" {
			valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("compactor.retentionResolutionRaw=%s", thanos.Compactor.RetentionResolutionRaw))
		}
		if thanos.Compactor.RetentionResolution5m != "" {
			valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("compactor.retentionResolution5m=%s", thanos.Compactor.RetentionResolution5m))
		}
		if thanos.Compactor.RetentionResolution1h != "" {
			valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("compactor.retentionResolution1h=%s", thanos.Compactor.RetentionResolution1h))
		}
		result := "{" + strings.Join(extraFlagsCompactor, ",") + "}"
		valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("compactor.extraFlags=%s", result))
	}
	// query frontend config
	if thanos.Querierfe.Name != "" {
		valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("queryFrontend.enabled=%s", "true"))
		if thanos.Querierfe.Cacheoption == "in-memory" {
			maxSizeMap := thanos.Querierfe.Config["maxSize"]
			maxItemSizeMap := thanos.Querierfe.Config["maxItemSize"]
			var maxSize, maxItemSize string
			// maxSize
			switch maxSizeMap := maxSizeMap.(type) {
			case int:
				maxSize = strconv.Itoa(maxSizeMap)
			case string:
				maxSize = maxSizeMap
			}
			// maxItemSize
			switch maxItemSizeMap := maxItemSizeMap.(type) {
			case int:
				maxItemSize = strconv.Itoa(maxItemSizeMap)
			default:
				return nil, errors.New("invalid maxItemSize type")
			}
			inMemConf := "--query-range.response-cache-config=" + `"config"` + ":\n  " +
				`"max_size": ` + maxSize + "\n  " +
				`"max_item_size": ` + maxItemSize + "\n" +
				`"type": "in-memory"`
			inMemConfResult := "{" + inMemConf + "}"
			valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("queryFrontend.extraFlags=%s", inMemConfResult))
		} else if thanos.Querierfe.Cacheoption == "memcached" {
			var ok bool
			var addressMap interface{}
			var address string
			if addressMap, ok = thanos.Querierfe.Config["address"]; !ok {
				return nil, errors.New("memcached address doesn't exist")
			}
			switch addressMap := addressMap.(type) {
			case string:
				address = addressMap
			default:
				return nil, errors.New("invalid memcached address type")
			}
			memCacheConf := "--query-range.response-cache-config=" + `"config"` + ":\n  " +
				`"addresses":` + "\n  " +
				`  - ` + `"` + "dnssrv+_grpc._tcp." + address + `"` + "\n  " +
				`"` + "dns_provider_update_interval" + `": ` + `"` + "10s" + `"` + "\n  " +
				`"` + "max_async_buffer_size" + `": ` + "10000" + "\n  " +
				`"` + "max_async_concurrency" + `": ` + "20" + "\n  " +
				`"` + "max_get_multi_batch_size" + `": ` + "0" + "\n  " +
				`"` + "max_get_multi_concurrency" + `": ` + "100" + "\n  " +
				`"` + "max_idle_connections" + `": ` + "100" + "\n  " +
				`"` + "timeout" + `": ` + `"` + "500ms" + `"` + "\n" +
				`"` + "type" + `": ` + `"` + "memcached" + `"`
			memCachedConfResult := "{" + memCacheConf + "}"
			valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("queryFrontend.extraFlags=%s", memCachedConfResult))
		}
	}

	// ruler config
	if thanos.Ruler.Name != "" {
		alertmanagers := "{" + strings.Join(thanos.Ruler.Alertmanagers, ",") + "}"
		valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("ruler.enabled=%s", "true"))
		valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("ruler.alertmanagers=%s", alertmanagers))
		valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("ruler.config=%s", thanos.Ruler.Config))
	}
	// pod annotations for adding secret
	valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("receive.podAnnotations.secret=%s", "krius-secret-"+hash))
	valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("compactor.podAnnotations.secret=%s", "krius-secret-"+hash))
	valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("ruler.podAnnotations.secret=%s", "krius-secret-"+hash))
	valueOpts.Values = append(valueOpts.Values, fmt.Sprintf("storegateway.podAnnotations.secret=%s", "krius-secret-"+hash))
	return &valueOpts, nil
}

func createSecretforObjStore(configType string, bucConfig map[string]interface{}) (map[string][]byte, error) {
	//create a secret for bucket config
	var obj Objspec
	obj.ConfigType = configType
	obj.Config = bucConfig
	objYaml, err := yaml.Marshal(obj)
	if err != nil {
		return nil, err
	}
	secretSpec := map[string][]byte{}
	secretSpec["objstore.yml"] = objYaml
	return secretSpec, nil
}
