# ssclient

[![PkgGoDev](https://pkg.go.dev/badge/go.artefactual.dev/ssclient)](https://pkg.go.dev/go.artefactual.dev/ssclient)
[![OpenAPI docs](https://img.shields.io/badge/openapi-docs-orange?logo=openapiinitiative&color=6BA539)][openapi-docs]

This repository provides the `go.artefactual.dev/ssclient` module. It does not
provide functionality beyond making the underlying REST API available.

The API is still **experimental**, breaking changes MAY occur.

## Usage

Check out [`example`], a small program that imports this module to print a list
of locations in Archivematica Storage Service.

### Working with the high-level client

The main entrypoint is `ssclient.New`, which constructs a client with small,
domain-oriented helpers such as `Locations()`, `Packages()`, and
`Pipelines()`. This is the recommended way to use the module for common
operations.

### Working with the generated client

This repository also ships a generated client based on the project's OpenAPI
description. If you need an endpoint that the high-level wrapper does not
expose yet, `Client.Raw()` returns that lower-level client as an escape hatch.

That generated API uses Kiota's request-builder pattern, including
`EmptyPathSegment()` for Storage Service endpoints that require a trailing
slash. See [`example`] for a side-by-side example using both the high-level
client and the generated client.

The generated client also inherits some Kiota limitations. The high-level
wrapper exists in part to smooth over those rough edges, so prefer it when a
wrapper method is available. For background, see [issue #20].

## OpenAPI specification

We use [TypeSpec] to describe the Storage Service API as an OpenAPI schema.

You can browse the schema in two forms:

- [Interactive OpenAPI docs][openapi-docs]
- [Raw OpenAPI YAML][openapi-schema]

The description is still incomplete, and we expect to extend it over time as we
cover more of the API.

### Notes

The Storage Service API itself is built with [TastyPie], which provides schema
introspection endpoints such as:

```sh
curl http://127.0.0.1:62081/api/v2/?fullschema=true
```

An example of that output is available at [ss-schema.json]. That data has been
useful as reference material while building the TypeSpec description.

If we want to expand coverage further, [django-tastypie-swagger] may be useful
prior art because it already maps Tastypie resources into Swagger-style schema
data.

[`example`]: ./example/main.go
[TypeSpec]: https://typespec.io
[openapi-schema]: https://raw.githubusercontent.com/artefactual-labs/ssclient-go/main/typespec/tsp-output/%40typespec/openapi3/openapi.v1.yaml
[openapi-docs]: https://artefactual-labs.github.io/ssclient-go/
[TastyPie]: https://django-tastypie.readthedocs.io/
[ss-schema.json]: https://gist.github.com/sevein/379f101ab9305235844c1e5101eeba04
[django-tastypie-swagger]: https://github.com/concentricsky/django-tastypie-swagger
[issue #20]: https://github.com/artefactual-labs/ssclient-go/issues/20
