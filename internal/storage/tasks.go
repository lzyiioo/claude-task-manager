// internal/storage/tasks.go
package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/yourname/claude-task-manager/pkg/models"
)

// TaskStore handles persistence of tasks
type TaskStore struct {
	mu    sync.RWMutex
	tasks map[string]*models.Task
	path  string
}

// NewTaskStore creates a new task store
func NewTaskStore() (*TaskStore, error) {
	path, err := taskStorePath()
	if err != nil {
		return nil, err
	}

	store := &TaskStore{
		tasks: make(map[string]*models.Task),
		path:  path,
	}

	if err := store.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return store, nil
}

// Add adds a new task to the store
func (s *TaskStore) Add(task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks[task.ID] = task
	return s.save()
}

// Get retrieves a task by ID
func (s *TaskStore) Get(id string) (*models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	return task, ok
}

// Update updates an existing task
func (s *TaskStore) Update(task *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[task.ID]; !ok {
		return os.ErrNotExist
	}

	s.tasks[task.ID] = task
	return s.save()
}

// Delete removes a task from the store
func (s *TaskStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.tasks, id)
	return s.save()
}

// List returns all tasks
func (s *TaskStore) List() []*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*models.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		tasks = append(tasks, t)
	}
	return tasks
}

// ListByStatus returns tasks with a specific status
func (s *TaskStore) ListByStatus(status models.TaskStatus) []*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var tasks []*models.Task
	for _, t := range s.tasks {
		if t.Status == status {
			tasks = append(tasks, t)
		}
	}
	return tasks
}

// load reads tasks from disk
func (s *TaskStore) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}

	var tasks []*models.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return err
	}

	for _, t := range tasks {
		s.tasks[t.ID] = t
	}
	return nil
}

// save writes tasks to disk
func (s *TaskStore) save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0755); err != nil {
		return err
	}

	tasks := make([]*models.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		tasks = append(tasks, t)
	}

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0644)
}

func taskStorePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".claude-task-manager", "tasks.json"), nil
}
