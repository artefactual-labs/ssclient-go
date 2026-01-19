all: gen

.PHONY: deps
deps:
	mise exec -- go list -u -m -json all | go-mod-outdated -update -direct

.PHONY: examplemocks
examplemocks: $(MOCKGEN)
	cd $(CURDIR)/example && mise exec -- mockgen -typed -destination=./adapter/adapter.go -package=adapter github.com/microsoft/kiota-abstractions-go RequestAdapter

.PHONY: ssclient
ssclient:
	mise exec -- kiota generate --language go --clean-output --class-name Client --namespace-name go.artefactual.dev/ssclient/kiota --openapi typespec/tsp-output/@typespec/openapi3/openapi.v1.yaml --output ./kiota

.PHONY: typespec
typespec:
	mise exec -- npm --prefix=$(CURDIR)/typespec clean-install
	mise exec -- npm --prefix=$(CURDIR)/typespec run compile

.PHONY: gen
gen: typespec ssclient examplemocks
