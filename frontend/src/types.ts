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
  // Existing settings
  cacheSize: number
  cacheSizeStr: string
  connectionsLimit: number
  downloadRate: number
  uploadRate: number
  preloadCache: number
  retrackersMode: number
  themeColor: string
  bgMusicVolume: number

  // Storage settings
  useDisk: boolean
  torrentsSavePath: string

  // Network settings
  peersListenPort: number

  // Protocol settings
  disableDHT: boolean
  disablePEX: boolean
  disableUTP: boolean
  disableUPNP: boolean
  disableTCP: boolean

  // Upload and encryption
  disableUpload: boolean
  forceEncrypt: boolean
  enableIPv6: boolean

  // Advanced settings
  enableDebug: boolean
  readerReadAHead: number
  removeCacheOnDrop: boolean
  responsiveMode: boolean
  torrentDisconnectTimeout: number
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

// RuTracker types
export interface RutrackerTorrent {
  topicId: string
  title: string
  category: string
  size: string
  seeds: number
  leeches: number
  author: string
  date: string
}

export interface CaptchaData {
  imageBase64: string
  sid: string
  codeField: string
}

export interface LoginData {
  username: string
  password: string
}

export interface RegistrationData {
  username: string
  password: string
  email: string
  captchaCode: string
  captchaSid: string
  codeField: string
}
