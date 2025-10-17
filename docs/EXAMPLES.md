# Build Script Examples

This document provides real-world examples of using the build scripts.

## Table of Contents

- [Local Development](#local-development)
- [Version Packages](#version-packages)
- [GitHub Releases](#github-releases)
- [Release Workflows](#release-workflows)
- [Special Cases](#special-cases)

## Local Development

### Quick Build for Testing

Just build the application without any versioning:

```bash
cd /root/TorrPlayer
bash scripts/build.sh
```

**Output:**
- `build/bin/torrplayer-merged.exe`
- `build/bin/libmpv-2.dll`

**Use case:** Quick local testing, no release artifacts needed.

### Development Build with Debug Tools

Build with console and DevTools for debugging:

```bash
bash scripts/build-dev.sh
```

**Output:**
- `build/bin/torrplayer-merged.exe` (debug version)
- Debug console window
- DevTools accessible with F12

**Use case:** Active development, debugging, testing new features.

## Version Packages

### Create Version Package (No GitHub Release)

Build and package without uploading to GitHub:

```bash
# Production build
bash scripts/build.sh -v v1.0.0

# Development build
bash scripts/build-dev.sh -v v1.0.0-dev
```

**Output:**
- `build/release/TorrPlayer-v1.0.0-windows-amd64.zip`
- Contains: exe, dll, README.md, LICENSE

**Use case:** Creating release artifacts for manual distribution or testing.

## GitHub Releases

### Prerequisites

```bash
# Install GitHub CLI (if not installed)
sudo apt install gh

# Authenticate
gh auth login
```

### Simple GitHub Release

Create a basic release with default notes:

```bash
bash scripts/build.sh -v v1.0.0 -r
```

**What happens:**
1. Builds the application
2. Creates ZIP package
3. Creates GitHub release with tag `v1.0.0`
4. Uploads ZIP file
5. Sets default release notes: "Release v1.0.0"

**GitHub URL:** https://github.com/german2285/TorrPlayer/releases/tag/v1.0.0

### Release with Custom Notes

Create a release with detailed release notes:

```bash
bash scripts/build.sh -v v1.0.0 -r -n "Initial stable release

Features:
- Stream torrents with MPV player
- Configurable cache settings
- Support for magnet links and .torrent files

Bug fixes:
- Fixed MPV initialization error
- Improved error handling"
```

**Use case:** Official releases with comprehensive changelog.

### Development Pre-Release

Create a pre-release for testing:

```bash
bash scripts/build-dev.sh -v v1.1.0-dev -r -n "Testing new cache algorithm

Changes:
- Experimental preload cache feature
- Enhanced torrent metadata loading
- Debug logging for performance

⚠️ This is a development build for testing only."
```

**What happens:**
1. Builds debug version
2. Adds DEV-BUILD-NOTICE.txt to package
3. Creates GitHub **pre-release** (marked as not production-ready)
4. Includes development warning in notes

## Release Workflows

### Workflow 1: Major Version Release

Complete workflow for releasing version 2.0.0:

```bash
# 1. Create release candidate for testing
bash scripts/build-dev.sh -v v2.0.0-rc -r -n "Release candidate for v2.0.0

New features:
- Feature A
- Feature B

Testing needed:
- Performance testing
- UI/UX feedback"

# 2. Test the RC build
# (Manual testing by team/community)

# 3. Create final stable release
bash scripts/build.sh -v v2.0.0 -r -n "TorrPlayer 2.0.0

Major new features:
- Feature A: Description
- Feature B: Description

Improvements:
- Improvement 1
- Improvement 2

Bug fixes:
- Fix 1
- Fix 2

Breaking changes:
- None"
```

### Workflow 2: Hotfix Release

Quick fix for critical bug in production:

```bash
# 1. Fix the bug in code
# (Edit source files)

# 2. Create hotfix release
bash scripts/build.sh -v v1.0.1 -r -n "Hotfix: Critical crash on startup

This hotfix addresses a critical bug that caused crashes when starting the application without libmpv-2.dll present.

Fixes:
- Fixed crash when libmpv-2.dll is missing
- Added proper error message
- Improved DLL loading logic

Users experiencing crashes should update immediately."
```

### Workflow 3: Beta Testing

Release beta version for community testing:

```bash
bash scripts/build.sh -v v1.1.0-beta -r -n "Beta: New features for testing

This is a BETA release. Please help us test these new features:

New features:
- Automatic torrent metadata caching
- Improved peer connection algorithm
- New settings page UI

Known issues:
- Settings page may have visual glitches
- Some trackers may timeout

Feedback:
Please report bugs at: https://github.com/german2285/TorrPlayer/issues

Thank you for testing!"
```

### Workflow 4: Feature Branch Preview

Preview build for specific feature:

```bash
bash scripts/build-dev.sh -v v1.1.0-feature-darkmode -r -n "Preview: Dark mode feature

This is a PREVIEW build showcasing the new dark mode feature.

Status: Work in progress
- Dark mode toggle in settings ✓
- Main UI dark theme ✓
- Player controls theme - in progress

Try it out and provide feedback!"
```

## Special Cases

### Multiple Builds Same Day

When creating multiple builds in one day:

```bash
# Morning build
bash scripts/build-dev.sh -v v1.0.0-dev.1 -r -n "Morning test build"

# Afternoon build with fixes
bash scripts/build-dev.sh -v v1.0.0-dev.2 -r -n "Afternoon build - fixed cache issue"

# Evening build ready for RC
bash scripts/build-dev.sh -v v1.0.0-rc -r -n "Release candidate"
```

### Release with Build Metadata

Include commit information in release notes:

```bash
COMMIT_HASH=$(git rev-parse --short HEAD)
COMMIT_MSG=$(git log -1 --pretty=%B)

bash scripts/build.sh -v v1.0.0 -r -n "Release v1.0.0

Build information:
- Commit: ${COMMIT_HASH}
- Message: ${COMMIT_MSG}
- Date: $(date)

Changes:
- Feature 1
- Feature 2"
```

### Patch Release Series

Creating patch releases (1.0.1, 1.0.2, 1.0.3):

```bash
# First patch
bash scripts/build.sh -v v1.0.1 -r -n "Patch: Fix torrent metadata loading"

# Second patch
bash scripts/build.sh -v v1.0.2 -r -n "Patch: Fix cache size calculation"

# Third patch
bash scripts/build.sh -v v1.0.3 -r -n "Patch: Fix settings persistence"
```

### Automated Release from CI/CD

Example for use in CI/CD pipeline:

```bash
#!/bin/bash
# ci-release.sh

# Get version from git tag
VERSION=${GITHUB_REF#refs/tags/}

# Generate release notes from commits
NOTES=$(git log --oneline --pretty=format:"- %s" $(git describe --tags --abbrev=0 @^)..@)

# Create release
bash scripts/build.sh -v "$VERSION" -r -n "Release $VERSION

Changes since last release:
$NOTES"
```

## Tips and Best Practices

### 1. Always Test Before Release

```bash
# Test with dev build first
bash scripts/build-dev.sh -v v1.0.0-rc -r

# After testing, create production release
bash scripts/build.sh -v v1.0.0 -r
```

### 2. Use Semantic Versioning

- **Major** (v2.0.0): Breaking changes
- **Minor** (v1.1.0): New features, backwards compatible
- **Patch** (v1.0.1): Bug fixes only

### 3. Write Meaningful Release Notes

Good example:
```bash
bash scripts/build.sh -v v1.1.0 -r -n "New features and improvements

Features:
- Added torrent search functionality
- Configurable download rate limiting

Improvements:
- Faster metadata loading (2x speedup)
- Better error messages

Bug fixes:
- Fixed crash on invalid torrent files
- Fixed memory leak in cache cleanup"
```

Bad example:
```bash
bash scripts/build.sh -v v1.1.0 -r -n "Updates"
```

### 4. Tag Releases Properly

```bash
# Create git tag before release
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# Then create GitHub release
bash scripts/build.sh -v v1.0.0 -r
```

### 5. Check Release Status

```bash
# View all releases
gh release list --repo german2285/TorrPlayer

# View specific release
gh release view v1.0.0 --repo german2285/TorrPlayer

# Delete release if needed
gh release delete v1.0.0 --repo german2285/TorrPlayer
```

## Troubleshooting Examples

### Release Already Exists

```bash
# Option 1: Delete and recreate
gh release delete v1.0.0 --repo german2285/TorrPlayer
bash scripts/build.sh -v v1.0.0 -r

# Option 2: Use different version
bash scripts/build.sh -v v1.0.1 -r
```

### Authentication Error

```bash
# Re-authenticate
gh auth logout
gh auth login

# Verify authentication
gh auth status
```

### Build Failed

```bash
# Check prerequisites
which wails
which npm
which go

# Check environment
echo $GOOS
echo $GOARCH
echo $CGO_ENABLED

# Clean and retry
rm -rf build/
bash scripts/build.sh
```

## Complete Example Session

Here's a complete example of developing and releasing v1.2.0:

```bash
# 1. Development phase
cd /root/TorrPlayer
bash scripts/build-dev.sh
# (Test and develop)

# 2. Create dev preview
bash scripts/build-dev.sh -v v1.2.0-dev -r -n "Development preview"

# 3. Create release candidate
bash scripts/build-dev.sh -v v1.2.0-rc -r -n "Release candidate for v1.2.0"
# (Community testing)

# 4. Fix any issues found in RC
# (Make code changes)

# 5. Create final release
bash scripts/build.sh -v v1.2.0 -r -n "TorrPlayer v1.2.0

New features:
- Enhanced cache management
- Improved UI responsiveness

Bug fixes:
- Fixed memory leaks
- Fixed settings persistence

Performance:
- 30% faster metadata loading
- Reduced memory usage"

# 6. Verify release
gh release view v1.2.0 --repo german2285/TorrPlayer

# 7. Start development on next version
bash scripts/build-dev.sh -v v1.3.0-dev
```

## Links

- [Build and Release Guide](BUILD_RELEASE.md)
- [Quick Reference](QUICK_REFERENCE.md)
- [GitHub Releases](https://github.com/german2285/TorrPlayer/releases)
