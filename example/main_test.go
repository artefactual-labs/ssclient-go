package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	kiotahttpgo "github.com/microsoft/kiota-http-go"
	"go.uber.org/mock/gomock"

	"go.artefactual.dev/ssclient"
	"go.artefactual.dev/ssclient/example/adapter"
	"go.artefactual.dev/ssclient/kiota"
	"go.artefactual.dev/ssclient/kiota/models"
)

func assertEqual(t *testing.T, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Mismatch found:\nGot: %#v\nWant: %#v", got, want)
	}
}

func TestRunUsage(t *testing.T) {
	ctx := context.Background()
	stdout := bytes.NewBuffer([]byte{})

	err := run(ctx, stdout, []string{"example"})
	assertEqual(t, err.Error(), "Usage: example -url=http://127.0.0.1:62081 -user=test -key=test")
}

func TestRun(t *testing.T) {
	ctx := context.Background()
	stdout := bytes.NewBuffer([]byte{})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, r.URL.Path, "/api/v2/location")
		assertEqual(t, r.Header, http.Header(map[string][]string{
			"Accept":          {"application/json"},
			"Accept-Encoding": {"gzip"},
			"Authorization":   {"ApiKey test:test"},
			"User-Agent":      {"kiota-go/" + kiotahttpgo.NewUserAgentHandlerOptions().ProductVersion},
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

	if err := run(ctx, stdout, []string{"-url=" + srv.URL, "-user=test", "-key=test"}); err != nil {
		t.Fatal(err)
	}

	want := `Found 1 locations!
» Location fff70864-a5d4-4ca6-ab29-b4ce67d8eeab with purpose DS.
`
	if got := stdout.String(); want != got {
		t.Fatalf("unexpected output; got:\n%v\nwant:\n%v\n", got, want)
	}
}

// TestApplication tests the application using a fake client.
func TestApplication(t *testing.T) {
	stdout := bytes.NewBuffer([]byte{})
	client, adapter := createClient(t)

	// Have our fake client return a list with a single location.
	{
		list := models.NewLocationList()

		count := int32(1)
		meta := models.NewListResponseMeta()
		meta.SetTotalCount(&count)
		list.SetMeta(meta)

		id := "fff70864-a5d4-4ca6-ab29-b4ce67d8eeab"
		ps := models.DS_LOCATIONPURPOSE
		location := models.NewLocation()
		location.SetUuid(&id)
		location.SetPurpose(&ps)
		list.SetObjects([]models.Locationable{location})

		param := gomock.Any()
		adapter.EXPECT().Send(param, param, param, param).Times(1).Return(list, nil)
	}

	app := application{client, stdout}
	if err := app.locations(context.Background()); err != nil {
		t.Fatal(err)
	}

	want := `Found 1 locations!
» Location fff70864-a5d4-4ca6-ab29-b4ce67d8eeab with purpose DS.
`
	if got := stdout.String(); want != got {
		t.Fatalf("unexpected output; got:\n%v\nwant:\n%v\n", got, want)
	}
}

func createClient(t *testing.T) (*kiota.Client, *adapter.MockRequestAdapter) {
	t.Helper()

	client, err := ssclient.New(&http.Client{}, "http://127.0.0.1:62081", "test", "test")
	if err != nil {
		t.Fatal(err)
	}

	adapter := adapter.NewMockRequestAdapter(gomock.NewController(t))
	client.RequestAdapter = adapter

	return client, adapter
}
