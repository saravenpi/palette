# palette

Terminal color palette generator and dotfile synchronization system.

`palette` extracts harmonious 16-color ANSI palettes from images and synchronizes them across Kitty, Tmux, Neovim, Sway, Foot, Waybar, Widgets, and Vim with zero startup latency and instant hot-reloading.

## Installation

Install directly with curl:

```bash
curl -fsSL https://raw.githubusercontent.com/saravenpi/palette/main/install.sh | bash
```

Or build from source:

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

### Palette modes

`palette` generates themes using Material Design 3 (Material You) with Celebi quantization and CAM16 scoring:

```bash
palette wallpaper.jpg -m tonal      # Default Material You tonal spot scheme
palette wallpaper.jpg -m wal        # Pure wallpaper cluster colors (pywal style)
palette wallpaper.jpg -m vibrant    # High-chroma dynamic scheme
palette wallpaper.jpg -m expressive # Complementary expressive accents
palette wallpaper.jpg -m material   # Full Material 3 dynamic color roles
palette wallpaper.jpg -m duo        # Two-accent high-contrast palette
palette wallpaper.jpg -m dominant   # Accents ordered by tonal lightness
```

Use `-l` or `--light` with any mode to generate a light theme.

### Export to specific formats

```bash
palette wallpaper.jpg -f kitty
palette wallpaper.jpg -f tmux
palette wallpaper.jpg -f nvim
palette wallpaper.jpg -f ghostty
palette wallpaper.jpg -f alacritty
palette wallpaper.jpg -f wezterm
palette wallpaper.jpg -f foot
palette wallpaper.jpg -f waybar
palette wallpaper.jpg -f widgets
palette wallpaper.jpg -f vim
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

Source the snippet from `~/.config/tmux/tmux.conf`:

```tmux
source-file -q ~/.config/tmux/palette.conf
```

The tmux snippet only exposes colors: `@palette_color0` through `@palette_color15`,
`@palette_bg`, `@palette_fg`, `@palette_cursor` and the `@paper_*` aliases. It sets no
style of its own, so your own status bar keeps its look. Reference a value as
`#{@palette_color8}` from any style or format.

### Foot

Add to `~/.config/foot/foot.ini`:

```ini
[main]
include=~/.config/foot/palette.ini
```

The snippet is a `[colors-dark]` section, the format foot 1.26 and later expect.
`[colors]` and `[cursor] color=` were removed upstream in 1.26.0 and 1.28.0.

### Waybar

`palette sync` writes the colors to `~/.config/waybar/palette.css` in two
forms. Import it at the top of `~/.config/waybar/style.css` and use whichever
your GTK understands:

```css
@import url("palette.css");

window#waybar {
    background-color: var(--palette-bg);
    color: var(--palette-fg);
}
```

`var()` and the `:root` custom properties behind it are GTK 4.16 and later
(and standard CSS everywhere else). On GTK 4.10 to 4.15 use the named colors
instead — `@palette_bg`, `@palette_fg`, `@palette_cursor` and `@palette_color0`
through `@palette_color15`, which the same file declares as a fallback. Both
spellings carry the same values foot, sway and tmux receive, so the bar always
matches the rest of the desktop.

### Widgets

`palette sync` writes the same stylesheet to `~/.config/widgets/palette.css`,
which the `statusbar`, `sidebar` and `launcher` binaries of
[widgets](https://github.com/Saravenpi/widgets) load at startup. It defines
`--palette-bg`, `--palette-fg`, `--palette-cursor` and `--palette-color0`
through `--palette-color15` as `:root` custom properties, plus the
`@palette_bg`-style named colors GTK 4.15 and older fall back to. Nothing to
wire up: run `palette sync` and restart a widget to pick up new colors.

### Neovim

Add to your Neovim configuration:

```lua
vim.cmd.colorscheme("palette")
```

Or consume the raw color table in your own theme:

```lua
local colors = require("palette").colors
```

### Vim

`palette sync` writes a complete colorscheme to `~/.vim/colors/palette.vim`.
Load it from `~/.vimrc` (the default runtimepath already covers `~/.vim`):

```vim
colorscheme palette
```

## How sync works

1. Canonical colors are stored in `~/.palette.yml`.
2. `palette sync` compiles native configuration snippets to:
   - Kitty: `~/.config/kitty/palette.conf`
   - Tmux: `~/.config/tmux/palette.conf`
   - Neovim: `~/.config/nvim/lua/palette.lua` and `~/.config/nvim/colors/palette.lua`
   - Sway: `~/.config/sway/palette.conf`
   - Foot: `~/.config/foot/palette.ini`
   - Waybar: `~/.config/waybar/palette.css`
   - Widgets: `~/.config/widgets/palette.css`
   - Vim: `~/.vim/colors/palette.vim`
3. Running applications reload live without restarting:
   - Kitty receives `SIGUSR1`.
   - Tmux runs `tmux source-file`.
   - Active Neovim sockets receive `<Cmd>colorscheme palette<CR>`.
   - Sway runs `swaymsg reload`.
   - Waybar receives `SIGUSR2`.
