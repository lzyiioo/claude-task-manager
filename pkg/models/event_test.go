// pkg/models/event_test.go
package models

import (
	"encoding/json"
	"testing"
)

func TestNewEvent(t *testing.T) {
	payload := json.RawMessage(`{"test": "data"}`)
	event := NewEvent("task-123", EventToolUse, payload)

	if event.ID == "" {
		t.Error("ID should not be empty")
	}
	if event.TaskID != "task-123" {
		t.Errorf("TaskID = %q, want %q", event.TaskID, "task-123")
	}
	if event.Type != EventToolUse {
		t.Errorf("Type = %q, want %q", event.Type, EventToolUse)
	}
	if event.Timestamp.IsZero() {
		t.Error("Timestamp should be set")
	}
}
