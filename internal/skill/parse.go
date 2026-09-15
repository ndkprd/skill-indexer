package skill

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// frontmatterDelim is the YAML frontmatter fence used at the top of a
// SKILL.md file.
const frontmatterDelim = "---"

// ParseFrontmatter reads the SKILL.md file at path, splits its leading
// YAML frontmatter from the Markdown body, and returns a populated Skill.
//
// The frontmatter's name and description keys map to Skill.Name and
// Skill.Description and are required. A nested metadata map, along with
// every other top-level frontmatter key (license, compatibility, version,
// author, dependencies, etc.), is folded as-is into Skill.Metadata.
func ParseFrontmatter(path string) (*Skill, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read skill file %s: %w", path, err)
	}

	frontmatter, body, err := splitFrontmatter(string(raw))
	if err != nil {
		return nil, fmt.Errorf("parse skill file %s: %w", path, err)
	}

	fields := make(map[string]any)
	if err := yaml.Unmarshal([]byte(frontmatter), &fields); err != nil {
		return nil, fmt.Errorf("parse skill file %s: unmarshal frontmatter: %w", path, err)
	}

	s, err := skillFromFields(fields)
	if err != nil {
		return nil, fmt.Errorf("parse skill file %s: %w", path, err)
	}
	s.Body = body

	dir := filepath.Dir(path)
	s.DirPath = dir
	s.DirName = filepath.Base(dir)

	return s, nil
}

// splitFrontmatter separates the leading "---"-delimited YAML block from
// the remainder of a SKILL.md file's content, returning the frontmatter
// YAML and the trailing Markdown body separately.
func splitFrontmatter(content string) (frontmatter, body string, err error) {
	normalized := strings.ReplaceAll(content, "\r\n", "\n")
	trimmed := strings.TrimPrefix(normalized, "\n")

	if !strings.HasPrefix(trimmed, frontmatterDelim) {
		return "", "", fmt.Errorf("missing opening %q frontmatter delimiter", frontmatterDelim)
	}

	rest := trimmed[len(frontmatterDelim):]
	rest = strings.TrimPrefix(rest, "\n")

	closeIdx := strings.Index(rest, "\n"+frontmatterDelim)
	if closeIdx == -1 {
		return "", "", fmt.Errorf("missing closing %q frontmatter delimiter", frontmatterDelim)
	}

	frontmatter = rest[:closeIdx]

	after := rest[closeIdx+len("\n"+frontmatterDelim):]
	if nl := strings.Index(after, "\n"); nl != -1 {
		body = after[nl+1:]
	} else {
		body = ""
	}

	return frontmatter, body, nil
}

// skillFromFields maps a raw frontmatter field map onto a Skill, requiring
// name and description, and folding every remaining key (including a
// nested metadata map) into Metadata.
func skillFromFields(fields map[string]any) (*Skill, error) {
	name, err := stringField(fields, "name")
	if err != nil {
		return nil, err
	}
	description, err := stringField(fields, "description")
	if err != nil {
		return nil, err
	}

	metadata := make(map[string]any)
	if nested, ok := fields["metadata"]; ok {
		for k, v := range coerceMap(nested) {
			metadata[k] = v
		}
	}
	for k, v := range fields {
		if k == "name" || k == "description" || k == "metadata" {
			continue
		}
		metadata[k] = v
	}

	return &Skill{
		Name:        name,
		Description: description,
		Metadata:    metadata,
	}, nil
}

// stringField extracts a required string-valued key from a frontmatter
// field map, returning an error if the key is absent or not a string.
func stringField(fields map[string]any, key string) (string, error) {
	v, ok := fields[key]
	if !ok {
		return "", fmt.Errorf("missing required frontmatter field %q", key)
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return "", fmt.Errorf("frontmatter field %q must be a non-empty string", key)
	}
	return s, nil
}

// coerceMap normalizes a decoded YAML mapping value into a
// map[string]any. yaml.v3 decodes mapping values into
// map[string]any when the target is any, so this mainly guards
// against a non-map value being provided for metadata.
func coerceMap(v any) map[string]any {
	if m, ok := v.(map[string]any); ok {
		return m
	}
	return nil
}
