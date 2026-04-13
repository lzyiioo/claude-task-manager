// pkg/models/task_test.go
package models

import (
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	task := NewTask("test task", "fix the bug", "/tmp/work", IsolationTmux)

	if task.ID == "" {
		t.Error("ID should not be empty")
	}
	if task.Name != "test task" {
		t.Errorf("Name = %q, want %q", task.Name, "test task")
	}
	if task.Status != TaskPending {
		t.Errorf("Status = %q, want %q", task.Status, TaskPending)
	}
	if task.Isolation != IsolationTmux {
		t.Errorf("Isolation = %q, want %q", task.Isolation, IsolationTmux)
	}
}

func TestTaskDuration(t *testing.T) {
	task := NewTask("test", "prompt", "/tmp", IsolationTmux)

	// No start time
	if d := task.Duration(); d != 0 {
		t.Errorf("Duration with no start = %v, want 0", d)
	}

	// Started but not finished
	start := time.Now().Add(-time.Minute)
	task.StartedAt = &start
	if d := task.Duration(); d < time.Minute {
		t.Errorf("Duration after start = %v, want >= 1m", d)
	}

	// Finished
	end := start.Add(30 * time.Second)
	task.FinishedAt = &end
	if d := task.Duration(); d != 30*time.Second {
		t.Errorf("Duration after finish = %v, want 30s", d)
	}
}

func TestTaskIsRunning(t *testing.T) {
	task := NewTask("test", "prompt", "/tmp", IsolationTmux)

	if task.IsRunning() {
		t.Error("pending task should not be running")
	}

	task.Status = TaskRunning
	if !task.IsRunning() {
		t.Error("running task should be running")
	}

	task.Status = TaskPaused
	if !task.IsRunning() {
		t.Error("paused task should be running")
	}
}
