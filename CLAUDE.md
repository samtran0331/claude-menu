# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

`claude-menu` is a tiny Go TUI that presents an arrow-key-navigable menu to launch Claude Code with different backend profiles (personal API, company Vertex, Ollama, LM Studio). It clears conflicting env vars, applies the selected profile's vars, then exec's the target command.

## Commands

```bash
# Build
go build -o claude-menu .

# Run
go run main.go

# Tidy dependencies
go mod tidy
```

## Architecture

Single file: `main.go`.

- `Target` struct — name, command, args, env overrides for one profile.
- `main()` — defines profiles, opens keyboard listener, renders menu, handles Up/Down/Enter/Esc.
- `render()` — clears screen, draws bordered list with cyan highlight on selected item.
- `runSelection()` — strips known conflicting env keys from the process, applies target's env, execs the command with inherited stdio.
- `clearScreen()` — ANSI escape on Unix, `cmd /c cls` on Windows.

## Dependency Note

The keyboard library import in `main.go` (`"://github.com"`) is malformed and `go.mod` is missing the dependency entry. The intended package is `github.com/eiannone/keyboard`. Run `go get github.com/eiannone/keyboard` and fix the import before building.
