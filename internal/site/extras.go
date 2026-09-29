package site

import (
	"fmt"
	"os"
	"path/filepath"
)

// extrasDir is where user-supplied CSS/JS files land, under assets/.
const extrasDir = "extra"

// copyExtras copies opts.ExtraCSS and opts.ExtraJS into
// <outputDir>/assets/extra, keeping each file's basename. It rejects
// directories, wrong extensions, and basenames used more than once across
// both lists.
func copyExtras(outputDir string, opts Options) error {
	seen := map[string]bool{}
	for _, group := range []struct {
		ext   string
		paths []string
	}{{".css", opts.ExtraCSS}, {".js", opts.ExtraJS}} {
		for _, src := range group.paths {
			name := filepath.Base(src)
			if filepath.Ext(name) != group.ext {
				return fmt.Errorf("extra file %q: want a %s file", src, group.ext)
			}
			if seen[name] {
				return fmt.Errorf("extra file %q: duplicate basename %q", src, name)
			}
			seen[name] = true
			if err := copyExtra(src, filepath.Join(outputDir, "assets", extrasDir, name)); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyExtra(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat extra file %q: %w", src, err)
	}
	if info.IsDir() {
		return fmt.Errorf("extra file %q is a directory", src)
	}
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read extra file %q: %w", src, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("create dir for %q: %w", dst, err)
	}
	if err := os.WriteFile(dst, content, 0o644); err != nil {
		return fmt.Errorf("write extra file %q: %w", dst, err)
	}
	return nil
}

func basenames(paths []string) []string {
	names := make([]string, len(paths))
	for i, p := range paths {
		names[i] = filepath.Base(p)
	}
	return names
}
