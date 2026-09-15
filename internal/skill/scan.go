package skill

import (
	"fmt"
	"os"
	"path/filepath"
)

// skillFileName is the frontmatter file every skill directory is expected
// to contain.
const skillFileName = "SKILL.md"

// ScanDir discovers skills under the immediate subdirectories of root.
//
// Each subdirectory's SKILL.md is parsed via ParseFrontmatter. A
// subdirectory missing SKILL.md or containing unparsable frontmatter is
// skipped and recorded as a warning rather than aborting the scan.
func ScanDir(root string) (skills []*Skill, warnings []error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, []error{fmt.Errorf("read skill dir %s: %w", root, err)}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		s, err := ParseFrontmatter(filepath.Join(root, entry.Name(), skillFileName))
		if err != nil {
			warnings = append(warnings, fmt.Errorf("skip skill %q: %w", entry.Name(), err))
			continue
		}

		skills = append(skills, s)
	}

	return skills, warnings
}
