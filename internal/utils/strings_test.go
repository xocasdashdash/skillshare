package utils

import "testing"

func TestIsHidden(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty string", "", false},
		{"hidden file", ".hidden", true},
		{"hidden directory", ".git", true},
		{"normal file", "file.txt", false},
		{"normal directory", "src", false},
		{"dot only", ".", true},
		{"double dot", "..", true},
		{"file starting with number", "1file", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsHidden(tt.input)
			if got != tt.expected {
				t.Errorf("IsHidden(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestHasTildePrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty string", "", false},
		{"tilde path", "~/Documents", true},
		{"tilde only", "~", true},
		{"absolute path", "/home/user", false},
		{"relative path", "./config", false},
		{"tilde in middle", "/home/~user", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasTildePrefix(tt.input)
			if got != tt.expected {
				t.Errorf("HasTildePrefix(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
