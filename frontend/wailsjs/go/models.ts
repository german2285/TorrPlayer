export namespace main {
	
	export class Settings {
	    cacheSize: number;
	    cacheSizeStr: string;
	    connectionsLimit: number;
	    downloadRate: number;
	    uploadRate: number;
	    preloadCache: number;
	    retrackersMode: number;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cacheSize = source["cacheSize"];
	        this.cacheSizeStr = source["cacheSizeStr"];
	        this.connectionsLimit = source["connectionsLimit"];
	        this.downloadRate = source["downloadRate"];
	        this.uploadRate = source["uploadRate"];
	        this.preloadCache = source["preloadCache"];
	        this.retrackersMode = source["retrackersMode"];
	    }
	}
	export class Torrent {
	    hash: string;
	    name: string;
	    title: string;
	    size: number;
	    sizeStr: string;
	    status: string;
	    progress: number;
	    downSpeed: number;
	    upSpeed: number;
	    downSpeedStr: string;
	    upSpeedStr: string;
	    peers: number;
	    seeders: number;
	    fileCount: number;
	    category: string;
	    poster: string;
	    timestamp: number;
	    loadingMeta: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Torrent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hash = source["hash"];
	        this.name = source["name"];
	        this.title = source["title"];
	        this.size = source["size"];
	        this.sizeStr = source["sizeStr"];
	        this.status = source["status"];
	        this.progress = source["progress"];
	        this.downSpeed = source["downSpeed"];
	        this.upSpeed = source["upSpeed"];
	        this.downSpeedStr = source["downSpeedStr"];
	        this.upSpeedStr = source["upSpeedStr"];
	        this.peers = source["peers"];
	        this.seeders = source["seeders"];
	        this.fileCount = source["fileCount"];
	        this.category = source["category"];
	        this.poster = source["poster"];
	        this.timestamp = source["timestamp"];
	        this.loadingMeta = source["loadingMeta"];
	    }
	}
	export class TorrentFile {
	    index: number;
	    path: string;
	    size: number;
	    sizeStr: string;
	
	    static createFrom(source: any = {}) {
	        return new TorrentFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.path = source["path"];
	        this.size = source["size"];
	        this.sizeStr = source["sizeStr"];
	    }
	}
	export class TorrentStats {
	    downSpeed: number;
	    upSpeed: number;
	    downSpeedStr: string;
	    upSpeedStr: string;
	    peers: number;
	    seeders: number;
	    downloaded: number;
	    downloadedStr: string;
	    cacheFilled: number;
	    cacheCapacity: number;
	    cacheFilledStr: string;
	    cacheCapacityStr: string;
	
	    static createFrom(source: any = {}) {
	        return new TorrentStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.downSpeed = source["downSpeed"];
	        this.upSpeed = source["upSpeed"];
	        this.downSpeedStr = source["downSpeedStr"];
	        this.upSpeedStr = source["upSpeedStr"];
	        this.peers = source["peers"];
	        this.seeders = source["seeders"];
	        this.downloaded = source["downloaded"];
	        this.downloadedStr = source["downloadedStr"];
	        this.cacheFilled = source["cacheFilled"];
	        this.cacheCapacity = source["cacheCapacity"];
	        this.cacheFilledStr = source["cacheFilledStr"];
	        this.cacheCapacityStr = source["cacheCapacityStr"];
	    }
	}

}

