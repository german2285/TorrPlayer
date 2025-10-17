// Torrent type matching Go backend
export interface Torrent {
  hash: string
  name: string
  title: string
  size: number
  sizeStr: string
  status: string
  progress: number
  downSpeed: number
  upSpeed: number
  downSpeedStr: string
  upSpeedStr: string
  peers: number
  seeders: number
  category: string
  poster: string
  timestamp: number
}

// TorrentFile type matching Go backend
export interface TorrentFile {
  index: number
  path: string
  size: number
  sizeStr: string
}

// Settings type matching Go backend
export interface Settings {
  cacheSize: number
  cacheSizeStr: string
  connectionsLimit: number
  downloadRate: number
  uploadRate: number
  preloadCache: number
  retrackersMode: number
}

// Legacy Movie type (for backwards compatibility)
export interface Movie {
  id: number
  title: string
  year: string
  duration: string
  poster: string
  category: string
  rating: string
  genres: string[]
  description: string
  ageRating: string
}
