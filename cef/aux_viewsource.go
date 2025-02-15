//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

// 辅助工具-显示网页源代码

package cef

import (
	"fmt"
	"github.com/energye/energy/v2/common"
	. "github.com/energye/energy/v2/consts"
	"github.com/energye/energy/v2/pkgs/assetserve"
	"github.com/energye/golcl/lcl"
)

const (
	viewSourceName = "view-source"
)

func (m *ICefBrowser) createBrowserViewSource() {
	if currentWindowInfo := BrowserWindow.GetWindowInfo(m.Identifier()); currentWindowInfo != nil {
		var frame = m.MainFrame()
		if currentWindowInfo.IsLCL() {
			var viewSourceUrl = fmt.Sprintf("view-source:%s", frame.Url())
			wp := NewWindowProperty()
			wp.Url = viewSourceUrl
			wp.Title = fmt.Sprintf("%s - %s", viewSourceName, frame.Url())
			wp.WindowType = WT_VIEW_SOURCE
			viewSourceWindow := NewLCLBrowserWindow(nil, wp)
			viewSourceWindow.SetWidth(800)
			viewSourceWindow.SetHeight(600)
			if common.IsDarwin() || common.IsLinux() {
				viewSourceWindow.Chromium().SetOnAfterCreated(func(sender lcl.IObject, browser *ICefBrowser) {
					viewSourceWindow.Chromium().LoadUrl(viewSourceUrl)
				})
			}
			if assetserve.AssetsServerHeaderKeyValue != "" {
				viewSourceWindow.Chromium().SetOnBeforeResourceLoad(func(sender lcl.IObject, browser *ICefBrowser, frame *ICefFrame, request *ICefRequest, callback *ICefCallback, result *TCefReturnValue) {
					request.SetHeaderByName(assetserve.AssetsServerHeaderKeyName, assetserve.AssetsServerHeaderKeyValue, true)
				})
			}
			viewSourceWindow.Chromium().SetOnBeforePopup(func(sender lcl.IObject, browser *ICefBrowser, frame *ICefFrame, beforePopupInfo *BeforePopupInfo, client *ICefClient, noJavascriptAccess *bool) bool {
				wp := NewWindowProperty()
				wp.Url = beforePopupInfo.TargetUrl
				wp.Title = beforePopupInfo.TargetUrl
				wp.WindowType = WT_VIEW_SOURCE
				bw := NewLCLBrowserWindow(nil, wp)
				bw.SetWidth(800)
				bw.SetHeight(600)
				bw.EnableDefaultCloseEvent()
				QueueAsyncCall(func(id int) { //main thread run
					bw.Show()
				})
				return true
			})
			viewSourceWindow.EnableDefaultCloseEvent()
			QueueAsyncCall(func(id int) { //main thread run
				viewSourceWindow.Show()
			})
		} else if currentWindowInfo.IsViewsFramework() {
			frame.ViewSource()
		}
	}
}
