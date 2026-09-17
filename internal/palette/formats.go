package palette

import (
	"encoding/json"
	"fmt"
	"strings"
)

func FormatTheme(t Theme, format string) (string, error) {
	switch strings.ToLower(format) {
	case "ghostty", "":
		return formatGhostty(t), nil
	case "kitty":
		return formatKitty(t), nil
	case "alacritty":
		return formatAlacritty(t), nil
	case "wezterm":
		return formatWezterm(t), nil
	case "foot":
		return formatFoot(t), nil
	case "xresources":
		return formatXresources(t), nil
	case "json":
		return formatJSON(t)
	case "hex":
		return formatHex(t), nil
	default:
		return "", fmt.Errorf("unknown format: %s (supported: ghostty, kitty, alacritty, wezterm, foot, xresources, json, hex)", format)
	}
}

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

func formatXresources(t Theme) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("*.background: %s\n", t.Background))
	b.WriteString(fmt.Sprintf("*.foreground: %s\n", t.Foreground))
	b.WriteString(fmt.Sprintf("*.cursorColor: %s\n", t.Cursor))
	for i, c := range t.Palette {
		b.WriteString(fmt.Sprintf("*.color%d: %s\n", i, c))
	}
	return b.String()
}

func formatJSON(t Theme) (string, error) {
	bytes, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes) + "\n", nil
}

func formatHex(t Theme) string {
	var b strings.Builder
	for _, c := range t.Palette {
		b.WriteString(c + "\n")
	}
	return b.String()
}
