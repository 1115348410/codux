package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/duxweb/codux/internal/app"
	"github.com/duxweb/codux/internal/services/git"
)

// RightPanelView 右侧面板视图
type RightPanelView struct {
	widget.BaseWidget
	window     fyne.Window
	store      *app.Store
	tabs       *container.AppTabs
	gitService *git.Service
	gitContent fyne.CanvasObject
	aiContent  fyne.CanvasObject
}

// NewRightPanelView 创建右侧面板视图
func NewRightPanelView(window fyne.Window, store *app.Store) *RightPanelView {
	rpv := &RightPanelView{
		window:     window,
		store:      store,
		gitService: git.NewService(),
	}
	rpv.ExtendBaseWidget(rpv)
	return rpv
}

// CreateRenderer 实现 fyne.Widget 接口
func (rpv *RightPanelView) CreateRenderer() fyne.WidgetRenderer {
	rpv.createGitPanel()
	rpv.createAIPanel()

	rpv.tabs = container.NewAppTabs(
		container.NewTabItem("Git", rpv.gitContent),
		container.NewTabItem("AI", rpv.aiContent),
	)

	rpv.tabs.SetTabLocation(container.TabLocationLeading)

	content := container.NewMax(rpv.tabs)
	return widget.NewSimpleRenderer(content)
}

// createGitPanel 创建 Git 面板
func (rpv *RightPanelView) createGitPanel() {
	infoLabel := widget.NewLabel("未检测到 Git 仓库")
	infoLabel.Alignment = fyne.TextAlignCenter
	infoLabel.TextStyle = fyne.TextStyle{Italic: true}

	fileList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.DocumentIcon()),
				widget.NewLabel("file.txt"),
				layout.NewSpacer(),
				widget.NewLabel("M"),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {},
	)

	statusBtn := widget.NewButtonWithIcon("刷新状态", theme.ViewRefreshIcon(), func() {
		rpv.refreshGitStatus()
	})

	stageAllBtn := widget.NewButton("全部暂存", func() {})
	commitEntry := widget.NewEntry()
	commitEntry.SetPlaceHolder("提交信息")
	commitEntry.MultiLine = true

	commitBtn := widget.NewButtonWithIcon("提交", theme.ContentPasteIcon(), func() {
		if commitEntry.Text != "" {
			rpv.doCommit(commitEntry.Text)
			commitEntry.Text = ""
			commitEntry.Refresh()
		}
	})

	form := container.NewVBox(
		container.NewHBox(statusBtn, stageAllBtn),
		widget.NewSeparator(),
		fileList,
		widget.NewSeparator(),
		commitEntry,
		commitBtn,
	)

	rpv.gitContent = container.NewVScroll(
		container.NewBorder(
			widget.NewLabel("Git 状态"),
			nil,
			nil,
			nil,
			container.NewVScroll(form),
		),
	)
}

// createAIPanel 创建 AI 统计面板
func (rpv *RightPanelView) createAIPanel() {
	placeholder := widget.NewLabel("AI 使用统计\n\n启动 AI 工具后会显示 Token 使用情况")
	placeholder.Alignment = fyne.TextAlignCenter
	placeholder.TextStyle = fyne.TextStyle{Italic: true}

	rpv.aiContent = container.NewCenter(placeholder)
}

// refreshGitStatus 刷新 Git 状态
func (rpv *RightPanelView) refreshGitStatus() {
	project := rpv.store.SelectedProject()
	if project == nil {
		return
	}

	_, err := rpv.gitService.GetRepositoryState(project.Path)
	if err != nil {
		return
	}

	// TODO: 更新 UI 显示状态
}

// doCommit 执行提交
func (rpv *RightPanelView) doCommit(message string) {
	project := rpv.store.SelectedProject()
	if project == nil {
		return
	}

	err := rpv.gitService.Commit(project.Path, message)
	if err != nil {
		dialog.ShowError(err, rpv.window)
		return
	}

	dialog.ShowInformation("提交成功", "代码已提交", rpv.window)
}
