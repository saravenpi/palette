package formats

import (
	"strings"
	"testing"
)

func coherenceTheme() Theme {
	return Theme{
		Background: "#19120c",
		Foreground: "#e8e1dd",
		Cursor:     "#ffc079",
		Palette: [16]string{
			"#261e15", "#ff8a70", "#a2b275", "#d6a163",
			"#5eb5e1", "#ca9dc7", "#79b7a8", "#ccc5c2",
			"#776b5f", "#ffcbc0", "#cede9e", "#ffcf99",
			"#a8deff", "#f9c8f4", "#a4e4d3", "#fcf5f4",
		},
	}
}

func TestFormatWaybar(t *testing.T) {
	theme := coherenceTheme()
	out, err := FormatTheme(theme, "waybar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertCustomProperties(t, "waybar", theme, out)
	assertNamedColors(t, "waybar", theme, out)
}

func TestAdaptersShareBackground(t *testing.T) {
	theme := coherenceTheme()
	checks := map[string]string{
		"kitty":   "background " + theme.Background,
		"foot":    "background=" + strings.TrimPrefix(theme.Background, "#"),
		"sway":    "set $palette_bg " + theme.Background,
		"tmux":    "@palette_bg \"" + theme.Background + "\"",
		"waybar":  "@define-color palette_bg " + theme.Background + ";",
		"widgets": "@define-color palette_bg " + theme.Background + ";",
		"vim":     "guifg=" + theme.Foreground + " guibg=" + theme.Background,
	}
	for format, want := range checks {
		out, err := FormatTheme(theme, format)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", format, err)
		}
		if !strings.Contains(out, want) {
			t.Errorf("%s output does not carry theme background %q:\n%s", format, theme.Background, out)
		}
	}
}
