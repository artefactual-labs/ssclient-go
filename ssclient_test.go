package ssclient_test

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	kabs "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"

	"go.artefactual.dev/ssclient"
	"go.artefactual.dev/ssclient/kiota/models"
)

// fakeRequestAdapter is a narrow Kiota adapter test double that lets each test
// override only the request-adapter entry points it exercises.
type fakeRequestAdapter struct {
	send                       func(context.Context, *kabs.RequestInformation, serialization.ParsableFactory, kabs.ErrorMappings) (serialization.Parsable, error)
	sendPrimitive              func(context.Context, *kabs.RequestInformation, string, kabs.ErrorMappings) (any, error)
	convertToNativeRequest     func(context.Context, *kabs.RequestInformation) (any, error)
	baseURL                    string
	serializationWriterFactory serialization.SerializationWriterFactory
}

// roundTripFunc adapts a function to http.RoundTripper for transport-level
// request assertions in tests.
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

func assertEqual(t interface {
	Helper()
	Errorf(string, ...any)
}, got, want any,
) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Mismatch found:\nGot: %#v\nWant: %#v", got, want)
	}
}

func ptr[T any](value T) *T {
	return &value
}

func TestClient(t *testing.T) {
	t.Run("New requires BaseURL", func(t *testing.T) {
		client, err := ssclient.New(ssclient.Config{})
		if err == nil {
			t.Fatal("expected error")
		}
		if got, want := err.Error(), "base URL is required"; got != want {
			t.Fatalf("unexpected error %q want %q", got, want)
		}
		if client != nil {
			t.Fatal("expected nil client")
		}
	})

	t.Run("Raw", func(t *testing.T) {
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
		if client.Async() == nil {
			t.Fatal("expected async service")
		}
		if client.Packages() == nil {
			t.Fatal("expected packages service")
		}
	})

	t.Run("Trailing slash URLs", func(t *testing.T) {
		const baseURL = "http://storage.service"

		t.Run("Wrapper methods", func(t *testing.T) {
			t.Run("PackagesGet", func(t *testing.T) {
				const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"

				client, err := ssclient.New(ssclient.Config{
					BaseURL:  baseURL,
					Username: "test",
					Key:      "test",
					HTTPClient: &http.Client{
						Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
							assertEqual(t, r.URL.String(), baseURL+"/api/v2/file/"+packageID+"/")
							return &http.Response{
								StatusCode: http.StatusOK,
								Header:     http.Header{"Content-Type": {"application/json"}},
								Body:       io.NopCloser(strings.NewReader(`{"uuid":"` + packageID + `","status":"UPLOADED"}`)),
								Request:    r,
							}, nil
						}),
					},
				})
				assertEqual(t, err, nil)

				_, err = client.Packages().Get(context.Background(), uuid.MustParse(packageID))
				assertEqual(t, err, nil)
			})

			t.Run("LocationsGet", func(t *testing.T) {
				const locationID = "fff70864-a5d4-4ca6-ab29-b4ce67d8eeab"

				client, err := ssclient.New(ssclient.Config{
					BaseURL:  baseURL,
					Username: "test",
					Key:      "test",
					HTTPClient: &http.Client{
						Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
							assertEqual(t, r.URL.String(), baseURL+"/api/v2/location/"+locationID+"/")
							return &http.Response{
								StatusCode: http.StatusOK,
								Header:     http.Header{"Content-Type": {"application/json"}},
								Body:       io.NopCloser(strings.NewReader(`{"uuid":"` + locationID + `"}`)),
								Request:    r,
							}, nil
						}),
					},
				})
				assertEqual(t, err, nil)

				_, err = client.Locations().Get(context.Background(), uuid.MustParse(locationID))
				assertEqual(t, err, nil)
			})

			t.Run("LocationsMove", func(t *testing.T) {
				const locationID = "fff70864-a5d4-4ca6-ab29-b4ce67d8eeab"

				client, err := ssclient.New(ssclient.Config{
					BaseURL:  baseURL,
					Username: "test",
					Key:      "test",
					HTTPClient: &http.Client{
						Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
							assertEqual(t, r.URL.String(), baseURL+"/api/v2/location/"+locationID+"/")
							return &http.Response{
								StatusCode: http.StatusOK,
								Header:     http.Header{"Content-Type": {"application/json"}},
								Body:       io.NopCloser(strings.NewReader("")),
								Request:    r,
							}, nil
						}),
					},
				})
				assertEqual(t, err, nil)

				body := models.NewMoveRequest()
				body.SetOriginLocation(ptr("/api/v2/location/origin/"))
				body.SetPipeline(ptr("/api/v2/pipeline/source/"))
				assertEqual(t, client.Locations().Move(context.Background(), uuid.MustParse(locationID), body), nil)
			})

			t.Run("PipelinesGet", func(t *testing.T) {
				const pipelineID = "a64e061a-5688-49b5-95c1-0b6885c40c04"

				client, err := ssclient.New(ssclient.Config{
					BaseURL:  baseURL,
					Username: "test",
					Key:      "test",
					HTTPClient: &http.Client{
						Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
							assertEqual(t, r.URL.String(), baseURL+"/api/v2/pipeline/"+pipelineID+"/")
							return &http.Response{
								StatusCode: http.StatusOK,
								Header:     http.Header{"Content-Type": {"application/json"}},
								Body:       io.NopCloser(strings.NewReader(`{"uuid":"` + pipelineID + `"}`)),
								Request:    r,
							}, nil
						}),
					},
				})
				assertEqual(t, err, nil)

				_, err = client.Pipelines().Get(context.Background(), uuid.MustParse(pipelineID))
				assertEqual(t, err, nil)
			})
		})

		t.Run("Generated item empty path segment builders", func(t *testing.T) {
			t.Run("PackageItem", func(t *testing.T) {
				const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"

				client, err := ssclient.New(ssclient.Config{
					BaseURL:  baseURL,
					Username: "test",
					Key:      "test",
					HTTPClient: &http.Client{
						Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
							assertEqual(t, r.URL.String(), baseURL+"/api/v2/file/"+packageID+"/")
							return &http.Response{
								StatusCode: http.StatusOK,
								Header:     http.Header{"Content-Type": {"application/json"}},
								Body:       io.NopCloser(strings.NewReader(`{"uuid":"` + packageID + `","status":"UPLOADED"}`)),
								Request:    r,
							}, nil
						}),
					},
				})
				assertEqual(t, err, nil)

				_, err = client.Raw().Api().V2().File().ByUuid(uuid.MustParse(packageID)).EmptyPathSegment().Get(context.Background(), nil)
				assertEqual(t, err, nil)
			})

			t.Run("LocationItem", func(t *testing.T) {
				const locationID = "fff70864-a5d4-4ca6-ab29-b4ce67d8eeab"

				client, err := ssclient.New(ssclient.Config{
					BaseURL:  baseURL,
					Username: "test",
					Key:      "test",
					HTTPClient: &http.Client{
						Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
							assertEqual(t, r.URL.String(), baseURL+"/api/v2/location/"+locationID+"/")
							return &http.Response{
								StatusCode: http.StatusOK,
								Header:     http.Header{"Content-Type": {"application/json"}},
								Body:       io.NopCloser(strings.NewReader("")),
								Request:    r,
							}, nil
						}),
					},
				})
				assertEqual(t, err, nil)

				body := models.NewMoveRequest()
				body.SetOriginLocation(ptr("/api/v2/location/origin/"))
				body.SetPipeline(ptr("/api/v2/pipeline/source/"))
				_, err = client.Raw().Api().V2().Location().ByUuid(uuid.MustParse(locationID)).EmptyPathSegment().Post(context.Background(), body, nil)
				assertEqual(t, err, nil)
			})

			t.Run("PipelineItem", func(t *testing.T) {
				const pipelineID = "a64e061a-5688-49b5-95c1-0b6885c40c04"

				client, err := ssclient.New(ssclient.Config{
					BaseURL:  baseURL,
					Username: "test",
					Key:      "test",
					HTTPClient: &http.Client{
						Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
							assertEqual(t, r.URL.String(), baseURL+"/api/v2/pipeline/"+pipelineID+"/")
							return &http.Response{
								StatusCode: http.StatusOK,
								Header:     http.Header{"Content-Type": {"application/json"}},
								Body:       io.NopCloser(strings.NewReader(`{"uuid":"` + pipelineID + `"}`)),
								Request:    r,
							}, nil
						}),
					},
				})
				assertEqual(t, err, nil)

				_, err = client.Raw().Api().V2().Pipeline().ByUuid(uuid.MustParse(pipelineID)).EmptyPathSegment().Get(context.Background(), nil)
				assertEqual(t, err, nil)
			})
		})
	})
}
