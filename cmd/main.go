package main

import (
	"fmt"
	"os"

	"github.com/bhusal-rj/go-sync/internal/sync"
	"github.com/spf13/cobra"
)

var (
	source      string
	destination string
	recursive   bool
	verbose     bool
)

// Th go-sync is the main entry point of the application
// Usage: go-sync flags
var rootCmd = &cobra.Command{
	Use:   "go-sync",
	Short: "A fast, lightweight file synchronization tool written in Go.",
	Long: `Go Sync is designed to efficiently synchronize files and directories 
between local and remote systems. It uses intelligent algorithms to detect 
changes and transfer only the necessary data, minimizing bandwidth and 
speeding up sync operations.`,
	Run: runSync,
}

// The sync command handles the actual synchronization process
// Usage: go-sync sync flags
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize files and directories",
	Long:  `Synchronize files and directories from source to destination`,
	Run:   runSync,
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display information about the tool",
	Long: `Display detailed information about the go-sync tool, including its features,
	usage instructions, and configuration options.`,
	Run: runInfo,
}

func init() {
	// Go automatically calls the init function before the main function
	addFlags(rootCmd)
	addFlags(syncCmd)

	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(infoCmd)
}

func runSync(cmd *cobra.Command, args []string) {
	sync.TraverseDirectory("../")
	return
	if _, err := os.Stat(source); os.IsNotExist(err) {
		println("Source path does not exist:", source)
		return
	}
	if verbose {
		fmt.Printf("Starting synchronization from %s to %s", source, destination)
		fmt.Printf("Recursive: %v", recursive)
		fmt.Printf("Verbose: %v", verbose)
		fmt.Printf("Starting sync...")
	} else {
		fmt.Printf("Syncing '%s' to '%s'", source, destination)
	}

	// TODO:- Call the actual sync function from the sync package

}

func runInfo(cmd *cobra.Command, args []string) {
	fmt.Println(("Go Sync is the tool for synchronizing files and directories."))
}

func addFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&source, "source", "s", "", "Source path (local or remote)")
	cmd.Flags().StringVarP(&destination, "destination", "d", "", "Destination path (local or remote)")
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively sync directories")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// cmd.MarkFlagRequired("source")
	// cmd.MarkFlagRequired("destination")
}

func main() {
	fmt.Println("Starting application...")
	if err := rootCmd.Execute(); err != nil {
		println("Error executing command:", err.Error())
		os.Exit(1)
	}
}
