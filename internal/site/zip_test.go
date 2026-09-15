package site

import (
	"archive/zip"
	"path/filepath"
	"strings"
	"testing"

	"skillstore/internal/skill"
)

func TestZipSkillDir(t *testing.T) {
	s := &skill.Skill{
		DirName: "vue",
		DirPath: "../../skills/vue",
	}

	destZipPath := filepath.Join(t.TempDir(), "vue.zip")
	if err := ZipSkillDir(s, destZipPath); err != nil {
		t.Fatalf("ZipSkillDir() error = %v", err)
	}

	r, err := zip.OpenReader(destZipPath)
	if err != nil {
		t.Fatalf("zip.OpenReader(%q) error = %v", destZipPath, err)
	}
	defer r.Close()

	var (
		hasSkillMD     bool
		hasReferenceMD bool
		hasBareSkillMD bool
	)
	for _, f := range r.File {
		switch {
		case f.Name == "vue/SKILL.md":
			hasSkillMD = true
		case f.Name == "SKILL.md":
			hasBareSkillMD = true
		case strings.HasPrefix(f.Name, "vue/references/") && strings.HasSuffix(f.Name, ".md"):
			hasReferenceMD = true
		}
	}

	if !hasSkillMD {
		t.Error("archive missing entry \"vue/SKILL.md\"")
	}
	if !hasReferenceMD {
		t.Error("archive missing any \"vue/references/*.md\" entry")
	}
	if hasBareSkillMD {
		t.Error("archive contains a bare \"SKILL.md\" entry at the root; entries must be prefixed with \"vue/\"")
	}
}
