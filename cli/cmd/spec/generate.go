package spec

import (
	"log"
	"os"

	client "github.com/infracloudio/krius/pkg/client"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a profile based on questions asked to user",
	Run:   createConfigYAML,
}

const defaultObjectStorageConfigName = "krius-bucketcluster"

func init() {
	specCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringP("file", "f", "", "file Path to genrate the config file")
	generateCmd.Flags().StringP("mode", "m", "", "Mode --> receiver/sidecar")
	err := generateCmd.MarkFlagRequired("mode")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func createFile(filePath string, content string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	_, err = file.WriteString("---\n" + content)
	if err != nil {
		log.Fatalf("Failed to write yaml file %s", err)
	}
	defer file.Close()
}

func createConfigYAML(cmd *cobra.Command, args []string) {
	mode, err := cmd.Flags().GetString("mode")
	if err != nil {
		log.Fatalf("error: %v", err)
	} else if mode != "receiver" && mode != "sidecar" {
		log.Fatalf("error: invalid mode: %s", mode)
	}

	bucketweb := client.Bucketweb{Enabled: true}

	receiver := client.Receiver{
		Name: "receiver",
	}

	ruler := client.Ruler{}
	ruler.Name = "ruler"
	ruler.Alertmanagers = []string{"http://kube-prometheus-alertmanager.monitoring.svc.cluster.local:9093"}
	rulerStr := map[string]interface{}{"group": map[string]interface{}{"name": "metamonitoring", "rules": map[string]interface{}{"alert": "PrometheusDown", "expr": "absent(up{prometheus='monitoring/kube-prometheus'})"}}}
	rulerStrYaml, err := yaml.Marshal(&rulerStr)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	ruler.Config = string(rulerStrYaml)

	compactor := client.Compactor{}
	compactor.Name = "compactor"
	compactor.Deduplication = true
	compactor.Deduplication = true
	compactor.RetentionResolution1h = "10y"
	compactor.RetentionResolution5m = "30d"
	compactor.RetentionResolutionRaw = "30d"

	querier := client.Querier{
		DedupEnbaled:    true,
		AutoDownsample:  true,
		PartialResponse: true,
		Name:            "querier",
	}

	querierfe := client.Querierfe{}
	querierfe.Name = "querierfe"
	querierfe.Cacheoption = "inMemory"
	querierfe.Config = map[string]interface{}{"maxSixe": 1}

	cluster1 := client.Cluster{}
	cluster1.Name = "prometheus"
	cluster1.Type = "prometheus"
	cluster1.Data = map[string]interface{}{"install": true, "name": "prometheus", "namespace": "default", "mode": mode, "objStoreConfig": defaultObjectStorageConfigName}

	cluster2 := client.Cluster{}
	cluster2.Name = "thanos"
	cluster2.Type = "thanos"
	cluster2.Data = map[string]interface{}{"install": true, "name": "thanos", "namespace": "default", "querier": querier, "querierFE": querierfe, "compactor": compactor, "ruler": ruler, "objStoreConfig": defaultObjectStorageConfigName}

	if mode == "receiver" {
		cluster2.Data["receiver"] = receiver
	}

	buckerconfig := make(map[string]interface{})
	buckerconfig["bucket"] = "Your s3 bucket name"
	buckerconfig["access_key"] = "Your AWS access key"
	buckerconfig["secret_key"] = "Your AWS secret key"
	buckerconfig["endpoint"] = "Your S3 bucket endpoint"
	buckerconfig["insecure"] = false
	buckerconfig["trace"] = true

	objstore := client.ObjStoreConfig{}
	objstore.Name = defaultObjectStorageConfigName
	objstore.Type = "S3"
	objstore.Bucketweb = bucketweb
	objstore.Config = buckerconfig

	config := client.Config{}
	config.Clusters = append(config.Clusters, cluster1)
	config.Clusters = append(config.Clusters, cluster2)
	config.ObjStoreConfigslist = append(config.ObjStoreConfigslist, objstore)

	configYAML, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	file, err := cmd.Flags().GetString("file")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if file == "" {
		createFile("config.yaml", string(configYAML))
	} else {
		createFile(file, string(configYAML))
	}
}
