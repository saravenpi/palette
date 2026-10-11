package formats

import (
	"fmt"
	"strings"
)

func formatSway(t Theme) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("set $palette_bg %s\n", t.Background))
	b.WriteString(fmt.Sprintf("set $palette_fg %s\n", t.Foreground))
	b.WriteString(fmt.Sprintf("set $palette_cursor %s\n", t.Cursor))
	for i, c := range t.Palette {
		b.WriteString(fmt.Sprintf("set $palette_color%d %s\n", i, c))
	}
	b.WriteString(fmt.Sprintf(
		"client.focused %s %s %s %s %s\n",
		t.Cursor, t.Background, t.Foreground, t.Cursor, t.Cursor,
	))
	b.WriteString(fmt.Sprintf(
		"client.focused_inactive %s %s %s %s %s\n",
		t.Palette[8], t.Background, t.Palette[7], t.Palette[8], t.Palette[8],
	))
	b.WriteString(fmt.Sprintf(
		"client.unfocused %s %s %s %s %s\n",
		t.Palette[8], t.Background, t.Palette[7], t.Palette[8], t.Palette[8],
	))
	b.WriteString(fmt.Sprintf(
		"client.urgent %s %s %s %s %s\n",
		t.Palette[1], t.Background, t.Palette[1], t.Palette[1], t.Palette[1],
	))
	b.WriteString(fmt.Sprintf(
		"client.placeholder %s %s %s %s %s\n",
		t.Background, t.Background, t.Foreground, t.Background, t.Background,
	))
	return b.String()
}
