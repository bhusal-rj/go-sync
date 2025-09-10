package sync

import (
	"fmt"
	"os"
	"path/filepath"
)

type SyncOptions struct {
	Source      string //The path of the file or directory to sync
	Destination string //Destination path
	Recursive   bool   //Whether to sync directory recursively
	Verbose     bool   //Enable verbose output
	Hidden      bool   // Inclue hidden files and directories
}

// BasicSync performs basic file/directory sync
func BasicSync(opts SyncOptions) error {
	sourceInfo, err := GetFileInfo(opts.Source)

	if err != nil {
		return fmt.Errorf("source path error: %w", err)
	}

	if sourceInfo.IsDir {
		return syncDirectory(opts.Source, opts.Destination, opts)
	} else {
		return syncFile(opts.Source, opts.Destination, opts)
	}
}

func syncDirectory(source, destination string, opts SyncOptions) error {
	if opts.Verbose {
		fmt.Printf("Syncing directory: %s -> %s", source, destination)
	}

	// Create the destination directory
	if err := CreateDirectory(destination); err != nil {
		return err
	}

	entries, err := os.ReadDir(source)

	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(source, entry.Name())
		dstPath := filepath.Join(destination, entry.Name())

		if entry.IsDir() {
			if opts.Recursive {
				if err := syncDirectory(srcPath, dstPath, opts); err != nil {
					return err
				}
			}
		} else {
			if err := syncFile(srcPath, dstPath, opts); err != nil {
				return nil
			}
		}
	}

	// Preserve directory metadata
	if err := PreserveDirectoryMetadata(source, destination); err != nil {
		if opts.Verbose {
			fmt.Printf("Warning: Could not preserve directory metadata: %v\n", err)
		}
	}

	return nil
}

func syncFile(source, destination string, opts SyncOptions) error {
	if opts.Verbose {
		fmt.Printf("Syncing file: %s -> %s\n", source, destination)
	}

	return CopyFile(source, destination)
}
