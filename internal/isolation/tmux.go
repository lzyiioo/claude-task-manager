package isolation

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/yourname/claude-task-manager/pkg/models"
)

type TmuxDriver struct{}

func NewTmuxDriver() *TmuxDriver {
	return &TmuxDriver{}
}

func (d *TmuxDriver) Create(taskID, workDir string) (string, error) {
	sessionName := fmt.Sprintf("ctm-%s", taskID)
	cmd := exec.Command("tmux", "new-session", "-d", "-s", sessionName, "-c", workDir)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return sessionName, nil
}

func (d *TmuxDriver) Start(sessionID, command string) error {
	cmd := exec.Command("tmux", "send-keys", "-t", sessionID, command, "Enter")
	return cmd.Run()
}

func (d *TmuxDriver) SendInput(sessionID, input string) error {
	cmd := exec.Command("tmux", "send-keys", "-t", sessionID, input, "Enter")
	return cmd.Run()
}

func (d *TmuxDriver) CaptureOutput(sessionID string) (string, error) {
	cmd := exec.Command("tmux", "capture-pane", "-t", sessionID, "-p")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func (d *TmuxDriver) Stop(sessionID string) error {
	cmd := exec.Command("tmux", "kill-session", "-t", sessionID)
	return cmd.Run()
}

func (d *TmuxDriver) IsRunning(sessionID string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", sessionID)
	return cmd.Run() == nil
}

func (d *TmuxDriver) Type() models.IsolationType {
	return models.IsolationTmux
}
