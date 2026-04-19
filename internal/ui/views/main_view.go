package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/duxweb/codux/internal/app"
)

// MainView 主视图
type MainView struct {
	widget.BaseWidget
	window     fyne.Window
	store      *app.Store
	sidebar    *SidebarView
	workspace  *WorkspaceView
	rightPanel *RightPanelView
	statusBar  *widget.Label
}

// NewMainView 创建主视图
func NewMainView(window fyne.Window, store *app.Store) *MainView {
	mv := &MainView{
		window: window,
		store:  store,
	}
	mv.ExtendBaseWidget(mv)
	return mv
}

// CreateRenderer 实现 fyne.Widget 接口
func (mv *MainView) CreateRenderer() fyne.WidgetRenderer {
	// 创建侧边栏
	mv.sidebar = NewSidebarView(mv.window, mv.store)

	// 创建工作区
	mv.workspace = NewWorkspaceView(mv.window, mv.store)

	// 创建右侧面板
	mv.rightPanel = NewRightPanelView(mv.window, mv.store)

	// 创建状态栏
	mv.statusBar = widget.NewLabel("就绪")
	mv.statusBar.TextStyle = fyne.TextStyle{Italic: true}

	// 主内容区域 (侧边栏 + 工作区 + 右面板)
	mainContent := container.NewHSplit(
		mv.sidebar,
		container.NewHSplit(mv.workspace, mv.rightPanel),
	)

	// 设置侧边栏比例
	mainContent.SetOffset(0.15)

	// 垂直布局 (主内容 + 状态栏)
	content := container.NewBorder(
		nil,          // 顶部
		mv.statusBar, // 底部
		nil,          // 左侧
		nil,          // 右侧
		mainContent,
	)

	return widget.NewSimpleRenderer(content)
}

// SetStatus 设置状态栏消息
func (mv *MainView) SetStatus(message string) {
	mv.statusBar.SetText(message)
}

// Refresh 刷新视图
func (mv *MainView) Refresh() {
	if mv.sidebar != nil {
		mv.sidebar.Refresh()
	}
	if mv.workspace != nil {
		mv.workspace.Refresh()
	}
}
