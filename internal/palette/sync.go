package palette

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func LoadPaletteFile(path string) (Theme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Theme{}, fmt.Errorf("palette file not found at %s: %w", path, err)
	}
	return ParseYAML(data)
}

func SavePaletteFile(path string, t Theme) error {
	data, err := formatYAML(t)
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(data), 0644)
}

func Sync(t Theme, reload bool) ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	var results []string

	kittyConf := filepath.Join(home, ".config", "kitty", "palette.conf")
	if err := writeFormatted(kittyConf, formatKitty(t)); err != nil {
		return nil, err
	}
	results = append(results, fmt.Sprintf("kitty -> %s", kittyConf))

	tmuxConf := filepath.Join(home, ".config", "tmux", "palette.conf")
	if err := writeFormatted(tmuxConf, formatTmux(t)); err != nil {
		return nil, err
	}
	results = append(results, fmt.Sprintf("tmux  -> %s", tmuxConf))

	nvimLua := filepath.Join(home, ".config", "nvim", "lua", "palette.lua")
	if err := writeFormatted(nvimLua, formatNvim(t)); err != nil {
		return nil, err
	}
	results = append(results, fmt.Sprintf("nvim  -> %s", nvimLua))

	nvimColors := filepath.Join(home, ".config", "nvim", "colors", "palette.lua")
	_ = writeFormatted(nvimColors, "package.loaded[\"palette\"] = nil\nrequire(\"palette\").load()\n")

	if reload {
		reloaded := reloadApplications(tmuxConf)
		results = append(results, reloaded...)
	}

	return results, nil
}

func writeFormatted(path string, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func reloadApplications(tmuxConf string) []string {
	var msgs []string

	uid := strconv.Itoa(os.Getuid())
	cmdKitty := exec.Command("pkill", "-SIGUSR1", "-u", uid, "kitty")
	if err := cmdKitty.Run(); err == nil {
		msgs = append(msgs, "reloaded kitty via SIGUSR1")
	}

	cmdTmux := exec.Command("tmux", "source-file", tmuxConf)
	if err := cmdTmux.Run(); err == nil {
		msgs = append(msgs, "reloaded tmux session")
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
		cmd := exec.Command("nvim", "--server", sock, "--remote-send", "<Cmd>colorscheme palette<CR>")
		_ = cmd.Run()
	}
}
