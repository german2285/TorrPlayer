# Quick Reference - Build Scripts

## Production Build (`build.sh`)

### Simple Build
```bash
bash scripts/build.sh
```
Output: `build/bin/torrplayer-merged.exe`

### Build with Version Package
```bash
bash scripts/build.sh -v v1.0.0
```
Output: `build/release/TorrPlayer-v1.0.0-windows-amd64.zip`

### Build and Release to GitHub
```bash
bash scripts/build.sh -v v1.0.0 -r
```

### Full Release with Notes
```bash
bash scripts/build.sh -v v1.0.0 -r -n "Release notes here"
```

## Development Build (`build-dev.sh`)

### Simple Dev Build
```bash
bash scripts/build-dev.sh
```
Output: `build/bin/torrplayer-merged.exe` (with debug features)

### Dev Build with Pre-Release
```bash
bash scripts/build-dev.sh -v v1.0.0-dev -r
```

### Full Dev Release with Notes
```bash
bash scripts/build-dev.sh -v v1.0.0-dev -r -n "Development build notes"
```

## Parameters

| Flag | Long Form | Description | Example |
|------|-----------|-------------|---------|
| `-v` | `--version` | Version tag | `-v v1.0.0` |
| `-r` | `--release` | Create GitHub release | `-r` |
| `-n` | `--notes` | Release notes | `-n "Bug fixes"` |

## Common Workflows

### 1. Local Development
```bash
bash scripts/build-dev.sh
```

### 2. Testing Before Release
```bash
bash scripts/build-dev.sh -v v1.0.0-rc -r -n "Release candidate for testing"
```

### 3. Stable Release
```bash
bash scripts/build.sh -v v1.0.0 -r -n "Stable release"
```

### 4. Hotfix
```bash
bash scripts/build.sh -v v1.0.1 -r -n "Hotfix: Critical bug"
```

## Version Naming

- Stable: `v1.0.0`
- Beta: `v1.0.0-beta`
- Dev: `v1.0.0-dev`
- RC: `v1.0.0-rc`

## Prerequisites

```bash
# Install GitHub CLI
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
sudo apt update
sudo apt install gh

# Authenticate
gh auth login

# Install Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Install MinGW (Linux only)
sudo apt-get install mingw-w64
```

## Troubleshooting

### Error: gh not found
```bash
sudo apt install gh
gh auth login
```

### Error: wails not found
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
export PATH=$PATH:/root/go/bin
```

### Error: Release already exists
```bash
gh release delete v1.0.0 --repo german2285/TorrPlayer
```

## Links

- [Full Build Guide](BUILD_RELEASE.md)
- [Repository](https://github.com/german2285/TorrPlayer)
- [Releases](https://github.com/german2285/TorrPlayer/releases)
