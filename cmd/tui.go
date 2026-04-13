// cmd/tui.go
package main

import (
	"github.com/spf13/cobra"
	"github.com/yourname/claude-task-manager/internal/core"
	"github.com/yourname/claude-task-manager/internal/storage"
	"github.com/yourname/claude-task-manager/internal/tui"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the terminal UI",
	Run:   runTUI,
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

func runTUI(cmd *cobra.Command, args []string) {
	taskStore, err := storage.NewTaskStore()
	if err != nil {
		panic(err)
	}

	eventStore, err := storage.NewEventStore()
	if err != nil {
		panic(err)
	}

	bus := core.NewEventBus()
	taskMgr := core.NewTaskManager(taskStore, eventStore, bus)

	if err := tui.Run(taskMgr); err != nil {
		panic(err)
	}
}
