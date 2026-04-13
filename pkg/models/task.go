// pkg/models/task.go
package models

import "time"

// TaskStatus represents the current state of a task
type TaskStatus string

const (
	TaskPending   TaskStatus = "pending"
	TaskRunning   TaskStatus = "running"
	TaskPaused    TaskStatus = "paused"    // waiting for permission
	TaskCompleted TaskStatus = "completed"
	TaskFailed    TaskStatus = "failed"
	TaskCancelled TaskStatus = "cancelled"
)

// IsolationType represents how the task is isolated
type IsolationType string

const (
	IsolationTmux   IsolationType = "tmux"
	IsolationDocker IsolationType = "docker"
)

// Task represents a single Claude Code execution task
type Task struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Prompt     string        `json:"prompt"`
	WorkDir    string        `json:"workDir"`
	Isolation  IsolationType `json:"isolation"`
	Status     TaskStatus    `json:"status"`
	CreatedAt  time.Time     `json:"createdAt"`
	StartedAt  *time.Time    `json:"startedAt,omitempty"`
	FinishedAt *time.Time    `json:"finishedAt,omitempty"`
	SessionID  string        `json:"sessionId,omitempty"`
}

// NewTask creates a new task with generated ID and pending status
func NewTask(name, prompt, workDir string, isolation IsolationType) *Task {
	return &Task{
		ID:        generateID(),
		Name:      name,
		Prompt:    prompt,
		WorkDir:   workDir,
		Isolation: isolation,
		Status:    TaskPending,
		CreatedAt: time.Now(),
	}
}

// Duration returns how long the task has been running
func (t *Task) Duration() time.Duration {
	if t.StartedAt == nil {
		return 0
	}
	end := time.Now()
	if t.FinishedAt != nil {
		end = *t.FinishedAt
	}
	return end.Sub(*t.StartedAt)
}

// IsRunning returns true if task is currently running or paused
func (t *Task) IsRunning() bool {
	return t.Status == TaskRunning || t.Status == TaskPaused
}
