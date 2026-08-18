# claude-menu

TUI agent selector for claude

## Install (users)

Prebuilt binaries for macOS and Windows are published on the GitHub Releases page.

> Note: the repo is currently **private**, so the `curl … | sh` / `irm … | iex`
> one-liners below won't work until it's made public (they fetch raw script files
> anonymously). Until then, clone the repo and run the scripts locally.

### macOS

```bash
curl -fsSL https://raw.githubusercontent.com/samtran0331/claude-menu/main/scripts/install.sh | sh
```

Installs to `/usr/local/bin` (or `~/.local/bin` if that isn't writable). Then run:

```bash
claude-menu
```

macOS Gatekeeper may block an unsigned binary on first run. If so:
`xattr -d com.apple.quarantine "$(command -v claude-menu)"` then run again.

### Windows

In PowerShell:

```powershell
irm https://raw.githubusercontent.com/samtran0331/claude-menu/main/scripts/install.ps1 | iex
```

Installs to `%LOCALAPPDATA%\Programs\claude-menu` and adds it to your user PATH.
Restart the terminal, then run:

```powershell
claude-menu
```

### Manual download

Grab the binary for your OS/arch from the latest release, make it executable
(`chmod +x` on macOS), put it on your PATH, and run `claude-menu`.

Asset names: `claude-menu_darwin_arm64`, `claude-menu_darwin_amd64`,
`claude-menu_windows_amd64.exe`, `claude-menu_windows_arm64.exe`.

## Configure the Kimi key (users)

The **My Kimi Subscription** profile reads its token from the `KIMI_API_KEY`
environment variable — no secret is stored in the app. Set it once:

### macOS / Linux

Add this line to `~/.zshrc` (or `~/.bashrc`), then reload:

```bash
echo 'export KIMI_API_KEY="sk-your-token-here"' >> ~/.zshrc
source ~/.zshrc
```

### Windows

Set a persistent user environment variable (PowerShell), then open a new terminal:

```powershell
setx KIMI_API_KEY "sk-your-token-here"
```

Or via the GUI: **Settings > System > About > Advanced system settings >
Environment Variables > New** (under *User variables*).

If the variable is missing, selecting the Kimi profile shows a reminder with these
steps instead of launching.

## Publish (developers)

> ⚠️ **Never hardcode tokens.** The Kimi token is read from `KIMI_API_KEY` at
> runtime, so binaries contain no secret. Keep it that way — don't commit tokens
> or bake them into builds.

Requires the [`gh`](https://cli.github.com) CLI, authenticated (`gh auth login`).

```bash
# 1. Cross-compile all release binaries + checksums into dist/
make release

# 2. Tag and cut a GitHub Release with the binaries attached
git tag v1.0.0
git push origin v1.0.0
gh release create v1.0.0 ./dist/* --title "v1.0.0" --generate-notes
```

The install scripts download from the `latest` release by default, so users get
the new build automatically once the release is published.

`make release` targets: `darwin/amd64`, `darwin/arm64`, `windows/amd64`,
`windows/arm64`. Add/remove platforms via the `PLATFORMS` variable in the `Makefile`.

## Run

```bash
go run main.go
```

## Build & run binary

```bash
go build -o claude-menu .
./claude-menu

# Or build + install to /usr/local/bin (override dir with PREFIX=~/.local)
make install
```

## Debug

```bash
# Verbose output — add GODEBUG prefix
GODEBUG=asyncpreemptoff=1 go run main.go

# Step through with Delve
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug main.go
```

## Add / edit profiles

Edit the `options` slice in `main()`. Each entry is a `Target`:

```go
{
    Name: "My Profile",
    Cmd:  "claude",
    Args: []string{},                          // optional extra CLI args
    Env:  map[string]string{"KEY": "value"},   // env vars applied before launch
}
```

Conflicting env vars (`ANTHROPIC_MODEL`, `CLAUDE_CODE_USE_VERTEX`, etc.) are automatically cleared before the selected profile's vars are applied.
