package ssclient

import (
	"context"
	"fmt"
	"math"
	"mime"
	"time"

	"go.artefactual.dev/ssclient/kiota/models"
)

const defaultAsyncPollInterval = time.Second

// AsyncService provides access to asynchronous operation tracking.
type AsyncService struct {
	client *Client
}

type waitOptions struct {
	pollInterval time.Duration
}

// WaitOption customizes the polling behavior of [AsyncService.Wait].
type WaitOption func(*waitOptions) error

// WithPollInterval sets the delay between successive polls performed by
// [AsyncService.Wait].
func WithPollInterval(interval time.Duration) WaitOption {
	return func(opts *waitOptions) error {
		if interval <= 0 {
			return fmt.Errorf("poll interval must be greater than zero")
		}
		opts.pollInterval = interval
		return nil
	}
}

// Get returns the current state of an asynchronous task by ID.
func (s *AsyncService) Get(ctx context.Context, id int) (*models.AsyncTask, error) {
	taskID, err := asyncIDToInt32(id)
	if err != nil {
		return nil, err
	}

	requestInfo, err := s.client.raw.Api().V2().Async().ById(taskID).EmptyPathSegment().ToGetRequestInformation(ctx, nil)
	if err != nil {
		return nil, normalizeError(err)
	}
	resp, err := s.client.execute(ctx, requestInfo)
	if err != nil {
		return nil, normalizeError(err)
	}

	switch resp.StatusCode {
	case 200:
		return decodeAsyncTask(resp)
	case 404:
		return nil, newResponseErrorFromSnapshot(resp, "async task not found")
	default:
		return nil, newResponseErrorFromSnapshot(resp, fmt.Sprintf("unexpected async task response status %d", resp.StatusCode))
	}
}

// Wait polls the asynchronous task until it completes, the context is canceled,
// or a transport/protocol failure occurs.
func (s *AsyncService) Wait(ctx context.Context, id int, opts ...WaitOption) (*models.AsyncTask, error) {
	cfg := waitOptions{
		pollInterval: defaultAsyncPollInterval,
	}
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}

	timer := time.NewTimer(0)
	defer func() {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
	}()

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timer.C:
		}
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		task, err := s.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		if task == nil {
			return nil, fmt.Errorf("async task response is empty")
		}

		if task.GetCompleted() != nil && *task.GetCompleted() {
			if task.GetWasError() != nil && *task.GetWasError() {
				return task, newAsyncTaskError(task)
			}
			return task, nil
		}

		// Safe to reset without stopping because the timer has already fired and
		// this iteration consumed from timer.C before entering the poll body.
		timer.Reset(cfg.pollInterval)
	}
}

func asyncIDToInt32(id int) (int32, error) {
	if id < 0 || id > math.MaxInt32 {
		return 0, fmt.Errorf("async task ID %d is out of range", id)
	}
	return int32(id), nil
}

func decodeAsyncTask(resp *responseSnapshot) (*models.AsyncTask, error) {
	if resp == nil || len(resp.Body) == 0 {
		return nil, fmt.Errorf("response body is empty")
	}

	contentType := "application/json"
	if header := resp.Headers.Get("Content-Type"); header != "" {
		if mediaType, _, err := mime.ParseMediaType(header); err == nil && mediaType != "" {
			contentType = mediaType
		}
	}

	root, err := newTolerantParseNodeFactory().GetRootParseNode(contentType, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse async task response: %w", err)
	}

	value, err := root.GetObjectValue(models.CreateAsyncTaskFromDiscriminatorValue)
	if err != nil {
		return nil, fmt.Errorf("decode async task response: %w", err)
	}

	task, ok := value.(*models.AsyncTask)
	if !ok {
		return nil, fmt.Errorf("unexpected async task type %T", value)
	}

	return task, nil
}
