package main

import (
	"github.com/dustin/go-humanize"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"torrplayer-merged/server/settings"
)

// GetSettings returns current settings
func (a *App) GetSettings() *Settings {
	btsets := settings.BTsets
	if btsets == nil {
		return &Settings{}
	}

	return &Settings{
		CacheSize:        btsets.CacheSize,
		CacheSizeStr:     humanize.Bytes(uint64(btsets.CacheSize)),
		ConnectionsLimit: btsets.ConnectionsLimit,
		DownloadRate:     btsets.DownloadRateLimit,
		UploadRate:       btsets.UploadRateLimit,
		PreloadCache:     btsets.PreloadCache,
		RetrackersMode:   btsets.RetrackersMode,
	}
}

// SetSettings updates settings
func (a *App) SetSettings(s *Settings) error {
	runtime.LogInfo(a.ctx, "Updating settings")

	btsets := settings.BTsets
	if btsets == nil {
		btsets = &settings.BTSets{}
	}

	btsets.CacheSize = s.CacheSize
	btsets.ConnectionsLimit = s.ConnectionsLimit
	btsets.DownloadRateLimit = s.DownloadRate
	btsets.UploadRateLimit = s.UploadRate
	btsets.PreloadCache = s.PreloadCache
	btsets.RetrackersMode = s.RetrackersMode

	settings.SetBTSets(btsets)

	runtime.LogInfo(a.ctx, "Settings updated")
	return nil
}
