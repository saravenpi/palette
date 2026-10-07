package formats

import (
	"fmt"
	"strings"
)

func formatGhostty(t Theme) string {
	var b strings.Builder
	for i, c := range t.Palette {
		b.WriteString(fmt.Sprintf("palette = %d=%s\n", i, c))
	}
	b.WriteString(fmt.Sprintf("background = %s\n", t.Background))
	b.WriteString(fmt.Sprintf("foreground = %s\n", t.Foreground))
	b.WriteString(fmt.Sprintf("cursor-color = %s\n", t.Cursor))
	return b.String()
}