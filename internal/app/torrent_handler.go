package app

import (
	"encoding/base64"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/dustin/go-humanize"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	torrserv "github.com/german2285/TorrPlayer/pkg/server/torr"
	"github.com/german2285/TorrPlayer/pkg/server/torr/state"
	"github.com/german2285/TorrPlayer/pkg/server/utils"
)

// AddTorrent adds a torrent by magnet link, .torrent file path, or hash
func (a *App) AddTorrent(input string) (*Torrent, error) {
	// Check if context is initialized
	if a.ctx == nil {
		return nil, fmt.Errorf("application not initialized yet")
	}

	runtime.LogInfo(a.ctx, fmt.Sprintf("Adding torrent: %s", input))

	// Parse input
	spec, err := a.parseTorrentInput(input)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to parse input: %v", err))
		return nil, err
	}

	// Add torrent
	tor, err := torrserv.AddTorrent(spec, "", "", "", "")
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to add torrent: %v", err))
		return nil, err
	}

	hashStr := tor.Hash().HexString()
	runtime.LogInfo(a.ctx, fmt.Sprintf("Torrent added to client: %s", hashStr))

	// Save to DB immediately with minimal info
	torrserv.SaveTorrentToDB(tor)
	runtime.LogInfo(a.ctx, fmt.Sprintf("Torrent saved to database: %s", hashStr))

	// Start background goroutine to wait for metadata
	go func() {
		timeout := time.After(60 * time.Second)
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-timeout:
				runtime.LogInfo(a.ctx, fmt.Sprintf("Timeout waiting for metadata for: %s", hashStr))
				return
			case <-ticker.C:
				if tor.GotInfo() {
					runtime.LogInfo(a.ctx, fmt.Sprintf("Metadata received for: %s - %s", hashStr, tor.Name()))
					// Update DB with full info
					torrserv.SaveTorrentToDB(tor)
					return
				}
			}
		}
	}()

	// Return immediately with partial info
	name := "Загрузка метаданных..."
	if spec.DisplayName != "" {
		name = spec.DisplayName
	}

	return &Torrent{
		Hash:      hashStr,
		Name:      name,
		Title:     name,
		Size:      0,
		SizeStr:   "Ожидание...",
		Status:    "loading",
		Category:  "Recent",
		Timestamp: time.Now().Unix(),
	}, nil
}

// GetTorrents returns list of all torrents
func (a *App) GetTorrents() []Torrent {
	runtime.LogDebug(a.ctx, "GetTorrents called")

	// Defer panic recovery
	defer func() {
		if r := recover(); r != nil {
			runtime.LogError(a.ctx, fmt.Sprintf("PANIC in GetTorrents: %v", r))
		}
	}()

	// Check if context is initialized (startup completed)
	if a.ctx == nil {
		runtime.LogWarning(a.ctx, "Context not initialized, returning empty list")
		return []Torrent{}
	}

	// Check if BT server is initialized
	if a.btServer == nil {
		runtime.LogWarning(a.ctx, "BT server not initialized yet, returning empty list")
		return []Torrent{}
	}

	runtime.LogDebug(a.ctx, "Calling ListTorrent...")
	list := torrserv.ListTorrent()
	runtime.LogDebug(a.ctx, "ListTorrent returned successfully")

	if list == nil {
		runtime.LogDebug(a.ctx, "No torrents in list, returning empty list")
		return []Torrent{}
	}

	listLen := len(list)
	runtime.LogInfo(a.ctx, fmt.Sprintf("GetTorrents called, found %d torrents", listLen))
	runtime.LogDebug(a.ctx, "Creating torrents slice...")
	torrents := make([]Torrent, 0, listLen)
	runtime.LogDebug(a.ctx, "Torrents slice created successfully")
	runtime.LogDebug(a.ctx, "Starting torrent iteration...")

	for _, tor := range list {
		if tor == nil {
			runtime.LogWarning(a.ctx, "Skipping nil torrent in list")
			continue
		}

		// Get hash first (safe operation)
		hashStr := "unknown"
		if tor.TorrentSpec != nil && tor.TorrentSpec.InfoHash != (metainfo.Hash{}) {
			hashStr = tor.TorrentSpec.InfoHash.HexString()
		}

		runtime.LogDebug(a.ctx, fmt.Sprintf("Processing torrent: %s", hashStr))

		// Always show the torrent - it's in DB or active
		// SAFETY: tor.Status() might panic if tor.Torrent is nil
		var st *state.TorrentStatus
		if tor.Torrent != nil {
			st = tor.Status()
		} else {
			// Create empty status if torrent not loaded yet
			st = &state.TorrentStatus{
				Name:        tor.Title,
				TorrentSize: tor.Size,
			}
		}

		// Get name and size from status to avoid blocking on inactive torrents
		name := st.Name
		if name == "" && tor.TorrentSpec != nil && tor.TorrentSpec.DisplayName != "" {
			name = tor.TorrentSpec.DisplayName
		}
		if name == "" {
			name = tor.Title
		}

		size := st.TorrentSize
		if size == 0 {
			size = tor.Size
		}

		// If we still don't have enough data, try to get it from DB
		if name == "" || size == 0 {
			// Safe hash extraction
			var hash metainfo.Hash
			if tor.TorrentSpec != nil {
				hash = tor.TorrentSpec.InfoHash
			}

			runtime.LogDebug(a.ctx, fmt.Sprintf("Checking DB for torrent %s (name empty: %v, size zero: %v)", hash.HexString(), name == "", size == 0))
			dbTor := torrserv.GetTorrentDB(hash)
			if dbTor != nil {
				runtime.LogDebug(a.ctx, fmt.Sprintf("Found in DB: Title=%s, Size=%d", dbTor.Title, dbTor.Size))
				if name == "" && dbTor.Title != "" {
					name = dbTor.Title
				}
				if size == 0 && dbTor.Size > 0 {
					size = dbTor.Size
				}
			} else {
				runtime.LogDebug(a.ctx, fmt.Sprintf("Not found in DB for hash %s", hash.HexString()))
			}
		}

		if name == "" {
			name = "Загрузка метаданных..."
		}

		progress := float64(0)
		if size > 0 {
			progress = float64(st.LoadedSize) / float64(size) * 100
		}

		// Determine status and metadata loading state
		status := "ready"
		loadingMeta := false
		sizeStr := humanize.Bytes(uint64(size))
		fileCount := len(st.FileStats)

		// Check if metadata is loaded (safe check)
		hasInfo := false
		if tor.Torrent != nil {
			hasInfo = tor.GotInfo()
		}

		if size == 0 || !hasInfo {
			status = "loading"
			loadingMeta = true
			if size == 0 {
				sizeStr = "Ожидание..."
			}
			if fileCount == 0 {
				fileCount = 0 // Будет отображаться как "загрузка..."
			}
		}

		// Get poster from DB
		poster := tor.Poster
		category := tor.Category
		timestamp := tor.Timestamp
		if timestamp == 0 {
			timestamp = time.Now().Unix()
		}

		// Safe hash extraction for result
		resultHash := hashStr
		if resultHash == "unknown" && tor.TorrentSpec != nil {
			resultHash = tor.TorrentSpec.InfoHash.HexString()
		}

		result := Torrent{
			Hash:         resultHash,
			Name:         name,
			Title:        tor.Title,
			Size:         size,
			SizeStr:      sizeStr,
			Status:       status,
			Progress:     progress,
			DownSpeed:    st.DownloadSpeed,
			UpSpeed:      st.UploadSpeed,
			DownSpeedStr: humanize.Bytes(uint64(st.DownloadSpeed)) + "/s",
			UpSpeedStr:   humanize.Bytes(uint64(st.UploadSpeed)) + "/s",
			Peers:        st.ActivePeers,
			Seeders:      st.ConnectedSeeders,
			FileCount:    fileCount,
			Category:     category,
			Poster:       poster,
			Timestamp:    timestamp,
			LoadingMeta:  loadingMeta,
		}

		runtime.LogDebug(a.ctx, fmt.Sprintf("Torrent: %s, Name: %s, Status: %s, Size: %d, FileCount: %d, LoadingMeta: %v", result.Hash, result.Name, result.Status, result.Size, result.FileCount, result.LoadingMeta))
		torrents = append(torrents, result)
	}

	runtime.LogInfo(a.ctx, fmt.Sprintf("Returning %d torrents to frontend", len(torrents)))
	runtime.LogDebug(a.ctx, "About to return from GetTorrents...")
	return torrents
}

// GetTorrentFiles returns list of files in a torrent
func (a *App) GetTorrentFiles(hash string) ([]TorrentFile, error) {
	// Check if context is initialized
	if a.ctx == nil {
		return nil, fmt.Errorf("application not initialized yet")
	}

	tor := torrserv.GetTorrent(hash)
	if tor == nil {
		return nil, fmt.Errorf("torrent not found")
	}

	files := tor.Files()
	result := make([]TorrentFile, 0, len(files))

	// Sort files by path
	sort.Slice(files, func(i, j int) bool {
		return utils.CompareStrings(files[i].Path(), files[j].Path())
	})

	for i, f := range files {
		result = append(result, TorrentFile{
			Index:   i + 1,
			Path:    f.Path(),
			Size:    f.Length(),
			SizeStr: humanize.Bytes(uint64(f.Length())),
		})
	}

	return result, nil
}

// RemoveTorrent removes a torrent
func (a *App) RemoveTorrent(hash string) error {
	// Check if context is initialized
	if a.ctx == nil {
		return fmt.Errorf("application not initialized yet")
	}

	runtime.LogInfo(a.ctx, fmt.Sprintf("Removing torrent: %s", hash))
	torrserv.RemTorrent(hash)
	return nil
}

// parseTorrentInput parses magnet link, .torrent file path, or hash
func (a *App) parseTorrentInput(input string) (*torrent.TorrentSpec, error) {
	// Check if it's a base64 encoded file (from drag & drop)
	if strings.HasPrefix(input, "data:") {
		// Parse data URL: data:application/x-bittorrent;base64,<content>
		parts := strings.SplitN(input, ",", 2)
		if len(parts) == 2 {
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				return nil, fmt.Errorf("failed to decode base64: %v", err)
			}

			// Create temporary file
			tmpFile, err := os.CreateTemp("", "torrent-*.torrent")
			if err != nil {
				return nil, fmt.Errorf("failed to create temp file: %v", err)
			}
			tmpFileName := tmpFile.Name()

			// Write data and close file
			if _, err := tmpFile.Write(decoded); err != nil {
				tmpFile.Close()
				os.Remove(tmpFileName)
				return nil, fmt.Errorf("failed to write temp file: %v", err)
			}
			tmpFile.Close()

			// Ensure cleanup after reading
			defer os.Remove(tmpFileName)

			minfo, err := metainfo.LoadFromFile(tmpFileName)
			if err != nil {
				return nil, fmt.Errorf("failed to load torrent file: %v", err)
			}

			info, err := minfo.UnmarshalInfo()
			if err != nil {
				return nil, fmt.Errorf("failed to parse torrent info: %v", err)
			}

			mag := minfo.Magnet(nil, &info)
			return &torrent.TorrentSpec{
				InfoBytes:   minfo.InfoBytes,
				Trackers:    [][]string{mag.Trackers},
				DisplayName: info.Name,
				InfoHash:    minfo.HashInfoBytes(),
			}, nil
		}
	}

	// Check if it's a magnet link
	if len(input) > 8 && input[:8] == "magnet:?" {
		mag, err := metainfo.ParseMagnetUri(input)
		if err != nil {
			return nil, fmt.Errorf("invalid magnet link: %v", err)
		}
		var trackers [][]string
		if len(mag.Trackers) > 0 {
			trackers = [][]string{mag.Trackers}
		}
		return &torrent.TorrentSpec{
			InfoHash:    mag.InfoHash,
			Trackers:    trackers,
			DisplayName: mag.DisplayName,
		}, nil
	}

	// Check if it's a file path
	if _, err := os.Stat(input); err == nil {
		minfo, err := metainfo.LoadFromFile(input)
		if err != nil {
			return nil, fmt.Errorf("failed to load torrent file: %v", err)
		}

		info, err := minfo.UnmarshalInfo()
		if err != nil {
			return nil, fmt.Errorf("failed to parse torrent info: %v", err)
		}

		mag := minfo.Magnet(nil, &info)
		return &torrent.TorrentSpec{
			InfoBytes:   minfo.InfoBytes,
			Trackers:    [][]string{mag.Trackers},
			DisplayName: info.Name,
			InfoHash:    minfo.HashInfoBytes(),
		}, nil
	}

	// Try as hash
	if len(input) == 40 || len(input) == 32 {
		var hash metainfo.Hash
		if err := hash.FromHexString(input); err == nil {
			return &torrent.TorrentSpec{
				InfoHash: hash,
			}, nil
		}
	}

	return nil, fmt.Errorf("invalid input: must be magnet link, torrent file path, or hash")
}
