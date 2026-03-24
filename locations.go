package ssclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	kabs "github.com/microsoft/kiota-abstractions-go"

	kapi "go.artefactual.dev/ssclient/kiota/api"
	"go.artefactual.dev/ssclient/kiota/models"
)

// LocationsService provides location-related API operations.
type LocationsService struct {
	client *Client
}

// ListLocationsQuery describes the supported query parameters for listing
// locations.
type ListLocationsQuery struct {
	Description  *string
	Limit        *int32
	Offset       *int32
	OrderBy      *string
	PipelineID   *uuid.UUID
	Purpose      *models.LocationPurpose
	Quota        *int32
	RelativePath *string
	Used         *int32
	ID           *uuid.UUID
}

// List returns a filtered list of locations.
func (s *LocationsService) List(ctx context.Context, query ListLocationsQuery) (*models.LocationList, error) {
	reqConfig := &kabs.RequestConfiguration[kapi.V2LocationEmptyPathSegmentRequestBuilderGetQueryParameters]{
		QueryParameters: &kapi.V2LocationEmptyPathSegmentRequestBuilderGetQueryParameters{
			Description:    query.Description,
			Limit:          query.Limit,
			Offset:         query.Offset,
			Order_by:       query.OrderBy,
			Pipeline__uuid: query.PipelineID,
			Purpose:        query.Purpose,
			Quota:          query.Quota,
			Relative_path:  query.RelativePath,
			Used:           query.Used,
			Uuid:           query.ID,
		},
	}

	res, err := s.client.raw.Api().V2().Location().EmptyPathSegment().Get(ctx, reqConfig)
	if err != nil {
		return nil, normalizeError(err)
	}
	if res == nil {
		return nil, nil
	}

	typed, ok := res.(*models.LocationList)
	if !ok {
		return nil, fmt.Errorf("unexpected location list type %T", res)
	}

	return typed, nil
}

// Get returns a location by ID.
func (s *LocationsService) Get(ctx context.Context, id uuid.UUID) (*models.Location, error) {
	res, err := s.client.raw.Api().V2().Location().ByUuid(id).EmptyPathSegment().Get(ctx, nil)
	if err != nil {
		return nil, normalizeError(err)
	}
	if res == nil {
		return nil, nil
	}

	typed, ok := res.(*models.Location)
	if !ok {
		return nil, fmt.Errorf("unexpected location type %T", res)
	}

	return typed, nil
}

// Default returns the Location header from the default-location redirect for a
// given purpose. The returned value is a resource URI; use [ParseResourceURI]
// to extract the resource identifier when needed.
func (s *LocationsService) Default(ctx context.Context, purpose models.LocationPurpose) (string, error) {
	builder := s.client.raw.Api().V2().Location().DefaultEscaped().ByPurpose(purpose.String()).EmptyPathSegment()
	requestInfo, err := builder.ToGetRequestInformation(ctx, nil)
	if err != nil {
		return "", normalizeError(err)
	}

	resp, err := s.client.execute(ctx, requestInfo)
	if err != nil {
		return "", normalizeError(err)
	}

	if resp.StatusCode != http.StatusFound {
		return "", newResponseErrorFromSnapshot(resp, fmt.Sprintf("unexpected default location response status %d", resp.StatusCode))
	}

	location := resp.Headers.Get("Location")
	if location == "" {
		return "", fmt.Errorf("missing Location header")
	}

	return location, nil
}

// Move moves files to the specified location.
func (s *LocationsService) Move(ctx context.Context, id uuid.UUID, body *models.MoveRequest) error {
	if body == nil {
		return fmt.Errorf("location move request is required")
	}

	_, err := s.client.raw.Api().V2().Location().ByUuid(id).EmptyPathSegment().Post(ctx, body, nil)
	return normalizeError(err)
}
