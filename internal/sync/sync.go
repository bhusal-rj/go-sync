package sync

import (
	"crypto/md5"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"syscall"
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

func CalculateFileChecksum(fileInfo FileInfo) string {
	// Placeholder for checksum calculation logic
	fileCheckSumPath := fileInfo.ModTime.String() +
		strconv.FormatInt(fileInfo.Size, 10)

	byteString := fmt.Sprintf("%x", md5.Sum([]byte(fileCheckSumPath)))

	return byteString
}

func CalculateFileDelta(path string) (string, error) {
	// Check if the destination exist
	destination_page, err := os.Stat(path)

	if err != nil {
		return "", nil
	}
	// If the destination exists check the delta of the file

	fileInfo := &FileInfo{
		Path:    path,
		Size:    destination_page.Size(),
		IsDir:   destination_page.IsDir(),
		ModTime: destination_page.ModTime(),
		Mode:    destination_page.Mode(),
		Uid:     int(destination_page.Sys().(*syscall.Stat_t).Uid),
		Gid:     int(destination_page.Sys().(*syscall.Stat_t).Gid),
	}

	// Calculate the checksum of the file
	checksum := CalculateFileChecksum(*fileInfo)
	return checksum, nil
}

func syncFile(source, destination string, opts SyncOptions) error {

	if opts.Verbose {
		fmt.Printf("Syncing file: %s -> %s\n", source, destination)
	}

	// Get the source calculatefilechecksum
	sourceChecksum, _ := CalculateFileDelta(source)
	destinationChecksum, err := CalculateFileDelta(path.Join(destination, source))

	// Check the checksum of the source and destination

	if err == nil && sourceChecksum == destinationChecksum {
		if opts.Verbose {
			fmt.Printf("File is up to date, skipping: %s\n", source)
		}
		return nil
	}

	return CopyFile(source, destination)
}
