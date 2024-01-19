package json

import (
	"encoding/json"

	"github.com/clebs/kubeglass/pkg/api"
)

// ToApiMap converts a JSON aggregated discovery into an APIMap
func ToAPIMap(data []byte) (*api.Map, error) {
	d := discovery{}
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}

	result := &api.Map{
		Groups: make(map[string]api.Resources, len(d.Items)),
	}

	for _, group := range d.Items {
		result.Groups[group.metadata.Name] = make(api.Resources)
		groupMap := result.Groups[group.metadata.Name]

		for _, version := range group.Versions {
			for _, resource := range version.Resources {
				groupMap[resource.Resource] = append(groupMap[resource.Resource], version.Version)
			}
		}
	}

	return result, nil
}

// discovery is a simplified representation of the apidiscovery.k8s.io JSON resource containing only the relevant information for the purpose at hand.
// This allows us to see all API changes per group/resource.
type discovery struct {
	Items []item `json:"items"`
}

type item struct {
	metadata `json:"metadata"`
	Versions []version `json:"versions"`
}

type metadata struct {
	Name string `json:"name"`
}

type version struct {
	Resources []resource `json:"resources"`
	Version   string     `json:"version"`
}

type resource struct {
	Resource string `json:"resource"`
}
