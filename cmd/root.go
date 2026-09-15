// Package cmd implements the skill-repo-store CLI.
package cmd

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
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
	Use:   "skill-repo-store",
	Short: "Generate a static skill marketplace site from a directory of Claude Code skills",
	Long: "skill-repo-store scans a directory of Claude Code skills, each with a SKILL.md\n" +
		"file, and generates a static, searchable marketplace site: a card grid, a\n" +
		"detail page per skill, zip downloads, and npx install commands.",
	RunE: runGenerate,
}

func runGenerate(cmd *cobra.Command, args []string) error {
	log := newLogger()
	log.Info().Str("event", "generate_start").Str("skill_dir", skillDir).Str("output_dir", outputDir).Msg("starting generation")
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
