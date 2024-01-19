package format

import (
	jsonEncoding "encoding/json"

	"github.com/clebs/kubeglass/pkg/api"
)

type json struct{}

func (json) Format(changes []api.Change) ([]byte, error) {
	return jsonEncoding.MarshalIndent(changes, "", "  ")
}
