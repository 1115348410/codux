package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// RightPanelView 右侧面板视图
type RightPanelView struct {
	widget.BaseWidget
	window fyne.Window
	tabs   *container.AppTabs
}

// NewRightPanelView 创建右侧面板视图
func NewRightPanelView(window fyne.Window) *RightPanelView {
	rpv := &RightPanelView{
		window: window,
	}
	rpv.ExtendBaseWidget(rpv)
	return rpv
}

// CreateRenderer 实现 fyne.Widget 接口
func (rpv *RightPanelView) CreateRenderer() fyne.WidgetRenderer {
	// 创建标签页
	rpv.tabs = container.NewAppTabs(
		container.NewTabItem("Git", widget.NewLabel("Git 面板")),
		container.NewTabItem("AI", widget.NewLabel("AI 统计")),
	)

	content := container.NewMax(rpv.tabs)

	return widget.NewSimpleRenderer(content)
}
