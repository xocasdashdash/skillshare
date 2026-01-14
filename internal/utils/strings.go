package utils

// IsHidden checks if a file/directory name starts with a dot.
// Returns false for empty strings to prevent panic on empty names.
func IsHidden(name string) bool {
	return len(name) > 0 && name[0] == '.'
}

// HasTildePrefix checks if a path starts with ~.
// Returns false for empty strings to prevent panic on empty paths.
func HasTildePrefix(path string) bool {
	return len(path) > 0 && path[0] == '~'
}
