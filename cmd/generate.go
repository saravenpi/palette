package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"palette/internal/palette"
)

type generateFlags struct {
	format  string
	mode    string
	light   bool
	preview bool
	raw     bool
	output  string
	maxDim  int
	sync    bool
}

func newGenerateCmd() *cobra.Command {
	var f generateFlags

	cmd := &cobra.Command{
		Use:   "generate <image-path>",
		Short: "Generate a terminal color palette from an image",
		Long: `palette generate extracts dominant and accent colors from an image to generate harmonious
16-color ANSI terminal palettes and themes for Ghostty, Kitty, Alacritty, WezTerm, and more.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerate(args[0], f)
		},
	}

	cmd.Flags().StringVarP(&f.format, "format", "f", "ghostty",
		"output format: ghostty, kitty, alacritty, wezterm, foot, sway,\n\t"+
			"xresources, tmux, nvim, waybar, widgets, vim, yaml, json, hex")
	cmd.Flags().StringVarP(&f.mode, "mode", "m", "tonal",
		"palette mode: tonal, wal, vibrant, expressive, duo,\n\tmaterial, dominant")
	cmd.Flags().BoolVarP(&f.light, "light", "l", false, "generate light theme")
	cmd.Flags().BoolVarP(&f.preview, "preview", "p", false, "show lipgloss preview card with actual colors")
	cmd.Flags().BoolVarP(&f.raw, "raw", "r", false, "output raw config without styling")
	cmd.Flags().StringVarP(&f.output, "output", "o", "", "write theme to file")
	cmd.Flags().IntVar(&f.maxDim, "max-dim", 400, "maximum image dimension before color extraction")
	cmd.Flags().BoolVarP(&f.sync, "sync", "s", false,
		"save theme to config and sync Kitty, Tmux, Neovim, Sway, Foot, Waybar, Widgets, Vim")

	return cmd
}

func runGenerate(imagePath string, f generateFlags) error {
	img, err := palette.LoadImage(imagePath)
	if err != nil {
		return err
	}

	scaled := palette.Downscale(img, f.maxDim)
	colors, dom, err := palette.ExtractColors(scaled, 16)
	if err != nil {
		return err
	}

	opts := palette.Options{Mode: f.mode, Light: f.light, MaxDim: f.maxDim}
	theme := palette.GenerateTheme(colors, dom, opts)
	rawFormatted, err := palette.FormatTheme(theme, f.format)
	if err != nil {
		return err
	}

	if f.sync {
		return runSync(theme, imagePath, f, rawFormatted)
	}
	return runOutput(theme, f, rawFormatted, imagePath)
}

func runSync(theme palette.Theme, imagePath string, f generateFlags, rawFormatted string) error {
	palettePath, err := palette.DefaultPalettePath()
	if err != nil {
		return err
	}
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

func runOutput(theme palette.Theme, f generateFlags, rawFormatted string, imagePath string) error {
	isTTY := term.IsTerminal(int(os.Stdout.Fd()))

	if f.output != "" {
		if err := os.WriteFile(f.output, []byte(rawFormatted), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		writePostOutput(theme, f, imagePath)
		return nil
	}

	if (!isTTY && !f.preview) || f.raw {
		fmt.Print(rawFormatted)
		return nil
	}

	if f.preview {
		fmt.Println(palette.RenderPreview(theme, imagePath, f.mode, f.light))
		fmt.Println()
	}

	fmt.Print(palette.RenderStyledConfig(theme, f.format))
	return nil
}

func writePostOutput(theme palette.Theme, f generateFlags, imagePath string) {
	printPreviewIfSet(theme, imagePath, f.mode, f.light, f.preview)
	isTTY := term.IsTerminal(int(os.Stdout.Fd()))
	if !isTTY || f.raw {
		return
	}
	fmt.Print(palette.RenderStyledConfig(theme, f.format))
	fmt.Printf("\nTheme written to %s\n", f.output)
}

func printPreviewIfSet(theme palette.Theme, imagePath, mode string, light, preview bool) {
	if !preview {
		return
	}
	fmt.Println(palette.RenderPreview(theme, imagePath, mode, light))
	fmt.Println()
}
