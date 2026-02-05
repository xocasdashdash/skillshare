package integration

import (
	"testing"

	"skillshare/internal/testutil"
)

func TestStatusProject_ShowsSyncedTargets(t *testing.T) {
	sb := testutil.NewSandbox(t)
	defer sb.Cleanup()
	projectRoot := sb.SetupProjectDir("claude-code", "cursor")
	sb.CreateProjectSkill(projectRoot, "skill-a", map[string]string{"SKILL.md": "# A"})

	// Sync first
	sb.RunCLIInDir(projectRoot, "sync", "-p").AssertSuccess(t)

	result := sb.RunCLIInDir(projectRoot, "status", "-p")
	result.AssertSuccess(t)
	result.AssertOutputContains(t, "synced")
	result.AssertOutputContains(t, "1/1")
}

func TestStatusProject_ShowsMissing(t *testing.T) {
	sb := testutil.NewSandbox(t)
	defer sb.Cleanup()
	projectRoot := sb.SetupProjectDir("claude-code")
	sb.CreateProjectSkill(projectRoot, "unsynced", map[string]string{"SKILL.md": "# U"})

	// Don't sync — should show missing
	result := sb.RunCLIInDir(projectRoot, "status", "-p")
	result.AssertSuccess(t)
	result.AssertOutputContains(t, "missing")
}

func TestStatusProject_ShowsProjectPath(t *testing.T) {
	sb := testutil.NewSandbox(t)
	defer sb.Cleanup()
	projectRoot := sb.SetupProjectDir("claude-code")

	result := sb.RunCLIInDir(projectRoot, "status", "-p")
	result.AssertSuccess(t)
	result.AssertOutputContains(t, "Project:")
	result.AssertOutputContains(t, "Source:")
}
