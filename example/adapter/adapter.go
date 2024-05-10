// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/microsoft/kiota-abstractions-go (interfaces: RequestAdapter)
//
// Generated by this command:
//
//	mockgen -typed -destination=./adapter/adapter.go -package=adapter github.com/microsoft/kiota-abstractions-go RequestAdapter
//

// Package adapter is a generated GoMock package.
package adapter

import (
	context "context"
	reflect "reflect"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	serialization "github.com/microsoft/kiota-abstractions-go/serialization"
	store "github.com/microsoft/kiota-abstractions-go/store"
	gomock "go.uber.org/mock/gomock"
)

// MockRequestAdapter is a mock of RequestAdapter interface.
type MockRequestAdapter struct {
	ctrl     *gomock.Controller
	recorder *MockRequestAdapterMockRecorder
}

// MockRequestAdapterMockRecorder is the mock recorder for MockRequestAdapter.
type MockRequestAdapterMockRecorder struct {
	mock *MockRequestAdapter
}

// NewMockRequestAdapter creates a new mock instance.
func NewMockRequestAdapter(ctrl *gomock.Controller) *MockRequestAdapter {
	mock := &MockRequestAdapter{ctrl: ctrl}
	mock.recorder = &MockRequestAdapterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRequestAdapter) EXPECT() *MockRequestAdapterMockRecorder {
	return m.recorder
}

// ConvertToNativeRequest mocks base method.
func (m *MockRequestAdapter) ConvertToNativeRequest(arg0 context.Context, arg1 *abstractions.RequestInformation) (any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConvertToNativeRequest", arg0, arg1)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConvertToNativeRequest indicates an expected call of ConvertToNativeRequest.
func (mr *MockRequestAdapterMockRecorder) ConvertToNativeRequest(arg0, arg1 any) *MockRequestAdapterConvertToNativeRequestCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConvertToNativeRequest", reflect.TypeOf((*MockRequestAdapter)(nil).ConvertToNativeRequest), arg0, arg1)
	return &MockRequestAdapterConvertToNativeRequestCall{Call: call}
}

// MockRequestAdapterConvertToNativeRequestCall wrap *gomock.Call
type MockRequestAdapterConvertToNativeRequestCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRequestAdapterConvertToNativeRequestCall) Return(arg0 any, arg1 error) *MockRequestAdapterConvertToNativeRequestCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRequestAdapterConvertToNativeRequestCall) Do(f func(context.Context, *abstractions.RequestInformation) (any, error)) *MockRequestAdapterConvertToNativeRequestCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRequestAdapterConvertToNativeRequestCall) DoAndReturn(f func(context.Context, *abstractions.RequestInformation) (any, error)) *MockRequestAdapterConvertToNativeRequestCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// EnableBackingStore mocks base method.
func (m *MockRequestAdapter) EnableBackingStore(arg0 store.BackingStoreFactory) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "EnableBackingStore", arg0)
}

// EnableBackingStore indicates an expected call of EnableBackingStore.
func (mr *MockRequestAdapterMockRecorder) EnableBackingStore(arg0 any) *MockRequestAdapterEnableBackingStoreCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableBackingStore", reflect.TypeOf((*MockRequestAdapter)(nil).EnableBackingStore), arg0)
	return &MockRequestAdapterEnableBackingStoreCall{Call: call}
}

// MockRequestAdapterEnableBackingStoreCall wrap *gomock.Call
type MockRequestAdapterEnableBackingStoreCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRequestAdapterEnableBackingStoreCall) Return() *MockRequestAdapterEnableBackingStoreCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRequestAdapterEnableBackingStoreCall) Do(f func(store.BackingStoreFactory)) *MockRequestAdapterEnableBackingStoreCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRequestAdapterEnableBackingStoreCall) DoAndReturn(f func(store.BackingStoreFactory)) *MockRequestAdapterEnableBackingStoreCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetBaseUrl mocks base method.
func (m *MockRequestAdapter) GetBaseUrl() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBaseUrl")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetBaseUrl indicates an expected call of GetBaseUrl.
func (mr *MockRequestAdapterMockRecorder) GetBaseUrl() *MockRequestAdapterGetBaseUrlCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBaseUrl", reflect.TypeOf((*MockRequestAdapter)(nil).GetBaseUrl))
	return &MockRequestAdapterGetBaseUrlCall{Call: call}
}

// MockRequestAdapterGetBaseUrlCall wrap *gomock.Call
type MockRequestAdapterGetBaseUrlCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRequestAdapterGetBaseUrlCall) Return(arg0 string) *MockRequestAdapterGetBaseUrlCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRequestAdapterGetBaseUrlCall) Do(f func() string) *MockRequestAdapterGetBaseUrlCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRequestAdapterGetBaseUrlCall) DoAndReturn(f func() string) *MockRequestAdapterGetBaseUrlCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetSerializationWriterFactory mocks base method.
func (m *MockRequestAdapter) GetSerializationWriterFactory() serialization.SerializationWriterFactory {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSerializationWriterFactory")
	ret0, _ := ret[0].(serialization.SerializationWriterFactory)
	return ret0
}

// GetSerializationWriterFactory indicates an expected call of GetSerializationWriterFactory.
func (mr *MockRequestAdapterMockRecorder) GetSerializationWriterFactory() *MockRequestAdapterGetSerializationWriterFactoryCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSerializationWriterFactory", reflect.TypeOf((*MockRequestAdapter)(nil).GetSerializationWriterFactory))
	return &MockRequestAdapterGetSerializationWriterFactoryCall{Call: call}
}

// MockRequestAdapterGetSerializationWriterFactoryCall wrap *gomock.Call
type MockRequestAdapterGetSerializationWriterFactoryCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRequestAdapterGetSerializationWriterFactoryCall) Return(arg0 serialization.SerializationWriterFactory) *MockRequestAdapterGetSerializationWriterFactoryCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRequestAdapterGetSerializationWriterFactoryCall) Do(f func() serialization.SerializationWriterFactory) *MockRequestAdapterGetSerializationWriterFactoryCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRequestAdapterGetSerializationWriterFactoryCall) DoAndReturn(f func() serialization.SerializationWriterFactory) *MockRequestAdapterGetSerializationWriterFactoryCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Send mocks base method.
func (m *MockRequestAdapter) Send(arg0 context.Context, arg1 *abstractions.RequestInformation, arg2 serialization.ParsableFactory, arg3 abstractions.ErrorMappings) (serialization.Parsable, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(serialization.Parsable)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Send indicates an expected call of Send.
func (mr *MockRequestAdapterMockRecorder) Send(arg0, arg1, arg2, arg3 any) *MockRequestAdapterSendCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockRequestAdapter)(nil).Send), arg0, arg1, arg2, arg3)
	return &MockRequestAdapterSendCall{Call: call}
}

// MockRequestAdapterSendCall wrap *gomock.Call
type MockRequestAdapterSendCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRequestAdapterSendCall) Return(arg0 serialization.Parsable, arg1 error) *MockRequestAdapterSendCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRequestAdapterSendCall) Do(f func(context.Context, *abstractions.RequestInformation, serialization.ParsableFactory, abstractions.ErrorMappings) (serialization.Parsable, error)) *MockRequestAdapterSendCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRequestAdapterSendCall) DoAndReturn(f func(context.Context, *abstractions.RequestInformation, serialization.ParsableFactory, abstractions.ErrorMappings) (serialization.Parsable, error)) *MockRequestAdapterSendCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SendCollection mocks base method.
func (m *MockRequestAdapter) SendCollection(arg0 context.Context, arg1 *abstractions.RequestInformation, arg2 serialization.ParsableFactory, arg3 abstractions.ErrorMappings) ([]serialization.Parsable, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendCollection", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]serialization.Parsable)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendCollection indicates an expected call of SendCollection.
func (mr *MockRequestAdapterMockRecorder) SendCollection(arg0, arg1, arg2, arg3 any) *MockRequestAdapterSendCollectionCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendCollection", reflect.TypeOf((*MockRequestAdapter)(nil).SendCollection), arg0, arg1, arg2, arg3)
	return &MockRequestAdapterSendCollectionCall{Call: call}
}

// MockRequestAdapterSendCollectionCall wrap *gomock.Call
type MockRequestAdapterSendCollectionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRequestAdapterSendCollectionCall) Return(arg0 []serialization.Parsable, arg1 error) *MockRequestAdapterSendCollectionCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRequestAdapterSendCollectionCall) Do(f func(context.Context, *abstractions.RequestInformation, serialization.ParsableFactory, abstractions.ErrorMappings) ([]serialization.Parsable, error)) *MockRequestAdapterSendCollectionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRequestAdapterSendCollectionCall) DoAndReturn(f func(context.Context, *abstractions.RequestInformation, serialization.ParsableFactory, abstractions.ErrorMappings) ([]serialization.Parsable, error)) *MockRequestAdapterSendCollectionCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SendEnum mocks base method.
func (m *MockRequestAdapter) SendEnum(arg0 context.Context, arg1 *abstractions.RequestInformation, arg2 serialization.EnumFactory, arg3 abstractions.ErrorMappings) (any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEnum", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendEnum indicates an expected call of SendEnum.
func (mr *MockRequestAdapterMockRecorder) SendEnum(arg0, arg1, arg2, arg3 any) *MockRequestAdapterSendEnumCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEnum", reflect.TypeOf((*MockRequestAdapter)(nil).SendEnum), arg0, arg1, arg2, arg3)
	return &MockRequestAdapterSendEnumCall{Call: call}
}

// MockRequestAdapterSendEnumCall wrap *gomock.Call
type MockRequestAdapterSendEnumCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRequestAdapterSendEnumCall) Return(arg0 any, arg1 error) *MockRequestAdapterSendEnumCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRequestAdapterSendEnumCall) Do(f func(context.Context, *abstractions.RequestInformation, serialization.EnumFactory, abstractions.ErrorMappings) (any, error)) *MockRequestAdapterSendEnumCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRequestAdapterSendEnumCall) DoAndReturn(f func(context.Context, *abstractions.RequestInformation, serialization.EnumFactory, abstractions.ErrorMappings) (any, error)) *MockRequestAdapterSendEnumCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SendEnumCollection mocks base method.
func (m *MockRequestAdapter) SendEnumCollection(arg0 context.Context, arg1 *abstractions.RequestInformation, arg2 serialization.EnumFactory, arg3 abstractions.ErrorMappings) ([]any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEnumCollection", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendEnumCollection indicates an expected call of SendEnumCollection.
func (mr *MockRequestAdapterMockRecorder) SendEnumCollection(arg0, arg1, arg2, arg3 any) *MockRequestAdapterSendEnumCollectionCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEnumCollection", reflect.TypeOf((*MockRequestAdapter)(nil).SendEnumCollection), arg0, arg1, arg2, arg3)
	return &MockRequestAdapterSendEnumCollectionCall{Call: call}
}

// MockRequestAdapterSendEnumCollectionCall wrap *gomock.Call
type MockRequestAdapterSendEnumCollectionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRequestAdapterSendEnumCollectionCall) Return(arg0 []any, arg1 error) *MockRequestAdapterSendEnumCollectionCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRequestAdapterSendEnumCollectionCall) Do(f func(context.Context, *abstractions.RequestInformation, serialization.EnumFactory, abstractions.ErrorMappings) ([]any, error)) *MockRequestAdapterSendEnumCollectionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRequestAdapterSendEnumCollectionCall) DoAndReturn(f func(context.Context, *abstractions.RequestInformation, serialization.EnumFactory, abstractions.ErrorMappings) ([]any, error)) *MockRequestAdapterSendEnumCollectionCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SendNoContent mocks base method.
func (m *MockRequestAdapter) SendNoContent(arg0 context.Context, arg1 *abstractions.RequestInformation, arg2 abstractions.ErrorMappings) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendNoContent", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendNoContent indicates an expected call of SendNoContent.
func (mr *MockRequestAdapterMockRecorder) SendNoContent(arg0, arg1, arg2 any) *MockRequestAdapterSendNoContentCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendNoContent", reflect.TypeOf((*MockRequestAdapter)(nil).SendNoContent), arg0, arg1, arg2)
	return &MockRequestAdapterSendNoContentCall{Call: call}
}

// MockRequestAdapterSendNoContentCall wrap *gomock.Call
type MockRequestAdapterSendNoContentCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRequestAdapterSendNoContentCall) Return(arg0 error) *MockRequestAdapterSendNoContentCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRequestAdapterSendNoContentCall) Do(f func(context.Context, *abstractions.RequestInformation, abstractions.ErrorMappings) error) *MockRequestAdapterSendNoContentCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRequestAdapterSendNoContentCall) DoAndReturn(f func(context.Context, *abstractions.RequestInformation, abstractions.ErrorMappings) error) *MockRequestAdapterSendNoContentCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SendPrimitive mocks base method.
func (m *MockRequestAdapter) SendPrimitive(arg0 context.Context, arg1 *abstractions.RequestInformation, arg2 string, arg3 abstractions.ErrorMappings) (any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendPrimitive", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendPrimitive indicates an expected call of SendPrimitive.
func (mr *MockRequestAdapterMockRecorder) SendPrimitive(arg0, arg1, arg2, arg3 any) *MockRequestAdapterSendPrimitiveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendPrimitive", reflect.TypeOf((*MockRequestAdapter)(nil).SendPrimitive), arg0, arg1, arg2, arg3)
	return &MockRequestAdapterSendPrimitiveCall{Call: call}
}

// MockRequestAdapterSendPrimitiveCall wrap *gomock.Call
type MockRequestAdapterSendPrimitiveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRequestAdapterSendPrimitiveCall) Return(arg0 any, arg1 error) *MockRequestAdapterSendPrimitiveCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRequestAdapterSendPrimitiveCall) Do(f func(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) (any, error)) *MockRequestAdapterSendPrimitiveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRequestAdapterSendPrimitiveCall) DoAndReturn(f func(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) (any, error)) *MockRequestAdapterSendPrimitiveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SendPrimitiveCollection mocks base method.
func (m *MockRequestAdapter) SendPrimitiveCollection(arg0 context.Context, arg1 *abstractions.RequestInformation, arg2 string, arg3 abstractions.ErrorMappings) ([]any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendPrimitiveCollection", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendPrimitiveCollection indicates an expected call of SendPrimitiveCollection.
func (mr *MockRequestAdapterMockRecorder) SendPrimitiveCollection(arg0, arg1, arg2, arg3 any) *MockRequestAdapterSendPrimitiveCollectionCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendPrimitiveCollection", reflect.TypeOf((*MockRequestAdapter)(nil).SendPrimitiveCollection), arg0, arg1, arg2, arg3)
	return &MockRequestAdapterSendPrimitiveCollectionCall{Call: call}
}

// MockRequestAdapterSendPrimitiveCollectionCall wrap *gomock.Call
type MockRequestAdapterSendPrimitiveCollectionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRequestAdapterSendPrimitiveCollectionCall) Return(arg0 []any, arg1 error) *MockRequestAdapterSendPrimitiveCollectionCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRequestAdapterSendPrimitiveCollectionCall) Do(f func(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) ([]any, error)) *MockRequestAdapterSendPrimitiveCollectionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRequestAdapterSendPrimitiveCollectionCall) DoAndReturn(f func(context.Context, *abstractions.RequestInformation, string, abstractions.ErrorMappings) ([]any, error)) *MockRequestAdapterSendPrimitiveCollectionCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetBaseUrl mocks base method.
func (m *MockRequestAdapter) SetBaseUrl(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetBaseUrl", arg0)
}

// SetBaseUrl indicates an expected call of SetBaseUrl.
func (mr *MockRequestAdapterMockRecorder) SetBaseUrl(arg0 any) *MockRequestAdapterSetBaseUrlCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBaseUrl", reflect.TypeOf((*MockRequestAdapter)(nil).SetBaseUrl), arg0)
	return &MockRequestAdapterSetBaseUrlCall{Call: call}
}

// MockRequestAdapterSetBaseUrlCall wrap *gomock.Call
type MockRequestAdapterSetBaseUrlCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRequestAdapterSetBaseUrlCall) Return() *MockRequestAdapterSetBaseUrlCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRequestAdapterSetBaseUrlCall) Do(f func(string)) *MockRequestAdapterSetBaseUrlCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRequestAdapterSetBaseUrlCall) DoAndReturn(f func(string)) *MockRequestAdapterSetBaseUrlCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}