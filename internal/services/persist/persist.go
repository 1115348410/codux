package persist

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/duxweb/codux/internal/models"
	_ "modernc.org/sqlite"
)

// PersistenceService 持久化服务
type PersistenceService struct {
	db *sql.DB
}

// NewPersistenceService 创建持久化服务
func NewPersistenceService() (*PersistenceService, error) {
	dbPath := getDatabasePath()

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	ps := &PersistenceService{db: db}
	if err := ps.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return ps, nil
}

// initSchema 初始化数据库表结构
func (ps *PersistenceService) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS projects (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		path TEXT NOT NULL,
		shell TEXT DEFAULT 'bash',
		default_command TEXT DEFAULT '',
		badge_text TEXT DEFAULT '',
		badge_symbol TEXT DEFAULT '',
		badge_color_hex TEXT DEFAULT '#007AFF',
		git_default_push_remote TEXT DEFAULT 'origin'
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		project_id TEXT NOT NULL,
		title TEXT DEFAULT '',
		command TEXT DEFAULT '',
		working_dir TEXT DEFAULT '',
		is_created BOOLEAN DEFAULT FALSE,
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS workspaces (
		project_id TEXT PRIMARY KEY,
		top_session_ids TEXT DEFAULT '[]',
		bottom_tab_session_ids TEXT DEFAULT '[]',
		selected_session_id TEXT DEFAULT '',
		selected_bottom_tab_id TEXT DEFAULT '',
		top_pane_ratios TEXT DEFAULT '[]',
		bottom_pane_height REAL DEFAULT 200.0,
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS ai_usage (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id TEXT NOT NULL,
		session_id TEXT NOT NULL,
		tool TEXT NOT NULL,
		tokens_input INTEGER DEFAULT 0,
		tokens_output INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_sessions_project_id ON sessions(project_id);
	CREATE INDEX IF NOT EXISTS idx_ai_usage_project_id ON ai_usage(project_id);
	CREATE INDEX IF NOT EXISTS idx_ai_usage_created_at ON ai_usage(created_at);
	`

	_, err := ps.db.Exec(schema)
	return err
}

// SaveProject 保存项目
func (ps *PersistenceService) SaveProject(project *models.Project) error {
	query := `
	INSERT OR REPLACE INTO projects (id, name, path, shell, default_command, badge_text, badge_symbol, badge_color_hex, git_default_push_remote)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := ps.db.Exec(query,
		project.ID.String(),
		project.Name,
		project.Path,
		project.Shell,
		project.DefaultCommand,
		project.BadgeText,
		project.BadgeSymbol,
		project.BadgeColorHex,
		project.GitDefaultPushRemote,
	)
	return err
}

// GetProject 获取项目
func (ps *PersistenceService) GetProject(id string) (*models.Project, error) {
	query := `SELECT id, name, path, shell, default_command, badge_text, badge_symbol, badge_color_hex, git_default_push_remote FROM projects WHERE id = ?`

	var project models.Project
	var idStr, badgeSymbol sql.NullString
	err := ps.db.QueryRow(query, id).Scan(
		&idStr,
		&project.Name,
		&project.Path,
		&project.Shell,
		&project.DefaultCommand,
		&project.BadgeText,
		&badgeSymbol,
		&project.BadgeColorHex,
		&project.GitDefaultPushRemote,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	project.ID = models.UUID{}
	project.ID.UnmarshalText([]byte(idStr.String))
	if badgeSymbol.Valid {
		project.BadgeSymbol = badgeSymbol.String
	}

	return &project, nil
}

// GetAllProjects 获取所有项目
func (ps *PersistenceService) GetAllProjects() ([]*models.Project, error) {
	query := `SELECT id, name, path, shell, default_command, badge_text, badge_symbol, badge_color_hex, git_default_push_remote FROM projects`

	rows, err := ps.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var project models.Project
		var idStr, badgeSymbol sql.NullString
		if err := rows.Scan(
			&idStr,
			&project.Name,
			&project.Path,
			&project.Shell,
			&project.DefaultCommand,
			&project.BadgeText,
			&badgeSymbol,
			&project.BadgeColorHex,
			&project.GitDefaultPushRemote,
		); err != nil {
			return nil, err
		}

		project.ID = models.UUID{}
		project.ID.UnmarshalText([]byte(idStr.String))
		if badgeSymbol.Valid {
			project.BadgeSymbol = badgeSymbol.String
		}

		projects = append(projects, &project)
	}

	return projects, rows.Err()
}

// DeleteProject 删除项目
func (ps *PersistenceService) DeleteProject(id string) error {
	_, err := ps.db.Exec("DELETE FROM projects WHERE id = ?", id)
	return err
}

// SaveWorkspace 保存工作区
func (ps *PersistenceService) SaveWorkspace(workspace *models.ProjectWorkspace) error {
	topIDs, _ := json.Marshal(uuidsToStrings(workspace.TopSessionIDs))
	bottomIDs, _ := json.Marshal(uuidsToStrings(workspace.BottomTabSessionIDs))
	ratios, _ := json.Marshal(workspace.TopPaneRatios)

	query := `
	INSERT OR REPLACE INTO workspaces (project_id, top_session_ids, bottom_tab_session_ids, selected_session_id, selected_bottom_tab_id, top_pane_ratios, bottom_pane_height)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := ps.db.Exec(query,
		workspace.ProjectID.String(),
		string(topIDs),
		string(bottomIDs),
		workspace.SelectedSessionID.String(),
		workspace.SelectedBottomTabID.String(),
		string(ratios),
		workspace.BottomPaneHeight,
	)
	return err
}

// GetWorkspace 获取工作区
func (ps *PersistenceService) GetWorkspace(projectID string) (*models.ProjectWorkspace, error) {
	query := `SELECT top_session_ids, bottom_tab_session_ids, selected_session_id, selected_bottom_tab_id, top_pane_ratios, bottom_pane_height FROM workspaces WHERE project_id = ?`

	var topIDsStr, bottomIDsStr, selectedIDStr, selectedBottomIDStr, ratiosStr string
	var bottomHeight float64

	err := ps.db.QueryRow(query, projectID).Scan(
		&topIDsStr,
		&bottomIDsStr,
		&selectedIDStr,
		&selectedBottomIDStr,
		&ratiosStr,
		&bottomHeight,
	)

	if err == sql.ErrNoRows {
		return models.NewProjectWorkspace(models.UUID{}), nil
	}
	if err != nil {
		return nil, err
	}

	workspace := models.NewProjectWorkspace(models.UUID{})
	workspace.ProjectID.UnmarshalText([]byte(projectID))

	json.Unmarshal([]byte(topIDsStr), &workspace.TopSessionIDs)
	json.Unmarshal([]byte(bottomIDsStr), &workspace.BottomTabSessionIDs)
	json.Unmarshal([]byte(ratiosStr), &workspace.TopPaneRatios)
	workspace.SelectedSessionID.UnmarshalText([]byte(selectedIDStr))
	workspace.SelectedBottomTabID.UnmarshalText([]byte(selectedBottomIDStr))
	workspace.BottomPaneHeight = bottomHeight

	return workspace, nil
}

// SaveSetting 保存设置
func (ps *PersistenceService) SaveSetting(key, value string) error {
	_, err := ps.db.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)", key, value)
	return err
}

// GetSetting 获取设置
func (ps *PersistenceService) GetSetting(key string) (string, error) {
	var value string
	err := ps.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

// SaveAISessionUsage 保存 AI 会话使用记录
func (ps *PersistenceService) SaveAISessionUsage(projectID, sessionID, tool string, tokensInput, tokensOutput int) error {
	query := `INSERT INTO ai_usage (project_id, session_id, tool, tokens_input, tokens_output) VALUES (?, ?, ?, ?, ?)`
	_, err := ps.db.Exec(query, projectID, sessionID, tool, tokensInput, tokensOutput)
	return err
}

// Close 关闭数据库连接
func (ps *PersistenceService) Close() error {
	return ps.db.Close()
}

// 辅助函数
func uuidsToStrings(uuids []models.UUID) []string {
	result := make([]string, len(uuids))
	for i, u := range uuids {
		result[i] = u.String()
	}
	return result
}

func getDatabasePath() string {
	configDir := getConfigDir()
	return filepath.Join(configDir, "codux.db")
}

func getConfigDir() string {
	// Linux
	if configHome := os.Getenv("XDG_CONFIG_HOME"); configHome != "" {
		dir := filepath.Join(configHome, "codux")
		os.MkdirAll(dir, 0755)
		return dir
	}

	// Fallback to home directory
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "codux")
}
