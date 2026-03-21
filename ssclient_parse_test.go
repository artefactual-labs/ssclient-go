package ssclient_test

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	kabs "github.com/microsoft/kiota-abstractions-go"
	khttp "github.com/microsoft/kiota-http-go"

	"go.artefactual.dev/ssclient"
)

func TestKiotaParsesTolerantTimestampsAsUTC(t *testing.T) {
	t.Run("FixityTimestamp", func(t *testing.T) {
		testCases := []struct {
			name      string
			timestamp string
			want      time.Time
		}{
			{
				name:      "RFC3339",
				timestamp: "2026-03-19T05:34:31.352112Z",
				want:      time.Date(2026, 3, 19, 5, 34, 31, 352112000, time.UTC),
			},
			{
				name:      "NaiveWithFraction",
				timestamp: "2026-03-19T05:34:31.352112",
				want:      time.Date(2026, 3, 19, 5, 34, 31, 352112000, time.UTC),
			},
			{
				name:      "NaiveWithoutFraction",
				timestamp: "2026-03-19T05:34:31",
				want:      time.Date(2026, 3, 19, 5, 34, 31, 0, time.UTC),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				client, err := ssclient.New(ssclient.Config{
					BaseURL:  "http://storage.service",
					Username: "test",
					Key:      "test",
					HTTPClient: &http.Client{
						Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
							body := `{"success":true,"message":"","timestamp":"` + tc.timestamp + `"}`
							return &http.Response{
								StatusCode:    http.StatusOK,
								Header:        http.Header{"Content-Type": {"application/json"}},
								ContentLength: int64(len(body)),
								Body:          ioNopCloserString(body),
								Request:       r,
							}, nil
						}),
					},
				})
				if err != nil {
					t.Fatal(err)
				}

				res, err := client.Packages().CheckFixity(context.Background(), uuid.MustParse("7c8a3549-2fe0-41d3-9d83-f485f1a43be3"), ssclient.CheckFixityOptions{})
				if err != nil {
					t.Fatal(err)
				}
				if res == nil || res.GetTimestamp() == nil {
					t.Fatal("expected parsed timestamp")
				}
				if got := *res.GetTimestamp(); !got.Equal(tc.want) || got.Location() != time.UTC {
					t.Fatalf("unexpected timestamp %#v want %#v", got, tc.want)
				}
			})
		}
	})

	t.Run("PackageStoredDate", func(t *testing.T) {
		client, err := ssclient.New(ssclient.Config{
			BaseURL:  "http://storage.service",
			Username: "test",
			Key:      "test",
			HTTPClient: &http.Client{
				Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					if !strings.HasSuffix(r.URL.Path, "/api/v2/file/7c8a3549-2fe0-41d3-9d83-f485f1a43be3/") {
						t.Fatalf("unexpected path %q", r.URL.Path)
					}
					body := `{"uuid":"7c8a3549-2fe0-41d3-9d83-f485f1a43be3","status":"UPLOADED","stored_date":"2026-03-19T05:34:31.352112"}`
					return &http.Response{
						StatusCode:    http.StatusOK,
						Header:        http.Header{"Content-Type": {"application/json"}},
						ContentLength: int64(len(body)),
						Body:          ioNopCloserString(body),
						Request:       r,
					}, nil
				}),
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		pkg, err := client.Packages().Get(context.Background(), uuid.MustParse("7c8a3549-2fe0-41d3-9d83-f485f1a43be3"))
		if err != nil {
			t.Fatal(err)
		}
		if pkg == nil || pkg.GetStoredDate() == nil {
			t.Fatal("expected parsed stored_date")
		}

		want := time.Date(2026, 3, 19, 5, 34, 31, 352112000, time.UTC)
		if got := *pkg.GetStoredDate(); !got.Equal(want) || got.Location() != time.UTC {
			t.Fatalf("unexpected stored_date %#v want %#v", got, want)
		}
	})
}

func ioNopCloserString(s string) *readCloserString {
	return &readCloserString{Reader: strings.NewReader(s)}
}

type readCloserString struct {
	*strings.Reader
}

func (r *readCloserString) Close() error { return nil }

func unwrapAdditionalDataValue(value any) any {
	switch v := value.(type) {
	case *string:
		if v == nil {
			return nil
		}
		return *v
	case *bool:
		if v == nil {
			return nil
		}
		return *v
	case *int32:
		if v == nil {
			return nil
		}
		return *v
	case *int64:
		if v == nil {
			return nil
		}
		return *v
	case *float32:
		if v == nil {
			return nil
		}
		return *v
	case *float64:
		if v == nil {
			return nil
		}
		return *v
	default:
		return value
	}
}

func newClientForJSONResponse(t *testing.T, body string) (*ssclient.Client, *http.Client) {
	t.Helper()

	httpClient := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode:    http.StatusOK,
				Header:        http.Header{"Content-Type": {"application/json"}},
				ContentLength: int64(len(body)),
				Body:          ioNopCloserString(body),
				Request:       r,
			}, nil
		}),
	}

	client, err := ssclient.New(ssclient.Config{
		BaseURL:    "http://storage.service",
		Username:   "test",
		Key:        "test",
		HTTPClient: httpClient,
	})
	if err != nil {
		t.Fatal(err)
	}

	return client, httpClient
}

func replaceWithStockKiotaAdapter(t *testing.T, client *ssclient.Client, httpClient *http.Client) {
	t.Helper()

	adapter, err := khttp.NewNetHttpRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		noopAuthProvider{},
		nil,
		nil,
		httpClient,
	)
	if err != nil {
		t.Fatal(err)
	}
	adapter.SetBaseUrl("http://storage.service")
	client.Raw().RequestAdapter = adapter
}

type noopAuthProvider struct{}

func (noopAuthProvider) AuthenticateRequest(context.Context, *kabs.RequestInformation, map[string]any) error {
	return nil
}

func TestKiotaPreservesScalarFieldBehavior(t *testing.T) {
	t.Run("PackageGetMatchesStockKiotaForNonTimeFields", func(t *testing.T) {
		body := `{
			"uuid":"7c8a3549-2fe0-41d3-9d83-f485f1a43be3",
			"status":"UPLOADED",
			"size":42,
			"current_path":"bags/aip.7z",
			"current_full_path":"/var/archivematica/sharedDirectory/www/AIPsStore/bags/aip.7z",
			"current_location":"/api/v2/location/154660b9-b4a3-4886-8d68-5e170c0923b8/",
			"encrypted":false,
			"origin_pipeline":"/api/v2/pipeline/154660b9-b4a3-4886-8d68-5e170c0923b8/",
			"package_type":"AIP",
			"related_packages":["/api/v2/file/a/","/api/v2/file/b/"],
			"replicas":["/api/v2/file/c/"],
			"replicated_package":"/api/v2/file/d/",
			"resource_uri":"/api/v2/file/7c8a3549-2fe0-41d3-9d83-f485f1a43be3/",
			"stored_date":"2026-03-19T05:34:31.352112",
			"unknown_field":"kept"
		}`

		tolerantClient, _ := newClientForJSONResponse(t, body)
		tolerantPkg, err := tolerantClient.Packages().Get(context.Background(), uuid.MustParse("7c8a3549-2fe0-41d3-9d83-f485f1a43be3"))
		if err != nil {
			t.Fatal(err)
		}

		stockClient, stockHTTPClient := newClientForJSONResponse(t, body)
		replaceWithStockKiotaAdapter(t, stockClient, stockHTTPClient)
		_, err = stockClient.Packages().Get(context.Background(), uuid.MustParse("7c8a3549-2fe0-41d3-9d83-f485f1a43be3"))
		if err == nil {
			t.Fatal("expected stock Kiota parsing to fail on naive timestamp")
		}

		if tolerantPkg == nil {
			t.Fatal("expected package")
		}
		if tolerantPkg.GetStoredDate() == nil {
			t.Fatal("expected stored_date")
		}

		if got, want := *tolerantPkg.GetStatus(), "UPLOADED"; got != want {
			t.Fatalf("status=%q want %q", got, want)
		}
		if got, want := *tolerantPkg.GetCurrentPath(), "bags/aip.7z"; got != want {
			t.Fatalf("current_path=%q want %q", got, want)
		}
		if got, want := *tolerantPkg.GetCurrentFullPath(), "/var/archivematica/sharedDirectory/www/AIPsStore/bags/aip.7z"; got != want {
			t.Fatalf("current_full_path=%q want %q", got, want)
		}
		if got, want := *tolerantPkg.GetCurrentLocation(), "/api/v2/location/154660b9-b4a3-4886-8d68-5e170c0923b8/"; got != want {
			t.Fatalf("current_location=%q want %q", got, want)
		}
		if got, want := *tolerantPkg.GetEncrypted(), false; got != want {
			t.Fatalf("encrypted=%v want %v", got, want)
		}
		if got, want := *tolerantPkg.GetOriginPipeline(), "/api/v2/pipeline/154660b9-b4a3-4886-8d68-5e170c0923b8/"; got != want {
			t.Fatalf("origin_pipeline=%q want %q", got, want)
		}
		if got, want := int(*tolerantPkg.GetSize()), 42; got != want {
			t.Fatalf("size=%d want %d", got, want)
		}
		if got, want := tolerantPkg.GetRelatedPackages(), []string{"/api/v2/file/a/", "/api/v2/file/b/"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("related_packages=%#v want %#v", got, want)
		}
		if got, want := tolerantPkg.GetReplicas(), []string{"/api/v2/file/c/"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("replicas=%#v want %#v", got, want)
		}
		if got, want := unwrapAdditionalDataValue(tolerantPkg.GetAdditionalData()["unknown_field"]), "kept"; got != want {
			t.Fatalf("additionalData[unknown_field]=%#v want %#v", got, want)
		}
	})
}

func TestKiotaPreservesCollectionBehavior(t *testing.T) {
	t.Run("LocationsListMatchesStockKiotaForNonTimeFields", func(t *testing.T) {
		body := `{
			"meta":{
				"limit":20,
				"next":"/api/v2/location/?limit=20&offset=20",
				"offset":0,
				"previous":null,
				"total_count":2,
				"meta_extra":"ok"
			},
			"objects":[
				{
					"uuid":"154660b9-b4a3-4886-8d68-5e170c0923b8",
					"description":"AIP store",
					"enabled":true,
					"path":"/var/aips",
					"pipeline":["/api/v2/pipeline/1/"],
					"purpose":"AS",
					"quota":100,
					"relative_path":"",
					"resource_uri":"/api/v2/location/154660b9-b4a3-4886-8d68-5e170c0923b8/",
					"space":"/api/v2/space/1/",
					"used":10,
					"object_extra":"one"
				},
				{
					"uuid":"254660b9-b4a3-4886-8d68-5e170c0923b8",
					"description":"Transfer store",
					"enabled":false,
					"path":"/var/transfers",
					"pipeline":["/api/v2/pipeline/2/","/api/v2/pipeline/3/"],
					"purpose":"TS",
					"quota":200,
					"relative_path":"incoming",
					"resource_uri":"/api/v2/location/254660b9-b4a3-4886-8d68-5e170c0923b8/",
					"space":"/api/v2/space/2/",
					"used":20
				}
			]
		}`

		tolerantClient, _ := newClientForJSONResponse(t, body)
		tolerantList, err := tolerantClient.Locations().List(context.Background(), ssclient.ListLocationsQuery{})
		if err != nil {
			t.Fatal(err)
		}

		stockClient, stockHTTPClient := newClientForJSONResponse(t, body)
		replaceWithStockKiotaAdapter(t, stockClient, stockHTTPClient)
		stockList, err := stockClient.Locations().List(context.Background(), ssclient.ListLocationsQuery{})
		if err != nil {
			t.Fatal(err)
		}

		if tolerantList == nil || stockList == nil {
			t.Fatal("expected location lists")
		}

		if got, want := *tolerantList.GetMeta().GetLimit(), *stockList.GetMeta().GetLimit(); got != want {
			t.Fatalf("meta.limit=%d want %d", got, want)
		}
		if got, want := *tolerantList.GetMeta().GetTotalCount(), *stockList.GetMeta().GetTotalCount(); got != want {
			t.Fatalf("meta.total_count=%d want %d", got, want)
		}
		if got, want := unwrapAdditionalDataValue(tolerantList.GetMeta().GetAdditionalData()["meta_extra"]), unwrapAdditionalDataValue(stockList.GetMeta().GetAdditionalData()["meta_extra"]); !reflect.DeepEqual(got, want) {
			t.Fatalf("meta additional data=%#v want %#v", got, want)
		}

		if got, want := len(tolerantList.GetObjects()), len(stockList.GetObjects()); got != want {
			t.Fatalf("objects length=%d want %d", got, want)
		}
		for i := range tolerantList.GetObjects() {
			gotObj := tolerantList.GetObjects()[i]
			wantObj := stockList.GetObjects()[i]
			if got, want := *gotObj.GetDescription(), *wantObj.GetDescription(); got != want {
				t.Fatalf("objects[%d].description=%q want %q", i, got, want)
			}
			if got, want := *gotObj.GetEnabled(), *wantObj.GetEnabled(); got != want {
				t.Fatalf("objects[%d].enabled=%v want %v", i, got, want)
			}
			if got, want := *gotObj.GetPath(), *wantObj.GetPath(); got != want {
				t.Fatalf("objects[%d].path=%q want %q", i, got, want)
			}
			if got, want := gotObj.GetPipeline(), wantObj.GetPipeline(); !reflect.DeepEqual(got, want) {
				t.Fatalf("objects[%d].pipeline=%#v want %#v", i, got, want)
			}
			if got, want := *gotObj.GetQuota(), *wantObj.GetQuota(); got != want {
				t.Fatalf("objects[%d].quota=%d want %d", i, got, want)
			}
			if got, want := *gotObj.GetUsed(), *wantObj.GetUsed(); got != want {
				t.Fatalf("objects[%d].used=%d want %d", i, got, want)
			}
		}

		if got, want := unwrapAdditionalDataValue(tolerantList.GetObjects()[0].GetAdditionalData()["object_extra"]), unwrapAdditionalDataValue(stockList.GetObjects()[0].GetAdditionalData()["object_extra"]); !reflect.DeepEqual(got, want) {
			t.Fatalf("object additional data=%#v want %#v", got, want)
		}
	})
}

func TestKiotaPreservesEnumBehavior(t *testing.T) {
	t.Run("LocationPurposeMatchesStockKiota", func(t *testing.T) {
		body := `{
			"uuid":"154660b9-b4a3-4886-8d68-5e170c0923b8",
			"description":"AIP store",
			"enabled":true,
			"path":"/var/aips",
			"pipeline":["/api/v2/pipeline/1/"],
			"purpose":"AS",
			"quota":100,
			"relative_path":"",
			"resource_uri":"/api/v2/location/154660b9-b4a3-4886-8d68-5e170c0923b8/",
			"space":"/api/v2/space/1/",
			"used":10
		}`

		tolerantClient, _ := newClientForJSONResponse(t, body)
		tolerantLocation, err := tolerantClient.Locations().Get(context.Background(), uuid.MustParse("154660b9-b4a3-4886-8d68-5e170c0923b8"))
		if err != nil {
			t.Fatal(err)
		}

		stockClient, stockHTTPClient := newClientForJSONResponse(t, body)
		replaceWithStockKiotaAdapter(t, stockClient, stockHTTPClient)
		stockLocation, err := stockClient.Locations().Get(context.Background(), uuid.MustParse("154660b9-b4a3-4886-8d68-5e170c0923b8"))
		if err != nil {
			t.Fatal(err)
		}

		if tolerantLocation == nil || stockLocation == nil {
			t.Fatal("expected locations")
		}
		if got, want := tolerantLocation.GetPurpose().String(), stockLocation.GetPurpose().String(); got != want {
			t.Fatalf("purpose=%q want %q", got, want)
		}
	})
}

func TestKiotaPreservesNullHandling(t *testing.T) {
	t.Run("ObjectNullFieldsMatchStockKiota", func(t *testing.T) {
		body := `{
			"uuid":"154660b9-b4a3-4886-8d68-5e170c0923b8",
			"description":"AIP store",
			"enabled":true,
			"path":"/var/aips",
			"pipeline":["/api/v2/pipeline/1/"],
			"purpose":"AS",
			"quota":null,
			"relative_path":"",
			"resource_uri":"/api/v2/location/154660b9-b4a3-4886-8d68-5e170c0923b8/",
			"space":"/api/v2/space/1/",
			"used":10
		}`

		tolerantClient, _ := newClientForJSONResponse(t, body)
		tolerantLocation, err := tolerantClient.Locations().Get(context.Background(), uuid.MustParse("154660b9-b4a3-4886-8d68-5e170c0923b8"))
		if err != nil {
			t.Fatal(err)
		}

		stockClient, stockHTTPClient := newClientForJSONResponse(t, body)
		replaceWithStockKiotaAdapter(t, stockClient, stockHTTPClient)
		stockLocation, err := stockClient.Locations().Get(context.Background(), uuid.MustParse("154660b9-b4a3-4886-8d68-5e170c0923b8"))
		if err != nil {
			t.Fatal(err)
		}

		if tolerantLocation == nil || stockLocation == nil {
			t.Fatal("expected locations")
		}
		if got, want := tolerantLocation.GetQuota(), stockLocation.GetQuota(); !reflect.DeepEqual(got, want) {
			t.Fatalf("quota=%#v want %#v", got, want)
		}
	})

	t.Run("NestedNullFieldsMatchStockKiota", func(t *testing.T) {
		body := `{
			"meta":{
				"limit":20,
				"next":"/api/v2/location/?limit=20&offset=20",
				"offset":0,
				"previous":null,
				"total_count":2
			},
			"objects":[]
		}`

		tolerantClient, _ := newClientForJSONResponse(t, body)
		tolerantList, err := tolerantClient.Locations().List(context.Background(), ssclient.ListLocationsQuery{})
		if err != nil {
			t.Fatal(err)
		}

		stockClient, stockHTTPClient := newClientForJSONResponse(t, body)
		replaceWithStockKiotaAdapter(t, stockClient, stockHTTPClient)
		stockList, err := stockClient.Locations().List(context.Background(), ssclient.ListLocationsQuery{})
		if err != nil {
			t.Fatal(err)
		}

		if tolerantList == nil || stockList == nil {
			t.Fatal("expected lists")
		}
		if got, want := tolerantList.GetMeta().GetPrevious(), stockList.GetMeta().GetPrevious(); !reflect.DeepEqual(got, want) {
			t.Fatalf("meta.previous=%#v want %#v", got, want)
		}
	})
}

func TestKiotaRejectsUnsupportedTimestampLayouts(t *testing.T) {
	t.Run("FixityTimestampRejectsUnsupportedLayout", func(t *testing.T) {
		body := `{"success":true,"message":"","timestamp":"Thu, 19 Mar 2026 05:34:31 UTC"}`

		client, _ := newClientForJSONResponse(t, body)
		res, err := client.Packages().CheckFixity(context.Background(), uuid.MustParse("7c8a3549-2fe0-41d3-9d83-f485f1a43be3"), ssclient.CheckFixityOptions{})
		if err == nil {
			t.Fatal("expected parse error")
		}
		if res != nil {
			t.Fatal("did not expect response")
		}
		if !strings.Contains(err.Error(), "unsupported timestamp layout") {
			t.Fatalf("unexpected error %v", err)
		}
	})

	t.Run("StoredDateRejectsUnsupportedLayout", func(t *testing.T) {
		body := `{
			"uuid":"7c8a3549-2fe0-41d3-9d83-f485f1a43be3",
			"status":"UPLOADED",
			"stored_date":"03/19/2026 05:34:31"
		}`

		client, _ := newClientForJSONResponse(t, body)
		pkg, err := client.Packages().Get(context.Background(), uuid.MustParse("7c8a3549-2fe0-41d3-9d83-f485f1a43be3"))
		if err == nil {
			t.Fatal("expected parse error")
		}
		if pkg != nil {
			t.Fatal("did not expect package")
		}
		if !strings.Contains(err.Error(), "unsupported timestamp layout") {
			t.Fatalf("unexpected error %v", err)
		}
	})
}
