package ssclient_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	kabs "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"

	"go.artefactual.dev/ssclient"
	"go.artefactual.dev/ssclient/kiota/models"
)

func closeBody(tb testing.TB, body io.ReadCloser) {
	tb.Helper()
	if body == nil {
		return
	}
	_ = body.Close()
}

func TestPackages(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"

		client, err := ssclient.New(ssclient.Config{
			BaseURL:    "http://storage.service",
			Username:   "test",
			Key:        "test",
			HTTPClient: &http.Client{},
		})
		assertEqual(t, err, nil)

		raw := client.Raw()
		raw.RequestAdapter = &fakeRequestAdapter{
			baseURL: "http://storage.service",
			send: func(ctx context.Context, requestInfo *kabs.RequestInformation, constructor serialization.ParsableFactory, errorMappings kabs.ErrorMappings) (serialization.Parsable, error) {
				if requestInfo == nil {
					t.Fatal("expected request info")
				}

				assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/file/{uuid}/")
				assertEqual(t, requestInfo.PathParameters["uuid"], packageID)
				assertEqual(t, requestInfo.Headers.Get("Accept"), []string{"application/json"})

				pkg := models.NewPackageEscaped()
				pkg.SetUuid(ptr(uuid.MustParse(packageID)))
				pkg.SetStatus(ptr("UPLOADED"))
				pkg.SetCurrentLocation(ptr("/api/v2/location/154660b9-b4a3-4886-8d68-5e170c0923b8/"))
				pkg.SetReplicas([]string{
					"/api/v2/file/96922350-ccde-4fb0-a999-d2010522028f/",
					"/api/v2/file/610bc407-ba6c-4dcd-8675-d2727a9aab18/",
				})

				return pkg, nil
			},
		}

		pkg, err := client.Packages().Get(context.Background(), uuid.MustParse(packageID))
		assertEqual(t, err, nil)
		if pkg == nil {
			t.Fatal("expected package")
		}

		assertEqual(t, pkg.GetUuid().String(), packageID)
		assertEqual(t, *pkg.GetStatus(), "UPLOADED")
		assertEqual(t, len(pkg.GetReplicas()), 2)
		assertEqual(t, *pkg.GetCurrentLocation(), "/api/v2/location/154660b9-b4a3-4886-8d68-5e170c0923b8/")
	})

	t.Run("DownloadPackage", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"

			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						assertEqual(t, r.Method, http.MethodGet)
						assertEqual(t, r.URL.String(), "http://storage.service/api/v2/file/"+packageID+"/download/")
						assertEqual(t, r.Header.Get("Accept"), "*/*")

						return &http.Response{
							StatusCode: http.StatusOK,
							Header: http.Header{
								"Content-Type":        {"application/zip"},
								"Content-Disposition": {`attachment; filename="working_bag.zip"`},
								"Content-Length":      {"9"},
							},
							Body:    io.NopCloser(strings.NewReader("zip-bytes")),
							Request: r,
						}, nil
					}),
				},
			})
			assertEqual(t, err, nil)

			raw := client.Raw()
			raw.RequestAdapter = &fakeRequestAdapter{
				baseURL: "http://storage.service",
				convertToNativeRequest: func(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
					assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/file/{uuid}/download/")
					assertEqual(t, requestInfo.PathParameters["uuid"], packageID)
					assertEqual(t, requestInfo.Headers.Get("Accept"), []string{"*/*"})

					return &http.Request{
						Method: http.MethodGet,
						URL: &url.URL{
							Scheme: "http",
							Host:   "storage.service",
							Path:   "/api/v2/file/" + packageID + "/download/",
						},
						Header: http.Header{
							"Accept": {"*/*"},
						},
					}, nil
				},
			}

			res, err := client.Packages().DownloadPackage(context.Background(), uuid.MustParse(packageID))
			assertEqual(t, err, nil)
			if res == nil {
				t.Fatal("expected download result")
			}
			defer closeBody(t, res.Body)

			body, err := io.ReadAll(res.Body)
			assertEqual(t, err, nil)

			assertEqual(t, res.StatusCode, http.StatusOK)
			assertEqual(t, res.ContentType, "application/zip")
			assertEqual(t, res.ContentLength, int64(9))
			assertEqual(t, res.ContentDisposition, `attachment; filename="working_bag.zip"`)
			assertEqual(t, res.Filename, "working_bag.zip")
			assertEqual(t, string(body), "zip-bytes")
		})

		t.Run("Unavailable", func(t *testing.T) {
			const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"

			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusAccepted,
							Header:     http.Header{"Content-Type": {"application/json"}},
							Body:       io.NopCloser(strings.NewReader(`{"message":"File is not locally available."}`)),
							Request:    r,
						}, nil
					}),
				},
			})
			assertEqual(t, err, nil)

			raw := client.Raw()
			raw.RequestAdapter = &fakeRequestAdapter{
				baseURL: "http://storage.service",
				convertToNativeRequest: func(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
					return &http.Request{
						Method: http.MethodGet,
						URL: &url.URL{
							Scheme: "http",
							Host:   "storage.service",
							Path:   "/api/v2/file/" + packageID + "/download/",
						},
						Header: http.Header{
							"Accept": {"*/*"},
						},
					}, nil
				},
			}

			res, err := client.Packages().DownloadPackage(context.Background(), uuid.MustParse(packageID))
			if err == nil {
				t.Fatal("expected download error")
			}
			if res != nil {
				t.Fatal("did not expect download result")
			}

			status, ok := ssclient.StatusCode(err)
			if !ok {
				t.Fatal("expected status code")
			}
			assertEqual(t, status, http.StatusAccepted)

			var unavailableErr *ssclient.NotAvailableError
			if !errors.As(err, &unavailableErr) {
				t.Fatalf("expected NotAvailableError, got %T", err)
			}
			assertEqual(t, unavailableErr.Message, "File is not locally available.")
			if !strings.Contains(err.Error(), "File is not locally available.") {
				t.Fatalf("unexpected error %v", err)
			}
		})
	})

	t.Run("DownloadFile", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"
			const relativePath = "working_bag/data/test.txt"

			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						assertEqual(t, r.Method, http.MethodGet)
						assertEqual(t, r.URL.String(), "http://storage.service/api/v2/file/"+packageID+"/extract_file/?relative_path_to_file=working_bag%2Fdata%2Ftest.txt")
						assertEqual(t, r.Header.Get("Accept"), "*/*")

						return &http.Response{
							StatusCode: http.StatusOK,
							Header: http.Header{
								"Content-Type":        {"text/plain"},
								"Content-Disposition": {`attachment; filename="test.txt"`},
								"Content-Length":      {"4"},
							},
							Body:    io.NopCloser(strings.NewReader("test")),
							Request: r,
						}, nil
					}),
				},
			})
			assertEqual(t, err, nil)

			raw := client.Raw()
			raw.RequestAdapter = &fakeRequestAdapter{
				baseURL: "http://storage.service",
				convertToNativeRequest: func(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
					assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/file/{uuid}/extract_file/?relative_path_to_file={relative_path_to_file}")
					assertEqual(t, requestInfo.PathParameters["uuid"], packageID)
					assertEqual(t, requestInfo.QueryParameters["relative_path_to_file"], relativePath)
					assertEqual(t, requestInfo.Headers.Get("Accept"), []string{"*/*"})

					return &http.Request{
						Method: http.MethodGet,
						URL: &url.URL{
							Scheme:   "http",
							Host:     "storage.service",
							Path:     "/api/v2/file/" + packageID + "/extract_file/",
							RawQuery: "relative_path_to_file=working_bag%2Fdata%2Ftest.txt",
						},
						Header: http.Header{
							"Accept": {"*/*"},
						},
					}, nil
				},
			}

			res, err := client.Packages().DownloadFile(context.Background(), uuid.MustParse(packageID), relativePath)
			assertEqual(t, err, nil)
			if res == nil {
				t.Fatal("expected extract file result")
			}
			defer closeBody(t, res.Body)

			body, err := io.ReadAll(res.Body)
			assertEqual(t, err, nil)

			assertEqual(t, res.StatusCode, http.StatusOK)
			assertEqual(t, res.ContentType, "text/plain")
			assertEqual(t, res.ContentLength, int64(4))
			assertEqual(t, res.ContentDisposition, `attachment; filename="test.txt"`)
			assertEqual(t, res.Filename, "test.txt")
			assertEqual(t, string(body), "test")
		})

		t.Run("EmptyPath", func(t *testing.T) {
			client, err := ssclient.New(ssclient.Config{
				BaseURL:    "http://storage.service",
				Username:   "test",
				Key:        "test",
				HTTPClient: &http.Client{},
			})
			assertEqual(t, err, nil)

			if _, err := client.Packages().DownloadFile(context.Background(), uuid.Nil, ""); err == nil {
				t.Fatal("expected error")
			}
		})
	})

	t.Run("DownloadPointerFile", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"

			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						assertEqual(t, r.Method, http.MethodGet)
						assertEqual(t, r.URL.String(), "http://storage.service/api/v2/file/"+packageID+"/pointer_file/")
						assertEqual(t, r.Header.Get("Accept"), "*/*")

						return &http.Response{
							StatusCode: http.StatusOK,
							Header: http.Header{
								"Content-Type":        {"application/xml"},
								"Content-Disposition": {`attachment; filename="pointer.xml"`},
								"Content-Length":      {"11"},
							},
							Body:    io.NopCloser(strings.NewReader("<pointer/>")),
							Request: r,
						}, nil
					}),
				},
			})
			assertEqual(t, err, nil)

			raw := client.Raw()
			raw.RequestAdapter = &fakeRequestAdapter{
				baseURL: "http://storage.service",
				convertToNativeRequest: func(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
					assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/file/{uuid}/pointer_file/")
					assertEqual(t, requestInfo.PathParameters["uuid"], packageID)
					assertEqual(t, requestInfo.Headers.Get("Accept"), []string{"*/*"})

					return &http.Request{
						Method: http.MethodGet,
						URL: &url.URL{
							Scheme: "http",
							Host:   "storage.service",
							Path:   "/api/v2/file/" + packageID + "/pointer_file/",
						},
						Header: http.Header{
							"Accept": {"*/*"},
						},
					}, nil
				},
			}

			res, err := client.Packages().DownloadPointerFile(context.Background(), uuid.MustParse(packageID))
			assertEqual(t, err, nil)
			if res == nil {
				t.Fatal("expected pointer file result")
			}
			defer closeBody(t, res.Body)

			body, err := io.ReadAll(res.Body)
			assertEqual(t, err, nil)

			assertEqual(t, res.StatusCode, http.StatusOK)
			assertEqual(t, res.ContentType, "application/xml")
			assertEqual(t, res.ContentLength, int64(11))
			assertEqual(t, res.ContentDisposition, `attachment; filename="pointer.xml"`)
			assertEqual(t, res.Filename, "pointer.xml")
			assertEqual(t, string(body), "<pointer/>")
		})

		t.Run("NotFound", func(t *testing.T) {
			const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"

			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusNotFound,
							Header:     http.Header{"Content-Type": {"text/plain; charset=utf-8"}},
							Body:       io.NopCloser(strings.NewReader("Resource with UUID does not have a pointer file")),
							Request:    r,
						}, nil
					}),
				},
			})
			assertEqual(t, err, nil)

			raw := client.Raw()
			raw.RequestAdapter = &fakeRequestAdapter{
				baseURL: "http://storage.service",
				convertToNativeRequest: func(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
					return &http.Request{
						Method: http.MethodGet,
						URL: &url.URL{
							Scheme: "http",
							Host:   "storage.service",
							Path:   "/api/v2/file/" + packageID + "/pointer_file/",
						},
						Header: http.Header{
							"Accept": {"*/*"},
						},
					}, nil
				},
			}

			res, err := client.Packages().DownloadPointerFile(context.Background(), uuid.MustParse(packageID))
			if err == nil {
				t.Fatal("expected pointer file error")
			}
			if res != nil {
				t.Fatal("did not expect pointer file result")
			}

			status, ok := ssclient.StatusCode(err)
			if !ok {
				t.Fatal("expected status code")
			}
			assertEqual(t, status, http.StatusNotFound)
			if !strings.Contains(err.Error(), "pointer file") {
				t.Fatalf("unexpected error %v", err)
			}
		})
	})

	t.Run("DeleteAIP", func(t *testing.T) {
		t.Run("Accepted", func(t *testing.T) {
			const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"

			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						assertEqual(t, r.Method, http.MethodPost)
						assertEqual(t, r.URL.String(), "http://storage.service/api/v2/file/"+packageID+"/delete_aip/")
						assertEqual(t, r.Header.Get("Accept"), "application/json")
						assertEqual(t, r.Header.Get("Authorization"), "ApiKey test:test")
						assertEqual(t, r.Header.Get("Content-Type"), "application/json")

						return &http.Response{
							StatusCode: http.StatusAccepted,
							Header:     http.Header{"Content-Type": {"application/json"}},
							Body:       io.NopCloser(strings.NewReader(`{"message":"Delete request created successfully.","id":17}`)),
							Request:    r,
						}, nil
					}),
				},
			})
			assertEqual(t, err, nil)

			raw := client.Raw()
			writerFactory := raw.RequestAdapter.GetSerializationWriterFactory()
			raw.RequestAdapter = &fakeRequestAdapter{
				serializationWriterFactory: writerFactory,
				convertToNativeRequest: func(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
					assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/file/{uuid}/delete_aip/")
					assertEqual(t, requestInfo.PathParameters["uuid"], packageID)
					assertEqual(t, requestInfo.Headers.Get("Accept"), []string{"application/json"})
					assertEqual(t, requestInfo.Headers.Get("Content-Type"), []string{"application/json"})

					got := string(requestInfo.Content)
					if !strings.Contains(got, `"event_reason":"Delete please!"`) {
						t.Fatalf("unexpected request body %q", got)
					}
					if !strings.Contains(got, `"pipeline":"4b9e8af5-b0af-4abf-80b8-4b7d76281f61"`) {
						t.Fatalf("unexpected request body %q", got)
					}
					if !strings.Contains(got, `"user_id":1`) {
						t.Fatalf("unexpected request body %q", got)
					}
					if !strings.Contains(got, `"user_email":"user@example.com"`) {
						t.Fatalf("unexpected request body %q", got)
					}

					return &http.Request{
						Method: http.MethodPost,
						URL: &url.URL{
							Scheme: "http",
							Host:   "storage.service",
							Path:   "/api/v2/file/" + packageID + "/delete_aip/",
						},
						Header: http.Header{
							"Accept":        {"application/json"},
							"Authorization": {"ApiKey test:test"},
							"Content-Type":  {"application/json"},
						},
						Body: io.NopCloser(strings.NewReader(got)),
					}, nil
				},
			}

			body := models.NewDeleteAipRequest()
			body.SetEventReason(ptr("Delete please!"))
			body.SetPipeline(ptr(uuid.MustParse("4b9e8af5-b0af-4abf-80b8-4b7d76281f61")))
			body.SetUserId(ptr(int32(1)))
			body.SetUserEmail(ptr("user@example.com"))

			res, err := client.Packages().DeleteAIP(context.Background(), uuid.MustParse(packageID), body)
			assertEqual(t, err, nil)
			if res == nil {
				t.Fatal("expected delete AIP result")
			}
			if !res.IsAccepted() {
				t.Fatal("expected accepted delete AIP result")
			}
			if res.HasExistingRequest() {
				t.Fatal("did not expect existing request result")
			}
			assertEqual(t, res.Accepted.Message, "Delete request created successfully.")
			assertEqual(t, res.Accepted.ID, int32(17))
		})

		t.Run("AlreadyExists", func(t *testing.T) {
			const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"

			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Header:     http.Header{"Content-Type": {"application/json"}},
							Body:       io.NopCloser(strings.NewReader(`{"error_message":"A deletion request already exists for this AIP."}`)),
							Request:    r,
						}, nil
					}),
				},
			})
			assertEqual(t, err, nil)

			raw := client.Raw()
			writerFactory := raw.RequestAdapter.GetSerializationWriterFactory()
			raw.RequestAdapter = &fakeRequestAdapter{
				serializationWriterFactory: writerFactory,
				convertToNativeRequest: func(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
					return &http.Request{
						Method: http.MethodPost,
						URL: &url.URL{
							Scheme: "http",
							Host:   "storage.service",
							Path:   "/api/v2/file/" + packageID + "/delete_aip/",
						},
						Header: http.Header{
							"Accept":        {"application/json"},
							"Authorization": {"ApiKey test:test"},
							"Content-Type":  {"application/json"},
						},
						Body: io.NopCloser(strings.NewReader(string(requestInfo.Content))),
					}, nil
				},
			}

			body := models.NewDeleteAipRequest()
			body.SetEventReason(ptr("Delete please!"))
			body.SetPipeline(ptr(uuid.MustParse("4b9e8af5-b0af-4abf-80b8-4b7d76281f61")))
			body.SetUserId(ptr(int32(1)))
			body.SetUserEmail(ptr("user@example.com"))

			res, err := client.Packages().DeleteAIP(context.Background(), uuid.MustParse(packageID), body)
			assertEqual(t, err, nil)
			if res == nil {
				t.Fatal("expected delete AIP result")
			}
			if res.IsAccepted() {
				t.Fatal("did not expect accepted delete AIP result")
			}
			if !res.HasExistingRequest() {
				t.Fatal("expected existing request result")
			}
			assertEqual(t, res.AlreadyExists.ErrorMessage, "A deletion request already exists for this AIP.")
		})

		t.Run("NilBody", func(t *testing.T) {
			client, err := ssclient.New(ssclient.Config{
				BaseURL:    "http://storage.service",
				Username:   "test",
				Key:        "test",
				HTTPClient: &http.Client{},
			})
			assertEqual(t, err, nil)

			if _, err := client.Packages().DeleteAIP(context.Background(), uuid.Nil, nil); err == nil {
				t.Fatal("expected error")
			}
		})
	})

	t.Run("CheckFixity", func(t *testing.T) {
		const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"

		client, err := ssclient.New(ssclient.Config{
			BaseURL:    "http://storage.service",
			Username:   "test",
			Key:        "test",
			HTTPClient: &http.Client{},
		})
		assertEqual(t, err, nil)

		raw := client.Raw()
		forceLocal := true
		raw.RequestAdapter = &fakeRequestAdapter{
			baseURL: "http://storage.service",
			send: func(ctx context.Context, requestInfo *kabs.RequestInformation, constructor serialization.ParsableFactory, errorMappings kabs.ErrorMappings) (serialization.Parsable, error) {
				assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/file/{uuid}/check_fixity/{?force_local}")
				assertEqual(t, requestInfo.PathParameters["uuid"], packageID)
				assertEqual(t, requestInfo.QueryParameters["force_local"], "true")
				assertEqual(t, requestInfo.Headers.Get("Accept"), []string{"application/json"})

				response := models.NewFixityResponse()
				response.SetSuccess(ptr(true))
				response.SetMessage(ptr("ok"))
				return response, nil
			},
		}

		res, err := client.Packages().CheckFixity(context.Background(), uuid.MustParse(packageID), ssclient.CheckFixityOptions{
			ForceLocal: &forceLocal,
		})
		assertEqual(t, err, nil)
		if res == nil {
			t.Fatal("expected fixity response")
		}
		assertEqual(t, *res.GetSuccess(), true)
		assertEqual(t, *res.GetMessage(), "ok")
	})

	t.Run("Move", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"
			const locationID = "154660b9-b4a3-4886-8d68-5e170c0923b8"

			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						assertEqual(t, r.Method, http.MethodPost)
						assertEqual(t, r.URL.String(), "http://storage.service/api/v2/file/"+packageID+"/move/")
						assertEqual(t, r.Header.Get("Accept"), "application/json")
						assertEqual(t, r.Header.Get("Content-Type"), "application/x-www-form-urlencoded")

						body, err := io.ReadAll(r.Body)
						assertEqual(t, err, nil)
						assertEqual(t, string(body), "location_uuid="+url.QueryEscape(locationID))

						return &http.Response{
							StatusCode: http.StatusAccepted,
							Header:     http.Header{"Location": {"/api/v2/async/1/"}},
							Body:       io.NopCloser(strings.NewReader("")),
							Request:    r,
						}, nil
					}),
				},
			})
			assertEqual(t, err, nil)

			raw := client.Raw()
			raw.RequestAdapter = &fakeRequestAdapter{
				baseURL: "http://storage.service",
				convertToNativeRequest: func(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
					assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/file/{uuid}/move/")
					assertEqual(t, requestInfo.PathParameters["uuid"], packageID)
					assertEqual(t, requestInfo.Headers.Get("Accept"), []string{"application/json"})
					assertEqual(t, requestInfo.Headers.Get("Content-Type"), []string{"application/x-www-form-urlencoded"})
					assertEqual(t, string(requestInfo.Content), "location_uuid="+url.QueryEscape(locationID))

					return &http.Request{
						Method: http.MethodPost,
						URL: &url.URL{
							Scheme: "http",
							Host:   "storage.service",
							Path:   "/api/v2/file/" + packageID + "/move/",
						},
						Header: http.Header{
							"Accept":       {"application/json"},
							"Content-Type": {"application/x-www-form-urlencoded"},
						},
						Body: io.NopCloser(strings.NewReader(string(requestInfo.Content))),
					}, nil
				},
			}

			res, err := client.Packages().Move(context.Background(), uuid.MustParse(packageID), uuid.MustParse(locationID))
			assertEqual(t, err, nil)
			if res == nil {
				t.Fatal("expected move result")
			}
			assertEqual(t, res.Location, "/api/v2/async/1/")
		})

		t.Run("MissingLocationHeader", func(t *testing.T) {
			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusAccepted,
							Header:     http.Header{},
							Body:       io.NopCloser(strings.NewReader("")),
							Request:    r,
						}, nil
					}),
				},
			})
			assertEqual(t, err, nil)

			raw := client.Raw()
			raw.RequestAdapter = &fakeRequestAdapter{
				baseURL: "http://storage.service",
				convertToNativeRequest: func(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
					return &http.Request{
						Method: http.MethodPost,
						URL: &url.URL{
							Scheme: "http",
							Host:   "storage.service",
							Path:   "/api/v2/file/" + uuid.Nil.String() + "/move/",
						},
						Header: http.Header{
							"Accept":       {"application/json"},
							"Content-Type": {"application/x-www-form-urlencoded"},
						},
						Body: io.NopCloser(strings.NewReader(string(requestInfo.Content))),
					}, nil
				},
			}

			if _, err := client.Packages().Move(context.Background(), uuid.Nil, uuid.Nil); err == nil {
				t.Fatal("expected error")
			}
		})
	})

	t.Run("ReviewAIPDeletion", func(t *testing.T) {
		const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"

		newReviewAIPDeletionRequestAdapter := func(t *testing.T) (*fakeRequestAdapter, *models.ReviewAipDeletionRequest) {
			t.Helper()

			client, err := ssclient.New(ssclient.Config{
				BaseURL:    "http://storage.service",
				Username:   "test",
				Key:        "test",
				HTTPClient: &http.Client{},
			})
			assertEqual(t, err, nil)

			raw := client.Raw()
			writerFactory := raw.RequestAdapter.GetSerializationWriterFactory()
			adapter := &fakeRequestAdapter{
				baseURL:                    "http://storage.service",
				serializationWriterFactory: writerFactory,
				convertToNativeRequest: func(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
					assertEqual(t, requestInfo.UrlTemplate, "{+baseurl}/api/v2/file/{uuid}/review_aip_deletion/")
					assertEqual(t, requestInfo.PathParameters["uuid"], packageID)
					assertEqual(t, requestInfo.Headers.Get("Accept"), []string{"application/json"})
					assertEqual(t, requestInfo.Headers.Get("Content-Type"), []string{"application/json"})

					got := string(requestInfo.Content)
					if !strings.Contains(got, `"event_id":99`) {
						t.Fatalf("unexpected request body %q", got)
					}
					if !strings.Contains(got, `"decision":"approve"`) {
						t.Fatalf("unexpected request body %q", got)
					}
					if !strings.Contains(got, `"reason":"approved by workflow"`) {
						t.Fatalf("unexpected request body %q", got)
					}

					return &http.Request{
						Method: http.MethodPost,
						URL: &url.URL{
							Scheme: "http",
							Host:   "storage.service",
							Path:   "/api/v2/file/" + packageID + "/review_aip_deletion/",
						},
						Header: http.Header{
							"Accept":       {"application/json"},
							"Content-Type": {"application/json"},
						},
						Body: io.NopCloser(strings.NewReader(got)),
					}, nil
				},
			}

			body := models.NewReviewAipDeletionRequest()
			body.SetEventId(ptr(int32(99)))
			decision := models.APPROVE_REVIEWAIPDELETIONDECISION
			body.SetDecision(&decision)
			body.SetReason(ptr("approved by workflow"))

			return adapter, body
		}

		t.Run("Success", func(t *testing.T) {
			requestAdapter, body := newReviewAIPDeletionRequestAdapter(t)

			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Header:     http.Header{"Content-Type": {"application/json"}},
							Body:       io.NopCloser(strings.NewReader(`{"message":"done"}`)),
							Request:    r,
						}, nil
					}),
				},
			})
			assertEqual(t, err, nil)
			client.Raw().RequestAdapter = requestAdapter

			res, err := client.Packages().ReviewAIPDeletion(context.Background(), uuid.MustParse(packageID), body)
			assertEqual(t, err, nil)
			if res == nil {
				t.Fatal("expected review deletion response")
			}
			assertEqual(t, res.Message, "done")
		})

		t.Run("SuccessWithDetail", func(t *testing.T) {
			requestAdapter, body := newReviewAIPDeletionRequestAdapter(t)

			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Header:     http.Header{"Content-Type": {"application/json"}},
							Body:       io.NopCloser(strings.NewReader(`{"message":"done","detail":"LOCKSS warning"}`)),
							Request:    r,
						}, nil
					}),
				},
			})
			assertEqual(t, err, nil)
			client.Raw().RequestAdapter = requestAdapter

			res, err := client.Packages().ReviewAIPDeletion(context.Background(), uuid.MustParse(packageID), body)
			assertEqual(t, err, nil)
			if res == nil {
				t.Fatal("expected review deletion response")
			}
			assertEqual(t, res.Message, "done")
			assertEqual(t, res.Detail, "LOCKSS warning")
		})

		t.Run("BusinessFailure", func(t *testing.T) {
			const packageID = "7c8a3549-2fe0-41d3-9d83-f485f1a43be3"

			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Header:     http.Header{"Content-Type": {"application/json"}},
							Body:       io.NopCloser(strings.NewReader(`{"error_message":"disk error"}`)),
							Request:    r,
						}, nil
					}),
				},
			})
			assertEqual(t, err, nil)

			raw := client.Raw()
			writerFactory := raw.RequestAdapter.GetSerializationWriterFactory()
			raw.RequestAdapter = &fakeRequestAdapter{
				baseURL:                    "http://storage.service",
				serializationWriterFactory: writerFactory,
				convertToNativeRequest: func(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
					return &http.Request{
						Method: http.MethodPost,
						URL: &url.URL{
							Scheme: "http",
							Host:   "storage.service",
							Path:   "/api/v2/file/" + packageID + "/review_aip_deletion/",
						},
						Header: http.Header{
							"Accept":       {"application/json"},
							"Content-Type": {"application/json"},
						},
						Body: io.NopCloser(strings.NewReader(string(requestInfo.Content))),
					}, nil
				},
			}

			body := models.NewReviewAipDeletionRequest()
			body.SetEventId(ptr(int32(99)))
			decision := models.APPROVE_REVIEWAIPDELETIONDECISION
			body.SetDecision(&decision)
			body.SetReason(ptr("approved by workflow"))

			res, err := client.Packages().ReviewAIPDeletion(context.Background(), uuid.MustParse(packageID), body)
			if err == nil {
				t.Fatal("expected review deletion error")
			}
			if res != nil {
				t.Fatalf("did not expect review deletion result %#v", res)
			}

			var reviewErr *ssclient.ReviewAIPDeletionError
			if !errors.As(err, &reviewErr) {
				t.Fatalf("expected ReviewAIPDeletionError, got %T", err)
			}
			assertEqual(t, reviewErr.ErrorMessage, "disk error")
			assertEqual(t, reviewErr.Detail, "")
			if !strings.Contains(err.Error(), "disk error") {
				t.Fatalf("unexpected error %v", err)
			}
		})

		t.Run("TransportError", func(t *testing.T) {
			requestAdapter, body := newReviewAIPDeletionRequestAdapter(t)

			client, err := ssclient.New(ssclient.Config{
				BaseURL:  "http://storage.service",
				Username: "test",
				Key:      "test",
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						return nil, errors.New("network down")
					}),
				},
			})
			assertEqual(t, err, nil)
			client.Raw().RequestAdapter = requestAdapter

			res, err := client.Packages().ReviewAIPDeletion(context.Background(), uuid.MustParse(packageID), body)
			if err == nil {
				t.Fatal("expected transport error")
			}
			if res != nil {
				t.Fatalf("did not expect review deletion result %#v", res)
			}

			var responseErr *ssclient.ResponseError
			if !errors.As(err, &responseErr) {
				t.Fatalf("expected ResponseError, got %T", err)
			}
			assertEqual(t, responseErr.StatusCode, 0)
			if !strings.Contains(responseErr.Message, "network down") {
				t.Fatalf("unexpected response error message %q", responseErr.Message)
			}
		})

		t.Run("NilBody", func(t *testing.T) {
			client, err := ssclient.New(ssclient.Config{
				BaseURL:    "http://storage.service",
				Username:   "test",
				Key:        "test",
				HTTPClient: &http.Client{},
			})
			assertEqual(t, err, nil)

			if _, err := client.Packages().ReviewAIPDeletion(context.Background(), uuid.Nil, nil); err == nil {
				t.Fatal("expected error")
			}
		})
	})
}
