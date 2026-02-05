package main

import (
	"fmt"
	"os"
	"path/filepath"

	"skillshare/internal/config"
	"skillshare/internal/install"
	"skillshare/internal/utils"
)

func reconcileProjectRemoteSkills(runtime *projectRuntime) error {
	entries, err := os.ReadDir(runtime.sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read project skills: %w", err)
	}

	changed := false
	index := map[string]int{}
	for i, skill := range runtime.config.Skills {
		index[skill.Name] = i
	}

	for _, entry := range entries {
		if !entry.IsDir() || utils.IsHidden(entry.Name()) {
			continue
		}

		skillName := entry.Name()
		skillPath := filepath.Join(runtime.sourcePath, skillName)
		meta, err := install.ReadMeta(skillPath)
		if err != nil || meta == nil || meta.Source == "" {
			continue
		}

		if existingIdx, ok := index[skillName]; ok {
			if runtime.config.Skills[existingIdx].Source != meta.Source {
				runtime.config.Skills[existingIdx].Source = meta.Source
				changed = true
			}
		} else {
			runtime.config.Skills = append(runtime.config.Skills, config.ProjectSkill{
				Name:   skillName,
				Source: meta.Source,
			})
			index[skillName] = len(runtime.config.Skills) - 1
			changed = true
		}

		if err := install.UpdateGitIgnore(filepath.Join(runtime.root, ".skillshare"), filepath.Join("skills", skillName)); err != nil {
			return fmt.Errorf("failed to update .skillshare/.gitignore: %w", err)
		}
	}

	if changed {
		if err := runtime.config.Save(runtime.root); err != nil {
			return err
		}
	}

	return nil
}
