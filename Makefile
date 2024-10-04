GIT_VERSION ?= $(shell git describe --tags --always --dirty)
GIT_HASH ?= $(shell git rev-parse HEAD)

PKG=sigs.k8s.io/release-utils/version
LDFLAGS=-X $(PKG).gitVersion=$(GIT_VERSION)

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)
GOLANG_CI_LINT = $(LOCALBIN)/golangci-lint
GOLANG_CI_LINT_VERSION := $(shell cat .github/workflows/lint.yml | grep [[:space:]]version: | sed 's/.*version: //')

.PHONY: build
build:
	CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o ./output/kubeglass .

.PHONY: test
test:
	go test ./... -coverprofile coverage.out -race -covermode=atomic

.PHONY: lint
lint:
	GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANG_CI_LINT_VERSION)
	$(LOCALBIN)/golangci-lint run

.PHONY: release
release:
	LDFLAGS="$(LDFLAGS)" goreleaser release --clean

# used when need to validate the goreleaser
.PHONY: snapshot
snapshot:
	LDFLAGS="$(LDFLAGS)" goreleaser release --skip=sign --skip=publish --snapshot --clean

.PHONY: clean
clean:
	rm -rf output/kubeglass
