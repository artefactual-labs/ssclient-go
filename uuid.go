package ssclient

import "github.com/google/uuid"

func parseOptionalUUID(value *string) (*uuid.UUID, error) {
	if value == nil {
		return nil, nil
	}

	parsed, err := uuid.Parse(*value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
