package ssclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"time"

	"github.com/google/uuid"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	absser "github.com/microsoft/kiota-abstractions-go/serialization"
	jsonser "github.com/microsoft/kiota-serialization-json-go"
)

var tolerantTimeLayouts = []string{
	time.RFC3339Nano,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05",
}

type tolerantJSONParseNodeFactory struct {
	inner absser.ParseNodeFactory
}

func newTolerantParseNodeFactory() absser.ParseNodeFactory {
	registry := absser.NewParseNodeFactoryRegistry()

	absser.DefaultParseNodeFactoryInstance.Lock()
	defer absser.DefaultParseNodeFactoryInstance.Unlock()

	maps.Copy(registry.ContentTypeAssociatedFactories, absser.DefaultParseNodeFactoryInstance.ContentTypeAssociatedFactories)

	jsonFactory, ok := registry.ContentTypeAssociatedFactories["application/json"]
	if !ok {
		jsonFactory = jsonser.NewJsonParseNodeFactory()
	}
	registry.ContentTypeAssociatedFactories["application/json"] = &tolerantJSONParseNodeFactory{
		inner: jsonFactory,
	}

	return registry
}

func (f *tolerantJSONParseNodeFactory) GetValidContentType() (string, error) {
	return f.inner.GetValidContentType()
}

func (f *tolerantJSONParseNodeFactory) GetRootParseNode(contentType string, content []byte) (absser.ParseNode, error) {
	node, err := f.inner.GetRootParseNode(contentType, content)
	if err != nil {
		return nil, err
	}
	return &tolerantParseNode{
		inner:       node,
		factory:     f,
		contentType: contentType,
	}, nil
}

type tolerantParseNode struct {
	inner       absser.ParseNode
	factory     absser.ParseNodeFactory
	contentType string
	before      absser.ParsableAction
	after       absser.ParsableAction
}

func isNilParseNode(node absser.ParseNode) bool {
	if node == nil {
		return true
	}
	value := reflect.ValueOf(node)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func (n *tolerantParseNode) GetChildNode(index string) (absser.ParseNode, error) {
	if n == nil || isNilParseNode(n.inner) {
		return nil, nil
	}
	child, err := n.inner.GetChildNode(index)
	if err != nil || child == nil {
		return child, err
	}
	return n.wrapChild(child)
}

func (n *tolerantParseNode) GetCollectionOfObjectValues(ctor absser.ParsableFactory) ([]absser.Parsable, error) {
	raw, err := n.inner.GetRawValue()
	if err != nil || raw == nil {
		return nil, err
	}

	items, ok := raw.([]any)
	if !ok {
		return nil, errors.New("value is not a collection")
	}

	result := make([]absser.Parsable, len(items))
	for i, item := range items {
		if item == nil {
			continue
		}
		child, err := n.childFromRaw(item)
		if err != nil {
			return nil, err
		}
		value, err := child.GetObjectValue(ctor)
		if err != nil {
			return nil, err
		}
		result[i] = value
	}

	return result, nil
}

func (n *tolerantParseNode) GetCollectionOfPrimitiveValues(targetType string) ([]any, error) {
	raw, err := n.inner.GetRawValue()
	if err != nil || raw == nil {
		return nil, err
	}

	items, ok := raw.([]any)
	if !ok {
		return nil, errors.New("value is not a collection")
	}

	result := make([]any, len(items))
	for i, item := range items {
		if item == nil {
			continue
		}
		child, err := n.childFromRaw(item)
		if err != nil {
			return nil, err
		}
		value, err := child.(*tolerantParseNode).getPrimitiveValue(targetType)
		if err != nil {
			return nil, err
		}
		result[i] = value
	}

	return result, nil
}

func (n *tolerantParseNode) GetCollectionOfEnumValues(parser absser.EnumFactory) ([]any, error) {
	raw, err := n.inner.GetRawValue()
	if err != nil || raw == nil {
		return nil, err
	}

	items, ok := raw.([]any)
	if !ok {
		return nil, errors.New("value is not a collection")
	}

	result := make([]any, len(items))
	for i, item := range items {
		if item == nil {
			continue
		}
		child, err := n.childFromRaw(item)
		if err != nil {
			return nil, err
		}
		value, err := child.GetEnumValue(parser)
		if err != nil {
			return nil, err
		}
		result[i] = value
	}

	return result, nil
}

func (n *tolerantParseNode) GetObjectValue(ctor absser.ParsableFactory) (absser.Parsable, error) {
	if ctor == nil {
		return nil, errors.New("constructor is nil")
	}

	raw, err := n.inner.GetRawValue()
	if err != nil || raw == nil {
		return nil, err
	}

	result, err := ctor(n)
	if err != nil {
		return nil, err
	}

	if _, isUntypedNode := result.(absser.UntypedNodeable); isUntypedNode {
		return toUntypedNode(raw), nil
	}

	abstractions.InvokeParsableAction(n.before, result)

	properties, ok := raw.(map[string]any)
	if ok {
		fields := result.GetFieldDeserializers()
		holder, isHolder := result.(absser.AdditionalDataHolder)
		var additionalData map[string]any
		if isHolder {
			additionalData = holder.GetAdditionalData()
			if additionalData == nil {
				additionalData = make(map[string]any)
				holder.SetAdditionalData(additionalData)
			}
		}

		for key, value := range properties {
			child, err := n.childFromRaw(value)
			if err != nil {
				return nil, err
			}

			field := fields[key]
			if field == nil {
				if isHolder {
					rawValue, err := child.GetRawValue()
					if err != nil {
						return nil, err
					}
					if rawValue != nil {
						additionalData[key] = rawValue
					}
				}
				continue
			}

			if err := field(child); err != nil {
				return nil, err
			}
		}
	}

	abstractions.InvokeParsableAction(n.after, result)
	return result, nil
}

func toUntypedNode(value any) absser.UntypedNodeable {
	switch v := value.(type) {
	case nil:
		return absser.NewUntypedNull()
	case bool:
		return absser.NewUntypedBoolean(v)
	case *bool:
		if v == nil {
			return absser.NewUntypedNull()
		}
		return absser.NewUntypedBoolean(*v)
	case string:
		return absser.NewUntypedString(v)
	case *string:
		if v == nil {
			return absser.NewUntypedNull()
		}
		return absser.NewUntypedString(*v)
	case float32:
		return absser.NewUntypedFloat(v)
	case *float32:
		if v == nil {
			return absser.NewUntypedNull()
		}
		return absser.NewUntypedFloat(*v)
	case float64:
		return absser.NewUntypedDouble(v)
	case *float64:
		if v == nil {
			return absser.NewUntypedNull()
		}
		return absser.NewUntypedDouble(*v)
	case int32:
		return absser.NewUntypedInteger(v)
	case *int32:
		if v == nil {
			return absser.NewUntypedNull()
		}
		return absser.NewUntypedInteger(*v)
	case int64:
		return absser.NewUntypedLong(v)
	case *int64:
		if v == nil {
			return absser.NewUntypedNull()
		}
		return absser.NewUntypedLong(*v)
	case map[string]any:
		properties := make(map[string]absser.UntypedNodeable, len(v))
		for key, element := range v {
			properties[key] = toUntypedNode(element)
		}
		return absser.NewUntypedObject(properties)
	case []any:
		items := make([]absser.UntypedNodeable, len(v))
		for i, element := range v {
			items[i] = toUntypedNode(element)
		}
		return absser.NewUntypedArray(items)
	default:
		return absser.NewUntypedNode(v)
	}
}

func (n *tolerantParseNode) GetStringValue() (*string, error) {
	if n == nil || isNilParseNode(n.inner) {
		return nil, nil
	}
	return n.inner.GetStringValue()
}

func (n *tolerantParseNode) GetBoolValue() (*bool, error) {
	if n == nil || isNilParseNode(n.inner) {
		return nil, nil
	}
	return n.inner.GetBoolValue()
}

func (n *tolerantParseNode) GetInt8Value() (*int8, error) {
	if n == nil || isNilParseNode(n.inner) {
		return nil, nil
	}
	return n.inner.GetInt8Value()
}

func (n *tolerantParseNode) GetByteValue() (*byte, error) {
	if n == nil || isNilParseNode(n.inner) {
		return nil, nil
	}
	return n.inner.GetByteValue()
}

func (n *tolerantParseNode) GetFloat32Value() (*float32, error) {
	if n == nil || isNilParseNode(n.inner) {
		return nil, nil
	}
	return n.inner.GetFloat32Value()
}

func (n *tolerantParseNode) GetFloat64Value() (*float64, error) {
	if n == nil || isNilParseNode(n.inner) {
		return nil, nil
	}
	return n.inner.GetFloat64Value()
}

func (n *tolerantParseNode) GetInt32Value() (*int32, error) {
	if n == nil || isNilParseNode(n.inner) {
		return nil, nil
	}
	return n.inner.GetInt32Value()
}

func (n *tolerantParseNode) GetInt64Value() (*int64, error) {
	if n == nil || isNilParseNode(n.inner) {
		return nil, nil
	}
	return n.inner.GetInt64Value()
}

func (n *tolerantParseNode) GetTimeValue() (*time.Time, error) {
	if n == nil || isNilParseNode(n.inner) {
		return nil, nil
	}
	v, err := n.inner.GetStringValue()
	if err != nil || v == nil {
		return nil, err
	}

	for _, layout := range tolerantTimeLayouts {
		parsed, parseErr := time.Parse(layout, *v)
		if parseErr == nil {
			if layout != time.RFC3339Nano {
				parsed = time.Date(
					parsed.Year(),
					parsed.Month(),
					parsed.Day(),
					parsed.Hour(),
					parsed.Minute(),
					parsed.Second(),
					parsed.Nanosecond(),
					time.UTC,
				)
			}
			return &parsed, nil
		}
	}

	return nil, fmt.Errorf("parse time %q: unsupported timestamp layout", *v)
}

func (n *tolerantParseNode) GetISODurationValue() (*absser.ISODuration, error) {
	if n == nil || isNilParseNode(n.inner) {
		return nil, nil
	}
	return n.inner.GetISODurationValue()
}

func (n *tolerantParseNode) GetTimeOnlyValue() (*absser.TimeOnly, error) {
	if n == nil || isNilParseNode(n.inner) {
		return nil, nil
	}
	return n.inner.GetTimeOnlyValue()
}

func (n *tolerantParseNode) GetDateOnlyValue() (*absser.DateOnly, error) {
	if n == nil || isNilParseNode(n.inner) {
		return nil, nil
	}
	return n.inner.GetDateOnlyValue()
}

func (n *tolerantParseNode) GetUUIDValue() (*uuid.UUID, error) {
	if n == nil || isNilParseNode(n.inner) {
		return nil, nil
	}
	return n.inner.GetUUIDValue()
}

func (n *tolerantParseNode) GetEnumValue(parser absser.EnumFactory) (any, error) {
	if n == nil || isNilParseNode(n.inner) {
		return nil, nil
	}
	return n.inner.GetEnumValue(parser)
}

func (n *tolerantParseNode) GetByteArrayValue() ([]byte, error) {
	if n == nil || isNilParseNode(n.inner) {
		return nil, nil
	}
	return n.inner.GetByteArrayValue()
}

func (n *tolerantParseNode) GetRawValue() (any, error) {
	if n == nil || isNilParseNode(n.inner) {
		return nil, nil
	}
	return n.inner.GetRawValue()
}

func (n *tolerantParseNode) GetOnBeforeAssignFieldValues() absser.ParsableAction {
	return n.before
}

func (n *tolerantParseNode) SetOnBeforeAssignFieldValues(action absser.ParsableAction) error {
	if n == nil || isNilParseNode(n.inner) {
		return nil
	}
	n.before = action
	return n.inner.SetOnBeforeAssignFieldValues(action)
}

func (n *tolerantParseNode) GetOnAfterAssignFieldValues() absser.ParsableAction {
	return n.after
}

func (n *tolerantParseNode) SetOnAfterAssignFieldValues(action absser.ParsableAction) error {
	if n == nil || isNilParseNode(n.inner) {
		return nil
	}
	n.after = action
	return n.inner.SetOnAfterAssignFieldValues(action)
}

func (n *tolerantParseNode) wrapChild(child absser.ParseNode) (absser.ParseNode, error) {
	wrapped := &tolerantParseNode{
		inner:       child,
		factory:     n.factory,
		contentType: n.contentType,
		before:      n.before,
		after:       n.after,
	}
	if err := wrapped.inner.SetOnBeforeAssignFieldValues(n.before); err != nil {
		return nil, err
	}
	if err := wrapped.inner.SetOnAfterAssignFieldValues(n.after); err != nil {
		return nil, err
	}
	return wrapped, nil
}

func (n *tolerantParseNode) childFromRaw(raw any) (absser.ParseNode, error) {
	content, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	child, err := n.factory.GetRootParseNode(n.contentType, content)
	if err != nil {
		return nil, err
	}
	if child == nil {
		return &tolerantParseNode{
			factory:     n.factory,
			contentType: n.contentType,
			before:      n.before,
			after:       n.after,
		}, nil
	}

	if err := child.SetOnBeforeAssignFieldValues(n.before); err != nil {
		return nil, err
	}
	if err := child.SetOnAfterAssignFieldValues(n.after); err != nil {
		return nil, err
	}

	return child, nil
}

func (n *tolerantParseNode) getPrimitiveValue(targetType string) (any, error) {
	switch targetType {
	case "string":
		return n.GetStringValue()
	case "bool":
		return n.GetBoolValue()
	case "uint8":
		return n.GetInt8Value()
	case "byte":
		return n.GetByteValue()
	case "float32":
		return n.GetFloat32Value()
	case "float64":
		return n.GetFloat64Value()
	case "int32":
		return n.GetInt32Value()
	case "int64":
		return n.GetInt64Value()
	case "time":
		return n.GetTimeValue()
	case "timeonly":
		return n.GetTimeOnlyValue()
	case "dateonly":
		return n.GetDateOnlyValue()
	case "isoduration":
		return n.GetISODurationValue()
	case "uuid":
		return n.GetUUIDValue()
	case "base64":
		return n.GetByteArrayValue()
	default:
		return nil, fmt.Errorf("targetType %s is not supported", targetType)
	}
}
