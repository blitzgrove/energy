//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

// energy ipc 事件驱动监听
//
// 事件监听、事件触发
package ipc

import (
	"github.com/energye/energy/common"
	"reflect"
	"sync"
)

const (
	Ln_IPC_GoEmitJS         = "IPCGoEmitJS"                   //Go执行Js on监听
	Ln_GET_BIND_FIELD_VALUE = "internal_GET_BIND_FIELD_VALUE" //browse进程监听获取字段值
	Ln_SET_BIND_FIELD_VALUE = "internal_SET_BIND_FIELD_VALUE" //browse进程监听设置字段值
	Ln_EXECUTE_BIND_FUNC    = "internal_EXECUTE_BIND_FUNC"    //browse进程监听执行绑定函数
	Ln_onConnectEvent       = "connect"
)

func InternalIPCNameCheck(name string) bool {
	if name == Ln_IPC_GoEmitJS || name == Ln_GET_BIND_FIELD_VALUE || name == Ln_SET_BIND_FIELD_VALUE || name == Ln_EXECUTE_BIND_FUNC || name == Ln_onConnectEvent {
		return true
	}
	return false
}

var (
	browser *browserIPC
)

// EmitContextCallback IPC 上下文件回调函数
type EmitContextCallback func(context IContext)

// EmitArgumentCallback 带有参数回调函数
type EmitArgumentCallback any

// browserIPC 主进程 IPC
type browserIPC struct {
	onEvent map[string]*callback
	lock    sync.Mutex
}

func init() {
	if common.Args.IsMain() {
		browser = &browserIPC{onEvent: make(map[string]*callback)}
	}
}

//On
// IPC GO 监听事件, 上下文参数
func On(name string, fn EmitContextCallback) {
	if browser != nil {
		browser.addOnEvent(name, &callback{context: &contextCallback{callback: fn}})
	}
}

//OnArguments
// IPC GO 监听事件, 自定义参数，仅支持对应 JavaScript 对应 Go 的常用类型
//
// 入参 不限制个数
//	 基本类型: int(int8 ~ uint64), bool, float(float32、float64), string
//
//   复杂类型: slice, map, struct
//
//   复杂类型限制示例: slice: []interface{}, map: map[string]interface{}
//
//   slice: 只能是 any 或 interface{}
//   map: key 只能 string 类型, value 基本类型+复杂类型
//   struct: 首字母大写, 字段类型匹配
//     type ArgsStructDemo struct {
//        Key1 string
//		  Key2 string
//		  Key3 string
//		  Key4 int
//		  Key5 float64
//		  Key6 bool
//		  Sub1  SubStructXXX
//		  Sub2  *SubStructXXX
//     }
//
// 出参
//
//   与同入参相同
func OnArguments(name string, fn EmitArgumentCallback) {
	if browser != nil {
		v := reflect.ValueOf(fn)
		// f must be a function
		if v.Kind() != reflect.Func {
			return
		}
		browser.addOnEvent(name, &callback{argument: &argumentCallback{callback: &v}})
	}
}

//RemoveOn
// IPC GO 移除监听事件
func RemoveOn(name string) {
	if browser == nil || name == "" {
		return
	}
	browser.lock.Lock()
	defer browser.lock.Unlock()
	delete(browser.onEvent, name)
}

//Emit
// IPC GO 中触发 JS 监听的事件
func Emit(name string, argument ...any) {

}

//
////EmitTarget
//// IPC GO 中触发指定目标 JS 监听的事件
//func EmitTarget(name string, argument ...any) {
//
//}

//CheckOnEvent
// IPC 检查 GO 中监听的事件是存在, 并返回回调函数
func CheckOnEvent(name string) *callback {
	if browser == nil || name == "" {
		return nil
	}
	browser.lock.Lock()
	defer browser.lock.Unlock()
	if fn, ok := browser.onEvent[name]; ok {
		return fn
	}
	return nil
}

// addOnEvent 添加监听事件
func (m *browserIPC) addOnEvent(name string, fn *callback) {
	if m == nil || name == "" || fn == nil {
		return
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	m.onEvent[name] = fn
}

// emitOnEvent 触发监听事件
func (m *browserIPC) emitOnEvent(name string, argumentList IArrayValue) {
	if m == nil || name == "" || argumentList == nil {
		return
	}
	m.lock.Lock()
	defer m.lock.Unlock()

}
