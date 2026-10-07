package formats

import (
	"fmt"
	"strings"
)

func formatAlacritty(t Theme) string {
	var b strings.Builder
	b.WriteString("[colors.primary]\n")
	b.WriteString(fmt.Sprintf("background = \"%s\"\n", t.Background))
	b.WriteString(fmt.Sprintf("foreground = \"%s\"\n", t.Foreground))
	b.WriteString("\n[colors.cursor]\n")
	b.WriteString(fmt.Sprintf("cursor = \"%s\"\n", t.Cursor))
	b.WriteString("\n[colors.normal]\n")
	b.WriteString(fmt.Sprintf("black = \"%s\"\n", t.Palette[0]))
	b.WriteString(fmt.Sprintf("red = \"%s\"\n", t.Palette[1]))
	b.WriteString(fmt.Sprintf("green = \"%s\"\n", t.Palette[2]))
	b.WriteString(fmt.Sprintf("yellow = \"%s\"\n", t.Palette[3]))
	b.WriteString(fmt.Sprintf("blue = \"%s\"\n", t.Palette[4]))
	b.WriteString(fmt.Sprintf("magenta = \"%s\"\n", t.Palette[5]))
	b.WriteString(fmt.Sprintf("cyan = \"%s\"\n", t.Palette[6]))
	b.WriteString(fmt.Sprintf("white = \"%s\"\n", t.Palette[7]))
	b.WriteString("\n[colors.bright]\n")
	b.WriteString(fmt.Sprintf("black = \"%s\"\n", t.Palette[8]))
	b.WriteString(fmt.Sprintf("red = \"%s\"\n", t.Palette[9]))
	b.WriteString(fmt.Sprintf("green = \"%s\"\n", t.Palette[10]))
	b.WriteString(fmt.Sprintf("yellow = \"%s\"\n", t.Palette[11]))
	b.WriteString(fmt.Sprintf("blue = \"%s\"\n", t.Palette[12]))
	b.WriteString(fmt.Sprintf("magenta = \"%s\"\n", t.Palette[13]))
	b.WriteString(fmt.Sprintf("cyan = \"%s\"\n", t.Palette[14]))
	b.WriteString(fmt.Sprintf("white = \"%s\"\n", t.Palette[15]))
	return b.String()
}