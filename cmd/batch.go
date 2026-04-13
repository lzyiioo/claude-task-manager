// cmd/batch.go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/yourname/claude-task-manager/internal/batch"
	"github.com/yourname/claude-task-manager/internal/core"
	"github.com/yourname/claude-task-manager/internal/storage"
	"github.com/yourname/claude-task-manager/pkg/models"
)

var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "Run batch execution of tasks",
	Long:  `Run multiple iterations of Claude Code tasks automatically.`,
}

var batchRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Start batch execution",
	Run:   runBatch,
}

var (
	batchConfig     string
	batchIterations int
	batchPrompt     string
	batchPromptFile string
	batchWorkDir    string
	batchAutoPerm   bool
	batchDryRun     bool
)

func init() {
	rootCmd.AddCommand(batchCmd)
	batchCmd.AddCommand(batchRunCmd)

	batchRunCmd.Flags().StringVarP(&batchConfig, "config", "c", "", "config file path")
	batchRunCmd.Flags().IntVarP(&batchIterations, "iterations", "n", 1, "number of iterations")
	batchRunCmd.Flags().StringVarP(&batchPrompt, "prompt", "p", "", "prompt for each iteration")
	batchRunCmd.Flags().StringVarP(&batchPromptFile, "prompt-file", "f", "", "file containing prompt")
	batchRunCmd.Flags().StringVarP(&batchWorkDir, "workdir", "w", ".", "working directory")
	batchRunCmd.Flags().BoolVarP(&batchAutoPerm, "auto-permission", "a", false, "auto-allow all permissions")
	batchRunCmd.Flags().BoolVarP(&batchDryRun, "dry-run", "d", false, "show what would be executed")
}

func runBatch(cmd *cobra.Command, args []string) {
	var config *models.BatchConfig

	if batchConfig != "" {
		var err error
		config, err = batch.LoadConfig(batchConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
	} else {
		config = batch.DefaultConfig()
		if batchIterations > 0 {
			config.Iterations = batchIterations
		}
		if batchPrompt != "" {
			config.Prompt = batchPrompt
		}
		if batchPromptFile != "" {
			config.PromptFile = batchPromptFile
		}
		config.WorkDir = batchWorkDir
		if batchAutoPerm {
			config.PermissionMode = models.PermissionAuto
		}
	}

	if config.Prompt == "" && config.PromptFile == "" {
		fmt.Fprintln(os.Stderr, "Error: --prompt or --prompt-file required")
		os.Exit(1)
	}

	if batchDryRun {
		fmt.Printf("Would run %d iterations with prompt: %s\n", config.Iterations, config.Prompt)
		return
	}

	// Initialize dependencies
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

	executor := batch.NewExecutor(config, taskMgr, eventStore)

	// Handle interrupts
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nStopping batch execution...")
		executor.Stop()
		cancel()
	}()

	fmt.Printf("Starting batch execution: %d iterations\n", config.Iterations)

	if err := executor.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Batch execution failed: %v\n", err)
		os.Exit(1)
	}

	status := executor.Status()
	fmt.Printf("\nBatch complete: %d/%d successful\n", status.Successful, status.Total)
}
