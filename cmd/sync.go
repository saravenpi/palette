package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"palette/internal/palette"
)

func newSyncCommand() *cobra.Command {
	var fileFlag string
	var noReloadFlag bool

	cmd := &cobra.Command{
		Use:   "sync [palette.yml]",
		Short: "Sync palette to Kitty, Tmux, and Neovim",
		Long: `sync reads a palette YAML file (defaults to ~/.palette.yml), compiles native config
fragments for Kitty, Tmux, and Neovim, writes them to their config paths, and signals running
applications to hot-reload.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := fileFlag
			if len(args) > 0 {
				path = args[0]
			}
			if path == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					return err
				}
				path = filepath.Join(home, ".palette.yml")
			}

			theme, err := palette.LoadPaletteFile(path)
			if err != nil {
				return err
			}

			results, err := palette.Sync(theme, !noReloadFlag)
			if err != nil {
				return err
			}

			for _, res := range results {
				fmt.Println(res)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&fileFlag, "file", "f", "", "path to palette file (default ~/.palette.yml)")
	cmd.Flags().BoolVar(&noReloadFlag, "no-reload", false, "write configs without sending reload signals")

	return cmd
}
