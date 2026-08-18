# claude-menu installer for Windows (PowerShell).
# Usage:
#   irm https://raw.githubusercontent.com/samtran0331/claude-menu/main/scripts/install.ps1 | iex
# Override the source repo or version:
#   $env:CLAUDE_MENU_REPO="samtran0331/claude-menu"; $env:VERSION="v1.0.0"; ./install.ps1
$ErrorActionPreference = "Stop"

$repo = if ($env:CLAUDE_MENU_REPO) { $env:CLAUDE_MENU_REPO } else { "samtran0331/claude-menu" }
$version = if ($env:VERSION) { $env:VERSION } else { "latest" }
$binary = "claude-menu"

$arch = switch ($env:PROCESSOR_ARCHITECTURE) {
  "AMD64" { "amd64" }
  "ARM64" { "arm64" }
  default { throw "Unsupported arch: $env:PROCESSOR_ARCHITECTURE" }
}

$asset = "${binary}_windows_${arch}.exe"
$dir = Join-Path $env:LOCALAPPDATA "Programs\claude-menu"
New-Item -ItemType Directory -Force -Path $dir | Out-Null
$dest = Join-Path $dir "$binary.exe"

if ($version -eq "latest") {
  $url = "https://github.com/$repo/releases/latest/download/$asset"
} else {
  $url = "https://github.com/$repo/releases/download/$version/$asset"
}

Write-Host "Downloading $url"
try {
  Invoke-WebRequest -Uri $url -OutFile $dest -UseBasicParsing
  Write-Host "Installed $dest"
} catch {
  if ((Get-Command go -ErrorAction SilentlyContinue) -and (Test-Path "go.mod")) {
    Write-Host "Download failed — building from source..."
    go build -ldflags "-s -w" -o $dest .
    Write-Host "Installed $dest"
  } else {
    throw "Download failed and no Go toolchain / source available. Set CLAUDE_MENU_REPO / VERSION, or run from the source tree with Go installed."
  }
}

# Add install dir to the user PATH if missing.
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$dir*") {
  [Environment]::SetEnvironmentVariable("Path", "$userPath;$dir", "User")
  Write-Host "Added $dir to your user PATH. Restart your terminal to pick it up."
}

Write-Host "Run: $binary"
