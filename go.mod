module github.com/infracloudio/krius

go 1.16

require (
	github.com/briandowns/spinner v1.16.0 // indirect
	github.com/gofrs/flock v0.8.0
	github.com/manifoldco/promptui v0.8.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/xeipuuv/gojsonschema v1.2.0
	gopkg.in/yaml.v2 v2.4.0
	helm.sh/helm/v3 v3.6.1
	k8s.io/api v0.21.0
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v0.21.0
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/yaml v1.2.0
)

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309
