package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// TerminalWidget 终端显示组件
type TerminalWidget struct {
	widget BaseWidget
	text   string
	lines  []string
	scroll *widget.Scroll
}

// NewTerminalWidget 创建终端组件
func NewTerminalWidget() *TerminalWidget {
	tw := &TerminalWidget{
		lines: make([]string, 0),
	}
	tw.ExtendBaseWidget(tw)
	return tw
}

// CreateRenderer 实现 fyne.Widget 接口
func (tw *TerminalWidget) CreateRenderer() fyne.WidgetRenderer {
	tw.scroll = widget.NewScroll(container.NewVBox())
	return widget.NewSimpleRenderer(tw.scroll)
}

// Append 追加文本
func (tw *TerminalWidget) Append(text string) {
	tw.lines = append(tw.lines, text)

	// 保持最多 1000 行
	if len(tw.lines) > 1000 {
		tw.lines = tw.lines[len(tw.lines)-1000:]
	}

	tw.updateContent()
}

// Clear 清空内容
func (tw *TerminalWidget) Clear() {
	tw.lines = make([]string, 0)
	tw.updateContent()
}

// SetContent 设置内容
func (tw *TerminalWidget) SetContent(content string) {
	tw.lines = splitLines(content)
	tw.updateContent()
}

// updateContent 更新显示内容
func (tw *TerminalWidget) updateContent() {
	if tw.scroll == nil {
		return
	}

	// 创建垂直盒子布局
	vbox := container.NewVBox()
	for _, line := range tw.lines {
		label := widget.NewLabel(line)
		label.TextStyle = fyne.TextStyle{Monospace: true}
		vbox.Add(label)
	}

	tw.scroll.Content = vbox
	tw.scroll.Refresh()

	// 滚动到底部
	tw.scroll.Offset.Y = tw.scroll.Content.MinSize().Height - tw.scroll.Size.Height
	if tw.scroll.Offset.Y < 0 {
		tw.scroll.Offset.Y = 0
	}
	tw.scroll.Refresh()
}

func splitLines(text string) []string {
	var lines []string
	current := ""
	for _, r := range text {
		if r == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}
