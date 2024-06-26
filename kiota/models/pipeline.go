package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Pipeline struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]any
    // The description property
    description *string
    // The remote_name property
    remote_name *string
    // The resource_uri property
    resource_uri *string
    // The uuid property
    uuid *string
}
// NewPipeline instantiates a new Pipeline and sets the default values.
func NewPipeline()(*Pipeline) {
    m := &Pipeline{
    }
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreatePipelineFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePipelineFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPipeline(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *Pipeline) GetAdditionalData()(map[string]any) {
    return m.additionalData
}
// GetDescription gets the description property value. The description property
// returns a *string when successful
func (m *Pipeline) GetDescription()(*string) {
    return m.description
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Pipeline) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
        }
        return nil
    }
    res["remote_name"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRemoteName(val)
        }
        return nil
    }
    res["resource_uri"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResourceUri(val)
        }
        return nil
    }
    res["uuid"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUuid(val)
        }
        return nil
    }
    return res
}
// GetRemoteName gets the remote_name property value. The remote_name property
// returns a *string when successful
func (m *Pipeline) GetRemoteName()(*string) {
    return m.remote_name
}
// GetResourceUri gets the resource_uri property value. The resource_uri property
// returns a *string when successful
func (m *Pipeline) GetResourceUri()(*string) {
    return m.resource_uri
}
// GetUuid gets the uuid property value. The uuid property
// returns a *string when successful
func (m *Pipeline) GetUuid()(*string) {
    return m.uuid
}
// Serialize serializes information the current object
func (m *Pipeline) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("remote_name", m.GetRemoteName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("resource_uri", m.GetResourceUri())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("uuid", m.GetUuid())
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
func (m *Pipeline) SetAdditionalData(value map[string]any)() {
    m.additionalData = value
}
// SetDescription sets the description property value. The description property
func (m *Pipeline) SetDescription(value *string)() {
    m.description = value
}
// SetRemoteName sets the remote_name property value. The remote_name property
func (m *Pipeline) SetRemoteName(value *string)() {
    m.remote_name = value
}
// SetResourceUri sets the resource_uri property value. The resource_uri property
func (m *Pipeline) SetResourceUri(value *string)() {
    m.resource_uri = value
}
// SetUuid sets the uuid property value. The uuid property
func (m *Pipeline) SetUuid(value *string)() {
    m.uuid = value
}
type Pipelineable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDescription()(*string)
    GetRemoteName()(*string)
    GetResourceUri()(*string)
    GetUuid()(*string)
    SetDescription(value *string)()
    SetRemoteName(value *string)()
    SetResourceUri(value *string)()
    SetUuid(value *string)()
}
