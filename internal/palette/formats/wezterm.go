package formats

import (
	"fmt"
	"strings"
)

func formatWezterm(t Theme) string {
	var b strings.Builder
	b.WriteString("return {\n")
	b.WriteString(fmt.Sprintf("  background = '%s',\n", t.Background))
	b.WriteString(fmt.Sprintf("  foreground = '%s',\n", t.Foreground))
	b.WriteString(fmt.Sprintf("  cursor_bg = '%s',\n", t.Cursor))
	b.WriteString("  ansi = {\n")
	for i := 0; i < 8; i++ {
		b.WriteString(fmt.Sprintf("    '%s',\n", t.Palette[i]))
	}
	b.WriteString("  },\n")
	b.WriteString("  brights = {\n")
	for i := 8; i < 16; i++ {
		b.WriteString(fmt.Sprintf("    '%s',\n", t.Palette[i]))
	}
	b.WriteString("  },\n")
	b.WriteString("}\n")
	return b.String()
}