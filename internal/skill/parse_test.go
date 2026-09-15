package skill

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	t.Run("nested metadata map", func(t *testing.T) {
		s, err := ParseFrontmatter("../../skills/vue/SKILL.md")
		if err != nil {
			t.Fatalf("ParseFrontmatter returned error: %v", err)
		}
		if s.Name != "vue" {
			t.Errorf("Name = %q, want %q", s.Name, "vue")
		}
		if s.Description == "" {
			t.Error("Description is empty, want non-empty")
		}
		if got, want := s.Metadata["author"], "Anthony Fu"; got != want {
			t.Errorf("Metadata[author] = %v, want %v", got, want)
		}
		if got, want := s.Metadata["version"], "2026.1.31"; got != want {
			t.Errorf("Metadata[version] = %v, want %v", got, want)
		}
		if s.Metadata["source"] == nil {
			t.Error("Metadata[source] is nil, want a value")
		}
		if !strings.Contains(s.Body, "# Vue") {
			t.Errorf("Body does not contain expected heading, got: %q", truncate(s.Body, 80))
		}
	})

	t.Run("top-level license and compatibility with 4-space metadata indent", func(t *testing.T) {
		s, err := ParseFrontmatter("../../skills/vueuse-functions/SKILL.md")
		if err != nil {
			t.Fatalf("ParseFrontmatter returned error: %v", err)
		}
		if s.Name != "vueuse-functions" {
			t.Errorf("Name = %q, want %q", s.Name, "vueuse-functions")
		}
		if got, want := s.Metadata["license"], "MIT"; got != want {
			t.Errorf("Metadata[license] = %v, want %v", got, want)
		}
		if got, want := s.Metadata["compatibility"], "Requires Vue 3 (or above) or Nuxt 3 (or above) project"; got != want {
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
		s, err := ParseFrontmatter("../../skills/asdp-create-plan/SKILL.md")
		if err != nil {
			t.Fatalf("ParseFrontmatter returned error: %v", err)
		}
		if s.Name != "asdp-create-plan" {
			t.Errorf("Name = %q, want %q", s.Name, "asdp-create-plan")
		}
		if s.Description == "" {
			t.Error("Description is empty, want non-empty")
		}
		if len(s.Metadata) != 0 {
			t.Errorf("Metadata = %v, want empty", s.Metadata)
		}
	})

	t.Run("top-level dependencies list folded into metadata", func(t *testing.T) {
		s, err := ParseFrontmatter("../../skills/gitlab-cli/SKILL.md")
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

	t.Run("top-level allowed-tools string folded into metadata", func(t *testing.T) {
		s, err := ParseFrontmatter("../../skills/playwright-cli/SKILL.md")
		if err != nil {
			t.Fatalf("ParseFrontmatter returned error: %v", err)
		}
		allowedTools, ok := s.Metadata["allowed-tools"].(string)
		if !ok {
			t.Fatalf("Metadata[allowed-tools] = %T(%v), want string", s.Metadata["allowed-tools"], s.Metadata["allowed-tools"])
		}
		if !strings.Contains(allowedTools, "playwright-cli") {
			t.Errorf("allowed-tools = %q, want it to mention playwright-cli", allowedTools)
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
