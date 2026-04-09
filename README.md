# ssclient

[![PkgGoDev](https://pkg.go.dev/badge/go.artefactual.dev/ssclient)](https://pkg.go.dev/go.artefactual.dev/ssclient)
[![OpenAPI docs](https://img.shields.io/badge/openapi-docs-orange?logo=openapiinitiative&color=6BA539)][openapi-docs]

This repository provides the `go.artefactual.dev/ssclient` module. It does not
provide functionality beyond making the underlying REST API available.

The API is still **experimental**, breaking changes MAY occur.

## Usage

Check out [`example`], a small program that imports this module to print a list
of locations and pipelines found in Archivematica Storage Service.

### Working with the high-level client

The main entrypoint is `ssclient.New`, which constructs a client with small,
domain-oriented helpers such as `Locations()`, `Packages()`, and `Pipelines()`.
This is the recommended way to use the module and should be the default path for
interacting with the API.

### Working with the generated client

This repository also ships a generated client based on the project's OpenAPI
description from the separate
[Archivematica Storage Service API Specification][spec-repo] repository. If you
need an endpoint that the high-level wrapper does not expose yet, `Client.Raw()`
returns that lower-level client as an escape hatch.

That generated API uses Kiota's request-builder pattern, including
`EmptyPathSegment()` for Storage Service endpoints that require a trailing
slash. See [`example`] for a side-by-side example using both the high-level
client and the generated client.

> [!WARNING]
> Prefer the high-level wrapper for normal use. `Client.Raw()` is an escape
> hatch for gaps in wrapper coverage while we continue to define that boundary.
> For background, see [issue #20].

## OpenAPI specification

The Storage Service API description lives in the separate
[Archivematica Storage Service API Specification][spec-repo] repository. That
repository owns the TypeSpec source, OpenAPI output, and related background
notes for the wire contract consumed by this module.

You can browse the published schema there via
[Interactive OpenAPI docs][openapi-docs] or [Raw OpenAPI YAML][openapi-schema].

## Contributor notes

This repository is wrapper-first. Generated code is supporting infrastructure
and a fallback escape hatch, not the primary interface we want callers to use.

The generated client in this repository depends on the `spec/` git submodule.
Initialize it before running generation commands:

```sh
git submodule update --init --recursive
```

When an endpoint's wire contract needs to be added or changed, the preferred
pattern is:

1. Update the specification in
   [archivematica-storage-service-api-specification][spec-repo] first so the
   OpenAPI remains accurate.
1. Regenerate Kiota without hand-editing generated files.
1. Expose the operation through the public client wrapper and treat the
   generated Go surface as supporting infrastructure.

Examples of wrapper-level normalization in this repository include:

- `Packages.Move`, where the client handles Storage Service's non-idiomatic
  package-move workflow while keeping the generated surface available as a
  fallback.
- `Packages.DeleteAIP`, where distinct non-error outcomes are preserved in a
  typed result instead of being flattened into one generic success path.
- `Packages.ReviewAIPDeletion`, where application-level failures can be raised
  as a typed Go error even when the server responds with HTTP `200`.
- `AsyncService.Wait`, where the wrapper builds a more Go-oriented polling API
  on top of the underlying async task resource.

[issue #20]: https://github.com/artefactual-labs/ssclient-go/issues/20
[openapi-docs]: https://editor.swagger.io/?url=https://raw.githubusercontent.com/archivematica/archivematica-storage-service-api-specification/main/openapi.v1.yaml
[openapi-schema]: https://raw.githubusercontent.com/archivematica/archivematica-storage-service-api-specification/main/openapi.v1.yaml
[spec-repo]: https://github.com/archivematica/archivematica-storage-service-api-specification
[`example`]: ./example/main.go
