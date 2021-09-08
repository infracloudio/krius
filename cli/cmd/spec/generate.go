package spec

import (
	"log"
	"os"
	"github.com/infracloudio/krius/pkg/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a profile based on questions asked to user",
	Run:   createConfigYAML,
}

func init() {
	specCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringP("file", "f", "", "file Path to genrate the config file")
	generateCmd.Flags().StringP("mode", "m", "", "Mode --> receiver/sidecar")
	generateCmd.MarkFlagRequired("mode")
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

	bucketweb := Bucketweb{Enabled: true}

	receiver := Receiver{
		Name: "receiver",
	}

	ruler := Ruler{}
	ruler.Alertmanagers = []string{"http://kube-prometheus-alertmanager.monitoring.svc.cluster.local:9093"}
	ruler_str := map[string]interface{}{"group": map[string]interface{}{"name": "metamonitoring", "rules": map[string]interface{}{"alert": "PrometheusDown", "expr": "absent(up{prometheus='monitoring/kube-prometheus'})"}}}
	ruler_str_yaml, err := yaml.Marshal(&ruler_str)
	ruler.Config = string(ruler_str_yaml)

	compactor := compactor{}
	compactor.Name = "compactor"
	compactor.Deduplication = true
	compactor.Deduplication = true
	compactor.RetentionResolution1h = "10y"
	compactor.RetentionResolution5m = "30d"
	compactor.RetentionResolutionRaw = "30d"

	querier := Querier{
		Dedupenbaled:    true,
		Autoownample:    true,
		Partialresponse: true,
		Name:            "querier",
	}

	querierfe := Querierfe{}
	querierfe.Name = "querierfe"
	querierfe.Cacheoption = "inMemory"

	cluster1 := Cluster{}
	cluster1.Name = "Prometheus"
	cluster1.Type = "prometheus"
	cluster1.Data = map[string]interface{}{"install": false, "name": "Prometheus", "namespace": "default", "objStoreConfig": "bucketcluster"}

	cluster2 := Cluster{}
	cluster2.Name = "Thanos"
	cluster2.Type = "thanos"
	cluster2.Data = map[string]interface{}{"name": "Thanos", "querier": querier, "querierFE": querierfe, "compactor": compactor, "ruler": ruler}

	if mode == "receiver" {
		cluster2.Data["receiver"] = receiver
	}

	buckerconfig := BucketConfig{}
	buckerconfig.BucketName = "Your s3 bucket name"
	buckerconfig.AccessKey = "Your AWS access key"
	buckerconfig.SecretKey = "Your AWS secret key"
	buckerconfig.Endpoint = "Your S3 bucket endpoint"
	buckerconfig.Insecure = false
	buckerconfig.Trace.Enable = true

	objstore := ObjStoreConfigslist{}
	objstore.Name = "bucketcluster"
	objstore.Type = "s3"
	objstore.Bucketweb = bucketweb
	objstore.Config = buckerconfig

	config := Config{}
	config.Clusters = append(config.Clusters, cluster1)
	config.Clusters = append(config.Clusters, cluster2)
	config.ObjStoreConfigslist = append(config.ObjStoreConfigslist, objstore)

	configYAML, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	file, err := cmd.Flags().GetString("file")
	if file == "" {
		createFile("config.yaml", string(configYAML))
	} else {
		createFile(file, string(configYAML))
	}
}
