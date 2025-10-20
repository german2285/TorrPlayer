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

// GetTorrents returns list of all torrents (simplified - no metadata)
func (a *App) GetTorrents() []Torrent {
	// Check if context is initialized
	if a.ctx == nil {
		return []Torrent{}
	}

	// Check if BT server is initialized
	if a.btServer == nil {
		return []Torrent{}
	}

	list := torrserv.ListTorrent()
	if list == nil {
		return []Torrent{}
	}

	torrents := make([]Torrent, 0, len(list))

	for _, tor := range list {
		if tor == nil {
			continue
		}

		// Get hash
		hashStr := "unknown"
		if tor.TorrentSpec != nil && tor.TorrentSpec.InfoHash != (metainfo.Hash{}) {
			hashStr = tor.TorrentSpec.InfoHash.HexString()
		}

		// Get name from DB or spec
		name := tor.Title
		if name == "" && tor.TorrentSpec != nil && tor.TorrentSpec.DisplayName != "" {
			name = tor.TorrentSpec.DisplayName
		}
		if name == "" {
			name = "Unknown"
		}

		// Get basic info from DB
		category := tor.Category
		poster := tor.Poster
		timestamp := tor.Timestamp
		if timestamp == 0 {
			timestamp = time.Now().Unix()
		}

		result := Torrent{
			Hash:      hashStr,
			Name:      name,
			Title:     tor.Title,
			Category:  category,
			Poster:    poster,
			Timestamp: timestamp,
			// No metadata fields
			Size:         0,
			SizeStr:      "",
			Peers:        0,
			Seeders:      0,
			FileCount:    0,
			Status:       "ready",
			Progress:     0,
			DownSpeed:    0,
			UpSpeed:      0,
			DownSpeedStr: "",
			UpSpeedStr:   "",
			LoadingMeta:  false,
		}

		torrents = append(torrents, result)
	}

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
