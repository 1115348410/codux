package app

import (
	"sync"

	"github.com/duxweb/codux/internal/models"
)

// SplitService 分屏布局服务
type SplitService struct {
	mu      sync.Mutex
	layouts map[string]*models.SplitLayout // projectID -> layout
	active  map[string]string              // projectID -> activeSessionID
}

// NewSplitService 创建分屏服务
func NewSplitService() *SplitService {
	return &SplitService{
		layouts: make(map[string]*models.SplitLayout),
		active:  make(map[string]string),
	}
}

// GetLayout 获取项目布局
func (s *SplitService) GetLayout(projectID string) *models.SplitLayout {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.layouts[projectID]
}

// SetLayout 设置项目布局
func (s *SplitService) SetLayout(projectID string, layout *models.SplitLayout) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.layouts[projectID] = layout
}

// GetActiveSession 获取活跃会话
func (s *SplitService) GetActiveSession(projectID string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.active[projectID]
}

// SetActiveSession 设置活跃会话
func (s *SplitService) SetActiveSession(projectID, sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.active[projectID] = sessionID
}

// CreateSplit 创建分屏
func (s *SplitService) CreateSplit(projectID string, sessionID string, horizontal bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current := s.layouts[projectID]
	if current == nil {
		s.layouts[projectID] = models.NewSingleLayout(sessionID)
		s.active[projectID] = sessionID
		return nil
	}

	newLayout := models.NewSingleLayout(sessionID)

	if horizontal {
		s.layouts[projectID] = models.NewHSplitLayout(current, newLayout, 0.5)
	} else {
		s.layouts[projectID] = models.NewVSplitLayout(current, newLayout, 0.5)
	}

	s.active[projectID] = sessionID
	return nil
}

// CloseSession 关闭会话
func (s *SplitService) CloseSession(projectID, sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO: 从布局中移除会话
	delete(s.active, projectID)
}

// ResizeSplit 调整分屏大小
func (s *SplitService) ResizeSplit(projectID string, sessionID string, ratio float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO: 实现分屏大小调整
}
