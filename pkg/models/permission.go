// pkg/models/permission.go
package models

import (
	"encoding/json"
	"time"
)

// PermissionRequest represents a request for user approval
type PermissionRequest struct {
	ID        string              `json:"id"`
	TaskID    string              `json:"taskId"`
	Tool      string              `json:"tool"`
	Input     json.RawMessage     `json:"input"`
	Risk      string              `json:"risk"`
	CreatedAt time.Time           `json:"createdAt"`
	Response  *PermissionResponse `json:"response,omitempty"`
}

// PermissionResponse represents user's response to a permission request
type PermissionResponse struct {
	Allow       bool      `json:"allow"`
	Reason      string    `json:"reason,omitempty"`
	RespondedAt time.Time `json:"respondedAt"`
}

// NewPermissionRequest creates a new permission request
func NewPermissionRequest(taskID, tool string, input json.RawMessage, risk string) *PermissionRequest {
	return &PermissionRequest{
		ID:        generateID(),
		TaskID:    taskID,
		Tool:      tool,
		Input:     input,
		Risk:      risk,
		CreatedAt: time.Now(),
	}
}

// Respond records a user response to the permission request
func (r *PermissionRequest) Respond(allow bool, reason string) {
	r.Response = &PermissionResponse{
		Allow:       allow,
		Reason:      reason,
		RespondedAt: time.Now(),
	}
}

// IsResponded returns true if the request has been responded to
func (r *PermissionRequest) IsResponded() bool {
	return r.Response != nil
}

// IsAllowed returns true if the request was allowed (false if denied or not responded)
func (r *PermissionRequest) IsAllowed() bool {
	return r.Response != nil && r.Response.Allow
}
