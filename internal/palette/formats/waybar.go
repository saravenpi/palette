package formats

import (
	"fmt"
	"strings"
)

func formatWaybar(t Theme) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("@define-color palette_bg %s;\n", t.Background))
	b.WriteString(fmt.Sprintf("@define-color palette_fg %s;\n", t.Foreground))
	b.WriteString(fmt.Sprintf("@define-color palette_cursor %s;\n", t.Cursor))
	for i, c := range t.Palette {
		b.WriteString(fmt.Sprintf("@define-color palette_color%d %s;\n", i, c))
	}
	return b.String()
}