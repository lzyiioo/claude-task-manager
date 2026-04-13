package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("62")).
			Padding(0, 1)

	taskStyle = lipgloss.NewStyle().
			Padding(0, 1)

	runningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	pendingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	completedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86"))

	failedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	dialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("196")).
			Padding(1, 2)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("62")).
			Padding(0, 2)

	buttonActiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("42")).
				Padding(0, 2).
				Bold(true)
)
