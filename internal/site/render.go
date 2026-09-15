package site

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"

	"skill-repo-store/internal/skill"
)

//go:embed templates/*.html.tmpl
var templatesFS embed.FS

//go:embed assets
var assetsFS embed.FS

// cardBadgePriority lists the metadata keys shown on a card, in order of
// preference, before falling back to the first remaining key alphabetically.
var cardBadgePriority = []string{"version", "author", "license", "compatibility"}

// maxCardBadges caps how many metadata badges a card shows at a glance.
const maxCardBadges = 2

var markdown = goldmark.New(goldmark.WithExtensions(extension.GFM))

// basePage carries fields every rendered page needs.
type basePage struct {
	Title string
}

type badge struct {
	Key   string
	Value string
}

type cardView struct {
	Name        string
	DirName     string
	Description string
	Badges      []badge
}

type indexPageData struct {
	basePage
	Skills []cardView
}

type detailView struct {
	Name        string
	DirName     string
	Description string
	BodyHTML    template.HTML
	Metadata    []badge
	ZipPath     string
}

type detailPageData struct {
	basePage
	Skill detailView
}

// Render writes the full static site (index page, one detail page per
// skill, and the vendored static assets) under outputDir.
//
// The index and detail pages are parsed as two separate template sets
// (each paired with the shared layout) because both define a "content"
// block under the same name; parsing them together would let one
// silently override the other.
func Render(skills []*skill.Skill, outputDir string) error {
	indexTmpl, err := template.ParseFS(templatesFS, "templates/layout.html.tmpl", "templates/index.html.tmpl")
	if err != nil {
		return fmt.Errorf("parse index template: %w", err)
	}
	skillTmpl, err := template.ParseFS(templatesFS, "templates/layout.html.tmpl", "templates/skill.html.tmpl")
	if err != nil {
		return fmt.Errorf("parse skill template: %w", err)
	}

	if err := renderIndex(indexTmpl, skills, outputDir); err != nil {
		return err
	}
	if err := renderDetailPages(skillTmpl, skills, outputDir); err != nil {
		return err
	}
	if err := writeAssets(outputDir); err != nil {
		return err
	}
	return nil
}

func renderIndex(tmpl *template.Template, skills []*skill.Skill, outputDir string) error {
	cards := make([]cardView, len(skills))
	for i, s := range skills {
		cards[i] = buildCardView(s)
	}
	sort.Slice(cards, func(i, j int) bool { return cards[i].Name < cards[j].Name })

	data := indexPageData{
		basePage: basePage{Title: "Skill Repo Store"},
		Skills:   cards,
	}
	return renderToFile(tmpl, filepath.Join(outputDir, "index.html"), data)
}

func renderDetailPages(tmpl *template.Template, skills []*skill.Skill, outputDir string) error {
	skillsDir := filepath.Join(outputDir, "skills")
	if err := os.MkdirAll(skillsDir, 0o755); err != nil {
		return fmt.Errorf("create skills output dir: %w", err)
	}

	for _, s := range skills {
		view, err := buildDetailView(s)
		if err != nil {
			return fmt.Errorf("render skill %q: %w", s.DirName, err)
		}
		data := detailPageData{
			basePage: basePage{Title: view.Name + " · Skill Repo Store"},
			Skill:    view,
		}
		outPath := filepath.Join(skillsDir, s.DirName+".html")
		if err := renderToFile(tmpl, outPath, data); err != nil {
			return err
		}
	}
	return nil
}

// renderToFile executes the "layout" template (which pulls in the page's
// own "content" block) and writes the result to outPath.
func renderToFile(tmpl *template.Template, outPath string, data any) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("create output dir for %s: %w", outPath, err)
	}
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create %s: %w", outPath, err)
	}
	defer f.Close()

	if err := tmpl.ExecuteTemplate(f, "layout", data); err != nil {
		return fmt.Errorf("execute layout template for %s: %w", outPath, err)
	}
	return nil
}

func writeAssets(outputDir string) error {
	assetsOut := filepath.Join(outputDir, "assets")
	if err := os.MkdirAll(assetsOut, 0o755); err != nil {
		return fmt.Errorf("create assets output dir: %w", err)
	}

	return fs.WalkDir(assetsFS, "assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel("assets", path)
		if err != nil {
			return fmt.Errorf("relativize asset path %q: %w", path, err)
		}
		content, err := assetsFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read embedded asset %q: %w", path, err)
		}
		outPath := filepath.Join(assetsOut, rel)
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return fmt.Errorf("create dir for asset %q: %w", outPath, err)
		}
		if err := os.WriteFile(outPath, content, 0o644); err != nil {
			return fmt.Errorf("write asset %q: %w", outPath, err)
		}
		return nil
	})
}

func buildCardView(s *skill.Skill) cardView {
	return cardView{
		Name:        s.Name,
		DirName:     s.DirName,
		Description: s.Description,
		Badges:      pickBadges(s.Metadata, cardBadgePriority, maxCardBadges),
	}
}

func buildDetailView(s *skill.Skill) (detailView, error) {
	bodyHTML, err := renderMarkdown(s.Body)
	if err != nil {
		return detailView{}, err
	}

	return detailView{
		Name:        s.Name,
		DirName:     s.DirName,
		Description: s.Description,
		BodyHTML:    bodyHTML,
		Metadata:    sortedMetadata(s.Metadata),
		ZipPath:     "/downloads/" + s.DirName + ".zip",
	}, nil
}

func renderMarkdown(body string) (template.HTML, error) {
	var buf bytes.Buffer
	if err := markdown.Convert([]byte(body), &buf); err != nil {
		return "", fmt.Errorf("render markdown body: %w", err)
	}
	return template.HTML(buf.String()), nil
}

// pickBadges selects up to max metadata entries for compact display,
// preferring keys in priority order and falling back to the first
// remaining key alphabetically.
func pickBadges(metadata map[string]any, priority []string, max int) []badge {
	if len(metadata) == 0 {
		return nil
	}

	used := make(map[string]bool, max)
	var badges []badge

	for _, key := range priority {
		if len(badges) >= max {
			break
		}
		v, ok := metadata[key]
		if !ok {
			continue
		}
		badges = append(badges, badge{Key: key, Value: formatMetadataValue(v)})
		used[key] = true
	}

	if len(badges) < max {
		for _, key := range sortedKeys(metadata) {
			if len(badges) >= max {
				break
			}
			if used[key] {
				continue
			}
			badges = append(badges, badge{Key: key, Value: formatMetadataValue(metadata[key])})
		}
	}

	return badges
}

// sortedMetadata returns every metadata entry as a badge, ordered
// alphabetically by key for deterministic rendering.
func sortedMetadata(metadata map[string]any) []badge {
	keys := sortedKeys(metadata)
	badges := make([]badge, len(keys))
	for i, k := range keys {
		badges[i] = badge{Key: k, Value: formatMetadataValue(metadata[k])}
	}
	return badges
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// formatMetadataValue renders an arbitrary YAML-decoded value (string,
// bool, list, or otherwise) as a single display string.
func formatMetadataValue(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case []any:
		parts := make([]string, len(val))
		for i, item := range val {
			parts[i] = fmt.Sprint(item)
		}
		return strings.Join(parts, ", ")
	default:
		return fmt.Sprint(val)
	}
}
