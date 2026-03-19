package ssclient_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	kabs "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"

	"go.artefactual.dev/ssclient"
	"go.artefactual.dev/ssclient/kiota/models"
)

func TestPipelines(t *testing.T) {
	t.Run("List", func(t *testing.T) {
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
				assertEqual(t, requestInfo.QueryParametersAny["uuid"], "a64e061a-5688-49b5-95c1-0b6885c40c04")

				list := models.NewPipelineList()
				pipeline := models.NewPipeline()
				pipeline.SetUuid(ptr(uuid.MustParse("a64e061a-5688-49b5-95c1-0b6885c40c04")))
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
	})

	t.Run("Get", func(t *testing.T) {
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
				assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/pipeline/{uuid}/")
				assertEqual(t, requestInfo.PathParameters["uuid"], pipelineID)

				pipeline := models.NewPipeline()
				pipeline.SetUuid(ptr(uuid.MustParse(pipelineID)))
				return pipeline, nil
			},
		}

		res, err := client.Pipelines().Get(context.Background(), pipelineID)
		assertEqual(t, err, nil)
		if res == nil {
			t.Fatal("expected pipeline")
		}
		assertEqual(t, res.GetUuid().String(), pipelineID)
	})
}
