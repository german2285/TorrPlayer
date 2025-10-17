package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/german2285/TorrPlayer/internal/player"
	torrserv "github.com/german2285/TorrPlayer/pkg/server/torr"
)

// PlayTorrentFile plays a specific file from a torrent
func (a *App) PlayTorrentFile(hash string, fileIndex int) error {
	runtime.LogInfo(a.ctx, fmt.Sprintf("Playing torrent %s file %d", hash, fileIndex))

	tor := torrserv.GetTorrent(hash)
	if tor == nil {
		return fmt.Errorf("torrent not found")
	}

	files := tor.Files()
	if fileIndex < 1 || fileIndex > len(files) {
		return fmt.Errorf("invalid file index")
	}

	// Start stream server
	port, err := a.startStreamServer(tor, fileIndex)
	if err != nil {
		return fmt.Errorf("failed to start stream server: %v", err)
	}
	a.streamPort = port

	streamURL := fmt.Sprintf("http://127.0.0.1:%s/stream", port)

	// Wait for buffer
	runtime.LogInfo(a.ctx, "Buffering...")
	a.waitForBuffer(tor, 5*time.Second)

	// Hide window and play
	runtime.WindowHide(a.ctx)
	runtime.LogInfo(a.ctx, "Starting playback...")

	err = player.PlayVideoWithMPV(streamURL)

	// Show window again
	runtime.WindowShow(a.ctx)

	// Stop stream server
	a.stopStreamServer()

	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Playback error: %v", err))
		return err
	}

	runtime.LogInfo(a.ctx, "Playback finished")
	return nil
}

// startStreamServer starts a local HTTP server for streaming
func (a *App) startStreamServer(tor *torrserv.Torrent, fileIndex int) (string, error) {
	// Create HTTP handler
	mux := http.NewServeMux()
	mux.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		tor.Stream(fileIndex, r, w)
	})

	server := &http.Server{
		Addr:    "127.0.0.1:0",
		Handler: mux,
	}

	// Create listener (OS will assign free port)
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return "", err
	}

	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)

	// Store server reference with mutex protection
	a.httpMu.Lock()
	a.httpServer = server
	a.httpMu.Unlock()

	// Start server with existing listener (no race condition)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				runtime.LogError(a.ctx, fmt.Sprintf("HTTP server panic: %v", r))
			}
		}()
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			runtime.LogError(a.ctx, fmt.Sprintf("HTTP server error: %v", err))
		}
	}()

	// Server is ready immediately since listener is already created and bound
	// Small delay to ensure goroutine starts (not strictly necessary but safer)
	time.Sleep(50 * time.Millisecond)
	return port, nil
}

// stopStreamServer stops the HTTP stream server
func (a *App) stopStreamServer() {
	a.httpMu.Lock()
	server := a.httpServer
	a.httpServer = nil
	a.httpMu.Unlock()

	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}
}

// waitForBuffer waits for minimum cache to be filled before playback
func (a *App) waitForBuffer(tor *torrserv.Torrent, maxWait time.Duration) {
	start := time.Now()
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	minCache := int64(32 << 20) // 32 MB minimum

	for {
		if time.Since(start) > maxWait {
			return
		}

		cache := tor.GetCache()
		if cache != nil {
			state := cache.GetState()
			if state.Filled >= minCache {
				return
			}
		}

		<-ticker.C
	}
}
