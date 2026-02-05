package main

import (
	"fmt"
	"os"

	"skillshare/internal/config"
	"skillshare/internal/sync"
	"skillshare/internal/ui"
)

func cmdSyncProject(args []string, root string) error {
	dryRun := false
	force := false

	for _, arg := range args {
		switch arg {
		case "--dry-run", "-n":
			dryRun = true
		case "--force", "-f":
			force = true
		}
	}

	if !projectConfigExists(root) {
		if err := performProjectInit(root, projectInitOptions{}); err != nil {
			return err
		}
	}

	runtime, err := loadProjectRuntime(root)
	if err != nil {
		return err
	}

	if _, err := os.Stat(runtime.sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("source directory does not exist: %s", runtime.sourcePath)
	}

	if dryRun {
		ui.Warning("Dry run mode - no changes will be made")
	}

	discoveredSkills, discoverErr := sync.DiscoverSourceSkills(runtime.sourcePath)
	if discoverErr == nil {
		collisions := sync.CheckNameCollisions(discoveredSkills)
		if len(collisions) > 0 {
			ui.Header("Name conflicts detected")
			for _, collision := range collisions {
				ui.Warning("Skill name '%s' is defined in multiple places:", collision.Name)
				for _, path := range collision.Paths {
					ui.Info("  - %s", path)
				}
			}
			ui.Info("CLI tools may not distinguish between them.")
			ui.Info("Suggestion: Rename one in SKILL.md (e.g., 'repo:skillname')")
			fmt.Println()
		}
	}

	for _, entry := range runtime.config.Targets {
		name := entry.Name
		target, ok := runtime.targets[name]
		if !ok {
			ui.Error("%s: target not found", name)
			continue
		}

		result, err := sync.SyncTargetMerge(name, target, runtime.sourcePath, dryRun, force)
		if err != nil {
			ui.Error("%s: %v", name, err)
			continue
		}

		if _, err := sync.PruneOrphanLinks(target.Path, runtime.sourcePath, dryRun); err != nil {
			ui.Warning("%s: prune failed: %v", name, err)
		}

		synced := len(result.Linked) + len(result.Updated)
		displayPath := projectTargetDisplayPath(entry)
		fmt.Printf("  %-12s %-16s ✓ %d skills synced\n", name, displayPath, synced)
	}

	return nil
}

func projectTargetDisplayPath(entry config.ProjectTargetEntry) string {
	if entry.Path != "" {
		return entry.Path
	}
	if known, ok := config.LookupProjectTarget(entry.Name); ok {
		return known.Path
	}
	return ""
}
