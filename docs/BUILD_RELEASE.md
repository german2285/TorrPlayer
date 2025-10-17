# Build and Release Guide

This guide explains how to build TorrPlayer and create GitHub releases.

## Prerequisites

### Required Tools

1. **Go** (1.23 or newer)
2. **Node.js** (16 or newer)
3. **Wails CLI**
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```
4. **MinGW-w64** (for cross-compilation on Linux)
   ```bash
   sudo apt-get install mingw-w64  # Ubuntu/Debian
   ```
5. **GitHub CLI** (for automatic releases)
   ```bash
   # Installation instructions: https://cli.github.com/
   # Ubuntu/Debian
   curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
   echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
   sudo apt update
   sudo apt install gh
   ```

### GitHub Authentication

Before creating releases, authenticate with GitHub:

```bash
gh auth login
```

Follow the prompts to authenticate with your GitHub account.

## Build Scripts

### Production Build (`build.sh`)

Creates a production-ready build without debug features.

#### Basic Build (No Release)

Build without creating a release package:

```bash
cd /root/TorrPlayer
bash scripts/build.sh
```

Output: `build/bin/torrplayer-merged.exe`

#### Build with Version Package

Create a versioned ZIP package without GitHub release:

```bash
bash scripts/build.sh -v v1.0.0
```

Output:
- `build/bin/torrplayer-merged.exe`
- `build/release/TorrPlayer-v1.0.0-windows-amd64.zip`

Package contents:
- `torrplayer-merged.exe`
- `libmpv-2.dll`
- `README.md`
- `LICENSE`

#### Build and Create GitHub Release

Create a full GitHub release with automatic upload:

```bash
# With default release notes
bash scripts/build.sh -v v1.0.0 -r

# With custom release notes
bash scripts/build.sh -v v1.0.0 -r -n "Initial stable release with streaming support"
```

This will:
1. Build the application
2. Create a ZIP package
3. Create a GitHub release with tag `v1.0.0`
4. Upload the ZIP file to the release
5. Display the release URL

### Development Build (`build-dev.sh`)

Creates a debug build with developer tools enabled.

#### Basic Development Build

```bash
bash scripts/build-dev.sh
```

#### Development Build with Pre-Release

```bash
# Create a pre-release on GitHub
bash scripts/build-dev.sh -v v1.0.0-dev -r

# With custom notes
bash scripts/build-dev.sh -v v1.0.0-dev -r -n "Testing new cache algorithm"
```

Development builds:
- Include debug console window
- Enable DevTools (F12)
- Marked as "pre-release" on GitHub
- Include `DEV-BUILD-NOTICE.txt` file

## Version Naming Convention

### Production Releases
- Stable: `v1.0.0`, `v1.2.3`
- Beta: `v1.0.0-beta`, `v1.0.0-beta.1`
- Release Candidate: `v1.0.0-rc`, `v1.0.0-rc.2`

### Development Builds
- Development: `v1.0.0-dev`, `v1.1.0-dev`
- Feature branch: `v1.0.0-feature-name`
- Hotfix: `v1.0.1-hotfix`

## Complete Examples

### Example 1: First Stable Release

```bash
cd /root/TorrPlayer
bash scripts/build.sh -v v1.0.0 -r -n "Initial stable release

Features:
- Stream torrents directly with MPV player
- In-memory caching system
- Configurable cache size and connection limits
- Support for magnet links and .torrent files

Requirements:
- Windows 10/11 (64-bit)
- libmpv-2.dll (included)"
```

### Example 2: Beta Release

```bash
bash scripts/build.sh -v v1.1.0-beta -r -n "Beta release with new features

New features:
- Improved cache management
- Enhanced torrent metadata loading
- Better error handling

Known issues:
- Settings page UI needs polishing
- Some trackers may timeout"
```

### Example 3: Development Build

```bash
bash scripts/build-dev.sh -v v1.2.0-dev -r -n "Testing improved streaming algorithm

Changes:
- Experimental cache preloading
- Debug logging for performance metrics

This is a development build for testing only."
```

### Example 4: Hotfix Release

```bash
bash scripts/build.sh -v v1.0.1 -r -n "Hotfix: Fix MPV initialization error

Fixes:
- Fixed crash when libmpv-2.dll is missing
- Improved error messages
- Fixed memory leak in cache cleanup"
```

## Script Parameters

### build.sh / build-dev.sh

| Parameter | Short | Description | Required |
|-----------|-------|-------------|----------|
| `--version` | `-v` | Version tag (e.g., v1.0.0) | No* |
| `--release` | `-r` | Create GitHub release | No |
| `--notes` | `-n` | Release notes | No |

\* Required when using `-r` flag

## Output Structure

```
build/
├── bin/
│   ├── torrplayer-merged.exe
│   └── libmpv-2.dll
└── release/
    ├── TorrPlayer-v1.0.0-windows-amd64/
    │   ├── torrplayer-merged.exe
    │   ├── libmpv-2.dll
    │   ├── README.md
    │   └── LICENSE
    └── TorrPlayer-v1.0.0-windows-amd64.zip
```

## Troubleshooting

### GitHub CLI Not Found

```bash
# Install GitHub CLI
# See: https://cli.github.com/
```

### Not Authenticated

```bash
gh auth login
```

### Wails Not Found

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
export PATH=$PATH:/root/go/bin
```

### MinGW-w64 Not Found (Linux)

```bash
sudo apt-get install mingw-w64
```

### Release Already Exists

If a release with the same version tag already exists, you'll need to either:

1. Delete the existing release:
   ```bash
   gh release delete v1.0.0 --repo german2285/TorrPlayer
   ```

2. Use a different version number:
   ```bash
   bash scripts/build.sh -v v1.0.1 -r
   ```

## CI/CD Integration

You can integrate these scripts into GitHub Actions or other CI/CD systems.

Example GitHub Actions workflow:

```yaml
name: Build and Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'

      - name: Install dependencies
        run: |
          sudo apt-get update
          sudo apt-get install -y mingw-w64
          go install github.com/wailsapp/wails/v2/cmd/wails@latest

      - name: Build and Release
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          bash scripts/build.sh -v ${GITHUB_REF#refs/tags/} -r
```

## Best Practices

1. **Semantic Versioning**: Follow semver (MAJOR.MINOR.PATCH)
2. **Changelog**: Maintain a CHANGELOG.md file
3. **Testing**: Test development builds before creating production releases
4. **Backup**: Keep backups of release artifacts
5. **Documentation**: Update documentation when releasing new features
6. **Git Tags**: Create git tags for releases
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

## Manual Release (Without Scripts)

If you need to create a release manually:

```bash
# Build
cd /root/TorrPlayer
bash scripts/build.sh -v v1.0.0

# Create release manually
gh release create v1.0.0 \
  build/release/TorrPlayer-v1.0.0-windows-amd64.zip \
  --title "TorrPlayer v1.0.0" \
  --notes "Release notes here" \
  --repo german2285/TorrPlayer
```

## Links

- GitHub Repository: https://github.com/german2285/TorrPlayer
- GitHub Releases: https://github.com/german2285/TorrPlayer/releases
- GitHub CLI Documentation: https://cli.github.com/manual/
- Wails Documentation: https://wails.io/docs/introduction/
