package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strings"
)

// DiffView Diff 查看器组件
type DiffView struct {
	widget.BaseWidget
	diff    string
	content *fyne.Container
	scroll  *container.Scroll
}

// NewDiffView 创建 Diff 查看器
func NewDiffView() *DiffView {
	dv := &DiffView{
		content: container.NewVBox(),
	}
	dv.ExtendBaseWidget(dv)
	return dv
}

// SetDiff 设置 Diff 内容
func (dv *DiffView) SetDiff(diff string) {
	dv.diff = diff
	dv.render()
}

// render 渲染 Diff
func (dv *DiffView) render() {
	dv.content.Objects = make([]fyne.CanvasObject, 0)

	if dv.diff == "" {
		dv.content.Objects = append(dv.content.Objects, widget.NewLabel("没有差异"))
		return
	}

	lines := strings.Split(dv.diff, "\n")
	for _, line := range lines {
		label := widget.NewLabel(line)
		label.TextStyle = fyne.TextStyle{Monospace: true}

		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			// 添加行 - 绿色背景
			dv.content.Objects = append(dv.content.Objects, container.NewStack(
				canvas.NewRectangle(color.NRGBA{R: 0, G: 100, B: 0, A: 50}),
				label,
			))
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			// 删除行 - 红色背景
			dv.content.Objects = append(dv.content.Objects, container.NewStack(
				canvas.NewRectangle(color.NRGBA{R: 100, G: 0, B: 0, A: 50}),
				label,
			))
		} else if strings.HasPrefix(line, "@@") {
			// 变更位置 - 紫色
			label.TextStyle.Bold = true
			dv.content.Objects = append(dv.content.Objects, label)
		} else if strings.HasPrefix(line, "diff") || strings.HasPrefix(line, "index") {
			// 文件头 - 灰色
			label.TextStyle.Italic = true
			dv.content.Objects = append(dv.content.Objects, label)
		} else {
			// 上下文 - 普通文本
			dv.content.Objects = append(dv.content.Objects, label)
		}
	}

	if dv.scroll != nil {
		dv.scroll.Refresh()
	}
}

// CreateRenderer 创建渲染器
func (dv *DiffView) CreateRenderer() fyne.WidgetRenderer {
	dv.scroll = container.NewScroll(dv.content)
	return widget.NewSimpleRenderer(dv.scroll)
}
