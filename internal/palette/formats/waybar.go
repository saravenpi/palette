package formats

import (
	"fmt"
	"strings"
)

// formatWaybar renders the stylesheet waybar imports. It carries the same two
// mechanisms as the widgets sheet: custom properties under `:root` for GTK 4.16
// and later (and for web consumers), then the @define-color declarations older
// GTK releases fall back to.
func formatWaybar(t Theme) string {
	var b strings.Builder
	b.WriteString("/* Custom properties for GTK 4.16 and later; the named color\n")
	b.WriteString("   declarations below are the fallback for older GTK releases. */\n\n")
	b.WriteString(":root {\n")
	b.WriteString(fmt.Sprintf("  --palette-bg: %s;\n", t.Background))
	b.WriteString(fmt.Sprintf("  --palette-fg: %s;\n", t.Foreground))
	b.WriteString(fmt.Sprintf("  --palette-cursor: %s;\n", t.Cursor))
	for i, c := range t.Palette {
		b.WriteString(fmt.Sprintf("  --palette-color%d: %s;\n", i, c))
	}
	b.WriteString("}\n\n")
	b.WriteString(fmt.Sprintf("@define-color palette_bg %s;\n", t.Background))
	b.WriteString(fmt.Sprintf("@define-color palette_fg %s;\n", t.Foreground))
	b.WriteString(fmt.Sprintf("@define-color palette_cursor %s;\n", t.Cursor))
	for i, c := range t.Palette {
		b.WriteString(fmt.Sprintf("@define-color palette_color%d %s;\n", i, c))
	}
	return b.String()
}
