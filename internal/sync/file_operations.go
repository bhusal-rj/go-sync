package sync

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/pkg/sftp"
)

type FileInfo struct {
	Path    string
	Size    int64
	IsDir   bool
	ModTime time.Time
	Mode    os.FileMode
	Uid     int
	Gid     int
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
		ModTime: fileInfo.ModTime(),
		Mode:    fileInfo.Mode(),
	}, nil

}

func CreateDirectory(path string) error {
	//Create the directory with the help of os command
	err := os.Mkdir(path, 07555)
	return err
}

func CreateServerDirectory(sftpClient *sftp.Client, path string) error {
	_, err := sftpClient.ReadDir(path)
	if err != nil {
		//Check the error is the os Doesnot exist error
		if os.IsNotExist(err) {
			// Create the directory
			err = sftpClient.MkdirAll(path)
			if err != nil {
				return fmt.Errorf("failed to create remote directory: %w", err)
			}
			fmt.Println("Remote directory created:", path)
			return nil
		} else {
			return fmt.Errorf("failed to read remote directory: %w", err)
		}
	}
	return nil
}

func CopyServerFile(parentDir string, sourceFile *os.File) error {
	err := SftpClient.MkdirAll(parentDir)
	if err != nil {
		return err
	}

	destFile, err := SftpClient.Create(parentDir)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)

	if err != nil {
		return err
	}
	return nil
}

func CopyLocalFile(path string, sourceFile *os.File) error {
	// Create the parent directory
	parentDir := filepath.Dir(path)

	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return err
	}
	//Create the destination file from the disk
	destFile, err := os.Create(path)

	if err != nil {
		return err
	}
	defer destFile.Close()

	//Copy the info from the source to the destination
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}
	return nil
}
func CopyFile(source string, destination string) error {

	//Open the source file from the disk
	sourceFile, err := os.Open(source)

	if err != nil {
		return err
	}
	defer sourceFile.Close()
	destinationPath := path.Join(destination, source)

	if IsServerSync {
		err = CopyServerFile(destinationPath, sourceFile)
	} else {
		err = CopyLocalFile(destinationPath, sourceFile)
	}

	//Ensure that all the data is written to the file

	return PreserveFileMetadata(source, destinationPath)
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
			ModTime: fileInfo.ModTime(),
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

// PreserveFileMetadata and permissions of the copied file
func PreserveFileMetadata(source, destination string) error {

	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// Preser the Server File Metadata
	if IsServerSync {
		// Preserve the permission
		if err := SftpClient.Chmod(destination, sourceInfo.Mode()); err != nil {
			return nil
		}

		// Preserve the ownership of the file and guid and uid
		if stat, ok := sourceInfo.Sys().(*syscall.Stat_t); ok {
			if err := SftpClient.Chown(destination, int(stat.Uid), int(stat.Gid)); err != nil {
				return nil
			}
		}
		return nil
	}
	// Preserve modification time
	if err := os.Chtimes(destination, time.Now(), sourceInfo.ModTime()); err != nil {
		return nil
	}

	// // Preserve the ownership of the file and guid and uid
	// if stat, ok := sourceInfo.Sys().(*syscall.Stat_t); ok {
	// 	if err := os.Chown(destination, int(stat.Uid), int(stat.Gid)); err != nil {
	// 		return nil
	// 	}
	// }
	return nil

}

func PreserveDirectoryMetadata(source, destination string) error {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return nil
	}

	// Preserve the permission
	if stat, ok := sourceInfo.Sys().(*syscall.Stat_t); ok {
		if err := os.Chown(destination, int(stat.Uid), int(stat.Gid)); err != nil {
			return nil
		}
	}
	return nil

}
