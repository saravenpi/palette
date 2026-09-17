# palette

Terminal color palette generator and dotfile synchronization system.

`palette` extracts harmonious 16-color ANSI palettes from images and synchronizes them across Kitty, Tmux, and Neovim with zero startup latency and instant hot-reloading.

## Installation

```bash
git clone https://github.com/saravenpi/palette.git
cd palette
./install.sh
```

This compiles and installs the binary to `~/.local/bin/palette`.

## Usage

### Extract from an image and sync instantly

```bash
palette wallpaper.jpg -s
```

This extracts dominant colors from the image, writes them to `~/.palette.yml`, and immediately updates Kitty, Tmux, and Neovim.

### Sync manually after editing colors

Edit `~/.palette.yml`, then run:

```bash
palette sync
```

### Preview colors without applying

```bash
palette wallpaper.jpg -p
```

### Export to specific formats

```bash
palette wallpaper.jpg -f kitty
palette wallpaper.jpg -f tmux
palette wallpaper.jpg -f nvim
palette wallpaper.jpg -f ghostty
palette wallpaper.jpg -f alacritty
palette wallpaper.jpg -f wezterm
palette wallpaper.jpg -f foot
palette wallpaper.jpg -f yaml
palette wallpaper.jpg -f json
```

## Dotfiles setup

### Kitty

Add to `~/.config/kitty/kitty.conf`:

```conf
include palette.conf
```

### Tmux

Add to the end of `~/.tmux.conf`:

```tmux
source-file -q ~/.config/tmux/palette.conf
```

### Neovim

Add to your Neovim configuration:

```lua
vim.cmd.colorscheme("palette")
```

Or consume the raw color table in your own theme:

```lua
local colors = require("palette").colors
```

## How sync works

1. Canonical colors are stored in `~/.palette.yml`.
2. `palette sync` compiles native configuration snippets to:
   - Kitty: `~/.config/kitty/palette.conf`
   - Tmux: `~/.config/tmux/palette.conf`
   - Neovim: `~/.config/nvim/lua/palette.lua` and `~/.config/nvim/colors/palette.lua`
3. Running applications reload live without restarting:
   - Kitty receives `SIGUSR1`.
   - Tmux runs `tmux source-file`.
   - Active Neovim sockets receive `<Cmd>colorscheme palette<CR>`.
