package client

// client to Preflight checks and installing tools on cluster
type Client interface {
	PreflightChecks(c *Config, clusterName string) ([]string, error)
	InstallClient(clusterName string, targets []string) (string, error)
	UninstallClient(clusterName string) error
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
	Trace      Trace  `yaml:"trace,omitempty"`
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
	Name        string                 `yaml:"name"`
	Cacheoption string                 `yaml:"cacheOption"`
	Config      map[string]interface{} `yaml:"config"`
}

// Querier
type Querier struct {
	Targets         []string `yaml:"targets,omitempty"`
	DedupEnbaled    bool     `yaml:"dedupEnbaled"`
	AutoDownsample  bool     `yaml:"autoDownSample"`
	PartialResponse bool     `yaml:"partialResponse"`
	Name            string   `yaml:"name"`
}

// Compactor
type Compactor struct {
	Name                   string `yaml:"name"`
	Downsampling           bool   `yaml:"downsampling"`
	Deduplication          bool   `yaml:"deduplication"`
	RetentionResolutionRaw string `yaml:"retentionResolutionRaw"`
	RetentionResolution5m  string `yaml:"retentionResolution5m"`
	RetentionResolution1h  string `yaml:"retentionResolution1h"`
}

// Ruler
type Ruler struct {
	Name          string   `yaml:"name"`
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
