// Package skill models a Claude Code skill parsed from a SKILL.md file.
package skill

// Skill is a single skill discovered under a skill directory.
type Skill struct {
	Name        string
	Description string
	Body        string
	Metadata    map[string]any
	DirName     string
	DirPath     string
}
