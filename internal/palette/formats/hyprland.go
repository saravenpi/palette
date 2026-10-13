package formats

import (
	"fmt"
	"strings"
)

// formatHyprland renders the colour fragment a Hyprland config sources.
//
// Hyprland spells a colour `rgb(RRGGBB)` (or `rgba(RRGGBBAA)`), not `#RRGGBB`,
// so every value is converted: a `$palette_color*` variable can then be dropped
// straight into a `general { col.active_border = ... }` line. The three border
// variables name the roles a bare palette leaves implicit — the focused window,
// a dimmed one and an urgent one — so the window rules never hard-code an index
// a palette refresh could change.
func formatHyprland(t Theme) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("$palette_bg = %s\n", hyprColor(t.Background)))
	b.WriteString(fmt.Sprintf("$palette_fg = %s\n", hyprColor(t.Foreground)))
	b.WriteString(fmt.Sprintf("$palette_cursor = %s\n", hyprColor(t.Cursor)))
	for i, c := range t.Palette {
		b.WriteString(fmt.Sprintf("$palette_color%d = %s\n", i, hyprColor(c)))
	}
	b.WriteString(fmt.Sprintf("$active_border = %s\n", hyprColor(t.Palette[4])))
	b.WriteString(fmt.Sprintf("$inactive_border = %s\n", hyprColor(t.Palette[8])))
	b.WriteString(fmt.Sprintf("$urgent_border = %s\n", hyprColor(t.Palette[1])))
	return b.String()
}

// hyprColor converts a `#RRGGBB` hex string to Hyprland's `rgb(RRGGBB)` form.
// A value that is already in Hyprland's syntax is returned unchanged.
func hyprColor(hex string) string {
	if strings.HasPrefix(hex, "rgb") {
		return hex
	}
	return fmt.Sprintf("rgb(%s)", strings.TrimPrefix(hex, "#"))
}
