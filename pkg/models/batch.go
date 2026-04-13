// pkg/models/batch.go
package models

// PermissionMode represents how permissions are handled in batch mode
type PermissionMode string

const (
	PermissionAuto       PermissionMode = "auto"       // auto-allow all
	PermissionInteractive PermissionMode = "interactive" // prompt for each
	PermissionAutoSafe   PermissionMode = "auto-safe"  // auto-allow read, prompt write
)

// OnCompleteAction represents what to do after each iteration completes
type OnCompleteAction string

const (
	OnCompleteCommit OnCompleteAction = "commit"
	OnCompletePush   OnCompleteAction = "push"
	OnCompleteNone   OnCompleteAction = "none"
)

// BatchConfig represents configuration for batch execution
type BatchConfig struct {
	Iterations     int              `json:"iterations"`
	Prompt         string           `json:"prompt,omitempty"`
	PromptFile     string           `json:"promptFile,omitempty"`
	WorkDir        string           `json:"workDir"`
	PermissionMode PermissionMode   `json:"permissionMode"`
	DelayBetween   int              `json:"delayBetween"` // seconds
	StopOnError    bool             `json:"stopOnError"`
	OnComplete     OnCompleteAction `json:"onComplete"`
}

// BatchStatus represents the current state of a batch run
type BatchStatus struct {
	Total       int `json:"total"`
	Completed   int `json:"completed"`
	Successful  int `json:"successful"`
	Failed      int `json:"failed"`
	CurrentTask int `json:"currentTask"` // 1-indexed
}

// Progress returns the completion percentage
func (s *BatchStatus) Progress() float64 {
	if s.Total == 0 {
		return 0
	}
	return float64(s.Completed) / float64(s.Total) * 100
}
