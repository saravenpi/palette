# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Releases before 0.6.0 predate this changelog.

## [Unreleased]

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
