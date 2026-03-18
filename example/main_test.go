package main

import (
	"bytes"
	"context"
	"net/http"
	"reflect"
	"strings"
	"testing"

	kabs "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
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

	locationList := models.NewLocationList()
	locationCount := int32(1)
	locationMeta := models.NewListResponseMeta()
	locationMeta.SetTotalCount(&locationCount)
	locationList.SetMeta(locationMeta)

	locationID := "fff70864-a5d4-4ca6-ab29-b4ce67d8eeab"
	locationPurpose := models.DS_LOCATIONPURPOSE
	location := models.NewLocation()
	location.SetUuid(&locationID)
	location.SetPurpose(&locationPurpose)
	locationList.SetObjects([]models.Locationable{location})

	pipelineList := models.NewPipelineList()
	pipelineCount := int32(1)
	pipelineMeta := models.NewListResponseMeta()
	pipelineMeta.SetTotalCount(&pipelineCount)
	pipelineList.SetMeta(pipelineMeta)

	pipelineID := "4b9e8af5-b0af-4abf-80b8-4b7d76281f61"
	pipelineName := "Archivematica"
	pipelineDescription := "Default pipeline"
	pipeline := models.NewPipeline()
	pipeline.SetUuid(&pipelineID)
	pipeline.SetRemoteName(&pipelineName)
	pipeline.SetDescription(&pipelineDescription)
	pipelineList.SetObjects([]models.Pipelineable{pipeline})

	adapter.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(4).
		DoAndReturn(func(_ context.Context, requestInfo *kabs.RequestInformation, _ serialization.ParsableFactory, _ kabs.ErrorMappings) (serialization.Parsable, error) {
			switch requestInfo.UrlTemplate {
			case "{+baseurl}/api/v2/location/{?description,limit,offset,order_by,pipeline__uuid,purpose,quota,relative_path,used,uuid}":
				return locationList, nil
			case "{+baseurl}/api/v2/pipeline/{?description,uuid}":
				return pipelineList, nil
			default:
				t.Fatalf("unexpected request template %q", requestInfo.UrlTemplate)
				return nil, nil
			}
		})

	app := application{
		client: client,
		raw:    client.Raw().Api().V2(),
		stdout: stdout,
	}
	if err := app.wrapped(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := app.rawClient(context.Background()); err != nil {
		t.Fatal(err)
	}

	want := strings.Join([]string{
		"",
		strings.Repeat("=", len("Using the wrapped client")),
		"Using the wrapped client",
		strings.Repeat("=", len("Using the wrapped client")),
		"",
		"Found 1 pipelines!",
		"» Pipeline 4b9e8af5-b0af-4abf-80b8-4b7d76281f61 with remote name Archivematica.",
		"Found 1 locations!",
		"» Location fff70864-a5d4-4ca6-ab29-b4ce67d8eeab with purpose DS.",
		"",
		strings.Repeat("=", len("Using the raw client")),
		"Using the raw client",
		strings.Repeat("=", len("Using the raw client")),
		"",
		"Found 1 pipelines!",
		"» Pipeline 4b9e8af5-b0af-4abf-80b8-4b7d76281f61 with remote name Archivematica.",
		"Found 1 locations!",
		"» Location fff70864-a5d4-4ca6-ab29-b4ce67d8eeab with purpose DS.",
		"",
	}, "\n")
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
