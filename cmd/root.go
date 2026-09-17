package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"palette/internal/palette"
)

var (
	formatFlag  string
	modeFlag    string
	lightFlag   bool
	previewFlag bool
	rawFlag     bool
	outputFlag  string
	maxDimFlag  int
	syncFlag    bool
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "palette <image-path>",
		Short: "Generate terminal color palettes from images",
		Long: `palette extracts dominant and accent colors from an image to generate harmonious
16-color ANSI terminal palettes and themes for Ghostty, Kitty, Alacritty, WezTerm, and more.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			imagePath := args[0]

			img, err := palette.LoadImage(imagePath)
			if err != nil {
				return err
			}

			scaled := palette.Downscale(img, maxDimFlag)
			colors, dom, err := palette.ExtractColors(scaled, 16)
			if err != nil {
				return err
			}

			opts := palette.Options{
				Mode:   modeFlag,
				Light:  lightFlag,
				MaxDim: maxDimFlag,
			}

			theme := palette.GenerateTheme(colors, dom, opts)
			rawFormatted, err := palette.FormatTheme(theme, formatFlag)
			if err != nil {
				return err
			}

			if syncFlag {
				home, err := os.UserHomeDir()
				if err != nil {
					return err
				}
				palettePath := filepath.Join(home, ".palette.yml")
				if err := palette.SavePaletteFile(palettePath, theme); err != nil {
					return err
				}
				syncResults, err := palette.Sync(theme, true)
				if err != nil {
					return err
				}
				fmt.Printf("saved to %s\n", palettePath)
				for _, r := range syncResults {
					fmt.Println(r)
				}
				return nil
			}

			if outputFlag != "" {
				if err := os.WriteFile(outputFlag, []byte(rawFormatted), 0644); err != nil {
					return fmt.Errorf("failed to write output file: %w", err)
				}
				if previewFlag {
					fmt.Println(palette.RenderPreview(theme, imagePath, modeFlag, lightFlag))
					fmt.Println()
				}
				if term.IsTerminal(int(os.Stdout.Fd())) && !rawFlag {
					fmt.Print(palette.RenderStyledConfig(theme, formatFlag))
					fmt.Printf("\nTheme written to %s\n", outputFlag)
				}
				return nil
			}

			isTTY := term.IsTerminal(int(os.Stdout.Fd()))
			if (!isTTY && !previewFlag) || rawFlag {
				fmt.Print(rawFormatted)
				return nil
			}

			if previewFlag {
				fmt.Println(palette.RenderPreview(theme, imagePath, modeFlag, lightFlag))
				fmt.Println()
			}

			fmt.Print(palette.RenderStyledConfig(theme, formatFlag))
			return nil
		},
	}

	cmd.Flags().StringVarP(&formatFlag, "format", "f", "ghostty", "output format (ghostty, kitty, alacritty, wezterm, foot, xresources, tmux, nvim, yaml, json, hex)")
	cmd.Flags().StringVarP(&modeFlag, "mode", "m", "material", "palette mode (material, vibrant, expressive, tonal, duo, ansi, dominant)")
	cmd.Flags().BoolVarP(&lightFlag, "light", "l", false, "generate light theme")
	cmd.Flags().BoolVarP(&previewFlag, "preview", "p", false, "show lipgloss preview card with actual colors")
	cmd.Flags().BoolVarP(&rawFlag, "raw", "r", false, "output raw config without styling")
	cmd.Flags().StringVarP(&outputFlag, "output", "o", "", "write theme to file")
	cmd.Flags().IntVar(&maxDimFlag, "max-dim", 400, "maximum image dimension before color extraction")
	cmd.Flags().BoolVarP(&syncFlag, "sync", "s", false, "save theme to ~/.palette.yml and sync Kitty, Tmux, Neovim")

	cmd.AddCommand(newPreviewCommand())
	cmd.AddCommand(newSyncCommand())

	return cmd
}
