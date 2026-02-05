package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"skillshare/internal/config"
	"skillshare/internal/ui"
)

type projectInitOptions struct {
	dryRun bool
}

type detectedProjectTarget struct {
	name         string
	path         string
	exists       bool
	parentExists bool
}

func parseProjectInitArgs(args []string) (projectInitOptions, bool, error) {
	opts := projectInitOptions{}
	for _, arg := range args {
		switch arg {
		case "--dry-run", "-n":
			opts.dryRun = true
		case "--help", "-h":
			return opts, true, nil
		default:
			if strings.HasPrefix(arg, "-") {
				return opts, false, fmt.Errorf("unknown option: %s", arg)
			}
		}
	}

	return opts, false, nil
}

func cmdInitProject(args []string) error {
	opts, showHelp, err := parseProjectInitArgs(args)
	if showHelp {
		fmt.Println("Usage: skillshare init -p [--dry-run]")
		return nil
	}
	if err != nil {
		return err
	}

	root, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot determine project directory: %w", err)
	}

	return performProjectInit(root, opts)
}

func performProjectInit(root string, opts projectInitOptions) error {
	projectDir := filepath.Join(root, ".skillshare")
	configPath := config.ProjectConfigPath(root)
	if _, err := os.Stat(projectDir); err == nil {
		return fmt.Errorf("project already initialized: %s", projectDir)
	}

	detected := detectProjectCLIDirectories(root)
	available := detected
	if len(available) == 0 {
		ui.Warning("No AI CLI directories detected.")
		available = listAllProjectTargets()
	}

	selected, err := promptProjectTargets(available)
	if err != nil {
		return err
	}

	if opts.dryRun {
		ui.Header("Dry run complete")
		ui.Info("Would create .skillshare/skills/")
		ui.Info("Would write config: %s", configPath)
		return nil
	}

	if err := os.MkdirAll(filepath.Join(root, ".skillshare", "skills"), 0755); err != nil {
		return fmt.Errorf("failed to create .skillshare/skills: %w", err)
	}

	if err := ensureProjectGitignore(root); err != nil {
		return err
	}

	cfg := &config.ProjectConfig{
		Targets: selected,
	}
	if err := cfg.Save(root); err != nil {
		return err
	}

	if err := createProjectTargetDirs(root, selected); err != nil {
		return err
	}

	ui.Success("Created .skillshare/config.yaml")
	ui.Success("Created .skillshare/skills/")
	ui.Success("Added %d target(s)", len(selected))
	fmt.Println()

	ui.Header("Initialized successfully")
	ui.Success("Source: .skillshare/skills/")
	ui.Success("Config: %s", config.ProjectConfigPath(root))
	fmt.Println()
	ui.Info("Next steps:")
	fmt.Println("  skillshare install <skill> -p    # Install a skill")
	fmt.Println("  skillshare sync                  # Sync to all targets")

	return nil
}

func detectProjectCLIDirectories(root string) []detectedProjectTarget {
	ui.Header("Detecting AI CLI directories")

	projectTargets := config.ProjectTargets()
	names := make([]string, 0, len(projectTargets))
	for name := range projectTargets {
		names = append(names, name)
	}
	sort.Strings(names)

	var detected []detectedProjectTarget
	for _, name := range names {
		target := projectTargets[name]
		relPath := filepath.FromSlash(target.Path)
		fullPath := filepath.Join(root, relPath)

		if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
			ui.Success("Found: %s (%s)", name, relPath)
			detected = append(detected, detectedProjectTarget{name: name, path: relPath, exists: true})
			continue
		}

		parentRel := filepath.Dir(relPath)
		if parentRel != "." {
			parentPath := filepath.Join(root, parentRel)
			if info, err := os.Stat(parentPath); err == nil && info.IsDir() {
				ui.Info("Found: %s (not initialized)", name)
				detected = append(detected, detectedProjectTarget{name: name, path: relPath, parentExists: true})
			}
		}
	}

	return detected
}

func listAllProjectTargets() []detectedProjectTarget {
	projectTargets := config.ProjectTargets()
	names := make([]string, 0, len(projectTargets))
	for name := range projectTargets {
		names = append(names, name)
	}
	sort.Strings(names)

	var available []detectedProjectTarget
	for _, name := range names {
		target := projectTargets[name]
		available = append(available, detectedProjectTarget{name: name, path: filepath.FromSlash(target.Path)})
	}
	return available
}

func promptProjectTargets(available []detectedProjectTarget) ([]config.ProjectTargetEntry, error) {
	ui.Header("Select targets to sync")

	options := make([]string, len(available))
	defaultIndices := []int{}
	for i, target := range available {
		status := ""
		if target.exists {
			status = ""
		} else if target.parentExists {
			status = "(not initialized)"
		}
		if status != "" {
			options[i] = fmt.Sprintf("%-12s %s %s", target.name, target.path, status)
		} else {
			options[i] = fmt.Sprintf("%-12s %s", target.name, target.path)
		}

		if target.exists || target.parentExists {
			defaultIndices = append(defaultIndices, i)
		}
	}

	var selectedIndices []int
	prompt := &survey.MultiSelect{
		Message:  "Select targets to sync:",
		Options:  options,
		Default:  defaultIndices,
		PageSize: 15,
		Help:     "Use arrow keys to navigate, space to select, enter to confirm",
	}

	if err := survey.AskOne(prompt, &selectedIndices); err != nil {
		return nil, nil
	}

	selected := make([]config.ProjectTargetEntry, 0, len(selectedIndices))
	for _, idx := range selectedIndices {
		name := available[idx].name
		selected = append(selected, config.ProjectTargetEntry{Name: name})
	}

	return selected, nil
}

func createProjectTargetDirs(root string, targets []config.ProjectTargetEntry) error {
	if len(targets) == 0 {
		return nil
	}

	knownTargets := config.ProjectTargets()
	created := map[string]bool{}

	for _, target := range targets {
		name := target.Name
		path := target.Path
		if path == "" {
			if known, ok := knownTargets[name]; ok {
				path = known.Path
			} else {
				continue
			}
		}

		absPath := path
		if !filepath.IsAbs(path) {
			absPath = filepath.Join(root, filepath.FromSlash(path))
		}

		if created[absPath] {
			continue
		}
		created[absPath] = true

		if err := os.MkdirAll(absPath, 0755); err != nil {
			return fmt.Errorf("failed to create target directory: %w", err)
		}
	}

	return nil
}

func ensureProjectGitignore(root string) error {
	gitignorePath := filepath.Join(root, ".skillshare", ".gitignore")
	if _, err := os.Stat(gitignorePath); err == nil {
		return nil
	}

	if err := os.WriteFile(gitignorePath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create .skillshare/.gitignore: %w", err)
	}

	return nil
}
