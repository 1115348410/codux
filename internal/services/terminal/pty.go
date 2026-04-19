package terminal

import (
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
)

// PTY 终端实例
type PTY struct {
	mu       sync.Mutex
	cmd      *exec.Cmd
	ptmx     *os.File
	isActive bool
}

// NewPTY 创建新的 PTY 实例
func NewPTY() *PTY {
	return &PTY{}
}

// Start 启动 PTY 进程
func (p *PTY) Start(command string, dir string, env []string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isActive {
		return nil
	}

	// 解析命令
	var cmd *exec.Cmd
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "bash"
	}

	if command != "" {
		cmd = exec.Command(shell, "-c", command)
	} else {
		cmd = exec.Command(shell, "-l")
	}

	if dir != "" {
		cmd.Dir = dir
	}

	if len(env) > 0 {
		cmd.Env = env
	} else {
		cmd.Env = os.Environ()
	}

	// 启动 PTY
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	p.cmd = cmd
	p.ptmx = ptmx
	p.isActive = true

	return nil
}

// Write 写入数据到 PTY
func (p *PTY) Write(data []byte) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isActive || p.ptmx == nil {
		return 0, nil
	}

	return p.ptmx.Write(data)
}

// Read 从 PTY 读取数据
func (p *PTY) Read(buffer []byte) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isActive || p.ptmx == nil {
		return 0, io.EOF
	}

	return p.ptmx.Read(buffer)
}

// Resize 调整 PTY 大小
func (p *PTY) Resize(height, width int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isActive || p.ptmx == nil {
		return nil
	}

	return pty.Setsize(p.ptmx, &pty.Winsize{
		Rows: uint16(height),
		Cols: uint16(width),
	})
}

// Close 关闭 PTY
func (p *PTY) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isActive {
		return nil
	}

	p.isActive = false

	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Kill()
	}

	if p.ptmx != nil {
		p.ptmx.Close()
	}

	return nil
}

// IsActive 检查 PTY 是否活跃
func (p *PTY) IsActive() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.isActive
}

// SetWinSize 设置窗口大小
func (p *PTY) SetWinSize(rows, cols uint16) error {
	return pty.Setsize(p.ptmx, &pty.Winsize{
		Rows: rows,
		Cols: cols,
	})
}
