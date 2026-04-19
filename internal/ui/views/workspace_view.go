package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/duxweb/codux/internal/app"
)

// WorkspaceView 工作区视图
type WorkspaceView struct {
	widget.BaseWidget
	window    fyne.Window
	store     *app.Store
	splitView *fyne.Container
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
	return widget.NewSimpleRenderer(wv.splitView)
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
		wv.splitView = emptyState
		return
	}

	// 创建工作区工具栏
	toolbar := wv.createToolbar()

	// 创建终端区域
	terminal := wv.createTerminal(project.ID.String())

	// 垂直布局
	wv.splitView = container.NewBorder(
		toolbar,
		nil,
		nil,
		nil,
		terminal,
	)
}

func (wv *WorkspaceView) createToolbar() fyne.CanvasObject {
	newSplitBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		// TODO: 创建分屏
	})

	closeSplitBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		// TODO: 关闭分屏
	})

	splitSelect := widget.NewSelect([]string{"水平分割", "垂直分割"}, func(value string) {})

	return container.NewHBox(
		widget.NewLabel("终端"),
		layout.NewSpacer(),
		splitSelect,
		newSplitBtn,
		closeSplitBtn,
	)
}

func (wv *WorkspaceView) createTerminal(projectID string) fyne.CanvasObject {
	// TODO: 集成 terminal widget
	terminal := widget.NewLabel("终端显示区域")
	terminal.TextStyle = fyne.TextStyle{Monospace: true}

	scroll := container.NewScroll(terminal)

	return scroll
}
