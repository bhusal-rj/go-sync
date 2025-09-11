# Go Sync

A fast, lightweight file synchronization tool written in Go. Supports local and remote syncing with intelligent delta updates and minimal bandwidth usage.

## Project Description
Go Sync is designed to efficiently synchronize files and directories between local and remote systems. It uses intelligent algorithms to detect changes and transfer only the necessary data, minimizing bandwidth and speeding up sync operations. The tool is suitable for backups, mirroring, and collaborative workflows, and aims to be easy to use, reliable, and extensible.

## Milestones

- **M1:** Basic file/directory sync (no delta)
  - [x] Design file/directory traversal
  - [x] Implement file copy/move logic
  - [x] Handle directory recursion
  - [x] Basic error handling

- **M2:** Preserve metadata (timestamps, permissions)
  - [x] Research metadata handling in Go
  - [x] Implement timestamp preservation
  - [x] Implement permission preservation
  - [x] Test metadata sync on all platforms

- **M3:** Checksum-based delta sync
  - [x] Select checksum algorithm
  - [x] Implement checksum calculation
  - [x] Integrate delta sync logic
  - [x] Test delta sync performance

- **M4:** Compression (optional gzip/zstd)
  - [ ] Research compression libraries
  - [ ] Implement gzip compression
  - [ ] Implement zstd compression
  - [ ] Add CLI/config option for compression

- **M5:** Remote sync over SSH or TCP
  - [ ] Integrate SSH/TCP libraries
  - [ ] Implement remote authentication
  - [ ] Implement remote file transfer
  - [ ] Test remote sync reliability

- **M6:** Progress bar, verbose logs
  - [ ] Design progress bar UI
  - [ ] Implement progress tracking
  - [ ] Add verbose logging option
  - [ ] Test logging and progress bar

- **M7:** Concurrent file transfers
  - [ ] Design concurrency model
  - [ ] Implement concurrent transfers
  - [ ] Handle concurrency errors
  - [ ] Benchmark concurrency performance

---
