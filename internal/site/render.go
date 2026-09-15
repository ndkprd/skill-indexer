package site

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"skill-indexer/internal/skill"
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

// Version is shown in the generated site's footer trademark line.
const Version = "0.1.0"

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
	Title   string
	Version string
	Skills  []cardView
}

// Render writes the single-page static site (card grid plus its vendored
// static assets) under outputDir. Skill detail is not rendered to separate
// pages: it's populated client-side, in a slide-in panel, from
// search-index.json (see BuildSearchIndex).
func Render(skills []*skill.Skill, outputDir string) error {
	tmpl, err := template.ParseFS(templatesFS, "templates/index.html.tmpl")
	if err != nil {
		return fmt.Errorf("parse index template: %w", err)
	}

	if err := renderIndex(tmpl, skills, outputDir); err != nil {
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
		Title:   "Skill Indexer",
		Version: Version,
		Skills:  cards,
	}

	outPath := filepath.Join(outputDir, "index.html")
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("create output dir for %s: %w", outPath, err)
	}
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create %s: %w", outPath, err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("execute index template: %w", err)
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
