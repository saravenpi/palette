package formats

import (
	"fmt"
	"strings"
	"testing"
)

func TestFormatHyprland(t *testing.T) {
	theme := coherenceTheme()
	out, err := FormatTheme(theme, "hyprland")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(out, "#") {
		t.Errorf("hyprland output carries hex colors; Hyprland wants rgb():\n%s", out)
	}

	for i, c := range theme.Palette {
		want := fmt.Sprintf("$palette_color%d = rgb(%s)", i, strings.TrimPrefix(c, "#"))
		if !strings.Contains(out, want) {
			t.Errorf("hyprland output missing %q:\n%s", want, out)
		}
	}

	base := []struct{ name, value string }{
		{"$palette_bg", theme.Background},
		{"$palette_fg", theme.Foreground},
		{"$palette_cursor", theme.Cursor},
	}
	for _, b := range base {
		want := fmt.Sprintf("%s = rgb(%s)", b.name, strings.TrimPrefix(b.value, "#"))
		if !strings.Contains(out, want) {
			t.Errorf("hyprland output missing %q:\n%s", want, out)
		}
	}

	// The border roles name the accent, the dim grey and the urgency color, so
	// the window rules never point at a raw index.
	borders := []struct{ name, value string }{
		{"$active_border", theme.Palette[4]},
		{"$inactive_border", theme.Palette[8]},
		{"$urgent_border", theme.Palette[1]},
	}
	for _, b := range borders {
		want := fmt.Sprintf("%s = rgb(%s)", b.name, strings.TrimPrefix(b.value, "#"))
		if !strings.Contains(out, want) {
			t.Errorf("hyprland output missing %q:\n%s", want, out)
		}
	}
}
