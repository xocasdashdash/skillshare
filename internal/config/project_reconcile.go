package config

import (
	"fmt"
	"os"
	"path/filepath"

	"skillshare/internal/install"
	"skillshare/internal/utils"
)

// ReconcileProjectSkills scans the project source directory for remotely-installed
// skills (those with install metadata) and ensures they are listed in ProjectConfig.Skills[].
// It also updates .skillshare/.gitignore for each tracked skill.
func ReconcileProjectSkills(projectRoot string, projectCfg *ProjectConfig, sourcePath string) error {
	entries, err := os.ReadDir(sourcePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // no skills dir yet
		}
		return fmt.Errorf("failed to read project skills: %w", err)
	}

	changed := false
	index := map[string]int{}
	for i, skill := range projectCfg.Skills {
		index[skill.Name] = i
	}

	for _, entry := range entries {
		if !entry.IsDir() || utils.IsHidden(entry.Name()) {
			continue
		}

		skillName := entry.Name()
		skillPath := filepath.Join(sourcePath, skillName)
		meta, err := install.ReadMeta(skillPath)
		if err != nil || meta == nil || meta.Source == "" {
			continue
		}

		if existingIdx, ok := index[skillName]; ok {
			if projectCfg.Skills[existingIdx].Source != meta.Source {
				projectCfg.Skills[existingIdx].Source = meta.Source
				changed = true
			}
		} else {
			projectCfg.Skills = append(projectCfg.Skills, ProjectSkill{
				Name:   skillName,
				Source: meta.Source,
			})
			index[skillName] = len(projectCfg.Skills) - 1
			changed = true
		}

		if err := install.UpdateGitIgnore(filepath.Join(projectRoot, ".skillshare"), filepath.Join("skills", skillName)); err != nil {
			return fmt.Errorf("failed to update .skillshare/.gitignore: %w", err)
		}
	}

	if changed {
		if err := projectCfg.Save(projectRoot); err != nil {
			return err
		}
	}

	return nil
}
