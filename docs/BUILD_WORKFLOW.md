# Build Workflow Diagrams

This document provides visual representations of the build and release workflows.

## Table of Contents
- [Simple Build Flow](#simple-build-flow)
- [Version Package Flow](#version-package-flow)
- [GitHub Release Flow](#github-release-flow)
- [Development vs Production](#development-vs-production)
- [Complete Release Workflow](#complete-release-workflow)

## Simple Build Flow

No parameters - basic build only:

```
┌─────────────────────────────────────────────────┐
│  bash scripts/build.sh                          │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  [1/3] Install frontend dependencies            │
│  $ cd frontend && npm install                   │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  [2/3] Build frontend                           │
│  $ npm run build                                │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  [3/3] Build with Wails                         │
│  $ wails build -clean -platform windows/amd64   │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  [4/5] Copy libmpv-2.dll                        │
│  $ cp third_party/libmpv-2.dll build/bin/       │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  OUTPUT:                                        │
│  • build/bin/torrplayer-merged.exe              │
│  • build/bin/libmpv-2.dll                       │
└─────────────────────────────────────────────────┘
```

## Version Package Flow

With `-v` parameter - creates ZIP package:

```
┌─────────────────────────────────────────────────┐
│  bash scripts/build.sh -v v1.0.0                │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  Build steps [1/3] to [4/5]                     │
│  (same as simple build)                         │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  [5/5] Create release package                   │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  Create release directory                       │
│  build/release/TorrPlayer-v1.0.0-windows-amd64/ │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  Copy files:                                    │
│  • torrplayer-merged.exe                        │
│  • libmpv-2.dll                                 │
│  • README.md                                    │
│  • LICENSE                                      │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  Create ZIP archive                             │
│  TorrPlayer-v1.0.0-windows-amd64.zip            │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  OUTPUT:                                        │
│  • build/bin/torrplayer-merged.exe              │
│  • build/bin/libmpv-2.dll                       │
│  • build/release/TorrPlayer-v1.0.0-...zip       │
└─────────────────────────────────────────────────┘
```

## GitHub Release Flow

With `-v` and `-r` parameters - creates GitHub release:

```
┌─────────────────────────────────────────────────┐
│  bash scripts/build.sh -v v1.0.0 -r             │
│                        -n "Release notes"       │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  Validate prerequisites                         │
│  ✓ Check if gh CLI is installed                │
│  ✓ Check GitHub authentication                 │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  Build and package                              │
│  (steps [1/3] to [5/5])                         │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  Create GitHub release                          │
│  $ gh release create v1.0.0 \                   │
│      build/release/TorrPlayer-v1.0.0-...zip \   │
│      --title "TorrPlayer v1.0.0" \              │
│      --notes "Release notes" \                  │
│      --repo german2285/TorrPlayer               │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  SUCCESS                                        │
│  Release created at:                            │
│  https://github.com/german2285/TorrPlayer/      │
│  releases/tag/v1.0.0                            │
└─────────────────────────────────────────────────┘
```

## Development vs Production

Side-by-side comparison:

```
┌────────────────────────────┬────────────────────────────┐
│   PRODUCTION BUILD         │   DEVELOPMENT BUILD        │
│   (build.sh)               │   (build-dev.sh)           │
├────────────────────────────┼────────────────────────────┤
│                            │                            │
│  Optimized binary          │  Debug binary              │
│  No console window         │  Console window ✓          │
│  No debug tools            │  DevTools enabled (F12) ✓  │
│  Production ready          │  For testing only          │
│                            │                            │
├────────────────────────────┼────────────────────────────┤
│                            │                            │
│  Package name:             │  Package name:             │
│  TorrPlayer-v1.0.0-...     │  TorrPlayer-v1.0.0-dev-... │
│                            │                            │
├────────────────────────────┼────────────────────────────┤
│                            │                            │
│  GitHub release:           │  GitHub release:           │
│  Full release              │  Pre-release ⚠️             │
│                            │  + DEV-BUILD-NOTICE.txt    │
│                            │                            │
├────────────────────────────┼────────────────────────────┤
│                            │                            │
│  Use for:                  │  Use for:                  │
│  • Stable releases         │  • Testing                 │
│  • Production deployment   │  • Development             │
│  • Public distribution     │  • Feature previews        │
│                            │  • Bug investigation       │
│                            │                            │
└────────────────────────────┴────────────────────────────┘
```

## Complete Release Workflow

End-to-end development to release:

```
                    START
                      │
                      ▼
      ┌───────────────────────────────┐
      │  Development Phase            │
      │  $ bash scripts/build-dev.sh  │
      │  (Local testing)              │
      └───────────┬───────────────────┘
                  │
                  ▼
      ┌───────────────────────────────┐
      │  Create Dev Preview           │
      │  $ bash scripts/build-dev.sh  │
      │    -v v1.0.0-dev -r           │
      └───────────┬───────────────────┘
                  │
                  ▼
      ┌───────────────────────────────┐
      │  Community Testing            │
      │  (Users download and test)    │
      └───────────┬───────────────────┘
                  │
                  ▼
      ┌───────────────────────────────┐
      │  Create Release Candidate     │
      │  $ bash scripts/build-dev.sh  │
      │    -v v1.0.0-rc -r            │
      └───────────┬───────────────────┘
                  │
                  ▼
      ┌───────────────────────────────┐
      │  Final Testing                │
      │  (Verify all features)        │
      └───────────┬───────────────────┘
                  │
                  ▼
      ┌───────────────────────────────┐
      │  Create Production Release    │
      │  $ bash scripts/build.sh      │
      │    -v v1.0.0 -r               │
      │    -n "Release notes"         │
      └───────────┬───────────────────┘
                  │
                  ▼
      ┌───────────────────────────────┐
      │  PUBLIC RELEASE               │
      │  Available on GitHub Releases │
      └───────────────────────────────┘
                  │
                  ▼
                 END
```

## Parameter Decision Tree

Choosing the right parameters:

```
                 START
                   │
                   ▼
           ┌───────────────┐
           │ Need version  │
           │   package?    │
           └───┬───────┬───┘
               │       │
              NO      YES
               │       │
               │       ▼
               │   ┌────────────┐
               │   │ Use -v     │
               │   │ parameter  │
               │   └─────┬──────┘
               │         │
               │         ▼
               │   ┌────────────┐
               │   │ Need GitHub│
               │   │  release?  │
               │   └─┬────────┬─┘
               │     │        │
               │    NO       YES
               │     │        │
               │     │        ▼
               │     │    ┌────────────┐
               │     │    │ Add -r     │
               │     │    │ flag       │
               │     │    └─────┬──────┘
               │     │          │
               │     │          ▼
               │     │    ┌────────────┐
               │     │    │ Custom     │
               │     │    │ notes?     │
               │     │    └─┬────────┬─┘
               │     │      │        │
               │     │     NO       YES
               │     │      │        │
               │     │      │        ▼
               │     │      │    ┌────────────┐
               │     │      │    │ Add -n     │
               │     │      │    │ parameter  │
               │     │      │    └─────┬──────┘
               ▼     ▼      ▼          ▼
         ┌─────────────────────────────────┐
         │     Execute Build Command       │
         └─────────────────────────────────┘

EXAMPLES:

Simple:       build.sh
With version: build.sh -v v1.0.0
With release: build.sh -v v1.0.0 -r
Full:         build.sh -v v1.0.0 -r -n "Notes"
```

## Error Handling Flow

How errors are handled:

```
      ┌────────────────────────┐
      │  Script Execution      │
      └───────────┬────────────┘
                  │
                  ▼
      ┌────────────────────────┐
      │  Parse Parameters      │
      └───┬──────────┬─────────┘
          │          │
       VALID    INVALID
          │          │
          │          ▼
          │     ┌─────────────────┐
          │     │ Show error      │
          │     │ Show usage      │
          │     │ Exit with code 1│
          │     └─────────────────┘
          │
          ▼
      ┌────────────────────────┐
      │  Check Prerequisites   │
      └───┬──────────┬─────────┘
          │          │
         OK      MISSING
          │          │
          │          ▼
          │     ┌─────────────────┐
          │     │ Show error msg  │
          │     │ Show install    │
          │     │ instructions    │
          │     │ Exit with code 1│
          │     └─────────────────┘
          │
          ▼
      ┌────────────────────────┐
      │  Check -r flag         │
      └───┬──────────┬─────────┘
          │          │
       NO -r      -r SET
          │          │
          │          ▼
          │     ┌─────────────────┐
          │     │ Check gh CLI    │
          │     └───┬─────┬───────┘
          │         │     │
          │        OK  MISSING
          │         │     │
          │         │     ▼
          │         │  ┌──────────────┐
          │         │  │ Error: gh    │
          │         │  │ not found    │
          │         │  │ Exit code 1  │
          │         │  └──────────────┘
          │         │
          │         ▼
          │     ┌─────────────────┐
          │     │ Check auth      │
          │     └───┬─────┬───────┘
          │         │     │
          │        OK  NOT AUTH
          │         │     │
          │         │     ▼
          │         │  ┌──────────────┐
          │         │  │ Error: run   │
          │         │  │ gh auth login│
          │         │  │ Exit code 1  │
          │         │  └──────────────┘
          │         │
          ▼         ▼
      ┌────────────────────────┐
      │  Execute Build         │
      └───┬──────────┬─────────┘
          │          │
      SUCCESS    FAILURE
          │          │
          │          ▼
          │     ┌─────────────────┐
          │     │ Show error      │
          │     │ Exit with code 1│
          │     └─────────────────┘
          │
          ▼
      ┌────────────────────────┐
      │  Success               │
      │  Exit with code 0      │
      └────────────────────────┘
```

## File Organization

Directory structure after build:

```
TorrPlayer/
│
├── scripts/
│   ├── build.sh              ← Production build script
│   └── build-dev.sh          ← Development build script
│
├── docs/
│   ├── BUILD_RELEASE.md      ← Complete guide
│   ├── QUICK_REFERENCE.md    ← Quick commands
│   ├── EXAMPLES.md           ← Usage examples
│   └── BUILD_WORKFLOW.md     ← This file
│
└── build/
    ├── bin/
    │   ├── torrplayer-merged.exe
    │   └── libmpv-2.dll
    │
    └── release/
        ├── TorrPlayer-v1.0.0-windows-amd64/
        │   ├── torrplayer-merged.exe
        │   ├── libmpv-2.dll
        │   ├── README.md
        │   └── LICENSE
        │
        └── TorrPlayer-v1.0.0-windows-amd64.zip
```

## Quick Decision Guide

```
┌─────────────────────────────────────────────────┐
│  What do you want to do?                        │
└─────────────────┬───────────────────────────────┘
                  │
    ┌─────────────┼─────────────┬─────────────┐
    │             │             │             │
    ▼             ▼             ▼             ▼
┌────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐
│ Local  │  │ Version  │  │ Test     │  │ Public   │
│ build  │  │ package  │  │ release  │  │ release  │
└───┬────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘
    │            │              │              │
    ▼            ▼              ▼              ▼
build.sh    build.sh -v    build-dev.sh   build.sh -v
                v1.0.0     -v v1.0.0-dev   v1.0.0 -r
                                -r         -n "Notes"
```

## Best Practices Flow

```
┌─────────────────────────────────────────────────┐
│  Before Creating Release                        │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  1. Test with dev build                         │
│     $ build-dev.sh                              │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  2. Create dev preview                          │
│     $ build-dev.sh -v v1.0.0-dev -r             │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  3. Get feedback                                │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  4. Create RC                                   │
│     $ build-dev.sh -v v1.0.0-rc -r              │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  5. Final testing                               │
└─────────────┬───────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────┐
│  6. Create production release                   │
│     $ build.sh -v v1.0.0 -r -n "Notes"          │
└─────────────────────────────────────────────────┘
```

## Links

- [Build Guide](BUILD_RELEASE.md)
- [Quick Reference](QUICK_REFERENCE.md)
- [Examples](EXAMPLES.md)
- [Main README](../README.md)
