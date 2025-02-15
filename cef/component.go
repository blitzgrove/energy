//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

// CEF 扩展组件

package cef

// NewTray 适用于 windows linux macos 系统托盘
func (m *LCLBrowserWindow) NewTray() ITray {
	if m.tray == nil {
		m.tray = newTray(m.TForm)
	}
	return m.tray
}

// NewSysTray LCL窗口组件,系统托盘
func (m *LCLBrowserWindow) NewSysTray() ITray {
	if m.tray == nil {
		m.tray = newSysTray()
	}
	return m.tray
}

// NewSysTray VF窗口组件,只适用于windows的无菜单托盘
func (m *ViewsFrameworkBrowserWindow) NewSysTray() ITray {
	if m == nil {
		return nil
	}
	if m.tray == nil {
		m.tray = newSysTray()
	}
	return m.tray
}
