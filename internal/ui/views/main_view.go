package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// MainView 主视图
type MainView struct {
	widget.BaseWidget
	window     fyne.Window
	sidebar    *SidebarView
	workspace  *WorkspaceView
	rightPanel *RightPanelView
	statusBar  *widget.Label
}

// NewMainView 创建主视图
func NewMainView(window fyne.Window) *MainView {
	mv := &MainView{
		window: window,
	}
	mv.ExtendBaseWidget(mv)
	return mv
}

// CreateRenderer 实现 fyne.Widget 接口
func (mv *MainView) CreateRenderer() fyne.WidgetRenderer {
	// 创建侧边栏
	mv.sidebar = NewSidebarView(mv.window)

	// 创建工作区
	mv.workspace = NewWorkspaceView(mv.window)

	// 创建右侧面板
	mv.rightPanel = NewRightPanelView(mv.window)

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

	// 获取右侧面板的引用以设置比例
	if rightSplit, ok := mainContent.Trailing.(*fyne.Container); ok {
		if innerSplit, ok := rightSplit.Objects[0].(*fyne.Container); ok {
			if split, ok := innerSplit.Objects[0].(*container.Split); ok {
				split.SetOffset(0.7)
			}
		}
	}

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
