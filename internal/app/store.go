package app

import (
	"sync"

	"github.com/duxweb/codux/internal/models"
	"github.com/duxweb/codux/internal/services/persist"
)

// Store 应用状态管理
type Store struct {
	mu sync.RWMutex

	projects          []*models.Project
	workspaces        map[string]*models.ProjectWorkspace
	selectedProjectID *models.UUID
	settings          *models.AppSettings

	persistence *persist.PersistenceService
}

// NewStore 创建应用 Store
func NewStore() (*Store, error) {
	ps, err := persist.NewPersistenceService()
	if err != nil {
		return nil, err
	}

	store := &Store{
		workspaces:  make(map[string]*models.ProjectWorkspace),
		persistence: ps,
	}

	if err := store.loadFromPersistence(); err != nil {
		return nil, err
	}

	return store, nil
}

// loadFromPersistence 从持久化层加载数据
func (s *Store) loadFromPersistence() error {
	// 加载项目
	projects, err := s.persistence.GetAllProjects()
	if err != nil {
		return err
	}
	s.projects = projects

	// 加载工作区
	for _, project := range projects {
		workspace, err := s.persistence.GetWorkspace(project.ID.String())
		if err != nil {
			continue
		}
		s.workspaces[project.ID.String()] = workspace
	}

	// 加载设置
	themeStr, _ := s.persistence.GetSetting("theme")
	if themeStr != "" {
		if s.settings == nil {
			s.settings = models.DefaultAppSettings()
		}
		s.settings.Theme = themeStr
	}

	// 如果没有加载到设置，使用默认值
	if s.settings == nil {
		s.settings = models.DefaultAppSettings()
	}

	// 如果有项目，选择第一个
	if len(s.projects) > 0 {
		s.selectedProjectID = &s.projects[0].ID
	}

	return nil
}

// Projects 获取所有项目
func (s *Store) Projects() []*models.Project {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.projects
}

// AddProject 添加项目
func (s *Store) AddProject(project *models.Project) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.persistence.SaveProject(project); err != nil {
		return err
	}

	s.projects = append(s.projects, project)

	// 创建工作区
	workspace := models.NewProjectWorkspace(project.ID)
	if err := s.persistence.SaveWorkspace(workspace); err != nil {
		return err
	}
	s.workspaces[project.ID.String()] = workspace

	// 如果是第一个项目，自动选中
	if len(s.projects) == 1 {
		s.selectedProjectID = &project.ID
	}

	return nil
}

// UpdateProject 更新项目
func (s *Store) UpdateProject(project *models.Project) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.persistence.SaveProject(project); err != nil {
		return err
	}

	for i, p := range s.projects {
		if p.ID.String() == project.ID.String() {
			s.projects[i] = project
			break
		}
	}

	return nil
}

// DeleteProject 删除项目
func (s *Store) DeleteProject(id models.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.persistence.DeleteProject(id.String()); err != nil {
		return err
	}

	// 从列表中移除
	newProjects := make([]*models.Project, 0, len(s.projects)-1)
	for _, p := range s.projects {
		if p.ID.String() != id.String() {
			newProjects = append(newProjects, p)
		}
	}
	s.projects = newProjects

	// 从工作区中移除
	delete(s.workspaces, id.String())

	// 如果删除的是当前选中的项目，选择第一个项目
	if s.selectedProjectID != nil && s.selectedProjectID.String() == id.String() {
		if len(s.projects) > 0 {
			s.selectedProjectID = &s.projects[0].ID
		} else {
			s.selectedProjectID = nil
		}
	}

	return nil
}

// SelectedProject 获取当前选中的项目
func (s *Store) SelectedProject() *models.Project {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.selectedProjectID == nil {
		return nil
	}

	for _, p := range s.projects {
		if p.ID.String() == s.selectedProjectID.String() {
			return p
		}
	}
	return nil
}

// SelectProject 选择项目
func (s *Store) SelectProject(id models.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.selectedProjectID = &id
}

// GetWorkspace 获取项目工作区
func (s *Store) GetWorkspace(projectID models.UUID) *models.ProjectWorkspace {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.workspaces[projectID.String()]
}

// SaveWorkspace 保存工作区
func (s *Store) SaveWorkspace(workspace *models.ProjectWorkspace) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.workspaces[workspace.ProjectID.String()] = workspace
	return s.persistence.SaveWorkspace(workspace)
}

// Settings 获取应用设置
func (s *Store) Settings() *models.AppSettings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.settings
}

// UpdateSettings 更新应用设置
func (s *Store) UpdateSettings(settings *models.AppSettings) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.settings = settings
	return s.persistence.SaveSetting("theme", settings.Theme)
}

// Close 关闭 Store
func (s *Store) Close() error {
	return s.persistence.Close()
}
