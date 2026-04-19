# Codux

A native terminal workspace for AI coding tools, built with Go and Fyne.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go](https://img.shields.io/badge/Go-1.21+-blue.svg)
![Fyne](https://img.shields.io/badge/Fyne-2.5+-blue.svg)

## Features

- **Multi-Project Workspace** - Organize all your projects in one place
- **Terminal Integration** - Built-in PTY terminal with split panes
- **Git Panel** - Built-in Git GUI for branch management, commits, and diffs
- **AI Usage Tracking** - Monitor token consumption for AI coding tools
- **Cross-Platform** - Linux, macOS, and Windows support
- **Native Performance** - No Electron, fast and lightweight

## Screenshots

![Codux Screenshot](docs/screenshot.png)

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

### System Requirements

**Linux:**
```bash
sudo apt-get install libgl1-mesa-dev xorg-dev
```

**macOS:**
```bash
xcode-select --install
```

**Windows:**
- Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or MinGW

## Keyboard Shortcuts

| Action | Shortcut |
|--------|----------|
| New Split | `Ctrl+T` |
| New Tab | `Ctrl+D` |
| Toggle Git Panel | `Ctrl+G` |
| Toggle AI Panel | `Ctrl+Y` |
| Close Split | `Ctrl+W` |
| Switch Project | `Alt+1-9` |

## Development

```bash
# Run in development mode
make run

# Build with debug info
make build-debug

# Run tests
make test

# Format code
make fmt

# Build for all platforms
make build-all
```

## Project Structure

```
codux/
├── main.go                 # Application entry point
├── internal/
│   ├── app/               # Application state management
│   ├── models/            # Data models
│   ├── services/          # Business services
│   │   ├── git/          # Git operations
│   │   ├── terminal/     # PTY terminal
│   │   ├── ai/           # AI usage tracking
│   │   └── persist/      # SQLite persistence
│   └── ui/               # User interface
│       ├── menu/         # Application menu
│       ├── settings/     # Settings window
│       ├── theme/        # Theme system
│       ├── views/        # View components
│       └── widgets/      # Custom widgets
└── cmd/                  # CLI tools
```

## Technology Stack

- **Language**: Go 1.21+
- **GUI Framework**: Fyne 2.5+
- **Terminal**: creack/pty
- **Database**: modernc.org/sqlite
- **State Management**: Custom Store pattern

## Roadmap

- [x] Project management
- [x] Theme system
- [x] Git integration
- [x] Terminal PTY
- [ ] Terminal emulation
- [ ] Split panes
- [ ] AI usage tracking
- [ ] Auto-update
- [ ] i18n support

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

Inspired by the original Swift/macOS version of Codux.
