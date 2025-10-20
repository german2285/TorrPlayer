package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/german2285/TorrPlayer/pkg/server/log"
	"github.com/german2285/TorrPlayer/pkg/server/settings"
	torrserv "github.com/german2285/TorrPlayer/pkg/server/torr"
)

// App struct
type App struct {
	ctx        context.Context
	btServer   *torrserv.BTServer
	httpServer *http.Server
	httpMu     sync.Mutex
	streamPort string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	runtime.LogInfo(ctx, "TorrPlayer starting...")

	// Initialize settings and logging
	// Use executable directory for config storage (settings.json will be next to .exe)
	exePath, err := os.Executable()
	if err != nil {
		runtime.LogError(ctx, fmt.Sprintf("Failed to get executable path: %v", err))
		exePath, _ = os.Getwd() // Fallback to current directory
	}

	// Get directory where executable is located
	configDir := filepath.Dir(exePath)

	settings.Path = configDir
	runtime.LogInfo(ctx, fmt.Sprintf("Config directory: %s", configDir))
	settings.InitSets(false, false) // readOnly=false, searchWA=false
	log.Init("", "")

	// Initialize DNS resolver
	a.initDNSResolver()

	// Initialize BitTorrent server
	runtime.LogInfo(ctx, "Initializing BitTorrent client...")
	a.btServer = torrserv.NewBTS()
	err = a.btServer.Connect()
	if err != nil {
		runtime.LogError(ctx, fmt.Sprintf("Failed to initialize BT client: %v", err))
		return
	}
	runtime.LogInfo(ctx, "BitTorrent client initialized successfully")
}

// Shutdown is called when the app is closing
func (a *App) Shutdown(ctx context.Context) {
	if a.btServer != nil {
		a.btServer.Disconnect()
	}
	if a.httpServer != nil {
		a.stopStreamServer()
	}
	settings.CloseDB()
}

func (a *App) initDNSResolver() {
	addrs, err := net.LookupHost("www.google.com")
	if len(addrs) == 0 {
		fn := func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", "1.1.1.1:53")
		}
		net.DefaultResolver = &net.Resolver{
			Dial: fn,
		}
	}
	_ = err
}
