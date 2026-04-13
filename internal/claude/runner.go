package claude

import (
	"context"
	"os"
	"os/exec"

	"github.com/yourname/claude-task-manager/pkg/models"
)

type Runner struct {
	taskID     string
	prompt     string
	workDir    string
	isolation  models.IsolationType
	sessionID  string
	permission string // "auto", "interactive"
}

func NewRunner(taskID, prompt, workDir string, isolation models.IsolationType) *Runner {
	return &Runner{
		taskID:     taskID,
		prompt:     prompt,
		workDir:    workDir,
		isolation:  isolation,
		permission: "auto",
	}
}

func (r *Runner) SetPermissionMode(mode string) {
	r.permission = mode
}

func (r *Runner) BuildCommand() *exec.Cmd {
	args := []string{"--print"}

	if r.permission == "auto" {
		args = append(args, "--dangerously-skip-permissions")
	}

	args = append(args, r.prompt)

	cmd := exec.Command("claude", args...)
	cmd.Dir = r.workDir
	cmd.Env = os.Environ()

	return cmd
}

func (r *Runner) Run(ctx context.Context) error {
	cmd := r.BuildCommand()
	return cmd.Run()
}
