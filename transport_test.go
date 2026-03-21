package ssclient

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	kabs "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"
	khttp "github.com/microsoft/kiota-http-go"
	"go.artefactual.dev/ssclient/kiota/models"
)

type executeTestAdapter struct {
	convert func(context.Context, *kabs.RequestInformation) (any, error)
}

func (a *executeTestAdapter) Send(context.Context, *kabs.RequestInformation, serialization.ParsableFactory, kabs.ErrorMappings) (serialization.Parsable, error) {
	panic("unexpected Send call")
}

func (a *executeTestAdapter) SendEnum(context.Context, *kabs.RequestInformation, serialization.EnumFactory, kabs.ErrorMappings) (any, error) {
	panic("unexpected SendEnum call")
}

func (a *executeTestAdapter) SendCollection(context.Context, *kabs.RequestInformation, serialization.ParsableFactory, kabs.ErrorMappings) ([]serialization.Parsable, error) {
	panic("unexpected SendCollection call")
}

func (a *executeTestAdapter) SendEnumCollection(context.Context, *kabs.RequestInformation, serialization.EnumFactory, kabs.ErrorMappings) ([]any, error) {
	panic("unexpected SendEnumCollection call")
}

func (a *executeTestAdapter) SendPrimitive(context.Context, *kabs.RequestInformation, string, kabs.ErrorMappings) (any, error) {
	panic("unexpected SendPrimitive call")
}

func (a *executeTestAdapter) SendPrimitiveCollection(context.Context, *kabs.RequestInformation, string, kabs.ErrorMappings) ([]any, error) {
	panic("unexpected SendPrimitiveCollection call")
}

func (a *executeTestAdapter) SendNoContent(context.Context, *kabs.RequestInformation, kabs.ErrorMappings) error {
	panic("unexpected SendNoContent call")
}

func (a *executeTestAdapter) GetSerializationWriterFactory() serialization.SerializationWriterFactory {
	return nil
}

func (a *executeTestAdapter) EnableBackingStore(store.BackingStoreFactory) {}

func (a *executeTestAdapter) SetBaseUrl(string) {}

func (a *executeTestAdapter) GetBaseUrl() string {
	return ""
}

func (a *executeTestAdapter) ConvertToNativeRequest(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
	return a.convert(ctx, requestInfo)
}

type executeRoundTripFunc func(*http.Request) (*http.Response, error)

func (f executeRoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

type executeMiddleware struct{}

func (m *executeMiddleware) Intercept(pipeline khttp.Pipeline, middlewareIndex int, req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Execute-Middleware", "true")
	return pipeline.Next(req, middlewareIndex)
}

func executePtr[T any](value T) *T {
	return &value
}

func TestExecuteCapturesStatusHeadersAndBody(t *testing.T) {
	client, err := New(Config{
		BaseURL:  "http://storage.service",
		Username: "test",
		Key:      "test",
		HTTPClient: &http.Client{
			Transport: executeRoundTripFunc(func(r *http.Request) (*http.Response, error) {
				if got, want := r.URL.String(), "http://storage.service/api/v2/file/delete"; got != want {
					t.Fatalf("unexpected URL %q want %q", got, want)
				}
				return &http.Response{
					StatusCode: http.StatusAccepted,
					Header:     http.Header{"Location": {"/jobs/123"}},
					Body:       io.NopCloser(strings.NewReader(`{"job_id":"123","state":"queued"}`)),
					Request:    r,
				}, nil
			}),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	client.adapter = &executeTestAdapter{
		convert: func(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
			if got, want := requestInfo.UrlTemplate, "{+baseurl}/api/v2/file/delete"; got != want {
				t.Fatalf("unexpected template %q want %q", got, want)
			}
			return &http.Request{
				Method: http.MethodPost,
				URL: &url.URL{
					Scheme: "http",
					Host:   "storage.service",
					Path:   "/api/v2/file/delete",
				},
				Header: http.Header{
					"Authorization": {"ApiKey test:test"},
				},
			}, nil
		},
	}

	reqInfo := kabs.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(kabs.POST, "{+baseurl}/api/v2/file/delete", map[string]string{
		"baseurl": "http://storage.service",
	})

	resp, err := client.execute(context.Background(), reqInfo)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := resp.StatusCode, http.StatusAccepted; got != want {
		t.Fatalf("unexpected status %d want %d", got, want)
	}
	if got, want := resp.Headers.Get("Location"), "/jobs/123"; got != want {
		t.Fatalf("unexpected Location %q want %q", got, want)
	}
	if got, want := string(resp.Body), `{"job_id":"123","state":"queued"}`; got != want {
		t.Fatalf("unexpected body %q want %q", got, want)
	}
}

func TestExecuteUsesConfiguredMiddlewares(t *testing.T) {
	client, err := New(Config{
		BaseURL:     "http://storage.service",
		Username:    "test",
		Key:         "test",
		Middlewares: []khttp.Middleware{&executeMiddleware{}},
		HTTPClient: &http.Client{
			Transport: executeRoundTripFunc(func(r *http.Request) (*http.Response, error) {
				if got, want := r.Header.Get("X-Execute-Middleware"), "true"; got != want {
					t.Fatalf("unexpected middleware header %q want %q", got, want)
				}
				return &http.Response{
					StatusCode: http.StatusNoContent,
					Header:     http.Header{},
					Body:       io.NopCloser(strings.NewReader("")),
					Request:    r,
				}, nil
			}),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	client.adapter = &executeTestAdapter{
		convert: func(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
			return &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "http",
					Host:   "storage.service",
					Path:   "/api/v2/location/default/DS/",
				},
				Header: http.Header{
					"Accept": {"application/json"},
				},
			}, nil
		},
	}

	reqInfo := kabs.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(kabs.GET, "{+baseurl}/api/v2/location/default/DS/", map[string]string{
		"baseurl": "http://storage.service",
	})

	if _, err := client.execute(context.Background(), reqInfo); err != nil {
		t.Fatal(err)
	}
}

func TestExecuteHonorsHTTPClientTimeout(t *testing.T) {
	client, err := New(Config{
		BaseURL:  "http://storage.service",
		Username: "test",
		Key:      "test",
		HTTPClient: &http.Client{
			Timeout: 20 * time.Millisecond,
			Transport: executeRoundTripFunc(func(r *http.Request) (*http.Response, error) {
				<-r.Context().Done()
				return nil, r.Context().Err()
			}),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	client.adapter = &executeTestAdapter{
		convert: func(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
			return &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Scheme: "http",
					Host:   "storage.service",
					Path:   "/api/v2/location/default/DS/",
				},
				Header: http.Header{
					"Accept": {"application/json"},
				},
			}, nil
		},
	}

	reqInfo := kabs.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(kabs.GET, "{+baseurl}/api/v2/location/default/DS/", map[string]string{
		"baseurl": "http://storage.service",
	})

	_, err = client.execute(context.Background(), reqInfo)
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded error, got %v", err)
	}
}

func TestExecuteSeedsBaseURLForRealAdapter(t *testing.T) {
	const baseURL = "http://storage.service"

	client, err := New(Config{
		BaseURL:  baseURL,
		Username: "test",
		Key:      "test",
		HTTPClient: &http.Client{
			Transport: executeRoundTripFunc(func(r *http.Request) (*http.Response, error) {
				if got, want := r.URL.String(), baseURL+"/api/v2/location/default/DS/"; got != want {
					t.Fatalf("unexpected URL %q want %q", got, want)
				}
				return &http.Response{
					StatusCode: http.StatusFound,
					Header:     http.Header{"Location": {baseURL + "/api/v2/location/123"}},
					Body:       io.NopCloser(strings.NewReader("")),
					Request:    r,
				}, nil
			}),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	reqInfo := kabs.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(kabs.GET, "{+baseurl}/api/v2/location/default/DS/", map[string]string{})

	resp, err := client.execute(context.Background(), reqInfo)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := resp.StatusCode, http.StatusFound; got != want {
		t.Fatalf("unexpected status %d want %d", got, want)
	}
	if got, want := reqInfo.PathParameters["baseurl"], baseURL; got != want {
		t.Fatalf("unexpected baseurl %q want %q", got, want)
	}
}

func TestExecuteStreamWithoutBaseURLFailsWithRealAdapterBeforeSeeding(t *testing.T) {
	const baseURL = "http://storage.service"

	client, err := New(Config{
		BaseURL:    baseURL,
		Username:   "test",
		Key:        "test",
		HTTPClient: &http.Client{},
	})
	if err != nil {
		t.Fatal(err)
	}

	reqInfo := kabs.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(kabs.GET, "{+baseurl}/api/v2/file/pkg/pointer_file/", map[string]string{})
	ensureRequestBaseURL(client.Adapter(), reqInfo)
	if _, err := client.Adapter().ConvertToNativeRequest(context.Background(), reqInfo); err != nil {
		t.Fatalf("expected baseurl-seeded native request, got %v", err)
	}
}

func TestEnsureRequestBaseURLNoopWithoutTemplate(t *testing.T) {
	reqInfo := kabs.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(kabs.GET, "http://storage.service/static", nil)
	ensureRequestBaseURL(&executeTestAdapter{}, reqInfo)
	if _, ok := reqInfo.PathParameters["baseurl"]; ok {
		t.Fatal("did not expect baseurl to be injected")
	}
}

func TestPackagesDeleteAndReviewWorkWithRealAdapter(t *testing.T) {
	const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"
	const baseURL = "http://storage.service"

	t.Run("DeleteAIP", func(t *testing.T) {
		client, err := New(Config{
			BaseURL:  baseURL,
			Username: "test",
			Key:      "test",
			HTTPClient: &http.Client{
				Transport: executeRoundTripFunc(func(r *http.Request) (*http.Response, error) {
					if got, want := r.URL.String(), baseURL+"/api/v2/file/"+packageID+"/delete_aip/"; got != want {
						t.Fatalf("unexpected URL %q want %q", got, want)
					}
					return &http.Response{
						StatusCode: http.StatusAccepted,
						Header:     http.Header{"Content-Type": {"application/json"}},
						Body:       io.NopCloser(strings.NewReader(`{"message":"Delete request created successfully.","id":17}`)),
						Request:    r,
					}, nil
				}),
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		body := models.NewDeleteAipRequest()
		body.SetEventReason(executePtr("Delete please!"))
		body.SetPipeline(executePtr(uuid.MustParse("4b9e8af5-b0af-4abf-80b8-4b7d76281f61")))
		body.SetUserId(executePtr(int32(1)))
		body.SetUserEmail(executePtr("user@example.com"))

		res, err := client.Packages().DeleteAIP(context.Background(), uuid.MustParse(packageID), body)
		if err != nil {
			t.Fatal(err)
		}
		if res == nil || !res.IsAccepted() {
			t.Fatalf("unexpected delete result %#v", res)
		}
	})

	t.Run("ReviewAIPDeletion", func(t *testing.T) {
		client, err := New(Config{
			BaseURL:  baseURL,
			Username: "test",
			Key:      "test",
			HTTPClient: &http.Client{
				Transport: executeRoundTripFunc(func(r *http.Request) (*http.Response, error) {
					if got, want := r.URL.String(), baseURL+"/api/v2/file/"+packageID+"/review_aip_deletion/"; got != want {
						t.Fatalf("unexpected URL %q want %q", got, want)
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     http.Header{"Content-Type": {"application/json"}},
						Body:       io.NopCloser(strings.NewReader(`{"message":"done"}`)),
						Request:    r,
					}, nil
				}),
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		body := models.NewReviewAipDeletionRequest()
		body.SetEventId(executePtr(int32(99)))
		decision := models.APPROVE_REVIEWAIPDELETIONDECISION
		body.SetDecision(&decision)
		body.SetReason(executePtr("approved by workflow"))

		res, err := client.Packages().ReviewAIPDeletion(context.Background(), uuid.MustParse(packageID), body)
		if err != nil {
			t.Fatal(err)
		}
		_ = res
	})
}

func TestResponseSnapshotDecodeJSON(t *testing.T) {
	resp := &responseSnapshot{
		StatusCode: http.StatusOK,
		Body:       []byte(`{"kind":"deleted","uuid":"1234"}`),
	}

	var payload struct {
		Kind string `json:"kind"`
		ID   string `json:"uuid"`
	}
	if err := resp.decodeJSON(&payload); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(payload, struct {
		Kind string `json:"kind"`
		ID   string `json:"uuid"`
	}{Kind: "deleted", ID: "1234"}) {
		t.Fatalf("unexpected payload %#v", payload)
	}
}
