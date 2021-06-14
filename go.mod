module github.com/infracloudio/krius

go 1.15

require (
	github.com/gofrs/flock v0.8.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/xeipuuv/gojsonschema v1.2.0
	gopkg.in/yaml.v2 v2.4.0
	google.golang.org/protobuf v1.25.0
	helm.sh/helm/v3 v3.5.1
	k8s.io/client-go v0.20.1 // indirect
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/yaml v1.2.0
)

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309
