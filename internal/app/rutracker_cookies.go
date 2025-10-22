package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// RutrackerCookies represents stored cookies for RuTracker
type RutrackerCookies struct {
	BBSession   string    `json:"bb_session"`
	BBSSL       string    `json:"bb_ssl"`
	BBT         string    `json:"bb_t"`
	LastUpdated time.Time `json:"lastUpdated"`
}

// GetCookiesPath returns path to cookies.json next to executable
func (a *App) GetCookiesPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %v", err)
	}

	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, "cookies.json"), nil
}

// SaveCookiesToFile saves HTTP cookies to cookies.json
func (a *App) SaveCookiesToFile(cookies []*http.Cookie) error {
	runtime.LogInfo(a.ctx, "========== SaveCookiesToFile START ==========")
	runtime.LogInfo(a.ctx, fmt.Sprintf("Saving %d cookies to file", len(cookies)))

	cookiesPath, err := a.GetCookiesPath()
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to get cookies path: %v", err))
		return err
	}
	runtime.LogInfo(a.ctx, fmt.Sprintf("Cookies path: %s", cookiesPath))

	// Extract RuTracker cookies
	rutrackerCookies := &RutrackerCookies{
		LastUpdated: time.Now(),
	}

	runtime.LogInfo(a.ctx, "Extracting RuTracker cookies...")
	for _, cookie := range cookies {
		runtime.LogInfo(a.ctx, fmt.Sprintf("Cookie: %s = %s (Domain: %s)", cookie.Name, cookie.Value[:min(20, len(cookie.Value))]+"...", cookie.Domain))
		switch cookie.Name {
		case "bb_session":
			rutrackerCookies.BBSession = cookie.Value
			runtime.LogInfo(a.ctx, "Found bb_session cookie")
		case "bb_ssl":
			rutrackerCookies.BBSSL = cookie.Value
			runtime.LogInfo(a.ctx, "Found bb_ssl cookie")
		case "bb_t":
			rutrackerCookies.BBT = cookie.Value
			runtime.LogInfo(a.ctx, "Found bb_t cookie")
		}
	}

	// Validate we have at least bb_session
	if rutrackerCookies.BBSession == "" {
		runtime.LogError(a.ctx, "No valid bb_session cookie found")
		return fmt.Errorf("no valid bb_session cookie found")
	}

	// Write to file
	runtime.LogInfo(a.ctx, "Marshaling cookies to JSON...")
	data, err := json.MarshalIndent(rutrackerCookies, "", "  ")
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to marshal cookies: %v", err))
		return fmt.Errorf("failed to marshal cookies: %v", err)
	}

	runtime.LogInfo(a.ctx, fmt.Sprintf("Writing %d bytes to file...", len(data)))
	err = os.WriteFile(cookiesPath, data, 0600)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to write cookies file: %v", err))
		return fmt.Errorf("failed to write cookies file: %v", err)
	}

	runtime.LogInfo(a.ctx, "Cookies saved successfully")
	runtime.LogInfo(a.ctx, "========== SaveCookiesToFile END ==========")
	return nil
}

// LoadCookiesFromFile loads cookies from cookies.json
func (a *App) LoadCookiesFromFile() ([]*http.Cookie, error) {
	runtime.LogInfo(a.ctx, "========== LoadCookiesFromFile START ==========")

	cookiesPath, err := a.GetCookiesPath()
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to get cookies path: %v", err))
		return nil, err
	}
	runtime.LogInfo(a.ctx, fmt.Sprintf("Cookies path: %s", cookiesPath))

	// Check if file exists
	runtime.LogInfo(a.ctx, "Checking if cookies file exists...")
	if _, err := os.Stat(cookiesPath); os.IsNotExist(err) {
		runtime.LogError(a.ctx, "Cookies file not found")
		return nil, fmt.Errorf("cookies file not found")
	}
	runtime.LogInfo(a.ctx, "Cookies file exists")

	// Read file
	runtime.LogInfo(a.ctx, "Reading cookies file...")
	data, err := os.ReadFile(cookiesPath)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to read cookies file: %v", err))
		return nil, fmt.Errorf("failed to read cookies file: %v", err)
	}
	runtime.LogInfo(a.ctx, fmt.Sprintf("Read %d bytes from cookies file", len(data)))

	// Parse JSON
	runtime.LogInfo(a.ctx, "Parsing JSON...")
	var rutrackerCookies RutrackerCookies
	err = json.Unmarshal(data, &rutrackerCookies)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to unmarshal cookies: %v", err))
		return nil, fmt.Errorf("failed to unmarshal cookies: %v", err)
	}
	runtime.LogInfo(a.ctx, fmt.Sprintf("Cookies last updated: %s", rutrackerCookies.LastUpdated))

	// Check if cookies are too old (30 days)
	age := time.Since(rutrackerCookies.LastUpdated)
	runtime.LogInfo(a.ctx, fmt.Sprintf("Cookies age: %v", age))
	if age > 30*24*time.Hour {
		runtime.LogError(a.ctx, fmt.Sprintf("Cookies expired (age: %v, max: 30 days)", age))
		return nil, fmt.Errorf("cookies expired (older than 30 days)")
	}

	// Convert to http.Cookie slice
	runtime.LogInfo(a.ctx, "Converting to http.Cookie slice...")
	var cookies []*http.Cookie

	if rutrackerCookies.BBSession != "" {
		cookies = append(cookies, &http.Cookie{
			Name:  "bb_session",
			Value: rutrackerCookies.BBSession,
		})
		runtime.LogInfo(a.ctx, "Added bb_session cookie")
	}

	if rutrackerCookies.BBSSL != "" {
		cookies = append(cookies, &http.Cookie{
			Name:  "bb_ssl",
			Value: rutrackerCookies.BBSSL,
		})
		runtime.LogInfo(a.ctx, "Added bb_ssl cookie")
	}

	if rutrackerCookies.BBT != "" {
		cookies = append(cookies, &http.Cookie{
			Name:  "bb_t",
			Value: rutrackerCookies.BBT,
		})
		runtime.LogInfo(a.ctx, "Added bb_t cookie")
	}

	if len(cookies) == 0 {
		runtime.LogError(a.ctx, "No valid cookies found in file")
		return nil, fmt.Errorf("no valid cookies found in file")
	}

	runtime.LogInfo(a.ctx, fmt.Sprintf("Loaded %d cookies successfully", len(cookies)))
	runtime.LogInfo(a.ctx, "========== LoadCookiesFromFile END ==========")
	return cookies, nil
}

// DeleteCookiesFile removes cookies.json
func (a *App) DeleteCookiesFile() error {
	cookiesPath, err := a.GetCookiesPath()
	if err != nil {
		return err
	}

	err = os.Remove(cookiesPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete cookies file: %v", err)
	}

	return nil
}
