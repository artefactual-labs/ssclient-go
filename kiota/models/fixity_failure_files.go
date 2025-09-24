package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type FixityFailureFiles struct {
    // Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
    additionalData map[string]any
    // The changed property
    changed []string
    // The missing property
    missing []string
    // The untracked property
    untracked []string
}
// NewFixityFailureFiles instantiates a new FixityFailureFiles and sets the default values.
func NewFixityFailureFiles()(*FixityFailureFiles) {
    m := &FixityFailureFiles{
    }
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateFixityFailureFilesFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateFixityFailureFilesFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFixityFailureFiles(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *FixityFailureFiles) GetAdditionalData()(map[string]any) {
    return m.additionalData
}
// GetChanged gets the changed property value. The changed property
// returns a []string when successful
func (m *FixityFailureFiles) GetChanged()([]string) {
    return m.changed
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *FixityFailureFiles) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["changed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetChanged(res)
        }
        return nil
    }
    res["missing"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetMissing(res)
        }
        return nil
    }
    res["untracked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetUntracked(res)
        }
        return nil
    }
    return res
}
// GetMissing gets the missing property value. The missing property
// returns a []string when successful
func (m *FixityFailureFiles) GetMissing()([]string) {
    return m.missing
}
// GetUntracked gets the untracked property value. The untracked property
// returns a []string when successful
func (m *FixityFailureFiles) GetUntracked()([]string) {
    return m.untracked
}
// Serialize serializes information the current object
func (m *FixityFailureFiles) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetChanged() != nil {
        err := writer.WriteCollectionOfStringValues("changed", m.GetChanged())
        if err != nil {
            return err
        }
    }
    if m.GetMissing() != nil {
        err := writer.WriteCollectionOfStringValues("missing", m.GetMissing())
        if err != nil {
            return err
        }
    }
    if m.GetUntracked() != nil {
        err := writer.WriteCollectionOfStringValues("untracked", m.GetUntracked())
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
func (m *FixityFailureFiles) SetAdditionalData(value map[string]any)() {
    m.additionalData = value
}
// SetChanged sets the changed property value. The changed property
func (m *FixityFailureFiles) SetChanged(value []string)() {
    m.changed = value
}
// SetMissing sets the missing property value. The missing property
func (m *FixityFailureFiles) SetMissing(value []string)() {
    m.missing = value
}
// SetUntracked sets the untracked property value. The untracked property
func (m *FixityFailureFiles) SetUntracked(value []string)() {
    m.untracked = value
}
type FixityFailureFilesable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetChanged()([]string)
    GetMissing()([]string)
    GetUntracked()([]string)
    SetChanged(value []string)()
    SetMissing(value []string)()
    SetUntracked(value []string)()
}
