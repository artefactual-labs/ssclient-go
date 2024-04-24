package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ListResponseMeta struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]any
    // The limit property
    limit *int32
    // The next property
    next *string
    // The offset property
    offset *int32
    // The previous property
    previous *string
    // The total_count property
    total_count *int32
}
// NewListResponseMeta instantiates a new ListResponseMeta and sets the default values.
func NewListResponseMeta()(*ListResponseMeta) {
    m := &ListResponseMeta{
    }
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateListResponseMetaFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateListResponseMetaFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewListResponseMeta(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *ListResponseMeta) GetAdditionalData()(map[string]any) {
    return m.additionalData
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ListResponseMeta) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["limit"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLimit(val)
        }
        return nil
    }
    res["next"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNext(val)
        }
        return nil
    }
    res["offset"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOffset(val)
        }
        return nil
    }
    res["previous"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrevious(val)
        }
        return nil
    }
    res["total_count"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalCount(val)
        }
        return nil
    }
    return res
}
// GetLimit gets the limit property value. The limit property
// returns a *int32 when successful
func (m *ListResponseMeta) GetLimit()(*int32) {
    return m.limit
}
// GetNext gets the next property value. The next property
// returns a *string when successful
func (m *ListResponseMeta) GetNext()(*string) {
    return m.next
}
// GetOffset gets the offset property value. The offset property
// returns a *int32 when successful
func (m *ListResponseMeta) GetOffset()(*int32) {
    return m.offset
}
// GetPrevious gets the previous property value. The previous property
// returns a *string when successful
func (m *ListResponseMeta) GetPrevious()(*string) {
    return m.previous
}
// GetTotalCount gets the total_count property value. The total_count property
// returns a *int32 when successful
func (m *ListResponseMeta) GetTotalCount()(*int32) {
    return m.total_count
}
// Serialize serializes information the current object
func (m *ListResponseMeta) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("limit", m.GetLimit())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("next", m.GetNext())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("offset", m.GetOffset())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("previous", m.GetPrevious())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("total_count", m.GetTotalCount())
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
func (m *ListResponseMeta) SetAdditionalData(value map[string]any)() {
    m.additionalData = value
}
// SetLimit sets the limit property value. The limit property
func (m *ListResponseMeta) SetLimit(value *int32)() {
    m.limit = value
}
// SetNext sets the next property value. The next property
func (m *ListResponseMeta) SetNext(value *string)() {
    m.next = value
}
// SetOffset sets the offset property value. The offset property
func (m *ListResponseMeta) SetOffset(value *int32)() {
    m.offset = value
}
// SetPrevious sets the previous property value. The previous property
func (m *ListResponseMeta) SetPrevious(value *string)() {
    m.previous = value
}
// SetTotalCount sets the total_count property value. The total_count property
func (m *ListResponseMeta) SetTotalCount(value *int32)() {
    m.total_count = value
}
type ListResponseMetaable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetLimit()(*int32)
    GetNext()(*string)
    GetOffset()(*int32)
    GetPrevious()(*string)
    GetTotalCount()(*int32)
    SetLimit(value *int32)()
    SetNext(value *string)()
    SetOffset(value *int32)()
    SetPrevious(value *string)()
    SetTotalCount(value *int32)()
}
