package site

import (
	"archive/zip"
	"path/filepath"
	"strings"
	"testing"

	"skill-indexer/internal/skill"
)

func TestZipSkillDir(t *testing.T) {
	s := &skill.Skill{
		DirName: "skill-creator",
		DirPath: "../../examples/skills/skill-creator",
	}

	destZipPath := filepath.Join(t.TempDir(), "skill-creator.zip")
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
		case f.Name == "skill-creator/SKILL.md":
			hasSkillMD = true
		case f.Name == "SKILL.md":
			hasBareSkillMD = true
		case strings.HasPrefix(f.Name, "skill-creator/references/") && strings.HasSuffix(f.Name, ".md"):
			hasReferenceMD = true
		}
	}

	if !hasSkillMD {
		t.Error("archive missing entry \"skill-creator/SKILL.md\"")
	}
	if !hasReferenceMD {
		t.Error("archive missing any \"skill-creator/references/*.md\" entry")
	}
	if hasBareSkillMD {
		t.Error("archive contains a bare \"SKILL.md\" entry at the root; entries must be prefixed with \"skill-creator/\"")
	}
}
