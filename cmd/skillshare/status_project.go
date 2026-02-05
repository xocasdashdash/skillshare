package main

import (
	"fmt"
	"os"
	"path/filepath"

	"skillshare/internal/sync"
	"skillshare/internal/utils"
)

func cmdStatusProject(root string) error {
	if !projectConfigExists(root) {
		if err := performProjectInit(root, projectInitOptions{}); err != nil {
			return err
		}
	}

	runtime, err := loadProjectRuntime(root)
	if err != nil {
		return err
	}

	skills, err := sync.DiscoverSourceSkills(runtime.sourcePath)
	if err != nil {
		return fmt.Errorf("cannot discover skills: %w", err)
	}

	fmt.Printf("Project: %s\n", root)
	fmt.Printf("Source: .skillshare/skills/ (%d skills)\n\n", len(skills))
	fmt.Println("Targets:")

	for _, entry := range runtime.config.Targets {
		target, ok := runtime.targets[entry.Name]
		if !ok {
			fmt.Printf("  %-12s %-16s ✗ target not found\n", entry.Name, projectTargetDisplayPath(entry))
			continue
		}

		missing := countMissingProjectSkills(target.Path, runtime.sourcePath, skills)
		displayPath := projectTargetDisplayPath(entry)
		if missing == 0 {
			fmt.Printf("  %-12s %-16s ✓ synced (%d/%d)\n", entry.Name, displayPath, len(skills), len(skills))
		} else {
			fmt.Printf("  %-12s %-16s ✗ %d missing\n", entry.Name, displayPath, missing)
		}
	}

	return nil
}

func countMissingProjectSkills(targetPath, sourcePath string, skills []sync.DiscoveredSkill) int {
	missing := 0
	for _, skill := range skills {
		targetSkillPath := filepath.Join(targetPath, skill.FlatName)
		info, err := os.Lstat(targetSkillPath)
		if err != nil {
			missing++
			continue
		}

		if info.Mode()&os.ModeSymlink == 0 {
			missing++
			continue
		}

		link, err := os.Readlink(targetSkillPath)
		if err != nil {
			missing++
			continue
		}

		absLink := link
		if !filepath.IsAbs(link) {
			absLink = filepath.Join(filepath.Dir(targetSkillPath), link)
		}
		absLink, _ = filepath.Abs(absLink)
		absSource, _ := filepath.Abs(skill.SourcePath)

		if !utils.PathsEqual(absLink, absSource) {
			missing++
		}
	}

	return missing
}
