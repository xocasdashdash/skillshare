//go:build windows

package sync

import (
	"os"
	"os/exec"
	"path/filepath"
)

// createLink creates a directory junction on Windows (no admin required).
// Falls back to symlink if junction fails.
func createLink(targetPath, sourcePath string) error {
	// Ensure absolute paths for junction
	absSource, err := filepath.Abs(sourcePath)
	if err != nil {
		return err
	}

	// Try junction first (no admin required)
	cmd := exec.Command("cmd", "/c", "mklink", "/J", targetPath, absSource)
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Fallback to symlink (requires admin/developer mode)
	return os.Symlink(sourcePath, targetPath)
}

// isJunctionOrSymlink checks if path is a junction or symlink
func isJunctionOrSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	// Both junctions and symlinks have ModeSymlink on Windows
	return info.Mode()&os.ModeSymlink != 0
}
