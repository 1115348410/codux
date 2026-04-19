package ai

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// UsageRecord AI 使用记录
type UsageRecord struct {
	ID           int64     `json:"id"`
	ProjectID    string    `json:"projectID"`
	SessionID    string    `json:"sessionID"`
	Tool         string    `json:"tool"`
	TokensInput  int64     `json:"tokensInput"`
	TokensOutput int64     `json:"tokensOutput"`
	Cost         float64   `json:"cost"`
	CreatedAt    time.Time `json:"createdAt"`
}

// DailyUsage 每日用量统计
type DailyUsage struct {
	Date         string  `json:"date"`
	TotalTokens  int64   `json:"totalTokens"`
	InputTokens  int64   `json:"inputTokens"`
	OutputTokens int64   `json:"outputTokens"`
	Cost         float64 `json:"cost"`
	SessionCount int     `json:"sessionCount"`
}

// ToolStats 工具统计
type ToolStats struct {
	Tool         string  `json:"tool"`
	TotalTokens  int64   `json:"totalTokens"`
	InputTokens  int64   `json:"inputTokens"`
	OutputTokens int64   `json:"outputTokens"`
	Cost         float64 `json:"cost"`
	SessionCount int     `json:"sessionCount"`
}

// Store AI 使用统计存储
type Store struct {
	mu       sync.RWMutex
	records  []*UsageRecord
	dataPath string
}

// NewStore 创建 AI 统计存储
func NewStore() (*Store, error) {
	dataDir := getDataDir()
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	store := &Store{
		records:  make([]*UsageRecord, 0),
		dataPath: filepath.Join(dataDir, "ai_usage.json"),
	}

	store.load()
	return store, nil
}

// Record 记录使用数据
func (s *Store) Record(projectID, sessionID, tool string, tokensInput, tokensOutput int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := &UsageRecord{
		ProjectID:    projectID,
		SessionID:    sessionID,
		Tool:         tool,
		TokensInput:  tokensInput,
		TokensOutput: tokensOutput,
		Cost:         calculateCost(tool, tokensInput, tokensOutput),
		CreatedAt:    time.Now(),
	}

	s.records = append(s.records, record)
	s.save()
}

// GetDailyUsage 获取每日用量
func (s *Store) GetDailyUsage(days int) []DailyUsage {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// TODO: 实现按日统计
	return []DailyUsage{}
}

// GetToolStats 获取工具统计
func (s *Store) GetToolStats() []ToolStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// TODO: 实现工具统计
	return []ToolStats{}
}

// GetSessions 获取会话列表
func (s *Store) GetSessions(projectID string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sessions := make(map[string]bool)
	for _, record := range s.records {
		if record.ProjectID == projectID {
			sessions[record.SessionID] = true
		}
	}

	result := make([]string, 0, len(sessions))
	for sessionID := range sessions {
		result = append(result, sessionID)
	}
	return result
}

// Clear 清空数据
func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records = make([]*UsageRecord, 0)
	s.save()
}

// load 加载数据
func (s *Store) load() {
	data, err := os.ReadFile(s.dataPath)
	if err != nil {
		return
	}

	json.Unmarshal(data, &s.records)
}

// save 保存数据
func (s *Store) save() {
	data, _ := json.MarshalIndent(s.records, "", "  ")
	os.WriteFile(s.dataPath, data, 0644)
}

func getDataDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "codux")
}

func calculateCost(tool string, input, output int64) float64 {
	// TODO: 根据工具和 token 数计算成本
	return 0
}
