package formats

import (
	"strings"
)

func formatHex(t Theme) string {
	var b strings.Builder
	for _, c := range t.Palette {
		b.WriteString(c + "\n")
	}
	return b.String()
}