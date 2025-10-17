# TorrPlayer Documentation

Welcome to the TorrPlayer documentation!

## Documentation Index

### Build and Release
- **[BUILD_RELEASE.md](BUILD_RELEASE.md)** - Complete guide for building and releasing TorrPlayer
  - Prerequisites and setup
  - Detailed build instructions
  - GitHub release automation
  - CI/CD integration examples
  - Troubleshooting guide

- **[QUICK_REFERENCE.md](QUICK_REFERENCE.md)** - Quick command reference
  - Common commands
  - Build script parameters
  - Version naming conventions
  - Quick troubleshooting

- **[EXAMPLES.md](EXAMPLES.md)** - Real-world examples and workflows
  - Local development examples
  - GitHub release examples
  - Complete release workflows
  - Special case scenarios
  - CI/CD automation examples

- **[BUILD_WORKFLOW.md](BUILD_WORKFLOW.md)** - Visual workflow diagrams
  - Build flow diagrams
  - Decision trees
  - Error handling flows
  - File organization charts
  - Best practices workflows

## Quick Start

### For Developers

**Build for local development:**
```bash
cd /root/TorrPlayer
bash scripts/build-dev.sh
```

**Create a production release:**
```bash
bash scripts/build.sh -v v1.0.0 -r -n "Release description"
```

### For Contributors

1. Read [BUILD_RELEASE.md](BUILD_RELEASE.md) for detailed instructions
2. Check [QUICK_REFERENCE.md](QUICK_REFERENCE.md) for command syntax
3. Follow semantic versioning for version numbers
4. Test with development builds before creating production releases

## Build Scripts

### Production Build (`scripts/build.sh`)
Creates optimized production builds without debug features.

**Usage:**
```bash
bash scripts/build.sh [-v VERSION] [-r] [-n "NOTES"]
```

**Features:**
- Production-optimized binary
- No console window
- No debug tools
- Creates full GitHub releases

### Development Build (`scripts/build-dev.sh`)
Creates debug builds with developer tools.

**Usage:**
```bash
bash scripts/build-dev.sh [-v VERSION] [-r] [-n "NOTES"]
```

**Features:**
- Debug console window
- DevTools enabled (F12)
- Verbose logging
- Creates GitHub pre-releases

## Documentation Structure

```
docs/
├── README.md                      # This file - Documentation index
├── BUILD_RELEASE.md               # Complete build and release guide
├── QUICK_REFERENCE.md             # Quick command reference
├── EXAMPLES.md                    # Real-world examples and workflows
├── BUILD_WORKFLOW.md              # Visual workflow diagrams
└── CHANGELOG_BUILD_SCRIPTS.md     # Build scripts version history
```

## External Resources

- [Main README](../README.md) - Project overview and usage
- [GitHub Repository](https://github.com/german2285/TorrPlayer)
- [Releases Page](https://github.com/german2285/TorrPlayer/releases)
- [Wails Documentation](https://wails.io/docs/introduction/)
- [GitHub CLI Manual](https://cli.github.com/manual/)

## Getting Help

- **Build Issues**: See [BUILD_RELEASE.md - Troubleshooting](BUILD_RELEASE.md#troubleshooting)
- **Quick Commands**: See [QUICK_REFERENCE.md](QUICK_REFERENCE.md)
- **GitHub Issues**: [Report a bug](https://github.com/german2285/TorrPlayer/issues)

## Contributing

When contributing to documentation:
1. Keep examples practical and tested
2. Update this index when adding new docs
3. Use clear, concise language
4. Include code examples where relevant
5. Test all commands before documenting

## License

This documentation is part of the TorrPlayer project and follows the same license (GPL-3.0).
