package app

import (
	"os/exec"
	"runtime"
)

// OpenInVSCode 在 VSCode 中打开项目
func OpenInVSCode(path string) error {
	cmd := exec.Command("code", path)
	return cmd.Start()
}

// OpenInTerminal 在终端中打开项目
func OpenInTerminal(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", "-a", "Terminal", path)
	case "linux":
		cmd = exec.Command("gnome-terminal", "--working-directory", path)
	case "windows":
		cmd = exec.Command("cmd", "/K", "cd", "/d", path)
	default:
		cmd = exec.Command("cd", path)
	}
	return cmd.Start()
}

// OpenInFinder 在文件管理器中显示
func OpenInFinder(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", "-R", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	case "windows":
		cmd = exec.Command("explorer", path)
	default:
		cmd = exec.Command("open", path)
	}
	return cmd.Start()
}

// OpenInXcode 在 Xcode 中打开
func OpenInXcode(path string) error {
	cmd := exec.Command("open", "-a", "Xcode", path)
	return cmd.Start()
}

// OpenInIterm 在 iTerm2 中打开
func OpenInIterm(path string) error {
	cmd := exec.Command("open", "-a", "iTerm", path)
	return cmd.Start()
}
