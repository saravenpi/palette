package formats

import (
	"fmt"
	"strings"
)

func formatTmux(t Theme) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("set -g status-style \"fg=%s,bg=%s\"\n", t.Foreground, t.Background))
	b.WriteString(fmt.Sprintf("set -g pane-border-style \"fg=%s,bg=%s\"\n", t.Palette[8], t.Background))
	b.WriteString(fmt.Sprintf("set -g pane-active-border-style \"fg=%s,bg=%s\"\n", t.Cursor, t.Background))
	b.WriteString(fmt.Sprintf("set -g mode-style \"fg=%s,bg=%s,bold\"\n", t.Background, t.Cursor))
	b.WriteString(fmt.Sprintf("set -g window-status-style \"fg=%s,bg=%s\"\n", t.Palette[7], t.Background))
	b.WriteString(fmt.Sprintf("set -g window-status-current-style \"fg=%s,bg=%s,bold\"\n", t.Cursor, t.Palette[0]))
	b.WriteString(fmt.Sprintf("set -g message-style \"fg=%s,bg=%s,bold\"\n", t.Foreground, t.Background))
	b.WriteString(fmt.Sprintf("set -g message-command-style \"fg=%s,bg=%s,bold\"\n", t.Cursor, t.Palette[0]))
	b.WriteString(fmt.Sprintf("set -g copy-mode-match-style \"fg=%s,bg=%s,bold\"\n", t.Background, t.Cursor))
	b.WriteString(fmt.Sprintf("set -g window-status-format \"#[fg=%s,bg=%s] #I #[fg=%s,bg=%s] #W #[fg=%s,bg=%s]\"\n", t.Palette[8], t.Background, t.Foreground, t.Palette[0], t.Palette[8], t.Background))
	b.WriteString(fmt.Sprintf("set -g window-status-current-format \"#[fg=%s,bg=%s,bold] #I #[fg=%s,bg=%s,bold] #W #[fg=%s,bg=%s,bold]\"\n", t.Cursor, t.Palette[8], t.Background, t.Cursor, t.Cursor, t.Palette[8]))
	b.WriteString(fmt.Sprintf("set -g status-left \"#[fg=%s,bg=%s,bold] #S #[fg=%s,bg=%s] \"\n", t.Background, t.Cursor, t.Cursor, t.Background))
	for i, c := range t.Palette {
		b.WriteString(fmt.Sprintf("set -g @palette_color%d \"%s\"\n", i, c))
	}
	b.WriteString(fmt.Sprintf("set -g @palette_bg \"%s\"\n", t.Background))
	b.WriteString(fmt.Sprintf("set -g @palette_fg \"%s\"\n", t.Foreground))
	b.WriteString(fmt.Sprintf("set -g @palette_cursor \"%s\"\n", t.Cursor))
	b.WriteString(fmt.Sprintf("set -g @paper_background \"%s\"\n", t.Background))
	b.WriteString(fmt.Sprintf("set -g @paper_foreground \"%s\"\n", t.Foreground))
	b.WriteString(fmt.Sprintf("set -g @paper_purple \"%s\"\n", t.Cursor))
	b.WriteString(fmt.Sprintf("set -g @paper_green \"%s\"\n", t.Palette[2]))
	b.WriteString(fmt.Sprintf("set -g @paper_orange \"%s\"\n", t.Palette[3]))
	b.WriteString(fmt.Sprintf("set -g @paper_yellow \"%s\"\n", t.Palette[3]))
	b.WriteString(fmt.Sprintf("set -g @paper_blue \"%s\"\n", t.Palette[4]))
	b.WriteString(fmt.Sprintf("set -g @paper_gray \"%s\"\n", t.Palette[8]))
	return b.String()
}