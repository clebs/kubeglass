# KubeGlass 🔮

## Overview
Make kubernetes API changes between versions as transparent as glass.

**NOTE:** This program only works with kubernetes versions 1.28+.

## Description
KubeGlass allows to compare 2 kubernetes versions and see what changed on the different APIs between them.

## Installation
Run `go install github.com/clebs/kubeglass`.

## Usage
### Online
To compare 2 kubernetes versions run `kubeglass --from=1.28 --to=1.29`. This will show the following output:
![kubeglass output](./assets/kubeglass-sample.png)

This usage requires internet connection as it downloads all necessary information to compare APIs based on the given versions.

### Offline
For those cases where an internet connection is not available, the aggregated discovery API files can be downloaded from the Kubernetes GitHub (e.g. [here](https://github.com/kubernetes/kubernetes/blob/release-1.29/api/discovery/aggregated_v2beta1.json)).
Then running kubeglass as follows:
`kubeglass --from=./relative/path/to/file/aggregated-api-1.28.json --to=/absolute/path/to/file/aggregated-api-1.29.json`

## Roadmap
- [ ] Persistent caching of downloaded data
- [x] Airgapped support
    - Allow `-f/--from` and `-t/--to` flags to take input files.
- [ ] Richtext and colored printing to stdout
- [ ] Krew support
- [ ] Scan running clusters via the `/apis` endpoint
