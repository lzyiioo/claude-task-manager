# Claude Task Manager (CTM)

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License">
  <img src="https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-blue?style=for-the-badge" alt="Platform">
</p>

> Terminal-based multi-task parallel execution platform for Claude Code

**Developer**: 小梁子 (xiaoliangzi)

---

## What is CTM?

Claude Task Manager (CTM) is a powerful terminal-based platform for running multiple Claude Code instances in parallel. It provides real-time monitoring, permission control, batch automation, and automatic resource cleanup.

### Use Cases

| Scenario | Problem | CTM Solution |
|----------|---------|--------------|
| **Batch Development** | Running Claude manually one by one is inefficient | Parallel batch execution with auto-iteration |
| **Multi-project Development** | Switching windows is messy, easy to get confused | Independent isolated environments |
| **Real-time Monitoring** | No visibility into task progress | TUI interface with live status |
| **Permission Approval** | Clicking Yes for every operation is tedious | Auto-approve, no manual interaction |
| **Resource Cleanup** | Forgetting to close processes/containers | Automatic cleanup, resource release |

### Target Users

- 🤖 **AI Developers** - Running multiple Claude Code instances simultaneously
- 📦 **DevOps Engineers** - Batch automation tasks
- 🧪 **Test Engineers** - Multi-scenario parallel testing
- 🚀 **Solo Developers** - Maximizing AI-assisted programming efficiency

---

## Features

### Core Features

| Feature | Description | Status |
|---------|-------------|--------|
| 🔄 **Parallel Tasks** | Run multiple Claude Code tasks simultaneously | ✅ |
| 🖥️ **TUI Interface** | Bubbletea terminal UI with real-time monitoring | ✅ |
| 🛡️ **Process Isolation** | tmux / Docker isolation modes | ✅ |
| 🔐 **Permission Control** | Real-time approval for Claude tool calls | ✅ |
| 📦 **Batch Execution** | Automated batch runs with multi-iteration support | ✅ |
| 📋 **Event Logging** | JSONL format for complete audit trail | ✅ |
| 🧹 **Auto Cleanup** | Automatic resource cleanup on task completion | ✅ |

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     CTM Architecture                         │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────┐   │
│  │              TUI Interface (Bubbletea)              │   │
│  │   Real-time status │ Keyboard shortcuts │ Live view │   │
│  └─────────────────────────────────────────────────────┘   │
│                           │                                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                   Task Manager                       │   │
│  │    Lifecycle │ Event Bus │ Permission Engine        │   │
│  └─────────────────────────────────────────────────────┘   │
│           │                               │                 │
│  ┌────────┴────────┐          ┌──────────┴──────────┐     │
│  │ Claude Runner   │          │   Batch Executor    │     │
│  │ Hook + Parser   │          │ Auto-iteration +    │     │
│  └─────────────────┘          │      Reports        │     │
│                           │                                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                   Isolation Layer                   │   │
│  │         ┌─────────┐        ┌─────────┐             │   │
│  │         │  tmux   │        │ Docker  │             │   │
│  │         └─────────┘        └─────────┘             │   │
│  └─────────────────────────────────────────────────────┘   │
│                           │                                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                   Storage Layer                      │   │
│  │    ┌─────────────┐      ┌─────────────┐            │   │
│  │    │  JSON       │      │  JSONL      │            │   │
│  │    │ (Tasks)     │      │ (Events)    │            │   │
│  │    └─────────────┘      └─────────────┘            │   │
│  └─────────────────────────────────────────────────────────────┘
└─────────────────────────────────────────────────────────────┘
```

---

## Documentation

### 1. TUI Interface

```bash
# Launch TUI
ctm tui
```

**Features**:
- Real-time task list display
- Task status at a glance (pending/running/completed/failed)
- Keyboard shortcuts

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `n` | New task |
| `s` | Stop selected task |
| `r` | Refresh list |
| `q` / `Ctrl+C` | Exit |

### 2. Batch Execution

```bash
# Basic usage - run 10 iterations
ctm batch run -p "Implement user login feature" -n 10

# Read prompt from file
ctm batch run -f prompts.txt -n 5

# Auto-approve all permissions (no manual Yes)
ctm batch run -p "Refactor code" -n 10 --auto-permission

# Dry run mode (don't actually execute)
ctm batch run -p "Test feature" --dry-run

# Specify working directory
ctm batch run -p "Code review" -w /path/to/project -n 3
```

**Batch Config File**:
```yaml
# batch-config.yaml
iterations: 10              # Number of iterations
prompt: "Implement new feature"  # Prompt content
promptFile: ""              # Or read from file
workDir: "."                # Working directory
permissionMode: "auto"      # Permission mode: ask/auto/auto-safe
isolation: "tmux"           # Isolation: tmux/docker
delayBetween: 5             # Delay between runs (seconds)
stopOnError: false          # Stop on error
onComplete: "commit"        # On complete action: commit/push/none
```

### 3. Process Isolation

| Mode | Description | Use Case |
|------|-------------|----------|
| **tmux** | Lightweight terminal multiplexing | Local dev, quick testing |
| **Docker** | Container-level isolation | Complete environment isolation |

```bash
# Use tmux isolation (default)
ctm batch run -p "task" --isolation tmux

# Use Docker isolation
ctm batch run -p "task" --isolation docker
```

### 4. Permission Control

```bash
# Manual approval mode (default)
ctm batch run -p "task" --permission-mode ask

# Auto-approve all permissions
ctm batch run -p "task" --permission-mode auto

# Safe mode - auto read, prompt write
ctm batch run -p "task" --permission-mode auto-safe
```

### 5. Auto Cleanup

**Automatic cleanup** (triggered on task complete/fail/cancel):
- Stop tmux session or Docker container
- Clean up associated resources
- Emit cleanup event log

**Manual cleanup**:
```bash
# Cleanup all resources
ctm cleanup

# Cleanup specific task
ctm cleanup -t <task-id>
```

### 6. Event Logging

**Location**: `~/.ctm/events.jsonl`

| Type | Description |
|------|-------------|
| `status` | Task status change |
| `tool_use` | Claude tool call |
| `tool_result` | Tool execution result |
| `permission` | Permission request |
| `cleanup` | Resource cleanup |

**View logs**:
```bash
# View last 20 events
tail -20 ~/.ctm/events.jsonl | jq .

# View specific task events
grep '"taskId":"task-xxx"' ~/.ctm/events.jsonl | jq .

# View cleanup events
grep '"type":"cleanup"' ~/.ctm/events.jsonl | jq .
```

---

## Quick Start

### Installation

```bash
# Clone repository
git clone https://github.com/lzyiioo/claude-task-manager.git
cd claude-task-manager

# Build
go build -o ctm ./cmd/...

# Or use Makefile
make build

# Add to PATH
# Linux/Mac:
sudo mv ctm /usr/local/bin/
# Windows:
# Move ctm.exe to a directory in PATH
```

### Dependencies

| Dependency | Required | Description |
|------------|----------|-------------|
| Go 1.21+ | ✅ | Build & run |
| tmux | Optional | tmux isolation mode |
| Docker | Optional | Docker isolation mode |

### Usage Examples

```bash
# 1. Launch TUI interface
ctm tui

# 2. Batch execute tasks
ctm batch run -p "Implement a calculator app" -n 5 --auto-permission

# 3. View task list
ctm task list

# 4. View event logs
tail -f ~/.ctm/events.jsonl

# 5. Cleanup resources
ctm cleanup
```

---

## Configuration

### Global Config

Location: `~/.ctm/config.yaml`

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

### Claude Code Permission Config

Configure auto-approve in `~/.claude/settings.json`:

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

## Production Deployment

### Startup Script

```bash
#!/bin/bash
# start-ctm.sh

# 1. Cleanup residual resources
ctm cleanup

# 2. Launch TUI
ctm tui

# 3. Auto cleanup on exit
trap 'ctm cleanup' EXIT
```

### Scheduled Cleanup (Cron)

```bash
# Cleanup hourly
0 * * * * ctm cleanup >> /var/log/ctm-cleanup.log 2>&1
```

### Health Check

```bash
#!/bin/bash
# health-check.sh

# Check residual sessions
TMUX_SESSIONS=$(tmux list-sessions 2>/dev/null | grep "ctm-" | wc -l)
if [ "$TMUX_SESSIONS" -gt 0 ]; then
    echo "WARNING: $TMUX_SESSIONS residual tmux sessions"
fi

# Check residual containers
DOCKER_CONTAINERS=$(docker ps --filter "name=ctm-" -q | wc -l)
if [ "$DOCKER_CONTAINERS" -gt 0 ]; then
    echo "WARNING: $DOCKER_CONTAINERS residual containers"
fi
```

---

## Project Structure

```
claude-task-manager/
├── cmd/                    # CLI commands
│   ├── main.go            # Entry point
│   ├── batch.go           # Batch commands
│   ├── cleanup.go         # Cleanup commands
│   └── tui.go             # TUI commands
├── internal/
│   ├── batch/             # Batch execution engine
│   ├── claude/            # Claude integration
│   ├── config/            # Configuration
│   ├── core/              # Core engine
│   ├── isolation/         # Isolation layer
│   ├── storage/           # Storage layer
│   └── tui/               # Terminal UI
├── pkg/models/            # Domain models
├── scripts/               # Helper scripts
├── Makefile               # Build script
└── go.mod                 # Go module
```

---

## Tech Stack

- **Language**: Go 1.21+
- **CLI Framework**: Cobra
- **TUI Framework**: Bubbletea + Lipgloss
- **Storage**: JSON + JSONL
- **Isolation**: tmux / Docker

---

## FAQ

### Q: How to batch generate code?
```bash
ctm batch run -p "Generate a TODO app" -n 10 --auto-permission
```

### Q: Task stuck怎么办?
```bash
# Stop specific task
ctm task stop <task-id>

# Or cleanup all resources
ctm cleanup
```

### Q: How to view task logs?
```bash
# Real-time events
tail -f ~/.ctm/events.jsonl | jq .

# Specific task
grep 'task-id' ~/.ctm/events.jsonl
```

### Q: Docker isolation fails?
```bash
# Ensure Docker is running
docker ps

# Use tmux mode instead
ctm batch run -p "task" --isolation tmux
```

---

## Changelog

### v1.0.0 (2026-04-13)
- ✅ Initial release
- ✅ Multi-task parallel execution
- ✅ TUI real-time interface
- ✅ tmux/Docker isolation
- ✅ Batch execution automation
- ✅ Auto permission approval
- ✅ Auto resource cleanup

---

## License

MIT License

---

## Contributing

Issues and Pull Requests are welcome!

---

**Developer**: 小梁子 (xiaoliangzi)  
**GitHub**: https://github.com/lzyiioo/claude-task-manager  
**Version**: 1.0.0
