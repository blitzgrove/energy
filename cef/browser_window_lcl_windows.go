//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

//go:build windows

// LCL窗口组件定义和实现-windows平台

package cef

import (
	"github.com/energye/energy/v2/cef/winapi"
	"github.com/energye/energy/v2/consts"
	"github.com/energye/energy/v2/consts/messages"
	et "github.com/energye/energy/v2/types"
	"github.com/energye/golcl/lcl"
	"github.com/energye/golcl/lcl/rtl"
	"github.com/energye/golcl/lcl/types"
	"github.com/energye/golcl/lcl/types/colors"
	"github.com/energye/golcl/lcl/win"
)

// 定义四角和边框范围
var (
	angleRange  int32 = 10 //四角
	borderRange int32 = 5  //四边框
)

// customWindowCaption 自定义窗口标题栏
//
// 隐藏窗口标题栏，通过html+css实现自定义窗口标题栏，实现窗口拖拽等
type customWindowCaption struct {
	bw                   *LCLBrowserWindow     //
	canCaption           bool                  //当前鼠标是否在标题栏区域
	canBorder            bool                  //当前鼠标是否在边框
	borderHT, borderWMSZ int                   //borderHT: 鼠标所在边框位置, borderWMSZ: 窗口改变大小边框方向 borderMD:
	borderMD             bool                  //borderMD: 鼠标调整窗口大小，已按下后，再次接收到132消息应该忽略该消息
	regions              *TCefDraggableRegions //窗口内html拖拽区域
	rgn                  *et.HRGN              //
}

// ShowTitle 显示标题栏
func (m *LCLBrowserWindow) ShowTitle() {
	m.WindowProperty().EnableHideCaption = false
	//win.SetWindowLong(m.Handle(), win.GWL_STYLE, uintptr(win.GetWindowLong(m.Handle(), win.GWL_STYLE)|win.WS_CAPTION))
	//win.SetWindowPos(m.Handle(), m.Handle(), 0, 0, 0, 0, win.SWP_NOSIZE|win.SWP_NOMOVE|win.SWP_NOZORDER|win.SWP_NOACTIVATE|win.SWP_FRAMECHANGED)
	m.EnabledMaximize(m.WindowProperty().EnableMaximize)
	m.EnabledMinimize(m.WindowProperty().EnableMinimize)
	m.SetBorderStyle(types.BsSizeable)
}

// HideTitle 隐藏标题栏 无边框样式
func (m *LCLBrowserWindow) HideTitle() {
	m.WindowProperty().EnableHideCaption = true
	m.SetBorderStyle(types.BsNone)
	//m.Frameless()
}

// SetRoundRectRgn 窗口无边框时圆角设置
func (m *LCLBrowserWindow) SetRoundRectRgn(rgn int) {
	if m.rgn == 0 && rgn > 0 {
		m.rgn = rgn
		m.SetOnPaint(func(sender lcl.IObject) {
			hnd := winapi.CreateRoundRectRgn(0, 0, et.LongInt(m.Width()), et.LongInt(m.Height()), et.LongInt(m.rgn), et.LongInt(m.rgn))
			winapi.SetWindowRgn(et.HWND(m.Handle()), hnd, true)
		})
	}
}

// FramelessForDefault 窗口四边框系统默认样式
//
//	TODO 窗口顶部有条线,
func (m *LCLBrowserWindow) FramelessForDefault() {
	gwlStyle := win.GetWindowLong(m.Handle(), win.GWL_STYLE)
	win.SetWindowLong(m.Handle(), win.GWL_STYLE, uintptr(gwlStyle&^win.WS_CAPTION&^win.WS_BORDER|win.WS_THICKFRAME))
	win.SetWindowPos(m.Handle(), 0, 0, 0, 0, 0, uint32(win.SWP_NOMOVE|win.SWP_NOSIZE|win.SWP_FRAMECHANGED))

	//winapi.SetClassLongPtr(this.Handle, GCL_STYLE, GetClassLong(this.Handle, GCL_STYLE)|CS_DropSHADOW)
	//gclStyle := winapi.GetClassLongPtr(et.HWND(m.Handle()), messages.GCL_STYLE)
	//winapi.WinSetClassLongPtr().SetClassLongPtr(et.HWND(m.Handle()), messages.GCL_STYLE, gclStyle|messages.CS_DROPSHADOW)
}

// FramelessForLine 窗口四边框是一条细线
func (m *LCLBrowserWindow) FramelessForLine() {
	gwlStyle := win.GetWindowLong(m.Handle(), win.GWL_STYLE)
	win.SetWindowLong(m.Handle(), win.GWL_STYLE, uintptr(gwlStyle&^win.WS_CAPTION&^win.WS_THICKFRAME|win.WS_BORDER))
	win.SetWindowPos(m.Handle(), 0, 0, 0, 0, 0, uint32(win.SWP_NOMOVE|win.SWP_NOSIZE|win.SWP_FRAMECHANGED))
}

// Frameless 无边框
func (m *LCLBrowserWindow) Frameless() {
	gwlStyle := win.GetWindowLong(m.Handle(), win.GWL_STYLE)
	win.SetWindowLong(m.Handle(), win.GWL_STYLE, uintptr(gwlStyle&^win.WS_CAPTION&^win.WS_THICKFRAME))
	win.SetWindowPos(m.Handle(), 0, 0, 0, 0, 0, uint32(win.SWP_NOMOVE|win.SWP_NOSIZE|win.SWP_FRAMECHANGED))
}

// windows无边框窗口任务栏处理
func (m *LCLBrowserWindow) taskMenu() {
	m.SetOnShow(func(sender lcl.IObject) bool {
		if m.WindowProperty().EnableHideCaption {
			gwlStyle := win.GetWindowLong(m.Handle(), win.GWL_STYLE)
			win.SetWindowLong(m.Handle(), win.GWL_STYLE, uintptr(gwlStyle|win.WS_SYSMENU|win.WS_MINIMIZEBOX))
		}
		return false
	})
}

// SetFocus
//
//	在窗口 (Visible = true) 显示之后设置窗口焦点
//	https://learn.microsoft.com/zh-cn/windows/win32/api/winuser/nf-winuser-showwindow
//	https://learn.microsoft.com/zh-cn/windows/win32/api/winuser/nf-winuser-setfocus
func (m *LCLBrowserWindow) SetFocus() {
	if m.TForm != nil {
		m.Visible()
		//窗口\激活在Z序中的下个顶层窗口
		m.Minimize()
		//激活窗口出现在前景
		m.Restore()
		//窗口设置焦点
		m.TForm.SetFocus()
	}
}

//func (m *LCLBrowserWindow) Frameless() {
//	var rect = &types.TRect{}
//	win.GetWindowRect(m.Handle(), rect)
//	win.SetWindowPos(m.Handle(), 0, rect.Left, rect.Top, rect.Right-rect.Left, rect.Bottom-rect.Top, win.SWP_FRAMECHANGED)
//}

// freeRgn
func (m *customWindowCaption) freeRgn() {
	if m.rgn != nil {
		winapi.SetRectRgn(m.rgn, 0, 0, 0, 0)
		winapi.DeleteObject(m.rgn)
		m.rgn.Free()
	}
}

// freeRegions
func (m *customWindowCaption) freeRegions() {
	if m.regions != nil {
		m.regions.regions = nil
		m.regions = nil
	}
}

// free
func (m *customWindowCaption) free() {
	if m != nil {
		m.freeRgn()
		m.freeRegions()
	}
}

// onNCMouseMove NC 非客户区鼠标移动
func (m *customWindowCaption) onNCMouseMove(message *types.TMessage, lResult *types.LRESULT, aHandled *bool) {
	if m.canCaption { // 当前在标题栏
	} else if m.canBorder { // 当前在边框
		*lResult = types.LRESULT(m.borderHT)
		*aHandled = true
	}
}

// onSetCursor 设置鼠标图标
func (m *customWindowCaption) onSetCursor(message *types.TMessage, lResult *types.LRESULT, aHandled *bool) {
	if m.canBorder { //当前在边框
		switch winapi.LOWORD(uint32(message.LParam)) {
		case messages.HTBOTTOMRIGHT, messages.HTTOPLEFT: //右下 左上
			*lResult = types.LRESULT(m.borderHT)
			*aHandled = true
			winapi.SetCursor(winapi.LoadCursor(0, messages.IDC_SIZENWSE))
		case messages.HTRIGHT, messages.HTLEFT: //右 左
			*lResult = types.LRESULT(m.borderHT)
			*aHandled = true
			winapi.SetCursor(winapi.LoadCursor(0, messages.IDC_SIZEWE))
		case messages.HTTOPRIGHT, messages.HTBOTTOMLEFT: //右上 左下
			*lResult = types.LRESULT(m.borderHT)
			*aHandled = true
			winapi.SetCursor(winapi.LoadCursor(0, messages.IDC_SIZENESW))
		case messages.HTTOP, messages.HTBOTTOM: //上 下
			*lResult = types.LRESULT(m.borderHT)
			*aHandled = true
			winapi.SetCursor(winapi.LoadCursor(0, messages.IDC_SIZENS))
		}
	}
}

// onCanBorder 鼠标是否在边框
func (m *customWindowCaption) onCanBorder(x, y int32, rect *types.TRect) (int, bool) {
	if m.canBorder = x <= rect.Width() && x >= rect.Width()-angleRange && y <= angleRange; m.canBorder { // 右上
		m.borderWMSZ = messages.WMSZ_TOPRIGHT
		m.borderHT = messages.HTTOPRIGHT
		return m.borderHT, true
	} else if m.canBorder = x <= rect.Width() && x >= rect.Width()-angleRange && y <= rect.Height() && y >= rect.Height()-angleRange; m.canBorder { // 右下
		m.borderWMSZ = messages.WMSZ_BOTTOMRIGHT
		m.borderHT = messages.HTBOTTOMRIGHT
		return m.borderHT, true
	} else if m.canBorder = x <= angleRange && y <= angleRange; m.canBorder { //左上
		m.borderWMSZ = messages.WMSZ_TOPLEFT
		m.borderHT = messages.HTTOPLEFT
		return m.borderHT, true
	} else if m.canBorder = x <= angleRange && y >= rect.Height()-angleRange; m.canBorder { //左下
		m.borderWMSZ = messages.WMSZ_BOTTOMLEFT
		m.borderHT = messages.HTBOTTOMLEFT
		return m.borderHT, true
	} else if m.canBorder = x > angleRange && x < rect.Width()-angleRange && y <= borderRange; m.canBorder { //上
		m.borderWMSZ = messages.WMSZ_TOP
		m.borderHT = messages.HTTOP
		return m.borderHT, true
	} else if m.canBorder = x > angleRange && x < rect.Width()-angleRange && y >= rect.Height()-borderRange; m.canBorder { //下
		m.borderWMSZ = messages.WMSZ_BOTTOM
		m.borderHT = messages.HTBOTTOM
		return m.borderHT, true
	} else if m.canBorder = x <= borderRange && y > angleRange && y < rect.Height()-angleRange; m.canBorder { //左
		m.borderWMSZ = messages.WMSZ_LEFT
		m.borderHT = messages.HTLEFT
		return m.borderHT, true
	} else if m.canBorder = x <= rect.Width() && x >= rect.Width()-borderRange && y > angleRange && y < rect.Height()-angleRange; m.canBorder { // 右
		m.borderWMSZ = messages.WMSZ_RIGHT
		m.borderHT = messages.HTRIGHT
		return m.borderHT, true
	}
	return 0, false
}

// onNCLButtonDown NC 鼠标左键按下
func (m *customWindowCaption) onNCLButtonDown(hWND types.HWND, message *types.TMessage, lResult *types.LRESULT, aHandled *bool) {
	if m.canCaption { // 标题栏
		*lResult = messages.HTCAPTION
		*aHandled = true
		//全屏时不能移动窗口
		if m.bw.WindowProperty().current.ws == types.WsFullScreen {
			return
		}
		x, y := m.toPoint(message)
		m.borderMD = true
		if win.ReleaseCapture() {
			win.PostMessage(hWND, messages.WM_NCLBUTTONDOWN, messages.HTCAPTION, rtl.MakeLParam(uint16(x), uint16(y)))
		}
	} else if m.canBorder { // 边框
		x, y := m.toPoint(message)
		*lResult = types.LRESULT(m.borderHT)
		*aHandled = true
		m.borderMD = true
		if win.ReleaseCapture() {
			win.PostMessage(hWND, messages.WM_SYSCOMMAND, uintptr(messages.SC_SIZE|m.borderWMSZ), rtl.MakeLParam(uint16(x), uint16(y)))
		}
	}
}

// toPoint 转换XY坐标
func (m *customWindowCaption) toPoint(message *types.TMessage) (x, y int32) {
	return winapi.GET_X_LPARAM(message.LParam), winapi.GET_Y_LPARAM(message.LParam)
}

// isCaption
// 鼠标是否在标题栏区域
//
// 如果启用了css拖拽则校验拖拽区域,否则只返回相对于浏览器窗口的x,y坐标
func (m *customWindowCaption) isCaption(hWND et.HWND, message *types.TMessage) (int32, int32, bool) {
	dx, dy := m.toPoint(message)
	p := &et.Point{
		X: dx,
		Y: dy,
	}
	winapi.ScreenToClient(hWND, p)
	p.X -= m.bw.WindowParent().Left()
	p.Y -= m.bw.WindowParent().Top()
	if m.bw.WindowProperty().EnableWebkitAppRegion && m.rgn != nil {
		m.canCaption = winapi.PtInRegion(m.rgn, p.X, p.Y)
	} else {
		m.canCaption = false
	}
	return p.X, p.Y, m.canCaption
}

// doOnRenderCompMsg
func (m *LCLBrowserWindow) doOnRenderCompMsg(message *types.TMessage, lResult *types.LRESULT, aHandled *bool) {
	switch message.Msg {
	case messages.WM_NCLBUTTONDBLCLK: // 163 NC left dclick
		//标题栏拖拽区域 双击最大化和还原
		if m.cwcap.canCaption && m.WindowProperty().EnableWebkitAppRegionDClk {
			*lResult = messages.HTCAPTION
			*aHandled = true
			if win.ReleaseCapture() {
				if m.WindowState() == types.WsNormal {
					win.PostMessage(m.Handle(), messages.WM_SYSCOMMAND, messages.SC_MAXIMIZE, 0)
				} else {
					win.PostMessage(m.Handle(), messages.WM_SYSCOMMAND, messages.SC_RESTORE, 0)
				}
				win.SendMessage(m.Handle(), messages.WM_NCLBUTTONUP, messages.HTCAPTION, 0)
			}
		}
	case messages.WM_NCLBUTTONDOWN: // 161 nc left down
		m.cwcap.onNCLButtonDown(m.Handle(), message, lResult, aHandled)
	case messages.WM_NCLBUTTONUP: // 162 nc l up
		if m.cwcap.canCaption {
			*lResult = messages.HTCAPTION
			*aHandled = true
		}
	case messages.WM_NCMOUSEMOVE: // 160 nc mouse move
		m.cwcap.onNCMouseMove(message, lResult, aHandled)
	case messages.WM_SETCURSOR: // 32 设置鼠标图标样式
		m.cwcap.onSetCursor(message, lResult, aHandled)
	case messages.WM_NCHITTEST: // 132 NCHITTEST
		if m.cwcap.borderMD { //TODO 测试windows7, 161消息之后再次处理132消息导致消息错误
			m.cwcap.borderMD = false
			return
		}
		//鼠标坐标是否在标题区域
		x, y, caption := m.cwcap.isCaption(et.HWND(m.Handle()), message)
		if caption { //窗口标题栏
			*lResult = messages.HTCAPTION
			*aHandled = true
		} else if m.WindowProperty().EnableHideCaption && m.WindowProperty().EnableResize && m.WindowState() == types.WsNormal { //1.窗口隐藏标题栏 2.启用了调整窗口大小 3.非最大化、最小化、全屏状态
			//全屏时不能调整窗口大小
			if m.WindowProperty().current.ws == types.WsFullScreen {
				return
			}
			rect := m.BoundsRect()
			if result, handled := m.cwcap.onCanBorder(x, y, &rect); handled {
				*lResult = types.LRESULT(result)
				*aHandled = true
			}
		}
	}
}

// setDraggableRegions
// 每一次拖拽区域改变都需要重新设置
func (m *LCLBrowserWindow) setDraggableRegions() {
	var scp float32
	// Windows 10 版本 1607 [仅限桌面应用]
	// Windows Server 2016 [仅限桌面应用]
	// 可动态调整
	dpi, err := winapi.GetDpiForWindow(et.HWND(m.Handle()))
	if err == nil {
		scp = float32(dpi) / 96.0
	} else {
		// 使用默认的，但不能动态调整
		scp = winapi.ScalePercent()
	}
	//在主线程中运行
	m.RunOnMainThread(func() {
		if m.cwcap.rgn == nil {
			//第一次时创建RGN
			m.cwcap.rgn = winapi.CreateRectRgn(0, 0, 0, 0)
		} else {
			//每次重置RGN
			winapi.SetRectRgn(m.cwcap.rgn, 0, 0, 0, 0)
		}
		// 重新根据缩放比计算新的区域位置
		for i := 0; i < m.cwcap.regions.RegionsCount(); i++ {
			region := m.cwcap.regions.Region(i)
			x := int32(float32(region.Bounds.X) * scp)
			y := int32(float32(region.Bounds.Y) * scp)
			w := int32(float32(region.Bounds.Width) * scp)
			h := int32(float32(region.Bounds.Height) * scp)
			creRGN := winapi.CreateRectRgn(x,
				y,
				x+w,
				y+h)
			if region.Draggable {
				winapi.CombineRgn(m.cwcap.rgn, m.cwcap.rgn, creRGN, consts.RGN_OR)
			} else {
				winapi.CombineRgn(m.cwcap.rgn, m.cwcap.rgn, creRGN, consts.RGN_DIFF)
			}
			winapi.DeleteObject(creRGN)
		}
	})
}

// registerWindowsCompMsgEvent
// 注册windows下CompMsg事件
func (m *LCLBrowserWindow) registerWindowsCompMsgEvent() {
	var bwEvent = BrowserWindow.browserEvent
	m.Chromium().SetOnRenderCompMsg(func(sender lcl.IObject, message *types.TMessage, lResult *types.LRESULT, aHandled *bool) {
		if bwEvent.onRenderCompMsg != nil {
			bwEvent.onRenderCompMsg(sender, message, lResult, aHandled)
		}
		if !*aHandled {
			m.doOnRenderCompMsg(message, lResult, aHandled)
		}
	})
	if m.WindowProperty().EnableWebkitAppRegion && m.WindowProperty().EnableWebkitAppRegionDClk {
		m.windowResize = func(sender lcl.IObject) bool {
			if m.WindowState() == types.WsMaximized && (m.WindowProperty().EnableHideCaption || m.BorderStyle() == types.BsNone || m.BorderStyle() == types.BsSingle) {
				var monitor = m.Monitor().WorkareaRect()
				m.SetBounds(monitor.Left, monitor.Top, monitor.Right-monitor.Left, monitor.Bottom-monitor.Top)
				m.SetWindowState(types.WsMaximized)
			}
			return false
		}
	}

	//if m.WindowProperty().EnableWebkitAppRegion {
	//
	//} else {
	//	if bwEvent.onRenderCompMsg != nil {
	//		m.chromium.SetOnRenderCompMsg(bwEvent.onRenderCompMsg)
	//	}
	//}
}

// Restore Windows平台，窗口还原
func (m *LCLBrowserWindow) Restore() {
	if m.TForm == nil {
		return
	}
	m.RunOnMainThread(func() {
		if win.ReleaseCapture() {
			win.SendMessage(m.Handle(), messages.WM_SYSCOMMAND, messages.SC_RESTORE, 0)
		}
	})
}

// Minimize Windows平台，窗口最小化
func (m *LCLBrowserWindow) Minimize() {
	if m.TForm == nil {
		return
	}
	m.RunOnMainThread(func() {
		if win.ReleaseCapture() {
			win.PostMessage(m.Handle(), messages.WM_SYSCOMMAND, messages.SC_MINIMIZE, 0)
		}
	})
}

// Maximize Windows平台，窗口最大化/还原
func (m *LCLBrowserWindow) Maximize() {
	if m.TForm == nil {
		return
	}
	m.RunOnMainThread(func() {
		if win.ReleaseCapture() {
			if m.WindowState() == types.WsNormal {
				win.PostMessage(m.Handle(), messages.WM_SYSCOMMAND, messages.SC_MAXIMIZE, 0)
			} else {
				win.SendMessage(m.Handle(), messages.WM_SYSCOMMAND, messages.SC_RESTORE, 0)
			}
		}
	})
}

// 窗口透明
//func (m *LCLBrowserWindow) SetTransparentColor() {
//	m.SetColor(colors.ClNavy)
//	Exstyle := win.GetWindowLong(m.Handle(), win.GWL_EXSTYLE)
//	Exstyle = Exstyle | win.WS_EX_LAYERED&^win.WS_EX_TRANSPARENT
//	win.SetWindowLong(m.Handle(), win.GWL_EXSTYLE, uintptr(Exstyle))
//	win.SetLayeredWindowAttributes(m.Handle(),
//		colors.ClNavy, //crKey 指定需要透明的背景颜色值
//		255,           //bAlpha 设置透明度,0完全透明，255不透明
//		//LWA_ALPHA: crKey无效,bAlpha有效
//		//LWA_COLORKEY: 窗体中的所有颜色为crKey的地方全透明，bAlpha无效
//		//LWA_ALPHA | LWA_COLORKEY: crKey的地方全透明，其它地方根据bAlpha确定透明度
//		win.LWA_ALPHA|win.LWA_COLORKEY)
//}

func (m *LCLBrowserWindow) doDrag() {
	// Windows Drag Window
	// m.drag != nil 时，这里处理的是 up 事件, 给标题栏标记为false
	if m.drag != nil {
		m.drag.drag()
	} else {
		// 全屏时不能拖拽窗口
		if m.WindowProperty().current.ws == types.WsFullScreen {
			return
		}
		// 此时是 down 事件, 拖拽窗口
		if win.ReleaseCapture() {
			win.PostMessage(m.Handle(), messages.WM_NCLBUTTONDOWN, messages.HTCAPTION, 0)
			m.cwcap.canCaption = true
		}
	}
}

// framelessWindowBroderTransparent 无边框窗口自定义窗口透明
func (m *WindowBroderForm) framelessWindowBroderTransparent(hWnd et.HWND) {
	exStyle := winapi.GetWindowLong(hWnd, win.GWL_EXSTYLE)
	exStyle = exStyle | win.WS_EX_LAYERED
	winapi.SetWindowLong(hWnd, win.GWL_EXSTYLE, exStyle)
	win.SetLayeredWindowAttributes(hWnd.ToPtr(), colors.ClWhite, m.alpha, win.LWA_ALPHA|win.LWA_COLORKEY)
}

// SetBorderColors 设置边框颜色
func (m *WindowBroderForm) SetBorderColors(borderColors []types.TColor) {
	m.borderColors = borderColors
}

// SetBorderRgn 设置边框圆角弧度
func (m *WindowBroderForm) SetBorderRgn(rgn int) {
	m.rgn = rgn
}

// SetBorderAlpha 边框透明度 0~255
func (m *WindowBroderForm) SetBorderAlpha(alpha uint8) {
	m.alpha = alpha
}

// Paint 根据设置的border颜色和四边圆角画边框, 在Show之前调用该函数
func (m *WindowBroderForm) Paint() {
	m.SetOnPaint(func(sender lcl.IObject) {
		r := m.ClientRect()
		canvas := m.Canvas()
		//canvas.Brush().SetStyle(types.BsClear)
		pen := canvas.Pen()
		pen.SetWidth(1)
		//pen.SetMode(types.PmNotMask)
		//pen.SetStyle(types.PsClear)
		for i := 1; i <= len(m.borderColors); i++ {
			pen.SetColor(m.borderColors[i-1])
			canvas.Rectangle(r.Left+int32(i), r.Top+int32(i), r.Right-int32(i), r.Bottom-int32(i))
		}
	})
	m.framelessWindowBroderTransparent(et.HWND(m.Handle()))
	m.RoundRectRgn()
}

func (m *WindowBroderForm) RoundRectRgn() {
	hnd := winapi.CreateRoundRectRgn(0, 0, et.LongInt(m.Width()), et.LongInt(m.Height()),
		et.LongInt(m.rgn), et.LongInt(m.rgn))
	winapi.SetWindowRgn(et.HWND(m.Handle()), hnd, true)
}

func NewWindowBorder(window *LCLBrowserWindow) (m *WindowBroderForm) {
	//var borderColors = []types.TColor{0xFEFEFE, 0xFDFDFD, 0xFCFCFC, 0xFBFBFB, 0xFAFAFA, 0xF7F7F7, 0xEDEDED, 0xE3E3E3, 0xD9D9D9}
	m = &WindowBroderForm{
		TForm:        lcl.NewForm(window),
		borderColors: []types.TColor{0xEDEDED, 0xE3E3E3, 0xD9D9D9, colors.ClBlack, colors.ClBlack},
		rgn:          30,
		alpha:        25,
	}

	m.borderWidth = int32(len(m.borderColors)) // 边框宽度, 像素个数
	m.SetPosition(window.Position())
	m.SetLeft(window.Left() - m.borderWidth)
	m.SetTop(window.Top() - m.borderWidth)
	m.SetWidth(window.Width() + m.borderWidth*2)   // 同步主窗口, 比主窗口大一圈, 填充边框颜色的像素个数
	m.SetHeight(window.Height() + m.borderWidth*2) // 同步主窗口, 比主窗口大一圈, 填充边框颜色的像素个数
	m.SetBorderStyle(types.BsNone)                 // 边框窗口无边框
	m.SetColor(colors.ClWhite)                     // 边框窗口的填充颜色, 做为透明色
	m.SetFormStyle(types.FsStayOnTop)
	return m
}
