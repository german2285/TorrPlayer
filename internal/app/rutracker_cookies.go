package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
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
	cookiesPath, err := a.GetCookiesPath()
	if err != nil {
		return err
	}

	// Extract RuTracker cookies
	rutrackerCookies := &RutrackerCookies{
		LastUpdated: time.Now(),
	}

	for _, cookie := range cookies {
		switch cookie.Name {
		case "bb_session":
			rutrackerCookies.BBSession = cookie.Value
		case "bb_ssl":
			rutrackerCookies.BBSSL = cookie.Value
		case "bb_t":
			rutrackerCookies.BBT = cookie.Value
		}
	}

	// Validate we have at least bb_session
	if rutrackerCookies.BBSession == "" {
		return fmt.Errorf("no valid bb_session cookie found")
	}

	// Write to file
	data, err := json.MarshalIndent(rutrackerCookies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cookies: %v", err)
	}

	err = os.WriteFile(cookiesPath, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write cookies file: %v", err)
	}

	return nil
}

// LoadCookiesFromFile loads cookies from cookies.json
func (a *App) LoadCookiesFromFile() ([]*http.Cookie, error) {
	cookiesPath, err := a.GetCookiesPath()
	if err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(cookiesPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("cookies file not found")
	}

	// Read file
	data, err := os.ReadFile(cookiesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cookies file: %v", err)
	}

	// Parse JSON
	var rutrackerCookies RutrackerCookies
	err = json.Unmarshal(data, &rutrackerCookies)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cookies: %v", err)
	}

	// Check if cookies are too old (30 days)
	if time.Since(rutrackerCookies.LastUpdated) > 30*24*time.Hour {
		return nil, fmt.Errorf("cookies expired (older than 30 days)")
	}

	// Convert to http.Cookie slice
	var cookies []*http.Cookie

	if rutrackerCookies.BBSession != "" {
		cookies = append(cookies, &http.Cookie{
			Name:  "bb_session",
			Value: rutrackerCookies.BBSession,
		})
	}

	if rutrackerCookies.BBSSL != "" {
		cookies = append(cookies, &http.Cookie{
			Name:  "bb_ssl",
			Value: rutrackerCookies.BBSSL,
		})
	}

	if rutrackerCookies.BBT != "" {
		cookies = append(cookies, &http.Cookie{
			Name:  "bb_t",
			Value: rutrackerCookies.BBT,
		})
	}

	if len(cookies) == 0 {
		return nil, fmt.Errorf("no valid cookies found in file")
	}

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
