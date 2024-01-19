package api

// Map contains kubernetes API information in a set of maps for efficient access, comparison and output.
// The structure is a root map containing the groups, then each group has a map of resources.
// Finally, each resource contains a set of versions it supports
type Map struct {
	Groups map[string]Resources
}

// Resources is a map where each key is one resource and the value is the set of versions it supports.
// Example { "Pod": ["v1alpha2", "v1"] }
type Resources map[string]Versions

// Versions is a list of versions (e.g. ["v1", "v2beta1"])
type Versions []string
