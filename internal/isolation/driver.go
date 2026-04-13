package isolation

import "github.com/yourname/claude-task-manager/pkg/models"

// Driver defines the interface for task isolation
type Driver interface {
	// Create creates a new isolated environment
	Create(taskID, workDir string) (sessionID string, err error)
	// Start starts a command in the session
	Start(sessionID, command string) error
	// SendInput sends input to the session
	SendInput(sessionID, input string) error
	// CaptureOutput captures current output
	CaptureOutput(sessionID string) (string, error)
	// Stop stops the session
	Stop(sessionID string) error
	// IsRunning checks if session is active
	IsRunning(sessionID string) bool
	// Type returns the driver type
	Type() models.IsolationType
}
