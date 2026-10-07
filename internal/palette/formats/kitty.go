package formats

import (
	"fmt"
	"strings"
)

func formatKitty(t Theme) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("background %s\n", t.Background))
	b.WriteString(fmt.Sprintf("foreground %s\n", t.Foreground))
	b.WriteString(fmt.Sprintf("cursor %s\n", t.Cursor))
	for i, c := range t.Palette {
		b.WriteString(fmt.Sprintf("color%d %s\n", i, c))
	}
	return b.String()
}