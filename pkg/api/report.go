package api

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Change represents a single change between 2 APIMaps for a given resource.
type Change struct {
	Type         ChangeType
	Group        string
	Resource     string
	VersionsFrom []string
	VersionsTo   []string
}

// ChangeType is an enummeration of all the possible kinds of changes between 2 APIs.
type ChangeType string

const (
	Addition     ChangeType = "added"
	Removal      ChangeType = "removed"
	Modification ChangeType = "modified"
)

// Report takes 2 mapped k8s APIs and builds a list of changes for each group, resource and version.
// The report can be represented in various formats for either human or machine consumption.
// How to present a change:
//
//		 group/res       From          change             To
//	  --------------------------------------------------------------------
//	  - group/
//	      resource:    NONE          added              versionA
//	      resource:    versionA      modified           versionB (or VersionA,VersionB)
//	      resource:    versionA      removed            NONE
func Report(from, to *Map) ([]Change, error) {
	result := make([]Change, 0)
	less := func(a, b string) bool { return a < b }

	for fromGroup, fromResources := range from.Groups {
		for fromRes, fromVersions := range fromResources {
			if _, exists := to.Groups[fromGroup][fromRes]; exists {
				// resource in both, add change if versions differ
				if !cmp.Equal(fromVersions, to.Groups[fromGroup][fromRes], cmpopts.SortSlices(less)) {
					result = append(result, Change{
						Type:         Modification,
						Group:        fromGroup,
						Resource:     fromRes,
						VersionsFrom: fromVersions,
						VersionsTo:   to.Groups[fromGroup][fromRes],
					})
				}
			} else {
				// resource in from but not in to => it was removed
				result = append(result, Change{
					Type:         Removal,
					Group:        fromGroup,
					Resource:     fromRes,
					VersionsFrom: fromVersions,
				})
			}
		}
	}

	// iterate "to" map to find newly added resources only
	for toGroup, toResources := range to.Groups {
		for toRes, toVersions := range toResources {
			if _, exists := from.Groups[toGroup][toRes]; !exists {
				result = append(result, Change{
					Type:       Addition,
					Group:      toGroup,
					Resource:   toRes,
					VersionsTo: toVersions,
				})
			}
		}
	}

	return result, nil
}
