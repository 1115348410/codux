package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/duxweb/codux/internal/app"
)

// RightPanelView 右侧面板视图
type RightPanelView struct {
	widget.BaseWidget
	window fyne.Window
	store  *app.Store
	tabs   *container.AppTabs
}

// NewRightPanelView 创建右侧面板视图
func NewRightPanelView(window fyne.Window, store *app.Store) *RightPanelView {
	rpv := &RightPanelView{
		window: window,
		store:  store,
	}
	rpv.ExtendBaseWidget(rpv)
	return rpv
}

// CreateRenderer 实现 fyne.Widget 接口
func (rpv *RightPanelView) CreateRenderer() fyne.WidgetRenderer {
	// Git 面板占位
	gitPanel := widget.NewLabel("Git 面板\n\n选择一个项目查看 Git 状态")
	gitPanel.Alignment = fyne.TextAlignCenter

	// AI 统计面板占位
	aiPanel := widget.NewLabel("AI 统计\n\nAI 工具使用统计将显示在这里")
	aiPanel.Alignment = fyne.TextAlignCenter

	// 创建标签页
	rpv.tabs = container.NewAppTabs(
		container.NewTabItem("Git", gitPanel),
		container.NewTabItem("AI", aiPanel),
	)

	content := container.NewMax(rpv.tabs)

	return widget.NewSimpleRenderer(content)
}
