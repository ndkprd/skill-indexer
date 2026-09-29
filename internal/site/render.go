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

// cardBadgeKeys lists the only metadata keys shown on a card, in order.
var cardBadgeKeys = []string{"version", "author"}

// Version is shown in the generated site's footer trademark line.
const Version = "0.1.1"

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

// Options configures the optional site features. The zero value keeps the
// default site: domain-root paths, zip install commands, no extra theming.
type Options struct {
	// BaseURL prefixes every root-relative reference (e.g. "/skills").
	BaseURL string
	// RepoURL, when set, switches install commands to the git repository
	// and adds an "install everything" banner.
	RepoURL string
	// SiteName, when set, shows a title row above the search field.
	SiteName string
	// ExtraCSS and ExtraJS are file paths copied into assets/extra and
	// linked after the built-in stylesheet and script.
	ExtraCSS []string
	ExtraJS  []string
}

type indexPageData struct {
	Title    string
	Version  string
	BaseURL  string
	RepoURL  string
	SiteName string
	ExtraCSS []string
	ExtraJS  []string
	Skills   []cardView
}

// Render writes the single-page static site (card grid plus its vendored
// static assets) under outputDir. Skill detail is not rendered to separate
// pages: it's populated client-side, in a slide-in panel, from
// search-index.json (see BuildSearchIndex). opts.BaseURL prefixes every
// root-relative reference the page emits (assets, wordmark link) so the
// site can be hosted under a subpath; leave it "" to keep the previous
// domain-root behavior.
func Render(skills []*skill.Skill, outputDir string, opts Options) error {
	tmpl, err := template.ParseFS(templatesFS, "templates/index.html.tmpl")
	if err != nil {
		return fmt.Errorf("parse index template: %w", err)
	}

	if err := renderIndex(tmpl, skills, outputDir, opts); err != nil {
		return err
	}
	if err := writeAssets(outputDir); err != nil {
		return err
	}
	return copyExtras(outputDir, opts)
}

func renderIndex(tmpl *template.Template, skills []*skill.Skill, outputDir string, opts Options) error {
	cards := make([]cardView, len(skills))
	for i, s := range skills {
		cards[i] = buildCardView(s)
	}
	sort.Slice(cards, func(i, j int) bool { return cards[i].Name < cards[j].Name })

	data := indexPageData{
		Title:    "Skill Indexer",
		Version:  Version,
		BaseURL:  opts.BaseURL,
		RepoURL:  opts.RepoURL,
		SiteName: opts.SiteName,
		ExtraCSS: basenames(opts.ExtraCSS),
		ExtraJS:  basenames(opts.ExtraJS),
		Skills:   cards,
	}
	if opts.SiteName != "" {
		data.Title = opts.SiteName
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
		Badges:      pickBadges(s.Metadata, cardBadgeKeys),
	}
}

// pickBadges returns the entries of metadata named in keys, in that order.
// Keys the skill doesn't define are skipped.
func pickBadges(metadata map[string]any, keys []string) []badge {
	var badges []badge
	for _, key := range keys {
		if v, ok := metadata[key]; ok {
			badges = append(badges, badge{Key: key, Value: formatMetadataValue(v)})
		}
	}
	return badges
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
