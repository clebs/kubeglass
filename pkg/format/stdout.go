package format

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/clebs/kubeglass/pkg/api"
)

type stdout struct{}

func (s stdout) Format(changes []api.Change) ([]byte, error) {
	out := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)

	// print headers
	fmt.Fprintf(out, "RESOURCE\tFROM\tCHANGE\tTO\n")

	for _, c := range changes {
		fromVersions := "NONE"
		if len(c.VersionsFrom) > 0 {
			fromVersions = strings.Join(c.VersionsFrom, ",")
		}

		toVersions := "NONE"
		if len(c.VersionsTo) > 0 {
			toVersions = strings.Join(c.VersionsTo, ",")
		}

		fmt.Fprintf(out, "%s/%s\t%s\t%s\t%s\n", c.Group, c.Resource, fromVersions, c.Type, toVersions)
	}
	out.Flush()
	return nil, nil
}
