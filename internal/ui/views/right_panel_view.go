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
	// 状态标签
	statusLabel := widget.NewLabel("未选择项目")
	statusLabel.TextStyle = fyne.TextStyle{Italic: true}

	// 文件列表
	fileList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.DocumentIcon()),
				widget.NewLabel("file"),
				layout.NewSpacer(),
				widget.NewLabel("M"),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {},
	)

	// 分支信息
	branchLabel := widget.NewLabel("Branch: main")

	// 按钮
	refreshBtn := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		rpv.refreshGitStatus()
	})

	stashBtn := widget.NewButton("Stash", func() {})
	fetchBtn := widget.NewButton("Fetch", func() {})
	pullBtn := widget.NewButton("Pull", func() {})
	pushBtn := widget.NewButton("Push", func() {})

	buttonRow := container.NewHBox(refreshBtn, stashBtn, fetchBtn, pullBtn, pushBtn)

	// 提交信息
	commitEntry := widget.NewEntry()
	commitEntry.SetPlaceHolder("提交信息")
	commitEntry.MultiLine = true

	commitBtn := widget.NewButtonWithIcon("提交", theme.ContentPasteIcon(), func() {
		rpv.doCommit(commitEntry.Text)
	})

	// 历史记录
	historyList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabel("commit hash"),
				widget.NewLabel("commit message"),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {},
	)

	// 组装面板
	content := container.NewVBox(
		container.NewHBox(branchLabel, layout.NewSpacer(), statusLabel),
		widget.NewSeparator(),
		buttonRow,
		widget.NewSeparator(),
		widget.NewLabel("变更文件"),
		fileList,
		widget.NewSeparator(),
		commitEntry,
		commitBtn,
		widget.NewSeparator(),
		widget.NewLabel("提交历史"),
		historyList,
	)

	rpv.gitContent = container.NewVScroll(content)
}

// createAIPanel 创建 AI 统计面板
func (rpv *RightPanelView) createAIPanel() {
	// 今日统计
	todayStats := widget.NewLabel("今日未使用 AI 工具")
	todayStats.TextStyle = fyne.TextStyle{Bold: true}

	// 等级显示
	levelLabel := widget.NewLabel("等级：Iron")
	levelLabel.TextStyle = fyne.TextStyle{Bold: true}

	// 进度指示
	progressContainer := container.NewHBox(
		widget.NewProgressBar(),
	)

	// 工具列表
	toolList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.ComputerIcon()),
				widget.NewLabel("Claude Code"),
				layout.NewSpacer(),
				widget.NewLabel("0 tokens"),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {},
	)

	// 刷新按钮
	refreshBtn := widget.NewButtonWithIcon("刷新", theme.ViewRefreshIcon(), func() {
		// TODO: 刷新 AI 统计
	})

	content := container.NewVBox(
		container.NewHBox(widget.NewLabel("今日统计"), layout.NewSpacer(), refreshBtn),
		todayStats,
		widget.NewSeparator(),
		levelLabel,
		progressContainer,
		widget.NewSeparator(),
		widget.NewLabel("工具使用"),
		toolList,
	)

	rpv.aiContent = container.NewVScroll(content)
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

	// TODO: 更新 UI
}

// doCommit 执行提交
func (rpv *RightPanelView) doCommit(message string) {
	project := rpv.store.SelectedProject()
	if project == nil {
		return
	}

	if message == "" {
		dialog.ShowError(nil, rpv.window)
		return
	}

	err := rpv.gitService.Commit(project.Path, message)
	if err != nil {
		dialog.ShowError(err, rpv.window)
		return
	}

	dialog.ShowInformation("提交成功", "代码已提交", rpv.window)
}
