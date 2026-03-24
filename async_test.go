package ssclient_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"go.artefactual.dev/ssclient"
)

func TestAsync(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						assertEqual(t, r.Method, http.MethodGet)
						assertEqual(t, r.URL.String(), "http://storage.service/api/v2/async/1/")
						return &http.Response{
							StatusCode: http.StatusOK,
							Header:     http.Header{"Content-Type": {"application/json"}},
							Body:       io.NopCloser(strings.NewReader(`{"id":1,"resource_uri":"/api/v2/async/1/","completed":true,"was_error":false,"created_time":"2026-03-24T10:00:00Z","updated_time":"2026-03-24T10:00:01Z","completed_time":"2026-03-24T10:00:01Z","result":"Package moved successfully"}`)),
							Request:    r,
						}, nil
					}),
				},
			})
			assertEqual(t, err, nil)

			task, err := client.Async().Get(context.Background(), 1)
			assertEqual(t, err, nil)
			if task == nil {
				t.Fatal("expected async task")
			}
			if task.GetId() == nil || *task.GetId() != 1 {
				t.Fatalf("unexpected task id %#v", task.GetId())
			}
			if task.GetCompleted() == nil || !*task.GetCompleted() {
				t.Fatalf("expected completed task %#v", task.GetCompleted())
			}
			if task.GetResult() == nil {
				t.Fatal("expected async result")
			}
		})

		t.Run("NotFound", func(t *testing.T) {
			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusNotFound,
							Header:     http.Header{"Content-Type": {"application/json"}},
							Body:       io.NopCloser(strings.NewReader(`{"error":"missing"}`)),
							Request:    r,
						}, nil
					}),
				},
			})
			assertEqual(t, err, nil)

			task, err := client.Async().Get(context.Background(), 1)
			if err == nil {
				t.Fatal("expected error")
			}
			if task != nil {
				t.Fatalf("did not expect async task %#v", task)
			}
			status, ok := ssclient.StatusCode(err)
			if !ok || status != http.StatusNotFound {
				t.Fatalf("unexpected status %d ok=%v err=%v", status, ok, err)
			}
		})
	})

	t.Run("Wait", func(t *testing.T) {
		t.Run("ImmediateSuccess", func(t *testing.T) {
			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Header:     http.Header{"Content-Type": {"application/json"}},
							Body:       io.NopCloser(strings.NewReader(`{"id":1,"resource_uri":"/api/v2/async/1/","completed":true,"was_error":false,"created_time":"2026-03-24T10:00:00Z","updated_time":"2026-03-24T10:00:01Z","completed_time":"2026-03-24T10:00:01Z","result":"Package moved successfully"}`)),
							Request:    r,
						}, nil
					}),
				},
			})
			assertEqual(t, err, nil)

			task, err := client.Async().Wait(context.Background(), 1)
			assertEqual(t, err, nil)
			if task == nil {
				t.Fatal("expected async task")
			}
		})

		t.Run("PendingThenSuccess", func(t *testing.T) {
			callCount := 0
			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						callCount++
						body := `{"id":1,"resource_uri":"/api/v2/async/1/","completed":false,"was_error":false,"created_time":"2026-03-24T10:00:00Z","updated_time":"2026-03-24T10:00:00Z","completed_time":null}`
						if callCount > 1 {
							body = `{"id":1,"resource_uri":"/api/v2/async/1/","completed":true,"was_error":false,"created_time":"2026-03-24T10:00:00Z","updated_time":"2026-03-24T10:00:01Z","completed_time":"2026-03-24T10:00:01Z","result":"Package moved successfully"}`
						}
						return &http.Response{
							StatusCode: http.StatusOK,
							Header:     http.Header{"Content-Type": {"application/json"}},
							Body:       io.NopCloser(strings.NewReader(body)),
							Request:    r,
						}, nil
					}),
				},
			})
			assertEqual(t, err, nil)

			task, err := client.Async().Wait(context.Background(), 1, ssclient.WithPollInterval(10*time.Millisecond))
			assertEqual(t, err, nil)
			if task == nil {
				t.Fatal("expected async task")
			}
			if callCount < 2 {
				t.Fatalf("expected multiple polls, got %d", callCount)
			}
		})

		t.Run("TerminalAsyncFailure", func(t *testing.T) {
			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Header:     http.Header{"Content-Type": {"application/json"}},
							Body:       io.NopCloser(strings.NewReader(`{"id":1,"resource_uri":"/api/v2/async/1/","completed":true,"was_error":true,"created_time":"2026-03-24T10:00:00Z","updated_time":"2026-03-24T10:00:01Z","completed_time":"2026-03-24T10:00:01Z","error":"task failed"}`)),
							Request:    r,
						}, nil
					}),
				},
			})
			assertEqual(t, err, nil)

			task, err := client.Async().Wait(context.Background(), 1)
			if err == nil {
				t.Fatal("expected error")
			}
			if task == nil {
				t.Fatal("expected terminal task")
			}
			var taskErr *ssclient.AsyncTaskError
			if !errors.As(err, &taskErr) {
				t.Fatalf("expected AsyncTaskError, got %T", err)
			}
			if taskErr.Task != task {
				t.Fatal("expected error to reference returned task")
			}
			assertEqual(t, taskErr.Message, "task failed")
		})

		t.Run("ContextCanceled", func(t *testing.T) {
			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Header:     http.Header{"Content-Type": {"application/json"}},
							Body:       io.NopCloser(strings.NewReader(`{"id":1,"resource_uri":"/api/v2/async/1/","completed":false,"was_error":false,"created_time":"2026-03-24T10:00:00Z","updated_time":"2026-03-24T10:00:00Z","completed_time":null}`)),
							Request:    r,
						}, nil
					}),
				},
			})
			assertEqual(t, err, nil)

			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
			defer cancel()

			task, err := client.Async().Wait(ctx, 1, ssclient.WithPollInterval(50*time.Millisecond))
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("expected deadline exceeded, got %v", err)
			}
			if task != nil {
				t.Fatalf("did not expect task %#v", task)
			}
		})

		t.Run("AlreadyCanceledContextDoesNotPoll", func(t *testing.T) {
			callCount := 0
			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						callCount++
						return &http.Response{
							StatusCode: http.StatusOK,
							Header:     http.Header{"Content-Type": {"application/json"}},
							Body:       io.NopCloser(strings.NewReader(`{"id":1,"resource_uri":"/api/v2/async/1/","completed":true,"was_error":false}`)),
							Request:    r,
						}, nil
					}),
				},
			})
			assertEqual(t, err, nil)

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			task, err := client.Async().Wait(ctx, 1)
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("expected context canceled, got %v", err)
			}
			if task != nil {
				t.Fatalf("did not expect task %#v", task)
			}
			if callCount != 0 {
				t.Fatalf("expected no polls after cancellation, got %d", callCount)
			}
		})

		t.Run("InvalidPollInterval", func(t *testing.T) {
			client, err := ssclient.New(ssclient.Config{
				BaseURL:    "http://storage.service",
				Username:   "test",
				Key:        "test",
				HTTPClient: &http.Client{},
			})
			assertEqual(t, err, nil)

			if _, err := client.Async().Wait(context.Background(), 1, ssclient.WithPollInterval(0)); err == nil {
				t.Fatal("expected error")
			}
		})
	})
}
