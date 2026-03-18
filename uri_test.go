package ssclient_test

import (
	"testing"

	"go.artefactual.dev/ssclient"
)

func TestParseResourceURI(t *testing.T) {
	tests := []struct {
		name         string
		resourceURI  string
		wantResource string
		wantUUID     string
		wantErr      bool
	}{
		{
			name:         "absolute URI",
			resourceURI:  "https://example.test/api/v2/file/96922350-ccde-4fb0-a999-d2010522028f/",
			wantResource: "file",
			wantUUID:     "96922350-ccde-4fb0-a999-d2010522028f",
		},
		{
			name:         "relative URI",
			resourceURI:  "/api/v2/location/610bc407-ba6c-4dcd-8675-d2727a9aab18/",
			wantResource: "location",
			wantUUID:     "610bc407-ba6c-4dcd-8675-d2727a9aab18",
		},
		{
			name:         "base path prefix",
			resourceURI:  "https://example.test/storage/api/v2/pipeline/a64e061a-5688-49b5-95c1-0b6885c40c04/",
			wantResource: "pipeline",
			wantUUID:     "a64e061a-5688-49b5-95c1-0b6885c40c04",
		},
		{
			name:        "reject endpoint URI",
			resourceURI: "/api/v2/file/96922350-ccde-4fb0-a999-d2010522028f/move/",
			wantErr:     true,
		},
		{
			name:        "reject nested path with second UUID",
			resourceURI: "/api/v2/file/96922350-ccde-4fb0-a999-d2010522028f/replica/610bc407-ba6c-4dcd-8675-d2727a9aab18/",
			wantErr:     true,
		},
		{
			name:        "reject empty URI",
			resourceURI: "",
			wantErr:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resource, uuid, err := ssclient.ParseResourceURI(test.resourceURI)
			if test.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}

			assertEqual(t, err, nil)
			assertEqual(t, resource, test.wantResource)
			assertEqual(t, uuid, test.wantUUID)
		})
	}
}

func TestMustParseResourceURI(t *testing.T) {
	t.Run("valid URI", func(t *testing.T) {
		resource, uuid := ssclient.MustParseResourceURI("/api/v2/space/141593ff-2a27-44a1-9de1-917573fa0f4a/")
		assertEqual(t, resource, "space")
		assertEqual(t, uuid, "141593ff-2a27-44a1-9de1-917573fa0f4a")
	})

	t.Run("panics on invalid URI", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic")
			}
		}()

		ssclient.MustParseResourceURI("/api/v2/file/not-a-resource/move/")
	})
}
