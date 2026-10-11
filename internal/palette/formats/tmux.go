package formats

import (
	"fmt"
	"strings"
)

func formatTmux(t Theme) string {
	return paletteLines(t) + paperVars(t)
}

func paletteLines(t Theme) string {
	var b strings.Builder
	for i, c := range t.Palette {
		b.WriteString(fmt.Sprintf("set -g @palette_color%d \"%s\"\n", i, c))
	}
	return b.String()
}

func paperVars(t Theme) string {
	return fmt.Sprintf(
		"set -g @paper_background \"%s\"\nset -g @paper_foreground \"%s\"\n"+
			"set -g @paper_purple \"%s\"\nset -g @paper_green \"%s\"\n"+
			"set -g @paper_orange \"%s\"\nset -g @paper_yellow \"%s\"\n"+
			"set -g @paper_blue \"%s\"\nset -g @paper_gray \"%s\"\nset -g @palette_bg \"%s\"\n"+
			"set -g @palette_fg \"%s\"\nset -g @palette_cursor \"%s\"\n",
		t.Background, t.Foreground, t.Cursor,
		t.Palette[2], t.Palette[3], t.Palette[3], t.Palette[4], t.Palette[8],
		t.Background, t.Foreground, t.Cursor,
	)
}
