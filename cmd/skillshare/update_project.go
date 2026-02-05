package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"skillshare/internal/git"
	"skillshare/internal/install"
	"skillshare/internal/ui"
	"skillshare/internal/utils"
)

func cmdUpdateProject(args []string, root string) error {
	var name string
	var updateAll bool
	var dryRun bool
	var force bool

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--all" || arg == "-a":
			updateAll = true
		case arg == "--dry-run" || arg == "-n":
			dryRun = true
		case arg == "--force" || arg == "-f":
			force = true
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
		return updateAllProjectSkills(sourcePath, dryRun, force)
	}

	return updateSingleProjectSkill(sourcePath, name, dryRun, force)
}

func updateSingleProjectSkill(sourcePath, name string, dryRun, force bool) error {
	// Normalize _ prefix for tracked repos
	repoName := name
	if !strings.HasPrefix(repoName, "_") {
		prefixed := filepath.Join(sourcePath, "_"+name)
		if install.IsGitRepo(prefixed) {
			repoName = "_" + name
		}
	}
	repoPath := filepath.Join(sourcePath, repoName)

	// Try as tracked repo first
	if install.IsGitRepo(repoPath) {
		return updateProjectTrackedRepo(repoName, repoPath, dryRun, force)
	}

	// Regular skill with metadata
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

	if dryRun {
		ui.Info("[dry-run] would update %s", name)
		return nil
	}

	spinner := ui.StartSpinner(fmt.Sprintf("Updating %s...", name))
	opts := install.InstallOptions{Force: true, Update: true}
	if _, err := install.Install(source, skillPath, opts); err != nil {
		spinner.Fail(fmt.Sprintf("%s failed: %v", name, err))
		return nil
	}
	spinner.Success(fmt.Sprintf("Updated %s", name))
	fmt.Println()
	ui.Info("Run 'skillshare sync' to distribute changes")
	return nil
}

func updateProjectTrackedRepo(repoName, repoPath string, dryRun, force bool) error {
	// Check for uncommitted changes
	if isDirty, _ := git.IsDirty(repoPath); isDirty {
		if !force {
			ui.Warning("%s has uncommitted changes (use --force to discard)", repoName)
			return fmt.Errorf("uncommitted changes in %s", repoName)
		}
		if !dryRun {
			if err := git.Restore(repoPath); err != nil {
				return fmt.Errorf("failed to discard changes: %w", err)
			}
		}
	}

	if dryRun {
		ui.Info("[dry-run] would git pull %s", repoName)
		return nil
	}

	spinner := ui.StartSpinner(fmt.Sprintf("Updating %s...", repoName))

	var info *git.UpdateInfo
	var err error
	if force {
		info, err = git.ForcePull(repoPath)
	} else {
		info, err = git.Pull(repoPath)
	}
	if err != nil {
		spinner.Fail(fmt.Sprintf("%s failed: %v", repoName, err))
		return nil
	}

	if info.UpToDate {
		spinner.Success(fmt.Sprintf("%s already up to date", repoName))
	} else {
		spinner.Success(fmt.Sprintf("%s %d commits, %d files", repoName, len(info.Commits), info.Stats.FilesChanged))
	}
	fmt.Println()
	ui.Info("Run 'skillshare sync' to distribute changes")
	return nil
}

func updateAllProjectSkills(sourcePath string, dryRun, force bool) error {
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

		// Tracked repo: git pull
		if install.IsGitRepo(skillPath) {
			if err := updateProjectTrackedRepo(skillName, skillPath, dryRun, force); err != nil {
				ui.Warning("%s: %v", skillName, err)
			} else {
				updated++
			}
			continue
		}

		// Regular skill with metadata: reinstall
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
