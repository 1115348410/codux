package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// SidebarView 侧边栏视图
type SidebarView struct {
	widget.BaseWidget
	window      fyne.Window
	projectList *widget.List
	projects    []string
}

// NewSidebarView 创建侧边栏视图
func NewSidebarView(window fyne.Window) *SidebarView {
	sv := &SidebarView{
		window:   window,
		projects: []string{},
	}
	sv.ExtendBaseWidget(sv)
	return sv
}

// CreateRenderer 实现 fyne.Widget 接口
func (sv *SidebarView) CreateRenderer() fyne.WidgetRenderer {
	// 侧边栏标题
	header := widget.NewLabel("项目")
	header.TextStyle = fyne.TextStyle{Bold: true}

	// 项目列表
	sv.projectList = widget.NewList(
		func() int {
			return len(sv.projects)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("示例项目")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if label, ok := item.(*widget.Label); ok {
				label.SetText(sv.projects[id])
			}
		},
	)

	// 添加项目按钮
	addBtn := widget.NewButtonWithIcon("添加项目", theme.ContentAddIcon(), func() {
		sv.showAddProjectDialog()
	})

	// 垂直布局
	content := container.NewVBox(
		header,
		sv.projectList,
		addBtn,
	)

	return widget.NewSimpleRenderer(content)
}

// showAddProjectDialog 显示添加项目对话框
func (sv *SidebarView) showAddProjectDialog() {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("项目名称")

	formItems := []*widget.FormItem{
		widget.NewFormItem("名称", nameEntry),
	}

	form := widget.NewForm(formItems...)
	form.OnSubmit = func() {
		if nameEntry.Text != "" {
			sv.projects = append(sv.projects, nameEntry.Text)
			sv.projectList.Refresh()
		}
	}

	d := dialog.NewForm("添加项目", "取消", "添加", formItems, func(confirmed bool) {
		if confirmed && nameEntry.Text != "" {
			sv.projects = append(sv.projects, nameEntry.Text)
			sv.projectList.Refresh()
		}
	}, sv.window)
	d.Show()
}
