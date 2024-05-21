package ssclient_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"go.artefactual.dev/ssclient"
)

func assertEqual(t *testing.T, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Mismatch found:\nGot: %#v\nWant: %#v", got, want)
	}
}

func TestClient(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, r.URL.Path, "/api/v2/location/")
		assertEqual(t, r.Header, http.Header(map[string][]string{
			"Accept":          {"application/json"},
			"Accept-Encoding": {"gzip"},
			"Authorization":   {"ApiKey test:test"},
			"User-Agent":      {"ssclient-go/v0"},
		}))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"meta": {
				"limit": 20,
				"next": null,
				"offset": 0,
				"previous": null,
				"total_count": 1
			},
			"objects": [
				{
					"enabled": true,
					"path": "/home",
					"pipeline": ["/api/v2/pipeline/a64e061a-5688-49b5-95c1-0b6885c40c04/"],
					"purpose": "DS",
					"quota": null,
					"relative_path": "home",
					"resource_uri": "/api/v2/location/fff70864-a5d4-4ca6-ab29-b4ce67d8eeab/",
					"space": "/api/v2/space/218caeb7-fd59-4b7b-99b1-f5771a2dd34f/",
					"uuid": "fff70864-a5d4-4ca6-ab29-b4ce67d8eeab"
				}
			]
		}`))
	}))
	t.Cleanup(func() { srv.Close() })

	client, err := ssclient.New(&http.Client{}, srv.URL, "test", "test")
	assertEqual(t, err, nil)

	locations, err := client.Api().V2().Location().Get(context.Background(), nil)
	assertEqual(t, err, nil)
	assertEqual(t, len(locations.GetObjects()), 1)
}
