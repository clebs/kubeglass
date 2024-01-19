package format

import "github.com/clebs/kubeglass/pkg/api"

type Formatter interface {
	Format(changes []api.Change) ([]byte, error)
}

func NewFormatter(t string) Formatter {
	switch t {
	case "json":
		return json{}
	case "yaml":
		return yaml{}
	default:
		// default and rich text stdout cases combined
		return stdout{}
	}
}
