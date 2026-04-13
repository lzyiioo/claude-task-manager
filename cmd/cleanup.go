// cmd/cleanup.go
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourname/claude-task-manager/internal/core"
	"github.com/yourname/claude-task-manager/internal/isolation"
	"github.com/yourname/claude-task-manager/internal/storage"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean up isolation resources",
	Long:  `Stop and remove all active isolation sessions (tmux/docker).`,
	Run:   runCleanup,
}

var cleanupTaskID string

func init() {
	rootCmd.AddCommand(cleanupCmd)
	cleanupCmd.Flags().StringVarP(&cleanupTaskID, "task", "t", "", "cleanup specific task only")
}

func runCleanup(cmd *cobra.Command, args []string) {
	taskStore, err := storage.NewTaskStore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating task store: %v\n", err)
		os.Exit(1)
	}

	eventStore, err := storage.NewEventStore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating event store: %v\n", err)
		os.Exit(1)
	}

	bus := core.NewEventBus()
	taskMgr := core.NewTaskManager(taskStore, eventStore, bus)

	// Register available drivers
	taskMgr.RegisterDriver(isolation.NewTmuxDriver())
	taskMgr.RegisterDriver(isolation.NewDockerDriver("claude-code:latest"))

	if cleanupTaskID != "" {
		// Cleanup specific task
		if err := taskMgr.CleanupTask(cleanupTaskID); err != nil {
			fmt.Fprintf(os.Stderr, "Error cleaning up task %s: %v\n", cleanupTaskID, err)
			os.Exit(1)
		}
		fmt.Printf("Cleaned up task %s\n", cleanupTaskID)
	} else {
		// Cleanup all
		taskMgr.CleanupAll()
		fmt.Println("Cleaned up all isolation resources")
	}
}
