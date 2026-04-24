# Makefile for podman
# See docs/make.md for usage

export GOPROXY ?= https://proxy.golang.org

GO ?= go
GOFLAGS ?= -trimpath
GOTAGS ?= \
	exclude_graphdriver_devicemapper \
	exclude_graphdriver_btrfs \
	btrfs_noversion \
	containers_image_openpgp

GOLD_FLAGS ?=
GO_BUILD_FLAGS ?= $(GOFLAGS) -tags '$(GOTAGSARGS)'
GOTAGSARGS ?= $(GOTAGSS)

PROJECT := github.com/containers/podman
BINDIR ?= $(DESTDIR)/usr/local/bin
LIBEXECDIR ?= $(DESTDIR)/usr/local/libexec
MANDIR ?= $(DESTDIR)/usr/local/share/man
SHAREDIR ?= $(DESTDIR)/usr/local/share
ETCDIR ?= $(DESTDIR)/etc
TMPFILESDIR ?= $(DESTDIR)/usr/lib/tmpfiles.d
SYSTEMDDIR ?= $(DESTDIR)/usr/lib/systemd/system
USERSYSTEMDDIR ?= $(DESTDIR)/usr/lib/systemd/user

BINARIES ?= bin/podman

# Version info
VERSION ?= $(shell cat version/rawversion)
GIT_COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_INFO ?= $(shell date +%s)

LDFLAGS_PODMAN ?= \
	-X $(PROJECT)/libpod/define.gitCommit=$(GIT_COMMIT) \
	-X $(PROJECT)/vendor/github.com/containers/common/pkg/config.additionalHelperBinariesDir=$(HELPER_BINARIES_DIR)

all: binaries

.PHONY: binaries
binaries: $(BINARIES)

bin/podman: $(shell find . -name '*.go' -not -path './vendor/*') go.mod go.sum
	@mkdir -p bin
	$(GO) build \
		$(GO_BUILD_FLAGS) \
		-ldflags '$(LDFLAGS_PODMAN) $(GOLD_FLAGS)' \
		-o $@ ./cmd/podman

.PHONY: clean
clean:
	rm -rf \
		$(BINARIES) \
		build \
		bin

.PHONY: install
install:
	install -d -m 755 $(BINDIR)
	install -m 755 bin/podman $(BINDIR)/podman

.PHONY: test
test: unit integration

.PHONY: unit
unit:
	$(GO) test -v ./...

.PHONY: integration
integration:
	$(GO) test -v -tags integration ./test/...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: vendor
vendor:
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: help
help:
	@echo 'Usage: make <target>'
	@echo ''
	@echo 'Targets:'
	@echo '  all        - Build all binaries (default)'
	@echo '  binaries   - Build podman binary'
	@echo '  clean      - Remove build artifacts'
	@echo '  install    - Install binaries'
	@echo '  test       - Run all tests'
	@echo '  unit       - Run unit tests'
	@echo '  integration - Run integration tests'
	@echo '  lint       - Run linter'
	@echo '  vendor     - Update vendor directory'
	@echo '  fmt        - Format Go source files'
	@echo '  vet        - Run go vet'
