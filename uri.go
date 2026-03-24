package ssclient

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ParseResourceURI extracts the resource name and UUID from a Storage Service
// resource URI.
func ParseResourceURI(resourceURI string) (resource, uuid string, err error) {
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
func MustParseResourceURI(resourceURI string) (resource, uuid string) {
	resource, uuid, err := ParseResourceURI(resourceURI)
	if err != nil {
		panic(err)
	}
	return resource, uuid
}

// ParseAsyncOperationURI extracts the task ID from a Storage Service async
// operation URI such as "/api/v2/async/1/".
func ParseAsyncOperationURI(resourceURI string) (int, error) {
	resource, identifier, err := ParseResourceURI(resourceURI)
	if err != nil {
		return 0, err
	}
	if resource != "async" {
		return 0, fmt.Errorf("invalid async operation URI %q", resourceURI)
	}

	id, err := strconv.Atoi(identifier)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid async operation URI %q", resourceURI)
	}
	return id, nil
}

// MustParseAsyncOperationURI extracts the task ID from a Storage Service async
// operation URI and panics if the URI is invalid.
func MustParseAsyncOperationURI(resourceURI string) int {
	id, err := ParseAsyncOperationURI(resourceURI)
	if err != nil {
		panic(err)
	}
	return id
}
