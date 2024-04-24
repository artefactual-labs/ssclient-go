package api

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use V2LocationItemWithUuPostResponseable instead.
type V2LocationItemWithUuResponse struct {
    V2LocationItemWithUuPostResponse
}
// NewV2LocationItemWithUuResponse instantiates a new V2LocationItemWithUuResponse and sets the default values.
func NewV2LocationItemWithUuResponse()(*V2LocationItemWithUuResponse) {
    m := &V2LocationItemWithUuResponse{
        V2LocationItemWithUuPostResponse: *NewV2LocationItemWithUuPostResponse(),
    }
    return m
}
// CreateV2LocationItemWithUuResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateV2LocationItemWithUuResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewV2LocationItemWithUuResponse(), nil
}
// Deprecated: This class is obsolete. Use V2LocationItemWithUuPostResponseable instead.
type V2LocationItemWithUuResponseable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    V2LocationItemWithUuPostResponseable
}
