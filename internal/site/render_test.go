package site

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"skill-indexer/internal/skill"
)

func fixtureSkills() []*skill.Skill {
	return []*skill.Skill{
		{
			Name:        "Alpha Tool",
			Description: "Does alpha things for developers.",
			Body:        "# Alpha\n\nSome **body** content.\n",
			Metadata:    map[string]any{"version": "1.0.0", "author": "octo"},
			DirName:     "alpha-tool",
			DirPath:     "alpha-tool",
		},
		{
			Name:        "Bravo Helper",
			Description: "Helps with bravo tasks.",
			Body:        "Bravo body.\n",
			Metadata:    map[string]any{},
			DirName:     "bravo-helper",
			DirPath:     "bravo-helper",
		},
	}
}

func TestRender(t *testing.T) {
	skills := fixtureSkills()
	outDir := t.TempDir()

	if err := Render(skills, outDir, ""); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	indexPath := filepath.Join(outDir, "index.html")
	indexHTML := readFile(t, indexPath)
	for _, s := range skills {
		if !strings.Contains(indexHTML, s.Name) {
			t.Errorf("index.html missing skill name %q", s.Name)
		}
		if !strings.Contains(indexHTML, `href="#`+s.DirName+`"`) {
			t.Errorf("index.html missing panel-trigger link for %q", s.DirName)
		}
	}

	if !strings.Contains(indexHTML, `id="skill-panel"`) {
		t.Error("index.html missing the detail panel skeleton")
	}
	if !strings.Contains(indexHTML, `id="theme-toggle"`) {
		t.Error("index.html missing the theme toggle button")
	}

	for _, asset := range []string{"style.css", "app.js", "fuse.min.js"} {
		if _, err := os.Stat(filepath.Join(outDir, "assets", asset)); err != nil {
			t.Errorf("expected asset %q to be written: %v", asset, err)
		}
	}
}

func TestRenderBaseURL(t *testing.T) {
	skills := fixtureSkills()
	outDir := t.TempDir()

	if err := Render(skills, outDir, "/skills"); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	indexHTML := readFile(t, filepath.Join(outDir, "index.html"))
	for _, want := range []string{
		`href="/skills/assets/style.css"`,
		`src="/skills/assets/fuse.min.js"`,
		`src="/skills/assets/app.js"`,
		`href="/skills/"`,
		`window.__BASE_URL__ = "/skills";`,
	} {
		if !strings.Contains(indexHTML, want) {
			t.Errorf("index.html missing %q with baseURL set", want)
		}
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
