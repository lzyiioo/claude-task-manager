// internal/storage/events.go
package storage

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/yourname/claude-task-manager/pkg/models"
)

// EventStore handles persistence of events
type EventStore struct {
	mu    sync.Mutex
	files map[string]*os.File // taskID -> file handle
	dir   string
}

// NewEventStore creates a new event store
func NewEventStore() (*EventStore, error) {
	dir, err := eventStoreDir()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	return &EventStore{
		files: make(map[string]*os.File),
		dir:   dir,
	}, nil
}

// Append adds an event to the event log for its task
func (s *EventStore) Append(event *models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := s.getFile(event.TaskID)
	if err != nil {
		return err
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = file.Write(append(data, '\n'))
	return err
}

// GetEvents retrieves all events for a task
func (s *EventStore) GetEvents(taskID string) ([]*models.Event, error) {
	path := s.eventFilePath(taskID)

	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []*models.Event
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event models.Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			continue // skip malformed lines
		}
		events = append(events, &event)
	}

	return events, scanner.Err()
}

// Close closes all open file handles
func (s *EventStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var lastErr error
	for _, file := range s.files {
		if err := file.Close(); err != nil {
			lastErr = err
		}
	}
	s.files = make(map[string]*os.File)
	return lastErr
}

// getFile returns the file handle for a task's events
func (s *EventStore) getFile(taskID string) (*os.File, error) {
	if file, ok := s.files[taskID]; ok {
		return file, nil
	}

	path := s.eventFilePath(taskID)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	s.files[taskID] = file
	return file, nil
}

// eventFilePath returns the path to a task's event file
func (s *EventStore) eventFilePath(taskID string) string {
	return filepath.Join(s.dir, "task-"+taskID+".jsonl")
}

func eventStoreDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".claude-task-manager", "events"), nil
}
