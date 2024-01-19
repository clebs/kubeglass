package format

import (
	yamlEncoding "gopkg.in/yaml.v3"

	"github.com/clebs/kubeglass/pkg/api"
)

type yaml struct{}

func (yaml) Format(changes []api.Change) ([]byte, error) {
	return yamlEncoding.Marshal(changes)
}
