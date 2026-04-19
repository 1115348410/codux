package models

// LayoutType 布局类型
type LayoutType string

const (
	LayoutSingle  LayoutType = "single"
	LayoutHSplit  LayoutType = "hsplit"
	LayoutVSplit  LayoutType = "vsplit"
	LayoutComplex LayoutType = "complex"
)

// SplitLayout 分屏布局
type SplitLayout struct {
	ID        string       `json:"id"`
	Type      LayoutType   `json:"type"`
	Left      *SplitLayout `json:"left,omitempty"`
	Right     *SplitLayout `json:"right,omitempty"`
	Top       *SplitLayout `json:"top,omitempty"`
	Bottom    *SplitLayout `json:"bottom,omitempty"`
	SessionID string       `json:"sessionID,omitempty"`
	Ratio     float64      `json:"ratio,omitempty"`
}

// NewSingleLayout 创建单个布局
func NewSingleLayout(sessionID string) *SplitLayout {
	return &SplitLayout{
		Type:      LayoutSingle,
		SessionID: sessionID,
	}
}

// NewHSplitLayout 创建水平分割布局
func NewHSplitLayout(left, right *SplitLayout, ratio float64) *SplitLayout {
	return &SplitLayout{
		Type:  LayoutHSplit,
		Left:  left,
		Right: right,
		Ratio: ratio,
	}
}

// NewVSplitLayout 创建垂直分割布局
func NewVSplitLayout(top, bottom *SplitLayout, ratio float64) *SplitLayout {
	return &SplitLayout{
		Type:   LayoutVSplit,
		Top:    top,
		Bottom: bottom,
		Ratio:  ratio,
	}
}

// GetSessionIDs 获取所有会话 ID
func (l *SplitLayout) GetSessionIDs() []string {
	if l == nil {
		return []string{}
	}

	if l.Type == LayoutSingle {
		return []string{l.SessionID}
	}

	var ids []string
	ids = append(ids, l.Left.GetSessionIDs()...)
	ids = append(ids, l.Right.GetSessionIDs()...)
	ids = append(ids, l.Top.GetSessionIDs()...)
	ids = append(ids, l.Bottom.GetSessionIDs()...)
	return ids
}

// FindLayout 查找包含指定会话的布局
func (l *SplitLayout) FindLayout(sessionID string) *SplitLayout {
	if l == nil {
		return nil
	}

	if l.Type == LayoutSingle && l.SessionID == sessionID {
		return l
	}

	if found := l.Left.FindLayout(sessionID); found != nil {
		return found
	}
	if found := l.Right.FindLayout(sessionID); found != nil {
		return found
	}
	if found := l.Top.FindLayout(sessionID); found != nil {
		return found
	}
	if found := l.Bottom.FindLayout(sessionID); found != nil {
		return found
	}

	return nil
}

// SessionLayout 会话布局信息
type SessionLayout struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectID"`
	Title     string `json:"title"`
	Command   string `json:"command"`
	Active    bool   `json:"active"`
}
