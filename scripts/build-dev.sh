#!/bin/bash

# Parse version argument
VERSION=""
CREATE_RELEASE=false
RELEASE_NOTES=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            echo "TorrPlayer Development Build Script"
            echo ""
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  -v, --version VERSION    Version tag (e.g., v1.0.0-dev)"
            echo "  -r, --release            Create GitHub pre-release"
            echo "  -n, --notes \"TEXT\"       Release notes (optional)"
            echo "  -h, --help               Show this help message"
            echo ""
            echo "Features:"
            echo "  - Debug console window enabled"
            echo "  - Developer tools (DevTools) enabled"
            echo "  - Verbose logging"
            echo "  - Marked as pre-release on GitHub"
            echo ""
            echo "Examples:"
            echo "  # Simple development build"
            echo "  $0"
            echo ""
            echo "  # Build with version package"
            echo "  $0 -v v1.0.0-dev"
            echo ""
            echo "  # Build and create GitHub pre-release"
            echo "  $0 -v v1.0.0-dev -r"
            echo ""
            echo "  # Full pre-release with custom notes"
            echo "  $0 -v v1.0.0-dev -r -n \"Testing new features\""
            echo ""
            echo "Documentation: docs/BUILD_RELEASE.md"
            exit 0
            ;;
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -r|--release)
            CREATE_RELEASE=true
            shift
            ;;
        -n|--notes)
            RELEASE_NOTES="$2"
            shift 2
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: $0 [-v|--version VERSION] [-r|--release] [-n|--notes \"Release notes\"]"
            echo "Example: $0 -v v1.0.0-dev -r -n \"Development build\""
            echo "Run '$0 --help' for more information"
            exit 1
            ;;
    esac
done

# Check if version is provided when creating release
if [ "$CREATE_RELEASE" = true ] && [ -z "$VERSION" ]; then
    echo "ERROR: Version is required when creating a release"
    echo "Usage: $0 -v VERSION -r [-n \"Release notes\"]"
    echo "Example: $0 -v v1.0.0-dev -r -n \"Development build\""
    exit 1
fi

echo "================================================"
echo "  Building TorrPlayer Merged (DEV MODE)"
if [ -n "$VERSION" ]; then
    echo "  Version: $VERSION"
fi
echo "================================================"
echo ""

# Check if running on Linux
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo "Cross-compiling for Windows from Linux..."
    export GOOS=windows
    export GOARCH=amd64
    export CGO_ENABLED=1
    export CC=x86_64-w64-mingw32-gcc
    export CXX=x86_64-w64-mingw32-g++
fi

# Add Go bin paths to PATH
export PATH=$PATH:/usr/local/go/bin:/root/go/bin

# Check if wails is available
if ! command -v wails &> /dev/null; then
    echo "ERROR: wails not found in PATH"
    echo "Please install wails or add it to PATH"
    exit 1
fi

# Check if gh CLI is available when creating release
if [ "$CREATE_RELEASE" = true ]; then
    if ! command -v gh &> /dev/null; then
        echo "ERROR: GitHub CLI (gh) not found in PATH"
        echo "Please install GitHub CLI from https://cli.github.com/"
        exit 1
    fi

    # Check if authenticated
    if ! gh auth status &> /dev/null; then
        echo "ERROR: Not authenticated with GitHub CLI"
        echo "Please run: gh auth login"
        exit 1
    fi
fi

echo "[1/3] Installing frontend dependencies..."
cd frontend
npm install
if [ $? -ne 0 ]; then
    echo "ERROR: npm install failed"
    exit 1
fi
cd ..

echo ""
echo "[2/3] Copying frontend dist for embedding..."
mkdir -p cmd/torrplayer/frontend
cp -r frontend/dist cmd/torrplayer/frontend/
if [ $? -ne 0 ]; then
    echo "ERROR: Failed to copy frontend dist"
    exit 1
fi
echo "Frontend dist copied successfully"

echo ""
echo "[3/3] Building Go application with Wails (DEV MODE)..."
cd cmd/torrplayer
wails build -clean -platform windows/amd64 -debug -devtools
if [ $? -ne 0 ]; then
    echo "ERROR: Wails build failed"
    exit 1
fi
cd ../..

echo ""
echo "[4/4] Copying libmpv-2.dll to build directory..."
if [ -f "third_party/libmpv-2.dll" ]; then
    cp third_party/libmpv-2.dll build/bin/
    echo "libmpv-2.dll copied successfully"
else
    echo "WARNING: libmpv-2.dll not found in third_party directory"
fi

# Create release package if version is specified
if [ -n "$VERSION" ]; then
    echo ""
    echo "[4/4] Creating development release package..."

    RELEASE_DIR="build/release"
    RELEASE_NAME="TorrPlayer-${VERSION}-windows-amd64-dev"
    RELEASE_PATH="${RELEASE_DIR}/${RELEASE_NAME}"

    # Clean and create release directory
    rm -rf "$RELEASE_DIR"
    mkdir -p "$RELEASE_PATH"

    # Copy files to release directory
    cp build/bin/torrplayer-merged.exe "${RELEASE_PATH}/"
    if [ -f "build/bin/libmpv-2.dll" ]; then
        cp build/bin/libmpv-2.dll "${RELEASE_PATH}/"
    fi

    # Copy additional files
    if [ -f "README.md" ]; then
        cp README.md "${RELEASE_PATH}/"
    fi
    if [ -f "LICENSE" ]; then
        cp LICENSE "${RELEASE_PATH}/"
    fi

    # Create development build notice
    cat > "${RELEASE_PATH}/DEV-BUILD-NOTICE.txt" << EOF
This is a DEVELOPMENT BUILD of TorrPlayer.

Features:
- Debug console window enabled
- Developer tools (DevTools) enabled (press F12)
- Verbose logging enabled

This build is intended for testing and development purposes only.
For production use, please download the stable release.
EOF

    # Create ZIP archive
    cd "$RELEASE_DIR"
    ZIP_FILE="${RELEASE_NAME}.zip"
    zip -r "$ZIP_FILE" "$RELEASE_NAME"
    cd ../..

    echo "Development release package created: ${RELEASE_DIR}/${ZIP_FILE}"

    # Create GitHub release if requested
    if [ "$CREATE_RELEASE" = true ]; then
        echo ""
        echo "Creating GitHub pre-release..."

        # Set default release notes if not provided
        if [ -z "$RELEASE_NOTES" ]; then
            RELEASE_NOTES="Development build ${VERSION}

This is a development build with debug features enabled:
- Debug console window
- Developer tools (F12)
- Verbose logging

⚠️ This build is for testing purposes only."
        fi

        # Create pre-release
        gh release create "$VERSION" \
            "${RELEASE_DIR}/${ZIP_FILE}" \
            --title "TorrPlayer ${VERSION} (Development)" \
            --notes "$RELEASE_NOTES" \
            --prerelease \
            --repo german2285/TorrPlayer

        if [ $? -eq 0 ]; then
            echo "GitHub pre-release created successfully!"
            echo "View release: https://github.com/german2285/TorrPlayer/releases/tag/${VERSION}"
        else
            echo "ERROR: Failed to create GitHub pre-release"
            exit 1
        fi
    fi
fi

echo ""
echo "================================================"
echo "  DEV Build completed successfully!"
echo "================================================"
echo ""
echo "Executable location: build/bin/torrplayer-merged.exe"
echo "Library location: build/bin/libmpv-2.dll"
if [ -n "$VERSION" ]; then
    echo "Release package: build/release/${RELEASE_NAME}.zip"
fi
echo ""
echo "NOTE: This is a DEBUG build with console window"
echo "      and developer tools enabled."
echo ""
