package rules

import (
	"fmt"
	"os"
)

type EnvVarRule struct {
	Key   string
	Value string
}

func (evr EnvVarRule) Name() string {
	return fmt.Sprintf(
		"envvar key:%s and value:%s",
		evr.Key,
		evr.Value)
}

func (evr EnvVarRule) Matches() (bool, error) {
	if val, ok := os.LookupEnv(evr.Key); ok && val == evr.Value {
		return true, nil
	}
	return false, nil
}
