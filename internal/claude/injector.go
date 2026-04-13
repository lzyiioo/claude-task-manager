package claude

import (
	"fmt"
	"io"
	"os/exec"
	"sync"
)

type InputInjector struct {
	mu     sync.Mutex
	stdin  io.WriteCloser
	taskID string
}

func NewInputInjector(taskID string) *InputInjector {
	return &InputInjector{taskID: taskID}
}

func (i *InputInjector) Attach(cmd *exec.Cmd) error {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	i.stdin = stdin
	return nil
}

func (i *InputInjector) Send(input string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.stdin == nil {
		return fmt.Errorf("not attached to process")
	}

	_, err := i.stdin.Write([]byte(input + "\n"))
	return err
}

func (i *InputInjector) SendPermissionResponse(allow bool, reason string) error {
	response := "n"
	if allow {
		response = "y"
	}
	return i.Send(response)
}

func (i *InputInjector) Close() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.stdin != nil {
		return i.stdin.Close()
	}
	return nil
}
