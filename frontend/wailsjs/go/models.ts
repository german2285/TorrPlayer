export namespace app {
	
	export class CaptchaData {
	    imageBase64: string;
	    sid: string;
	    codeField: string;
	
	    static createFrom(source: any = {}) {
	        return new CaptchaData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.imageBase64 = source["imageBase64"];
	        this.sid = source["sid"];
	        this.codeField = source["codeField"];
	    }
	}
	export class LoginData {
	    username: string;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new LoginData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.password = source["password"];
	    }
	}
	export class RegistrationData {
	    username: string;
	    password: string;
	    email: string;
	    captchaCode: string;
	    captchaSid: string;
	    codeField: string;
	
	    static createFrom(source: any = {}) {
	        return new RegistrationData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.password = source["password"];
	        this.email = source["email"];
	        this.captchaCode = source["captchaCode"];
	        this.captchaSid = source["captchaSid"];
	        this.codeField = source["codeField"];
	    }
	}
	export class RutrackerTorrent {
	    topicId: string;
	    title: string;
	    category: string;
	    size: string;
	    seeds: number;
	    leeches: number;
	    author: string;
	    date: string;
	
	    static createFrom(source: any = {}) {
	        return new RutrackerTorrent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.topicId = source["topicId"];
	        this.title = source["title"];
	        this.category = source["category"];
	        this.size = source["size"];
	        this.seeds = source["seeds"];
	        this.leeches = source["leeches"];
	        this.author = source["author"];
	        this.date = source["date"];
	    }
	}
	export class Settings {
	    cacheSize: number;
	    cacheSizeStr: string;
	    connectionsLimit: number;
	    downloadRate: number;
	    uploadRate: number;
	    preloadCache: number;
	    retrackersMode: number;
	    themeColor: string;
	    bgMusicVolume: number;
	
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
	        this.themeColor = source["themeColor"];
	        this.bgMusicVolume = source["bgMusicVolume"];
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

export namespace http {
	
	export class Cookie {
	    Name: string;
	    Value: string;
	    Quoted: boolean;
	    Path: string;
	    Domain: string;
	    // Go type: time
	    Expires: any;
	    RawExpires: string;
	    MaxAge: number;
	    Secure: boolean;
	    HttpOnly: boolean;
	    SameSite: number;
	    Partitioned: boolean;
	    Raw: string;
	    Unparsed: string[];
	
	    static createFrom(source: any = {}) {
	        return new Cookie(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Value = source["Value"];
	        this.Quoted = source["Quoted"];
	        this.Path = source["Path"];
	        this.Domain = source["Domain"];
	        this.Expires = this.convertValues(source["Expires"], null);
	        this.RawExpires = source["RawExpires"];
	        this.MaxAge = source["MaxAge"];
	        this.Secure = source["Secure"];
	        this.HttpOnly = source["HttpOnly"];
	        this.SameSite = source["SameSite"];
	        this.Partitioned = source["Partitioned"];
	        this.Raw = source["Raw"];
	        this.Unparsed = source["Unparsed"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

