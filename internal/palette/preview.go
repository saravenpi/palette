package palette

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/lucasb-eyer/go-colorful"
)

func ContrastColor(hex string) color.Color {
	c, err := colorful.Hex(hex)
	if err != nil {
		return lipgloss.Color("#ffffff")
	}
	_, _, l := c.Hsl()
	if l > 0.55 {
		return lipgloss.Color("#000000")
	}
	return lipgloss.Color("#ffffff")
}

func RenderPreview(t Theme, imagePath string, mode string, light bool) string {
	themeMode := "dark"
	if light {
		themeMode = "light"
	}

	accent1 := lipgloss.Color(t.Palette[1])
	accent2 := lipgloss.Color(t.Palette[4])
	bg := lipgloss.Color(t.Background)
	fg := lipgloss.Color(t.Foreground)
	dim := lipgloss.Color(t.Palette[8])

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ContrastColor(t.Palette[4])).
		Background(accent2).
		Padding(0, 1)

	tagStyle := lipgloss.NewStyle().
		Foreground(dim).
		Padding(0, 1)

	header := lipgloss.JoinHorizontal(lipgloss.Center,
		titleStyle.Render("PALETTE"),
		tagStyle.Render(fmt.Sprintf("%s • mode: %s • %s", filepath.Base(imagePath), mode, themeMode)),
	)

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accent1).
		Padding(1, 2)

	swatch := func(label string, hex string) string {
		col := lipgloss.Color(hex)
		block := lipgloss.NewStyle().Foreground(col).Render("████")
		pill := lipgloss.NewStyle().Background(col).Foreground(ContrastColor(hex)).Bold(true).Render(" " + hex + " ")
		text := lipgloss.NewStyle().Foreground(fg).Render(fmt.Sprintf("%-3s", label))
		return fmt.Sprintf("%s %s %s", text, block, pill)
	}

	normRow1 := fmt.Sprintf("  %s    %s",
		swatch("0", t.Palette[0]),
		swatch("1", t.Palette[1]),
	)
	normRow2 := fmt.Sprintf("  %s    %s",
		swatch("2", t.Palette[2]),
		swatch("3", t.Palette[3]),
	)
	normRow3 := fmt.Sprintf("  %s    %s",
		swatch("4", t.Palette[4]),
		swatch("5", t.Palette[5]),
	)
	normRow4 := fmt.Sprintf("  %s    %s",
		swatch("6", t.Palette[6]),
		swatch("7", t.Palette[7]),
	)

	brightRow1 := fmt.Sprintf("  %s    %s",
		swatch("8", t.Palette[8]),
		swatch("9", t.Palette[9]),
	)
	brightRow2 := fmt.Sprintf("  %s    %s",
		swatch("10", t.Palette[10]),
		swatch("11", t.Palette[11]),
	)
	brightRow3 := fmt.Sprintf("  %s    %s",
		swatch("12", t.Palette[12]),
		swatch("13", t.Palette[13]),
	)
	brightRow4 := fmt.Sprintf("  %s    %s",
		swatch("14", t.Palette[14]),
		swatch("15", t.Palette[15]),
	)

	specRow1 := fmt.Sprintf("  %s    %s",
		swatch("bg", t.Background),
		swatch("fg", t.Foreground),
	)
	specRow2 := fmt.Sprintf("  %s",
		swatch("cur", t.Cursor),
	)

	termBox := lipgloss.NewStyle().
		Background(bg).
		Foreground(fg).
		Padding(1, 2).
		Border(lipgloss.NormalBorder()).
		BorderForeground(dim)

	promptSym := lipgloss.NewStyle().Foreground(accent1).Render("❯ ")
	cmdText := lipgloss.NewStyle().Foreground(fg).Render("git status")
	branchText := lipgloss.NewStyle().Foreground(dim).Render("# On branch main")
	modText := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Palette[2])).Render("  modified:   main.go")
	untrackText := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Palette[1])).Render("  untracked:  " + filepath.Base(imagePath))

	termContent := lipgloss.JoinVertical(lipgloss.Left,
		promptSym+cmdText,
		branchText,
		modText,
		untrackText,
	)

	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		lipgloss.NewStyle().Bold(true).Foreground(fg).Render("Normal ANSI Colors:"),
		normRow1,
		normRow2,
		normRow3,
		normRow4,
		"",
		lipgloss.NewStyle().Bold(true).Foreground(fg).Render("Bright ANSI Colors:"),
		brightRow1,
		brightRow2,
		brightRow3,
		brightRow4,
		"",
		lipgloss.NewStyle().Bold(true).Foreground(fg).Render("Special Colors:"),
		specRow1,
		specRow2,
		"",
		lipgloss.NewStyle().Bold(true).Foreground(fg).Render("Terminal Preview:"),
		termBox.Render(termContent),
	)

	return cardStyle.Render(content)
}

func RenderStyledConfig(t Theme, format string) string {
	if strings.ToLower(format) != "ghostty" && format != "" {
		out, _ := FormatTheme(t, format)
		return out
	}

	var b strings.Builder
	prefixStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	for i, hex := range t.Palette {
		col := lipgloss.Color(hex)
		badge := lipgloss.NewStyle().Background(col).Foreground(ContrastColor(hex)).Bold(true).Render(" " + hex + " ")
		block := lipgloss.NewStyle().Foreground(col).Render("████")
		valStyle := lipgloss.NewStyle().Foreground(col).Bold(true)
		line := fmt.Sprintf("%s%s  %s  %s\n", prefixStyle.Render(fmt.Sprintf("palette = %2d=", i)), valStyle.Render(hex), block, badge)
		b.WriteString(line)
	}

	specials := []struct{ name, hex string }{
		{"background  ", t.Background},
		{"foreground  ", t.Foreground},
		{"cursor-color", t.Cursor},
	}
	for _, s := range specials {
		col := lipgloss.Color(s.hex)
		badge := lipgloss.NewStyle().Background(col).Foreground(ContrastColor(s.hex)).Bold(true).Render(" " + s.hex + " ")
		block := lipgloss.NewStyle().Foreground(col).Render("████")
		valStyle := lipgloss.NewStyle().Foreground(col).Bold(true)
		line := fmt.Sprintf("%s%s  %s  %s\n", prefixStyle.Render(fmt.Sprintf("%s = ", s.name)), valStyle.Render(s.hex), block, badge)
		b.WriteString(line)
	}

	return b.String()
}
