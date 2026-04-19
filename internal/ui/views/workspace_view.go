package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/duxweb/codux/internal/app"
)

// WorkspaceView 工作区视图
type WorkspaceView struct {
	widget.BaseWidget
	window  fyne.Window
	store   *app.Store
	content fyne.CanvasObject
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
	wv.updateContent()
	return widget.NewSimpleRenderer(wv.content)
}

// Refresh 刷新视图
func (wv *WorkspaceView) Refresh() {
	wv.updateContent()
}

func (wv *WorkspaceView) updateContent() {
	project := wv.store.SelectedProject()

	if project == nil {
		// 空状态
		emptyState := container.NewCenter(
			container.NewVBox(
				widget.NewIcon(theme.FolderIcon()),
				widget.NewLabel("选择一个项目开始工作"),
			),
		)
		wv.content = emptyState
	} else {
		// 显示终端区域
		terminalArea := container.NewMax(
			widget.NewLabel("终端区域 - 项目：" + project.Name),
		)
		wv.content = terminalArea
	}
}
