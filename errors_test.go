package ssclient_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	kabs "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"go.artefactual.dev/ssclient"
	"go.artefactual.dev/ssclient/kiota/models"
)

func TestPackagesGetNormalizesKiotaError(t *testing.T) {
	const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"

	client, err := ssclient.New(ssclient.Config{
		BaseURL:    "http://storage.service",
		Username:   "test",
		Key:        "test",
		HTTPClient: &http.Client{},
	})
	assertEqual(t, err, nil)

	headers := kabs.NewResponseHeaders()
	headers.Add("X-Request-Id", "req-123")

	raw := client.Raw()
	raw.RequestAdapter = &fakeRequestAdapter{
		baseURL: "http://storage.service",
		send: func(ctx context.Context, requestInfo *kabs.RequestInformation, constructor serialization.ParsableFactory, errorMappings kabs.ErrorMappings) (serialization.Parsable, error) {
			errModel := models.NewErrorEscaped()
			errModel.SetErrorEscaped(ptr("bad request"))
			errModel.SetStatusCode(http.StatusBadRequest)
			errModel.SetResponseHeaders(headers)
			return nil, errModel
		},
	}

	_, err = client.Packages().Get(context.Background(), packageID)
	if err == nil {
		t.Fatal("expected error")
	}

	var responseErr *ssclient.ResponseError
	if !errors.As(err, &responseErr) {
		t.Fatalf("expected ResponseError, got %T", err)
	}

	assertEqual(t, responseErr.StatusCode, http.StatusBadRequest)
	assertEqual(t, responseErr.Message, "bad request")
	assertEqual(t, responseErr.Headers.Get("X-Request-Id"), "req-123")

	status, ok := ssclient.StatusCode(err)
	assertEqual(t, ok, true)
	assertEqual(t, status, http.StatusBadRequest)
}

func TestPackagesDeleteAIPNormalizesManualError(t *testing.T) {
	const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"

	client, err := ssclient.New(ssclient.Config{
		BaseURL:  "http://storage.service",
		Username: "test",
		Key:      "test",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusMethodNotAllowed,
					Header:     http.Header{"Content-Type": {"application/json"}, "X-Request-Id": {"req-405"}},
					Body:       io.NopCloser(strings.NewReader(`{"message":"Deletes not allowed on this package type."}`)),
					Request:    r,
				}, nil
			}),
		},
	})
	assertEqual(t, err, nil)

	raw := client.Raw()
	writerFactory := raw.RequestAdapter.GetSerializationWriterFactory()
	raw.RequestAdapter = &fakeRequestAdapter{
		serializationWriterFactory: writerFactory,
		convertToNativeRequest: func(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
			return &http.Request{
				Method: http.MethodPost,
				URL: &url.URL{
					Scheme: "http",
					Host:   "storage.service",
					Path:   "/api/v2/file/" + packageID + "/delete_aip/",
				},
				Header: http.Header{
					"Accept":        {"application/json"},
					"Authorization": {"ApiKey test:test"},
					"Content-Type":  {"application/json"},
				},
				Body: io.NopCloser(strings.NewReader(string(requestInfo.Content))),
			}, nil
		},
	}

	body := models.NewDeleteAipRequest()
	body.SetEventReason(ptr("Delete please!"))
	body.SetPipeline(ptr("4b9e8af5-b0af-4abf-80b8-4b7d76281f61"))
	body.SetUserId(ptr(int32(1)))
	body.SetUserEmail(ptr("user@example.com"))

	_, err = client.Packages().DeleteAIP(context.Background(), packageID, body)
	if err == nil {
		t.Fatal("expected error")
	}

	var responseErr *ssclient.ResponseError
	if !errors.As(err, &responseErr) {
		t.Fatalf("expected ResponseError, got %T", err)
	}

	assertEqual(t, responseErr.StatusCode, http.StatusMethodNotAllowed)
	assertEqual(t, responseErr.Message, "Deletes not allowed on this package type.")
	assertEqual(t, string(responseErr.Body), `{"message":"Deletes not allowed on this package type."}`)
	assertEqual(t, responseErr.Headers.Get("X-Request-Id"), "req-405")
}
