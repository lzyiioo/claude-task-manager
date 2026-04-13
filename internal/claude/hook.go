package claude

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type HookConfig struct {
	PreToolUse   string `json:"pre_tool_use"`
	PostToolUse  string `json:"post_tool_use"`
	Notification string `json:"notification"`
}

type HookHandler struct {
	taskID     string
	socketPath string
}

func NewHookHandler(taskID string) *HookHandler {
	socketPath := filepath.Join(os.TempDir(), fmt.Sprintf("ctm-%s.sock", taskID))
	return &HookHandler{
		taskID:     taskID,
		socketPath: socketPath,
	}
}

func (h *HookHandler) GenerateHookScript() string {
	return fmt.Sprintf(`#!/bin/sh
# Claude Task Manager Hook for task %s
exec nc -U %s
`, h.taskID, h.socketPath)
}

func (h *HookHandler) GetEnvVars() map[string]string {
	hookScript := h.GenerateHookScript()
	hookPath := filepath.Join(os.TempDir(), fmt.Sprintf("ctm-hook-%s.sh", h.taskID))

	os.WriteFile(hookPath, []byte(hookScript), 0755)

	hooks := HookConfig{
		PreToolUse:   hookPath,
		PostToolUse:  hookPath,
		Notification: hookPath,
	}
	hooksJSON, _ := json.Marshal(hooks)

	return map[string]string{
		"CLAUDE_HOOKS": string(hooksJSON),
	}
}
