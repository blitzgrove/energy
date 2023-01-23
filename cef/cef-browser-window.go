//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under GNU General Public License v3.0
//
//----------------------------------------

package cef

import (
	"github.com/energye/energy/consts"
	"github.com/energye/golcl/lcl"
	"github.com/energye/golcl/lcl/types"
)

//辅助工具
type auxTools struct {
	devToolsWindow   IBrowserWindow //devTools
	devToolsX        int32          //上次改变的窗体位置，宽度
	devToolsY        int32          //
	devToolsWidth    int32          //
	devToolsHeight   int32          //
	viewSourceWindow IBrowserWindow //viewSource
	viewSourceUrl    string         //
	viewSourceX      int32          //上次改变的窗体位置，宽度
	viewSourceY      int32          //
	viewSourceWidth  int32          //
	viewSourceHeight int32          //
}

//窗口属性配置器
//
//部分属性配置并不支持所有平台
type WindowProperty struct {
	IsShowModel               bool               //是否以模态窗口显示
	WindowState               types.TWindowState //窗口 状态
	Title                     string             //窗口 标题
	Url                       string             //默认打开URL
	Icon                      string             //窗口图标 加载本地图标 local > /app/resources/icon.ico
	IconFS                    string             //窗口图标 加载emfs内置图标 emfs > resources/icon.ico
	_EnableHideCaption        bool               //窗口 是否隐藏标题栏
	EnableMinimize            bool               //窗口 是否启用最小化 default: true
	EnableMaximize            bool               //窗口 是否启用最大化 default: true
	EnableResize              bool               //窗口 是否允许调整大小 default: true
	EnableClose               bool               //窗口 关闭时是否关闭窗口 default: true
	EnableCenterWindow        bool               //窗口 居中显示 default: true
	EnableDragFile            bool               //窗口 是否允许向窗口内拖拽文件
	EnableWebkitAppRegion     bool               //窗口 html元素中设置css属性 -webkit-app-region: drag/no-drag 是否允许拖拽窗口 default: true
	EnableCaptionDClkMaximize bool               //窗口 是否启用标题栏双击放大还原 default: true
	AlwaysOnTop               bool               //窗口 窗口置顶
	X                         int32              //窗口 EnableCenterWindow=false X坐标 default: 100
	Y                         int32              //窗口 EnableCenterWindow=false Y坐标 default: 100
	Width                     int32              //窗口 宽 default: 1024
	Height                    int32              //窗口 高 default: 768
}

//浏览器窗口基础接口
//
//定义了常用函数
type IBrowserWindow interface {
	Id() int32                                                   //窗口ID
	Handle() types.HWND                                          //
	Show()                                                       //显示窗口
	Hide()                                                       //隐藏窗口
	Maximize()                                                   //窗口最大化
	Minimize()                                                   //窗口最小化
	Close()                                                      //关闭窗口
	CloseBrowserWindow()                                         //关闭浏览器窗口
	WindowType() consts.WINDOW_TYPE                              //窗口类型
	SetWindowType(windowType consts.WINDOW_TYPE)                 //设置窗口类型
	Browser() *ICefBrowser                                       //窗口内的Browser对象
	Chromium() IChromium                                         //窗口内的Chromium对象
	DisableMaximize()                                            //禁用最大化
	DisableMinimize()                                            //禁用最小化
	DisableResize()                                              //禁用窗口大小调整
	EnableMaximize()                                             //启用最大化
	EnableMinimize()                                             //启用最小化
	EnableResize()                                               //启用窗口大小调用
	IsClosing() bool                                             //窗口是否已状态
	AsViewsFrameworkBrowserWindow() IViewsFrameworkBrowserWindow //转换为ViewsFramework窗口接口
	AsLCLBrowserWindow() ILCLBrowserWindow                       //转换为LCL窗口接口
	Frames() TCEFFrame                                           //窗口内的所有Frame
	addFrame(frame *ICefFrame)                                   //
	setBrowser(browser *ICefBrowser)                             //
	createAuxTools()                                             //
	getAuxTools() *auxTools                                      //
	EnableAllDefaultEvent()                                      //启用所有默认事件
	SetTitle(title string)                                       //设置窗口标题栏标题
	IsViewsFramework() bool                                      //是否为 IViewsFrameworkBrowserWindow 窗口
	IsLCL() bool                                                 //是否为 ILCLBrowserWindow 窗口
	WindowProperty() *WindowProperty                             //窗口常用属性
	SetWidth(value int32)                                        //设置窗口宽
	SetHeight(value int32)                                       //设置窗口高
	Point() *TCefPoint                                           //窗口坐标
	Size() *TCefSize                                             //窗口宽高
	Bounds() *TCefRect                                           //窗口坐标和宽高
	SetPoint(x, y int32)                                         //设置窗口坐标
	SetSize(width, height int32)                                 //设置窗口宽高
	SetBounds(x, y, width, height int32)                         //设置窗口坐标和宽高
	SetCenterWindow(value bool)                                  //设置窗口居中
	NewCefTray(width, height int32, url string) ITray            //创建托盘CEF自定义html, 实现4种系统托盘，1: LCL原生, 2: CEF基于LCL, 3: VF(views framework), 4:系统原生
	ShowTitle()                                                  //显示窗口标题栏
	HideTitle()                                                  //隐藏窗口标题栏
	SetDefaultInTaskBar()                                        //窗口默认在任务栏上显示图标
	SetShowInTaskBar()                                           //强制窗口在任务栏上显示图标
	SetNotInTaskBar()                                            //强制不在任务栏上显示窗口图标
}

//浏览器 LCLBrowserWindow 窗口接口 继承 IBrowserWindow
//
//定义了LCL常用函数
type ILCLBrowserWindow interface {
	IBrowserWindow
	BrowserWindow() *LCLBrowserWindow //返回 LCLBrowserWindow 窗口结构
	EnableDefaultCloseEvent()         //启用默认关闭事件
	WindowParent() ITCefWindowParent  //浏览器父窗口组件
	DisableTransparent()              //禁用窗口透明
	EnableTransparent(value uint8)    //用用并设置窗口透明
	DisableSystemMenu()               //禁用标题栏系统菜单
	DisableHelp()                     //禁用标题栏帮助
	EnableSystemMenu()                //启用标题栏系统菜单
	EnableHelp()                      //启用标题栏帮助
	NewTray() ITray                   //创建LCL的系统托盘
}

//浏览器 ViewsFrameworkBrowserWindow 窗口接口 继承 IBrowserWindow
//
//定义了ViewsFramework常用函数
type IViewsFrameworkBrowserWindow interface {
	IBrowserWindow
	BrowserWindow() *ViewsFrameworkBrowserWindow                       //返回 ViewsFrameworkBrowserWindow 窗口结构
	CreateTopLevelWindow()                                             //创建窗口, 在窗口组件中需要默认调用Show函数
	CenterWindow(size *TCefSize)                                       //设置窗口居中，同时指定窗口大小
	Component() lcl.IComponent                                         //窗口父组件
	WindowComponent() *TCEFWindowComponent                             //窗口组件
	BrowserViewComponent() *TCEFBrowserViewComponent                   //窗口浏览器组件
	SetOnWindowCreated(onWindowCreated WindowComponentOnWindowCreated) //设置窗口默认的创建回调事件函数
}

//创建一个属性配置器，带有窗口默认属性值
func NewWindowProperty() WindowProperty {
	return WindowProperty{
		Title:                     "Energy",
		Url:                       "about:blank",
		EnableMinimize:            true,
		EnableMaximize:            true,
		EnableResize:              true,
		EnableClose:               true,
		EnableCenterWindow:        true,
		EnableCaptionDClkMaximize: true,
		EnableWebkitAppRegion:     true,
		X:                         100,
		Y:                         100,
		Width:                     1024,
		Height:                    768,
	}
}
