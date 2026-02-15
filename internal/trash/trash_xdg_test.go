package trash

import (
	"path/filepath"
	"testing"
)

func TestTrashDir_RespectsXDGConfigHome(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/custom/config")
	t.Setenv("SKILLSHARE_CONFIG", "")

	got := TrashDir()
	want := filepath.Join("/custom/config", "skillshare", "trash")
	if got != want {
		t.Errorf("TrashDir() = %q, want %q", got, want)
	}
}
