package app

// Torrent represents a torrent in the UI
type Torrent struct {
	Hash         string  `json:"hash"`
	Name         string  `json:"name"`
	Title        string  `json:"title"`
	Size         int64   `json:"size"`
	SizeStr      string  `json:"sizeStr"`
	Status       string  `json:"status"`
	Progress     float64 `json:"progress"`
	DownSpeed    float64 `json:"downSpeed"`
	UpSpeed      float64 `json:"upSpeed"`
	DownSpeedStr string  `json:"downSpeedStr"`
	UpSpeedStr   string  `json:"upSpeedStr"`
	Peers        int     `json:"peers"`
	Seeders      int     `json:"seeders"`
	FileCount    int     `json:"fileCount"`    // Количество файлов
	Category     string  `json:"category"`
	Poster       string  `json:"poster"`
	Timestamp    int64   `json:"timestamp"`
	LoadingMeta  bool    `json:"loadingMeta"`  // Флаг загрузки метаданных
}

// TorrentFile represents a file inside a torrent
type TorrentFile struct {
	Index   int    `json:"index"`
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	SizeStr string `json:"sizeStr"`
}

// TorrentStats represents real-time statistics
type TorrentStats struct {
	DownSpeed        float64 `json:"downSpeed"`
	UpSpeed          float64 `json:"upSpeed"`
	DownSpeedStr     string  `json:"downSpeedStr"`
	UpSpeedStr       string  `json:"upSpeedStr"`
	Peers            int     `json:"peers"`
	Seeders          int     `json:"seeders"`
	Downloaded       int64   `json:"downloaded"`
	DownloadedStr    string  `json:"downloadedStr"`
	CacheFilled      int64   `json:"cacheFilled"`
	CacheCapacity    int64   `json:"cacheCapacity"`
	CacheFilledStr   string  `json:"cacheFilledStr"`
	CacheCapacityStr string  `json:"cacheCapacityStr"`
}

// Settings represents app settings
type Settings struct {
	CacheSize        int64  `json:"cacheSize"`
	CacheSizeStr     string `json:"cacheSizeStr"`
	ConnectionsLimit int    `json:"connectionsLimit"`
	DownloadRate     int    `json:"downloadRate"`
	UploadRate       int    `json:"uploadRate"`
	PreloadCache     int    `json:"preloadCache"`
	RetrackersMode   int    `json:"retrackersMode"`
	ThemeColor       string `json:"themeColor"`
	BgMusicVolume    int    `json:"bgMusicVolume"`
}
