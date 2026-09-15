package skill

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	t.Run("nested metadata map", func(t *testing.T) {
		s, err := ParseFrontmatter("../../examples/skills/web-design-guidelines/SKILL.md")
		if err != nil {
			t.Fatalf("ParseFrontmatter returned error: %v", err)
		}
		if s.Name != "web-design-guidelines" {
			t.Errorf("Name = %q, want %q", s.Name, "web-design-guidelines")
		}
		if s.Description == "" {
			t.Error("Description is empty, want non-empty")
		}
		if got, want := s.Metadata["author"], "vercel"; got != want {
			t.Errorf("Metadata[author] = %v, want %v", got, want)
		}
		if got, want := s.Metadata["version"], "1.0.0"; got != want {
			t.Errorf("Metadata[version] = %v, want %v", got, want)
		}
		if s.Metadata["argument-hint"] == nil {
			t.Error("Metadata[argument-hint] is nil, want a value")
		}
		if !strings.Contains(s.Body, "# Web Interface Guidelines") {
			t.Errorf("Body does not contain expected heading, got: %q", truncate(s.Body, 80))
		}
	})

	t.Run("top-level license and compatibility with 4-space metadata indent", func(t *testing.T) {
		// Synthetic rather than a real fixture: no single real skill in
		// examples/skills/ combines top-level license/compatibility with a
		// 4-space-indented metadata block, and the fixture set has already
		// changed shape twice. Constructing the exact shape here decouples
		// this edge case from whichever skills happen to exist.
		dir := t.TempDir()
		path := filepath.Join(dir, "SKILL.md")
		writeFile(t, path, "---\n"+
			"name: synthetic-skill\n"+
			"description: A synthetic skill for testing 4-space metadata indent.\n"+
			"license: MIT\n"+
			"compatibility: Requires a compatible thing\n"+
			"metadata:\n"+
			"    author: someone\n"+
			"    version: \"1.1\"\n"+
			"---\n\nBody.\n")

		s, err := ParseFrontmatter(path)
		if err != nil {
			t.Fatalf("ParseFrontmatter returned error: %v", err)
		}
		if got, want := s.Metadata["license"], "MIT"; got != want {
			t.Errorf("Metadata[license] = %v, want %v", got, want)
		}
		if got, want := s.Metadata["compatibility"], "Requires a compatible thing"; got != want {
			t.Errorf("Metadata[compatibility] = %v, want %v", got, want)
		}
		if s.Metadata["author"] == nil {
			t.Error("Metadata[author] is nil, want a value from the 4-space-indented metadata block")
		}
		if got, want := s.Metadata["version"], "1.1"; got != want {
			t.Errorf("Metadata[version] = %v, want %v", got, want)
		}
	})

	t.Run("minimal frontmatter with only name and description", func(t *testing.T) {
		s, err := ParseFrontmatter("../../examples/skills/find-skills/SKILL.md")
		if err != nil {
			t.Fatalf("ParseFrontmatter returned error: %v", err)
		}
		if s.Name != "find-skills" {
			t.Errorf("Name = %q, want %q", s.Name, "find-skills")
		}
		if s.Description == "" {
			t.Error("Description is empty, want non-empty")
		}
		if len(s.Metadata) != 0 {
			t.Errorf("Metadata = %v, want empty", s.Metadata)
		}
	})

	t.Run("top-level dependencies list folded into metadata", func(t *testing.T) {
		s, err := ParseFrontmatter("../../examples/skills/gitlab-cli/SKILL.md")
		if err != nil {
			t.Fatalf("ParseFrontmatter returned error: %v", err)
		}
		deps, ok := s.Metadata["dependencies"].([]any)
		if !ok {
			t.Fatalf("Metadata[dependencies] = %T(%v), want []any", s.Metadata["dependencies"], s.Metadata["dependencies"])
		}
		if len(deps) == 0 {
			t.Error("dependencies list is empty, want at least one entry")
		}
		if got, want := deps[0], "glab"; got != want {
			t.Errorf("dependencies[0] = %v, want %v", got, want)
		}
	})

	t.Run("top-level string key folded into metadata", func(t *testing.T) {
		s, err := ParseFrontmatter("../../examples/skills/frontend-design/SKILL.md")
		if err != nil {
			t.Fatalf("ParseFrontmatter returned error: %v", err)
		}
		license, ok := s.Metadata["license"].(string)
		if !ok {
			t.Fatalf("Metadata[license] = %T(%v), want string", s.Metadata["license"], s.Metadata["license"])
		}
		if !strings.Contains(license, "LICENSE") {
			t.Errorf("license = %q, want it to mention LICENSE", license)
		}
	})

	t.Run("missing name returns error", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "SKILL.md")
		writeFile(t, path, "---\ndescription: no name here\n---\n\nBody.\n")

		if _, err := ParseFrontmatter(path); err == nil {
			t.Fatal("ParseFrontmatter returned nil error, want error for missing name")
		}
	})

	t.Run("malformed YAML returns error", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "SKILL.md")
		writeFile(t, path, "---\nname: broken\n  description: [unterminated\n---\n\nBody.\n")

		if _, err := ParseFrontmatter(path); err == nil {
			t.Fatal("ParseFrontmatter returned nil error, want error for malformed YAML")
		}
	})

	t.Run("missing file returns error", func(t *testing.T) {
		if _, err := ParseFrontmatter(filepath.Join(t.TempDir(), "SKILL.md")); err == nil {
			t.Fatal("ParseFrontmatter returned nil error, want error for missing file")
		}
	})
}

// truncate shortens s to at most n bytes, for readable test failure output.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
