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

## OpenAPI specification

This module was partially generated using an API client generator ([Kiota]) and
the [OpenAPI-described API][openapi-schema], which we've built using [TypeSpec].
The API has not been fully described yet, but we'll be extending support as
needed.

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
