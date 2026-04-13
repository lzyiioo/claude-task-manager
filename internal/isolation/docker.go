package isolation

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/yourname/claude-task-manager/pkg/models"
)

type DockerDriver struct {
	image string
}

func NewDockerDriver(image string) *DockerDriver {
	if image == "" {
		image = "claude-code:latest"
	}
	return &DockerDriver{image: image}
}

func (d *DockerDriver) Create(taskID, workDir string) (string, error) {
	containerName := fmt.Sprintf("ctm-%s", taskID)
	absWorkDir := workDir
	cmd := exec.Command("docker", "run", "-d",
		"--name", containerName,
		"-v", fmt.Sprintf("%s:/workspace", absWorkDir),
		"-w", "/workspace",
		d.image, "tail", "-f", "/dev/null",
	)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func (d *DockerDriver) Start(sessionID, command string) error {
	cmd := exec.Command("docker", "exec", sessionID, "sh", "-c", command)
	return cmd.Run()
}

func (d *DockerDriver) SendInput(sessionID, input string) error {
	cmd := exec.Command("docker", "exec", "-i", sessionID, "sh", "-c", input)
	return cmd.Run()
}

func (d *DockerDriver) CaptureOutput(sessionID string) (string, error) {
	cmd := exec.Command("docker", "logs", sessionID)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out.String(), nil
}

func (d *DockerDriver) Stop(sessionID string) error {
	cmd := exec.Command("docker", "stop", sessionID)
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("docker", "rm", sessionID)
	return cmd.Run()
}

func (d *DockerDriver) IsRunning(sessionID string) bool {
	cmd := exec.Command("docker", "inspect", "-f", "{{.State.Running}}", sessionID)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return false
	}
	return strings.TrimSpace(out.String()) == "true"
}

func (d *DockerDriver) Type() models.IsolationType {
	return models.IsolationDocker
}
