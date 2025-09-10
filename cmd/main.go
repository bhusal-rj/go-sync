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
	hidden      bool
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
	// addFlags(rootCmd)
	addFlags(syncCmd)

	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(infoCmd)
}

func runSync(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(source); os.IsNotExist(err) {
		fmt.Printf("Source path does not exist: %s\n", source)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("Starting sync from %s to %s\n", source, destination)
		fmt.Printf("Recursive: %v\n", recursive)
		fmt.Printf("Verbose: %v\n", verbose)
		fmt.Printf("Hidden: %v\n", hidden)
	}

	opts := sync.SyncOptions{
		Source:      source,
		Destination: destination,
		Recursive:   recursive,
		Verbose:     verbose,
		Hidden:      hidden,
	}

	if err := sync.BasicSync(opts); err != nil {
		fmt.Printf("Sync failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Sync completed successfully!")
}

func runInfo(cmd *cobra.Command, args []string) {
	fmt.Println(("Go Sync is the tool for synchronizing files and directories."))
}

func addFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&source, "source", "s", "", "Source path (local or remote)")
	cmd.Flags().StringVarP(&destination, "destination", "d", "", "Destination path (local or remote)")
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively sync directories")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	if cmd == syncCmd {

		cmd.Flags().BoolVarP(&hidden, "hidden", "H", false, "Include the hidden directories")
	}

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
