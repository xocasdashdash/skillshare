package main

import (
	"fmt"
	"os"
	"strings"

	"skillshare/internal/config"
	"skillshare/internal/install"
	"skillshare/internal/ui"
	"skillshare/internal/utils"
)

func cmdList(args []string) error {
	var verbose bool

	// Parse arguments
	for _, arg := range args {
		switch arg {
		case "--verbose", "-v":
			verbose = true
		case "--help", "-h":
			printListHelp()
			return nil
		default:
			if strings.HasPrefix(arg, "-") {
				return fmt.Errorf("unknown option: %s", arg)
			}
		}
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Read source directory
	entries, err := os.ReadDir(cfg.Source)
	if err != nil {
		return fmt.Errorf("cannot read source directory: %w", err)
	}

	// Collect skills
	var skills []skillEntry
	for _, e := range entries {
		if !e.IsDir() || utils.IsHidden(e.Name()) {
			continue
		}

		entry := skillEntry{Name: e.Name()}
		skillPath := cfg.Source + "/" + e.Name()

		// Read metadata if available
		if meta, err := install.ReadMeta(skillPath); err == nil && meta != nil {
			entry.Source = meta.Source
			entry.Type = meta.Type
			entry.InstalledAt = meta.InstalledAt.Format("2006-01-02")
		}

		skills = append(skills, entry)
	}

	if len(skills) == 0 {
		ui.Info("No skills installed")
		ui.Info("Use 'skillshare install <source>' to install a skill")
		return nil
	}

	// Display
	ui.Header("Installed skills")
	fmt.Println(strings.Repeat("-", 50))

	for _, s := range skills {
		if verbose {
			fmt.Printf("  %s\n", s.Name)
			if s.Source != "" {
				fmt.Printf("    Source: %s\n", s.Source)
				fmt.Printf("    Type: %s\n", s.Type)
				fmt.Printf("    Installed: %s\n", s.InstalledAt)
			} else {
				fmt.Printf("    Source: (local - no metadata)\n")
			}
			fmt.Println()
		} else if s.Source != "" {
			// Show abbreviated source
			source := abbreviateSource(s.Source)
			fmt.Printf("  %-25s  %s\n", s.Name, source)
		} else {
			fmt.Printf("  %-25s  (local)\n", s.Name)
		}
	}

	if !verbose {
		fmt.Println()
		ui.Info("Use --verbose for more details")
	}

	return nil
}

type skillEntry struct {
	Name        string
	Source      string
	Type        string
	InstalledAt string
}

// abbreviateSource shortens long sources for display
func abbreviateSource(source string) string {
	// Remove https:// prefix
	source = strings.TrimPrefix(source, "https://")
	source = strings.TrimPrefix(source, "http://")

	// Truncate if too long
	if len(source) > 50 {
		return source[:47] + "..."
	}
	return source
}

func printListHelp() {
	fmt.Println(`Usage: skillshare list [options]

List all installed skills in the source directory.

Options:
  --verbose, -v   Show detailed information (source, type, install date)
  --help, -h      Show this help

Examples:
  skillshare list
  skillshare list --verbose`)
}
