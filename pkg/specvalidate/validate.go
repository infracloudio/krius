package specvalidate

import (
	"fmt"
	"io/ioutil"

	"github.com/xeipuuv/gojsonschema"
	"sigs.k8s.io/yaml"
)

func GetLoaders(configPath string) (gojsonschema.JSONLoader, gojsonschema.JSONLoader, error) {
	ruleSchemaLoader := gojsonschema.NewBytesLoader([]byte(RuleSchema))
	loader, err := jsonLoader(configPath)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: unable to load schema ref: %s", configPath, err)
	}
	return loader, ruleSchemaLoader, nil
}

func ValidateYML(loader gojsonschema.JSONLoader, schemaLoader gojsonschema.JSONLoader) (bool, []gojsonschema.ResultError) {
	result, err := gojsonschema.Validate(schemaLoader, loader)
	if err != nil {
		panic(err.Error())
	}
	return result.Valid(), result.Errors()
}

func jsonLoader(path string) (gojsonschema.JSONLoader, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	buf, err = yaml.YAMLToJSON(buf)
	if err != nil {
		return nil, err
	}
	return gojsonschema.NewBytesLoader(buf), nil
}
