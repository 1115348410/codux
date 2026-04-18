# Codux Windows 版本技术规格说明书

## 1. 项目概述

### 1.1 项目名称
**Codux for Windows** - AI 编程工具的原生 Windows 终端工作区

### 1.2 核心功能概要
一个专为 AI CLI 工具打造的原生 Windows 终端工作区，支持多项目、多分屏、内置 Git 面板和 AI 用量追踪。

### 1.3 项目类型
Windows 原生桌面应用程序 (WinUI 3 / Windows App SDK)

### 1.4 目标平台
- Windows 10 (version 1809, build 17763) 或更高
- Windows 11

---

## 2. 技术架构

### 2.1 技术栈

| 层级 | 技术选型 | 说明 |
|------|---------|------|
| UI 框架 | **WinUI 3** | Windows App SDK 的现代化 UI 框架 |
| 应用框架 | **Windows Application Packaging Project** | MSIX 打包 |
| 语言 | **C# / C++/WinRT** | 主要使用 C#，底层交互用 C++/WinRT |
| 终端 | **Windows Terminal / ConPTY** | 通过 ConPTY 接口集成 |
| 状态管理 | **MVVM** | CommunityToolkit.Mvvm |
| 依赖注入 | **Microsoft.Extensions.DependencyInjection** | |
| 日志 | **Serilog** | 跨平台日志 |
| 更新 | **WinGet** | Windows 程序包管理器更新 |

### 2.2 项目结构

```
Codux/
├── src/
│   ├── Codux.WinUI/                    # 主应用程序 (WinUI 3)
│   │   ├── App.xaml / App.xaml.cs
│   │   ├── MainWindow.xaml
│   │   ├── Views/                      # 页面视图
│   │   │   ├── WorkspaceView.xaml
│   │   │   ├── SettingsView.xaml
│   │   │   ├── GitPanelView.xaml
│   │   │   └── AIStatsView.xaml
│   │   ├── ViewModels/                 # 视图模型
│   │   ├── Models/                     # 数据模型
│   │   ├── Services/                  # 业务服务
│   │   ├── Platform/                   # Windows 平台特定代码
│   │   └── Resources/                 # 资源文件
│   │
│   ├── Codux.Core/                    # 核心业务逻辑（可跨平台）
│   │   ├── Session management
│   │   ├── Project management
│   │   ├── Git operations
│   │   └── AI usage tracking
│   │
│   └── Codux.Terminal/                # 终端抽象层
│       ├── ITerminalService           # 终端接口
│       └── WindowsTerminalService     # Windows Terminal 实现
│
├── SPEC.md
└── README.md
```

### 2.3 架构分层

```
┌─────────────────────────────────────────────────────────────┐
│                        WinUI 3 Views                        │
│  (WorkspaceView, SettingsView, GitPanelView, AIStatsView)  │
├─────────────────────────────────────────────────────────────┤
│                      ViewModels (MVVM)                      │
│   (WorkspaceViewModel, SettingsViewModel, GitPanelVM...)   │
├─────────────────────────────────────────────────────────────┤
│                       Core Services                         │
│  (SessionService, ProjectService, GitService, UsageService)│
├─────────────────────────────────────────────────────────────┤
│                    Platform Abstraction                     │
│         (ITerminalService, INotificationService)           │
├─────────────────────────────────────────────────────────────┤
│                     Windows Platform                        │
│    (WindowsTerminalService, ToastNotificationService)      │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. 功能模块设计

### 3.1 多项目工作区

| 功能 | 描述 | WinUI 3 实现 |
|------|------|-------------|
| 项目管理 | 添加/移除项目，切换活跃项目 | ListView + NavigationView |
| 会话保存 | 自动保存/恢复工作区状态 | ApplicationData API |
| 状态监控 | 实时显示 AI 工具运行状态 | Polling / Event hooks |

### 3.2 分屏布局

| 功能 | 描述 | WinUI 3 实现 |
|------|------|-------------|
| 水平分屏 | 左右并排显示终端 | Grid / SplitView |
| 垂直分屏 | 上下堆叠显示终端 | Grid / SplitView |
| 标签页 | 底部标签切换 | TabView |
| 拖拽调整 | 拖拽分隔线调整大小 | GridSplitter |

### 3.3 内置 Git 面板

| 功能 | 描述 | WinUI 3 实现 |
|------|------|-------------|
| 分支管理 | 切换分支、创建分支 | TreeView |
| 变更暂存 | 暂存/取消暂存文件 | CheckBox + ListView |
| 差异查看 | 显示文件变更 | 自定义 DiffViewer |
| 提交历史 | 显示提交记录 | ListView + DataGrid |
| 远程同步 | Pull/Push/Fetch | Button + ProgressBar |

### 3.4 AI 用量仪表盘

| 功能 | 描述 | WinUI 3 实现 |
|------|------|-------------|
| Token 统计 | 追踪 Token 消耗 | ProgressRing / DataGrid |
| 模型用量 | 各模型使用占比 | PieChart / BarChart |
| 每日趋势 | 每日用量趋势 | LineChart |
| 等级系统 | AI 使用等级展示 | Badge / ProgressBar |

### 3.5 终端功能

| 功能 | 描述 | WinUI 3 实现 |
|------|------|-------------|
| 终端模拟 | 执行 shell 命令 | ConPTY (Windows Pseudo Console) |
| 多终端 | 支持多个并发终端 | 每个 Session 对应一个 ConPTY |
| 主题支持 | 终端颜色主题 | Windows Terminal Settings |
| 快捷键 | 终端快捷键 | KeyDown 事件处理 |

---

## 4. 平台适配层设计

### 4.1 终端服务接口

```csharp
public interface ITerminalService
{
    Task<TerminalSession> CreateSessionAsync(string workingDirectory, CancellationToken ct);
    Task WriteAsync(Guid sessionId, string input, CancellationToken ct);
    Task<string> ReadAsync(Guid sessionId, CancellationToken ct);
    Task ResizeAsync(Guid sessionId, int width, int height, CancellationToken ct);
    Task CloseSessionAsync(Guid sessionId, CancellationToken ct);
    event EventHandler<TerminalOutputEventArgs> OutputReceived;
}
```

### 4.2 Windows 实现 (ConPTY)

```csharp
public class WindowsTerminalService : ITerminalService
{
    // 使用 Windows ConPTY API
    // Pseudo Console API (ConPTY)
}
```

### 4.3 通知服务接口

```csharp
public interface INotificationService
{
    Task ShowToastAsync(string title, string message, CancellationToken ct);
    Task ShowProgressAsync(string title, double progress, CancellationToken ct);
}
```

### 4.4 Windows 实现

```csharp
public class ToastNotificationService : INotificationService
{
    // 使用 Windows.UI.Notifications
    // ToastNotificationManager
}
```

---

## 5. 数据持久化

### 5.1 本地存储

| 数据 | 存储位置 | API |
|------|---------|-----|
| 应用设置 | %LOCALAPPDATA%\Codux\Settings | ApplicationDataContainer |
| 项目状态 | %LOCALAPPDATA%\Codux\Workspaces | ApplicationDataContainer |
| 日志文件 | %LOCALAPPDATA%\Codux\Logs | File API |
| AI 用量数据 | %LOCALAPPDATA%\Codux\Usage | SQLite / ApplicationData |

### 5.2 数据模型

```csharp
public class WorkspaceState
{
    public Guid Id { get; set; }
    public string Name { get; set; }
    public string Path { get; set; }
    public List<PaneState> Panes { get; set; }
    public Guid? ActivePaneId { get; set; }
    public DateTime LastAccessed { get; set; }
}

public class PaneState
{
    public Guid Id { get; set; }
    public Orientation Orientation { get; set; }
    public double SplitRatio { get; set; }
    public List<TerminalSessionState> Sessions { get; set; }
}

public class TerminalSessionState
{
    public Guid Id { get; set; }
    public Guid? TerminalPaneId { get; set; }
    public string WorkingDirectory { get; set; }
    public int Width { get; set; }
    public int Height { get; set; }
}
```

---

## 6. 依赖项清单

### 6.1 NuGet 包

| 包名 | 版本 | 用途 |
|------|-----|------|
| **Microsoft.WindowsAppSDK** | 1.5+ | WinUI 3 框架 |
| **CommunityToolkit.Mvvm** | 8.2+ | MVVM 工具 |
| **Microsoft.Extensions.DependencyInjection** | 8.0+ | 依赖注入 |
| **Serilog** + **Serilog.Sinks.File** | 3.0+ | 日志 |
| **LibGit2Sharp** | 0.30+ | Git 操作 |
| **Microsoft.Toolkit.Uwp.Notifications** | 7.0+ | Toast 通知 |
| **Newtonsoft.Json** | 13.0+ | JSON 序列化 |

### 6.2 Windows SDK

| 组件 | 用途 |
|------|------|
| **ConPTY API** | 终端仿真 |
| **Windows.UI.Notifications** | 通知 |
| **Windows.Storage** | 文件存储 |

---

## 7. 构建与发布

### 7.1 构建配置

- **框架**: net8.0-windows10.0.19041.0
- **目标**: Windows 10 1809+ (build 17763)
- **输出**: MSIX 包 + 可执行文件

### 7.2 发布方式

1. **MSIX 包** (Microsoft Store)
2. **独立 exe** ( sideload 安装)
3. **WinGet** (未来支持)

---

## 8. 已知限制与风险

| 问题 | 影响 | 缓解方案 |
|------|-----|---------|
| ConPTY 复杂性 | 终端功能开发难度高 | 使用 Windows Terminal 作为底层 |
| WinUI 3 成熟度 | 某些控件可能不完善 | 回退到 UWP 控件 |
| LibGit2Sharp 兼容性 | Git 操作可能有问题 | 测试主流 Git 操作 |
| 性能 | 多个 ConPTY 可能占用较高 | 资源监控和优化 |

---

## 9. 后续工作

1. 创建 WinUI 3 项目骨架
2. 实现基础 UI 布局（NavigationView + 分屏）
3. 集成 ConPTY 终端
4. 实现 Git 服务
5. 实现 AI 用量追踪
6. 完整功能和 UI 打磨
7. 打包和发布
