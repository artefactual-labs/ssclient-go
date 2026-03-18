package ssclient

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"

	kabs "github.com/microsoft/kiota-abstractions-go"
	kapi "go.artefactual.dev/ssclient/kiota/api"
	"go.artefactual.dev/ssclient/kiota/models"
)

// PackagesService provides package-related API operations.
type PackagesService struct {
	client *Client
}

// CheckFixityOptions configures a package fixity check request.
type CheckFixityOptions struct {
	ForceLocal *bool
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
	StatusCode    int
	Accepted      *DeleteAIPAccepted
	AlreadyExists *DeleteAIPAlreadyExists
}

// FileStream captures a successful streamed file response.
type FileStream struct {
	StatusCode    int
	ContentType   string
	ContentLength int64
	Filename      string
	Body          io.ReadCloser
}

type extractFileQuery struct {
	RelativePathToFile *string `uriparametername:"relative_path_to_file"`
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

// Get returns a package by UUID.
func (s *PackagesService) Get(ctx context.Context, uuid string) (*models.PackageEscaped, error) {
	pkg, err := s.client.raw.Api().V2().File().ByUuid(uuid).Get(ctx, nil)
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
func (s *PackagesService) DownloadPackage(ctx context.Context, uuid string) (*FileStream, error) {
	requestInfo, err := s.client.raw.Api().V2().File().ByUuid(uuid).ToGetRequestInformation(ctx, nil)
	if err != nil {
		return nil, normalizeError(err)
	}

	requestInfo.UrlTemplate = "{+baseurl}/api/v2/file/{uuid}/download/"
	requestInfo.Headers.Remove("Accept")
	requestInfo.Headers.Add("Accept", "*/*")

	return s.streamPackageRequest(ctx, requestInfo, "download")
}

// DownloadFile extracts a file from a package and streams it back.
func (s *PackagesService) DownloadFile(ctx context.Context, uuid string, relativePathToFile string) (*FileStream, error) {
	if relativePathToFile == "" {
		return nil, fmt.Errorf("relative path to file is required")
	}

	requestInfo := kabs.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(
		kabs.GET,
		"{+baseurl}/api/v2/file/{uuid}/extract_file/{?relative_path_to_file}",
		map[string]string{"uuid": uuid},
	)
	requestInfo.AddQueryParameters(extractFileQuery{
		RelativePathToFile: &relativePathToFile,
	})
	requestInfo.Headers.Add("Accept", "*/*")

	return s.streamPackageRequest(ctx, requestInfo, "extract file")
}

// DownloadPointerFile returns the package pointer file as a stream.
func (s *PackagesService) DownloadPointerFile(ctx context.Context, uuid string) (*FileStream, error) {
	requestInfo, err := s.client.raw.Api().V2().File().ByUuid(uuid).ToGetRequestInformation(ctx, nil)
	if err != nil {
		return nil, normalizeError(err)
	}

	requestInfo.UrlTemplate = "{+baseurl}/api/v2/file/{uuid}/pointer_file/"
	requestInfo.Headers.Remove("Accept")
	requestInfo.Headers.Add("Accept", "*/*")

	return s.streamPackageRequest(ctx, requestInfo, "pointer file")
}

func (s *PackagesService) streamPackageRequest(ctx context.Context, requestInfo *kabs.RequestInformation, action string) (*FileStream, error) {
	resp, err := s.client.executeStream(ctx, requestInfo)
	if err != nil {
		return nil, normalizeError(err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()

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
		StatusCode:    resp.StatusCode,
		ContentType:   resp.Headers.Get("Content-Type"),
		ContentLength: parseContentLength(resp.Headers.Get("Content-Length")),
		Filename:      parseFilename(resp.Headers.Get("Content-Disposition")),
		Body:          resp.Body,
	}, nil
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

// DeleteAIP creates an AIP deletion request for the given package. The server
// exposes two non-error outcomes: one for a newly created request and one for
// an already pending request.
func (s *PackagesService) DeleteAIP(ctx context.Context, uuid string, body *models.DeleteAipRequest) (*DeleteAIPResult, error) {
	if body == nil {
		return nil, fmt.Errorf("delete AIP request is required")
	}

	builder := s.client.raw.Api().V2().File().ByUuid(uuid).Delete_aip().EmptyPathSegment()
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
			StatusCode: http.StatusAccepted,
			Accepted:   &accepted,
		}, nil
	case http.StatusOK:
		var existing DeleteAIPAlreadyExists
		if err := resp.decodeJSON(&existing); err != nil {
			return nil, normalizeError(err)
		}
		return &DeleteAIPResult{
			StatusCode:    http.StatusOK,
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

// CheckFixity runs a fixity check for the given package UUID.
func (s *PackagesService) CheckFixity(ctx context.Context, uuid string, opts CheckFixityOptions) (*models.FixityResponse, error) {
	reqConfig := &kapi.V2FileItemCheck_fixityEmptyPathSegmentRequestBuilderGetRequestConfiguration{}
	if opts.ForceLocal != nil {
		reqConfig.QueryParameters = &kapi.V2FileItemCheck_fixityEmptyPathSegmentRequestBuilderGetQueryParameters{
			Force_local: opts.ForceLocal,
		}
	}

	res, err := s.client.raw.Api().V2().File().ByUuid(uuid).Check_fixity().EmptyPathSegment().Get(ctx, reqConfig)
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

// Move moves a package to a different storage location.
func (s *PackagesService) Move(ctx context.Context, uuid string, body *models.PackageMoveRequest) error {
	if body == nil {
		return fmt.Errorf("package move request is required")
	}

	_, err := s.client.raw.Api().V2().File().ByUuid(uuid).Move().EmptyPathSegment().Post(ctx, body, nil)
	return normalizeError(err)
}

// ReviewAIPDeletion approves or rejects an AIP deletion request associated with
// a package.
func (s *PackagesService) ReviewAIPDeletion(ctx context.Context, uuid string, body *models.ReviewAipDeletionRequest) (*models.ReviewAipDeletionSuccess, error) {
	if body == nil {
		return nil, fmt.Errorf("review AIP deletion request is required")
	}

	res, err := s.client.raw.Api().V2().File().ByUuid(uuid).Review_aip_deletion().EmptyPathSegment().Post(ctx, body, nil)
	if err != nil {
		return nil, normalizeError(err)
	}
	if res == nil {
		return nil, nil
	}

	typed, ok := res.(*models.ReviewAipDeletionSuccess)
	if !ok {
		return nil, fmt.Errorf("unexpected review AIP deletion response type %T", res)
	}

	return typed, nil
}
