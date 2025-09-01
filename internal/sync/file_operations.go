package sync

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Path    string
	Size    int64
	IsDir   bool
	ModTime int64
	Mode    os.FileMode
}

var files []FileInfo

func TraverseDirectory(rooPath string) ([]FileInfo, error) {

	// List all the directories within this folder
	dirEntry, err := os.ReadDir(rooPath)

	if err != nil {
		return nil, err
	}
	for _, entry := range dirEntry {
		fileInfo, err := entry.Info()
		if err != nil {
			return nil, err
		}
		files = append(files, FileInfo{
			Path:    filepath.Join(rooPath, entry.Name()),
			Size:    fileInfo.Size(),
			IsDir:   fileInfo.IsDir(),
			ModTime: fileInfo.ModTime().Unix(),
			Mode:    fileInfo.Mode(),
		})
		if !entry.IsDir() {
			fullPath := filepath.Join(rooPath, entry.Name())
			TraverseDirectory(fullPath)
		}
	}
	fmt.Println("Directory entries:", files)
	return nil, nil
}
