package main

// TorrentMetadataLoadedEvent represents metadata loaded event
type TorrentMetadataLoadedEvent struct {
	Hash      string `json:"hash"`
	Title     string `json:"title"`
	Peers     int    `json:"peers"`
	Seeders   int    `json:"seeders"`
	FileCount int    `json:"fileCount"`
	TotalSize int64  `json:"totalSize"`
	SizeStr   string `json:"sizeStr"`
	Loaded    bool   `json:"loaded"`
}
