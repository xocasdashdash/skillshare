package main

import (
	"fmt"
	"os"
	"path/filepath"

	"skillshare/internal/config"
	"skillshare/internal/ui"
)

func cmdUpdate(args []string) error {
	dryRun := false
	force := false

	// Parse args
	for _, arg := range args {
		switch arg {
		case "--dry-run", "-n":
			dryRun = true
		case "--force", "-f":
			force = true
		}
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config not found: run 'skillshare init' first")
	}

	skillshareSkillDir := filepath.Join(cfg.Source, "skillshare")
	skillshareSkillFile := filepath.Join(skillshareSkillDir, "SKILL.md")

	ui.Header("Updating skillshare skill")

	// Check if skill exists
	exists := false
	if _, err := os.Stat(skillshareSkillFile); err == nil {
		exists = true
	}

	if dryRun {
		ui.Warning("Dry run mode - no changes will be made")
		fmt.Println()
		if exists {
			ui.Info("Would update: %s", skillshareSkillFile)
		} else {
			ui.Info("Would download: %s", skillshareSkillFile)
		}
		ui.Info("Source: %s", skillshareSkillURL)
		return nil
	}

	// Confirm if exists and not forced
	if exists && !force {
		fmt.Print("  Overwrite existing skillshare skill? [y/N]: ")
		var input string
		fmt.Scanln(&input)
		if input != "y" && input != "Y" && input != "yes" {
			ui.Info("Cancelled")
			return nil
		}
	}

	// Create directory if needed
	if err := os.MkdirAll(skillshareSkillDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Download
	ui.Info("Downloading from GitHub...")
	if err := downloadSkillshareSkill(skillshareSkillFile); err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	ui.Success("Updated skillshare skill")
	ui.Info("Path: %s", skillshareSkillFile)
	ui.Info("")
	ui.Info("Run 'skillshare sync' to distribute to all targets")

	return nil
}
