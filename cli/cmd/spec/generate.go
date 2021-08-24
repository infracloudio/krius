package spec

import (
	"log"
	"os"

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
	generateCmd.Flags().StringP("mode", "m", "", "Mode --> reciever/sidecar")
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
	} else if mode != "reciever" && mode != "sidecar" {
		log.Fatalf("error: invalid mode: %s", mode)
	}

	bucketweb := Bucketweb{Enabled: true}

	reciver := Reciever{Name: "#Name of your Reciver"}

	setup := Setup{}
	setup.Enabled = true
	setup.Name = "#Name of Your Grafana Setup"
	setup.Namespace = "#Your Grafana Namesapce"

	ruler := Ruler{}
	ruler.Alertmanagers = []string{"http://kube-prometheus-alertmanager.monitoring.svc.cluster.local:9093"}
	ruler_str := map[string]interface{}{"group": map[string]interface{}{"name": "metamonitoring", "rules": map[string]interface{}{"alert": "PrometheusDown", "expr": "absent(up{prometheus='monitoring/kube-prometheus'})"}}}
	ruler_str_yaml, err := yaml.Marshal(&ruler_str)
	ruler.Config = string(ruler_str_yaml)

	compactor := Compactor{Name: "#Name of your compactor"}

	querier := Querier{
		Targets:         "",
		Dedupenbaled:    "",
		Autoownample:    "",
		Partialresponse: "",
		Name:            "",
	}

	querierfe := Querierfe{}

	grafana := Grafana{}
	grafana.Setup = setup

	cluster1 := Cluster{}
	cluster1.Name = "#Name of your Prometheus Cluster"
	cluster1.Type = "prometheus"
	cluster1.Data = map[string]interface{}{"install": false, "name": "#Name Of Prometheus Cluster", "namespace": "#Namespace Name", "objStoreConfig": "bucketcluster", "mode": mode}

	cluster2 := Cluster{}
	cluster2.Name = "#Name of your Thanos Cluster"
	cluster2.Type = "thanos"
	cluster2.Data = map[string]interface{}{"name": "#Name Of Thanos Cluster", "querier": querier, "querierFE": querierfe, "compactor": compactor, "ruler": ruler}

	if mode == "reciever" {
		cluster2.Data["reciever"] = reciver
	}

	prom := Prometheus{}
	prom.Name = "cluster1"
	prom.Install = true
	prom.Mode = "sidecar"
	prom.Namespace = "monitoring"
	prom.ObjStoreConfig = "bucketcluster"

	thanos := Thanos{}
	thanos.Name = "Thanos"
	thanos.Compactor = compactor
	thanos.Querier = querier
	thanos.Querierfe = querierfe
	thanos.Reciever = reciver
	thanos.Ruler = ruler

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
