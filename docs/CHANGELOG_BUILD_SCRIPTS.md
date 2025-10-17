# Build Scripts Changelog

This file tracks changes to the build scripts (`scripts/build.sh` and `scripts/build-dev.sh`).

## [2.0.0] - 2025-10-17

### Added
- **Version parameter** (`-v|--version`): Specify version for builds
- **GitHub release automation** (`-r|--release`): Automatically create GitHub releases
- **Custom release notes** (`-n|--notes`): Add custom release notes
- **Help flag** (`-h|--help`): Display usage information and examples
- **ZIP packaging**: Automatic creation of release packages
  - Includes executable, DLL, README, and LICENSE
- **GitHub CLI integration**: Checks for `gh` availability and authentication
- **Pre-release support**: Development builds create GitHub pre-releases
- **DEV-BUILD-NOTICE.txt**: Automatic notice file for development builds
- **Comprehensive documentation**: BUILD_RELEASE.md, QUICK_REFERENCE.md, EXAMPLES.md

### Changed
- **Script location**: Moved from root to `scripts/` directory
- **libmpv-2.dll path**: Updated to use `third_party/libmpv-2.dll`
- **Output structure**: Added `build/release/` directory for release artifacts
- **Build steps**: Updated numbering to include packaging step
- **Error handling**: Improved error messages and validation

### Enhanced
- **User experience**: Clear progress messages and help output
- **Validation**: Check for required tools (wails, gh)
- **Authentication**: Verify GitHub CLI authentication before release
- **Documentation**: Detailed guides with examples

### Technical Details

#### build.sh (Production)
- Creates optimized production builds
- Generates release packages: `TorrPlayer-{VERSION}-windows-amd64.zip`
- Creates full GitHub releases
- Default release notes if none provided

#### build-dev.sh (Development)
- Creates debug builds with DevTools
- Generates dev packages: `TorrPlayer-{VERSION}-windows-amd64-dev.zip`
- Creates GitHub pre-releases
- Includes development build notice
- Enhanced default release notes for dev builds

### File Structure
```
scripts/
├── build.sh           # Production build script
└── build-dev.sh       # Development build script

docs/
├── BUILD_RELEASE.md   # Complete build guide
├── QUICK_REFERENCE.md # Command reference
├── EXAMPLES.md        # Usage examples
└── README.md          # Documentation index

build/
├── bin/               # Build output
│   ├── torrplayer-merged.exe
│   └── libmpv-2.dll
└── release/           # Release packages
    ├── TorrPlayer-{VERSION}-windows-amd64/
    └── TorrPlayer-{VERSION}-windows-amd64.zip
```

## [1.0.0] - Original Version

### Original Features
- Basic build functionality
- Frontend npm install and build
- Wails build for Windows
- libmpv-2.dll copying
- Cross-compilation support (Linux to Windows)

### Original Limitations
- No version management
- Manual release creation required
- No package automation
- Limited error handling
- Basic documentation

## Migration Guide

### From v1.0.0 to v2.0.0

**Old way (v1.0.0):**
```bash
bash build.sh
# Manually create ZIP
# Manually create GitHub release
# Manually upload files
```

**New way (v2.0.0):**
```bash
bash scripts/build.sh -v v1.0.0 -r -n "Release notes"
# Everything automated!
```

### Breaking Changes
None - backward compatible. Scripts still work without parameters.

## Future Enhancements

### Planned Features
- [ ] Multi-platform builds (Linux, macOS)
- [ ] Automatic changelog generation
- [ ] Digital signature support
- [ ] Checksum file generation
- [ ] Docker-based builds
- [ ] Parallel build support
- [ ] Build caching
- [ ] Incremental builds

### Potential Improvements
- [ ] Build time measurement
- [ ] Build artifact verification
- [ ] Automatic version bumping
- [ ] Integration with semantic-release
- [ ] Automated testing before release
- [ ] Rollback support
- [ ] Build notifications (Slack, Discord)

## Usage Examples

### Simple Build
```bash
bash scripts/build.sh
```

### Versioned Build
```bash
bash scripts/build.sh -v v1.0.0
```

### Full Release
```bash
bash scripts/build.sh -v v1.0.0 -r -n "Release notes"
```

### Development Build
```bash
bash scripts/build-dev.sh -v v1.0.0-dev -r
```

## Links

- [Build Guide](BUILD_RELEASE.md)
- [Quick Reference](QUICK_REFERENCE.md)
- [Examples](EXAMPLES.md)
- [GitHub Repository](https://github.com/german2285/TorrPlayer)

## Contributing

When modifying build scripts:
1. Update version in this changelog
2. Test all parameter combinations
3. Update documentation
4. Test on Linux (cross-compilation)
5. Verify GitHub release creation

## Support

For issues or questions:
- Check [BUILD_RELEASE.md - Troubleshooting](BUILD_RELEASE.md#troubleshooting)
- Open an issue: https://github.com/german2285/TorrPlayer/issues
