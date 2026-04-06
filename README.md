# gitree

A minimal terminal UI for git — like lazygit, but only the parts you actually need.

![Go](https://img.shields.io/badge/Go-1.26-blue)

## Features

- **Changes panel** — file tree view of uncommitted changes, color-coded by status
- **Branches panel** — list all branches, switch with enter, create new ones with `B`
- **Commits panel** — recent commits for the selected branch, updates as you navigate
- **Diff preview** — full-screen diff view for any changed file

## Install

```sh
go install github.com/morizkraemer/gitree@latest
```

Or build from source:

```sh
git clone https://github.com/morizkraemer/gitree.git
cd gitree
go build -o gitree .
```

## Usage

Run `gitree` inside any git repository.

### Keybindings

| Key | Action |
|-----|--------|
| `tab` / `shift+tab` | Switch between panels |
| `j` / `k` | Navigate up/down |
| `enter` | Switch branch / open diff preview |
| `B` | Create new branch (in branches panel) |
| `r` | Refresh all data |
| `d` / `u` | Page down/up (in diff view) |
| `q` / `esc` | Quit / close diff |

## Built with

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — styling
