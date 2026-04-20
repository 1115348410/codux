package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// SplitContainer 可拖拽的分屏容器
type SplitContainer struct {
	widget.BaseWidget
	left       fyne.CanvasObject
	right      fyne.CanvasObject
	splitter   fyne.CanvasObject
	ratio      float64
	horizontal bool
	dragging   bool
	lastX      float32
	lastY      float32
	onChange   func(float64)
}

// NewSplitContainer 创建分屏容器
func NewSplitContainer(left, right fyne.CanvasObject, horizontal bool, ratio float64) *SplitContainer {
	sc := &SplitContainer{
		left:       left,
		right:      right,
		ratio:      ratio,
		horizontal: horizontal,
	}
	sc.ExtendBaseWidget(sc)
	return sc
}

// SetOnChange 设置变化回调
func (sc *SplitContainer) SetOnChange(fn func(float64)) {
	sc.onChange = fn
}

// CreateRenderer 创建渲染器
func (sc *SplitContainer) CreateRenderer() fyne.WidgetRenderer {
	sc.splitter = sc.createSplitter()
	return &splitContainerRenderer{
		container: sc,
		objects:   []fyne.CanvasObject{sc.left, sc.splitter, sc.right},
	}
}

// createSplitter 创建分隔条
func (sc *SplitContainer) createSplitter() fyne.CanvasObject {
	s := &splitter{
		isHorizontal: sc.horizontal,
		onDrag: func(deltaX, deltaY float32) {
			sc.handleDrag(deltaX, deltaY)
		},
	}
	s.ExtendBaseWidget(s)
	return s
}

// handleDrag 处理拖拽
func (sc *SplitContainer) handleDrag(deltaX, deltaY float32) {
	size := sc.Size()
	if sc.horizontal {
		deltaRatio := float64(deltaX) / float64(size.Width)
		sc.ratio += deltaRatio
		if sc.ratio < 0.1 {
			sc.ratio = 0.1
		}
		if sc.ratio > 0.9 {
			sc.ratio = 0.9
		}
	} else {
		deltaRatio := float64(deltaY) / float64(size.Height)
		sc.ratio += deltaRatio
		if sc.ratio < 0.1 {
			sc.ratio = 0.1
		}
		if sc.ratio > 0.9 {
			sc.ratio = 0.9
		}
	}

	if sc.onChange != nil {
		sc.onChange(sc.ratio)
	}

	sc.Refresh()
}

func (sc *SplitContainer) leftSize(size fyne.Size) fyne.Size {
	if sc.horizontal {
		return fyne.NewSize(size.Width*float32(sc.ratio), size.Height)
	}
	return fyne.NewSize(size.Width, size.Height*float32(sc.ratio))
}

func (sc *SplitContainer) rightSize(size fyne.Size) fyne.Size {
	if sc.horizontal {
		return fyne.NewSize(size.Width*(1-float32(sc.ratio)), size.Height)
	}
	return fyne.NewSize(size.Width, size.Height*(1-float32(sc.ratio)))
}

func (sc *SplitContainer) leftPos(size fyne.Size) fyne.Position {
	return fyne.NewPos(0, 0)
}

func (sc *SplitContainer) rightPos(size fyne.Size) fyne.Position {
	if sc.horizontal {
		return fyne.NewPos(size.Width*float32(sc.ratio)+8, 0) // 8 = splitter width
	}
	return fyne.NewPos(0, size.Height*float32(sc.ratio)+8)
}

func (sc *SplitContainer) splitterPos(size fyne.Size) fyne.Position {
	if sc.horizontal {
		return fyne.NewPos(size.Width*float32(sc.ratio), 0)
	}
	return fyne.NewPos(0, size.Height*float32(sc.ratio))
}

// splitter 分隔条
type splitter struct {
	widget.BaseWidget
	isHorizontal bool
	onDrag       func(float32, float32)
}

func (s *splitter) CreateRenderer() fyne.WidgetRenderer {
	rect := canvas.NewRectangle(theme.ShadowColor())
	rect.Resize(fyne.NewSize(8, 8))
	return widget.NewSimpleRenderer(rect)
}

func (s *splitter) MinSize() fyne.Size {
	if s.isHorizontal {
		return fyne.NewSize(8, 100)
	}
	return fyne.NewSize(100, 8)
}

func (s *splitter) MouseIn(e *desktop.MouseEvent) {
	// 改变光标样式
}

func (s *splitter) MouseMoved(e *desktop.MouseEvent) {
	if s.onDrag != nil {
		s.onDrag(e.Position.X, e.Position.Y)
	}
}

func (s *splitter) MouseOut() {
	// 恢复光标
}

func (s *splitter) Dragged(e *fyne.DragEvent) {
	if s.onDrag != nil {
		s.onDrag(e.Dragged.DX, e.Dragged.DY)
	}
}

func (s *splitter) DragEnd() {
	// 拖拽结束
}

// splitContainerRenderer 分屏容器渲染器
type splitContainerRenderer struct {
	container *SplitContainer
	objects   []fyne.CanvasObject
}

func (r *splitContainerRenderer) Layout(size fyne.Size) {
	leftSize := r.container.leftSize(size)
	rightSize := r.container.rightSize(size)
	leftPos := r.container.leftPos(size)
	rightPos := r.container.rightPos(size)
	splitterPos := r.container.splitterPos(size)

	r.objects[0].Resize(leftSize)
	r.objects[0].Move(leftPos)

	// Splitter size
	var splitterSize fyne.Size
	if r.container.horizontal {
		splitterSize = fyne.NewSize(8, size.Height)
	} else {
		splitterSize = fyne.NewSize(size.Width, 8)
	}
	r.objects[1].Resize(splitterSize)
	r.objects[1].Move(splitterPos)

	r.objects[2].Resize(rightSize)
	r.objects[2].Move(rightPos)
}

func (r *splitContainerRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, 200)
}

func (r *splitContainerRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *splitContainerRenderer) Refresh() {
	r.container.Refresh()
}

func (r *splitContainerRenderer) Destroy() {
}
