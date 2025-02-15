//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

package cef

import (
	"github.com/energye/energy/v2/cef/internal/def"
	"github.com/energye/energy/v2/common/imports"
	"github.com/energye/energy/v2/consts"
	"github.com/energye/golcl/lcl"
	"github.com/energye/golcl/lcl/api"
	"unsafe"
)

// MenuButtonDelegateRef -> ICefMenuModelDelegate
var MenuButtonDelegateRef menuButtonDelegate

type menuButtonDelegate uintptr

func (*menuButtonDelegate) New() *ICefMenuButtonDelegate {
	var result uintptr
	imports.Proc(def.MenuButtonDelegateRef_Create).Call(uintptr(unsafe.Pointer(&result)))
	if result != 0 {
		return &ICefMenuButtonDelegate{&ICefButtonDelegate{&ICefViewDelegate{
			instance: getInstance(result),
		}}}
	}
	return nil
}

func (*menuButtonDelegate) NewForCustom(menuButtonComponent *TCEFMenuButtonComponent) *ICefMenuButtonDelegate {
	var result uintptr
	imports.Proc(def.MenuButtonDelegateRef_CreateForCustom).Call(menuButtonComponent.Instance(), uintptr(unsafe.Pointer(&result)))
	if result != 0 {
		return &ICefMenuButtonDelegate{&ICefButtonDelegate{&ICefViewDelegate{
			instance: getInstance(result),
			ct:       consts.CtOther,
		}}}
	}
	return nil
}

func (m *ICefMenuButtonDelegate) SetOnMenuButtonPressed(fn onMenuButtonPressed) {
	if !m.IsValid() || m.IsOtherEvent() {
		return
	}
	imports.Proc(def.MenuButtonDelegate_SetOnMenuButtonPressed).Call(m.Instance(), api.MakeEventDataPtr(fn))
}

type onMenuButtonPressed func(button *ICefMenuButton, screenPoint *TCefPoint, buttonPressedLock *ICefMenuButtonPressedLock)

func init() {
	lcl.RegisterExtEventCallback(func(fn interface{}, getVal func(idx int) uintptr) bool {
		getPtr := func(i int) unsafe.Pointer {
			return unsafe.Pointer(getVal(i))
		}
		switch fn.(type) {
		case onMenuButtonPressed:
			button := &ICefMenuButton{&ICefLabelButton{&ICefButton{&ICefView{instance: getPtr(0)}}}}
			screenPoint := (*TCefPoint)(getPtr(1))
			buttonPressedLock := &ICefMenuButtonPressedLock{base: TCefBaseRefCounted{instance: getPtr(3)}}
			fn.(onMenuButtonPressed)(button, screenPoint, buttonPressedLock)
		default:
			return false
		}
		return true
	})
}
