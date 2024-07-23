PROJECT := ssclient
.DEFAULT_GOAL := ssclient
UNAME_OS := $(shell uname -s)
UNAME_ARCH := $(shell uname -m)
CACHE_BASE ?= $(HOME)/.cache/$(PROJECT)
CACHE := $(CACHE_BASE)/$(UNAME_OS)/$(UNAME_ARCH)
CACHE_BIN := $(CACHE)/bin
CACHE_VERSIONS := $(CACHE)/versions
CACHE_GOBIN := $(CACHE)/gobin
CACHE_GOCACHE := $(CACHE)/gocache
export PATH := $(abspath $(CACHE_BIN)):$(PATH)
SHELL := /usr/bin/env bash -o pipefail
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-print-directory

MOCKGEN_VERSION ?= 0.4.0
MOCKGEN := $(CACHE_VERSIONS)/mockgen/$(MOCKGEN_VERSION)
$(MOCKGEN):
	@rm -f $(CACHE_BIN)/mockgen
	@mkdir -p $(CACHE_BIN)
	@env GOBIN=$(CACHE_BIN) go install go.uber.org/mock/mockgen@v$(MOCKGEN_VERSION)
	@chmod +x $(CACHE_BIN)/mockgen
	@rm -rf $(dir $(MOCKGEN))
	@mkdir -p $(dir $(MOCKGEN))
	@touch $(MOCKGEN)

.PHONY: examplemocks
examplemocks: $(MOCKGEN)
	cd $(CURDIR)/example && mockgen -typed -destination=./adapter/adapter.go -package=adapter github.com/microsoft/kiota-abstractions-go RequestAdapter

.PHONY: ssclient
ssclient:
	@which kiota > /dev/null 2>&1 || (echo "kiota not found in PATH, download v1.14.0 from: https://learn.microsoft.com/en-ca/openapi/kiota/install" && exit 1)
	kiota generate --language go --clean-output --class-name Client --namespace-name xgo.artefactual.dev/ssclient/kiota --openapi typespec/tsp-output/@typespec/openapi3/openapi.v1.yaml --output ./kiota
	$(MAKE) update-kiota-imports

.PHONY: update-kiota-imports
update-kiota-imports:
	# We use `artefactual.dev/ssclient/kiota` as the namespace because kiota
	# excapes `go.artefactual.dev/ssclient/kiota`. See https://github.com/microsoft/kiota/issues/5012 for more details.
	find ./kiota -type f -name '*.go' -exec sed -i 's|xgo.artefactual.dev/ssclient/kiota/|go.artefactual.dev/ssclient/kiota/|g' {} +

.PHONY: typespec
typespec:
	npm --prefix=$(CURDIR)/typespec run compile

.PHONY: gen
gen: typespec ssclient examplemocks
