// internal/core/permission.go
package core

import (
	"sync"
	"time"

	"github.com/yourname/claude-task-manager/pkg/models"
)

// PermissionEngine handles permission requests and responses
type PermissionEngine struct {
	mu       sync.Mutex
	requests map[string]*models.PermissionRequest
	timeout  time.Duration
}

// NewPermissionEngine creates a new permission engine
func NewPermissionEngine() *PermissionEngine {
	return &PermissionEngine{
		requests: make(map[string]*models.PermissionRequest),
		timeout:  30 * time.Second,
	}
}

// Request creates a new permission request
func (e *PermissionEngine) Request(taskID, tool string, input []byte, risk string) *models.PermissionRequest {
	req := models.NewPermissionRequest(taskID, tool, input, risk)
	e.mu.Lock()
	e.requests[req.ID] = req
	e.mu.Unlock()
	return req
}

// Respond records a response to a permission request
func (e *PermissionEngine) Respond(reqID string, allow bool, reason string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	req, ok := e.requests[reqID]
	if !ok {
		return ErrRequestNotFound
	}

	req.Respond(allow, reason)
	return nil
}

// Get retrieves a permission request by ID
func (e *PermissionEngine) Get(reqID string) (*models.PermissionRequest, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	req, ok := e.requests[reqID]
	return req, ok
}

// WaitForResponse blocks until a response is received or timeout
func (e *PermissionEngine) WaitForResponse(reqID string) (*models.PermissionResponse, error) {
	deadline := time.Now().Add(e.timeout)
	for time.Now().Before(deadline) {
		if req, ok := e.Get(reqID); ok && req.IsResponded() {
			return req.Response, nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil, ErrTimeout
}
