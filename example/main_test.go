package main

import (
	"bytes"
	"context"
	"net/http"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"go.artefactual.dev/ssclient"
	"go.artefactual.dev/ssclient/example/adapter"
	"go.artefactual.dev/ssclient/kiota/models"
)

func assertEqual(t *testing.T, got, want interface{}) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Mismatch found:\nGot: %#v\nWant: %#v", got, want)
	}
}

func TestRunUsage(t *testing.T) {
	ctx := context.Background()
	stdout := bytes.NewBuffer([]byte{})

	err := run(ctx, stdout, []string{})
	assertEqual(t, err.Error(), "flag: help requested")
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
		adapter.EXPECT().Send(param, param, param, param).Times(2).Return(list, nil)
	}

	app := application{
		client: client,
		raw:    client.Raw().Api().V2(),
		stdout: stdout,
	}
	if err := app.locations(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := app.locationsRaw(context.Background()); err != nil {
		t.Fatal(err)
	}

	want := `Found 1 locations!
» Location fff70864-a5d4-4ca6-ab29-b4ce67d8eeab with purpose DS.
Found 1 locations!
» Location fff70864-a5d4-4ca6-ab29-b4ce67d8eeab with purpose DS.
`
	if got := stdout.String(); want != got {
		t.Fatalf("unexpected output; got:\n%v\nwant:\n%v\n", got, want)
	}
}

func createClient(t *testing.T) (*ssclient.Client, *adapter.MockRequestAdapter) {
	t.Helper()

	client, err := ssclient.New(ssclient.Config{
		BaseURL:    "http://127.0.0.1:62081",
		Username:   "test",
		Key:        "test",
		HTTPClient: &http.Client{},
	})
	if err != nil {
		t.Fatal(err)
	}

	adapter := adapter.NewMockRequestAdapter(gomock.NewController(t))
	client.Raw().RequestAdapter = adapter

	return client, adapter
}
