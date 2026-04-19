package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/duxweb/codux/internal/app"
)

// WorkspaceView 工作区视图
type WorkspaceView struct {
	widget.BaseWidget
	window fyne.Window
	store  *app.Store
}

// NewWorkspaceView 创建工作区视图
func NewWorkspaceView(window fyne.Window, store *app.Store) *WorkspaceView {
	wv := &WorkspaceView{
		window: window,
		store:  store,
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

// Refresh 刷新视图
func (wv *WorkspaceView) Refresh() {
	// TODO: 根据选中的项目更新工作区
}
