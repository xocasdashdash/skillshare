package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"skillshare/internal/config"
	"skillshare/internal/install"
	"skillshare/internal/ui"
)

func cmdUninstall(args []string) error {
	var skillName string
	var force, dryRun bool

	// Parse arguments
	i := 0
	for i < len(args) {
		arg := args[i]
		switch {
		case arg == "--force" || arg == "-f":
			force = true
		case arg == "--dry-run" || arg == "-n":
			dryRun = true
		case arg == "--help" || arg == "-h":
			printUninstallHelp()
			return nil
		case strings.HasPrefix(arg, "-"):
			return fmt.Errorf("unknown option: %s", arg)
		default:
			if skillName != "" {
				return fmt.Errorf("unexpected argument: %s", arg)
			}
			skillName = arg
		}
		i++
	}

	if skillName == "" {
		printUninstallHelp()
		return fmt.Errorf("skill name is required")
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if skill exists
	skillPath := filepath.Join(cfg.Source, skillName)
	info, err := os.Stat(skillPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("skill '%s' not found in source", skillName)
		}
		return fmt.Errorf("cannot access skill: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("'%s' is not a directory", skillName)
	}

	// Display info
	ui.Header("Uninstalling skill")
	fmt.Println(strings.Repeat("-", 45))
	ui.Info("Skill: %s", skillName)
	ui.Info("Path: %s", skillPath)

	// Show metadata if available
	if meta, err := install.ReadMeta(skillPath); err == nil && meta != nil {
		ui.Info("Source: %s", meta.Source)
		ui.Info("Installed: %s", meta.InstalledAt.Format("2006-01-02 15:04"))
	}
	fmt.Println()

	if dryRun {
		ui.Warning("[dry-run] would remove %s", skillPath)
		return nil
	}

	// Confirm unless --force
	if !force {
		fmt.Print("Are you sure you want to uninstall this skill? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		input = strings.TrimSpace(strings.ToLower(input))
		if input != "y" && input != "yes" {
			ui.Info("Cancelled")
			return nil
		}
	}

	// Remove the skill
	if err := os.RemoveAll(skillPath); err != nil {
		return fmt.Errorf("failed to remove skill: %w", err)
	}

	ui.Success("Uninstalled: %s", skillName)
	fmt.Println()
	ui.Info("Run 'skillshare sync' to update all targets")

	return nil
}

func printUninstallHelp() {
	fmt.Println(`Usage: skillshare uninstall <skill-name> [options]

Remove a skill from the source directory.

Options:
  --force, -f     Skip confirmation prompt
  --dry-run, -n   Preview without making changes
  --help, -h      Show this help

Examples:
  skillshare uninstall my-skill
  skillshare uninstall my-skill --force
  skillshare uninstall my-skill --dry-run`)
}
