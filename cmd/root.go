// Package cmd implements the skillstore CLI.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"skillstore/internal/site"
	"skillstore/internal/skill"
)

var (
	skillDir  string
	outputDir string
)

func newLogger() zerolog.Logger {
	writer := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	return zerolog.New(writer).With().Timestamp().Logger()
}

var rootCmd = &cobra.Command{
	Use:   "skillstore",
	Short: "Generate a static skill marketplace site from a directory of Claude Code skills",
	Long: "skillstore scans a directory of Claude Code skills, each with a SKILL.md\n" +
		"file, and generates a static, searchable marketplace site: a card grid, a\n" +
		"detail page per skill, zip downloads, and npx install commands.",
	RunE: runGenerate,
}

func runGenerate(cmd *cobra.Command, args []string) error {
	log := newLogger()
	log.Info().Str("event", "generate_start").Str("skill_dir", skillDir).Str("output_dir", outputDir).Msg("starting generation")

	if err := os.RemoveAll(outputDir); err != nil {
		return fmt.Errorf("clear output dir %q: %w", outputDir, err)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir %q: %w", outputDir, err)
	}

	skills, warnings := skill.ScanDir(skillDir)
	for _, w := range warnings {
		log.Warn().Str("event", "skill_skipped").Err(w).Msg("skipping invalid skill")
	}
	log.Info().Str("event", "skills_scanned").Int("valid_count", len(skills)).Int("skipped_count", len(warnings)).Msg("scanned skill directory")

	if err := site.Render(skills, outputDir); err != nil {
		return fmt.Errorf("render site: %w", err)
	}

	if err := writeDownloads(skills, outputDir); err != nil {
		return err
	}

	if err := writeSearchIndex(skills, outputDir); err != nil {
		return err
	}

	log.Info().Str("event", "generate_done").Int("skill_count", len(skills)).Str("output_dir", outputDir).Msg("generation complete")
	return nil
}

func writeDownloads(skills []*skill.Skill, outputDir string) error {
	downloadsDir := filepath.Join(outputDir, "downloads")
	if err := os.MkdirAll(downloadsDir, 0o755); err != nil {
		return fmt.Errorf("create downloads dir %q: %w", downloadsDir, err)
	}
	for _, s := range skills {
		zipPath := filepath.Join(downloadsDir, s.DirName+".zip")
		if err := site.ZipSkillDir(s, zipPath); err != nil {
			return fmt.Errorf("zip skill %q: %w", s.DirName, err)
		}
	}
	return nil
}

func writeSearchIndex(skills []*skill.Skill, outputDir string) error {
	data, err := site.BuildSearchIndex(skills)
	if err != nil {
		return fmt.Errorf("build search index: %w", err)
	}
	indexPath := filepath.Join(outputDir, "search-index.json")
	if err := os.WriteFile(indexPath, data, 0o644); err != nil {
		return fmt.Errorf("write search index %q: %w", indexPath, err)
	}
	return nil
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVar(&skillDir, "skill-dir", "skills", "directory containing skill subdirectories to scan")
	rootCmd.Flags().StringVar(&outputDir, "output-dir", "public", "directory to write the generated static site into")
}
