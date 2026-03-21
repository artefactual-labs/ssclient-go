package ssclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/textproto"

	kabs "github.com/microsoft/kiota-abstractions-go"
)

// ResponseError describes a non-success HTTP response returned by Storage
// Service.
type ResponseError struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Message    string
	Detail     string
	Cause      error
}

// NotAvailableError describes a 202 response indicating the requested content
// is not locally available for download yet.
type NotAvailableError struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Message    string
}

// Error returns a readable representation of the failed response.
func (e *ResponseError) Error() string {
	if e == nil {
		return "<nil>"
	}

	message := e.Message
	if message == "" && e.Cause != nil {
		message = e.Cause.Error()
	}

	switch {
	case e.StatusCode > 0 && message != "" && e.Detail != "":
		return fmt.Sprintf("storage service request failed with status %d: %s (%s)", e.StatusCode, message, e.Detail)
	case e.StatusCode > 0 && message != "":
		return fmt.Sprintf("storage service request failed with status %d: %s", e.StatusCode, message)
	case e.StatusCode > 0:
		return fmt.Sprintf("storage service request failed with status %d", e.StatusCode)
	case message != "" && e.Detail != "":
		return fmt.Sprintf("%s (%s)", message, e.Detail)
	case message != "":
		return message
	case e.Cause != nil:
		return e.Cause.Error()
	default:
		return "storage service request failed"
	}
}

// Unwrap returns the original error when this error wraps one.
func (e *ResponseError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// Error returns a readable representation of the unavailable response.
func (e *NotAvailableError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Message != "" {
		return e.Message
	}
	if e.StatusCode > 0 {
		return fmt.Sprintf("storage service content unavailable with status %d", e.StatusCode)
	}
	return "storage service content unavailable"
}

// StatusCode returns the HTTP status code from a supported Storage Service
// error, if present.
func StatusCode(err error) (int, bool) {
	var responseErr *ResponseError
	if errors.As(err, &responseErr) && responseErr.StatusCode > 0 {
		return responseErr.StatusCode, true
	}
	var unavailableErr *NotAvailableError
	if errors.As(err, &unavailableErr) && unavailableErr.StatusCode > 0 {
		return unavailableErr.StatusCode, true
	}
	return 0, false
}

// IsNotFound reports whether err represents a 404 response.
func IsNotFound(err error) bool {
	status, ok := StatusCode(err)
	return ok && status == http.StatusNotFound
}

func normalizeError(err error) error {
	if err == nil {
		return nil
	}

	var responseErr *ResponseError
	if errors.As(err, &responseErr) {
		return err
	}

	normalized := &ResponseError{
		Cause:   err,
		Message: extractErrorMessage(err),
		Detail:  extractErrorDetail(err),
	}

	var apiErr kabs.ApiErrorable
	if errors.As(err, &apiErr) {
		normalized.StatusCode = apiErr.GetStatusCode()
		normalized.Headers = responseHeadersToHTTPHeader(apiErr.GetResponseHeaders())
	}

	if normalized.Message == "" {
		normalized.Message = err.Error()
	}

	return normalized
}

func newResponseErrorFromSnapshot(resp *responseSnapshot, fallbackMessage string) *ResponseError {
	if resp == nil {
		return &ResponseError{Message: fallbackMessage}
	}

	message, detail := decodeErrorPayload(resp.Body)
	if message == "" {
		message = fallbackMessage
	}

	return &ResponseError{
		StatusCode: resp.StatusCode,
		Headers:    cloneHeaders(resp.Headers),
		Body:       append([]byte(nil), resp.Body...),
		Message:    message,
		Detail:     detail,
	}
}

func newNotAvailableErrorFromSnapshot(resp *responseSnapshot, fallbackMessage string) *NotAvailableError {
	if resp == nil {
		return &NotAvailableError{Message: fallbackMessage}
	}

	message, _ := decodeErrorPayload(resp.Body)
	if message == "" {
		message = fallbackMessage
	}

	return &NotAvailableError{
		StatusCode: resp.StatusCode,
		Headers:    cloneHeaders(resp.Headers),
		Body:       append([]byte(nil), resp.Body...),
		Message:    message,
	}
}

func extractErrorMessage(err error) string {
	type errorGetter interface {
		GetErrorEscaped() *string
	}
	var withError errorGetter
	if errors.As(err, &withError) && withError.GetErrorEscaped() != nil {
		return *withError.GetErrorEscaped()
	}

	type errorMessageGetter interface {
		GetErrorMessage() *string
	}
	var withErrorMessage errorMessageGetter
	if errors.As(err, &withErrorMessage) && withErrorMessage.GetErrorMessage() != nil {
		return *withErrorMessage.GetErrorMessage()
	}

	type messageGetter interface {
		GetMessage() *string
	}
	var withMessage messageGetter
	if errors.As(err, &withMessage) && withMessage.GetMessage() != nil {
		return *withMessage.GetMessage()
	}

	return ""
}

func extractErrorDetail(err error) string {
	type detailGetter interface {
		GetDetail() *string
	}
	var withDetail detailGetter
	if errors.As(err, &withDetail) && withDetail.GetDetail() != nil {
		return *withDetail.GetDetail()
	}

	return ""
}

func decodeErrorPayload(body []byte) (message, detail string) {
	if len(body) == 0 {
		return "", ""
	}

	var payload map[string]string
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", ""
	}

	switch {
	case payload["error_message"] != "":
		message = payload["error_message"]
	case payload["error"] != "":
		message = payload["error"]
	case payload["message"] != "":
		message = payload["message"]
	}

	detail = payload["detail"]
	return message, detail
}

func responseHeadersToHTTPHeader(headers *kabs.ResponseHeaders) http.Header {
	if headers == nil {
		return nil
	}

	result := make(http.Header, len(headers.ListKeys()))
	for _, key := range headers.ListKeys() {
		result[textproto.CanonicalMIMEHeaderKey(key)] = append([]string(nil), headers.Get(key)...)
	}

	return result
}
