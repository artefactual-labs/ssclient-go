package ssclient

import (
	"context"
	"fmt"
	"net/http"

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
	PipelineUUID *string
	Purpose      *models.LocationPurpose
	Quota        *int32
	RelativePath *string
	Used         *int32
	UUID         *string
}

// List returns a filtered list of locations.
func (s *LocationsService) List(ctx context.Context, query ListLocationsQuery) (*models.LocationList, error) {
	pipelineUUID, err := parseOptionalUUID(query.PipelineUUID)
	if err != nil {
		return nil, fmt.Errorf("invalid pipeline UUID %q: %w", *query.PipelineUUID, err)
	}

	locationUUID, err := parseOptionalUUID(query.UUID)
	if err != nil {
		return nil, fmt.Errorf("invalid location UUID %q: %w", *query.UUID, err)
	}

	reqConfig := &kapi.V2LocationEmptyPathSegmentRequestBuilderGetRequestConfiguration{
		QueryParameters: &kapi.V2LocationEmptyPathSegmentRequestBuilderGetQueryParameters{
			Description:              query.Description,
			Limit:                    query.Limit,
			Offset:                   query.Offset,
			Order_by:                 query.OrderBy,
			Pipeline__uuid:           pipelineUUID,
			PurposeAsLocationPurpose: query.Purpose,
			Quota:                    query.Quota,
			Relative_path:            query.RelativePath,
			Used:                     query.Used,
			Uuid:                     locationUUID,
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

// Get returns a location by UUID.
func (s *LocationsService) Get(ctx context.Context, uuid string) (*models.Location, error) {
	res, err := s.client.raw.Api().V2().Location().ByUuid(uuid).EmptyPathSegment().Get(ctx, nil)
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
// given purpose.
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
func (s *LocationsService) Move(ctx context.Context, uuid string, body *models.MoveRequest) error {
	if body == nil {
		return fmt.Errorf("location move request is required")
	}

	_, err := s.client.raw.Api().V2().Location().ByUuid(uuid).EmptyPathSegment().Post(ctx, body, nil)
	return normalizeError(err)
}
