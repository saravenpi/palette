package cmd

import (
	"fmt"
	"os"

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
			formatted, err := palette.FormatTheme(theme, formatFlag)
			if err != nil {
				return err
			}

			if outputFlag != "" {
				if err := os.WriteFile(outputFlag, []byte(formatted), 0644); err != nil {
					return fmt.Errorf("failed to write output file: %w", err)
				}
				if term.IsTerminal(int(os.Stdout.Fd())) && !rawFlag {
					fmt.Println(palette.RenderPreview(theme, imagePath, modeFlag, lightFlag))
					fmt.Printf("Theme written to %s\n", outputFlag)
				}
				return nil
			}

			if previewFlag {
				fmt.Println(palette.RenderPreview(theme, imagePath, modeFlag, lightFlag))
				if !rawFlag {
					fmt.Println()
				}
			}

			fmt.Print(formatted)
			return nil
		},
	}

	cmd.Flags().StringVarP(&formatFlag, "format", "f", "ghostty", "output format (ghostty, kitty, alacritty, wezterm, foot, xresources, json, hex)")
	cmd.Flags().StringVarP(&modeFlag, "mode", "m", "duo", "palette mode (duo, ansi, dominant)")
	cmd.Flags().BoolVarP(&lightFlag, "light", "l", false, "generate light theme")
	cmd.Flags().BoolVarP(&previewFlag, "preview", "p", false, "show lipgloss preview")
	cmd.Flags().BoolVarP(&rawFlag, "raw", "r", false, "output raw config without preview")
	cmd.Flags().StringVarP(&outputFlag, "output", "o", "", "write theme to file")
	cmd.Flags().IntVar(&maxDimFlag, "max-dim", 400, "maximum image dimension before color extraction")

	cmd.AddCommand(newPreviewCommand())

	return cmd
}
