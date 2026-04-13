// internal/batch/executor.go
package batch

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/yourname/claude-task-manager/internal/core"
	"github.com/yourname/claude-task-manager/internal/storage"
	"github.com/yourname/claude-task-manager/pkg/models"
)

// Executor handles batch execution of tasks
type Executor struct {
	config     *models.BatchConfig
	taskMgr    *core.TaskManager
	eventStore *storage.EventStore
	status     *models.BatchStatus
	mu         sync.RWMutex
	cancel     context.CancelFunc
}

// NewExecutor creates a new batch executor
func NewExecutor(config *models.BatchConfig, taskMgr *core.TaskManager, eventStore *storage.EventStore) *Executor {
	return &Executor{
		config:     config,
		taskMgr:    taskMgr,
		eventStore: eventStore,
		status: &models.BatchStatus{
			Total: config.Iterations,
		},
	}
}

// Run executes the batch iterations
func (e *Executor) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	e.cancel = cancel
	defer cancel()

	for i := 1; i <= e.config.Iterations; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		e.mu.Lock()
		e.status.CurrentTask = i
		e.mu.Unlock()

		err := e.runIteration(ctx, i)

		e.mu.Lock()
		e.status.Completed++
		if err != nil {
			e.status.Failed++
			if e.config.StopOnError {
				e.mu.Unlock()
				return fmt.Errorf("iteration %d failed: %w", i, err)
			}
		} else {
			e.status.Successful++
		}
		e.mu.Unlock()

		if i < e.config.Iterations && e.config.DelayBetween > 0 {
			time.Sleep(time.Duration(e.config.DelayBetween) * time.Second)
		}
	}

	return nil
}

// runIteration executes a single iteration
func (e *Executor) runIteration(ctx context.Context, iteration int) error {
	taskName := fmt.Sprintf("batch-%03d", iteration)
	prompt := e.generatePrompt(iteration)

	task, err := e.taskMgr.Create(taskName, prompt, e.config.WorkDir, models.IsolationTmux)
	if err != nil {
		return err
	}

	if err := e.taskMgr.Start(task.ID); err != nil {
		return err
	}

	// Wait for completion (simplified - real impl would monitor events)
	time.Sleep(5 * time.Second)

	return nil
}

// generatePrompt creates the prompt for an iteration
func (e *Executor) generatePrompt(iteration int) string {
	return fmt.Sprintf("You are starting iteration %d of %d.\n\nTask: %s\n\nPlease:\n1. Create a new todo item for this task\n2. Implement the required changes\n3. Run tests to verify\n4. Commit your changes with a descriptive message",
		iteration, e.config.Iterations, e.config.Prompt)
}

// Stop cancels the batch execution
func (e *Executor) Stop() {
	if e.cancel != nil {
		e.cancel()
	}
}

// Status returns the current batch status
func (e *Executor) Status() *models.BatchStatus {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.status
}
