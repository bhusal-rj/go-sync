package sync

import (
	"io"
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

func GetFileInfo(path string) (*FileInfo, error) {
	// Get the absolute path if the path is not available
	file_path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// Get the relevant info about the file
	fileInfo, err := os.Stat(file_path)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Path:    file_path,
		Size:    fileInfo.Size(),
		IsDir:   fileInfo.IsDir(),
		ModTime: fileInfo.ModTime().Unix(),
		Mode:    fileInfo.Mode(),
	}, nil

}

func CreateDirectory(path string) error {
	//Create the directory with the help of os command
	err := os.Mkdir(path, 07555)
	return err
}

func CopyFile(source string, destination string) error {

	//Open the source file from the disk
	sourceFile, err := os.Open(source)

	if err != nil {
		return err
	}
	defer sourceFile.Close()

	//Create the destination file from the disk
	destFile, err := os.Create(destination)

	if err != nil {
		return err
	}
	defer destFile.Close()

	//Copy the info from the source to the destination
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	//Ensure that all the data is written to the file
	return destFile.Sync()
}
func TraverseDirectory(rooPath string, hidden bool) ([]FileInfo, error) {

	var files []FileInfo
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
		if entry.IsDir() {
			if !hidden && entry.Name()[0] == '.' {
				continue
			}
			fullPath := filepath.Join(rooPath, entry.Name())
			filesFound, err := TraverseDirectory(fullPath, hidden)
			if err != nil {
				return nil, err
			}
			files = append(files, filesFound...)
		}
	}

	return files, nil
}
