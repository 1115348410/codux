# Codux for Windows

一个为 AI 编程工具打造的原生 Windows 终端工作区。

## 项目状态

**此项目正在开发中** - 这是 macOS 版 Codux 的 Windows (WPF) 移植版本。

## 技术栈

- **UI 框架**: WPF
- **语言**: C# (.NET 8)
- **架构**: MVVM
- **终端**: ConPTY (Windows Pseudo Console)
- **Git**: LibGit2Sharp
- **日志**: Serilog

## 构建要求

- Windows 10 version 1809 (build 17763) 或更高
- Windows 11
- .NET 8 SDK
- Visual Studio 2022 (推荐) 或 VS Code

## 项目结构

```
Sources/DmuxWorkspace.WinUI/
├── App.xaml / App.xaml.cs     # 应用入口
├── MainWindow.xaml           # 主窗口
├── Views/                    # 页面视图
├── ViewModels/               # 视图模型 (MVVM)
├── Services/                 # 业务服务
│   ├── ITerminalService.cs   # 终端接口
│   ├── WindowsTerminalService.cs
│   ├── GitService.cs
│   └── ...
├── Models/                   # 数据模型
├── Converters/               # 值转换器
└── Resources/                # 资源文件
```

## 核心功能 (规划中)

- [ ] 多项目工作区
- [ ] 灵活的分屏布局
- [ ] 内置 Git 面板
- [ ] AI 用量仪表盘
- [ ] 终端模拟 (ConPTY)

## 构建

```bash
dotnet restore
dotnet build
```

## 运行

```bash
dotnet run
```

## 许可证

参见 [LICENSE](../LICENSE)
