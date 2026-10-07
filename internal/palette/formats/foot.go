package formats

import (
	"fmt"
	"strings"
)

func formatFoot(t Theme) string {
	var b strings.Builder
	b.WriteString("[colors]\n")
	b.WriteString(fmt.Sprintf("background=%s\n", strings.TrimPrefix(t.Background, "#")))
	b.WriteString(fmt.Sprintf("foreground=%s\n", strings.TrimPrefix(t.Foreground, "#")))
	for i := 0; i < 8; i++ {
		b.WriteString(fmt.Sprintf("regular%d=%s\n", i, strings.TrimPrefix(t.Palette[i], "#")))
	}
	for i := 8; i < 16; i++ {
		b.WriteString(fmt.Sprintf("bright%d=%s\n", i-8, strings.TrimPrefix(t.Palette[i], "#")))
	}
	return b.String()
}