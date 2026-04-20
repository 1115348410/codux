package views

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/duxweb/codux/internal/app"
	"github.com/duxweb/codux/internal/ui/widgets"
)

// WorkspaceView 工作区视图
type WorkspaceView struct {
	widget.BaseWidget
	window    fyne.Window
	store     *app.Store
	splitView fyne.CanvasObject
	terminal  *widgets.TerminalWidget
	sessionID string
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
	// 使用项目 ID 作为会话 ID
	wv.sessionID = projectID

	// 创建终端组件
	wv.terminal = widgets.NewTerminalWidget()
	wv.terminal.SetOnInput(func(data []byte) {
		// 发送输入到 PTY
		wv.store.Terminal().WriteToSession(wv.sessionID, data)
	})

	// 启动 PTY 进程（如果尚未启动）
	termService := wv.store.Terminal()
	pty := termService.GetSession(wv.sessionID)
	if pty == nil {
		// 创建新的终端会话，默认启动 shell
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/bash"
		}
		err := termService.CreateSession(wv.sessionID, shell, "")
		if err != nil {
			// 创建失败，显示错误信息
			errLabel := widget.NewLabel("Failed to start terminal: " + err.Error())
			errLabel.TextStyle = fyne.TextStyle{Monospace: true}
			return container.NewScroll(errLabel)
		}
		pty = termService.GetSession(wv.sessionID)
	}

	// 读取终端输出
	if pty != nil {
		go func() {
			buffer := make([]byte, 4096)
			for {
				n, err := pty.Read(buffer)
				if err != nil || n == 0 {
					break
				}
				// 在主线程中更新 UI
				fyne.Do(func() {
					wv.terminal.Append(string(buffer[:n]))
				})
			}
		}()
	}

	scroll := container.NewScroll(wv.terminal)
	return scroll
}
