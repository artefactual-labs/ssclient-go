package ssclient

import (
	"context"
	"fmt"

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
	UUID        *string
}

// List returns a filtered list of pipelines.
func (s *PipelinesService) List(ctx context.Context, query ListPipelinesQuery) (*models.PipelineList, error) {
	reqConfig := &kapi.V2PipelineEmptyPathSegmentRequestBuilderGetRequestConfiguration{
		QueryParameters: &kapi.V2PipelineEmptyPathSegmentRequestBuilderGetQueryParameters{
			Description: query.Description,
			Uuid:        query.UUID,
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

// Get returns a pipeline by UUID.
func (s *PipelinesService) Get(ctx context.Context, uuid string) (*models.Pipeline, error) {
	res, err := s.client.raw.Api().V2().Pipeline().ByUuid(uuid).Get(ctx, nil)
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
