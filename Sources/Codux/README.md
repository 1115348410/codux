# Codux - Windows Tauri Version

AI 编程工具的原生 Windows 终端工作区 (Tauri/Rust 版本)

## 技术栈

- **后端**: Rust + Tauri 2.0
- **前端**: HTML/CSS/JavaScript (原生 Web)
- **终端**: portable-pty (支持 PTY)
- **打包**: Windows exe

## 项目结构

```
Sources/Codux/
├── src/                    # 前端 (HTML/CSS/JS)
├── src-tauri/
│   ├── src/
│   │   ├── main.rs         # 应用入口
│   │   ├── terminal.rs     # 终端服务
│   │   └── project.rs      # 项目管理
│   ├── Cargo.toml         # Rust 依赖
│   ├── tauri.conf.json    # Tauri 配置
│   └── build.rs           # 构建脚本
└── README.md
```

## 构建要求

- Rust 1.70+
- Node.js 18+
- Windows 10/11

## 构建步骤

```bash
# 1. 安装 Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# 2. 安装 Node.js 依赖
cd Sources/Codux
npm install

# 3. 构建
npm run tauri build
```

## 功能

- [x] 多项目工作区
- [x] 终端模拟 (PTY)
- [ ] 分屏布局
- [ ] 内置 Git 面板
- [ ] AI 用量追踪

## 许可证

MIT
