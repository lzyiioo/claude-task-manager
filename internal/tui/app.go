package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourname/claude-task-manager/internal/core"
	"github.com/yourname/claude-task-manager/pkg/models"
)

type Model struct {
	taskMgr   *core.TaskManager
	tasks     []*models.Task
	selected  int
	width     int
	height    int
	showPerm  bool
	permReq   *models.PermissionRequest
	permFocus int // 0=allow, 1=deny
}

func NewModel(taskMgr *core.TaskManager) Model {
	return Model{
		taskMgr: taskMgr,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.loadTasks,
	)
}

func (m Model) loadTasks() tea.Msg {
	m.tasks = m.taskMgr.List()
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.showPerm {
			return m.updatePermDialog(msg)
		}
		return m.updateMain(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m Model) updateMain(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "j", "down":
		if m.selected < len(m.tasks)-1 {
			m.selected++
		}
	case "k", "up":
		if m.selected > 0 {
			m.selected--
		}
	case "n":
		// TODO: new task dialog
	case "s":
		if len(m.tasks) > 0 && m.selected < len(m.tasks) {
			m.taskMgr.Cancel(m.tasks[m.selected].ID)
			m.tasks = m.taskMgr.List()
		}
	case "r":
		m.tasks = m.taskMgr.List()
	}
	return m, nil
}

func (m Model) updatePermDialog(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "right", "left":
		m.permFocus = 1 - m.permFocus
	case "enter":
		// allow := m.permFocus == 0
		// TODO: respond to permission
		m.showPerm = false
	case "y":
		m.permFocus = 0
	case "n":
		m.permFocus = 1
	}
	return m, nil
}

func (m Model) View() string {
	if m.showPerm {
		return m.viewPermDialog()
	}
	return m.viewMain()
}

func (m Model) viewMain() string {
	var b strings.Builder

	// Title
	title := titleStyle.Render(" Claude Task Manager ")
	help := helpStyle.Render("[h]elp [q]uit")
	header := lipgloss.JoinHorizontal(lipgloss.Top, title, " ", help)
	b.WriteString(header + "\n\n")

	// Task list
	b.WriteString("Tasks\n")
	for i, task := range m.tasks {
		style := taskStyle
		if i == m.selected {
			style = style.Background(lipgloss.Color("236"))
		}

		status := string(task.Status)
		switch task.Status {
		case models.TaskRunning:
			status = runningStyle.Render("* running")
		case models.TaskPending:
			status = pendingStyle.Render("o pending")
		case models.TaskCompleted:
			status = completedStyle.Render("+ completed")
		case models.TaskFailed:
			status = failedStyle.Render("x failed")
		case models.TaskPaused:
			status = pendingStyle.Render("|| paused")
		case models.TaskCancelled:
			status = failedStyle.Render("x cancelled")
		}

		line := fmt.Sprintf("  %s %s", status, task.Name)
		b.WriteString(style.Render(line) + "\n")
	}

	if len(m.tasks) == 0 {
		b.WriteString(helpStyle.Render("  No tasks. Press 'n' to create one.\n"))
	}

	// Footer
	b.WriteString("\n" + helpStyle.Render("[n]ew [s]top [r]efresh [a]ttach"))

	return b.String()
}

func (m Model) viewPermDialog() string {
	var b strings.Builder
	b.WriteString("Permission Required\n\n")
	if m.permReq != nil {
		b.WriteString(fmt.Sprintf("Tool: %s\n", m.permReq.Tool))
		b.WriteString(fmt.Sprintf("Risk: %s\n", m.permReq.Risk))
	}
	b.WriteString("\n")

	allowBtn := buttonStyle.Render("[y] Allow")
	denyBtn := buttonStyle.Render("[n] Deny")

	if m.permFocus == 0 {
		allowBtn = buttonActiveStyle.Render("[y] Allow")
	} else {
		denyBtn = buttonActiveStyle.Render("[n] Deny")
	}

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, allowBtn, "  ", denyBtn))

	return dialogStyle.Render(b.String())
}

func Run(taskMgr *core.TaskManager) error {
	p := tea.NewProgram(NewModel(taskMgr))
	_, err := p.Run()
	return err
}
