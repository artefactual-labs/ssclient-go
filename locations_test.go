package ssclient_test

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	kabs "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"

	"go.artefactual.dev/ssclient"
	"go.artefactual.dev/ssclient/kiota/models"
)

func TestLocations(t *testing.T) {
	t.Run("List", func(t *testing.T) {
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
				location.SetUuid(ptr(uuid.MustParse("fff70864-a5d4-4ca6-ab29-b4ce67d8eeab")))
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
	})

	t.Run("Get", func(t *testing.T) {
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
				assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/location/{uuid}/")
				assertEqual(t, requestInfo.PathParameters["uuid"], locationID)

				location := models.NewLocation()
				location.SetUuid(ptr(uuid.MustParse(locationID)))
				return location, nil
			},
		}

		res, err := client.Locations().Get(context.Background(), uuid.MustParse(locationID))
		assertEqual(t, err, nil)
		if res == nil {
			t.Fatal("expected location")
		}
		assertEqual(t, res.GetUuid().String(), locationID)
	})

	t.Run("Default", func(t *testing.T) {
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

		client.Raw().RequestAdapter = &fakeRequestAdapter{
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

		location, err := client.Locations().Default(context.Background(), models.DS_LOCATIONPURPOSE)
		assertEqual(t, err, nil)
		assertEqual(t, location, "/api/v2/location/154660b9-b4a3-4886-8d68-5e170c0923b8/")
	})

	t.Run("Move", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
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
					assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/location/{uuid}/")
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
			assertEqual(t, client.Locations().Move(context.Background(), uuid.MustParse(locationID), body), nil)
		})

		t.Run("NilBody", func(t *testing.T) {
			client, err := ssclient.New(ssclient.Config{
				BaseURL:    "http://storage.service",
				Username:   "test",
				Key:        "test",
				HTTPClient: &http.Client{},
			})
			assertEqual(t, err, nil)

			if err := client.Locations().Move(context.Background(), uuid.Nil, nil); err == nil {
				t.Fatal("expected error")
			}
		})
	})
}
