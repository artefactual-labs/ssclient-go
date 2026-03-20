package ssclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	kabs "github.com/microsoft/kiota-abstractions-go"
	khttp "github.com/microsoft/kiota-http-go"
)

// responseSnapshot captures the native HTTP response for wrapper methods that
// cannot rely on Kiota's generated method signatures, for example endpoints
// with multiple successful 2xx response shapes.
type responseSnapshot struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

type responseStream struct {
	StatusCode int
	Headers    http.Header
	Body       io.ReadCloser
}

func (r *responseSnapshot) decodeJSON(dst any) error {
	if dst == nil {
		return fmt.Errorf("decode target is required")
	}
	if len(r.Body) == 0 {
		return fmt.Errorf("response body is empty")
	}

	if err := json.Unmarshal(r.Body, dst); err != nil {
		return fmt.Errorf("decode JSON response: %w", err)
	}

	return nil
}

// execute fully consumes the response body and therefore closes it before
// returning the buffered snapshot.
func (c *Client) execute(ctx context.Context, requestInfo *kabs.RequestInformation) (*responseSnapshot, error) {
	resp, err := c.executeStream(ctx, requestInfo)
	if err != nil {
		return nil, err
	}
	defer closeBody(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	return &responseSnapshot{
		StatusCode: resp.StatusCode,
		Headers:    cloneHeaders(resp.Headers),
		Body:       body,
	}, nil
}

// executeStream returns the live response body to the caller, so the caller is
// responsible for closing it on successful requests.
//
//nolint:contextcheck // The request inherits the caller context and appends Kiota request options before dispatch.
func (c *Client) executeStream(ctx context.Context, requestInfo *kabs.RequestInformation) (*responseStream, error) {
	requestInfo.AddRequestOptions([]kabs.RequestOption{
		&khttp.RedirectHandlerOptions{
			ShouldRedirect: func(req *http.Request, res *http.Response) bool {
				return false
			},
		},
	})

	requestAdapter := c.adapter
	if c.raw != nil && c.raw.RequestAdapter != nil {
		requestAdapter = c.raw.RequestAdapter
	}

	ensureRequestBaseURL(requestAdapter, requestInfo)

	native, err := requestAdapter.ConvertToNativeRequest(ctx, requestInfo)
	if err != nil {
		return nil, err
	}

	req, ok := native.(*http.Request)
	if !ok {
		return nil, fmt.Errorf("unexpected native request type %T", native)
	}

	if ctx == nil {
		ctx = req.Context()
	}
	for _, value := range requestInfo.GetRequestOptions() {
		ctx = context.WithValue(ctx, value.GetKey(), value)
	}
	req = req.WithContext(ctx)

	client := cloneHTTPClient(c.httpClient)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	//nolint:bodyclose // The caller owns the streamed response body on success.
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return &responseStream{
		StatusCode: resp.StatusCode,
		Headers:    cloneHeaders(resp.Header),
		Body:       resp.Body,
	}, nil
}

func ensureRequestBaseURL(adapter kabs.RequestAdapter, requestInfo *kabs.RequestInformation) {
	if adapter == nil || requestInfo == nil {
		return
	}
	if !strings.Contains(strings.ToLower(requestInfo.UrlTemplate), "{+baseurl}") {
		return
	}
	if requestInfo.PathParameters == nil {
		requestInfo.PathParameters = make(map[string]string)
	}
	if requestInfo.PathParameters["baseurl"] != "" {
		return
	}

	requestInfo.PathParameters["baseurl"] = adapter.GetBaseUrl()
}

func cloneHeaders(headers http.Header) http.Header {
	cloned := make(http.Header, len(headers))
	for key, values := range headers {
		cloned[key] = append([]string(nil), values...)
	}
	return cloned
}

func closeBody(body io.ReadCloser) {
	if body == nil {
		return
	}
	_ = body.Close()
}
