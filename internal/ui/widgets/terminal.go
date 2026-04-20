package widgets

import (
	"image/color"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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
	scroll        *container.Scroll
	content       *fyne.Container
	cursorLabel   *widget.Label
}

// TerminalStyle 终端样式
type TerminalStyle struct {
	Foreground color.Color
	Background color.Color
	Bold       bool
	Underline  bool
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
	tw.scroll = container.NewScroll(tw.content)
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

// Append 追加文本（支持 ANSI）
func (tw *TerminalWidget) Append(text string) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	// 解析 ANSI 并添加带样式的行
	tw.parseAndAppend(text)

	// 限制行数
	maxLines := 1000
	if len(tw.lines) > maxLines {
		tw.lines = tw.lines[len(tw.lines)-maxLines:]
	}

	tw.updateContent()
}

// parseAndAppend 解析 ANSI 并追加
func (tw *TerminalWidget) parseAndAppend(text string) {
	style := TerminalStyle{
		Foreground: tw.fgColor,
		Background: tw.bgColor,
	}

	var currentLine strings.Builder
	i := 0
	for i < len(text) {
		if text[i] == '\x1b' && i+1 < len(text) && text[i+1] == '[' {
			// ANSI 转义序列
			end := strings.IndexAny(text[i:], "mKJHfABCDsSu")
			if end == -1 {
				currentLine.WriteByte(text[i])
				i++
				continue
			}

			seq := text[i : i+end+1]
			tw.parseANSI(seq, &style)
			i += end + 1
		} else if text[i] == '\n' {
			// 换行
			tw.lines = append(tw.lines, currentLine.String())
			currentLine.Reset()
			i++
		} else if text[i] == '\r' {
			// 回车
			currentLine.Reset()
			i++
		} else if text[i] == '\t' {
			// 制表符
			currentLine.WriteString("    ")
			i++
		} else {
			currentLine.WriteByte(text[i])
			i++
		}
	}

	// 添加剩余内容
	if currentLine.Len() > 0 {
		tw.lines = append(tw.lines, currentLine.String())
	}
}

// parseANSI 解析 ANSI 序列
func (tw *TerminalWidget) parseANSI(seq string, style *TerminalStyle) {
	if len(seq) < 2 || seq[0] != '\x1b' || seq[1] != '[' {
		return
	}

	// 移除 ESC[ 和后缀
	code := strings.TrimPrefix(seq, "\x1b[")
	code = strings.TrimRight(code, "mKJHfABCDsSu")

	// 处理控制命令
	if strings.HasSuffix(seq, "K") {
		// 清行 - 暂不实现
		return
	}
	if strings.HasSuffix(seq, "J") {
		// 清屏 - 暂不实现
		return
	}
	if strings.HasSuffix(seq, "A") {
		// 光标的上
		tw.CursorUp()
		return
	}
	if strings.HasSuffix(seq, "B") {
		// 光标的下
		tw.CursorDown()
		return
	}
	if strings.HasSuffix(seq, "C") {
		// 光标右移
		tw.CursorRight()
		return
	}
	if strings.HasSuffix(seq, "D") {
		// 光标左移
		tw.CursorLeft()
		return
	}

	// 解析颜色代码
	if code == "" || code == "0" {
		// 重置
		style.Foreground = tw.fgColor
		style.Background = tw.bgColor
		style.Bold = false
		style.Underline = false
		return
	}

	parts := strings.Split(code, ";")
	for _, p := range parts {
		switch p {
		case "1":
			style.Bold = true
		case "4":
			style.Underline = true
		case "30":
			style.Foreground = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
		case "31":
			style.Foreground = color.NRGBA{R: 255, G: 85, B: 85, A: 255}
		case "32":
			style.Foreground = color.NRGBA{R: 85, G: 255, B: 85, A: 255}
		case "33":
			style.Foreground = color.NRGBA{R: 255, G: 255, B: 85, A: 255}
		case "34":
			style.Foreground = color.NRGBA{R: 85, G: 85, B: 255, A: 255}
		case "35":
			style.Foreground = color.NRGBA{R: 255, G: 85, B: 255, A: 255}
		case "36":
			style.Foreground = color.NRGBA{R: 85, G: 255, B: 255, A: 255}
		case "37":
			style.Foreground = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		case "40":
			style.Background = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
		case "41":
			style.Background = color.NRGBA{R: 255, G: 85, B: 85, A: 255}
		case "42":
			style.Background = color.NRGBA{R: 85, G: 255, B: 85, A: 255}
		case "43":
			style.Background = color.NRGBA{R: 255, G: 255, B: 85, A: 255}
		case "44":
			style.Background = color.NRGBA{R: 85, G: 85, B: 255, A: 255}
		case "45":
			style.Background = color.NRGBA{R: 255, G: 85, B: 255, A: 255}
		case "46":
			style.Background = color.NRGBA{R: 85, G: 255, B: 255, A: 255}
		case "47":
			style.Background = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		}
	}
}

// Clear 清空内容
func (tw *TerminalWidget) Clear() {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	tw.lines = make([]string, 0, 500)
	tw.cursorX = 0
	tw.cursorY = 0
	tw.updateContent()
}

// ClearScreen 清屏
func (tw *TerminalWidget) ClearScreen() {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	tw.lines = make([]string, 0, 500)
	tw.cursorX = 0
	tw.cursorY = 0
	tw.updateContent()
}

// MoveCursor 移动光标到指定位置
func (tw *TerminalWidget) MoveCursor(x, y int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if x >= 0 {
		tw.cursorX = x
	}
	if y >= 0 && y < len(tw.lines) {
		tw.cursorY = y
	}
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
	go func() {
		if tw.scroll.Content != nil {
			contentHeight := tw.scroll.Content.MinSize().Height
			viewHeight := tw.scroll.Size().Height
			if contentHeight > viewHeight {
				tw.scroll.Offset.Y = contentHeight - viewHeight
				tw.scroll.Refresh()
			}
		}
	}()
}

// CursorUp 光标的上移
func (tw *TerminalWidget) CursorUp() {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.cursorY > 0 {
		tw.cursorY--
	}
}

// CursorDown 光标的下移
func (tw *TerminalWidget) CursorDown() {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.cursorY < len(tw.lines)-1 {
		tw.cursorY++
	}
}

// CursorLeft 光标左移
func (tw *TerminalWidget) CursorLeft() {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.cursorX > 0 {
		tw.cursorX--
	}
}

// CursorRight 光标右移
func (tw *TerminalWidget) CursorRight() {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if len(tw.lines) > tw.cursorY && tw.cursorX < len(tw.lines[tw.cursorY]) {
		tw.cursorX++
	}
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
