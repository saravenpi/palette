package formats

import (
	"fmt"
	"strings"
)

// formatWidgets renders the stylesheet widgets loads at startup. Every color is
// written twice: as a CSS custom property under `:root`, which GTK 4.16 and
// later — and any web consumer — reads with var(), and then as a @define-color
// declaration, which is all GTK 4.10 through 4.15 understand.
func formatWidgets(t Theme) string {
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
