package api

import (
    i2b7a3625368152c59661ed1a63c26960f7e9cda05d0fbc8e5d79ac57ca250e0a "go.artefactual.dev/ssclient/kiota/models"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type V2LocationItemWithUuPostResponse struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]any
    // List of objects containing `source` and `destination`. The source and destination are paths relative to their Location of the files to be moved.
    files []i2b7a3625368152c59661ed1a63c26960f7e9cda05d0fbc8e5d79ac57ca250e0a.MoveFileable
    // URI of the Location the files should be moved from.
    origin_location *string
    // URI of the
    pipeline *string
}
// NewV2LocationItemWithUuPostResponse instantiates a new V2LocationItemWithUuPostResponse and sets the default values.
func NewV2LocationItemWithUuPostResponse()(*V2LocationItemWithUuPostResponse) {
    m := &V2LocationItemWithUuPostResponse{
    }
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateV2LocationItemWithUuPostResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateV2LocationItemWithUuPostResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewV2LocationItemWithUuPostResponse(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *V2LocationItemWithUuPostResponse) GetAdditionalData()(map[string]any) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *V2LocationItemWithUuPostResponse) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["files"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(i2b7a3625368152c59661ed1a63c26960f7e9cda05d0fbc8e5d79ac57ca250e0a.CreateMoveFileFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]i2b7a3625368152c59661ed1a63c26960f7e9cda05d0fbc8e5d79ac57ca250e0a.MoveFileable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(i2b7a3625368152c59661ed1a63c26960f7e9cda05d0fbc8e5d79ac57ca250e0a.MoveFileable)
                }
            }
            m.SetFiles(res)
        }
        return nil
    }
    res["origin_location"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOriginLocation(val)
        }
        return nil
    }
    res["pipeline"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPipeline(val)
        }
        return nil
    }
    return res
}
// GetFiles gets the files property value. List of objects containing `source` and `destination`. The source and destination are paths relative to their Location of the files to be moved.
// returns a []MoveFileable when successful
func (m *V2LocationItemWithUuPostResponse) GetFiles()([]i2b7a3625368152c59661ed1a63c26960f7e9cda05d0fbc8e5d79ac57ca250e0a.MoveFileable) {
    return m.files
}
// GetOriginLocation gets the origin_location property value. URI of the Location the files should be moved from.
// returns a *string when successful
func (m *V2LocationItemWithUuPostResponse) GetOriginLocation()(*string) {
    return m.origin_location
}
// GetPipeline gets the pipeline property value. URI of the
// returns a *string when successful
func (m *V2LocationItemWithUuPostResponse) GetPipeline()(*string) {
    return m.pipeline
}
// Serialize serializes information the current object
func (m *V2LocationItemWithUuPostResponse) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetFiles() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetFiles()))
        for i, v := range m.GetFiles() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("files", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("origin_location", m.GetOriginLocation())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("pipeline", m.GetPipeline())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *V2LocationItemWithUuPostResponse) SetAdditionalData(value map[string]any)() {
    m.additionalData = value
}
// SetFiles sets the files property value. List of objects containing `source` and `destination`. The source and destination are paths relative to their Location of the files to be moved.
func (m *V2LocationItemWithUuPostResponse) SetFiles(value []i2b7a3625368152c59661ed1a63c26960f7e9cda05d0fbc8e5d79ac57ca250e0a.MoveFileable)() {
    m.files = value
}
// SetOriginLocation sets the origin_location property value. URI of the Location the files should be moved from.
func (m *V2LocationItemWithUuPostResponse) SetOriginLocation(value *string)() {
    m.origin_location = value
}
// SetPipeline sets the pipeline property value. URI of the
func (m *V2LocationItemWithUuPostResponse) SetPipeline(value *string)() {
    m.pipeline = value
}
type V2LocationItemWithUuPostResponseable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetFiles()([]i2b7a3625368152c59661ed1a63c26960f7e9cda05d0fbc8e5d79ac57ca250e0a.MoveFileable)
    GetOriginLocation()(*string)
    GetPipeline()(*string)
    SetFiles(value []i2b7a3625368152c59661ed1a63c26960f7e9cda05d0fbc8e5d79ac57ca250e0a.MoveFileable)()
    SetOriginLocation(value *string)()
    SetPipeline(value *string)()
}
