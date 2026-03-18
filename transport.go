package ssclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

func (c *Client) execute(ctx context.Context, requestInfo *kabs.RequestInformation) (*responseSnapshot, error) {
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

	native, err := requestAdapter.ConvertToNativeRequest(ctx, requestInfo)
	if err != nil {
		return nil, err
	}

	req, ok := native.(*http.Request)
	if !ok {
		return nil, fmt.Errorf("unexpected native request type %T", native)
	}

	reqCtx := ctx
	if reqCtx == nil {
		reqCtx = context.Background()
	}
	for _, value := range requestInfo.GetRequestOptions() {
		reqCtx = context.WithValue(reqCtx, value.GetKey(), value)
	}
	req = req.WithContext(reqCtx)

	client := cloneHTTPClient(c.httpClient)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	return &responseSnapshot{
		StatusCode: resp.StatusCode,
		Headers:    cloneHeaders(resp.Header),
		Body:       body,
	}, nil
}

func cloneHeaders(headers http.Header) http.Header {
	cloned := make(http.Header, len(headers))
	for key, values := range headers {
		cloned[key] = append([]string(nil), values...)
	}
	return cloned
}
