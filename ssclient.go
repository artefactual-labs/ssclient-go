package ssclient

import (
	"context"
	"fmt"
	"net/http"
	"time"

	kabs "github.com/microsoft/kiota-abstractions-go"
	khttp "github.com/microsoft/kiota-http-go"

	"go.artefactual.dev/ssclient/kiota"
)

const (
	defaultUserAgentProductName    = "ssclient-go"
	defaultUserAgentProductVersion = "v0"
)

// Config configures a Storage Service client.
type Config struct {
	BaseURL     string
	Username    string
	Key         string
	HTTPClient  *http.Client
	Middlewares []khttp.Middleware
}

// Client is the public entrypoint for the wrapper around the generated
// Kiota client.
type Client struct {
	raw        *kiota.Client
	adapter    kabs.RequestAdapter
	httpClient *http.Client

	locations *LocationsService
	packages  *PackagesService
	pipelines *PipelinesService
}

// New constructs a Storage Service client backed by the generated Kiota client.
func New(cfg Config) (*Client, error) {
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	httpClient = cloneHTTPClient(httpClient)

	if err := configureMiddleware(httpClient, cfg.Middlewares); err != nil {
		return nil, fmt.Errorf("configure client middleware: %w", err)
	}

	adapter, err := khttp.NewNetHttpRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		&authProvider{username: cfg.Username, key: cfg.Key},
		newTolerantParseNodeFactory(),
		nil,
		httpClient,
	)
	if err != nil {
		return nil, fmt.Errorf("create client adapter: %w", err)
	}

	adapter.SetBaseUrl(cfg.BaseURL)

	client := &Client{
		raw:        kiota.NewClient(adapter),
		adapter:    adapter,
		httpClient: httpClient,
	}
	client.locations = &LocationsService{client: client}
	client.packages = &PackagesService{client: client}
	client.pipelines = &PipelinesService{client: client}

	return client, nil
}

// Raw returns the generated Kiota client as an escape hatch.
func (c *Client) Raw() *kiota.Client {
	return c.raw
}

// Adapter returns the underlying request adapter.
func (c *Client) Adapter() kabs.RequestAdapter {
	return c.adapter
}

// Locations returns location-related operations.
func (c *Client) Locations() *LocationsService {
	return c.locations
}

// Packages returns package-related operations. The Storage Service transport
// uses "/file/" for these resources, but the public API exposes them as
// packages.
func (c *Client) Packages() *PackagesService {
	return c.packages
}

// Pipelines returns pipeline-related operations.
func (c *Client) Pipelines() *PipelinesService {
	return c.pipelines
}

// configureMiddleware installs the middlewares needed by this client.
func configureMiddleware(client *http.Client, extra []khttp.Middleware) error {
	userAgentOpts := khttp.UserAgentHandlerOptions{
		Enabled:        true,
		ProductName:    defaultUserAgentProductName,
		ProductVersion: defaultUserAgentProductVersion,
	}
	compressionOpts := khttp.NewCompressionOptionsReference(false)
	retryOpts := khttp.RetryHandlerOptions{
		ShouldRetry: func(delay time.Duration, executionCount int, request *http.Request, response *http.Response) bool {
			return false
		},
	}

	// We rely on the default set of middlewares provided by Kiota with a small
	// number of customizations.
	middlewares, err := khttp.GetDefaultMiddlewaresWithOptions(
		&userAgentOpts,
		compressionOpts,
		&retryOpts,
	)
	if err != nil {
		return err
	}
	middlewares = append(middlewares, extra...)

	client.Transport = khttp.NewCustomTransportWithParentTransport(client.Transport, middlewares...)

	return nil
}

func cloneHTTPClient(client *http.Client) *http.Client {
	clone := *client
	return &clone
}

type authProvider struct {
	username string
	key      string
}

func (p *authProvider) AuthenticateRequest(ctx context.Context, request *kabs.RequestInformation, additionalAuthenticationContext map[string]any) error {
	request.Headers.Add("Authorization", fmt.Sprintf("ApiKey %s:%s", p.username, p.key))

	return nil
}
