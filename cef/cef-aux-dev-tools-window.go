//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under GNU General Public License v3.0
//
//----------------------------------------

package cef

import (
	"fmt"
	. "github.com/energye/energy/consts"
	"github.com/energye/golcl/lcl"
	"github.com/energye/golcl/lcl/api"
	"github.com/energye/golcl/lcl/types"
)

const (
	dev_tools_name = "DevTools"
)

func updateBrowserDevTools(browser *ICefBrowser, title string) {
	if browserWinInfo := BrowserWindow.GetWindowInfo(browser.Identifier()); browserWinInfo != nil {
		if browserWinInfo.getAuxTools() != nil && browserWinInfo.getAuxTools().devToolsWindow != nil {
			QueueAsyncCall(func(id int) {
				browserWinInfo.getAuxTools().devToolsWindow.SetTitle(fmt.Sprintf("%s - %s", dev_tools_name, browser.MainFrame().Url))
			})
		}
	}
}

func (m *ICefBrowser) createBrowserDevTools(browserWinInfo IBrowserWindow) {
	if browserWinInfo.IsLCL() {
		QueueAsyncCall(func(id int) {
			window := browserWinInfo.AsLCLBrowserWindow().BrowserWindow()
			window.createAuxTools()
			winAuxTools := window.auxTools
			if winAuxTools.devToolsWindow != nil {
				winAuxTools.devToolsWindow.Show()
				return
			}
			devToolsWindow := &LCLBrowserWindow{}
			winAuxTools.devToolsWindow = devToolsWindow
			devToolsWindow.SetWindowType(WT_DEV_TOOLS)
			devToolsWindow.TForm = lcl.NewForm(window)
			devToolsWindow.SetTitle(fmt.Sprintf("%s - %s", dev_tools_name, m.MainFrame().Url))
			devToolsWindow.FormCreate()
			devToolsWindow.defaultWindowEvent()
			devToolsWindow.defaultWindowCloseEvent()
			winAuxTools.devToolsWindow.SetSize(1024, 768)
			winAuxTools.devToolsWindow.SetShowInTaskBar()
			devToolsWindow.SetOnResize(func(sender lcl.IObject) {
				winAuxTools.devToolsX = devToolsWindow.Left()
				winAuxTools.devToolsY = devToolsWindow.Top()
				winAuxTools.devToolsWidth = devToolsWindow.Width()
				winAuxTools.devToolsHeight = devToolsWindow.Height()

				if devToolsWindow.isClosing {
					return
				}
				if devToolsWindow.chromium != nil {
					devToolsWindow.chromium.NotifyMoveOrResizeStarted()
				}
				if devToolsWindow.windowParent != nil {
					devToolsWindow.windowParent.UpdateSize()
				}
			})
			devToolsWindow.SetOnClose(func(sender lcl.IObject, action *types.TCloseAction) {
				if devToolsWindow.isClosing {
					return
				}
				*action = types.CaFree
			})
			devToolsWindow.SetOnCloseQuery(func(sender lcl.IObject, canClose *bool) {
				if devToolsWindow.isClosing {
					return
				}
				devToolsWindow.isClosing = true
				BrowserWindow.removeWindowInfo(devToolsWindow.windowId)
			})

			devToolsWindow.ChromiumCreate(nil, "")
			devToolsWindow.putChromiumWindowInfo()
			devToolsWindow.defaultChromiumEvent()
			winAuxTools.devToolsWindow.Show()
			//明确的生成下一个窗体序号
			BrowserWindow.setOrIncNextWindowNum()
			_CEFBrowser_ShowDevTools(devToolsWindow.chromium.Instance(), uintptr(m.Identifier()), devToolsWindow.windowParent.Instance(), api.PascalStr(dev_tools_name))
		})
	} else if browserWinInfo.IsViewsFramework() {
		if application.cfg.remoteDebuggingPort > 1024 && application.cfg.remoteDebuggingPort < 65535 {
			window := browserWinInfo.AsViewsFrameworkBrowserWindow().BrowserWindow()
			window.createAuxTools()
			winAuxTools := window.auxTools
			if winAuxTools.devToolsWindow != nil {
				winAuxTools.devToolsWindow.Show()
				return
			}
			wp := NewWindowProperty()
			wp.Url = ""
			//devToolsWindow := NewViewsFrameworkBrowserWindow(nil, nil)
		}
	}
}
