import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

using TypeSpec.Http;
using TypeSpec.Rest;

@service({
  title: "Archivematica Storage Service",
})
@route("/api/v2")
@useAuth(ApiKeyAuth<ApiKeyLocation.header, "ApiKey">)
namespace Service {
  @route("/pipeline")
  namespace Pipelines {
    @get
    op read(@path uuid: string): {
      @statusCode statusCode: 200;
      @body pipeline: Pipeline;
    } | {
      @statusCode statusCode: 400;
      @body error: Error;
    } | {
      @statusCode statusCode: 404;
    };
  }
  @route("/location")
  namespace Locations {
    @get
    op list(): {
      @statusCode statusCode: 200;
      @body list: LocationList;
    };

    @get
    op read(@path uuid: uuid): {
      @statusCode statusCode: 200;
      @body location: Location;
    } | {
      @statusCode statusCode: 400;
      @body error: Error;
    } | {
      @statusCode statusCode: 404;
    };

    @route("/default/{purpose}")
    @get
    op default(@path purpose: LocationPurpose): {
      @statusCode statusCode: 302;
    } | {
      @statusCode statusCode: 400;
      @body error: Error;
    } | {
      @statusCode statusCode: 404;
    };

    @doc("Move files to the specified location.")
    @post
    op move(@path uuid: uuid): {
      @doc("URI of the Location the files should be moved from.")
      origin_location: uri;

      @doc("URI of the")
      pipeline: uri;

      @doc("List of objects containing `source` and `destination`. The source and destination are paths relative to their Location of the files to be moved.")
      files: moveFilesList;
    } | {
      @statusCode statusCode: 400;
      @body error: Error;
    } | {
      @statusCode statusCode: 404;
    };
  }
}

@format("UUIDv4, e.g.: \"e3bae040-7682-4457-aef8-20a419aed0bd\".")
scalar uuid extends string;

@doc("Uniform Resource Identifier, e.g. \"/api/v2/space/141593ff-2a27-44a1-9de1-917573fa0f4a/\".")
scalar uri extends string;

model Pipeline {
  description: string;
  remote_name: string;
  resource_uri: string;
  uuid: uuid;
}

enum LocationPurpose {
  AIP_RECOVERY: "AR",
  AIP_STORAGE: "AS",
  CURRENTLY_PROCESSING: "CP",
  DIP_STORAGE: "DS",
  SWORD_DEPOSIT: "SD",
  STORAGE_SERVICE_INTERNAL: "SS",
  BACKLOG: "BL",
  TRANSFER_SOURCE: "TS",
  REPLICATOR: "RP",
}

model Location {
  description: string;
  enabled: boolean;
  path: string;
  pipeline: string[];
  purpose: LocationPurpose;
  quota: integer | null;
  relative_path: string;
  resource_uri: string;
  space: string;
  used: integer;
  uuid: uuid;
}

model LocationList {
  ...ListResponse<Location>;
}

alias moveFilesList = moveFile[];

model moveFile {
  source: string;
  destination: string;
}

model ListResponseMeta {
  limit: integer;
  next: string | null;
  offset: integer;
  previous: string | null;
  total_count: integer;
}

model ListResponse<Item> {
  meta: ListResponseMeta;
  objects: Item[];
}

@error
model Error {
  error: string;
}
