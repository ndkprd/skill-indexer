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

	if err := Render(skills, outDir, Options{}); err != nil {
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

	if err := Render(skills, outDir, Options{BaseURL: "/skills"}); err != nil {
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

func TestRenderDefaultsOmitOptionalParts(t *testing.T) {
	outDir := t.TempDir()
	if err := Render(fixtureSkills(), outDir, Options{}); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	indexHTML := readFile(t, filepath.Join(outDir, "index.html"))
	for _, notWant := range []string{`site-title`, `__REPO_URL__`, `data-install-cmd="all"`, `assets/extra/`} {
		if strings.Contains(indexHTML, notWant) {
			t.Errorf("default index.html unexpectedly contains %q", notWant)
		}
	}
	if !strings.Contains(indexHTML, `data-view="list"`) || !strings.Contains(indexHTML, `data-view="tile"`) {
		t.Error("index.html missing the tile/list view toggle")
	}
	if !strings.Contains(indexHTML, `id="panel-download-skill"`) {
		t.Error("index.html missing the .skill download link")
	}
}

func TestRenderSiteNameAndRepoURL(t *testing.T) {
	outDir := t.TempDir()
	opts := Options{SiteName: "Example Skills", RepoURL: "https://example.com/g/skills.git"}
	if err := Render(fixtureSkills(), outDir, opts); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	indexHTML := readFile(t, filepath.Join(outDir, "index.html"))
	for _, want := range []string{
		`<h1 class="site-title">`,
		`<title>Example Skills</title>`,
		`window.__REPO_URL__ = "https://example.com/g/skills.git";`,
		`data-install-cmd="all"`,
	} {
		if !strings.Contains(indexHTML, want) {
			t.Errorf("index.html missing %q", want)
		}
	}
	main := indexHTML[strings.Index(indexHTML, "<main"):strings.Index(indexHTML, "</main>")]
	if !strings.Contains(main, `data-install-cmd="all"`) {
		t.Error("install-everything block must live inside main (scrolls away, inert with it)")
	}
}

func TestRenderExtras(t *testing.T) {
	src := t.TempDir()
	css := writeTempFile(t, src, "theme.css", "body{}")
	js := writeTempFile(t, src, "theme.js", "void 0")
	outDir := t.TempDir()

	opts := Options{BaseURL: "/skills", ExtraCSS: []string{css}, ExtraJS: []string{js}}
	if err := Render(fixtureSkills(), outDir, opts); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	indexHTML := readFile(t, filepath.Join(outDir, "index.html"))
	for _, want := range []string{
		`href="/skills/assets/extra/theme.css"`,
		`src="/skills/assets/extra/theme.js"`,
	} {
		if !strings.Contains(indexHTML, want) {
			t.Errorf("index.html missing %q", want)
		}
	}
	if got := readFile(t, filepath.Join(outDir, "assets", "extra", "theme.css")); got != "body{}" {
		t.Errorf("copied css = %q", got)
	}
	if strings.Index(indexHTML, "assets/style.css") > strings.Index(indexHTML, "assets/extra/theme.css") {
		t.Error("extra css must load after style.css")
	}
}

func TestRenderExtrasErrors(t *testing.T) {
	src := t.TempDir()
	css := writeTempFile(t, src, "a.css", "")
	otherDir := t.TempDir()
	dupCSS := writeTempFile(t, otherDir, "a.css", "")
	txt := writeTempFile(t, src, "a.txt", "")

	tests := []struct {
		name string
		opts Options
	}{
		{"missing file", Options{ExtraCSS: []string{filepath.Join(src, "nope.css")}}},
		{"directory", Options{ExtraCSS: []string{src + "/dir.css"}}},
		{"wrong extension", Options{ExtraCSS: []string{txt}}},
		{"css given as js", Options{ExtraJS: []string{css}}},
		{"duplicate basename", Options{ExtraCSS: []string{css, dupCSS}}},
	}
	if err := os.Mkdir(src+"/dir.css", 0o755); err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Render(fixtureSkills(), t.TempDir(), tt.opts); err == nil {
				t.Error("Render() error = nil, want error")
			}
		})
	}
}

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
	return path
}
