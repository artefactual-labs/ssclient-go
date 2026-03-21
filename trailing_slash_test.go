package ssclient_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"go.artefactual.dev/ssclient"
	"go.artefactual.dev/ssclient/kiota/models"
)

func TestWrapperMethodsEmitTrailingSlashURLs(t *testing.T) {
	const baseURL = "http://storage.service"

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
}

func TestGeneratedItemEmptyPathSegmentBuildersEmitTrailingSlashURLs(t *testing.T) {
	const baseURL = "http://storage.service"

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
}
