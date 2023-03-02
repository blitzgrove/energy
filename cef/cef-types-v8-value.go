//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

// CEF v8 value V8所有类型
package cef

import (
	"github.com/energye/energy/common"
	"github.com/energye/energy/common/imports"
	"github.com/energye/energy/consts"
	"github.com/energye/golcl/lcl/api"
	"time"
	"unsafe"
)

// isInternalKey 内部key不允许使用
func isInternalKey(key string) bool {
	return key == internalIPCKey ||
		key == internalEmit ||
		key == internalOn
}

func (m *ICefV8Value) Instance() uintptr {
	if m == nil {
		return 0
	}
	return uintptr(m.instance)
}

func (m *ICefV8Value) IsValid() bool {
	if m == nil || m.instance == nil {
		return false
	}
	r1, _, _ := imports.Proc(internale_CefV8Value_IsValid).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) IsUndefined() bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_IsUndefined).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) IsNull() bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_IsNull).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) IsBool() bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_IsBool).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) IsInt() bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_IsInt).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) IsUInt() bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_IsUInt).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) IsDouble() bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_IsDouble).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) IsDate() bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_IsDate).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) IsString() bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_IsString).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) IsObject() bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_IsObject).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) IsArray() bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_IsArray).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) IsArrayBuffer() bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_IsArrayBuffer).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) IsFunction() bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_IsFunction).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) IsPromise() bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_IsPromise).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) IsSame() bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_IsSame).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) GetBoolValue() bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_GetBoolValue).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) GetIntValue() int32 {
	r1, _, _ := imports.Proc(internale_CefV8Value_GetIntValue).Call(m.Instance())
	return int32(r1)
}

func (m *ICefV8Value) GetUIntValue() uint32 {
	r1, _, _ := imports.Proc(internale_CefV8Value_GetIntValue).Call(m.Instance())
	return uint32(r1)
}

func (m *ICefV8Value) GetDoubleValue() (result float64) {
	imports.Proc(internale_CefV8Value_GetDoubleValue).Call(m.Instance(), uintptr(unsafe.Pointer(&result)))
	return result
}

func (m *ICefV8Value) GetDateValue() time.Time {
	var result float64
	imports.Proc(internale_CefV8Value_GetDateValue).Call(m.Instance(), uintptr(unsafe.Pointer(&result)))
	return common.DDateTimeToGoDateTime(result)
}

func (m *ICefV8Value) GetStringValue() string {
	r1, _, _ := imports.Proc(internale_CefV8Value_GetStringValue).Call(m.Instance())
	return api.GoStr(r1)
}

func (m *ICefV8Value) IsUserCreated() bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_IsUserCreated).Call(m.Instance())
	return api.GoBool(r1)
}

//func (m *ICefV8Value) HasException() bool {
//	r1, _, _ := imports.Proc(internale_CefV8Value_HasException).Call(m.Instance())
//	return api.GoBool(r1)
//}

//func (m *ICefV8Value) GetException() {
//
//}

func (m *ICefV8Value) ClearException() bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_ClearException).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) WillRethrowExceptions() bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_WillRethrowExceptions).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) SetRethrowExceptions(reThrow bool) bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_SetRethrowExceptions).Call(m.Instance(), api.PascalBool(reThrow))
	return api.GoBool(r1)
}

func (m *ICefV8Value) HasValueByKey(key string) bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_HasValueByKey).Call(m.Instance(), api.PascalStr(key))
	return api.GoBool(r1)
}

func (m *ICefV8Value) HasValueByIndex(index int32) bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_HasValueByIndex).Call(m.Instance(), uintptr(index))
	return api.GoBool(r1)
}

// deleteValueByKey internal
func (m *ICefV8Value) deleteValueByKey(key string) bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_DeleteValueByKey).Call(m.Instance(), api.PascalStr(key))
	return api.GoBool(r1)
}

// DeleteValueByKey export
func (m *ICefV8Value) DeleteValueByKey(key string) bool {
	if isInternalKey(key) {
		return false
	}
	return m.deleteValueByKey(key)
}

func (m *ICefV8Value) DeleteValueByIndex(index int32) bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_DeleteValueByIndex).Call(m.Instance(), uintptr(index))
	return api.GoBool(r1)
}

// getValueByKey internal
func (m *ICefV8Value) getValueByKey(key string) *ICefV8Value {
	var result uintptr
	imports.Proc(internale_CefV8Value_GetValueByKey).Call(m.Instance(), api.PascalStr(key), uintptr(unsafe.Pointer(&result)))
	return &ICefV8Value{
		instance: unsafe.Pointer(result),
	}
}

// GetValueByKey export
func (m *ICefV8Value) GetValueByKey(key string) *ICefV8Value {
	if key == internalIPCKey {
		return nil
	}
	return m.getValueByKey(key)
}

func (m *ICefV8Value) GetValueByIndex(index int32) *ICefV8Value {
	var result uintptr
	imports.Proc(internale_CefV8Value_GetValueByIndex).Call(m.Instance(), uintptr(index), uintptr(unsafe.Pointer(&result)))
	return &ICefV8Value{
		instance: unsafe.Pointer(result),
	}
}

// setValueByKey internal
func (m *ICefV8Value) setValueByKey(key string, value *ICefV8Value, attribute consts.TCefV8PropertyAttributes) bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_SetValueByKey).Call(m.Instance(), api.PascalStr(key), value.Instance(), attribute.ToPtr())
	return api.GoBool(r1)
}

// SetValueByKey export
func (m *ICefV8Value) SetValueByKey(key string, value *ICefV8Value, attribute consts.TCefV8PropertyAttributes) bool {
	if isInternalKey(key) {
		return false
	}
	return m.setValueByKey(key, value, attribute)
}

func (m *ICefV8Value) SetValueByIndex(index int32, value *ICefV8Value) bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_SetValueByIndex).Call(m.Instance(), uintptr(index), value.Instance())
	return api.GoBool(r1)
}

// SetValueByAccessor internal
func (m *ICefV8Value) setValueByAccessor(key string, settings consts.TCefV8AccessControls, attribute consts.TCefV8PropertyAttributes) bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_SetValueByAccessor).Call(m.Instance(), api.PascalStr(key), settings.ToPtr(), attribute.ToPtr())
	return api.GoBool(r1)
}

// SetValueByAccessor export
func (m *ICefV8Value) SetValueByAccessor(key string, settings consts.TCefV8AccessControls, attribute consts.TCefV8PropertyAttributes) bool {
	if isInternalKey(key) {
		return false
	}
	return m.setValueByAccessor(key, settings, attribute)
}

//func (m *ICefV8Value) GetKeys(keys *lcl.TStrings) int32 {
//	r1, _, _ := imports.Proc(internale_CefV8Value_GetKeys).Call(m.Instance(), keys.Instance())
//	return int32(r1)
//}

func (m *ICefV8Value) SetUserData(data *ICefV8Value) bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_SetUserData).Call(m.Instance(), data.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) GetUserData() *ICefV8Value {
	var result uintptr
	imports.Proc(internale_CefV8Value_GetUserData).Call(m.Instance(), uintptr(unsafe.Pointer(&result)))
	return &ICefV8Value{
		instance: unsafe.Pointer(result),
	}
}

func (m *ICefV8Value) GetExternallyAllocatedMemory() int32 {
	r1, _, _ := imports.Proc(internale_CefV8Value_GetExternallyAllocatedMemory).Call(m.Instance())
	return int32(r1)
}

func (m *ICefV8Value) AdjustExternallyAllocatedMemory(changeInBytes int32) int32 {
	r1, _, _ := imports.Proc(internale_CefV8Value_AdjustExternallyAllocatedMemory).Call(m.Instance(), uintptr(changeInBytes))
	return int32(r1)
}

func (m *ICefV8Value) GetArrayLength() int32 {
	r1, _, _ := imports.Proc(internale_CefV8Value_GetArrayLength).Call(m.Instance())
	return int32(r1)
}

//func (m *ICefV8Value) GetArrayBufferReleaseCallback() {
//
//}

func (m *ICefV8Value) NeuterArrayBuffer() bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_NeuterArrayBuffer).Call(m.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) GetFunctionName() string {
	r1, _, _ := imports.Proc(internale_CefV8Value_GetFunctionName).Call(m.Instance())
	return api.GoStr(r1)
}

func (m *ICefV8Value) GetFunctionHandler() *ICefV8Handler {
	var result uintptr
	imports.Proc(internale_CefV8Value_GetFunctionHandler).Call(m.Instance(), uintptr(unsafe.Pointer(&result)))
	return &ICefV8Handler{
		instance: unsafe.Pointer(result),
	}
}

func (m *ICefV8Value) ExecuteFunction(obj *ICefV8Value, arguments []*ICefV8Value) *ICefV8Value {
	var result uintptr
	imports.Proc(internale_CefV8Value_ExecuteFunction).Call(m.Instance(), obj.Instance(), uintptr(unsafe.Pointer(&result)), uintptr(unsafe.Pointer(arguments[0])), uintptr(int32(len(arguments))))
	return &ICefV8Value{
		instance: unsafe.Pointer(result),
	}
}

func (m *ICefV8Value) ExecuteFunctionWithContext(v8Context *ICefV8Context, obj *ICefV8Value, arguments []*ICefV8Value) *ICefV8Value {
	var result uintptr
	imports.Proc(internale_CefV8Value_ExecuteFunctionWithContext).Call(m.Instance(), v8Context.Instance(), obj.Instance(), uintptr(unsafe.Pointer(&result)), uintptr(unsafe.Pointer(arguments[0])), uintptr(int32(len(arguments))))
	return &ICefV8Value{
		instance: unsafe.Pointer(result),
	}
}

func (m *ICefV8Value) ResolvePromise(arg *ICefV8Value) bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_ResolvePromise).Call(m.Instance(), arg.Instance())
	return api.GoBool(r1)
}

func (m *ICefV8Value) RejectPromise(errorMsg string) bool {
	r1, _, _ := imports.Proc(internale_CefV8Value_RejectPromise).Call(m.Instance(), api.PascalStr(errorMsg))
	return api.GoBool(r1)
}

func (m *ICefV8Value) Free() {
	m.instance = nil
}

// ResultV8Value 返回 ICefV8Value 的替代结构
type ResultV8Value struct {
	v8value *ICefV8Value
}

// SetResult 设置 ICefV8Value 返回值
func (m *ResultV8Value) SetResult(v8value *ICefV8Value) {
	m.v8value = v8value
}

// Get 根据下标获取 ICefV8Value
func (m *TCefV8ValueArray) Get(index int) *ICefV8Value {
	if index < m.argumentsLength {
		return &ICefV8Value{instance: unsafe.Pointer(common.GetParamOf(index, m.arguments))}
	}
	return nil
}

// Size 返回 ICefV8Value 数组长度
func (m *TCefV8ValueArray) Size() int {
	return m.argumentsLength
}

// V8ValueRef -> ICefV8Value
var V8ValueRef cefV8Value

// cefV8Value
type cefV8Value uintptr

func (*cefV8Value) NewUndefined() *ICefV8Value {
	var result uintptr
	imports.Proc(internale_CefV8ValueRef_NewUndefined).Call(uintptr(unsafe.Pointer(&result)))
	return &ICefV8Value{
		instance: unsafe.Pointer(result),
	}
}

func (*cefV8Value) NewNull() *ICefV8Value {
	var result uintptr
	imports.Proc(internale_CefV8ValueRef_NewNull).Call(uintptr(unsafe.Pointer(&result)))
	return &ICefV8Value{
		instance: unsafe.Pointer(result),
	}
}

func (*cefV8Value) NewBool(value bool) *ICefV8Value {
	var result uintptr
	imports.Proc(internale_CefV8ValueRef_NewBool).Call(api.PascalBool(value), uintptr(unsafe.Pointer(&result)))
	return &ICefV8Value{
		instance: unsafe.Pointer(result),
	}
}
func (*cefV8Value) NewInt(value int32) *ICefV8Value {
	var result uintptr
	imports.Proc(internale_CefV8ValueRef_NewInt).Call(uintptr(value), uintptr(unsafe.Pointer(&result)))
	return &ICefV8Value{
		instance: unsafe.Pointer(result),
	}
}

func (*cefV8Value) NewUInt(value uint32) *ICefV8Value {
	var result uintptr
	imports.Proc(internale_CefV8ValueRef_NewUInt).Call(uintptr(value), uintptr(unsafe.Pointer(&result)))
	return &ICefV8Value{
		instance: unsafe.Pointer(result),
	}
}

func (*cefV8Value) NewDouble(value float64) *ICefV8Value {
	var result uintptr
	imports.Proc(internale_CefV8ValueRef_NewDouble).Call(uintptr(unsafe.Pointer(&value)), uintptr(unsafe.Pointer(&result)))
	return &ICefV8Value{
		instance: unsafe.Pointer(result),
	}
}

func (*cefV8Value) NewDate(value time.Time) *ICefV8Value {
	var result uintptr
	val := common.GoDateTimeToDDateTime(value)
	imports.Proc(internale_CefV8ValueRef_NewDate).Call(uintptr(unsafe.Pointer(&val)), uintptr(unsafe.Pointer(&result)))
	return &ICefV8Value{
		instance: unsafe.Pointer(result),
	}
}

func (*cefV8Value) NewString(value string) *ICefV8Value {
	var result uintptr
	imports.Proc(internale_CefV8ValueRef_NewString).Call(api.PascalStr(value), uintptr(unsafe.Pointer(&result)))
	return &ICefV8Value{
		instance: unsafe.Pointer(result),
	}
}

//func (*cefV8Value) NewObject(accessor *ICefV8Accessor, interceptor *ICefV8Interceptor) *ICefV8Value {
func (*cefV8Value) NewObject(accessor *ICefV8Accessor) *ICefV8Value {
	var result uintptr
	imports.Proc(internale_CefV8ValueRef_NewObject).Call(accessor.Instance(), uintptr(0), uintptr(unsafe.Pointer(&result)))
	return &ICefV8Value{
		instance: unsafe.Pointer(result),
	}
}

func (*cefV8Value) NewArray(len int32) *ICefV8Value {
	var result uintptr
	imports.Proc(internale_CefV8ValueRef_NewArray).Call(uintptr(len), uintptr(unsafe.Pointer(&result)))
	return &ICefV8Value{
		instance: unsafe.Pointer(result),
	}
}

func (*cefV8Value) NewArrayBuffer(buffer []byte, callback *ICefV8ArrayBufferReleaseCallback) *ICefV8Value {
	var bufLen = len(buffer)
	if bufLen == 0 {
		return nil
	}
	if callback == nil {
		callback = V8ArrayBufferReleaseCallbackRef.New()
	}
	var result uintptr
	imports.Proc(internale_CefV8ValueRef_NewArrayBuffer).Call(uintptr(unsafe.Pointer(&buffer[0])), uintptr(int32(bufLen)), callback.Instance(), uintptr(unsafe.Pointer(&result)))
	return &ICefV8Value{
		instance: unsafe.Pointer(result),
	}
}

// newFunction internal
func (*cefV8Value) newFunction(name string, handler *ICefV8Handler) *ICefV8Value {
	var result uintptr
	imports.Proc(internale_CefV8ValueRef_NewFunction).Call(api.PascalStr(name), handler.Instance(), uintptr(unsafe.Pointer(&result)))
	return &ICefV8Value{
		instance: unsafe.Pointer(result),
	}
}

// NewFunction export
func (m *cefV8Value) NewFunction(name string, handler *ICefV8Handler) *ICefV8Value {
	if isInternalKey(name) {
		return nil
	}
	return m.newFunction(name, handler)
}

func (*cefV8Value) NewPromise() *ICefV8Value {
	var result uintptr
	imports.Proc(internale_CefV8ValueRef_NewPromise).Call(uintptr(unsafe.Pointer(&result)))
	return &ICefV8Value{
		instance: unsafe.Pointer(result),
	}
}
