# ssclient

[![PkgGoDev](https://pkg.go.dev/badge/go.artefactual.dev/ssclient)](https://pkg.go.dev/go.artefactual.dev/ssclient)
[![OpenAPI spec](https://img.shields.io/badge/openapi-spec-orange?logo=openapiinitiative&color=6BA539)][openapi-schema-editor]

This repository provides the `go.artefactual.dev/ssclient` module.

The API is still **experimental**, breaking changes MAY occur.

## Inspecting the API schema

This repository describes the Archivematica Storage Service API using TypeSpec.
The OpenAPI specification has been generated and is available at
[openapi.v1.yaml][openapi-schema]. It is not yet complete but we'll be extending
support as needed.

The API service is built with [TastyPie] which includes built-in schema
inspection capabilities which has been instrumental for this project, for
example:

    curl http://127.0.0.1:62081/api/v2/?fullschema=true

Visit [ss-schema.json] to see the output. We should explore options for
using this feature. [django-tastypie-swagger] could be a really good start since
it's already doing all the mapping. TypeSpec could be a target by using the
[emitter framework].

[openapi-schema]: https://raw.githubusercontent.com/artefactual-labs/ssclient-go/main/typespec/tsp-output/%40typespec/openapi3/openapi.v1.yaml
[openapi-schema-editor]: https://editor.swagger.io/?url=https://raw.githubusercontent.com/artefactual-labs/ssclient-go/main/typespec/tsp-output/%40typespec/openapi3/openapi.v1.yaml
[TastyPie]: https://django-tastypie.readthedocs.io/
[ss-schema.json]: https://gist.github.com/sevein/379f101ab9305235844c1e5101eeba04
[django-tastypie-swagger]: https://github.com/concentricsky/django-tastypie-swagger
[emitter framework]: https://typespec.io/docs/next/extending-typespec/emitter-framework
