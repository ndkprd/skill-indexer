// Package site renders the static skill hub site (HTML pages, zip
// downloads, and the client-side search index) from parsed skills.
package site

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"skill-repo-store/internal/skill"
)

// ZipSkillDir archives the entire contents of s.DirPath into a new zip
// file at destZipPath. Every entry is written under a top-level
// "<s.DirName>/" prefix (e.g. "vue/SKILL.md", "vue/references/foo.md") so
// that extracting the archive at some destination reproduces a
// "<dest>/<DirName>/..." layout matching the --skill-dir convention.
func ZipSkillDir(s *skill.Skill, destZipPath string) error {
	out, err := os.Create(destZipPath)
	if err != nil {
		return fmt.Errorf("create zip file %q: %w", destZipPath, err)
	}
	defer out.Close()

	zw := zip.NewWriter(out)
	defer zw.Close()

	if err := addDirToZip(zw, s.DirPath, s.DirName); err != nil {
		return fmt.Errorf("zip skill dir %q: %w", s.DirPath, err)
	}

	if err := zw.Close(); err != nil {
		return fmt.Errorf("finalize zip file %q: %w", destZipPath, err)
	}
	return nil
}

// addDirToZip walks every regular file under srcDir and writes it into zw
// as an entry named "<prefix>/<path relative to srcDir>".
func addDirToZip(zw *zip.Writer, srcDir, prefix string) error {
	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk %q: %w", path, err)
		}
		if d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return fmt.Errorf("relativize %q against %q: %w", path, srcDir, err)
		}

		entryName := filepath.ToSlash(filepath.Join(prefix, rel))
		if err := writeZipEntry(zw, path, entryName); err != nil {
			return fmt.Errorf("write zip entry %q: %w", entryName, err)
		}
		return nil
	})
}

// writeZipEntry copies the file at srcPath into zw as a new entry named
// entryName.
func writeZipEntry(zw *zip.Writer, srcPath, entryName string) error {
	f, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open %q: %w", srcPath, err)
	}
	defer f.Close()

	w, err := zw.Create(entryName)
	if err != nil {
		return fmt.Errorf("create zip entry: %w", err)
	}

	if _, err := io.Copy(w, f); err != nil {
		return fmt.Errorf("copy contents: %w", err)
	}
	return nil
}
