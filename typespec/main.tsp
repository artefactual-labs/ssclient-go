import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";
import "@typespec/versioning";

using TypeSpec.Http;
using TypeSpec.Rest;
using TypeSpec.Versioning;

@service({
  title: "Archivematica Storage Service",
})
@route("/api/v2/")
@useAuth(ApiKeyAuth<ApiKeyLocation.header, "ApiKey">)
@versioned(Versions)
namespace Service {
  enum Versions {
    v1,
  }

  @route("/pipeline/")
  namespace Pipelines {
    @get
    op list(
      @query
      description?: string,

      @query
      uuid?: uuid,
    ): {
      @statusCode statusCode: 200;
      @body list: PipelineList;
    };

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

  @route("/location/")
  namespace Locations {
    @get
    op list(
      @query
      description?: string,

      @query
      purpose?: LocationPurpose,

      @query
      quota?: integer,

      @query("relative_path")
      relativePath?: string,

      @query
      used?: integer,

      @query
      uuid?: uuid,

      // TODO: enable all ORM filters, including across relationships.
      // "space": ALL_WITH_RELATIONS
      // "pipeline": ALL_WITH_RELATIONS

      @query("pipeline__uuid")
      pipelineId?: uuid,

      ...ListOpQueries,
    ): {
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

    @route("/default/{purpose}/")
    @get
    op default(@path purpose: LocationPurpose): {
      @statusCode
      statusCode: 302;

      @header("Location")
      location: string;
    } | {
      @statusCode
      statusCode: 400;

      @body
      error: Error;
    } | {
      @statusCode
      statusCode: 404;
    };

    @doc("Move files to the specified location.")
    @post
    op move(
      @path
      uuid: uuid,

      @body
      body: MoveRequest,
    ): {
      @statusCode statusCode: 200;
    } | {
      @statusCode statusCode: 400;
      @body error: Error;
    } | {
      @statusCode statusCode: 404;
    };
  }

  @route("/file/")
  namespace Packages {
    @get
    op read(@path uuid: uuid): {
      @statusCode statusCode: 200;
      @body package: Package;
    } | {
      @statusCode statusCode: 400;
      @body error: Error;
    } | {
      @statusCode statusCode: 404;
    };

    @doc("Move the package to a different storage location.")
    @route("{uuid}/move/")
    @post
    op move(
      @path
      uuid: uuid,

      @body
      body: PackageMoveRequest,
    ): {
      @statusCode statusCode: 200;
    } | {
      @statusCode statusCode: 400;
      @body error: Error;
    } | {
      @statusCode statusCode: 404;
    };

    @doc("Run a fixity check for the specified package.")
    @route("{uuid}/check_fixity/")
    @get
    op checkFixity(
      @path
      uuid: uuid,

      @query("force_local")
      forceLocal?: boolean,
    ): {
      @statusCode statusCode: 200;
      @body result: FixityResponse;
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

model ListOpQueries {
  @query("order_by")
  orderBy?: string;

  @query
  limit?: integer;

  @query
  offset?: integer;
}

model Pipeline {
  description: string;
  remote_name: string;
  resource_uri: string;
  uuid: uuid;
}

model PipelineList {
  ...ListResponse<Pipeline>;
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

model MoveRequest {
  @doc("URI of the Location the files should be moved from.")
  origin_location: uri;

  @doc("URI of the pipeline.")
  pipeline: uri;

  @doc("List of objects containing `source` and `destination`. The source and destination are paths relative to their Location of the files to be moved.")
  files: moveFilesList;
}

alias moveFilesList = moveFile[];

model Package {
  uuid: uuid;
  current_full_path: string;
  current_location: uri;
  current_path: string;
  encrypted: boolean;
  origin_pipeline: uri;
  package_type: PackageType;
  related_packages: string[];
  replicas: string[];
  replicated_package: string | null;
  resource_uri: uri;
  size: integer;
  status: string;
  stored_date: string | null;
}

model PackageMoveRequest {
  location_uuid: uuid;
}

model FixityResponse {
  success: boolean | null;
  message: string;
  timestamp: string | null;
  failures?: FixityFailures;
}

model FixityFailures {
  files?: FixityFailureFiles;
}

model FixityFailureFiles {
  missing?: string[];
  changed?: string[];
  untracked?: string[];
}

enum PackageType {
  AIP: "AIP";
  AIC: "AIC";
  DIP: "DIP";
  TRANSFER: "transfer";
  SIP: "SIP";
  FILE: "file";
  DEPOSIT: "deposit";
}

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
