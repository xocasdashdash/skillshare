package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"skillshare/internal/install"
	"skillshare/internal/ui"
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

	if len(skills) == 0 {
		ui.Info("No skills installed")
		ui.Info("Use 'skillshare install -p <source>' to install a skill")
		return nil
	}

	ui.Header("Installed skills (project)")

	maxNameLen := 0
	for _, skill := range skills {
		if len(skill.Name) > maxNameLen {
			maxNameLen = len(skill.Name)
		}
	}

	for _, skill := range skills {
		suffix := "local"
		if skill.Remote {
			suffix = abbreviateSource(skill.Source)
		}
		format := fmt.Sprintf("  %s→%s %%-%ds  %s%%s%s\n", ui.Cyan, ui.Reset, maxNameLen, ui.Gray, ui.Reset)
		fmt.Printf(format, skill.Name, suffix)
	}

	fmt.Println()
	remoteCount := 0
	for _, skill := range skills {
		if skill.Remote {
			remoteCount++
		}
	}
	ui.Info("%d skill(s): %d remote, %d local", len(skills), remoteCount, len(skills)-remoteCount)
	return nil
}
