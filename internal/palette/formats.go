package palette

import (
	"palette/internal/palette/formats"
)

// Theme represents a terminal color palette with 16 ANSI colors plus background,
// foreground, and cursor colors.
type Theme = formats.Theme

// Options configures how a theme is generated.
type Options = formats.Options

// FormatTheme dispatches theme formatting to format-specific handlers based on the format string.
func FormatTheme(t Theme, format string) (string, error) {
	return formats.FormatTheme(t, format)
}