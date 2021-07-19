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
}

func createConfigYAML(cmd *cobra.Command, args []string) {
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
	cluster1.Name = "cluster1"
	cluster1.Type = "prometheus"
	cluster1.Data = map[string]interface{}{"install": false, "name": "prometheus1", "namespace": "monitoring", "mode": "sidecar", "objStoreConfig": "bucketcluster"}

	cluster2 := Cluster{}
	cluster2.Name = "cluster2"
	cluster2.Type = "prometheus"
	cluster2.Data = map[string]interface{}{"install": false, "name": "prometheus2", "namespace": "monitoring", "mode": "sidecar", "objStoreConfig": "bucketcluster"}

	cluster3 := Cluster{}
	cluster3.Name = "cluster3"
	cluster3.Type = "thanos"
	cluster3.Data = map[string]interface{}{"name": "thanos-ag1", "querier": querier, "querier-fe": querierfe, "reciever": reciver, "compactor": compactor, "ruler": ruler}

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
	objstore.Bucketweb = bucketweb
	objstore.Config = buckerconfig

	config := Config{}
	config.Clusters = append(config.Clusters, cluster1)
	config.Clusters = append(config.Clusters, cluster2)
	config.Clusters = append(config.Clusters, cluster3)
	config.ObjStoreConfigslist = append(config.ObjStoreConfigslist, objstore)

	configYAML, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	configFile, err := os.Create("config.yaml")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	_, err = configFile.WriteString("---\n" + string(configYAML))

	if err != nil {
		log.Fatalf("Failed to write yaml file %s", err)
	}
	defer configFile.Close()
}
