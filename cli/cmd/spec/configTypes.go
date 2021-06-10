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

// type Data2 map[string]map[string]interface{}
type Data map[string]interface{}

// Objstoresonfiglist
type ObjStoreConfigslist struct {
	Bucketweb Bucketweb `yaml:"bucketweb"`
	Name      string    `yaml:"name"`
	Type      string    `yaml:"type"`
}

// Prometheus
type Prometheus struct {
	Install          bool   `yaml:"install"`
	Name             string `yaml:"name"`
	Namespace        string `yaml:"namespace"`
	Mode             string `yaml:"mode"`
	ReceiveReference string `yaml:"receiveReference"`
	Objstoresonfig   string `yaml:"objStoreConfig"`
}

// Thanos
type Thanos struct {
	Name      string    `yaml:"name"`
	Querier   Querier   `yaml:"querier"`
	Querierfe Querierfe `yaml:"querierFE"`
	Reciever  Reciever  `yaml:"receiver"`
	Compactor Compactor `yaml:"compactor"`
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
	Targets         string `yaml:"targets"`
	Dedupenbaled    string `yaml:"dedupEnbaled"`
	Autoownample    string `yaml:"autoDownSample"`
	Partialresponse string `yaml:"partialResponse"`
	Name            string `yaml:"name"`
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

// Reciever
type Reciever struct {
	Name string `yaml:"name"`
}

// Bucketweb
type Bucketweb struct {
	Enabled bool `yaml:"enabled"`
}