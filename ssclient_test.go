package ssclient_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	kabs "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoft/kiota-abstractions-go/store"

	"go.artefactual.dev/ssclient"
)

type fakeRequestAdapter struct {
	send                       func(context.Context, *kabs.RequestInformation, serialization.ParsableFactory, kabs.ErrorMappings) (serialization.Parsable, error)
	sendPrimitive              func(context.Context, *kabs.RequestInformation, string, kabs.ErrorMappings) (any, error)
	convertToNativeRequest     func(context.Context, *kabs.RequestInformation) (any, error)
	baseURL                    string
	serializationWriterFactory serialization.SerializationWriterFactory
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func (f *fakeRequestAdapter) Send(ctx context.Context, requestInfo *kabs.RequestInformation, constructor serialization.ParsableFactory, errorMappings kabs.ErrorMappings) (serialization.Parsable, error) {
	if f.send == nil {
		panic("unexpected Send call")
	}
	return f.send(ctx, requestInfo, constructor, errorMappings)
}

func (f *fakeRequestAdapter) SendEnum(context.Context, *kabs.RequestInformation, serialization.EnumFactory, kabs.ErrorMappings) (any, error) {
	panic("unexpected SendEnum call")
}

func (f *fakeRequestAdapter) SendCollection(context.Context, *kabs.RequestInformation, serialization.ParsableFactory, kabs.ErrorMappings) ([]serialization.Parsable, error) {
	panic("unexpected SendCollection call")
}

func (f *fakeRequestAdapter) SendEnumCollection(context.Context, *kabs.RequestInformation, serialization.EnumFactory, kabs.ErrorMappings) ([]any, error) {
	panic("unexpected SendEnumCollection call")
}

func (f *fakeRequestAdapter) SendPrimitive(ctx context.Context, requestInfo *kabs.RequestInformation, typeName string, errorMappings kabs.ErrorMappings) (any, error) {
	if f.sendPrimitive == nil {
		panic("unexpected SendPrimitive call")
	}
	return f.sendPrimitive(ctx, requestInfo, typeName, errorMappings)
}

func (f *fakeRequestAdapter) SendPrimitiveCollection(context.Context, *kabs.RequestInformation, string, kabs.ErrorMappings) ([]any, error) {
	panic("unexpected SendPrimitiveCollection call")
}

func (f *fakeRequestAdapter) SendNoContent(context.Context, *kabs.RequestInformation, kabs.ErrorMappings) error {
	panic("unexpected SendNoContent call")
}

func (f *fakeRequestAdapter) GetSerializationWriterFactory() serialization.SerializationWriterFactory {
	return f.serializationWriterFactory
}

func (f *fakeRequestAdapter) EnableBackingStore(store.BackingStoreFactory) {}

func (f *fakeRequestAdapter) SetBaseUrl(baseURL string) {
	f.baseURL = baseURL
}

func (f *fakeRequestAdapter) GetBaseUrl() string {
	return f.baseURL
}

func (f *fakeRequestAdapter) ConvertToNativeRequest(ctx context.Context, requestInfo *kabs.RequestInformation) (any, error) {
	if f.convertToNativeRequest == nil {
		panic("unexpected ConvertToNativeRequest call")
	}
	return f.convertToNativeRequest(ctx, requestInfo)
}

func assertEqual(t interface {
	Helper()
	Errorf(string, ...any)
}, got, want any) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Mismatch found:\nGot: %#v\nWant: %#v", got, want)
	}
}

func ptr[T any](value T) *T {
	return &value
}

func TestClient(t *testing.T) {
	t.Run("New requires BaseURL", func(t *testing.T) {
		client, err := ssclient.New(ssclient.Config{})
		if err == nil {
			t.Fatal("expected error")
		}
		if got, want := err.Error(), "base URL is required"; got != want {
			t.Fatalf("unexpected error %q want %q", got, want)
		}
		if client != nil {
			t.Fatal("expected nil client")
		}
	})

	t.Run("Raw", func(t *testing.T) {
		client, err := ssclient.New(ssclient.Config{
			BaseURL:    "http://storage.service",
			Username:   "test",
			Key:        "test",
			HTTPClient: &http.Client{},
		})
		assertEqual(t, err, nil)

		if client.Raw() == nil {
			t.Fatal("expected raw client")
		}
		if client.Packages() == nil {
			t.Fatal("expected packages service")
		}
	})
}
