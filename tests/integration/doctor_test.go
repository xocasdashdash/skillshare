package integration

import (
	"os"
	"path/filepath"
	"testing"

	"skillshare/internal/testutil"
)

func TestDoctor_AllGood_PassesAll(t *testing.T) {
	sb := testutil.NewSandbox(t)
	defer sb.Cleanup()

	sb.CreateSkill("skill1", map[string]string{"SKILL.md": "# Skill 1"})
	targetPath := sb.CreateTarget("claude")

	// Create synced state
	os.Symlink(filepath.Join(sb.SourcePath, "skill1"), filepath.Join(targetPath, "skill1"))

	sb.WriteConfig(`source: ` + sb.SourcePath + `
targets:
  claude:
    path: ` + targetPath + `
`)

	result := sb.RunCLI("doctor")

	result.AssertSuccess(t)
	result.AssertOutputContains(t, "All checks passed")
}

func TestDoctor_NoConfig_ShowsError(t *testing.T) {
	sb := testutil.NewSandbox(t)
	defer sb.Cleanup()

	os.Remove(sb.ConfigPath)

	result := sb.RunCLI("doctor")

	result.AssertSuccess(t) // doctor doesn't return error, it reports issues
	result.AssertOutputContains(t, "not found")
}

func TestDoctor_NoSource_ShowsError(t *testing.T) {
	sb := testutil.NewSandbox(t)
	defer sb.Cleanup()

	// Remove source directory
	os.RemoveAll(sb.SourcePath)

	sb.WriteConfig(`source: ` + sb.SourcePath + `
targets: {}
`)

	result := sb.RunCLI("doctor")

	result.AssertSuccess(t)
	result.AssertOutputContains(t, "Source not found")
}

func TestDoctor_ChecksSymlinkSupport(t *testing.T) {
	sb := testutil.NewSandbox(t)
	defer sb.Cleanup()

	sb.WriteConfig(`source: ` + sb.SourcePath + `
targets: {}
`)

	result := sb.RunCLI("doctor")

	result.AssertSuccess(t)
	result.AssertOutputContains(t, "Symlink")
}

func TestDoctor_TargetIssues_ShowsProblems(t *testing.T) {
	sb := testutil.NewSandbox(t)
	defer sb.Cleanup()

	// Point to non-existent directory with non-existent parent
	sb.WriteConfig(`source: ` + sb.SourcePath + `
targets:
  broken:
    path: /nonexistent/path/skills
`)

	result := sb.RunCLI("doctor")

	result.AssertSuccess(t)
	result.AssertAnyOutputContains(t, "issue")
}

func TestDoctor_WrongSymlink_ShowsWarning(t *testing.T) {
	sb := testutil.NewSandbox(t)
	defer sb.Cleanup()

	// Create wrong symlink
	wrongSource := filepath.Join(sb.Home, "wrong-source")
	os.MkdirAll(wrongSource, 0755)

	targetPath := filepath.Join(sb.Home, ".claude", "skills")
	os.MkdirAll(filepath.Dir(targetPath), 0755)
	os.Symlink(wrongSource, targetPath)

	sb.WriteConfig(`source: ` + sb.SourcePath + `
targets:
  claude:
    path: ` + targetPath + `
    mode: symlink
`)

	result := sb.RunCLI("doctor")

	result.AssertSuccess(t)
	result.AssertAnyOutputContains(t, "wrong location")
}

func TestDoctor_ShowsSkillCount(t *testing.T) {
	sb := testutil.NewSandbox(t)
	defer sb.Cleanup()

	sb.CreateSkill("skill1", map[string]string{"SKILL.md": "# Skill 1"})
	sb.CreateSkill("skill2", map[string]string{"SKILL.md": "# Skill 2"})
	sb.CreateSkill("skill3", map[string]string{"SKILL.md": "# Skill 3"})

	sb.WriteConfig(`source: ` + sb.SourcePath + `
targets: {}
`)

	result := sb.RunCLI("doctor")

	result.AssertSuccess(t)
	result.AssertOutputContains(t, "3 skills")
}
