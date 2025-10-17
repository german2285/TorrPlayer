#!/bin/bash

echo "================================================"
echo "  Building TorrPlayer Merged for Windows"
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

echo "[1/3] Installing frontend dependencies..."
cd frontend
npm install
if [ $? -ne 0 ]; then
    echo "ERROR: npm install failed"
    exit 1
fi

echo ""
echo "[2/3] Building frontend..."
npm run build
if [ $? -ne 0 ]; then
    echo "ERROR: Frontend build failed"
    exit 1
fi
cd ..

echo ""
echo "[3/3] Building Go application with Wails..."
wails build -clean -platform windows/amd64
if [ $? -ne 0 ]; then
    echo "ERROR: Wails build failed"
    exit 1
fi

echo ""
echo "[4/4] Copying libmpv-2.dll to build directory..."
if [ -f "libmpv-2.dll" ]; then
    cp libmpv-2.dll build/bin/
    echo "libmpv-2.dll copied successfully"
else
    echo "WARNING: libmpv-2.dll not found in current directory"
fi

echo ""
echo "================================================"
echo "  Build completed successfully!"
echo "================================================"
echo ""
echo "Executable location: build/bin/torrplayer-merged.exe"
echo "Library location: build/bin/libmpv-2.dll"
echo ""
