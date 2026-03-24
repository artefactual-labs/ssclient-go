package ssclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"

	kabs "github.com/microsoft/kiota-abstractions-go"
	khttp "github.com/microsoft/kiota-http-go"
)

// responseSnapshot buffers an HTTP response so higher-level wrappers can
// inspect the status, headers, and body after the transport body has been
// consumed. This is primarily used for endpoints whose success responses do
// not map cleanly to a single Kiota-generated shape.
type responseSnapshot struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// decodeJSON unmarshals a buffered JSON response into dst.
//
// The caller must provide a non-nil destination and the snapshot must contain
// a body, otherwise the error should stay explicit because these are client
// misuse cases rather than server responses.
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

// responseStream preserves transport metadata while leaving the body open for
// the caller to consume.
type responseStream struct {
	StatusCode int
	Headers    http.Header
	Body       io.ReadCloser
}

// FileStream describes a successful streamed file response and exposes common
// metadata derived from the HTTP headers alongside the live body reader.
type FileStream struct {
	StatusCode         int
	ContentType        string
	ContentLength      int64
	ContentDisposition string
	Filename           string
	Body               io.ReadCloser
}

// execute sends requestInfo and fully buffers the response body. The body is
// always closed before the method returns.
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

// executeStream sends requestInfo and returns the native HTTP response metadata
// plus the live body stream. Successful callers take ownership of closing the
// returned body.
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

// streamRequest executes a request whose success path is expected to be an
// HTTP 200 response with a streamed body. Non-200 responses are normalized into
// the package's domain errors after buffering the response body.
func (c *Client) streamRequest(ctx context.Context, requestInfo *kabs.RequestInformation, action string) (*FileStream, error) {
	resp, err := c.executeStream(ctx, requestInfo)
	if err != nil {
		return nil, normalizeError(err)
	}

	if resp.StatusCode != http.StatusOK {
		defer closeBody(resp.Body)

		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, normalizeError(fmt.Errorf("read %s error response body: %w", action, readErr))
		}

		snapshot := &responseSnapshot{
			StatusCode: resp.StatusCode,
			Headers:    resp.Headers,
			Body:       body,
		}
		if resp.StatusCode == http.StatusAccepted {
			return nil, newNotAvailableErrorFromSnapshot(snapshot, fmt.Sprintf("%s not locally available", action))
		}

		return nil, newResponseErrorFromSnapshot(snapshot, fmt.Sprintf("unexpected %s response status %d", action, resp.StatusCode))
	}

	return &FileStream{
		StatusCode:         resp.StatusCode,
		ContentType:        resp.Headers.Get("Content-Type"),
		ContentLength:      parseContentLength(resp.Headers.Get("Content-Length")),
		ContentDisposition: resp.Headers.Get("Content-Disposition"),
		Filename:           parseFilename(resp.Headers.Get("Content-Disposition")),
		Body:               resp.Body,
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

func parseContentLength(value string) int64 {
	if value == "" {
		return -1
	}

	length, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return -1
	}

	return length
}

func parseFilename(value string) string {
	if value == "" {
		return ""
	}

	_, params, err := mime.ParseMediaType(value)
	if err != nil {
		return ""
	}

	return params["filename"]
}
