# ssclient

[![PkgGoDev](https://pkg.go.dev/badge/go.artefactual.dev/ssclient)](https://pkg.go.dev/go.artefactual.dev/ssclient)
[![OpenAPI docs](https://img.shields.io/badge/openapi-docs-orange?logo=openapiinitiative&color=6BA539)][openapi-docs]

This repository provides the `go.artefactual.dev/ssclient` module. It does not
provide functionality beyond making the underlying REST API available.

The API is still **experimental**, breaking changes MAY occur.

## Usage

Check out [`example`], a small program that imports this module to print a list
of locations in Archivematica Storage Service.

For a more detailed example, refer to CCP's [ssclient package][ccp-ssclient],
which offers additional features such as retrieving default locations through
header inspection, paging results, and more. This could potentially become a
separate package in this repository.

## Working with the generated client

This repository ships the raw Kiota output. The generated API follows Kiota's
builder pattern, so you'll often see fluent chains such as:

```go
client, _ := ssclient.New(http.DefaultClient, url, user, key)
locations, _ := client.Api().V2().Location().EmptyPathSegment().
    Get(ctx, &api.V2LocationEmptyPathSegmentRequestBuilderGetRequestConfiguration{})

fixity, _ := client.Api().V2().File().ByUuid(id).
    CheckFixity().EmptyPathSegment().
    Get(ctx, &api.V2FileItemCheckFixityEmptyPathSegmentRequestBuilderGetRequestConfiguration{})
```

Endpoints in Archivematica Storage Service expect a trailing slash, so Kiota
represents that final `/` as an additional fluent step named
`EmptyPathSegment()`. If you're new to Kiota, take a look at the
[`example`](./example/main.go) program for a complete, working walkthrough.
For a friendlier experience you can wrap these builders in your own helper
functions, but we keep the generated tree here so you can choose the style that
fits your project.

Kiota itself aims to provide strongly typed building blocks—request builders,
models, middleware hooks—without committing to a domain-specific shape. Most
teams layer thin helpers or full SDKs on top to present a friendlier surface.
We would like to do the same here by adding a small Go package that wraps the
generated client (so you can call `pkg.CheckFixity(...)` instead of poking
through the builder chain), but that work has not happened yet.

## OpenAPI specification

This module was partially generated using an API client generator ([Kiota]) and
the [OpenAPI-described API][openapi-schema], which we've built using [TypeSpec].
You can browse an interactive rendering of that schema at [our published
OpenAPI docs][openapi-docs]. The API has not been fully described yet, but
we'll be extending support as needed.

Furthermore, the API service is built with [TastyPie], and old webservice API
framework for Django. It includes built-in schema inspection capabilities which
has been instrumental for this project, e.g.:

    curl http://127.0.0.1:62081/api/v2/?fullschema=true

Visit [ss-schema.json] to see the output. We should explore options for
using this feature in order to describe the API further.
[django-tastypie-swagger] could be a really good start since it's already
doing all the mapping. TypeSpec could be a target by using the [emitter
framework].

[`example`]: ./example/main.go
[Kiota]: https://learn.microsoft.com/en-us/openapi/kiota/overview
[TypeSpec]: https://typespec.io
[openapi-schema]: https://raw.githubusercontent.com/artefactual-labs/ssclient-go/main/typespec/tsp-output/%40typespec/openapi3/openapi.v1.yaml
[openapi-docs]: https://artefactual-labs.github.io/ssclient-go/
[TastyPie]: https://django-tastypie.readthedocs.io/
[ss-schema.json]: https://gist.github.com/sevein/379f101ab9305235844c1e5101eeba04
[django-tastypie-swagger]: https://github.com/concentricsky/django-tastypie-swagger
[emitter framework]: https://typespec.io/docs/next/extending-typespec/emitter-framework
[ccp-ssclient]: https://github.com/artefactual-labs/ccp/tree/ccp/hack/ccp/internal/ssclient
