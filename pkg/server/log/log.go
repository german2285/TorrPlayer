package log

import (
	"log"
	"os"
)

var (
	logPath    = ""
	webLogPath = ""
)

var webLog *log.Logger

var (
	logFile    *os.File
	webLogFile *os.File
)

func Init(path, webpath string) {
	webLogPath = webpath
	logPath = path

	if webpath != "" {
		ff, err := os.OpenFile(webLogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
		if err != nil {
			TLogln("Error create web log file:", err)
		} else {
			webLogFile = ff
			webLog = log.New(ff, " ", log.LstdFlags)
		}
	}

	if path != "" {
		if fi, err := os.Lstat(path); err == nil {
			if fi.Size() >= 100*1024*1024 { // 100MB
				os.Remove(path)
			}
		}
		ff, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
		if err != nil {
			TLogln("Error create log file:", err)
			return
		}
		logFile = ff
		os.Stdout = ff
		os.Stderr = ff
		// var timeFmt string
		// var ok bool
		// timeFmt, ok = os.LookupEnv("GO_LOG_TIME_FMT")
		// if !ok {
		// 	timeFmt = "2006-01-02T15:04:05-0700"
		// }
		// log.SetFlags(log.Lmsgprefix)
		// log.SetPrefix(time.Now().Format(timeFmt) + " TSM ")
		log.SetFlags(log.LstdFlags | log.LUTC | log.Lmsgprefix)
		log.SetPrefix("UTC0 ")
		log.SetOutput(ff)
	}
}

func Close() {
	if logFile != nil {
		logFile.Close()
	}
	if webLogFile != nil {
		webLogFile.Close()
	}
}

func TLogln(v ...interface{}) {
	log.Println(v...)
}

func WebLogln(v ...interface{}) {
	if webLog != nil {
		webLog.Println(v...)
	}
}
