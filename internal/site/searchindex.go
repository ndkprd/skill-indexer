package site

import (
	"encoding/json"
	"fmt"

	"skill-indexer/internal/skill"
)

// searchEntry is one row of the client-side search index. DirName lets the
// browser correlate a match back to its card in the rendered grid, and also
// drives the slide-in detail panel: ZipPath and Metadata carry everything
// the panel needs, so no per-skill page or extra fetch is required.
type searchEntry struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	DirName     string         `json:"dirName"`
	ZipPath     string         `json:"zipPath"`
	Metadata    map[string]any `json:"metadata"`
}

// BuildSearchIndex serializes skills into the JSON array consumed by the
// client-side Fuse.js search and the detail panel on the index page.
// baseURL prefixes each entry's ZipPath so downloads resolve correctly when
// the site is hosted under a subpath; pass "" to keep root-relative paths.
func BuildSearchIndex(skills []*skill.Skill, baseURL string) ([]byte, error) {
	entries := make([]searchEntry, len(skills))
	for i, s := range skills {
		entries[i] = searchEntry{
			Name:        s.Name,
			Description: s.Description,
			DirName:     s.DirName,
			ZipPath:     baseURL + "/downloads/" + s.DirName + ".zip",
			Metadata:    s.Metadata,
		}
	}

	data, err := json.Marshal(entries)
	if err != nil {
		return nil, fmt.Errorf("marshal search index: %w", err)
	}
	return data, nil
}
