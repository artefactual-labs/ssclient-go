SHELL := /bin/bash
.DEFAULT_GOAL := help
.PHONY: *
SPEC_FILE := $(CURDIR)/spec/openapi.v1.yaml

DBG_MAKEFILE ?=
ifeq ($(DBG_MAKEFILE),1)
    $(warning ***** starting Makefile for goal(s) "$(MAKECMDGOALS)")
    $(warning ***** $(shell date))
else
    # If we're not debugging the Makefile, don't echo recipes.
    MAKEFLAGS += -s
endif

deps: # @HELP List available direct dependency updates.
	mise exec -- go list -u -m -json all | go-mod-outdated -update -direct

examplemocks: # @HELP Generate the example app mocks.
examplemocks: $(MOCKGEN)
	cd $(CURDIR)/example && mise exec -- mockgen -typed -destination=./adapter/adapter.go -package=adapter github.com/microsoft/kiota-abstractions-go RequestAdapter

lint: # @HELP Lint the project Go files with golangci-lint (linters + formatters).
lint: LINT_FLAGS ?= --fix=1
lint:
	mise exec -- golangci-lint run $(LINT_FLAGS)

spec-check: # @HELP Ensure the spec submodule is initialized.
	@test -f "$(SPEC_FILE)" || { \
		printf '%s\n' \
			'OpenAPI spec not found at $(SPEC_FILE).' \
			'Run: git submodule update --init --recursive'; \
		exit 1; \
	}

ssclient: # @HELP Generate the Kiota client from the OpenAPI spec.
ssclient: spec-check
	KIOTA_TUTORIAL_ENABLED=false mise exec -- kiota generate --language go --clean-output --class-name Client --namespace-name go.artefactual.dev/ssclient/kiota --openapi $(SPEC_FILE) --output ./kiota --exclude-backward-compatible
	@printf '%s\n' \
		'-------------------------------------------------------------------------' \
		'# WARNING: kiota emits some expected warnings:                          #' \
		'#                                                                       #' \
		'# - Polymorphic schema without discriminator:                           #' \
		'#   /api/v2/file/{uuid}/review_aip_deletion/ is polymorphic but does    #' \
		'#   not define a discriminator, which may cause serialization issues.   #' \
		'#                                                                       #' \
		'# - Multiple divergent success schemas:                                 #' \
		'#   /api/v2/file/{uuid}/delete_aip/ defines multiple success responses  #' \
		'#   (e.g. 200 and 202) with different schemas. Kiota uses the lowest    #' \
		'#   success status code schema.                                         #' \
		'#                                                                       #' \
		'# - Missing typed error models for download endpoints:                  #' \
		'#   Some endpoints (e.g. Packages_downloadFile and                      #' \
		'#   Packages_downloadPointerFile) return plain text or unsupported      #' \
		'#   error responses (400, 404, 501), so Kiota cannot generate typed     #' \
		'#   error models for them.                                              #' \
		'#                                                                       #' \
		'-------------------------------------------------------------------------'

ssclient-deps: # @HELP Show the Kiota generation dependency graph.
ssclient-deps: spec-check
	KIOTA_TUTORIAL_ENABLED=false mise exec -- kiota info --openapi $(SPEC_FILE) --language Go

gen: # @HELP Generate code.
gen: ssclient examplemocks

test: # @HELP Run tests.
test:
	go test -v ./...
	cd $(CURDIR)/example && go test -v ./...

help: # @HELP Print this message.
help:
	echo "TARGETS:"
	grep -E '^.*: *# *@HELP' Makefile             \
	    | awk '                                   \
	        BEGIN {FS = ": *# *@HELP"};           \
	        { printf "  %-30s %s\n", $$1, $$2 };  \
	    '
