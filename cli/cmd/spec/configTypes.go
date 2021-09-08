package spec

// KriusConfig
type Config struct {
	Clusters            []Cluster             `yaml:"clusters"`
	ObjStoreConfigslist []ObjStoreConfigslist `yaml:"objStoreConfigslist"`
}

// Cluster
type Cluster struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Data Data   `yaml:"data"`
}

type Data map[string]interface{}

// Objstoresonfiglist
type ObjStoreConfigslist struct {
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
	Name      string    `yaml:"name"`
	Querier   Querier   `yaml:"querier"`
	Querierfe Querierfe `yaml:"querierFE"`
	Receiver  Receiver  `yaml:"receiver"`
	Compactor compactor `yaml:"compactor"`
	Ruler     Ruler     `yaml:"ruler"`
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
	Targets         string `yaml:"targets,omitempty"`
	Dedupenbaled    bool   `yaml:"dedupEnbaled"`
	Autoownample    bool   `yaml:"autoDownSample"`
	Partialresponse bool   `yaml:"partialResponse"`
	Name            string `yaml:"name"`
}

// Compactor
type compactor struct {
	Name                   string `yaml:"name"`
	Downsampling           bool   `yaml:"downsampling"`
	Deduplication          bool   `yaml:"deduplication"`
	RetentionResolutionRaw string `yaml:"retentionResolutionRaw"`
	RetentionResolution5m  string `yaml:"retentionResolution5m"`
	RetentionResolution1h  string `yaml:"retentionResolution1h"`
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

// Reciever
type Receiver struct {
	Name string `yaml:"name"`
}

// Bucketweb
type Bucketweb struct {
	Enabled bool `yaml:"enabled"`
}
