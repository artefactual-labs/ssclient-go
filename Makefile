SHELL := /bin/bash
.DEFAULT_GOAL := help
.PHONY: *

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

ssclient: # @HELP Generate the Kiota client from the OpenAPI spec.
	KIOTA_TUTORIAL_ENABLED=false mise exec -- kiota generate --language go --clean-output --class-name Client --namespace-name go.artefactual.dev/ssclient/kiota --openapi typespec/tsp-output/@typespec/openapi3/openapi.v1.yaml --output ./kiota
	@printf '%s\n' \
		'-------------------------------------------------------------------------' \
		'# WARNING: kiota emits some expected warnings:                          #' \
		'#                                                                       #' \
		'# - multiple success schemas: the API really has two different success  #' \
		'#   payloads (200 and 202), and the high-level wrapper handles that     #' \
		'#   endpoint manually.                                                  #' \
		'# - plain-text error responses on download endpoints: SS really returns #' \
		'#   text/plain bodies for some 400/404/501 cases, so Kiota cannot       #' \
		'#   generate typed error models for them.                               #' \
		'#                                                                       #' \
		'-------------------------------------------------------------------------'

ssclient-deps: # @HELP Show the Kiota generation dependency graph.
	KIOTA_TUTORIAL_ENABLED=false mise exec -- kiota info --openapi typespec/tsp-output/@typespec/openapi3/openapi.v1.yaml --language Go

typespec: # @HELP Install TypeSpec deps and compile the API spec.
	mise exec -- npm --prefix=$(CURDIR)/typespec clean-install
	mise exec -- npm --prefix=$(CURDIR)/typespec run compile

gen: # @HELP Generate code.
gen: typespec ssclient examplemocks

help: # @HELP Print this message.
help:
	echo "TARGETS:"
	grep -E '^.*: *# *@HELP' Makefile             \
	    | awk '                                   \
	        BEGIN {FS = ": *# *@HELP"};           \
	        { printf "  %-30s %s\n", $$1, $$2 };  \
	    '
