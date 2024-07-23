package ssclient

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	kabs "github.com/microsoft/kiota-abstractions-go"
	khttp "github.com/microsoft/kiota-http-go"

	"go.artefactual.dev/ssclient/kiota"
)

func New(httpClient *http.Client, baseURL, username, key string) (*kiota.Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if err := configureMiddleware(httpClient); err != nil {
		return nil, fmt.Errorf("configure client middleware: %v", err)
	}

	adapter, err := khttp.NewNetHttpRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		&authProvider{username: username, key: key},
		nil,
		nil,
		httpClient,
	)
	if err != nil {
		return nil, fmt.Errorf("create client adapter: %v", err)
	}

	adapter.SetBaseUrl(baseURL)

	return kiota.NewClient(adapter), nil
}

// configureMiddleware installs the middlewares needed by this client.
func configureMiddleware(client *http.Client) error {
	var middlewares []khttp.Middleware

	userAgentOpts := khttp.UserAgentHandlerOptions{
		Enabled:        true,
		ProductName:    "ssclient-go",
		ProductVersion: "v0",
	}
	compressionOpts := khttp.NewCompressionOptions(false)
	retryOpts := khttp.RetryHandlerOptions{
		ShouldRetry: func(delay time.Duration, executionCount int, request *http.Request, response *http.Response) bool {
			// TODO: we use go-retryablehttp but this could be provided instead.
			return false
		},
	}

	// We rely on the default set of middlewares provided by Kiota with a small
	// number of customizations.
	middlewares, err := khttp.GetDefaultMiddlewaresWithOptions(
		&userAgentOpts,
		&compressionOpts,
		&retryOpts,
	)
	if err != nil {
		return err
	}

	// Introduce our middleware to inject the trailing slash.
	middlewares = append(middlewares, appendTrailingSlashHandler{})

	client.Transport = khttp.NewCustomTransportWithParentTransport(client.Transport, middlewares...)

	return nil
}

type authProvider struct {
	username string
	key      string
}

func NewStorageServiceAuthProvider(username, key string) *authProvider {
	return &authProvider{
		username: username,
		key:      key,
	}
}

func (p *authProvider) AuthenticateRequest(ctx context.Context, request *kabs.RequestInformation, additionalAuthenticationContext map[string]interface{}) error {
	request.Headers.Add("Authorization", fmt.Sprintf("ApiKey %s:%s", p.username, p.key))

	return nil
}

// appendTrailingSlashHandler is a middleware that ensures that the path has a
// trailing slash which is the expected in Archivematica Storage Service API.
// With TypeSpec we can sucessfully describe paths with a trailing slash, but
// Kiota drops it during generation (issue #4291).
type appendTrailingSlashHandler struct{}

func (middleware appendTrailingSlashHandler) Intercept(pipeline khttp.Pipeline, middlewareIndex int, req *http.Request) (*http.Response, error) {
	if !strings.HasSuffix(req.URL.Path, "/") {
		req.URL.Path += "/"
	}

	return pipeline.Next(req, middlewareIndex)
}
