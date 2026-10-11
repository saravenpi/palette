# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Releases before 0.6.0 predate this changelog.

## [Unreleased]

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
