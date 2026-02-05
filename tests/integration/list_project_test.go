package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"skillshare/internal/testutil"
)

func TestListProject_ShowsLocalAndRemote(t *testing.T) {
	sb := testutil.NewSandbox(t)
	defer sb.Cleanup()
	projectRoot := sb.SetupProjectDir("claude-code")

	// Local skill (no meta)
	sb.CreateProjectSkill(projectRoot, "local-skill", map[string]string{
		"SKILL.md": "# Local",
	})

	// Remote skill (with meta)
	skillDir := sb.CreateProjectSkill(projectRoot, "remote-skill", map[string]string{
		"SKILL.md": "# Remote",
	})
	meta := map[string]interface{}{
		"source": "someone/skills/remote-skill",
		"type":   "github",
	}
	metaJSON, _ := json.Marshal(meta)
	os.WriteFile(filepath.Join(skillDir, ".skillshare-meta.json"), metaJSON, 0644)

	result := sb.RunCLIInDir(projectRoot, "list", "-p")
	result.AssertSuccess(t)
	result.AssertOutputContains(t, "local-skill")
	result.AssertOutputContains(t, "local")
	result.AssertOutputContains(t, "remote-skill")
	result.AssertOutputContains(t, "remote")
}

func TestListProject_Empty(t *testing.T) {
	sb := testutil.NewSandbox(t)
	defer sb.Cleanup()
	projectRoot := sb.SetupProjectDir("claude-code")

	result := sb.RunCLIInDir(projectRoot, "list", "-p")
	result.AssertSuccess(t)
	result.AssertOutputContains(t, "No skills installed")
}

func TestListProject_AutoDetectsMode(t *testing.T) {
	sb := testutil.NewSandbox(t)
	defer sb.Cleanup()
	projectRoot := sb.SetupProjectDir("claude-code")
	sb.CreateProjectSkill(projectRoot, "skill", map[string]string{"SKILL.md": "# S"})

	result := sb.RunCLIInDir(projectRoot, "list")
	result.AssertSuccess(t)
	result.AssertOutputContains(t, "Installed skills (project)")
}
