package ssclient

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	kabs "github.com/microsoft/kiota-abstractions-go"

	kapi "go.artefactual.dev/ssclient/kiota/api"
	"go.artefactual.dev/ssclient/kiota/models"
)

// PipelinesService provides pipeline-related API operations.
type PipelinesService struct {
	client *Client
}

// ListPipelinesQuery describes the supported query parameters for listing
// pipelines.
type ListPipelinesQuery struct {
	Description *string
	ID          *uuid.UUID
}

// List returns a filtered list of pipelines.
func (s *PipelinesService) List(ctx context.Context, query ListPipelinesQuery) (*models.PipelineList, error) {
	reqConfig := &kabs.RequestConfiguration[kapi.V2PipelineEmptyPathSegmentRequestBuilderGetQueryParameters]{
		QueryParameters: &kapi.V2PipelineEmptyPathSegmentRequestBuilderGetQueryParameters{
			Description: query.Description,
			Uuid:        query.ID,
		},
	}

	res, err := s.client.raw.Api().V2().Pipeline().EmptyPathSegment().Get(ctx, reqConfig)
	if err != nil {
		return nil, normalizeError(err)
	}
	if res == nil {
		return nil, nil
	}

	typed, ok := res.(*models.PipelineList)
	if !ok {
		return nil, fmt.Errorf("unexpected pipeline list type %T", res)
	}

	return typed, nil
}

// Get returns a pipeline by ID.
func (s *PipelinesService) Get(ctx context.Context, id uuid.UUID) (*models.Pipeline, error) {
	res, err := s.client.raw.Api().V2().Pipeline().ByUuid(id).EmptyPathSegment().Get(ctx, nil)
	if err != nil {
		return nil, normalizeError(err)
	}
	if res == nil {
		return nil, nil
	}

	typed, ok := res.(*models.Pipeline)
	if !ok {
		return nil, fmt.Errorf("unexpected pipeline type %T", res)
	}

	return typed, nil
}
