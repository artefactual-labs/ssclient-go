package ssclient

import (
	"fmt"
	"net/url"
	"strings"
)

// ParseResourceURI extracts the resource name and UUID from a Storage Service
// resource URI.
func ParseResourceURI(resourceURI string) (resource string, uuid string, err error) {
	if resourceURI == "" {
		return "", "", fmt.Errorf("empty resource URI")
	}

	parsed, err := url.Parse(resourceURI)
	if err != nil {
		return "", "", fmt.Errorf("parse resource URI: %w", err)
	}

	segments := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(segments) < 4 {
		return "", "", fmt.Errorf("invalid resource URI %q", resourceURI)
	}

	offset := len(segments) - 4
	if segments[offset] != "api" || segments[offset+1] != "v2" {
		return "", "", fmt.Errorf("invalid resource URI %q", resourceURI)
	}

	resource = segments[offset+2]
	uuid = segments[offset+3]
	if resource == "" || uuid == "" {
		return "", "", fmt.Errorf("invalid resource URI %q", resourceURI)
	}

	return resource, uuid, nil
}

// MustParseResourceURI extracts the resource name and UUID from a Storage
// Service resource URI and panics if the URI is invalid.
func MustParseResourceURI(resourceURI string) (resource string, uuid string) {
	resource, uuid, err := ParseResourceURI(resourceURI)
	if err != nil {
		panic(err)
	}
	return resource, uuid
}
