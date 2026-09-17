package palette

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
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
	case "tmux":
		return formatTmux(t), nil
	case "nvim", "lua":
		return formatNvim(t), nil
	case "yaml", "yml":
		return formatYAML(t)
	case "json":
		return formatJSON(t)
	case "hex":
		return formatHex(t), nil
	default:
		return "", fmt.Errorf("unknown format: %s (supported: ghostty, kitty, alacritty, wezterm, foot, xresources, tmux, nvim, yaml, json, hex)", format)
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

func formatNvim(t Theme) string {
	var b strings.Builder
	b.WriteString("local M = {}\n\n")
	b.WriteString("M.colors = {\n")
	b.WriteString(fmt.Sprintf("  bg = \"%s\",\n", t.Background))
	b.WriteString(fmt.Sprintf("  fg = \"%s\",\n", t.Foreground))
	b.WriteString(fmt.Sprintf("  cursor = \"%s\",\n", t.Cursor))
	names := []string{
		"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white",
		"bright_black", "bright_red", "bright_green", "bright_yellow", "bright_blue", "bright_magenta", "bright_cyan", "bright_white",
	}
	for i, name := range names {
		b.WriteString(fmt.Sprintf("  %s = \"%s\",\n", name, t.Palette[i]))
	}
	b.WriteString("}\n\n")
	b.WriteString(`function M.load()
  if vim.g.colors_name then
    vim.cmd("hi clear")
  end
  vim.o.termguicolors = true
  vim.g.colors_name = "palette"
  local c = M.colors
  local hl = vim.api.nvim_set_hl

  hl(0, "Normal", { fg = c.fg, bg = c.bg })
  hl(0, "NormalFloat", { fg = c.fg, bg = c.bg })
  hl(0, "Cursor", { fg = c.bg, bg = c.cursor })
  hl(0, "CursorLine", { bg = c.black })
  hl(0, "CursorLineNr", { fg = c.cursor, bold = true })
  hl(0, "LineNr", { fg = c.bright_black })
  hl(0, "Visual", { fg = c.bg, bg = c.cursor })
  hl(0, "Search", { fg = c.bg, bg = c.yellow })
  hl(0, "IncSearch", { fg = c.bg, bg = c.cursor })
  hl(0, "StatusLine", { fg = c.fg, bg = c.black })
  hl(0, "StatusLineNC", { fg = c.bright_black, bg = c.bg })
  hl(0, "VertSplit", { fg = c.bright_black })
  hl(0, "WinSeparator", { fg = c.bright_black })
  hl(0, "Pmenu", { fg = c.fg, bg = c.black })
  hl(0, "PmenuSel", { fg = c.bg, bg = c.cursor })
  hl(0, "PmenuSbar", { bg = c.black })
  hl(0, "PmenuThumb", { bg = c.bright_black })

  hl(0, "Comment", { fg = c.bright_black, italic = true })
  hl(0, "Constant", { fg = c.yellow })
  hl(0, "String", { fg = c.green })
  hl(0, "Character", { fg = c.green })
  hl(0, "Number", { fg = c.yellow })
  hl(0, "Boolean", { fg = c.yellow })
  hl(0, "Float", { fg = c.yellow })

  hl(0, "Identifier", { fg = c.fg })
  hl(0, "Function", { fg = c.blue, bold = true })
  hl(0, "Statement", { fg = c.magenta, bold = true })
  hl(0, "Conditional", { fg = c.magenta, bold = true })
  hl(0, "Repeat", { fg = c.magenta, bold = true })
  hl(0, "Label", { fg = c.magenta })
  hl(0, "Operator", { fg = c.cyan })
  hl(0, "Keyword", { fg = c.magenta, bold = true })
  hl(0, "Exception", { fg = c.red, bold = true })
  hl(0, "PreProc", { fg = c.cyan })
  hl(0, "Type", { fg = c.yellow })
  hl(0, "Special", { fg = c.blue })
  hl(0, "Underlined", { underline = true })
  hl(0, "Error", { fg = c.red, bg = c.bg })
  hl(0, "Todo", { fg = c.yellow, bg = c.bg, bold = true })

  hl(0, "DiagnosticError", { fg = c.red })
  hl(0, "DiagnosticWarn", { fg = c.yellow })
  hl(0, "DiagnosticInfo", { fg = c.blue })
  hl(0, "DiagnosticHint", { fg = c.cyan })

  hl(0, "@variable", { fg = c.fg })
  hl(0, "@function", { fg = c.blue, bold = true })
  hl(0, "@keyword", { fg = c.magenta, bold = true })
  hl(0, "@string", { fg = c.green })
  hl(0, "@comment", { fg = c.bright_black, italic = true })
end

return M
`)
	return b.String()
}

func formatYAML(t Theme) (string, error) {
	bytes, err := yaml.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ParseYAML(data []byte) (Theme, error) {
	var t Theme
	err := yaml.Unmarshal(data, &t)
	if err != nil {
		return Theme{}, err
	}
	return t, nil
}
