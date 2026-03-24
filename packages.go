package ssclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	kabs "github.com/microsoft/kiota-abstractions-go"

	kapi "go.artefactual.dev/ssclient/kiota/api"
	"go.artefactual.dev/ssclient/kiota/models"
)

// PackagesService provides package-related API operations.
type PackagesService struct {
	client *Client
}

// Get returns a package by ID.
func (s *PackagesService) Get(ctx context.Context, id uuid.UUID) (*models.PackageEscaped, error) {
	pkg, err := s.client.raw.Api().V2().File().ByUuid(id).EmptyPathSegment().Get(ctx, nil)
	if err != nil {
		return nil, normalizeError(err)
	}
	if pkg == nil {
		return nil, nil
	}

	typed, ok := pkg.(*models.PackageEscaped)
	if !ok {
		return nil, fmt.Errorf("unexpected package type %T", pkg)
	}

	return typed, nil
}

// DownloadPackage returns the package archive as downloaded by Storage Service.
func (s *PackagesService) DownloadPackage(ctx context.Context, id uuid.UUID) (*FileStream, error) {
	requestInfo, err := s.client.raw.Api().V2().File().ByUuid(id).Download().EmptyPathSegment().ToGetRequestInformation(ctx, nil)
	if err != nil {
		return nil, normalizeError(err)
	}
	requestInfo.Headers.Remove("Accept")
	requestInfo.Headers.Add("Accept", "*/*")

	return s.client.streamRequest(ctx, requestInfo, "download")
}

// DownloadFile extracts a file from a package and streams it back.
func (s *PackagesService) DownloadFile(ctx context.Context, id uuid.UUID, relativePathToFile string) (*FileStream, error) {
	if relativePathToFile == "" {
		return nil, fmt.Errorf("relative path to file is required")
	}

	requestInfo, err := s.client.raw.Api().V2().File().ByUuid(id).Extract_file().EmptyPathSegment().ToGetRequestInformation(ctx, &kabs.RequestConfiguration[kapi.V2FileItemExtract_fileEmptyPathSegmentRequestBuilderGetQueryParameters]{
		QueryParameters: &kapi.V2FileItemExtract_fileEmptyPathSegmentRequestBuilderGetQueryParameters{
			Relative_path_to_file: &relativePathToFile,
		},
	})
	if err != nil {
		return nil, normalizeError(err)
	}
	requestInfo.Headers.Remove("Accept")
	requestInfo.Headers.Add("Accept", "*/*")

	return s.client.streamRequest(ctx, requestInfo, "extract file")
}

// DownloadPointerFile returns the package pointer file as a stream.
func (s *PackagesService) DownloadPointerFile(ctx context.Context, id uuid.UUID) (*FileStream, error) {
	requestInfo, err := s.client.raw.Api().V2().File().ByUuid(id).Pointer_file().EmptyPathSegment().ToGetRequestInformation(ctx, nil)
	if err != nil {
		return nil, normalizeError(err)
	}
	requestInfo.Headers.Remove("Accept")
	requestInfo.Headers.Add("Accept", "*/*")

	return s.client.streamRequest(ctx, requestInfo, "pointer file")
}

// DeleteAIPAccepted is returned when the server creates a new deletion request.
type DeleteAIPAccepted struct {
	Message string `json:"message"`
	ID      int32  `json:"id"`
}

// DeleteAIPAlreadyExists is returned when a pending deletion request already
// exists for the package.
type DeleteAIPAlreadyExists struct {
	ErrorMessage string `json:"error_message"`
}

// DeleteAIPResult preserves the two non-error outcomes exposed by the Storage
// Service delete request endpoint.
type DeleteAIPResult struct {
	Accepted      *DeleteAIPAccepted
	AlreadyExists *DeleteAIPAlreadyExists
}

// IsAccepted reports whether the request created a new deletion event.
func (r *DeleteAIPResult) IsAccepted() bool {
	return r != nil && r.Accepted != nil
}

// HasExistingRequest reports whether the package already had a pending
// deletion request.
func (r *DeleteAIPResult) HasExistingRequest() bool {
	return r != nil && r.AlreadyExists != nil
}

// DeleteAIP creates an AIP deletion request for the given package. The server
// exposes two non-error outcomes: one for a newly created request and one for
// an already pending request.
func (s *PackagesService) DeleteAIP(ctx context.Context, id uuid.UUID, body *models.DeleteAipRequest) (*DeleteAIPResult, error) {
	if body == nil {
		return nil, fmt.Errorf("delete AIP request is required")
	}

	builder := s.client.raw.Api().V2().File().ByUuid(id).Delete_aip().EmptyPathSegment()
	requestInfo, err := builder.ToPostRequestInformation(ctx, body, nil)
	if err != nil {
		return nil, normalizeError(err)
	}

	resp, err := s.client.execute(ctx, requestInfo)
	if err != nil {
		return nil, normalizeError(err)
	}

	switch resp.StatusCode {
	case http.StatusAccepted:
		var accepted DeleteAIPAccepted
		if err := resp.decodeJSON(&accepted); err != nil {
			return nil, normalizeError(err)
		}
		return &DeleteAIPResult{
			Accepted: &accepted,
		}, nil
	case http.StatusOK:
		var existing DeleteAIPAlreadyExists
		if err := resp.decodeJSON(&existing); err != nil {
			return nil, normalizeError(err)
		}
		return &DeleteAIPResult{
			AlreadyExists: &existing,
		}, nil
	case http.StatusBadRequest:
		return nil, newResponseErrorFromSnapshot(resp, "delete AIP request bad request")
	case http.StatusMethodNotAllowed:
		return nil, newResponseErrorFromSnapshot(resp, "delete AIP request not allowed")
	case http.StatusNotFound:
		return nil, newResponseErrorFromSnapshot(resp, "delete AIP request not found")
	default:
		return nil, newResponseErrorFromSnapshot(resp, fmt.Sprintf("unexpected delete AIP response status %d", resp.StatusCode))
	}
}

// CheckFixityOptions configures a package fixity check request.
type CheckFixityOptions struct {
	ForceLocal *bool
}

// CheckFixity runs a fixity check for the given package ID.
func (s *PackagesService) CheckFixity(ctx context.Context, id uuid.UUID, opts CheckFixityOptions) (*models.FixityResponse, error) {
	reqConfig := &kabs.RequestConfiguration[kapi.V2FileItemCheck_fixityEmptyPathSegmentRequestBuilderGetQueryParameters]{}
	if opts.ForceLocal != nil {
		reqConfig.QueryParameters = &kapi.V2FileItemCheck_fixityEmptyPathSegmentRequestBuilderGetQueryParameters{
			Force_local: opts.ForceLocal,
		}
	}

	res, err := s.client.raw.Api().V2().File().ByUuid(id).Check_fixity().EmptyPathSegment().Get(ctx, reqConfig)
	if err != nil {
		return nil, normalizeError(err)
	}
	if res == nil {
		return nil, nil
	}

	typed, ok := res.(*models.FixityResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected fixity response type %T", res)
	}

	return typed, nil
}

// MoveResult captures the accepted async package move response.
type MoveResult struct {
	Location string
}

// Move starts an asynchronous package move to a different storage location.
// The returned Location header points at the async operation resource.
func (s *PackagesService) Move(ctx context.Context, id, locationID uuid.UUID) (*MoveResult, error) {
	values := url.Values{}
	values.Set("location_uuid", locationID.String())

	requestInfo := kabs.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(
		kabs.POST,
		"{+baseurl}/api/v2/file/{uuid}/move/",
		map[string]string{
			"uuid": id.String(),
		},
	)
	requestInfo.Headers.TryAdd("Accept", "application/json")
	requestInfo.Headers.TryAdd("Content-Type", "application/x-www-form-urlencoded")
	requestInfo.Content = []byte(values.Encode())

	resp, err := s.client.execute(ctx, requestInfo)
	if err != nil {
		return nil, normalizeError(err)
	}

	switch resp.StatusCode {
	case http.StatusAccepted:
		location := resp.Headers.Get("Location")
		if location == "" {
			return nil, fmt.Errorf("missing Location header")
		}
		return &MoveResult{
			Location: location,
		}, nil
	case http.StatusBadRequest:
		return nil, newResponseErrorFromSnapshot(resp, "package move request bad request")
	case http.StatusNotFound:
		return nil, newResponseErrorFromSnapshot(resp, "package move request not found")
	default:
		return nil, newResponseErrorFromSnapshot(resp, fmt.Sprintf("unexpected package move response status %d", resp.StatusCode))
	}
}

// ReviewAIPDeletionSuccess captures a successful review response.
type ReviewAIPDeletionSuccess struct {
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// ReviewAIPDeletionError describes an application-level review failure
// returned by the review AIP deletion endpoint with HTTP 200.
type ReviewAIPDeletionError struct {
	ErrorMessage string
	Detail       string
}

// Error returns a readable representation of the application-level review
// failure.
func (e *ReviewAIPDeletionError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.ErrorMessage != "" && e.Detail != "" {
		return fmt.Sprintf("%s (%s)", e.ErrorMessage, e.Detail)
	}
	if e.ErrorMessage != "" {
		return e.ErrorMessage
	}
	if e.Detail != "" {
		return e.Detail
	}
	return "review AIP deletion failed"
}

// ReviewAIPDeletion approves or rejects an AIP deletion request associated with
// a package.
//
// Storage Service reports two distinct business outcomes for this endpoint
// using the same HTTP 200 status code and content type:
//   - success bodies use {"message": ...}
//   - application-level failure bodies use {"error_message": ...}
//
// That shape is part of the deployed API, but it is awkward for generated
// clients because there is no discriminator beyond the JSON fields themselves.
// The wrapper therefore executes the request and inspects the response body
// directly so callers can rely on err for both transport/protocol failures and
// application-level review failures. Callers that need structured failure
// details can use errors.As with *ReviewAIPDeletionError.
func (s *PackagesService) ReviewAIPDeletion(ctx context.Context, id uuid.UUID, body *models.ReviewAipDeletionRequest) (*ReviewAIPDeletionSuccess, error) {
	if body == nil {
		return nil, fmt.Errorf("review AIP deletion request is required")
	}

	builder := s.client.raw.Api().V2().File().ByUuid(id).Review_aip_deletion().EmptyPathSegment()
	requestInfo, err := builder.ToPostRequestInformation(ctx, body, nil)
	if err != nil {
		return nil, normalizeError(err)
	}

	resp, err := s.client.execute(ctx, requestInfo)
	if err != nil {
		return nil, normalizeError(err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		if hasJSONField(resp.Body, "error_message") {
			var failure struct {
				ErrorMessage string `json:"error_message"`
				Detail       string `json:"detail,omitempty"`
			}
			if err := resp.decodeJSON(&failure); err != nil {
				return nil, normalizeError(err)
			}
			return nil, &ReviewAIPDeletionError{
				ErrorMessage: failure.ErrorMessage,
				Detail:       failure.Detail,
			}
		}

		var success ReviewAIPDeletionSuccess
		if err := resp.decodeJSON(&success); err != nil {
			return nil, normalizeError(err)
		}
		return &success, nil
	case http.StatusBadRequest:
		return nil, newResponseErrorFromSnapshot(resp, "review AIP deletion request bad request")
	case http.StatusForbidden:
		return nil, newResponseErrorFromSnapshot(resp, "review AIP deletion request forbidden")
	case http.StatusNotFound:
		return nil, newResponseErrorFromSnapshot(resp, "review AIP deletion request not found")
	default:
		return nil, newResponseErrorFromSnapshot(resp, fmt.Sprintf("unexpected review AIP deletion response status %d", resp.StatusCode))
	}
}

func hasJSONField(body []byte, field string) bool {
	if len(body) == 0 {
		return false
	}

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(body, &payload); err != nil {
		return false
	}

	_, ok := payload[field]
	return ok
}
