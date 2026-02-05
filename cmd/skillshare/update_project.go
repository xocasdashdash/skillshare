package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"skillshare/internal/install"
	"skillshare/internal/ui"
	"skillshare/internal/utils"
)

func cmdUpdateProject(args []string, root string) error {
	var name string
	var updateAll bool
	var dryRun bool

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--all" || arg == "-a":
			updateAll = true
		case arg == "--dry-run" || arg == "-n":
			dryRun = true
		case arg == "--help" || arg == "-h":
			printUpdateHelp()
			return nil
		case strings.HasPrefix(arg, "-"):
			return fmt.Errorf("unknown option: %s", arg)
		default:
			if name != "" {
				return fmt.Errorf("unexpected argument: %s", arg)
			}
			name = arg
		}
	}

	if name == "" && !updateAll {
		updateAll = true
	}

	if !projectConfigExists(root) {
		if err := performProjectInit(root, projectInitOptions{}); err != nil {
			return err
		}
	}

	sourcePath := filepath.Join(root, ".skillshare", "skills")

	if updateAll {
		return updateAllProjectSkills(sourcePath, dryRun)
	}

	return updateSingleProjectSkill(sourcePath, name, dryRun)
}

func updateSingleProjectSkill(sourcePath, name string, dryRun bool) error {
	skillPath := filepath.Join(sourcePath, name)
	if _, err := os.Stat(skillPath); err != nil {
		return fmt.Errorf("skill '%s' not found", name)
	}

	meta, err := install.ReadMeta(skillPath)
	if err != nil || meta == nil {
		return fmt.Errorf("%s is a local skill, nothing to update", name)
	}

	source, err := install.ParseSource(meta.Source)
	if err != nil {
		return fmt.Errorf("invalid source for %s: %w", name, err)
	}

	opts := install.InstallOptions{Force: true, Update: true, DryRun: dryRun}
	if dryRun {
		ui.Info("[dry-run] would update %s", name)
		return nil
	}

	spinner := ui.StartSpinner(fmt.Sprintf("Updating %s...", name))
	if _, err := install.Install(source, skillPath, opts); err != nil {
		spinner.Fail(fmt.Sprintf("%s failed: %v", name, err))
		return nil
	}
	spinner.Success(fmt.Sprintf("Updated %s", name))
	fmt.Println()
	ui.Info("Run 'skillshare sync' to distribute changes")
	return nil
}

func updateAllProjectSkills(sourcePath string, dryRun bool) error {
	entries, err := os.ReadDir(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read project skills: %w", err)
	}

	if dryRun {
		ui.Warning("Dry run mode - no changes will be made")
	}

	updated := 0
	for _, entry := range entries {
		if !entry.IsDir() || utils.IsHidden(entry.Name()) {
			continue
		}

		skillName := entry.Name()
		skillPath := filepath.Join(sourcePath, skillName)
		meta, err := install.ReadMeta(skillPath)
		if err != nil || meta == nil {
			continue
		}

		source, err := install.ParseSource(meta.Source)
		if err != nil {
			ui.Warning("%s invalid source: %v", skillName, err)
			continue
		}

		if dryRun {
			ui.Info("[dry-run] would update %s", skillName)
			continue
		}

		spinner := ui.StartSpinner(fmt.Sprintf("Updating %s...", skillName))
		if _, err := install.Install(source, skillPath, install.InstallOptions{Force: true, Update: true}); err != nil {
			spinner.Fail(fmt.Sprintf("%s failed: %v", skillName, err))
			continue
		}
		spinner.Success(fmt.Sprintf("Updated %s", skillName))
		updated++
	}

	if updated > 0 && !dryRun {
		fmt.Println()
		ui.Info("Run 'skillshare sync' to distribute changes")
	}

	return nil
}
