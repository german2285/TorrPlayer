package main

import (
	"fmt"

	"github.com/dustin/go-humanize"

	torrserv "torrplayer-merged/server/torr"
)

// GetTorrentStats returns real-time statistics for a torrent
func (a *App) GetTorrentStats(hash string) (*TorrentStats, error) {
	tor := torrserv.GetTorrent(hash)
	if tor == nil {
		return nil, fmt.Errorf("torrent not found")
	}

	st := tor.Status()
	cache := tor.GetCache()

	stats := &TorrentStats{
		DownSpeed:     st.DownloadSpeed,
		UpSpeed:       st.UploadSpeed,
		DownSpeedStr:  humanize.Bytes(uint64(st.DownloadSpeed)) + "/s",
		UpSpeedStr:    humanize.Bytes(uint64(st.UploadSpeed)) + "/s",
		Peers:         st.ActivePeers,
		Seeders:       st.ConnectedSeeders,
		Downloaded:    st.LoadedSize,
		DownloadedStr: humanize.Bytes(uint64(st.LoadedSize)),
	}

	if cache != nil {
		cacheState := cache.GetState()
		stats.CacheFilled = cacheState.Filled
		stats.CacheCapacity = cacheState.Capacity
		stats.CacheFilledStr = humanize.Bytes(uint64(cacheState.Filled))
		stats.CacheCapacityStr = humanize.Bytes(uint64(cacheState.Capacity))
	}

	return stats, nil
}
