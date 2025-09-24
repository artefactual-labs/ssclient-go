package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type FixityResponse struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]any
    // The failures property
    failures FixityFailuresable
    // The message property
    message *string
    // The success property
    success *bool
    // The timestamp property
    timestamp *string
}
// NewFixityResponse instantiates a new FixityResponse and sets the default values.
func NewFixityResponse()(*FixityResponse) {
    m := &FixityResponse{
    }
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateFixityResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateFixityResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFixityResponse(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *FixityResponse) GetAdditionalData()(map[string]any) {
    return m.additionalData
}
// GetFailures gets the failures property value. The failures property
// returns a FixityFailuresable when successful
func (m *FixityResponse) GetFailures()(FixityFailuresable) {
    return m.failures
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *FixityResponse) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["failures"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFixityFailuresFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFailures(val.(FixityFailuresable))
        }
        return nil
    }
    res["message"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMessage(val)
        }
        return nil
    }
    res["success"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSuccess(val)
        }
        return nil
    }
    res["timestamp"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTimestamp(val)
        }
        return nil
    }
    return res
}
// GetMessage gets the message property value. The message property
// returns a *string when successful
func (m *FixityResponse) GetMessage()(*string) {
    return m.message
}
// GetSuccess gets the success property value. The success property
// returns a *bool when successful
func (m *FixityResponse) GetSuccess()(*bool) {
    return m.success
}
// GetTimestamp gets the timestamp property value. The timestamp property
// returns a *string when successful
func (m *FixityResponse) GetTimestamp()(*string) {
    return m.timestamp
}
// Serialize serializes information the current object
func (m *FixityResponse) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("failures", m.GetFailures())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("message", m.GetMessage())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("success", m.GetSuccess())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("timestamp", m.GetTimestamp())
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
func (m *FixityResponse) SetAdditionalData(value map[string]any)() {
    m.additionalData = value
}
// SetFailures sets the failures property value. The failures property
func (m *FixityResponse) SetFailures(value FixityFailuresable)() {
    m.failures = value
}
// SetMessage sets the message property value. The message property
func (m *FixityResponse) SetMessage(value *string)() {
    m.message = value
}
// SetSuccess sets the success property value. The success property
func (m *FixityResponse) SetSuccess(value *bool)() {
    m.success = value
}
// SetTimestamp sets the timestamp property value. The timestamp property
func (m *FixityResponse) SetTimestamp(value *string)() {
    m.timestamp = value
}
type FixityResponseable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetFailures()(FixityFailuresable)
    GetMessage()(*string)
    GetSuccess()(*bool)
    GetTimestamp()(*string)
    SetFailures(value FixityFailuresable)()
    SetMessage(value *string)()
    SetSuccess(value *bool)()
    SetTimestamp(value *string)()
}
