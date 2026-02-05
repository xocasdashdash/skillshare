package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"skillshare/internal/install"
	"skillshare/internal/utils"
)

type projectSkillEntry struct {
	Name   string
	Source string
	Remote bool
}

func cmdListProject(root string) error {
	if !projectConfigExists(root) {
		if err := performProjectInit(root, projectInitOptions{}); err != nil {
			return err
		}
	}

	sourcePath := filepath.Join(root, ".skillshare", "skills")
	entries, err := os.ReadDir(sourcePath)
	if err != nil {
		return fmt.Errorf("cannot read project skills: %w", err)
	}

	var skills []projectSkillEntry
	for _, entry := range entries {
		if !entry.IsDir() || utils.IsHidden(entry.Name()) {
			continue
		}

		skillPath := filepath.Join(sourcePath, entry.Name())
		if _, err := os.Stat(filepath.Join(skillPath, "SKILL.md")); err != nil {
			continue
		}

		meta, _ := install.ReadMeta(skillPath)
		if meta != nil && meta.Source != "" {
			skills = append(skills, projectSkillEntry{Name: entry.Name(), Source: meta.Source, Remote: true})
		} else {
			skills = append(skills, projectSkillEntry{Name: entry.Name(), Remote: false})
		}
	}

	sort.Slice(skills, func(i, j int) bool {
		return skills[i].Name < skills[j].Name
	})

	fmt.Println("Skills (project):")
	if len(skills) == 0 {
		fmt.Println("  (none)")
		return nil
	}

	maxLen := 0
	for _, skill := range skills {
		if len(skill.Name) > maxLen {
			maxLen = len(skill.Name)
		}
	}

	remoteCount := 0
	for _, skill := range skills {
		if skill.Remote {
			remoteCount++
			fmt.Printf("  %-*s  remote  %s\n", maxLen, skill.Name, skill.Source)
		} else {
			fmt.Printf("  %-*s  local\n", maxLen, skill.Name)
		}
	}

	fmt.Println()
	fmt.Printf("%d skill(s): %d remote, %d local\n", len(skills), remoteCount, len(skills)-remoteCount)
	return nil
}
