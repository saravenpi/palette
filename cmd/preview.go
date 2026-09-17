package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"palette/internal/palette"
)

func newPreviewCommand() *cobra.Command {
	var (
		mode  string
		light bool
	)

	cmd := &cobra.Command{
		Use:   "preview <image-path>",
		Short: "Display a visual terminal preview of the palette",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			imagePath := args[0]

			img, err := palette.LoadImage(imagePath)
			if err != nil {
				return err
			}

			scaled := palette.Downscale(img, 400)
			colors, dom, err := palette.ExtractColors(scaled, 16)
			if err != nil {
				return err
			}

			opts := palette.Options{
				Mode:  mode,
				Light: light,
			}

			theme := palette.GenerateTheme(colors, dom, opts)
			fmt.Println(palette.RenderPreview(theme, imagePath, mode, light))
			return nil
		},
	}

	cmd.Flags().StringVarP(&mode, "mode", "m", "tonal", "palette mode (tonal, material, vibrant, expressive, duo, ansi, dominant)")
	cmd.Flags().BoolVarP(&light, "light", "l", false, "generate light theme")

	return cmd
}
