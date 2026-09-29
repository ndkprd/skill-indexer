package cmd

import "testing"

func TestNormalizeRepoURL(t *testing.T) {
	tests := map[string]string{
		"":                                "",
		"  https://example.com/g/s.git  ": "https://example.com/g/s.git",
		"https://example.com/g/s.git/":    "https://example.com/g/s.git",
		"git@gitlab.com:group/skills.git": "git@gitlab.com:group/skills.git",
	}
	for in, want := range tests {
		if got := normalizeRepoURL(in); got != want {
			t.Errorf("normalizeRepoURL(%q) = %q, want %q", in, got, want)
		}
	}
}
