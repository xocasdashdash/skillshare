package backup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// BackupDir returns the backup directory path
func BackupDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "skillshare", "backups")
}

// Create creates a backup of the target directory
// Returns the backup path
func Create(targetName, targetPath string) (string, error) {
	// Check if target exists and has content
	info, err := os.Lstat(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // Nothing to backup
		}
		return "", err
	}

	// Skip if it's already a symlink (no local data to backup)
	if info.Mode()&os.ModeSymlink != 0 {
		return "", nil
	}

	// Check if directory has any content
	entries, err := os.ReadDir(targetPath)
	if err != nil || len(entries) == 0 {
		return "", nil // Empty, nothing to backup
	}

	// Create backup directory with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupPath := filepath.Join(BackupDir(), timestamp, targetName)

	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Copy target contents to backup
	if err := copyDir(targetPath, backupPath); err != nil {
		return "", fmt.Errorf("failed to backup: %w", err)
	}

	return backupPath, nil
}

// List returns all backups sorted by date (newest first)
func List() ([]BackupInfo, error) {
	backupDir := BackupDir()

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var backups []BackupInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		backupPath := filepath.Join(backupDir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// List targets in this backup
		targetEntries, _ := os.ReadDir(backupPath)
		var targets []string
		for _, t := range targetEntries {
			if t.IsDir() {
				targets = append(targets, t.Name())
			}
		}

		backups = append(backups, BackupInfo{
			Timestamp: entry.Name(),
			Path:      backupPath,
			Targets:   targets,
			Date:      info.ModTime(),
		})
	}

	// Sort by date (newest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Date.After(backups[j].Date)
	})

	return backups, nil
}

// BackupInfo holds information about a backup
type BackupInfo struct {
	Timestamp string
	Path      string
	Targets   []string
	Date      time.Time
}

// copyDir copies a directory recursively
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip symlinks - we only backup real files
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		relPath, _ := filepath.Rel(src, path)
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath)
	})
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
