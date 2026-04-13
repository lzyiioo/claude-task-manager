# Claude Task Manager (CTM)

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License">
  <img src="https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-blue?style=for-the-badge" alt="Platform">
</p>

> 终端多任务并行执行平台 | Terminal-based multi-task parallel execution platform for Claude Code

**开发者**: 小梁子

---

## 作用之处

### 解决什么问题？

| 场景 | 痛点 | CTM 解决方案 |
|------|------|-------------|
| **批量开发任务** | 手动一个个运行 Claude，效率低 | 并行批量执行，自动迭代 |
| **多项目同时开发** | 切换窗口麻烦，容易混淆 | 独立隔离环境，互不干扰 |
| **需要实时监控** | 不知道任务进度 | TUI 实时界面，一目了然 |
| **权限审批繁琐** | 每个操作都要点 Yes | 自动批准，无需手动 |
| **资源残留清理** | 进程/容器忘记关闭 | 自动清理，释放资源 |

### 适用人群

- 🤖 **AI 开发者** - 同时运行多个 Claude Code 实例
- 📦 **DevOps 工程师** - 批量自动化任务
- 🧪 **测试工程师** - 多场景并行测试
- 🚀 **独立开发者** - 高效利用 AI 辅助编程

---

## 功能特点

### 核心功能

| 功能 | 说明 | 状态 |
|------|------|------|
| 🔄 **多任务并行** | 同时运行多个 Claude Code 任务 | ✅ |
| 🖥️ **TUI 界面** | Bubbletea 终端界面，实时监控 | ✅ |
| 🛡️ **进程隔离** | tmux / Docker 两种隔离方式 | ✅ |
| 🔐 **权限控制** | 实时审批 Claude 工具调用 | ✅ |
| 📦 **批量执行** | 自动化批量运行，支持多轮迭代 | ✅ |
| 📋 **事件日志** | JSONL 格式完整审计追踪 | ✅ |
| 🧹 **自动清理** | 任务完成后自动释放资源 | ✅ |

### 亮点特性

```
┌─────────────────────────────────────────────────────────────┐
│                     CTM 架构图                               │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────┐   │
│  │              TUI 界面 (Bubbletea)                    │   │
│  │     实时显示任务状态 │ 键盘快捷操作 │ 状态一目了然    │   │
│  └─────────────────────────────────────────────────────┘   │
│                           │                                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                   任务管理器                          │   │
│  │    生命周期管理 │ 事件总线 │ 权限引擎                │   │
│  └─────────────────────────────────────────────────────┘   │
│           │                               │                 │
│  ┌────────┴────────┐          ┌──────────┴──────────┐     │
│  │   Claude 运行器  │          │    批量执行器       │     │
│  │ Hook + 解析器    │          │  自动迭代 + 报告    │     │
│  └─────────────────┘          └─────────────────────┘     │
│                           │                                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                   隔离层                             │   │
│  │         ┌─────────┐        ┌─────────┐             │   │
│  │         │  tmux   │        │ Docker  │             │   │
│  │         └─────────┘        └─────────┘             │   │
│  └─────────────────────────────────────────────────────┘   │
│                           │                                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                   存储层                             │   │
│  │    ┌─────────────┐      ┌─────────────┐            │   │
│  │    │  JSON       │      │  JSONL      │            │   │
│  │    │ (任务持久化) │      │ (事件日志)  │            │   │
│  │    └─────────────┘      └─────────────┘            │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

## 功能介绍

### 1. TUI 终端界面

```bash
# 启动 TUI
ctm tui
```

**功能**:
- 任务列表实时显示
- 任务状态一目了然 (pending/running/completed/failed)
- 键盘快捷操作

**快捷键**:
| 按键 | 功能 |
|------|------|
| `j` / `↓` | 向下移动 |
| `k` / `↑` | 向上移动 |
| `n` | 新建任务 |
| `s` | 停止选中任务 |
| `r` | 刷新列表 |
| `q` / `Ctrl+C` | 退出 |

### 2. 批量执行

```bash
# 基本用法 - 执行 10 次
ctm batch run -p "实现用户登录功能" -n 10

# 从文件读取 prompt
ctm batch run -f prompts.txt -n 5

# 自动批准所有权限 (无需手动点 Yes)
ctm batch run -p "重构代码" -n 10 --auto-permission

# 预览模式 (不实际执行)
ctm batch run -p "测试功能" --dry-run

# 指定工作目录
ctm batch run -p "代码审查" -w /path/to/project -n 3
```

**批量配置文件**:
```yaml
# batch-config.yaml
iterations: 10              # 迭代次数
prompt: "实现新功能"         # prompt 内容
promptFile: ""              # 或从文件读取
workDir: "."                # 工作目录
permissionMode: "auto"      # 权限模式: ask/auto/auto-safe
isolation: "tmux"           # 隔离方式: tmux/docker
delayBetween: 5             # 每次执行间隔(秒)
stopOnError: false          # 错误时停止
onComplete: "commit"        # 完成动作: commit/push/none
```

### 3. 进程隔离

| 模式 | 说明 | 适用场景 |
|------|------|---------|
| **tmux** | 轻量级终端复用 | 本地开发、快速测试 |
| **Docker** | 容器级隔离 | 需要完全隔离环境 |

```bash
# 使用 tmux 隔离 (默认)
ctm batch run -p "任务" --isolation tmux

# 使用 Docker 隔离
ctm batch run -p "任务" --isolation docker
```

### 4. 权限控制

```bash
# 手动审批模式 (默认)
ctm batch run -p "任务" --permission-mode ask

# 自动批准所有权限
ctm batch run -p "任务" --permission-mode auto

# 安全模式 - 读操作自动，写操作审批
ctm batch run -p "任务" --permission-mode auto-safe
```

### 5. 自动清理

**自动清理** (任务完成/失败/取消时自动触发):
- 停止 tmux session 或 Docker 容器
- 清理关联资源
- 发送清理事件日志

**手动清理**:
```bash
# 清理所有资源
ctm cleanup

# 清理指定任务
ctm cleanup -t <task-id>
```

### 6. 事件日志

**存储位置**: `~/.ctm/events.jsonl`

**事件类型**:
| 类型 | 说明 |
|------|------|
| `status` | 任务状态变更 |
| `tool_use` | Claude 工具调用 |
| `tool_result` | 工具执行结果 |
| `permission` | 权限请求 |
| `cleanup` | 资源清理 |

**查看日志**:
```bash
# 查看最近 20 条
tail -20 ~/.ctm/events.jsonl | jq .

# 查看特定任务事件
grep '"taskId":"task-xxx"' ~/.ctm/events.jsonl | jq .

# 查看清理事件
grep '"type":"cleanup"' ~/.ctm/events.jsonl | jq .
```

---

## 快速开始

### 安装

```bash
# 克隆仓库
git clone https://github.com/lzyiioo/claude-task-manager.git
cd claude-task-manager

# 编译
go build -o ctm ./cmd/...

# 或使用 Makefile
make build

# 添加到 PATH
# Linux/Mac:
sudo mv ctm /usr/local/bin/
# Windows:
# 将 ctm.exe 移到 PATH 中的目录
```

### 依赖

| 依赖 | 必需 | 说明 |
|------|------|------|
| Go 1.21+ | ✅ | 编译运行 |
| tmux | 可选 | tmux 隔离模式 |
| Docker | 可选 | Docker 隔离模式 |

### 使用示例

```bash
# 1. 启动 TUI 界面
ctm tui

# 2. 批量执行任务
ctm batch run -p "实现一个计算器功能" -n 5 --auto-permission

# 3. 查看任务列表
ctm task list

# 4. 查看事件日志
tail -f ~/.ctm/events.jsonl

# 5. 清理资源
ctm cleanup
```

---

## 配置

### 全局配置

位置: `~/.ctm/config.yaml`

```yaml
storage:
  tasksPath: ~/.ctm/tasks.json
  eventsPath: ~/.ctm/events.jsonl

isolation:
  default: tmux
  tmux:
    socketPath: /tmp/ctm-tmux.sock
  docker:
    image: claude-code:latest
    network: ctm-network

cleanup:
  autoOnComplete: true
  autoOnFail: true
  autoOnCancel: true
  timeout: 30s

batch:
  maxConcurrency: 5
  retryCount: 3
  retryDelay: 5s

logging:
  level: info
  format: json
  output: ~/.ctm/events.jsonl
```

### Claude Code 权限配置

在 `~/.claude/settings.json` 中配置自动批准:

```json
{
  "permissions": {
    "autoApprove": true,
    "allowCommands": ["**/*"],
    "allowRead": ["**/*"],
    "allowWrite": ["**/*"]
  }
}
```

---

## 生产环境部署

### 启动脚本

```bash
#!/bin/bash
# start-ctm.sh

# 1. 清理残留资源
ctm cleanup

# 2. 启动 TUI
ctm tui

# 3. 退出时自动清理
trap 'ctm cleanup' EXIT
```

### 定时清理 (Cron)

```bash
# 每小时清理超时任务
0 * * * * ctm cleanup >> /var/log/ctm-cleanup.log 2>&1
```

### 健康检查

```bash
#!/bin/bash
# health-check.sh

# 检查残留 session
TMUX_SESSIONS=$(tmux list-sessions 2>/dev/null | grep "ctm-" | wc -l)
if [ "$TMUX_SESSIONS" -gt 0 ]; then
    echo "WARNING: $TMUX_SESSIONS 个残留 tmux session"
fi

# 检查残留容器
DOCKER_CONTAINERS=$(docker ps --filter "name=ctm-" -q | wc -l)
if [ "$DOCKER_CONTAINERS" -gt 0 ]; then
    echo "WARNING: $DOCKER_CONTAINERS 个残留容器"
fi
```

---

## 项目结构

```
claude-task-manager/
├── cmd/                    # CLI 命令
│   ├── main.go            # 入口
│   ├── batch.go           # 批量命令
│   ├── cleanup.go         # 清理命令
│   └── tui.go             # TUI 命令
├── internal/
│   ├── batch/             # 批量执行引擎
│   ├── claude/            # Claude 集成
│   ├── config/            # 配置管理
│   ├── core/              # 核心引擎
│   ├── isolation/         # 隔离层
│   ├── storage/           # 存储层
│   └── tui/               # 终端界面
├── pkg/models/            # 领域模型
├── scripts/               # 辅助脚本
├── Makefile               # 构建脚本
└── go.mod                 # Go 模块
```

---

## 技术栈

- **语言**: Go 1.21+
- **CLI 框架**: Cobra
- **TUI 框架**: Bubbletea + Lipgloss
- **存储**: JSON + JSONL
- **隔离**: tmux / Docker

---

## 常见问题

### Q: 如何批量生成代码?
```bash
ctm batch run -p "生成一个 TODO 应用" -n 10 --auto-permission
```

### Q: 任务卡住了怎么办?
```bash
# 停止指定任务
ctm task stop <task-id>

# 或清理所有资源
ctm cleanup
```

### Q: 如何查看任务日志?
```bash
# 实时查看事件
tail -f ~/.ctm/events.jsonl | jq .

# 查看特定任务
grep 'task-id' ~/.ctm/events.jsonl
```

### Q: Docker 隔离启动失败?
```bash
# 确保 Docker 正在运行
docker ps

# 使用 tmux 模式
ctm batch run -p "任务" --isolation tmux
```

---

## 更新日志

### v1.0.0 (2026-04-13)
- ✅ 初始版本发布
- ✅ 多任务并行执行
- ✅ TUI 实时界面
- ✅ tmux/Docker 隔离
- ✅ 批量执行自动化
- ✅ 权限自动批准
- ✅ 自动资源清理

---

## 许可证

MIT License

---

## 贡献

欢迎提交 Issue 和 Pull Request！

---

**开发者**: 小梁子  
**GitHub**: https://github.com/lzyiioo/claude-task-manager  
**版本**: 1.0.0
