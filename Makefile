all: gen

.PHONY: deps
deps:
	mise exec -- go list -u -m -json all | go-mod-outdated -update -direct

.PHONY: examplemocks
examplemocks: $(MOCKGEN)
	cd $(CURDIR)/example && mise exec -- mockgen -typed -destination=./adapter/adapter.go -package=adapter github.com/microsoft/kiota-abstractions-go RequestAdapter

.PHONY: ssclient
ssclient:
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

.PHONY: ssclient-deps
ssclient-deps:
	KIOTA_TUTORIAL_ENABLED=false mise exec -- kiota info --openapi typespec/tsp-output/@typespec/openapi3/openapi.v1.yaml --language Go

.PHONY: typespec
typespec:
	mise exec -- npm --prefix=$(CURDIR)/typespec clean-install
	mise exec -- npm --prefix=$(CURDIR)/typespec run compile

.PHONY: gen
gen: typespec ssclient examplemocks
