package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/duxweb/codux/internal/app"
	"github.com/duxweb/codux/internal/models"
)

// SidebarView 侧边栏视图
type SidebarView struct {
	widget.BaseWidget
	window      fyne.Window
	store       *app.Store
	projectList *widget.List
	projects    []*models.Project
	selectedID  *models.UUID
}

// NewSidebarView 创建侧边栏视图
func NewSidebarView(window fyne.Window, store *app.Store) *SidebarView {
	sv := &SidebarView{
		window: window,
		store:  store,
	}
	sv.ExtendBaseWidget(sv)
	return sv
}

// CreateRenderer 实现 fyne.Widget 接口
func (sv *SidebarView) CreateRenderer() fyne.WidgetRenderer {
	sv.loadProjects()

	// 侧边栏标题
	header := widget.NewLabel("项目")
	header.TextStyle = fyne.TextStyle{Bold: true}

	// 项目列表
	sv.projectList = widget.NewList(
		func() int {
			return len(sv.projects)
		},
		func() fyne.CanvasObject {
			// 创建列表项容器
			icon := widget.NewIcon(theme.FolderIcon())
			name := widget.NewLabel("Project Name")
			name.TextStyle = fyne.TextStyle{Bold: true}

			return container.NewHBox(icon, name)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id >= 0 && id < len(sv.projects) {
				project := sv.projects[id]
				if hbox, ok := item.(*fyne.Container); ok {
					if icon, ok := hbox.Objects[0].(*widget.Icon); ok {
						icon.SetResource(theme.FolderIcon())
					}
					if label, ok := hbox.Objects[1].(*widget.Label); ok {
						label.SetText(project.Name)
						if sv.selectedID != nil && sv.selectedID.String() == project.ID.String() {
							label.TextStyle = fyne.TextStyle{Bold: true}
						} else {
							label.TextStyle = fyne.TextStyle{}
						}
					}
				}
			}
		},
	)

	sv.projectList.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(sv.projects) {
			project := sv.projects[id]
			sv.store.SelectProject(project.ID)
			sv.selectedID = &project.ID
			sv.projectList.Refresh()
		}
	}

	// 添加项目按钮
	addBtn := widget.NewButtonWithIcon("添加项目", theme.ContentAddIcon(), func() {
		sv.showAddProjectDialog()
	})

	// 垂直布局
	content := container.NewVBox(
		header,
		sv.projectList,
		layout.NewSpacer(),
		addBtn,
	)

	return widget.NewSimpleRenderer(content)
}

// loadProjects 加载项目列表
func (sv *SidebarView) loadProjects() {
	sv.projects = sv.store.Projects()
	selected := sv.store.SelectedProject()
	if selected != nil {
		sv.selectedID = &selected.ID
	}
}

// Refresh 刷新列表
func (sv *SidebarView) Refresh() {
	sv.loadProjects()
	if sv.projectList != nil {
		sv.projectList.Refresh()
	}
}

// showAddProjectDialog 显示添加项目对话框
func (sv *SidebarView) showAddProjectDialog() {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("项目名称")

	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("/path/to/project")

	// 目录选择按钮
	browseBtn := widget.NewButtonWithIcon("浏览...", theme.FolderOpenIcon(), func() {
		d := dialog.NewFolderOpen(func(f fyne.ListableURI, err error) {
			if err != nil {
				return
			}
			if f != nil {
				pathEntry.SetText(f.Path())
			}
		}, sv.window)
		d.Show()
	})

	formItems := []*widget.FormItem{
		widget.NewFormItem("名称", nameEntry),
		widget.NewFormItem("路径", container.NewHBox(pathEntry, browseBtn)),
	}

	d := dialog.NewForm("添加项目", "取消", "添加", formItems, func(confirmed bool) {
		if confirmed && nameEntry.Text != "" && pathEntry.Text != "" {
			project := models.NewProject(nameEntry.Text, pathEntry.Text)
			if err := sv.store.AddProject(project); err != nil {
				dialog.ShowError(err, sv.window)
				return
			}
			sv.Refresh()
		}
	}, sv.window)
	d.Show()
}
