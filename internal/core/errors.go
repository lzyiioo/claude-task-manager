// internal/core/errors.go
package core

import "errors"

var (
	// ErrTaskNotFound indicates the task does not exist
	ErrTaskNotFound = errors.New("task not found")

	// ErrInvalidStatus indicates the task is not in a valid state for the operation
	ErrInvalidStatus = errors.New("invalid task status for operation")

	// ErrTaskNotRunning indicates the task is not currently running
	ErrTaskNotRunning = errors.New("task not running")

	// ErrSessionNotFound indicates the session does not exist
	ErrSessionNotFound = errors.New("session not found")

	// ErrRequestNotFound indicates the permission request does not exist
	ErrRequestNotFound = errors.New("permission request not found")

	// ErrTimeout indicates a permission request timed out
	ErrTimeout = errors.New("permission request timed out")
)
