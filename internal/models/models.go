package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Project 项目模型
type Project struct {
	ID                   UUID   `json:"id"`
	Name                 string `json:"name"`
	Path                 string `json:"path"`
	Shell                string `json:"shell"`
	DefaultCommand       string `json:"defaultCommand"`
	BadgeText            string `json:"badgeText"`
	BadgeSymbol          string `json:"badgeSymbol"`
	BadgeColorHex        string `json:"badgeColorHex"`
	GitDefaultPushRemote string `json:"gitDefaultPushRemote"`
}

// UUID 包装类型
type UUID struct {
	uuid.UUID
}

// String 实现 Stringer 接口
func (u UUID) String() string {
	return u.UUID.String()
}

// MarshalText 实现 encoding.TextMarshaler
func (u UUID) MarshalText() ([]byte, error) {
	return u.UUID.MarshalText()
}

// UnmarshalText 实现 encoding.TextUnmarshaler
func (u *UUID) UnmarshalText(text []byte) error {
	return u.UUID.UnmarshalText(text)
}

// MarshalJSON 实现 json.Marshaler
func (u UUID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + u.String() + `"`), nil
}

// UnmarshalJSON 实现 json.Unmarshaler
func (u *UUID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	return u.UnmarshalText([]byte(s))
}

// UUID 包装类型
type UUID struct {
	uuid.UUID
}

// NewProject 创建新项目
func NewProject(name, path string) *Project {
	return &Project{
		ID:             UUID{UUID: uuid.New()},
		Name:           name,
		Path:           path,
		Shell:          "bash",
		DefaultCommand: "",
		BadgeColorHex:  "#007AFF",
	}
}

// TerminalSession 终端会话模型
type TerminalSession struct {
	ID         UUID   `json:"id"`
	ProjectID  UUID   `json:"projectID"`
	Title      string `json:"title"`
	Command    string `json:"command"`
	WorkingDir string `json:"workingDir"`
	IsCreated  bool   `json:"isCreated"`
}

// NewTerminalSession 创建新会话
func NewTerminalSession(projectID UUID, command string) *TerminalSession {
	return &TerminalSession{
		ID:        UUID{UUID: uuid.New()},
		ProjectID: projectID,
		Command:   command,
		IsCreated: false,
	}
}

// ProjectWorkspace 项目工作区模型
type ProjectWorkspace struct {
	ProjectID           UUID               `json:"projectID"`
	Sessions            []*TerminalSession `json:"sessions"`
	TopSessionIDs       []UUID             `json:"topSessionIDs"`
	BottomTabSessionIDs []UUID             `json:"bottomTabSessionIDs"`
	SelectedSessionID   UUID               `json:"selectedSessionID"`
	SelectedBottomTabID UUID               `json:"selectedBottomTabID"`
	TopPaneRatios       []float64          `json:"topPaneRatios"`
	BottomPaneHeight    float64            `json:"bottomPaneHeight"`
}

// NewProjectWorkspace 创建工作区
func NewProjectWorkspace(projectID UUID) *ProjectWorkspace {
	return &ProjectWorkspace{
		ProjectID:           projectID,
		Sessions:            make([]*TerminalSession, 0),
		TopSessionIDs:       make([]UUID, 0),
		BottomTabSessionIDs: make([]UUID, 0),
		TopPaneRatios:       make([]float64, 0),
	}
}

// AppSettings 应用设置
type AppSettings struct {
	Theme                   string  `json:"theme"`
	TerminalFontSize        float64 `json:"terminalFontSize"`
	TerminalLineHeight      float64 `json:"terminalLineHeight"`
	TerminalFontFamily      string  `json:"terminalFontFamily"`
	GitAutoRefreshInterval  int     `json:"gitAutoRefreshInterval"`
	AIAutoRefreshInterval   int     `json:"aiAutoRefreshInterval"`
	AIBackgroundRefresh     bool    `json:"aiBackgroundRefresh"`
	ShowsPerformanceMonitor bool    `json:"showsPerformanceMonitor"`
}

// DefaultAppSettings 默认设置
func DefaultAppSettings() *AppSettings {
	return &AppSettings{
		Theme:                  "system",
		TerminalFontSize:       14,
		TerminalLineHeight:     1.2,
		TerminalFontFamily:     "Monospace",
		GitAutoRefreshInterval: 30,
		AIAutoRefreshInterval:  5,
		AIBackgroundRefresh:    true,
	}
}

// AppSnapshot 应用状态快照
type AppSnapshot struct {
	Projects          []*Project          `json:"projects"`
	Workspaces        []*ProjectWorkspace `json:"workspaces"`
	SelectedProjectID *UUID               `json:"selectedProjectID"`
	AppSettings       *AppSettings        `json:"appSettings"`
}
