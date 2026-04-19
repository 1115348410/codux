package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// WorkspaceView 工作区视图
type WorkspaceView struct {
	widget.BaseWidget
	window fyne.Window
}

// NewWorkspaceView 创建工作区视图
func NewWorkspaceView(window fyne.Window) *WorkspaceView {
	wv := &WorkspaceView{
		window: window,
	}
	wv.ExtendBaseWidget(wv)
	return wv
}

// CreateRenderer 实现 fyne.Widget 接口
func (wv *WorkspaceView) CreateRenderer() fyne.WidgetRenderer {
	// 空状态提示
	emptyState := widget.NewLabel("选择一个项目开始工作")
	emptyState.TextStyle = fyne.TextStyle{Italic: true}
	emptyState.Alignment = fyne.TextAlignCenter

	content := container.NewCenter(emptyState)

	return widget.NewSimpleRenderer(content)
}
