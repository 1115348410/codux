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
	"github.com/duxweb/codux/internal/ui/widgets"
)

// RightPanelView 右侧面板视图
type RightPanelView struct {
	widget.BaseWidget
	window       fyne.Window
	store        *app.Store
	tabs         *container.AppTabs
	gitService   *git.Service
	gitContent   fyne.CanvasObject
	aiContent    fyne.CanvasObject
	diffView     *widgets.DiffView
	fileList     *widget.List
	selectedFile string
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

	// 分支信息
	branchLabel := widget.NewLabel("Branch: -")

	// 按钮
	refreshBtn := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		rpv.refreshGitStatus()
	})

	fetchBtn := widget.NewButton("Fetch", func() {
		rpv.doFetch()
	})
	pullBtn := widget.NewButton("Pull", func() {
		rpv.doPull()
	})
	pushBtn := widget.NewButton("Push", func() {
		rpv.doPush()
	})

	buttonRow := container.NewHBox(refreshBtn, fetchBtn, pullBtn, pushBtn)

	// Diff 查看器
	rpv.diffView = widgets.NewDiffView()

	// 文件列表
	rpv.fileList = widget.NewList(
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

	// 提交信息
	commitEntry := widget.NewEntry()
	commitEntry.SetPlaceHolder("提交信息")
	commitEntry.MultiLine = true

	commitBtn := widget.NewButtonWithIcon("提交", theme.ContentPasteIcon(), func() {
		rpv.doCommit(commitEntry.Text)
	})

	// 组装面板
	content := container.NewVBox(
		container.NewHBox(branchLabel, layout.NewSpacer(), statusLabel),
		widget.NewSeparator(),
		buttonRow,
		widget.NewSeparator(),
		widget.NewLabel("变更文件"),
		rpv.fileList,
		widget.NewSeparator(),
		rpv.diffView,
		widget.NewSeparator(),
		commitEntry,
		commitBtn,
	)

	rpv.gitContent = container.NewBorder(nil, nil, nil, nil, container.NewVScroll(content))
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

	state, err := rpv.gitService.GetRepositoryState(project.Path)
	if err != nil {
		dialog.ShowError(err, rpv.window)
		return
	}

	// 更新文件列表
	allFiles := append(state.ModifiedFiles, state.StagedFiles...)
	allFiles = append(allFiles, state.UntrackedFiles...)

	fileListData := allFiles
	rpv.fileList.Length = func() int { return len(fileListData) }
	rpv.fileList.CreateItem = func() fyne.CanvasObject {
		return container.NewHBox(
			widget.NewIcon(theme.DocumentIcon()),
			widget.NewLabel(""),
			layout.NewSpacer(),
			widget.NewLabel(""),
		)
	}
	rpv.fileList.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
		if id < len(fileListData) {
			file := fileListData[id]
			if hBox, ok := item.(*fyne.Container); ok {
				if label, ok := hBox.Objects[1].(*widget.Label); ok {
					label.SetText(file.Path)
				}
				if statusLabel, ok := hBox.Objects[3].(*widget.Label); ok {
					statusLabel.SetText(file.Status)
				}
			}
		}
	}
	rpv.fileList.Refresh()

	// 文件选择回调
	rpv.fileList.OnSelected = func(id widget.ListItemID) {
		if id < len(fileListData) {
			file := fileListData[id]
			rpv.selectedFile = file.Path
			diff, _ := rpv.gitService.GetDiff(project.Path, file.Path)
			if rpv.diffView != nil {
				rpv.diffView.SetDiff(diff)
			}
		}
	}
}

// doFetch 执行 Fetch
func (rpv *RightPanelView) doFetch() {
	project := rpv.store.SelectedProject()
	if project == nil {
		return
	}
	err := rpv.gitService.Fetch(project.Path)
	if err != nil {
		dialog.ShowError(err, rpv.window)
	}
}

// doPull 执行 Pull
func (rpv *RightPanelView) doPull() {
	project := rpv.store.SelectedProject()
	if project == nil {
		return
	}
	err := rpv.gitService.Pull(project.Path)
	if err != nil {
		dialog.ShowError(err, rpv.window)
	} else {
		dialog.ShowInformation("Pull 成功", "代码已更新", rpv.window)
	}
}

// doPush 执行 Push
func (rpv *RightPanelView) doPush() {
	project := rpv.store.SelectedProject()
	if project == nil {
		return
	}
	err := rpv.gitService.Push(project.Path)
	if err != nil {
		dialog.ShowError(err, rpv.window)
	} else {
		dialog.ShowInformation("Push 成功", "代码已推送", rpv.window)
	}
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
	rpv.refreshGitStatus()
}
