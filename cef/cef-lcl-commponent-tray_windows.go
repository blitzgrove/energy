//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under GNU General Public License v3.0
//
//----------------------------------------

//go:build windows
// +build windows

package cef

import (
	. "github.com/energye/energy/common"
	"github.com/energye/energy/common/assetserve"
	. "github.com/energye/energy/consts"
	"github.com/energye/energy/ipc"
	"github.com/energye/energy/logger"
	"github.com/energye/golcl/lcl"
	"github.com/energye/golcl/lcl/types"
)

//Cef托盘
type tLCLTrayWindow struct {
	*lcl.TForm
	owner        lcl.IComponent
	trayIcon     *lcl.TTrayIcon
	chromium     IChromium
	windowParent ITCefWindowParent
	x, y, w, h   int32
	mouseUp      TMouseEvent
	isClosing    bool
	url          string
}

func newLCLTrayWindow(owner lcl.IComponent, width, height int32, url string) *tLCLTrayWindow {
	var trayForm *tLCLTrayWindow
	lcl.Application.CreateForm(&trayForm)
	trayForm.trayIcon = lcl.NewTrayIcon(owner)
	trayForm.trayIcon.SetVisible(true)
	trayForm.owner = owner
	trayForm.x = -width
	trayForm.y = -height
	trayForm.w = width
	trayForm.h = height
	trayForm.url = url
	trayForm.onmMouse()
	trayForm.createCefTrayWindow()
	return trayForm
}

func (m *tLCLTrayWindow) OnFormCreate(sender lcl.IObject) {
	m.SetShowInTaskBar(types.StNever)
}

func (m *tLCLTrayWindow) Tray() *Tray {
	return nil
}

func (m *tLCLTrayWindow) Show() {
	if BrowserWindow.mainBrowserWindow.Chromium() == nil || !BrowserWindow.mainBrowserWindow.Chromium().Initialized() {
		return
	}
	m.TForm.Show()
}

func (m *tLCLTrayWindow) Hide() {
	m.TForm.Hide()
}

func (m *tLCLTrayWindow) close() {
	if m.isClosing {
		return
	}
	m.Hide()
	m.TForm.Close()
}

func (m *tLCLTrayWindow) SetOnDblClick(fn lcl.TNotifyEvent) {
	m.trayIcon.SetOnDblClick(fn)
}

func (m *tLCLTrayWindow) SetOnClick(fn lcl.TNotifyEvent) {
	m.trayIcon.SetOnClick(fn)
}

func (m *tLCLTrayWindow) SetOnMouseUp(fn TMouseEvent) {
	m.mouseUp = fn
}
func (m *tLCLTrayWindow) SetOnMouseDown(fn lcl.TMouseEvent) {
	m.trayIcon.SetOnMouseDown(fn)
}
func (m *tLCLTrayWindow) SetOnMouseMove(fn lcl.TMouseMoveEvent) {
	m.trayIcon.SetOnMouseMove(fn)
}

func (m *tLCLTrayWindow) Visible() bool {
	return m.TForm.Visible()
}

func (m *tLCLTrayWindow) SetVisible(v bool) {
	m.trayIcon.SetVisible(v)
}

func (m *tLCLTrayWindow) SetHint(value string) {
	m.trayIcon.SetHint(value)
}

func (m *tLCLTrayWindow) SetTitle(title string) {
	m.TForm.SetCaption(title)
}

func (m *tLCLTrayWindow) onmMouse() {
	QueueAsyncCall(func(id int) {
		m.trayIcon.SetOnMouseUp(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) {
			var monitor = m.TForm.Monitor()
			var monitorWidth = monitor.Width()
			width, height := m.TForm.Width(), m.TForm.Height()
			var mx = x + width
			var my = y + height
			if mx < monitorWidth {
				mx = x
			} else {
				mx = x - width
			}
			if my > m.h {
				my = y
			}
			if my > height {
				my = y - height
			}
			m.TForm.SetBounds(mx, my, width, height)
			var ret bool
			if m.mouseUp != nil {
				ret = m.mouseUp(sender, button, shift, x, y)
			}
			if !ret {
				if button == types.MbRight {
					m.Show()
				}
			}
		})
	})
}

//设置托盘气泡
//title 气泡标题
//content 气泡内容
//timeout 显示时间(毫秒)
func (m *tLCLTrayWindow) SetBalloon(title, content string, timeout int32) ITray {
	m.trayIcon.SetBalloonTitle(title)
	m.trayIcon.SetBalloonHint(content)
	m.trayIcon.SetBalloonTimeout(timeout)
	return m
}

//显示托盘气泡
func (m *tLCLTrayWindow) ShowBalloon() {
	m.trayIcon.ShowBalloonHint()
}

func (m *tLCLTrayWindow) createCefTrayWindow() {
	m.TForm.SetBorderStyle(types.BsNone)
	m.TForm.SetFormStyle(types.FsStayOnTop)
	m.TForm.SetBounds(-(m.w * 2), -(m.h * 2), m.w, m.h)
	m.TForm.SetOnActivate(func(sender lcl.IObject) {
		m.chromium.Initialized()
		m.chromium.CreateBrowser(m.windowParent)
	})
	m.TForm.SetOnWndProc(func(msg *types.TMessage) {
		m.TForm.InheritedWndProc(msg)
		if msg.Msg == 6 && msg.WParam == 0 && msg.LParam == 0 {
			QueueAsyncCall(func(id int) {
				if m.isClosing {
					return
				}
				m.TForm.Hide()
			})
		}
	})
	m.TForm.SetOnDeactivate(func(sender lcl.IObject) {
		if m.isClosing {
			return
		}
		m.TForm.Hide()
	})

	m.TForm.SetOnCloseQuery(func(sender lcl.IObject, canClose *bool) {
		*canClose = true
		logger.Debug("tray.window.onCloseQuery canClose:", *canClose)
		if m.isClosing {
			return
		}
		m.isClosing = true
		m.Hide()
		m.chromium.CloseBrowser(true)
		m.trayIcon.Free()
	})
	m.TForm.SetOnClose(func(sender lcl.IObject, action *types.TCloseAction) {
		*action = types.CaFree
		logger.Debug("tray.window.onClose action:", *action)
	})
	m.TForm.SetOnShow(func(sender lcl.IObject) {
		if m.windowParent != nil {
			QueueAsyncCall(func(id int) {
				m.windowParent.UpdateSize()
			})
		}
	})
	m.windowParent = NewCEFWindow(m.TForm)
	m.windowParent.SetParent(m.TForm)
	m.windowParent.SetAlign(types.AlClient)
	m.windowParent.SetAnchors(types.NewSet(types.AkTop, types.AkLeft, types.AkRight, types.AkBottom))
	m.chromium = NewChromium(m.windowParent, nil)
	m.chromium.SetOnBeforeContextMenu(func(sender lcl.IObject, browser *ICefBrowser, frame *ICefFrame, params *ICefContextMenuParams, model *ICefMenuModel) {
		model.Clear()
	})
	m.chromium.SetOnBeforeBrowser(func(sender lcl.IObject, browser *ICefBrowser, frame *ICefFrame) bool {
		BrowserWindow.setOrIncNextWindowNum(browser.Identifier() + 1)
		return false
	})
	m.chromium.SetOnBeforeResourceLoad(func(sender lcl.IObject, browser *ICefBrowser, frame *ICefFrame, request *ICefRequest, callback *ICefCallback, result *TCefReturnValue) {
		if assetserve.AssetsServerHeaderKeyValue != "" {
			request.SetHeaderByName(assetserve.AssetsServerHeaderKeyName, assetserve.AssetsServerHeaderKeyValue, true)
		}
	})
	m.chromium.SetOnClose(func(sender lcl.IObject, browser *ICefBrowser, aAction *TCefCloseBrowsesAction) {
		logger.Debug("tray.chromium.onClose")
		if IsDarwin() {
			m.windowParent.DestroyChildWindow()
		}
		*aAction = CbaClose
	})
	m.chromium.SetOnBeforeClose(func(sender lcl.IObject, browser *ICefBrowser) {
		logger.Debug("tray.chromium.onBeforeClose")
	})
	m.chromium.SetOnProcessMessageReceived(func(sender lcl.IObject, browser *ICefBrowser, frame *ICefFrame, sourceProcess CefProcessId, message *ipc.ICefProcessMessage) bool {
		return false
	})
	m.windowParent.SetChromium(m.chromium, 0)
	m.chromium.SetDefaultURL(m.url)
}

//设置托盘图标
func (m *tLCLTrayWindow) SetIcon(iconResourcePath string) {
	m.trayIcon.Icon().LoadFromFSFile(iconResourcePath)
}
