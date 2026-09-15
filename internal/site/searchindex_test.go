package site

import (
	"encoding/json"
	"testing"
)

func TestBuildSearchIndex(t *testing.T) {
	skills := fixtureSkills()

	data, err := BuildSearchIndex(skills)
	if err != nil {
		t.Fatalf("BuildSearchIndex() error = %v", err)
	}

	if !json.Valid(data) {
		t.Fatalf("BuildSearchIndex() did not produce valid JSON: %s", data)
	}

	var entries []searchEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("unmarshal search index: %v", err)
	}

	if len(entries) != len(skills) {
		t.Fatalf("got %d entries, want %d", len(entries), len(skills))
	}

	for i, s := range skills {
		if entries[i].Name != s.Name || entries[i].DirName != s.DirName {
			t.Errorf("entry %d = %+v, want name=%q dirName=%q", i, entries[i], s.Name, s.DirName)
		}
	}
}
