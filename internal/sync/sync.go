package sync

import (
	"crypto/md5"
	"fmt"

	"os"
	"path"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SyncOptions struct {
	Source      string //The path of the file or directory to sync
	Destination string //Destination path
	Recursive   bool   //Whether to sync directory recursively
	Verbose     bool   //Enable verbose output
	Hidden      bool   // Inclue hidden files and directories
	Username    string // Username for remote connection
	Host        string // Host for remote connection
	SSHKey      string // Path to SSH private key for remote connections
}

var SftpClient *sftp.Client
var IsServerSync bool
var SSHConn *ssh.Client

// Server Sync
func GetSFTPClient(opts SyncOptions) error {
	// Connect to the server via ssh

	key, err := os.ReadFile(opts.SSHKey)
	if err != nil {
		return fmt.Errorf("failed to read SSH key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("failed to parse SSH key: %w", err)
	}

	config := &ssh.ClientConfig{
		User: opts.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Establish the connection
	sshConn, err := ssh.Dial("tcp", opts.Host+":22", config)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", opts.Host, err)
	}
	fmt.Println("SSH connection established")
	// defer sshConn.Close()

	sftpClient, err := sftp.NewClient(sshConn)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}
	SSHConn = sshConn
	// defer sftpClient.Close()

	SftpClient = sftpClient

	fmt.Println("SFTP client created")
	return nil
}

// BasicSync performs basic file/directory sync
func BasicSync(opts SyncOptions) error {

	sourceInfo, err := GetFileInfo(opts.Source)

	if err != nil {
		return fmt.Errorf("source path error: %w", err)
	}

	IsServerSync = opts.Host != "" && opts.Username != "" && opts.SSHKey != ""
	if IsServerSync {
		err := GetSFTPClient(opts)
		if err != nil {
			return fmt.Errorf("failed to get SFTP client: %w", err)
		}

		// return nil
	}

	if sourceInfo.IsDir {
		syncDirectory(opts.Source, opts.Destination, opts)
	} else {

		syncFile(opts.Source, opts.Destination, opts)
	}

	fmt.Println("Sync completed successfully")
	return nil
}

func syncDirectory(source, destination string, opts SyncOptions) error {
	if opts.Verbose {
		fmt.Printf("Syncing directory: %s -> %s", source, destination)
	}

	if IsServerSync {
		if err := CreateServerDirectory(SftpClient, destination); err != nil {
			return err
		}
	} else {
		// Create the destination directory
		if err := CreateDirectory(destination); err != nil {
			return err
		}

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

func CalculateFileDelta(destination_page os.FileInfo, path string) (string, error) {
	// Check if the destination exist
	// destination_page, err := os.Stat(path)

	// if err != nil {
	// 	return "", nil
	// }
	// If the destination exists check the delta of the file

	fileInfo := &FileInfo{
		Path:    path,
		Size:    destination_page.Size(),
		IsDir:   destination_page.IsDir(),
		ModTime: destination_page.ModTime(),
		Mode:    destination_page.Mode(),
	}

	sys := destination_page.Sys()

	switch stat := sys.(type) {
	case *syscall.Stat_t:
		fileInfo.Uid = int(stat.Uid)
		fileInfo.Gid = int(stat.Gid)
	case *sftp.FileStat:
		fileInfo.Uid = int(stat.UID)
		fileInfo.Gid = int(stat.GID)
	default:
		fileInfo.Uid = 0
		fileInfo.Gid = 0
	}
	// Calculate the checksum of the file
	checksum := CalculateFileChecksum(*fileInfo)
	return checksum, nil
}

func getServerFileStat(path string) (os.FileInfo, error) {
	fileStat, err := SftpClient.Stat(path)
	if err != nil {
		return nil, err
	}
	return fileStat, nil
}
func syncFile(source, destination string, opts SyncOptions) error {

	if opts.Verbose {
		fmt.Printf("Syncing file: %s -> %s\n", source, destination)
	}

	var sourceChecksum, destinationChecksum string
	var destinationPage os.FileInfo
	// Get the source calculatefilechecksum

	sourcePage, err := os.Stat(source)

	if err == nil {
		sourceChecksum, _ = CalculateFileDelta(sourcePage, source)
	}

	if IsServerSync {
		destinationPage, err = getServerFileStat(path.Join(destination, path.Base(source)))

	} else {
		destinationPage, err = os.Stat(destination)
	}

	if err == nil {
		destinationChecksum, _ = CalculateFileDelta(destinationPage, path.Join(destination, source))
	}
	// Check the checksum of the source and destination

	if err == nil && sourceChecksum == destinationChecksum {
		if opts.Verbose {
			fmt.Printf("File is up to date, skipping: %s\n", source)
		}
		return nil
	}

	return CopyFile(source, destination)
}
