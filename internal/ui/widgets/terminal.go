package widgets

import (
	"image/color"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// TerminalWidget 终端显示组件
type TerminalWidget struct {
	widget.BaseWidget
	mu            sync.RWMutex
	lines         []string
	cursorX       int
	cursorY       int
	cursorVisible bool
	fgColor       color.Color
	bgColor       color.Color
	onInput       func([]byte)
	scroll        *widget.Scroll
	content       *fyne.Container
}

// NewTerminalWidget 创建终端组件
func NewTerminalWidget() *TerminalWidget {
	tw := &TerminalWidget{
		lines:         make([]string, 0, 500),
		cursorX:       0,
		cursorY:       0,
		cursorVisible: true,
		fgColor:       color.NRGBA{R: 235, G: 235, B: 235, A: 255},
		bgColor:       color.NRGBA{R: 30, G: 30, B: 30, A: 255},
	}
	tw.ExtendBaseWidget(tw)
	return tw
}

// CreateRenderer 实现 fyne.Widget 接口
func (tw *TerminalWidget) CreateRenderer() fyne.WidgetRenderer {
	tw.content = container.NewVBox()
	tw.scroll = widget.NewScroll(tw.content)
	tw.scroll.Direction = widget.ScrollBoth
	return widget.NewSimpleRenderer(tw.scroll)
}

// SetOnInput 设置输入回调
func (tw *TerminalWidget) SetOnInput(fn func([]byte)) {
	tw.onInput = fn
}

// TypedRune 处理字符输入
func (tw *TerminalWidget) TypedRune(r rune) {
	if tw.onInput != nil {
		tw.onInput([]byte(string(r)))
	}
}

// TypedKey 处理按键
func (tw *TerminalWidget) TypedKey(key *fyne.KeyEvent) {
	if tw.onInput == nil {
		return
	}

	switch key.Name {
	case fyne.KeyEnter:
		tw.onInput([]byte{0x0D})
	case fyne.KeyBackspace:
		tw.onInput([]byte{0x7F})
	case fyne.KeyTab:
		tw.onInput([]byte{0x09})
	case fyne.KeyEscape:
		tw.onInput([]byte{0x1B})
	case fyne.KeyUp:
		tw.onInput([]byte{0x1B, 0x5B, 0x41})
	case fyne.KeyDown:
		tw.onInput([]byte{0x1B, 0x5B, 0x42})
	case fyne.KeyRight:
		tw.onInput([]byte{0x1B, 0x5B, 0x43})
	case fyne.KeyLeft:
		tw.onInput([]byte{0x1B, 0x5B, 0x44})
	}
}

// Append 追加文本
func (tw *TerminalWidget) Append(text string) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	// 处理 ANSI 转义序列
	clean := stripANSI(text)

	// 分割行
	parts := strings.Split(clean, "\n")
	for i, part := range parts {
		if i == 0 && len(tw.lines) > 0 {
			tw.lines[len(tw.lines)-1] += part
		} else if part != "" || i > 0 {
			tw.lines = append(tw.lines, part)
		}
	}

	// 限制行数
	maxLines := 1000
	if len(tw.lines) > maxLines {
		tw.lines = tw.lines[len(tw.lines)-maxLines:]
	}

	tw.updateContent()
}

// Clear 清空内容
func (tw *TerminalWidget) Clear() {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	tw.lines = make([]string, 0, 500)
	tw.updateContent()
}

// SetContent 设置内容
func (tw *TerminalWidget) SetContent(content string) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	tw.lines = strings.Split(content, "\n")
	tw.updateContent()
}

// GetLines 获取所有行
func (tw *TerminalWidget) GetLines() []string {
	tw.mu.RLock()
	defer tw.mu.RUnlock()
	result := make([]string, len(tw.lines))
	copy(result, tw.lines)
	return result
}

// updateContent 更新显示内容
func (tw *TerminalWidget) updateContent() {
	if tw.content == nil {
		return
	}

	tw.content.Objects = make([]fyne.CanvasObject, 0, len(tw.lines))
	for _, line := range tw.lines {
		label := widget.NewLabel(line)
		label.TextStyle = fyne.TextStyle{Monospace: true}
		tw.content.Objects = append(tw.content.Objects, label)
	}

	tw.content.Refresh()
	tw.scroll.Refresh()

	// 滚动到底部
	if tw.scroll.Content != nil {
		contentHeight := tw.scroll.Content.MinSize().Height
		viewHeight := tw.scroll.Size.Height
		if contentHeight > viewHeight {
			tw.scroll.Offset.Y = contentHeight - viewHeight
		}
	}
	tw.scroll.Refresh()
}

// CursorUp 光标的上移
func (tw *TerminalWidget) CursorUp() {
	// TODO: 实现光标移动
}

// CursorDown 光标的下移
func (tw *TerminalWidget) CursorDown() {
	// TODO: 实现光标移动
}

// CursorLeft 光标左移
func (tw *TerminalWidget) CursorLeft() {
	// TODO: 实现光标移动
}

// CursorRight 光标右移
func (tw *TerminalWidget) CursorRight() {
	// TODO: 实现光标移动
}

// SetColor 设置颜色
func (tw *TerminalWidget) SetColor(fg, bg color.Color) {
	tw.fgColor = fg
	tw.bgColor = bg
}

// Focusable 实现 Focusable 接口
func (tw *TerminalWidget) FocusGained() {
	// TODO: 显示光标
}

func (tw *TerminalWidget) FocusLost() {
	// TODO: 隐藏光标
}

// 移除 ANSI 转义序列
func stripANSI(s string) string {
	var result strings.Builder
	inEscape := false
	for _, r := range s {
		if r == 0x1B {
			inEscape = true
			continue
		}
		if inEscape && (r >= 0x30 && r <= 0x7E) {
			inEscape = false
			continue
		}
		if inEscape && (r == 0x5B || r == 0x28 || r == 0x29) {
			continue
		}
		if !inEscape {
			result.WriteRune(r)
		}
	}
	return result.String()
}
