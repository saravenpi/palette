package palette

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"palette/internal/palette/formats"
)

type syncTarget struct {
	name   string
	format string
	path   string
}

func syncTargets(home string) []syncTarget {
	return []syncTarget{
		{name: "kitty", format: "kitty", path: filepath.Join(home, ".config", "kitty", "palette.conf")},
		{name: "tmux", format: "tmux", path: filepath.Join(home, ".config", "tmux", "palette.conf")},
		{name: "nvim", format: "nvim", path: filepath.Join(home, ".config", "nvim", "lua", "palette.lua")},
		{name: "sway", format: "sway", path: filepath.Join(home, ".config", "sway", "palette.conf")},
		{name: "foot", format: "foot", path: filepath.Join(home, ".config", "foot", "palette.ini")},
		{name: "waybar", format: "waybar", path: filepath.Join(home, ".config", "waybar", "palette.css")},
		{name: "widgets", format: "widgets", path: filepath.Join(home, ".config", "widgets", "palette.css")},
		{name: "vim", format: "vim", path: filepath.Join(home, ".vim", "colors", "palette.vim")},
	}
}

// DefaultPalettePath returns the resolved path to the user's palette file.
func DefaultPalettePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(home, ".config", "palette", "config.yml")
	if _, err := os.Stat(configPath); err == nil {
		return configPath, nil
	}
	legacyPath := filepath.Join(home, ".palette.yml")
	if _, err := os.Stat(legacyPath); err == nil {
		return legacyPath, nil
	}
	return configPath, nil
}

// LoadPaletteFile reads and parses a theme YAML file at the given path.
func LoadPaletteFile(path string) (Theme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Theme{}, fmt.Errorf("palette file not found at %s: %w", path, err)
	}
	return formats.ParseYAML(data)
}

// SavePaletteFile serializes a theme to YAML and writes it to path.
func SavePaletteFile(path string, t Theme) error {
	data, err := formats.FormatTheme(t, "yaml")
	if err != nil {
		return err
	}
	if err := ensureDir(filepath.Dir(path)); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(data), 0644)
}

// ensureDir makes sure dir exists as a directory, creating it and its parents
// only when they are missing. os.MkdirAll fails with a "file exists" error when
// dir is a symlink whose target is missing — the layout dotfile managers leave
// behind — so the dangling link is resolved and the directory it points at is
// created instead.
func ensureDir(dir string) error {
	info, err := os.Stat(dir)
	if err == nil {
		if info.IsDir() {
			return nil
		}
		return fmt.Errorf("%s exists and is not a directory", dir)
	}
	if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(resolveSymlink(dir), 0755)
}

// resolveSymlink follows the symlink chain at path and returns the final target,
// whether or not that target exists. A path that is not a symlink is returned
// unchanged.
func resolveSymlink(path string) string {
	for range 32 {
		info, err := os.Lstat(path)
		if err != nil || info.Mode()&os.ModeSymlink == 0 {
			return path
		}
		target, err := os.Readlink(path)
		if err != nil {
			return path
		}
		if !filepath.IsAbs(target) {
			target = filepath.Join(filepath.Dir(path), target)
		}
		path = target
	}
	return path
}

// Sync compiles theme snippets for supported applications and hot-reloads them.
func Sync(t Theme, reload bool) ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	var results []string
	for _, tgt := range syncTargets(home) {
		out, err := formats.FormatTheme(t, tgt.format)
		if err != nil {
			return nil, err
		}
		if err := writeFormatted(tgt.path, out); err != nil {
			return nil, err
		}
		results = append(results, fmt.Sprintf("%-5s -> %s", tgt.name, tgt.path))
	}
	nvimColors := filepath.Join(home, ".config", "nvim", "colors", "palette.lua")
	if err := writeFormatted(nvimColors, "package.loaded[\"palette\"] = nil\nrequire(\"palette\").load()\n"); err != nil {
		return nil, err
	}
	if reload {
		results = append(results, reloadApplications(filepath.Join(home, ".config", "tmux", "palette.conf"))...)
	}
	return results, nil
}

func writeFormatted(path string, content string) error {
	if err := ensureDir(filepath.Dir(path)); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func reloadApplications(tmuxConf string) []string {
	var msgs []string
	uid := strconv.Itoa(os.Getuid())
	if err := exec.Command("pkill", "-SIGUSR1", "-u", uid, "kitty").Run(); err == nil {
		msgs = append(msgs, "reloaded kitty via SIGUSR1")
	}
	if err := exec.Command("tmux", "source-file", tmuxConf).Run(); err == nil {
		msgs = append(msgs, "reloaded tmux session")
	}
	if err := exec.Command("swaymsg", "reload").Run(); err == nil {
		msgs = append(msgs, "reloaded sway session")
	}
	if err := exec.Command("pkill", "-SIGUSR2", "-u", uid, "waybar").Run(); err == nil {
		msgs = append(msgs, "reloaded waybar via SIGUSR2")
	}
	reloadNvimSockets()
	return msgs
}

func reloadNvimSockets() {
	var sockets []string
	tmpSockets, _ := filepath.Glob("/tmp/nvim*/*")
	sockets = append(sockets, tmpSockets...)
	xdgRuntime := os.Getenv("XDG_RUNTIME_DIR")
	if xdgRuntime != "" {
		xdgSockets, _ := filepath.Glob(filepath.Join(xdgRuntime, "nvim*"))
		sockets = append(sockets, xdgSockets...)
		xdgSubSockets, _ := filepath.Glob(filepath.Join(xdgRuntime, "nvim*", "*"))
		sockets = append(sockets, xdgSubSockets...)
	}
	for _, sock := range sockets {
		if err := exec.Command("nvim", "--server", sock, "--remote-send", "<Cmd>colorscheme palette<CR>").Run(); err != nil {
			continue
		}
	}
}