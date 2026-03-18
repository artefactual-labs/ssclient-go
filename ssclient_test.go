package ssclient_test

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	kabs "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"

	"go.artefactual.dev/ssclient"
	"go.artefactual.dev/ssclient/kiota/models"
)

type fakeRequestAdapter struct {
	send                       func(context.Context, *kabs.RequestInformation, serialization.ParsableFactory, kabs.ErrorMappings) (serialization.Parsable, error)
	sendPrimitive              func(context.Context, *kabs.RequestInformation, string, kabs.ErrorMappings) (any, error)
	convertToNativeRequest     func(context.Context, *kabs.RequestInformation) (any, error)
	baseURL                    string
	serializationWriterFactory serialization.SerializationWriterFactory
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func (f *fakeRequestAdapter) Send(ctx context.Context, requestInfo *kabs.RequestInformation, constructor serialization.ParsableFactory, errorMappings kabs.ErrorMappings) (serialization.Parsable, error) {
	if f.send == nil {
		panic("unexpected Send call")
	}
	return f.send(ctx, requestInfo, constructor, errorMappings)
}

func (f *fakeRequestAdapter) SendEnum(context.Context, *kabs.RequestInformation, serialization.EnumFactory, kabs.ErrorMappings) (any, error) {
	panic("unexpected SendEnum call")
}

func (f *fakeRequestAdapter) SendCollection(context.Context, *kabs.RequestInformation, serialization.ParsableFactory, kabs.ErrorMappings) ([]serialization.Parsable, error) {
	panic("unexpected SendCollection call")
}

func (f *fakeRequestAdapter) SendEnumCollection(context.Context, *kabs.RequestInformation, serialization.EnumFactory, kabs.ErrorMappings) ([]any, error) {
	panic("unexpected SendEnumCollection call")
}

func (f *fakeRequestAdapter) SendPrimitive(ctx context.Context, requestInfo *kabs.RequestInformation, typeName string, errorMappings kabs.ErrorMappings) (any, error) {
	if f.sendPrimitive == nil {
		panic("unexpected SendPrimitive call")
	}
	return f.sendPrimitive(ctx, requestInfo, typeName, errorMappings)
}

func (f *fakeRequestAdapter) SendPrimitiveCollection(context.Context, *kabs.RequestInformation, string, kabs.ErrorMappings) ([]any, error) {
	panic("unexpected SendPrimitiveCollection call")
}

func (f *fakeRequestAdapter) SendNoContent(context.Context, *kabs.RequestInformation, kabs.ErrorMappings) error {
	panic("unexpected SendNoContent call")
}

func (f *fakeRequestAdapter) GetSerializationWriterFactory() serialization.SerializationWriterFactory {
	return f.serializationWriterFactory
}

func (f *fakeRequestAdapter) EnableBackingStore(store.BackingStoreFactory) {}

func (f *fakeRequestAdapter) SetBaseUrl(baseURL string) {
	f.baseURL = baseURL
}

func (f *fakeRequestAdapter) GetBaseUrl() string {
	return f.baseURL
}

func (f *fakeRequestAdapter) ConvertToNativeRequest(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
	if f.convertToNativeRequest == nil {
		panic("unexpected ConvertToNativeRequest call")
	}
	return f.convertToNativeRequest(ctx, requestInfo)
}

func assertEqual(t *testing.T, got, want interface{}) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Mismatch found:\nGot: %#v\nWant: %#v", got, want)
	}
}

func TestClientRaw(t *testing.T) {
	client, err := ssclient.New(ssclient.Config{
		BaseURL:    "http://storage.service",
		Username:   "test",
		Key:        "test",
		HTTPClient: &http.Client{},
	})
	assertEqual(t, err, nil)

	if client.Raw() == nil {
		t.Fatal("expected raw client")
	}
	if client.Packages() == nil {
		t.Fatal("expected packages service")
	}
}

func TestPackagesGet(t *testing.T) {
	const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"

	client, err := ssclient.New(ssclient.Config{
		BaseURL:    "http://storage.service",
		Username:   "test",
		Key:        "test",
		HTTPClient: &http.Client{},
	})
	assertEqual(t, err, nil)

	raw := client.Raw()
	raw.RequestAdapter = &fakeRequestAdapter{
		baseURL: "http://storage.service",
		send: func(ctx context.Context, requestInfo *kabs.RequestInformation, constructor serialization.ParsableFactory, errorMappings kabs.ErrorMappings) (serialization.Parsable, error) {
			if requestInfo == nil {
				t.Fatal("expected request info")
			}

			assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/file/{uuid}")
			assertEqual(t, requestInfo.PathParameters["uuid"], packageID)
			assertEqual(t, requestInfo.Headers.Get("Accept"), []string{"application/json"})

			pkg := models.NewPackageEscaped()
			pkg.SetUuid(ptr(packageID))
			pkg.SetStatus(ptr("UPLOADED"))
			pkg.SetCurrentLocation(ptr("/api/v2/location/154660b9-b4a3-4886-8d68-5e170c0923b8/"))
			pkg.SetReplicas([]string{
				"/api/v2/file/96922350-ccde-4fb0-a999-d2010522028f/",
				"/api/v2/file/610bc407-ba6c-4dcd-8675-d2727a9aab18/",
			})

			return pkg, nil
		},
	}

	pkg, err := client.Packages().Get(context.Background(), packageID)
	assertEqual(t, err, nil)
	if pkg == nil {
		t.Fatal("expected package")
	}

	assertEqual(t, *pkg.GetUuid(), packageID)
	assertEqual(t, *pkg.GetStatus(), "UPLOADED")
	assertEqual(t, len(pkg.GetReplicas()), 2)
	assertEqual(t, *pkg.GetCurrentLocation(), "/api/v2/location/154660b9-b4a3-4886-8d68-5e170c0923b8/")
}

func TestPackagesDeleteAIPAccepted(t *testing.T) {
	const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"

	client, err := ssclient.New(ssclient.Config{
		BaseURL:  "http://storage.service",
		Username: "test",
		Key:      "test",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				assertEqual(t, r.Method, http.MethodPost)
				assertEqual(t, r.URL.String(), "http://storage.service/api/v2/file/"+packageID+"/delete_aip/")
				assertEqual(t, r.Header.Get("Accept"), "application/json")
				assertEqual(t, r.Header.Get("Authorization"), "ApiKey test:test")
				assertEqual(t, r.Header.Get("Content-Type"), "application/json")

				return &http.Response{
					StatusCode: http.StatusAccepted,
					Header:     http.Header{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{"message":"Delete request created successfully.","id":17}`)),
					Request:    r,
				}, nil
			}),
		},
	})
	assertEqual(t, err, nil)

	raw := client.Raw()
	writerFactory := raw.RequestAdapter.GetSerializationWriterFactory()
	raw.RequestAdapter = &fakeRequestAdapter{
		serializationWriterFactory: writerFactory,
		convertToNativeRequest: func(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
			assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/file/{uuid}/delete_aip/")
			assertEqual(t, requestInfo.PathParameters["uuid"], packageID)
			assertEqual(t, requestInfo.Headers.Get("Accept"), []string{"application/json"})
			assertEqual(t, requestInfo.Headers.Get("Content-Type"), []string{"application/json"})

			got := string(requestInfo.Content)
			if !strings.Contains(got, `"event_reason":"Delete please!"`) {
				t.Fatalf("unexpected request body %q", got)
			}
			if !strings.Contains(got, `"pipeline":"4b9e8af5-b0af-4abf-80b8-4b7d76281f61"`) {
				t.Fatalf("unexpected request body %q", got)
			}
			if !strings.Contains(got, `"user_id":1`) {
				t.Fatalf("unexpected request body %q", got)
			}
			if !strings.Contains(got, `"user_email":"user@example.com"`) {
				t.Fatalf("unexpected request body %q", got)
			}

			return &http.Request{
				Method: http.MethodPost,
				URL: &url.URL{
					Scheme: "http",
					Host:   "storage.service",
					Path:   "/api/v2/file/" + packageID + "/delete_aip/",
				},
				Header: http.Header{
					"Accept":        {"application/json"},
					"Authorization": {"ApiKey test:test"},
					"Content-Type":  {"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(got)),
			}, nil
		},
	}

	body := models.NewDeleteAipRequest()
	body.SetEventReason(ptr("Delete please!"))
	body.SetPipeline(ptr("4b9e8af5-b0af-4abf-80b8-4b7d76281f61"))
	body.SetUserId(ptr(int32(1)))
	body.SetUserEmail(ptr("user@example.com"))

	res, err := client.Packages().DeleteAIP(context.Background(), packageID, body)
	assertEqual(t, err, nil)
	if res == nil {
		t.Fatal("expected delete AIP result")
	}
	if !res.IsAccepted() {
		t.Fatal("expected accepted delete AIP result")
	}
	if res.HasExistingRequest() {
		t.Fatal("did not expect existing request result")
	}
	assertEqual(t, res.StatusCode, http.StatusAccepted)
	assertEqual(t, res.Accepted.Message, "Delete request created successfully.")
	assertEqual(t, res.Accepted.ID, int32(17))
}

func TestPackagesDeleteAIPAlreadyExists(t *testing.T) {
	const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"

	client, err := ssclient.New(ssclient.Config{
		BaseURL:  "http://storage.service",
		Username: "test",
		Key:      "test",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": {"application/json"}},
					Body:       io.NopCloser(strings.NewReader(`{"error_message":"A deletion request already exists for this AIP."}`)),
					Request:    r,
				}, nil
			}),
		},
	})
	assertEqual(t, err, nil)

	raw := client.Raw()
	writerFactory := raw.RequestAdapter.GetSerializationWriterFactory()
	raw.RequestAdapter = &fakeRequestAdapter{
		serializationWriterFactory: writerFactory,
		convertToNativeRequest: func(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
			return &http.Request{
				Method: http.MethodPost,
				URL: &url.URL{
					Scheme: "http",
					Host:   "storage.service",
					Path:   "/api/v2/file/" + packageID + "/delete_aip/",
				},
				Header: http.Header{
					"Accept":        {"application/json"},
					"Authorization": {"ApiKey test:test"},
					"Content-Type":  {"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(string(requestInfo.Content))),
			}, nil
		},
	}

	body := models.NewDeleteAipRequest()
	body.SetEventReason(ptr("Delete please!"))
	body.SetPipeline(ptr("4b9e8af5-b0af-4abf-80b8-4b7d76281f61"))
	body.SetUserId(ptr(int32(1)))
	body.SetUserEmail(ptr("user@example.com"))

	res, err := client.Packages().DeleteAIP(context.Background(), packageID, body)
	assertEqual(t, err, nil)
	if res == nil {
		t.Fatal("expected delete AIP result")
	}
	if res.IsAccepted() {
		t.Fatal("did not expect accepted delete AIP result")
	}
	if !res.HasExistingRequest() {
		t.Fatal("expected existing request result")
	}
	assertEqual(t, res.StatusCode, http.StatusOK)
	assertEqual(t, res.AlreadyExists.ErrorMessage, "A deletion request already exists for this AIP.")
}

func TestPackagesDeleteAIPNilBody(t *testing.T) {
	client, err := ssclient.New(ssclient.Config{
		BaseURL:    "http://storage.service",
		Username:   "test",
		Key:        "test",
		HTTPClient: &http.Client{},
	})
	assertEqual(t, err, nil)

	if _, err := client.Packages().DeleteAIP(context.Background(), "pkg", nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestPackagesCheckFixity(t *testing.T) {
	const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"

	client, err := ssclient.New(ssclient.Config{
		BaseURL:    "http://storage.service",
		Username:   "test",
		Key:        "test",
		HTTPClient: &http.Client{},
	})
	assertEqual(t, err, nil)

	raw := client.Raw()
	forceLocal := true
	raw.RequestAdapter = &fakeRequestAdapter{
		baseURL: "http://storage.service",
		send: func(ctx context.Context, requestInfo *kabs.RequestInformation, constructor serialization.ParsableFactory, errorMappings kabs.ErrorMappings) (serialization.Parsable, error) {
			assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/file/{uuid}/check_fixity/{?force_local}")
			assertEqual(t, requestInfo.PathParameters["uuid"], packageID)
			assertEqual(t, requestInfo.QueryParameters["force_local"], "true")
			assertEqual(t, requestInfo.Headers.Get("Accept"), []string{"application/json"})

			response := models.NewFixityResponse()
			response.SetSuccess(ptr(true))
			response.SetMessage(ptr("ok"))
			return response, nil
		},
	}

	res, err := client.Packages().CheckFixity(context.Background(), packageID, ssclient.CheckFixityOptions{
		ForceLocal: &forceLocal,
	})
	assertEqual(t, err, nil)
	if res == nil {
		t.Fatal("expected fixity response")
	}
	assertEqual(t, *res.GetSuccess(), true)
	assertEqual(t, *res.GetMessage(), "ok")
}

func TestPackagesMove(t *testing.T) {
	const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"
	const locationID = "154660b9-b4a3-4886-8d68-5e170c0923b8"

	client, err := ssclient.New(ssclient.Config{
		BaseURL:    "http://storage.service",
		Username:   "test",
		Key:        "test",
		HTTPClient: &http.Client{},
	})
	assertEqual(t, err, nil)

	raw := client.Raw()
	writerFactory := raw.RequestAdapter.GetSerializationWriterFactory()
	raw.RequestAdapter = &fakeRequestAdapter{
		baseURL:                    "http://storage.service",
		serializationWriterFactory: writerFactory,
		sendPrimitive: func(ctx context.Context, requestInfo *kabs.RequestInformation, typeName string, errorMappings kabs.ErrorMappings) (any, error) {
			assertEqual(t, typeName, "[]byte")
			assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/file/{uuid}/move/")
			assertEqual(t, requestInfo.PathParameters["uuid"], packageID)
			assertEqual(t, requestInfo.Headers.Get("Accept"), []string{"application/json"})
			assertEqual(t, requestInfo.Headers.Get("Content-Type"), []string{"application/json"})
			if got := string(requestInfo.Content); !strings.Contains(got, `"location_uuid":"`+locationID+`"`) {
				t.Fatalf("unexpected request body %q", got)
			}

			return nil, nil
		},
	}

	body := models.NewPackageMoveRequest()
	body.SetLocationUuid(ptr(locationID))
	assertEqual(t, client.Packages().Move(context.Background(), packageID, body), nil)
}

func TestPackagesMoveNilBody(t *testing.T) {
	client, err := ssclient.New(ssclient.Config{
		BaseURL:    "http://storage.service",
		Username:   "test",
		Key:        "test",
		HTTPClient: &http.Client{},
	})
	assertEqual(t, err, nil)

	if err := client.Packages().Move(context.Background(), "pkg", nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestPackagesReviewAIPDeletion(t *testing.T) {
	const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"

	client, err := ssclient.New(ssclient.Config{
		BaseURL:    "http://storage.service",
		Username:   "test",
		Key:        "test",
		HTTPClient: &http.Client{},
	})
	assertEqual(t, err, nil)

	raw := client.Raw()
	writerFactory := raw.RequestAdapter.GetSerializationWriterFactory()
	raw.RequestAdapter = &fakeRequestAdapter{
		baseURL:                    "http://storage.service",
		serializationWriterFactory: writerFactory,
		send: func(ctx context.Context, requestInfo *kabs.RequestInformation, constructor serialization.ParsableFactory, errorMappings kabs.ErrorMappings) (serialization.Parsable, error) {
			assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/file/{uuid}/review_aip_deletion/")
			assertEqual(t, requestInfo.PathParameters["uuid"], packageID)
			assertEqual(t, requestInfo.Headers.Get("Accept"), []string{"application/json"})
			assertEqual(t, requestInfo.Headers.Get("Content-Type"), []string{"application/json"})

			got := string(requestInfo.Content)
			if !strings.Contains(got, `"event_id":99`) {
				t.Fatalf("unexpected request body %q", got)
			}
			if !strings.Contains(got, `"decision":"approve"`) {
				t.Fatalf("unexpected request body %q", got)
			}
			if !strings.Contains(got, `"reason":"approved by workflow"`) {
				t.Fatalf("unexpected request body %q", got)
			}

			res := models.NewReviewAipDeletionSuccess()
			res.SetMessage(ptr("done"))
			return res, nil
		},
	}

	body := models.NewReviewAipDeletionRequest()
	body.SetEventId(ptr(int32(99)))
	decision := models.APPROVE_REVIEWAIPDELETIONDECISION
	body.SetDecision(&decision)
	body.SetReason(ptr("approved by workflow"))

	res, err := client.Packages().ReviewAIPDeletion(context.Background(), packageID, body)
	assertEqual(t, err, nil)
	if res == nil {
		t.Fatal("expected review deletion response")
	}
	assertEqual(t, *res.GetMessage(), "done")
}

func TestPackagesReviewAIPDeletionNilBody(t *testing.T) {
	client, err := ssclient.New(ssclient.Config{
		BaseURL:    "http://storage.service",
		Username:   "test",
		Key:        "test",
		HTTPClient: &http.Client{},
	})
	assertEqual(t, err, nil)

	if _, err := client.Packages().ReviewAIPDeletion(context.Background(), "pkg", nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestLocationsList(t *testing.T) {
	client, err := ssclient.New(ssclient.Config{
		BaseURL:    "http://storage.service",
		Username:   "test",
		Key:        "test",
		HTTPClient: &http.Client{},
	})
	assertEqual(t, err, nil)

	raw := client.Raw()
	purpose := models.DS_LOCATIONPURPOSE
	raw.RequestAdapter = &fakeRequestAdapter{
		baseURL: "http://storage.service",
		send: func(ctx context.Context, requestInfo *kabs.RequestInformation, constructor serialization.ParsableFactory, errorMappings kabs.ErrorMappings) (serialization.Parsable, error) {
			assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/location/{?description,limit,offset,order_by,pipeline__uuid,purpose,quota,relative_path,used,uuid}")
			assertEqual(t, requestInfo.QueryParameters["description"], "test")
			assertEqual(t, requestInfo.QueryParametersAny["purpose"], "DS")
			assertEqual(t, requestInfo.QueryParameters["limit"], "10")
			assertEqual(t, requestInfo.QueryParameters["order_by"], "description")

			list := models.NewLocationList()
			location := models.NewLocation()
			location.SetUuid(ptr("fff70864-a5d4-4ca6-ab29-b4ce67d8eeab"))
			list.SetObjects([]models.Locationable{location})
			return list, nil
		},
	}

	description := "test"
	limit := int32(10)
	orderBy := "description"
	res, err := client.Locations().List(context.Background(), ssclient.ListLocationsQuery{
		Description: &description,
		Purpose:     &purpose,
		Limit:       &limit,
		OrderBy:     &orderBy,
	})
	assertEqual(t, err, nil)
	if res == nil {
		t.Fatal("expected location list")
	}
	assertEqual(t, len(res.GetObjects()), 1)
}

func TestLocationsGet(t *testing.T) {
	const locationID = "fff70864-a5d4-4ca6-ab29-b4ce67d8eeab"

	client, err := ssclient.New(ssclient.Config{
		BaseURL:    "http://storage.service",
		Username:   "test",
		Key:        "test",
		HTTPClient: &http.Client{},
	})
	assertEqual(t, err, nil)

	raw := client.Raw()
	raw.RequestAdapter = &fakeRequestAdapter{
		baseURL: "http://storage.service",
		send: func(ctx context.Context, requestInfo *kabs.RequestInformation, constructor serialization.ParsableFactory, errorMappings kabs.ErrorMappings) (serialization.Parsable, error) {
			assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/location/{uuid}")
			assertEqual(t, requestInfo.PathParameters["uuid"], locationID)

			location := models.NewLocation()
			location.SetUuid(ptr(locationID))
			return location, nil
		},
	}

	res, err := client.Locations().Get(context.Background(), locationID)
	assertEqual(t, err, nil)
	if res == nil {
		t.Fatal("expected location")
	}
	assertEqual(t, *res.GetUuid(), locationID)
}

func TestLocationsDefault(t *testing.T) {
	client, err := ssclient.New(ssclient.Config{
		BaseURL:  "http://storage.service",
		Username: "test",
		Key:      "test",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				assertEqual(t, r.Method, http.MethodGet)
				assertEqual(t, r.URL.String(), "http://storage.service/api/v2/location/default/DS/")
				assertEqual(t, r.Header.Get("Accept"), "application/json")
				assertEqual(t, r.Header.Get("Authorization"), "ApiKey test:test")

				return &http.Response{
					StatusCode: http.StatusFound,
					Header:     http.Header{"Location": {"/api/v2/location/154660b9-b4a3-4886-8d68-5e170c0923b8/"}},
					Body:       io.NopCloser(strings.NewReader("")),
					Request:    r,
				}, nil
			}),
		},
	})
	assertEqual(t, err, nil)

	adapter := &fakeRequestAdapter{
		convertToNativeRequest: func(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
			assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/location/default/{purpose}/")
			assertEqual(t, requestInfo.PathParameters["purpose"], "DS")
			return &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "http",
					Host:   "storage.service",
					Path:   "/api/v2/location/default/DS/",
				},
				Header: http.Header{
					"Accept":        {"application/json"},
					"Authorization": {"ApiKey test:test"},
				},
			}, nil
		},
	}
	client.Raw().RequestAdapter = adapter

	location, err := client.Locations().Default(context.Background(), models.DS_LOCATIONPURPOSE)
	assertEqual(t, err, nil)
	assertEqual(t, location, "/api/v2/location/154660b9-b4a3-4886-8d68-5e170c0923b8/")
}

func TestLocationsMove(t *testing.T) {
	const locationID = "fff70864-a5d4-4ca6-ab29-b4ce67d8eeab"

	client, err := ssclient.New(ssclient.Config{
		BaseURL:    "http://storage.service",
		Username:   "test",
		Key:        "test",
		HTTPClient: &http.Client{},
	})
	assertEqual(t, err, nil)

	raw := client.Raw()
	writerFactory := raw.RequestAdapter.GetSerializationWriterFactory()
	raw.RequestAdapter = &fakeRequestAdapter{
		serializationWriterFactory: writerFactory,
		sendPrimitive: func(ctx context.Context, requestInfo *kabs.RequestInformation, typeName string, errorMappings kabs.ErrorMappings) (any, error) {
			assertEqual(t, typeName, "[]byte")
			assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/location/{uuid}")
			assertEqual(t, requestInfo.PathParameters["uuid"], locationID)
			assertEqual(t, requestInfo.Headers.Get("Content-Type"), []string{"application/json"})

			got := string(requestInfo.Content)
			if !strings.Contains(got, `"origin_location":"/api/v2/location/origin/"`) {
				t.Fatalf("unexpected request body %q", got)
			}
			if !strings.Contains(got, `"pipeline":"/api/v2/pipeline/source/"`) {
				t.Fatalf("unexpected request body %q", got)
			}

			return nil, nil
		},
	}

	body := models.NewMoveRequest()
	body.SetOriginLocation(ptr("/api/v2/location/origin/"))
	body.SetPipeline(ptr("/api/v2/pipeline/source/"))
	assertEqual(t, client.Locations().Move(context.Background(), locationID, body), nil)
}

func TestLocationsMoveNilBody(t *testing.T) {
	client, err := ssclient.New(ssclient.Config{
		BaseURL:    "http://storage.service",
		Username:   "test",
		Key:        "test",
		HTTPClient: &http.Client{},
	})
	assertEqual(t, err, nil)

	if err := client.Locations().Move(context.Background(), "location", nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestPipelinesList(t *testing.T) {
	client, err := ssclient.New(ssclient.Config{
		BaseURL:    "http://storage.service",
		Username:   "test",
		Key:        "test",
		HTTPClient: &http.Client{},
	})
	assertEqual(t, err, nil)

	raw := client.Raw()
	raw.RequestAdapter = &fakeRequestAdapter{
		baseURL: "http://storage.service",
		send: func(ctx context.Context, requestInfo *kabs.RequestInformation, constructor serialization.ParsableFactory, errorMappings kabs.ErrorMappings) (serialization.Parsable, error) {
			assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/pipeline/{?description,uuid}")
			assertEqual(t, requestInfo.QueryParameters["description"], "pipeline")
			assertEqual(t, requestInfo.QueryParameters["uuid"], "a64e061a-5688-49b5-95c1-0b6885c40c04")

			list := models.NewPipelineList()
			pipeline := models.NewPipeline()
			pipeline.SetUuid(ptr("a64e061a-5688-49b5-95c1-0b6885c40c04"))
			list.SetObjects([]models.Pipelineable{pipeline})
			return list, nil
		},
	}

	description := "pipeline"
	uuid := "a64e061a-5688-49b5-95c1-0b6885c40c04"
	res, err := client.Pipelines().List(context.Background(), ssclient.ListPipelinesQuery{
		Description: &description,
		UUID:        &uuid,
	})
	assertEqual(t, err, nil)
	if res == nil {
		t.Fatal("expected pipeline list")
	}
	assertEqual(t, len(res.GetObjects()), 1)
}

func TestPipelinesGet(t *testing.T) {
	const pipelineID = "a64e061a-5688-49b5-95c1-0b6885c40c04"

	client, err := ssclient.New(ssclient.Config{
		BaseURL:    "http://storage.service",
		Username:   "test",
		Key:        "test",
		HTTPClient: &http.Client{},
	})
	assertEqual(t, err, nil)

	raw := client.Raw()
	raw.RequestAdapter = &fakeRequestAdapter{
		baseURL: "http://storage.service",
		send: func(ctx context.Context, requestInfo *kabs.RequestInformation, constructor serialization.ParsableFactory, errorMappings kabs.ErrorMappings) (serialization.Parsable, error) {
			assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/pipeline/{uuid}")
			assertEqual(t, requestInfo.PathParameters["uuid"], pipelineID)

			pipeline := models.NewPipeline()
			pipeline.SetUuid(ptr(pipelineID))
			return pipeline, nil
		},
	}

	res, err := client.Pipelines().Get(context.Background(), pipelineID)
	assertEqual(t, err, nil)
	if res == nil {
		t.Fatal("expected pipeline")
	}
	assertEqual(t, *res.GetUuid(), pipelineID)
}

func ptr[T any](value T) *T {
	return &value
}
