# Codux

A native terminal workspace for AI coding tools, built with Go and Fyne.

## Features

- **Multi-Project Workspace**: Organize all your projects in one place
- **Terminal Integration**: Built-in terminal emulator with split panes
- **Git Panel**: Built-in Git GUI for branch management, commits, and diffs
- **AI Usage Tracking**: Monitor token consumption for Claude Code, Codex, and Gemini CLI

## Requirements

- Go 1.21+
- Linux: X11, OpenGL development libraries
- macOS: Xcode command line tools
- Windows: GCC (TDM-GCC or MinGW)

## Installation

### Build from Source

```bash
# Clone the repository
git clone https://github.com/duxweb/codux.git
cd codux

# Build
go build -o codux .

# Run
./codux
```

### Install Dependencies

**Linux (Debian/Ubuntu)**:
```bash
sudo apt-get install libgl1-mesa-dev xorg-dev
```

**Linux (Fedora)**:
```bash
sudo dnf install mesa-libGL-devel libXext-devel libXrender-devel
```

**macOS**:
```bash
xcode-select --install
```

## Keyboard Shortcuts

| Action | Shortcut |
|--------|----------|
| New Split | `Ctrl+T` |
| New Tab | `Ctrl+D` |
| Toggle Git Panel | `Ctrl+G` |
| Toggle AI Panel | `Ctrl+Y` |

## Development

```bash
# Run in development mode
go run .

# Build with debug info
go build -tags debug -o codux .

# Run tests
go test ./...
```

## Project Structure

```
.
├── main.go                 # Application entry point
├── internal/
│   ├── app/               # Application state management
│   ├── models/            # Data models
│   ├── services/          # Business logic services
│   │   ├── git/          # Git operations
│   │   ├── ai/           # AI runtime probing
│   │   ├── terminal/     # PTY management
│   │   └── persist/      # Data persistence
│   └── ui/               # User interface
│       ├── theme/        # Theme system
│       ├── widgets/      # Custom widgets
│       └── views/        # View components
└── resources/            # Static assets
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

Inspired by the original Swift/macOS version of Codux.
