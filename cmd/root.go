package cmd

import (
	"context"
	"os"

	"charm.land/fang/v2"
	"github.com/spf13/cobra"
)

var version = "dev"

func newRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "palette",
		Short: "Generate terminal color palettes from images",
		Long: `palette extracts dominant and accent colors from an image to generate harmonious
16-color ANSI terminal palettes and themes for Ghostty, Kitty, Alacritty, WezTerm, and more.

Without a subcommand it prints this help.`,
		DisableAutoGenTag: true,
	}
	rootCmd.SetVersionTemplate("{{.Name}} {{.Version}}\n")
	rootCmd.AddCommand(newGenerateCmd())
	rootCmd.AddCommand(newPreviewCommand())
	rootCmd.AddCommand(newSyncCommand())
	return rootCmd
}

// Execute runs the palette command tree through fang (styled help and errors,
// --version, man and completion commands) and exits non-zero on error. Fang
// already renders the styled error, so the wrapper only forwards the exit code.
func Execute() {
	if fang.Execute(context.Background(), newRootCommand(), fang.WithVersion(version)) != nil {
		os.Exit(1)
	}
}
