// cmd/main.go
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "ctm",
	Short: "Claude Task Manager - Multi-task execution platform for Claude Code",
	Long: `Claude Task Manager (ctm) is a terminal-based platform for running
multiple Claude Code instances in parallel with real-time monitoring,
permission control, and batch automation.`,
	Version: version,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
