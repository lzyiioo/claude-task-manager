// pkg/models/permission_test.go
package models

import (
	"encoding/json"
	"testing"
)

func TestNewPermissionRequest(t *testing.T) {
	input := json.RawMessage(`{"command": "ls"}`)
	req := NewPermissionRequest("task-123", "Bash", input, "File system access")

	if req.ID == "" {
		t.Error("ID should not be empty")
	}
	if req.TaskID != "task-123" {
		t.Errorf("TaskID = %q, want %q", req.TaskID, "task-123")
	}
	if req.IsResponded() {
		t.Error("new request should not be responded")
	}
}

func TestPermissionRequestRespond(t *testing.T) {
	input := json.RawMessage(`{}`)
	req := NewPermissionRequest("task-123", "Bash", input, "test")

	req.Respond(true, "looks safe")

	if !req.IsResponded() {
		t.Error("request should be responded")
	}
	if !req.IsAllowed() {
		t.Error("request should be allowed")
	}
	if req.Response.Reason != "looks safe" {
		t.Errorf("Reason = %q, want %q", req.Response.Reason, "looks safe")
	}
}
