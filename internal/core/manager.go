// internal/core/manager.go
package core

import (
	"time"

	"github.com/yourname/claude-task-manager/internal/storage"
	"github.com/yourname/claude-task-manager/pkg/models"
)

// TaskManager handles task lifecycle operations
type TaskManager struct {
	taskStore    *storage.TaskStore
	eventStore   *storage.EventStore
	eventBus     *EventBus
	now          func() time.Time
}

// NewTaskManager creates a new task manager
func NewTaskManager(taskStore *storage.TaskStore, eventStore *storage.EventStore, eventBus *EventBus) *TaskManager {
	return &TaskManager{
		taskStore:  taskStore,
		eventStore: eventStore,
		eventBus:   eventBus,
		now:        time.Now,
	}
}

// Create creates a new task
func (m *TaskManager) Create(name, prompt, workDir string, isolation models.IsolationType) (*models.Task, error) {
	task := models.NewTask(name, prompt, workDir, isolation)
	if err := m.taskStore.Add(task); err != nil {
		return nil, err
	}
	return task, nil
}

// Start begins execution of a task
func (m *TaskManager) Start(taskID string) error {
	task, ok := m.taskStore.Get(taskID)
	if !ok {
		return ErrTaskNotFound
	}

	now := m.now()
	task.Status = models.TaskRunning
	task.StartedAt = &now

	if err := m.taskStore.Update(task); err != nil {
		return err
	}

	m.emitStatusChange(task, models.TaskPending, models.TaskRunning)
	return nil
}

// Pause pauses a running task
func (m *TaskManager) Pause(taskID string) error {
	task, ok := m.taskStore.Get(taskID)
	if !ok {
		return ErrTaskNotFound
	}

	if task.Status != models.TaskRunning {
		return ErrInvalidStatus
	}

	task.Status = models.TaskPaused
	if err := m.taskStore.Update(task); err != nil {
		return err
	}

	m.emitStatusChange(task, models.TaskRunning, models.TaskPaused)
	return nil
}

// Resume resumes a paused task
func (m *TaskManager) Resume(taskID string) error {
	task, ok := m.taskStore.Get(taskID)
	if !ok {
		return ErrTaskNotFound
	}

	if task.Status != models.TaskPaused {
		return ErrInvalidStatus
	}

	task.Status = models.TaskRunning
	if err := m.taskStore.Update(task); err != nil {
		return err
	}

	m.emitStatusChange(task, models.TaskPaused, models.TaskRunning)
	return nil
}

// Complete marks a task as completed
func (m *TaskManager) Complete(taskID string) error {
	task, ok := m.taskStore.Get(taskID)
	if !ok {
		return ErrTaskNotFound
	}

	now := m.now()
	task.Status = models.TaskCompleted
	task.FinishedAt = &now

	if err := m.taskStore.Update(task); err != nil {
		return err
	}

	m.emitStatusChange(task, task.Status, models.TaskCompleted)
	return nil
}

// Fail marks a task as failed
func (m *TaskManager) Fail(taskID string) error {
	task, ok := m.taskStore.Get(taskID)
	if !ok {
		return ErrTaskNotFound
	}

	now := m.now()
	task.Status = models.TaskFailed
	task.FinishedAt = &now

	if err := m.taskStore.Update(task); err != nil {
		return err
	}

	m.emitStatusChange(task, task.Status, models.TaskFailed)
	return nil
}

// Cancel cancels a task
func (m *TaskManager) Cancel(taskID string) error {
	task, ok := m.taskStore.Get(taskID)
	if !ok {
		return ErrTaskNotFound
	}

	now := m.now()
	task.Status = models.TaskCancelled
	task.FinishedAt = &now

	if err := m.taskStore.Update(task); err != nil {
		return err
	}

	m.emitStatusChange(task, task.Status, models.TaskCancelled)
	return nil
}

// Get retrieves a task by ID
func (m *TaskManager) Get(taskID string) (*models.Task, bool) {
	return m.taskStore.Get(taskID)
}

// List returns all tasks
func (m *TaskManager) List() []*models.Task {
	return m.taskStore.List()
}

// ListByStatus returns tasks with a specific status
func (m *TaskManager) ListByStatus(status models.TaskStatus) []*models.Task {
	return m.taskStore.ListByStatus(status)
}

// emitStatusChange emits a status change event
func (m *TaskManager) emitStatusChange(task *models.Task, oldStatus, newStatus models.TaskStatus) {
	payload := models.StatusPayload{
		OldStatus: oldStatus,
		NewStatus: newStatus,
	}
	event := models.NewEvent(task.ID, models.EventStatus, mustMarshal(payload))
	m.eventStore.Append(event)
	m.eventBus.Publish(event)
}
