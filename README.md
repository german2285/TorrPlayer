# TorrPlayer

Desktop application for streaming torrents with integrated MPV player.

## What it does

TorrPlayer allows you to watch video content from torrents without waiting for the full download. The application loads torrent metadata, caches chunks in RAM, and streams video directly to the built-in MPV player.

## Requirements

### For building:
- Go 1.23 or newer
- Node.js 16 or newer
- Wails CLI
- MinGW-w64 (for Windows builds with CGO)

### For running:
- Windows 10/11 (64-bit)
- libmpv-2.dll (must be placed next to the executable)

## Installation

Install dependencies:

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Install MinGW-w64 on Linux for cross-compilation
sudo apt-get install mingw-w64  # Ubuntu/Debian
```

## Building

On Windows:
```batch
build.bat
```

On Linux (cross-compile for Windows):
```bash
bash build.sh
```

The compiled application will be in `build/bin/torrplayer-merged.exe`

After building, copy `libmpv-2.dll` to the `build/bin/` directory.

## Usage

### Starting the application

1. Run `torrplayer-merged.exe`
2. Make sure `libmpv-2.dll` is in the same directory

### Adding a torrent

1. Click "Add Torrent" button
2. Enter one of the following:
   - Magnet link: `magnet:?xt=urn:btih:HASH...`
   - Path to .torrent file: `C:\Downloads\movie.torrent`
   - Torrent hash (40 characters)
3. Wait for torrent metadata to load

### Playing files

1. Click on a torrent card
2. Select a file from the list
3. Click "Play"
4. The application will start caching and launch MPV player

### Settings

Open settings page to configure:

- Cache Size: RAM cache size (64 MB - 2 GB)
- Connections Limit: Maximum peer connections (10-100)
- Download Rate: Download speed limit in KB/s (0 = unlimited)
- Upload Rate: Upload speed limit in KB/s (0 = unlimited)
- Preload Cache: Cache preload percentage (0-100%)
- Retrackers Mode:
  - 0 = do not add
  - 1 = add retrackers (default)
  - 2 = remove retrackers
  - 3 = replace retrackers

## MPV player controls

During playback:
- Space: pause/play
- Left/Right arrows: seek 5 seconds
- Up/Down arrows: volume control
- F: fullscreen mode
- M: mute audio
- S: take screenshot
- Q: quit

## How it works

The application uses a BitTorrent client to download torrent pieces. Instead of writing to disk, pieces are cached in RAM. When you play a file, a local HTTP server streams the cached data to the MPV player. The player can start immediately while the torrent continues downloading in the background.

Database stores torrent metadata, user settings, and viewing history using BoltDB and JSON files.

## Project structure

```
TorrPlayer/
├── cmd/
│   └── torrplayer/       - Application entry point
│       ├── main.go       - Main function
│       └── wails.json    - Wails configuration
├── internal/
│   ├── app/              - Application logic
│   │   ├── app.go        - Core app structure
│   │   ├── types.go      - Type definitions
│   │   ├── events.go     - Event handling
│   │   ├── torrent_handler.go   - Torrent operations
│   │   ├── settings_handler.go  - Settings management
│   │   ├── stats_handler.go     - Statistics
│   │   └── stream_server.go     - HTTP streaming
│   └── player/           - MPV player integration
│       ├── player.go     - MPV implementation (Windows)
│       └── player_stub.go - MPV stub (non-Windows)
├── pkg/
│   └── server/           - TorrServer modules
│       ├── torr/         - BitTorrent client
│       ├── settings/     - Settings and database
│       ├── log/          - Logging
│       ├── utils/        - Utilities
│       ├── mimetype/     - MIME type detection
│       └── ffprobe/      - Media analysis
├── frontend/             - Vue.js interface
│   ├── src/
│   ├── index.html
│   └── package.json
├── third_party/          - External dependencies
│   ├── mpv/
│   │   └── client.h      - libmpv header
│   ├── libmpv-2.a        - MPV static library
│   └── libmpv-2.dll      - MPV dynamic library
├── scripts/              - Build scripts
│   ├── build.bat         - Windows build
│   ├── build.sh          - Linux cross-compile
│   └── build-dev.sh      - Development build
├── docs/                 - Documentation
├── build/                - Build output
│   └── windows/
│       └── icon.ico
└── go.mod                - Go dependencies
```

## Architecture

- Frontend: Vue.js 3 with Vite
- Backend: Go with Wails v2 framework
- BitTorrent: anacrolix/torrent library
- Cache: In-memory with automatic cleanup
- Streaming: Local HTTP server
- Player: MPV via CGO and libmpv-2.dll
- Database: BoltDB and JSON

## Troubleshooting

### MPV not starting

1. Check that libmpv-2.dll is in the same directory as the executable
2. Verify DLL is 64-bit version
3. Install Visual C++ Redistributable
4. Check logs in application DevTools (F12)

### Build errors

**gcc not found**
- Install MinGW-w64 and add to PATH

**wails command not found**
- Install Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

**mpv/client.h not found**
- Ensure mpv/client.h exists in project directory

**npm not found**
- Install Node.js and add to PATH

## License

Based on TorrServer (GPL-3.0) and Wails framework (MIT).

## Links

- TorrServer: https://github.com/YouROK/TorrServer
- Wails: https://wails.io/
- MPV: https://mpv.io/
- libmpv for Windows: https://sourceforge.net/projects/mpv-player-windows/files/libmpv/
