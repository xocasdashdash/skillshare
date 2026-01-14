package validate

import (
	"strings"
	"testing"
)

func TestTargetName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid lowercase", "claude", false},
		{"valid with numbers", "claude2", false},
		{"valid with hyphen", "my-target", false},
		{"valid with underscore", "my_target", false},
		{"valid uppercase", "Claude", false},
		{"valid mixed case", "MyTarget", false},
		{"empty", "", true},
		{"starts with number", "2claude", true},
		{"reserved word add", "add", true},
		{"reserved word remove", "remove", true},
		{"reserved word list", "list", true},
		{"reserved word all", "all", true},
		{"reserved word case insensitive", "ADD", true},
		{"too long", strings.Repeat("a", 65), true},
		{"special chars at", "my@target", true},
		{"special chars space", "my target", true},
		{"special chars dot", "my.target", true},
		{"only hyphen", "-", true},
		{"starts with hyphen", "-target", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := TargetName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("TargetName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestPath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid absolute", "/home/user/skills", false},
		{"valid relative", "./skills", false},
		{"valid with tilde", "~/skills", false},
		{"empty", "", true},
		{"null byte", "/home/user\x00/skills", true},
		{"too long", "/" + strings.Repeat("a", 4096), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Path(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Path(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestIsLikelySkillsPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"ends with skills", "/home/user/skills", true},
		{"claude skills", "/home/user/.claude/skills", true},
		{"codex skills", "/home/user/.codex/skills", true},
		{"cursor skills", "/home/user/.cursor/skills", true},
		{"gemini skills", "/home/user/.gemini/antigravity/skills", true},
		{"opencode skills", "/home/user/.config/opencode/skills", true},
		{"random directory", "/home/user/documents", false},
		{"contains skills but not ending", "/home/skills/other", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsLikelySkillsPath(tt.input)
			if got != tt.expected {
				t.Errorf("IsLikelySkillsPath(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
