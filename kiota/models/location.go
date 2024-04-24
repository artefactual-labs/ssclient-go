package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Location struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]any
    // The description property
    description *string
    // The enabled property
    enabled *bool
    // The path property
    path *string
    // The pipeline property
    pipeline []string
    // The purpose property
    purpose *LocationPurpose
    // The quota property
    quota *int32
    // The relative_path property
    relative_path *string
    // The resource_uri property
    resource_uri *string
    // The space property
    space *string
    // The used property
    used *int32
    // The uuid property
    uuid *string
}
// NewLocation instantiates a new Location and sets the default values.
func NewLocation()(*Location) {
    m := &Location{
    }
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateLocationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateLocationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewLocation(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *Location) GetAdditionalData()(map[string]any) {
    return m.additionalData
}
// GetDescription gets the description property value. The description property
// returns a *string when successful
func (m *Location) GetDescription()(*string) {
    return m.description
}
// GetEnabled gets the enabled property value. The enabled property
// returns a *bool when successful
func (m *Location) GetEnabled()(*bool) {
    return m.enabled
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Location) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["enabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnabled(val)
        }
        return nil
    }
    res["path"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPath(val)
        }
        return nil
    }
    res["pipeline"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetPipeline(res)
        }
        return nil
    }
    res["purpose"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseLocationPurpose)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPurpose(val.(*LocationPurpose))
        }
        return nil
    }
    res["quota"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQuota(val)
        }
        return nil
    }
    res["relative_path"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRelativePath(val)
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
    res["space"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSpace(val)
        }
        return nil
    }
    res["used"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUsed(val)
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
// GetPath gets the path property value. The path property
// returns a *string when successful
func (m *Location) GetPath()(*string) {
    return m.path
}
// GetPipeline gets the pipeline property value. The pipeline property
// returns a []string when successful
func (m *Location) GetPipeline()([]string) {
    return m.pipeline
}
// GetPurpose gets the purpose property value. The purpose property
// returns a *LocationPurpose when successful
func (m *Location) GetPurpose()(*LocationPurpose) {
    return m.purpose
}
// GetQuota gets the quota property value. The quota property
// returns a *int32 when successful
func (m *Location) GetQuota()(*int32) {
    return m.quota
}
// GetRelativePath gets the relative_path property value. The relative_path property
// returns a *string when successful
func (m *Location) GetRelativePath()(*string) {
    return m.relative_path
}
// GetResourceUri gets the resource_uri property value. The resource_uri property
// returns a *string when successful
func (m *Location) GetResourceUri()(*string) {
    return m.resource_uri
}
// GetSpace gets the space property value. The space property
// returns a *string when successful
func (m *Location) GetSpace()(*string) {
    return m.space
}
// GetUsed gets the used property value. The used property
// returns a *int32 when successful
func (m *Location) GetUsed()(*int32) {
    return m.used
}
// GetUuid gets the uuid property value. The uuid property
// returns a *string when successful
func (m *Location) GetUuid()(*string) {
    return m.uuid
}
// Serialize serializes information the current object
func (m *Location) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("enabled", m.GetEnabled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("path", m.GetPath())
        if err != nil {
            return err
        }
    }
    if m.GetPipeline() != nil {
        err := writer.WriteCollectionOfStringValues("pipeline", m.GetPipeline())
        if err != nil {
            return err
        }
    }
    if m.GetPurpose() != nil {
        cast := (*m.GetPurpose()).String()
        err := writer.WriteStringValue("purpose", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("quota", m.GetQuota())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("relative_path", m.GetRelativePath())
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
        err := writer.WriteStringValue("space", m.GetSpace())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("used", m.GetUsed())
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
func (m *Location) SetAdditionalData(value map[string]any)() {
    m.additionalData = value
}
// SetDescription sets the description property value. The description property
func (m *Location) SetDescription(value *string)() {
    m.description = value
}
// SetEnabled sets the enabled property value. The enabled property
func (m *Location) SetEnabled(value *bool)() {
    m.enabled = value
}
// SetPath sets the path property value. The path property
func (m *Location) SetPath(value *string)() {
    m.path = value
}
// SetPipeline sets the pipeline property value. The pipeline property
func (m *Location) SetPipeline(value []string)() {
    m.pipeline = value
}
// SetPurpose sets the purpose property value. The purpose property
func (m *Location) SetPurpose(value *LocationPurpose)() {
    m.purpose = value
}
// SetQuota sets the quota property value. The quota property
func (m *Location) SetQuota(value *int32)() {
    m.quota = value
}
// SetRelativePath sets the relative_path property value. The relative_path property
func (m *Location) SetRelativePath(value *string)() {
    m.relative_path = value
}
// SetResourceUri sets the resource_uri property value. The resource_uri property
func (m *Location) SetResourceUri(value *string)() {
    m.resource_uri = value
}
// SetSpace sets the space property value. The space property
func (m *Location) SetSpace(value *string)() {
    m.space = value
}
// SetUsed sets the used property value. The used property
func (m *Location) SetUsed(value *int32)() {
    m.used = value
}
// SetUuid sets the uuid property value. The uuid property
func (m *Location) SetUuid(value *string)() {
    m.uuid = value
}
type Locationable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDescription()(*string)
    GetEnabled()(*bool)
    GetPath()(*string)
    GetPipeline()([]string)
    GetPurpose()(*LocationPurpose)
    GetQuota()(*int32)
    GetRelativePath()(*string)
    GetResourceUri()(*string)
    GetSpace()(*string)
    GetUsed()(*int32)
    GetUuid()(*string)
    SetDescription(value *string)()
    SetEnabled(value *bool)()
    SetPath(value *string)()
    SetPipeline(value []string)()
    SetPurpose(value *LocationPurpose)()
    SetQuota(value *int32)()
    SetRelativePath(value *string)()
    SetResourceUri(value *string)()
    SetSpace(value *string)()
    SetUsed(value *int32)()
    SetUuid(value *string)()
}
