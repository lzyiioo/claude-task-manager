// pkg/models/event.go
package models

import (
	"encoding/json"
	"time"
)

// EventType represents the type of event
type EventType string

const (
	EventToolUse    EventType = "tool_use"
	EventToolResult EventType = "tool_result"
	EventPermission EventType = "permission"
	EventMessage    EventType = "message"
	EventError      EventType = "error"
	EventStatus     EventType = "status" // task status change
)

// Event represents a single event in the event stream
type Event struct {
	ID        string          `json:"id"`
	TaskID    string          `json:"taskId"`
	Type      EventType       `json:"type"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

// NewEvent creates a new event with generated ID and current timestamp
func NewEvent(taskID string, eventType EventType, payload json.RawMessage) *Event {
	return &Event{
		ID:        generateID(),
		TaskID:    taskID,
		Type:      eventType,
		Timestamp: time.Now(),
		Payload:   payload,
	}
}

// ToolUsePayload represents payload for tool_use events
type ToolUsePayload struct {
	Tool   string          `json:"tool"`
	Input  json.RawMessage `json:"input"`
}

// ToolResultPayload represents payload for tool_result events
type ToolResultPayload struct {
	Tool    string `json:"tool"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// MessagePayload represents payload for message events
type MessagePayload struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// StatusPayload represents payload for status change events
type StatusPayload struct {
	OldStatus TaskStatus `json:"oldStatus"`
	NewStatus TaskStatus `json:"newStatus"`
}
