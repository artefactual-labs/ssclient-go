PROJECT := ssclient
.DEFAULT_GOAL := ssclienttest
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

.PHONY: ssclienttest
ssclienttest: $(MOCKGEN)

.PHONY: ssclient
ssclient:
	rm -rf $(CURDIR)/kiota
	/tmp/kiota/kiota generate --language go --class-name Client --namespace-name go.artefactual.dev/ssclient/kiota --openapi typespec/tsp-output/@typespec/openapi3/openapi.yaml --output ./kiota
	mv $(CURDIR)/kiota/goescaped/artefactual/dev/ssclient/kiota/* kiota/
	rm -rf $(CURDIR)/kiota/goescaped/
	$(MAKE) update-kiota-imports

.PHONY: update-kiota-imports
update-kiota-imports:
	find ./kiota -type f -name '*.go' -exec sed -i 's|goescaped.artefactual.dev/ssclient/kiota/|go.artefactual.dev/ssclient/kiota/|g' {} +

