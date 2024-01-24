package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	// Github API URL to list files in discovery to find out the name of the aggregated json. Add release version.
	listDiscoveryFilesURL = `https://api.github.com/repos/kubernetes/kubernetes/contents/api/discovery?ref=release-%s`
)

// SemVer fetches the aggregated discovery spec for a given kubernetes version
func SemVer(v string) ([]byte, error) {
	v = strings.TrimPrefix(v, "v")
	// TODO: consider caching this to avoid exhausting rates.
	lsResp, err := http.DefaultClient.Get(fmt.Sprintf(listDiscoveryFilesURL, v))
	if err != nil {
		return nil, fmt.Errorf("could not get the kubernetes discovery file from github: %w", err)
	}
	defer lsResp.Body.Close()

	out := files{}
	if err := json.NewDecoder(lsResp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("could not decode the github kubernetes discovery file list: %w", err)
	}

	download := ""
	for _, f := range out {
		if strings.Contains(f.Name, "aggregated") {
			download = f.DownloadURL
		}
	}

	if download == "" {
		return nil, fmt.Errorf("could not find aggregated discovery api file")
	}

	dlResp, err := http.DefaultClient.Get(download)
	if err != nil {
		return nil, fmt.Errorf("could not get the aggregated discovery file [%s]: %w", download, err)
	}
	defer dlResp.Body.Close()

	return io.ReadAll(dlResp.Body)
}

type files []struct {
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
}
