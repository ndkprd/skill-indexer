package skill

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanDir(t *testing.T) {
	t.Run("real fixture directory", func(t *testing.T) {
		skills, warnings := ScanDir("../../skills")

		if len(warnings) != 0 {
			t.Errorf("warnings = %v, want none", warnings)
		}
		if got, want := len(skills), 46; got != want {
			t.Errorf("len(skills) = %d, want %d", got, want)
		}
	})

	t.Run("skips a broken skill directory with one warning", func(t *testing.T) {
		root := t.TempDir()

		copyDir(t, "../../skills/vue", filepath.Join(root, "vue"))
		copyDir(t, "../../skills/gitlab-cli", filepath.Join(root, "gitlab-cli"))

		// Broken: no SKILL.md at all.
		if err := os.MkdirAll(filepath.Join(root, "no-skill-md"), 0o755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}

		skills, warnings := ScanDir(root)

		if got, want := len(warnings), 1; got != want {
			t.Fatalf("len(warnings) = %d, want %d (warnings: %v)", got, want, warnings)
		}
		if got, want := len(skills), 2; got != want {
			t.Fatalf("len(skills) = %d, want %d", got, want)
		}

		names := map[string]bool{}
		for _, s := range skills {
			names[s.DirName] = true
		}
		if !names["vue"] || !names["gitlab-cli"] {
			t.Errorf("skills = %v, want vue and gitlab-cli present", names)
		}
	})
}
