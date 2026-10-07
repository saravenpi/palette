package formats

import (
	"fmt"
	"strings"
)

// FormatTheme dispatches theme formatting to format-specific handlers based on the format string.
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