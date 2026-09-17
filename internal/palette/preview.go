package palette

import (
	"fmt"
	"image/color"
	"path/filepath"

	"charm.land/lipgloss/v2"
)

func RenderPreview(t Theme, imagePath string, mode string, light bool) string {
	themeMode := "dark"
	if light {
		themeMode = "light"
	}

	accent1 := lipgloss.Color(t.Palette[1])
	accent2 := lipgloss.Color(t.Palette[4])
	curColor := lipgloss.Color(t.Cursor)
	bg := lipgloss.Color(t.Background)
	fg := lipgloss.Color(t.Foreground)
	dim := lipgloss.Color(t.Palette[8])

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ffffff")).
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

	swatch := func(label string, col color.Color, hex string) string {
		block := lipgloss.NewStyle().Foreground(col).Render("██")
		text := lipgloss.NewStyle().Foreground(fg).Render(fmt.Sprintf(" %-3s %s", label, hex))
		return block + text
	}

	normRow1 := fmt.Sprintf("%s   %s   %s   %s",
		swatch("0", lipgloss.Color(t.Palette[0]), t.Palette[0]),
		swatch("1", lipgloss.Color(t.Palette[1]), t.Palette[1]),
		swatch("2", lipgloss.Color(t.Palette[2]), t.Palette[2]),
		swatch("3", lipgloss.Color(t.Palette[3]), t.Palette[3]),
	)
	normRow2 := fmt.Sprintf("%s   %s   %s   %s",
		swatch("4", lipgloss.Color(t.Palette[4]), t.Palette[4]),
		swatch("5", lipgloss.Color(t.Palette[5]), t.Palette[5]),
		swatch("6", lipgloss.Color(t.Palette[6]), t.Palette[6]),
		swatch("7", lipgloss.Color(t.Palette[7]), t.Palette[7]),
	)

	brightRow1 := fmt.Sprintf("%s   %s   %s   %s",
		swatch("8", lipgloss.Color(t.Palette[8]), t.Palette[8]),
		swatch("9", lipgloss.Color(t.Palette[9]), t.Palette[9]),
		swatch("10", lipgloss.Color(t.Palette[10]), t.Palette[10]),
		swatch("11", lipgloss.Color(t.Palette[11]), t.Palette[11]),
	)
	brightRow2 := fmt.Sprintf("%s   %s   %s   %s",
		swatch("12", lipgloss.Color(t.Palette[12]), t.Palette[12]),
		swatch("13", lipgloss.Color(t.Palette[13]), t.Palette[13]),
		swatch("14", lipgloss.Color(t.Palette[14]), t.Palette[14]),
		swatch("15", lipgloss.Color(t.Palette[15]), t.Palette[15]),
	)

	specRow := fmt.Sprintf("%s   %s   %s",
		swatch("bg", bg, t.Background),
		swatch("fg", fg, t.Foreground),
		swatch("cur", curColor, t.Cursor),
	)

	promptSym := lipgloss.NewStyle().Foreground(accent1).Render("❯ ")
	cmdText := lipgloss.NewStyle().Foreground(fg).Render("git status")
	branchText := lipgloss.NewStyle().Foreground(dim).Render("# On branch main")
	modText := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Palette[2])).Render("  modified:   main.go")
	untrackText := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Palette[1])).Render("  untracked:  " + filepath.Base(imagePath))

	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		lipgloss.NewStyle().Bold(true).Foreground(fg).Render("Normal ANSI Colors:"),
		normRow1,
		normRow2,
		"",
		lipgloss.NewStyle().Bold(true).Foreground(fg).Render("Bright ANSI Colors:"),
		brightRow1,
		brightRow2,
		"",
		lipgloss.NewStyle().Bold(true).Foreground(fg).Render("Special Colors:"),
		specRow,
		"",
		lipgloss.NewStyle().Bold(true).Foreground(fg).Render("Terminal Preview:"),
		promptSym+cmdText,
		branchText,
		modText,
		untrackText,
	)

	return cardStyle.Render(content)
}
