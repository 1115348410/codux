# Codux 开发文档

## 项目结构

```
codux/
├── main.go                     # 应用入口
├── go.mod                      # Go 模块定义
├── go.sum                      # 依赖版本锁定
├── Makefile                    # 构建脚本
├── DEVELOPMENT.md              # 开发文档（本文件）
├── TASKS.md                    # 任务列表
├── internal/                   # 内部包
│   ├── app/                    # 应用状态管理
│   │   ├── store.go            # 应用 Store
│   │   └── terminal.go         # 终端服务
│   ├── models/                 # 数据模型
│   │   └── models.go           # 核心模型定义
│   ├── resources/              # 静态资源
│   │   └── icons/              # 图标资源
│   ├── services/               # 业务服务
│   │   ├── git/                # Git 服务
│   │   ├── terminal/           # PTY 终端服务
│   │   └── persist/            # 数据持久化
│   └── ui/                     # 用户界面
│       ├── menu/               # 应用菜单
│       ├── settings/           # 设置窗口
│       ├── theme/              # 主题系统
│       ├── views/              # 视图组件
│       │   ├── main_view.go    # 主视图
│       │   ├── sidebar_view.go # 侧边栏
│       │   ├── workspace_view.go # 工作区
│       │   └── right_panel_view.go # 右面板
│       └── widgets/            # 自定义组件
│           ├── terminal.go     # 终端组件
│           └── widgets.go      # 基础组件
└── cmd/                        # 命令行工具
```

## 架构设计

### MVC 模式

- **Model**: `internal/models/` - 数据模型和状态
- **View**: `internal/ui/views/` - UI 视图组件
- **Controller**: `internal/app/` - 状态管理和业务逻辑

### 数据流

```
用户操作 → View → Store (Controller) → Service → Model → 持久化
                                    ↓
                                View 刷新
```

### 核心组件

#### 1. Store (internal/app/store.go)

应用状态管理中心，负责：
- 项目管理（增删改查）
- 工作区管理
- 设置管理
- 数据持久化

#### 2. Services

- **Git Service**: Git 仓库操作
- **Terminal Service**: PTY 终端管理
- **Persistence Service**: SQLite 数据存储

#### 3. UI Components

- **MainView**: 三栏布局容器
- **SidebarView**: 项目列表侧边栏
- **WorkspaceView**: 终端工作区
- **RightPanelView**: Git/AI 面板

## 开发指南

### 添加新功能

1. 在 `internal/models/` 定义数据模型
2. 在 `internal/services/` 实现业务逻辑
3. 在 `internal/ui/views/` 创建视图组件
4. 在 `Store` 中注册状态管理

### 代码规范

- 使用 `go fmt` 格式化代码
- 使用 `go vet` 检查代码
- 遵循 Go 命名规范
- 导出函数/类型首字母大写

### 测试

```bash
# 运行所有测试
make test

# 运行带覆盖率的测试
make test-coverage

# 运行单个包的测试
go test ./internal/services/git/...
```

### 构建

```bash
# 开发模式
make build-debug

# 生产模式
make build

# 跨平台构建
make build-all
```

## 依赖管理

### 主要依赖

- `fyne.io/fyne/v2` - GUI 框架
- `modernc.org/sqlite` - SQLite 驱动
- `github.com/creack/pty` - PTY 终端
- `github.com/google/uuid` - UUID 生成

### 更新依赖

```bash
# 更新所有依赖
go get -u ./...
go mod tidy

# 更新特定依赖
go get -u fyne.io/fyne/v2@latest
```

## 发布流程

1. 更新 `CHANGELOG.md`
2. 打标签 `git tag -a v0.1.0 -m "Release v0.1.0"`
3. 推送标签 `git push origin v0.1.0`
4. GitHub Actions 自动构建发布

## 故障排查

### 常见问题

#### 编译错误：缺少 X11 库

```bash
# Ubuntu/Debian
sudo apt-get install libgl1-mesa-dev xorg-dev

# Fedora
sudo dnf install mesa-libGL-devel libXext-devel
```

#### 运行时错误：无法打开显示

```bash
# 确保在 X11 环境运行
export DISPLAY=:0
```

### 调试技巧

1. 启用调试标签：`go build -tags debug`
2. 查看日志输出
3. 使用 `fmt.Println` 调试关键路径

## 性能优化

1. 避免在主线程执行耗时操作
2. 使用异步加载大数据
3. 缓存频繁访问的数据
4. 限制列表项数量（最多 1000 行）

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

MIT License - 见 LICENSE 文件

