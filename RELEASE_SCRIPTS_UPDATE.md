# Build Scripts Update Summary

This document summarizes the updates made to the TorrPlayer build scripts to support automated versioning and GitHub releases.

## Overview

The build scripts (`build.sh` and `build-dev.sh`) have been enhanced with:
- ✅ Version parameter support
- ✅ Automatic GitHub release creation
- ✅ Custom release notes
- ✅ ZIP package creation
- ✅ Help documentation
- ✅ Comprehensive error handling

## What's New

### 1. Enhanced Build Scripts

#### Production Build (`scripts/build.sh`)
**New capabilities:**
- Accept version parameter: `-v v1.0.0`
- Create GitHub releases: `-r`
- Add custom release notes: `-n "Release notes"`
- Display help: `--help`
- Automatic ZIP packaging
- GitHub CLI integration

**Example usage:**
```bash
# Simple build
bash scripts/build.sh

# Build with version
bash scripts/build.sh -v v1.0.0

# Full release to GitHub
bash scripts/build.sh -v v1.0.0 -r -n "Initial release"
```

#### Development Build (`scripts/build-dev.sh`)
**New capabilities:**
- Same features as production build
- Creates GitHub **pre-releases** (marked as not stable)
- Includes DEV-BUILD-NOTICE.txt in package
- Enhanced default release notes for dev builds

**Example usage:**
```bash
# Development build
bash scripts/build-dev.sh

# Dev pre-release
bash scripts/build-dev.sh -v v1.0.0-dev -r
```

### 2. New Documentation

Created comprehensive documentation in `docs/` folder:

#### BUILD_RELEASE.md
Complete guide covering:
- Prerequisites and setup
- Build instructions
- GitHub release automation
- Version naming conventions
- Troubleshooting guide
- CI/CD integration examples

#### QUICK_REFERENCE.md
Quick command reference with:
- Common commands
- Parameter descriptions
- Version naming guide
- Quick troubleshooting

#### EXAMPLES.md
Real-world examples including:
- Local development workflows
- GitHub release examples
- Complete release workflows
- Special case scenarios
- CI/CD automation examples

#### CHANGELOG_BUILD_SCRIPTS.md
Tracks changes to build scripts:
- Version history
- Added features
- Breaking changes
- Migration guide

### 3. Updated Files

#### Updated `README.md`
Added section for production releases:
```markdown
### Production Release

Create a versioned release and publish to GitHub:

bash scripts/build.sh -v v1.0.0 -r -n "Release notes"
bash scripts/build-dev.sh -v v1.0.0-dev -r

**See [BUILD_RELEASE.md](docs/BUILD_RELEASE.md) for detailed instructions.**
```

#### Updated `docs/README.md`
Documentation index with links to all guides.

## Features Breakdown

### Version Management
```bash
# Specify version for build
-v v1.0.0        # Production release
-v v1.0.0-beta   # Beta release
-v v1.0.0-dev    # Development build
-v v1.0.0-rc     # Release candidate
```

### GitHub Release Automation
```bash
# Create GitHub release automatically
-r

# With custom release notes
-n "Release notes here"

# Complete example
bash scripts/build.sh -v v1.0.0 -r -n "Initial stable release"
```

### Output Structure
```
build/
├── bin/
│   ├── torrplayer-merged.exe  # Executable
│   └── libmpv-2.dll            # Required library
└── release/
    ├── TorrPlayer-v1.0.0-windows-amd64/
    │   ├── torrplayer-merged.exe
    │   ├── libmpv-2.dll
    │   ├── README.md
    │   └── LICENSE
    └── TorrPlayer-v1.0.0-windows-amd64.zip  # Release package
```

### Help System
```bash
# Display help for production build
bash scripts/build.sh --help

# Display help for development build
bash scripts/build-dev.sh --help
```

Output includes:
- Usage syntax
- Parameter descriptions
- Examples
- Link to documentation

## Prerequisites

### Required Tools

1. **Go** (1.23+)
2. **Node.js** (16+)
3. **Wails CLI**
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

4. **MinGW-w64** (Linux cross-compilation)
   ```bash
   sudo apt-get install mingw-w64
   ```

5. **GitHub CLI** (for releases)
   ```bash
   sudo apt install gh
   gh auth login
   ```

## Usage Examples

### Example 1: Local Development
```bash
cd /root/TorrPlayer
bash scripts/build-dev.sh
```

### Example 2: Create Version Package
```bash
bash scripts/build.sh -v v1.0.0
# Output: build/release/TorrPlayer-v1.0.0-windows-amd64.zip
```

### Example 3: Full GitHub Release
```bash
bash scripts/build.sh -v v1.0.0 -r -n "Initial stable release

Features:
- Stream torrents with MPV player
- Configurable cache settings
- Support for magnet links

Requirements:
- Windows 10/11 (64-bit)"
```

### Example 4: Development Pre-Release
```bash
bash scripts/build-dev.sh -v v1.1.0-dev -r -n "Testing new cache algorithm"
```

## Script Parameters

| Parameter | Short | Description | Required |
|-----------|-------|-------------|----------|
| `--version` | `-v` | Version tag (e.g., v1.0.0) | No* |
| `--release` | `-r` | Create GitHub release | No |
| `--notes` | `-n` | Release notes | No |
| `--help` | `-h` | Show help message | No |

\* Required when using `-r` flag

## Key Improvements

### 1. Automation
- **Before**: Manual ZIP creation, manual GitHub release
- **After**: One command creates everything

### 2. Consistency
- Standardized version naming
- Consistent package structure
- Predictable output paths

### 3. Safety
- Validates GitHub CLI availability
- Checks authentication
- Verifies required tools
- Clear error messages

### 4. Documentation
- Comprehensive guides
- Real-world examples
- Quick reference
- Inline help

### 5. Flexibility
- Works without parameters (backward compatible)
- Optional version packaging
- Optional GitHub release
- Custom release notes

## Workflow Examples

### Release Workflow
```bash
# 1. Development
bash scripts/build-dev.sh

# 2. Testing preview
bash scripts/build-dev.sh -v v1.0.0-dev -r

# 3. Release candidate
bash scripts/build-dev.sh -v v1.0.0-rc -r

# 4. Final release
bash scripts/build.sh -v v1.0.0 -r -n "Production release"
```

### Hotfix Workflow
```bash
# 1. Fix bug in code
# 2. Create hotfix release
bash scripts/build.sh -v v1.0.1 -r -n "Hotfix: Critical bug fix"
```

## Testing

The scripts have been tested with:
- ✅ No parameters (backward compatibility)
- ✅ Version parameter only
- ✅ Version + release flag
- ✅ Full command with release notes
- ✅ Help flag
- ✅ Invalid parameters (error handling)

## File Changes Summary

### Modified Files
- `scripts/build.sh` - Enhanced with new features
- `scripts/build-dev.sh` - Enhanced with new features
- `README.md` - Added production release section
- `docs/README.md` - Updated with new docs

### New Files
- `docs/BUILD_RELEASE.md` - Complete build guide (7.4 KB)
- `docs/QUICK_REFERENCE.md` - Quick reference (2.7 KB)
- `docs/EXAMPLES.md` - Usage examples (9.4 KB)
- `docs/CHANGELOG_BUILD_SCRIPTS.md` - Change tracking (4.8 KB)
- `RELEASE_SCRIPTS_UPDATE.md` - This file

### Total Documentation Added
- ~24 KB of comprehensive documentation
- 4 new documentation files
- Updated 3 existing files

## GitHub Repository Integration

The scripts integrate with: https://github.com/german2285/TorrPlayer

**Automatic actions:**
1. Creates release with specified tag
2. Uploads ZIP package as release asset
3. Sets release title: "TorrPlayer {VERSION}"
4. Adds custom or default release notes
5. Marks dev builds as pre-release

**Release URL format:**
```
https://github.com/german2285/TorrPlayer/releases/tag/{VERSION}
```

## Migration Guide

### From Old Scripts
```bash
# Old way
bash build.sh
# Manual packaging
# Manual GitHub release

# New way
bash scripts/build.sh -v v1.0.0 -r -n "Release notes"
```

### Backward Compatibility
Old usage still works:
```bash
bash scripts/build.sh
# Produces: build/bin/torrplayer-merged.exe
```

## Troubleshooting

### Issue: gh command not found
```bash
sudo apt install gh
gh auth login
```

### Issue: Release already exists
```bash
gh release delete v1.0.0 --repo german2285/TorrPlayer
# Then re-run build script
```

### Issue: wails not found
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
export PATH=$PATH:/root/go/bin
```

## Next Steps

1. **Test the scripts:**
   ```bash
   bash scripts/build.sh --help
   bash scripts/build-dev.sh --help
   ```

2. **Read documentation:**
   - Start with `docs/QUICK_REFERENCE.md`
   - Review `docs/EXAMPLES.md` for workflows
   - Check `docs/BUILD_RELEASE.md` for details

3. **Try a release:**
   ```bash
   # Without GitHub release
   bash scripts/build.sh -v v1.0.0

   # With GitHub release (when ready)
   bash scripts/build.sh -v v1.0.0 -r
   ```

## Links

- **Documentation:** `docs/README.md`
- **Build Guide:** `docs/BUILD_RELEASE.md`
- **Quick Reference:** `docs/QUICK_REFERENCE.md`
- **Examples:** `docs/EXAMPLES.md`
- **Changelog:** `docs/CHANGELOG_BUILD_SCRIPTS.md`
- **Repository:** https://github.com/german2285/TorrPlayer
- **Releases:** https://github.com/german2285/TorrPlayer/releases

## Support

For questions or issues:
1. Check documentation in `docs/`
2. Run `bash scripts/build.sh --help`
3. Open an issue on GitHub

## Conclusion

The build scripts now provide a complete automated workflow for:
- ✅ Building TorrPlayer
- ✅ Creating release packages
- ✅ Publishing to GitHub
- ✅ Version management
- ✅ Development and production builds

All with comprehensive documentation and examples!
