package client

// client to Preflight checks and installing tools on cluster
type Client interface {
	PreflightChecks(c *Config, clusterName string) ([]string, error)
	InstallClient(clusterName string, targets []string) (string, error)
}

// KriusConfig
type Config struct {
	Clusters            []Cluster        `yaml:"clusters"`
	ObjStoreConfigslist []ObjStoreConfig `yaml:"objStoreConfigslist"`
	Order               int              //if 1 then mode is sidecar else mode is receiver
}

// Cluster
type Cluster struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Data Data   `yaml:"data"`
}

type Data map[string]interface{}

// Objstoresonfiglist
type ObjStoreConfig struct {
	Bucketweb Bucketweb    `yaml:"bucketweb"`
	Name      string       `yaml:"name"`
	Type      string       `yaml:"type"`
	Config    BucketConfig `yaml:"config"`
}

// Bucket Config
type BucketConfig struct {
	BucketName string `yaml:"bucket"`
	Endpoint   string `yaml:"endpoint"`
	AccessKey  string `yaml:"accessKey"`
	SecretKey  string `yaml:"secretKey"`
	Insecure   bool   `yaml:"insecure"`
	Trace      Trace  `yaml:"trace"`
}

type Trace struct {
	Enable bool `yaml:"enable"`
}

// Prometheus
type Prometheus struct {
	Install          bool   `yaml:"install"`
	Name             string `yaml:"name"`
	Namespace        string `yaml:"namespace"`
	Mode             string `yaml:"mode"`
	ReceiveReference string `yaml:"receiveReference"`
	ObjStoreConfig   string `yaml:"objStoreConfig"`
}

// Thanos
type Thanos struct {
	Name           string    `yaml:"name"`
	Namespace      string    `yaml:"namespace"`
	ObjStoreConfig string    `yaml:"objStoreConfig"`
	Querier        Querier   `yaml:"querier"`
	Querierfe      Querierfe `yaml:"querierFE"`
	Receiver       Receiver  `yaml:"receiver"`
	Compactor      Compactor `yaml:"compactor"`
	Ruler          Ruler     `yaml:"ruler"`
}

// Grafana
type Grafana struct {
	Name  string `yaml:"name"`
	Setup Setup  `yaml:"setup"`
}

// Querierfe
type Querierfe struct {
	Name        string `yaml:"name"`
	Cacheoption string `yaml:"cacheOption"`
}

// Querier
type Querier struct {
	Targets         []string `yaml:"targets"`
	DedupEnbaled    string   `yaml:"dedupEnbaled"`
	AutoDownsample  bool     `yaml:"autoDownSample"`
	PartialResponse bool     `yaml:"partialResponse"`
	Name            string   `yaml:"name"`
	ExtraFlags      []string
}

// Compactor
type Compactor struct {
	Name string `yaml:"name"`
}

// Ruler
type Ruler struct {
	Alertmanagers []string `yaml:"alertManagers"`
	Config        string   `yaml:"config"`
}

// Setup
type Setup struct {
	Enabled   bool   `yaml:"enabled"`
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

// Receiver
type Receiver struct {
	Name string `yaml:"name"`
}

// Bucketweb
type Bucketweb struct {
	Enabled bool `yaml:"enabled"`
}
