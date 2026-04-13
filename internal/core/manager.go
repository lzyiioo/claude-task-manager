// internal/core/manager.go
package core

import (
	"sync"
	"time"

	"github.com/yourname/claude-task-manager/internal/isolation"
	"github.com/yourname/claude-task-manager/internal/storage"
	"github.com/yourname/claude-task-manager/pkg/models"
)

// TaskManager handles task lifecycle operations
type TaskManager struct {
	taskStore    *storage.TaskStore
	eventStore   *storage.EventStore
	eventBus     *EventBus
	drivers      map[models.IsolationType]isolation.Driver
	sessions     map[string]string // taskID -> sessionID
	sessionMu    sync.RWMutex
	now          func() time.Time
}

// NewTaskManager creates a new task manager
func NewTaskManager(taskStore *storage.TaskStore, eventStore *storage.EventStore, eventBus *EventBus) *TaskManager {
	return &TaskManager{
		taskStore:  taskStore,
		eventStore: eventStore,
		eventBus:   eventBus,
		drivers:    make(map[models.IsolationType]isolation.Driver),
		sessions:   make(map[string]string),
		now:        time.Now,
	}
}

// RegisterDriver registers an isolation driver
func (m *TaskManager) RegisterDriver(driver isolation.Driver) {
	m.drivers[driver.Type()] = driver
}

// RegisterSession associates a session ID with a task
func (m *TaskManager) RegisterSession(taskID, sessionID string) {
	m.sessionMu.Lock()
	defer m.sessionMu.Unlock()
	m.sessions[taskID] = sessionID
}

// GetSession returns the session ID for a task
func (m *TaskManager) GetSession(taskID string) (string, bool) {
	m.sessionMu.RLock()
	defer m.sessionMu.RUnlock()
	sessionID, ok := m.sessions[taskID]
	return sessionID, ok
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

// Complete marks a task as completed and cleans up resources
func (m *TaskManager) Complete(taskID string) error {
	task, ok := m.taskStore.Get(taskID)
	if !ok {
		return ErrTaskNotFound
	}

	// Cleanup isolation resources first
	m.cleanupSession(taskID, task.Isolation)

	now := m.now()
	task.Status = models.TaskCompleted
	task.FinishedAt = &now

	if err := m.taskStore.Update(task); err != nil {
		return err
	}

	m.emitStatusChange(task, task.Status, models.TaskCompleted)
	return nil
}

// Fail marks a task as failed and cleans up resources
func (m *TaskManager) Fail(taskID string) error {
	task, ok := m.taskStore.Get(taskID)
	if !ok {
		return ErrTaskNotFound
	}

	// Cleanup isolation resources first
	m.cleanupSession(taskID, task.Isolation)

	now := m.now()
	task.Status = models.TaskFailed
	task.FinishedAt = &now

	if err := m.taskStore.Update(task); err != nil {
		return err
	}

	m.emitStatusChange(task, task.Status, models.TaskFailed)
	return nil
}

// Cancel cancels a task and cleans up resources
func (m *TaskManager) Cancel(taskID string) error {
	task, ok := m.taskStore.Get(taskID)
	if !ok {
		return ErrTaskNotFound
	}

	// Cleanup isolation resources first
	m.cleanupSession(taskID, task.Isolation)

	now := m.now()
	task.Status = models.TaskCancelled
	task.FinishedAt = &now

	if err := m.taskStore.Update(task); err != nil {
		return err
	}

	m.emitStatusChange(task, task.Status, models.TaskCancelled)
	return nil
}

// cleanupSession stops and removes the isolation session for a task
func (m *TaskManager) cleanupSession(taskID string, isolationType models.IsolationType) {
	m.sessionMu.Lock()
	sessionID, exists := m.sessions[taskID]
	if exists {
		delete(m.sessions, taskID)
	}
	m.sessionMu.Unlock()

	if !exists || sessionID == "" {
		return
	}

	driver, ok := m.drivers[isolationType]
	if !ok {
		return
	}

	// Stop the session
	if driver.IsRunning(sessionID) {
		driver.Stop(sessionID)
	}

	// Emit cleanup event
	event := models.NewEvent(taskID, models.EventCleanup, mustMarshal(map[string]string{
		"sessionId": sessionID,
		"type":      string(isolationType),
	}))
	m.eventStore.Append(event)
	m.eventBus.Publish(event)
}

// CleanupAll cleans up all active sessions
func (m *TaskManager) CleanupAll() {
	m.sessionMu.RLock()
	taskSessions := make(map[string]string)
	for k, v := range m.sessions {
		taskSessions[k] = v
	}
	m.sessionMu.RUnlock()

	for taskID := range taskSessions {
		task, ok := m.taskStore.Get(taskID)
		if !ok {
			continue
		}
		m.cleanupSession(taskID, task.Isolation)
	}
}

// CleanupTask explicitly cleans up a task's resources
func (m *TaskManager) CleanupTask(taskID string) error {
	task, ok := m.taskStore.Get(taskID)
	if !ok {
		return ErrTaskNotFound
	}
	m.cleanupSession(taskID, task.Isolation)
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
