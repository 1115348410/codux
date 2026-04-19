package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// Service Git 服务
type Service struct {
	mu sync.Mutex
}

// NewService 创建 Git 服务
func NewService() *Service {
	return &Service{}
}

// RepositoryState 仓库状态
type RepositoryState struct {
	IsGitRepo      bool
	CurrentBranch  string
	RemoteBranch   string
	AheadBehind    string
	ModifiedFiles  []FileStatus
	StagedFiles    []FileStatus
	UntrackedFiles []FileStatus
	Branches       []string
	Remotes        []Remote
}

// FileStatus 文件状态
type FileStatus struct {
	Path   string
	Status string
}

// Remote 远程仓库
type Remote struct {
	Name string
	URL  string
}

// GetRepositoryState 获取仓库状态
func (s *Service) GetRepositoryState(repoPath string) (*RepositoryState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isGitRepository(repoPath) {
		return &RepositoryState{IsGitRepo: false}, nil
	}

	state := &RepositoryState{
		IsGitRepo: true,
	}

	// 获取当前分支
	branch, err := s.runGitCommand(repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	if err == nil {
		state.CurrentBranch = strings.TrimSpace(branch)
	}

	// 获取上游分支
	upstream, err := s.runGitCommand(repoPath, "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if err == nil {
		state.RemoteBranch = strings.TrimSpace(upstream)
	}

	// 获取 ahead/behind 状态
	if state.RemoteBranch != "" {
		aheadBehind, err := s.runGitCommand(repoPath, "rev-list", "--left-right", "--count", "HEAD..."+state.RemoteBranch)
		if err == nil {
			state.AheadBehind = strings.TrimSpace(aheadBehind)
		}
	}

	// 获取文件状态
	status, err := s.runGitCommand(repoPath, "status", "--porcelain")
	if err == nil {
		s.parseStatus(status, state)
	}

	// 获取分支列表
	branches, err := s.runGitCommand(repoPath, "branch")
	if err == nil {
		for _, line := range strings.Split(branches, "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				state.Branches = append(state.Branches, strings.TrimPrefix(line, "* "))
			}
		}
	}

	// 获取远程列表
	remotes, err := s.runGitCommand(repoPath, "remote", "-v")
	if err == nil {
		s.parseRemotes(remotes, state)
	}

	return state, nil
}

func (s *Service) isGitRepository(path string) bool {
	gitDir := filepath.Join(path, ".git")
	info, err := os.Stat(gitDir)
	return err == nil && info.IsDir()
}

func (s *Service) parseStatus(output string, state *RepositoryState) {
	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			continue
		}

		status := line[:2]
		path := strings.TrimSpace(line[3:])

		file := FileStatus{Path: path, Status: status}

		switch {
		case strings.HasPrefix(status, " "):
			state.StagedFiles = append(state.StagedFiles, file)
		case strings.HasPrefix(status, "??"):
			state.UntrackedFiles = append(state.UntrackedFiles, file)
		default:
			state.ModifiedFiles = append(state.ModifiedFiles, file)
		}
	}
}

func (s *Service) parseRemotes(output string, state *RepositoryState) {
	seen := make(map[string]bool)
	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			name := parts[1]
			if !seen[name] {
				seen[name] = true
				state.Remotes = append(state.Remotes, Remote{Name: name, URL: parts[0]})
			}
		}
	}
}

func (s *Service) runGitCommand(repoPath string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	output, err := cmd.Output()
	return string(output), err
}

// StageFile 暂存文件
func (s *Service) StageFile(repoPath, file string) error {
	_, err := s.runGitCommand(repoPath, "add", file)
	return err
}

// UnstageFile 取消暂存
func (s *Service) UnstageFile(repoPath, file string) error {
	_, err := s.runGitCommand(repoPath, "restore", "--staged", file)
	return err
}

// DiscardFile 丢弃更改
func (s *Service) DiscardFile(repoPath, file string) error {
	_, err := s.runGitCommand(repoPath, "restore", file)
	return err
}

// Commit 提交
func (s *Service) Commit(repoPath, message string) error {
	_, err := s.runGitCommand(repoPath, "commit", "-m", message)
	return err
}

// Pull 拉取
func (s *Service) Pull(repoPath string) error {
	_, err := s.runGitCommand(repoPath, "pull")
	return err
}

// Push 推送
func (s *Service) Push(repoPath string) error {
	_, err := s.runGitCommand(repoPath, "push")
	return err
}

// Fetch 获取
func (s *Service) Fetch(repoPath string) error {
	_, err := s.runGitCommand(repoPath, "fetch")
	return err
}

// CreateBranch 创建分支
func (s *Service) CreateBranch(repoPath, branch string) error {
	_, err := s.runGitCommand(repoPath, "checkout", "-b", branch)
	return err
}

// CheckoutBranch 切换分支
func (s *Service) CheckoutBranch(repoPath, branch string) error {
	_, err := s.runGitCommand(repoPath, "checkout", branch)
	return err
}

// GetDiff 获取差异
func (s *Service) GetDiff(repoPath, file string) (string, error) {
	return s.runGitCommand(repoPath, "diff", "--", file)
}

// GetStagedDiff 获取暂存区差异
func (s *Service) GetStagedDiff(repoPath, file string) (string, error) {
	return s.runGitCommand(repoPath, "diff", "--staged", "--", file)
}
