package rules

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type SimpleYamlRule struct {
	Key   string
	Value string
	Path  string
}

func (yr SimpleYamlRule) Name() string {
	return fmt.Sprintf(
		"yaml key:%s value:%s path%s", yr.Key, yr.Value, yr.Path)
}

func (yr SimpleYamlRule) Matches() (bool, error) {
	data, err := ioutil.ReadFile(yr.Path)
	if err != nil {
		return false, err
	}

	// assuming the key value pair will be present at the root
	var yamlData map[string]interface{}
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return false, err
	}

	if val, ok := yamlData[yr.Key]; !ok {
		return false, nil
	} else {
		cval, ok := val.(string)
		if !ok {
			return false, nil
		} else {
			return yr.Value == cval, nil
		}
	}
}
