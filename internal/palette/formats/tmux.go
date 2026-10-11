package formats

import (
	"fmt"
	"strings"
)

func formatTmux(t Theme) string {
	paletteColors := paletteLines(t)
	styles := []string{
		fmt.Sprintf("set -g status-style \"fg=%s,bg=%s\"", t.Foreground, t.Background),
		fmt.Sprintf("set -g pane-border-style \"fg=%s,bg=%s\"", t.Palette[8], t.Background),
		fmt.Sprintf("set -g pane-active-border-style \"fg=%s,bg=%s\"", t.Cursor, t.Background),
		fmt.Sprintf("set -g mode-style \"fg=%s,bg=%s,bold\"", t.Background, t.Cursor),
		fmt.Sprintf("set -g window-status-style \"fg=%s,bg=%s\"", t.Palette[7], t.Background),
		fmt.Sprintf("set -g window-status-current-style \"fg=%s,bg=%s,bold\"", t.Cursor, t.Palette[0]),
		fmt.Sprintf("set -g message-style \"fg=%s,bg=%s,bold\"", t.Foreground, t.Background),
		fmt.Sprintf("set -g message-command-style \"fg=%s,bg=%s,bold\"", t.Cursor, t.Palette[0]),
		fmt.Sprintf("set -g copy-mode-match-style \"fg=%s,bg=%s,bold\"", t.Background, t.Cursor),
		fmt.Sprintf(
			"set -g window-status-format \"#[fg=%s,bg=%s] #I #[fg=%s,bg=%s] #W #[fg=%s,bg=%s]\"",
			t.Palette[8], t.Background, t.Foreground, t.Palette[0], t.Palette[8], t.Background,
		),
		fmt.Sprintf(
			"set -g window-status-current-format \"#[fg=%s,bg=%s,bold] #I #[fg=%s,bg=%s,bold] #W #[fg=%s,bg=%s,bold]\"",
			t.Cursor, t.Palette[8], t.Background, t.Cursor, t.Cursor, t.Palette[8],
		),
		fmt.Sprintf(
			"set -g status-left \"#[fg=%s,bg=%s,bold] #S #[fg=%s,bg=%s] \"",
			t.Background, t.Cursor, t.Cursor, t.Background,
		),
		paperVars(t),
	}
	return strings.Join(styles, "\n") + "\n" + paletteColors
}

func paletteLines(t Theme) string {
	var b strings.Builder
	for i, c := range t.Palette {
		b.WriteString(fmt.Sprintf("set -g @palette_color%d \"%s\"\n", i, c))
	}
	return b.String()
}

func paperVars(t Theme) string {
	return fmt.Sprintf(
		"set -g @paper_background \"%s\"\nset -g @paper_foreground \"%s\"\n"+
			"set -g @paper_purple \"%s\"\nset -g @paper_green \"%s\"\n"+
			"set -g @paper_orange \"%s\"\nset -g @paper_yellow \"%s\"\n"+
			"set -g @paper_blue \"%s\"\nset -g @paper_gray \"%s\"\nset -g @palette_bg \"%s\"\n"+
			"set -g @palette_fg \"%s\"\nset -g @palette_cursor \"%s\"\n",
		t.Background, t.Foreground, t.Cursor,
		t.Palette[2], t.Palette[3], t.Palette[3], t.Palette[4], t.Palette[8],
		t.Background, t.Foreground, t.Cursor,
	)
}