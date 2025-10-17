package torr

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"

	"github.com/german2285/TorrPlayer/pkg/server/log"
	sets "github.com/german2285/TorrPlayer/pkg/server/settings"
)

var bts *BTServer
var metadataLoadedCallback func(hash string, peers int, fileCount int, totalSize int64)

func InitApiHelper(bt *BTServer) {
	bts = bt
}

// SetMetadataLoadedCallback sets callback for metadata loaded events
func SetMetadataLoadedCallback(callback func(hash string, peers int, fileCount int, totalSize int64)) {
	metadataLoadedCallback = callback
}

func LoadTorrent(tor *Torrent) *Torrent {
	if tor.TorrentSpec == nil {
		return nil
	}
	tr, err := NewTorrent(tor.TorrentSpec, bts)
	if err != nil {
		return nil
	}
	if !tr.WaitInfo() {
		return nil
	}
	tr.Title = tor.Title
	tr.Poster = tor.Poster
	tr.Data = tor.Data
	return tr
}

func AddTorrent(spec *torrent.TorrentSpec, title, poster string, data string, category string) (*Torrent, error) {
	torr, err := NewTorrent(spec, bts)
	if err != nil {
		log.TLogln("error add torrent:", err)
		return nil, err
	}

	torDB := GetTorrentDB(spec.InfoHash)

	if torr.Title == "" {
		torr.Title = title
		if title == "" && torDB != nil {
			torr.Title = torDB.Title
		}
		if torr.Title == "" && torr.Torrent != nil && torr.Torrent.Info() != nil {
			torr.Title = torr.Info().Name
		}
	}

	if torr.Category == "" {
		torr.Category = category
		if torr.Category == "" && torDB != nil {
			torr.Category = torDB.Category
		}
	}

	if torr.Poster == "" {
		torr.Poster = poster
		if torr.Poster == "" && torDB != nil {
			torr.Poster = torDB.Poster
		}
	}

	if torr.Data == "" {
		torr.Data = data
		if torr.Data == "" && torDB != nil {
			torr.Data = torDB.Data
		}
	}

	return torr, nil
}

func SaveTorrentToDB(torr *Torrent) {
	log.TLogln("save to db:", torr.Hash())
	AddTorrentDB(torr)
}

// LoadTorrentsFromDBInstant loads torrent metadata from DB instantly without waiting for DHT
// Returns list of torrents with basic info (Title, Poster, Category, etc.)
func LoadTorrentsFromDBInstant() []*Torrent {
	log.TLogln("Loading torrents from DB (instant mode)...")
	dblist := ListTorrentsDB()
	loaded := make([]*Torrent, 0)

	for _, dbTor := range dblist {
		// Just return the DB record without loading into BTServer
		// This allows instant UI display
		log.TLogln("Loaded torrent metadata from DB:", dbTor.TorrentSpec.InfoHash.HexString(), "Title:", dbTor.Title)
		loaded = append(loaded, dbTor)
	}

	log.TLogln("Loaded", len(loaded), "torrents metadata from DB (instant)")
	return loaded
}

// LoadTorrentMetadataAsync loads full torrent metadata asynchronously in background
// This function starts goroutines for each torrent to fetch DHT metadata
func LoadTorrentMetadataAsync(hash string) {
	go func() {
		log.TLogln("Loading torrent metadata asynchronously:", hash)
		hashObj := metainfo.NewHashFromHex(hash)

		// Skip if already loaded
		if bts.GetTorrent(hashObj) != nil {
			log.TLogln("Torrent already in memory:", hash)
			return
		}

		// Get torrent from DB
		dbTor := GetTorrentDB(hashObj)
		if dbTor == nil {
			log.TLogln("Torrent not found in DB:", hash)
			return
		}

		// Load torrent into BTServer (waits for DHT metadata)
		log.TLogln("Loading torrent into BTServer:", hash, "Title:", dbTor.Title)
		tor := LoadTorrent(dbTor)
		if tor != nil {
			log.TLogln("Torrent metadata loaded successfully:", hash)

			// Collect metadata
			peers := 0
			fileCount := 0
			totalSize := int64(0)

			if tor.Torrent != nil {
				stats := tor.Torrent.Stats()
				peers = stats.ConnectedSeeders + stats.ActivePeers

				if tor.Torrent.Info() != nil {
					files := tor.Torrent.Files()
					fileCount = len(files)
					totalSize = tor.Torrent.Length()
				}
			}

			log.TLogln("Torrent metadata - Hash:", hash, "Peers:", peers, "Files:", fileCount, "Size:", totalSize)

			// Call callback if set
			if metadataLoadedCallback != nil {
				metadataLoadedCallback(hash, peers, fileCount, totalSize)
			}
		} else {
			log.TLogln("Failed to load torrent metadata:", hash)
		}
	}()
}

// LoadTorrentsFromDB - DEPRECATED: Use LoadTorrentsFromDBInstant + LoadTorrentMetadataAsync
// This function is kept for backward compatibility
func LoadTorrentsFromDB() []*Torrent {
	log.TLogln("Loading torrents from DB...")
	dblist := ListTorrentsDB()
	loaded := make([]*Torrent, 0)

	// Check if BT server is initialized
	if bts == nil {
		log.TLogln("BT server not initialized, returning DB list only")
		for _, dbTor := range dblist {
			loaded = append(loaded, dbTor)
		}
		return loaded
	}

	for hash, dbTor := range dblist {
		// Skip if already in memory
		if bts.GetTorrent(hash) != nil {
			log.TLogln("Torrent already in memory:", hash.HexString())
			continue
		}

		// Load torrent into BTServer
		log.TLogln("Loading torrent from DB:", hash.HexString(), "Title:", dbTor.Title)
		tor := LoadTorrent(dbTor)
		if tor != nil {
			loaded = append(loaded, tor)
			log.TLogln("Torrent loaded successfully:", hash.HexString())
		} else {
			log.TLogln("Failed to load torrent:", hash.HexString())
		}
	}

	log.TLogln("Loaded", len(loaded), "torrents from DB")
	return loaded
}

func GetTorrent(hashHex string) *Torrent {
	hash := metainfo.NewHashFromHex(hashHex)
	timeout := time.Second * time.Duration(sets.BTsets.TorrentDisconnectTimeout)
	if timeout > time.Minute {
		timeout = time.Minute
	}

	// Check if BT server is initialized
	var tor *Torrent
	if bts != nil {
		tor = bts.GetTorrent(hash)
		if tor != nil {
			tor.AddExpiredTime(timeout)
			return tor
		}
	}

	tr := GetTorrentDB(hash)
	if tr != nil {
		tor = tr
		// Only start background loading if BT server is ready
		if bts != nil {
			go func() {
				log.TLogln("New torrent", tor.Hash())
				tr, _ := NewTorrent(tor.TorrentSpec, bts)
				if tr != nil {
					tr.Title = tor.Title
					tr.Poster = tor.Poster
					tr.Data = tor.Data
					tr.Size = tor.Size
					tr.Timestamp = tor.Timestamp
					tr.Category = tor.Category
					tr.GotInfo()
				}
			}()
		}
	}
	return tor
}

func SetTorrent(hashHex, title, poster, category string, data string) *Torrent {
	hash := metainfo.NewHashFromHex(hashHex)
	var torr *Torrent
	if bts != nil {
		torr = bts.GetTorrent(hash)
	}
	torrDb := GetTorrentDB(hash)

	if title == "" && torr == nil && torrDb != nil {
		torr = GetTorrent(hashHex)
		torr.GotInfo()
		if torr.Torrent != nil && torr.Torrent.Info() != nil {
			title = torr.Info().Name
		}
	}

	if torr != nil {
		if title == "" && torr.Torrent != nil && torr.Torrent.Info() != nil {
			title = torr.Info().Name
		}
		torr.Title = title
		torr.Poster = poster
		torr.Category = category
		if data != "" {
			torr.Data = data
		}
	}
	// update torrent data in DB
	if torrDb != nil {
		torrDb.Title = title
		torrDb.Poster = poster
		torrDb.Category = category
		if data != "" {
			torrDb.Data = data
		}
		AddTorrentDB(torrDb)
	}
	if torr != nil {
		return torr
	} else {
		return torrDb
	}
}

func RemTorrent(hashHex string) {
	if sets.ReadOnly {
		log.TLogln("API RemTorrent: Read-only DB mode!", hashHex)
		return
	}
	hash := metainfo.NewHashFromHex(hashHex)

	// Remove from BT server if initialized
	removed := false
	if bts != nil {
		removed = bts.RemoveTorrent(hash)
	}

	if removed {
		if sets.BTsets.UseDisk && hashHex != "" && hashHex != "/" {
			name := filepath.Join(sets.BTsets.TorrentsSavePath, hashHex)
			ff, _ := os.ReadDir(name)
			for _, f := range ff {
				os.Remove(filepath.Join(name, f.Name()))
			}
			err := os.Remove(name)
			if err != nil {
				log.TLogln("Error remove cache:", err)
			}
		}
	}

	// Always remove from DB
	RemTorrentDB(hash)
}

func ListTorrent() []*Torrent {
	// Check if BT server is initialized
	btlist := make(map[metainfo.Hash]*Torrent)
	if bts != nil {
		btlist = bts.ListTorrents()
	}

	dblist := ListTorrentsDB()

	for hash, t := range dblist {
		if _, ok := btlist[hash]; !ok {
			btlist[hash] = t
		}
	}
	var ret []*Torrent

	for _, t := range btlist {
		ret = append(ret, t)
	}

	sort.Slice(ret, func(i, j int) bool {
		if ret[i].Timestamp != ret[j].Timestamp {
			return ret[i].Timestamp > ret[j].Timestamp
		} else {
			return ret[i].Title > ret[j].Title
		}
	})

	return ret
}

func DropTorrent(hashHex string) {
	if bts == nil {
		log.TLogln("DropTorrent: BT server not initialized")
		return
	}
	hash := metainfo.NewHashFromHex(hashHex)
	bts.RemoveTorrent(hash)
}

func SetSettings(set *sets.BTSets) {
	if sets.ReadOnly {
		log.TLogln("API SetSettings: Read-only DB mode!")
		return
	}
	sets.SetBTSets(set)
	log.TLogln("drop all torrents")
	dropAllTorrent()
	time.Sleep(time.Second * 1)
	log.TLogln("disconect")
	bts.Disconnect()
	log.TLogln("connect")
	bts.Connect()
	time.Sleep(time.Second * 1)
	log.TLogln("end set settings")
}

func SetDefSettings() {
	if sets.ReadOnly {
		log.TLogln("API SetDefSettings: Read-only DB mode!")
		return
	}
	sets.SetDefaultConfig()
	log.TLogln("drop all torrents")
	dropAllTorrent()
	time.Sleep(time.Second * 1)
	log.TLogln("disconect")
	bts.Disconnect()
	log.TLogln("connect")
	bts.Connect()
	time.Sleep(time.Second * 1)
	log.TLogln("end set default settings")
}

func dropAllTorrent() {
	for _, torr := range bts.torrents {
		torr.drop()
		<-torr.closed
	}
}

func Shutdown() {
	bts.Disconnect()
	sets.CloseDB()
	log.TLogln("Received shutdown. Quit")
	os.Exit(0)
}

func WriteStatus(w io.Writer) {
	bts.client.WriteStatus(w)
}

func Preload(torr *Torrent, index int) {
	cache := float32(sets.BTsets.CacheSize)
	preload := float32(sets.BTsets.PreloadCache)
	size := int64((cache / 100.0) * preload)
	if size <= 0 {
		return
	}
	if size > sets.BTsets.CacheSize {
		size = sets.BTsets.CacheSize
	}
	torr.Preload(index, size)
}
