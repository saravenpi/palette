# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Releases before 0.6.0 predate this changelog.

## [0.12.0] - 2026-10-19

### Added

- `yazi`: new sync target and export format writing a complete
  `~/.config/yazi/theme.toml`. Yazi merges it over its built-in dark/light
  preset, so layout, icons and keymaps stay in `yazi.toml` and `keymap.toml`
  while the file list, markers, tabs, modes, indicator, status bar, dialogs and
  per-file-type rules follow the palette. Yazi reads the theme at startup and
  has no reload hook, so a running instance needs a restart.

## [0.11.0] - 2026-10-13

### Added

- `hyprland`: new sync target and export format writing a colour fragment to
  `~/.config/hypr/palette.conf`. Hyprland spells colours `rgb(RRGGBB)`, so every
  value is converted from the palette's hex, and three role variables
  (`$active_border`, `$inactive_border`, `$urgent_border`) name the focused,
  dimmed and urgent border colours so window rules never hard-code an index.
- `palette sync` reloads a running Hyprland with `hyprctl reload` after writing.

## [0.10.0] - 2026-10-13

### Added

- `mako`: new sync target and export format writing the palette-driven mako
  colours (`background-color`, `text-color`, `border-color`, `progress-color`
  and an `[urgency=critical]` accent) to `~/.config/mako/palette.conf`, the
  fragment a mako config pulls in with `include=`.
- `palette sync` reloads a running mako with `makoctl reload` after writing.

## [0.9.2] - 2026-10-13

### Changed

- README: the widgets section no longer tells you to restart a widget after
  `palette sync`. The widgets watch `~/.config/widgets/palette.css`, so the bar,
  sidebar and launcher restyle themselves when the file is rewritten.

## [0.9.1] - 2026-10-13

### Fixed

- `sync`: a config directory that is a symlink whose target is missing — the
  layout dotfile managers routinely leave behind — made `os.MkdirAll` fail with
  "file exists" and aborted the whole sync. The link is now resolved and its
  target directory created only when absent, so every target is written.

## [0.9.0] - 2026-10-13

### Changed

- `widgets` and `waybar`: the generated stylesheet now declares every color as
  a `:root` custom property (`--palette-bg`, `--palette-fg`, `--palette-cursor`,
  `--palette-color0`..`--palette-color15`) so stylesheets can use `var()`, and
  keeps the `@define-color` declarations below it as the fallback. Custom
  properties, `:root` and `var()` are GTK 4.16 and later; GTK 4.10 to 4.15
  ignore the new block and read the named colors exactly as before, so existing
  consumers need no change.

## [0.8.0] - 2026-10-13

### Added

- `widgets`: new sync target and export format writing GTK color definitions
  (`@define-color palette_bg`, `palette_fg`, `palette_cursor` and
  `palette_color0`..`palette_color15`) to `~/.config/widgets/palette.css`, the
  file the statusbar, sidebar and launcher of
  [widgets](https://github.com/Saravenpi/widgets) load at startup.
- `vim`: new sync target and export format writing a complete Vim colorscheme
  (`g:colors_name = "palette"`, GUI hex colors plus cterm 0-15 indices) to
  `~/.vim/colors/palette.vim`.

## [0.7.0] - 2026-10-12

### Added

- `waybar`: new sync target and export format writing GTK color definitions
  (`@define-color palette_bg`, `palette_fg`, `palette_cursor` and
  `palette_color0`..`palette_color15`) to `~/.config/waybar/palette.css`, so
  waybar carries the same colors as foot, sway and tmux.
- `palette sync` reloads a running waybar with `SIGUSR2` after writing.

## [0.6.0] - 2026-10-11

### Changed

- `tmux`: the synced fragment now contains only color variables
  (`@palette_color0..15`, `@palette_bg`, `@palette_fg`, `@palette_cursor` and
  the `@paper_*` aliases). It no longer sets `status-style`,
  `window-status-style`, `window-status-format`, `status-left`, the pane border
  styles, `mode-style`, `message-style` or `copy-mode-match-style`.
- `foot`: the snippet is a `[colors-dark]` section with `cursor=<text> <cursor>`
  inside it instead of a `[cursor]` section.

### Fixed

- `foot`: `[colors]` was removed in foot 1.28.0 and `[cursor] color=` in 1.26.0,
  so foot rejected the generated file with "unknown option".
