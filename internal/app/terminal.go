package app

import (
	"sync"

	"github.com/duxweb/codux/internal/services/terminal"
)

// TerminalService 终端管理服务
type TerminalService struct {
	mu       sync.Mutex
	sessions map[string]*terminal.PTY
}

// NewTerminalService 创建终端服务
func NewTerminalService() *TerminalService {
	return &TerminalService{
		sessions: make(map[string]*terminal.PTY),
	}
}

// CreateSession 创建终端会话
func (s *TerminalService) CreateSession(id string, command string, workingDir string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.sessions[id]; exists {
		return nil
	}

	pty := terminal.NewPTY()
	if err := pty.Start(command, workingDir, nil); err != nil {
		return err
	}

	s.sessions[id] = pty
	return nil
}

// GetSession 获取终端会话
func (s *TerminalService) GetSession(id string) *terminal.PTY {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sessions[id]
}

// CloseSession 关闭终端会话
func (s *TerminalService) CloseSession(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if pty, exists := s.sessions[id]; exists {
		pty.Close()
		delete(s.sessions, id)
	}
	return nil
}

// ResizeSession 调整终端大小
func (s *TerminalService) ResizeSession(id string, rows, cols uint16) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if pty, exists := s.sessions[id]; exists {
		return pty.SetWinSize(rows, cols)
	}
	return nil
}

// WriteToSession 写入数据到会话
func (s *TerminalService) WriteToSession(id string, data []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if pty, exists := s.sessions[id]; exists {
		return pty.Write(data)
	}
	return 0, nil
}

// CloseAll 关闭所有会话
func (s *TerminalService) CloseAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, pty := range s.sessions {
		pty.Close()
		delete(s.sessions, id)
	}
}
