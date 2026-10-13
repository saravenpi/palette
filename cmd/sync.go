package cmd

import (
	"github.com/spf13/cobra"

	"palette/internal/palette"
)

// resolvePalettePath determines the palette file to sync. It checks positional
// args first, then the --file flag, and finally falls back to the default palette path.
func resolvePalettePath(fileFlag string, args []string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}
	if fileFlag != "" {
		return fileFlag, nil
	}
	return palette.DefaultPalettePath()
}

// printSyncResults prints each line reported by palette.Sync.
func printSyncResults(results []string) {
	for _, r := range results {
		println(r)
	}
}

// runSyncCommand executes the sync command: resolve the palette path,
// load the theme, sync it, and print results.
func runSyncCommand(cmd *cobra.Command, args []string) error {
	fileFlag, _ := cmd.Flags().GetString("file")
	path, err := resolvePalettePath(fileFlag, args)
	if err != nil {
		return err
	}
	theme, err := palette.LoadPaletteFile(path)
	if err != nil {
		return err
	}
	noReload, _ := cmd.Flags().GetBool("no-reload")
	results, err := palette.Sync(theme, !noReload)
	if err != nil {
		return err
	}
	printSyncResults(results)
	return nil
}

func newSyncCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync [palette.yml]",
		Short: "Sync palette to Kitty, Tmux, Neovim, Sway, Foot, Waybar, Widgets, and Vim",
		Long: `sync reads a palette YAML file (defaults to ~/.config/palette/config.yml or ~/.palette.yml),
compiles native config fragments for Kitty, Tmux, Neovim, Sway, Foot, Waybar, Widgets, and Vim,
writes them to their config paths, and signals running applications to hot-reload.`,
		Args: cobra.MaximumNArgs(1),
		RunE: runSyncCommand,
	}
	cmd.Flags().StringVarP(
		new(string),
		"file",
		"f",
		"",
		"path to palette file (default ~/.config/palette/config.yml or ~/.palette.yml)",
	)
	cmd.Flags().BoolVar(new(bool), "no-reload", false, "write configs without sending reload signals")
	return cmd
}
