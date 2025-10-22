package app

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/text/encoding/charmap"
)

// RutrackerTorrent represents a torrent from RuTracker search results
type RutrackerTorrent struct {
	TopicID  string `json:"topicId"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Size     string `json:"size"`
	Seeds    int    `json:"seeds"`
	Leeches  int    `json:"leeches"`
	Author   string `json:"author"`
	Date     string `json:"date"`
}

// SearchRuTracker searches for torrents on RuTracker
func (a *App) SearchRuTracker(query string) ([]RutrackerTorrent, error) {
	runtime.LogInfo(a.ctx, fmt.Sprintf("Searching RuTracker for: %s", query))

	if query == "" {
		return nil, fmt.Errorf("search query is empty")
	}

	// Load cookies
	cookies, err := a.LoadCookiesFromFile()
	if err != nil {
		return nil, fmt.Errorf("not authenticated: %v", err)
	}

	// Create HTTP client and load cookies
	client, err := NewRutrackerClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %v", err)
	}

	err = client.LoadCookies(cookies)
	if err != nil {
		return nil, fmt.Errorf("failed to load cookies: %v", err)
	}

	// Encode query to Windows-1251
	queryWin1251, err := charmap.Windows1251.NewEncoder().String(query)
	if err != nil {
		runtime.LogWarning(a.ctx, fmt.Sprintf("Failed to encode to Windows-1251, using UTF-8: %v", err))
		queryWin1251 = query
	}

	// Prepare form data
	formData := url.Values{}
	formData.Set("max", "1")
	formData.Set("nm", queryWin1251)

	// POST tracker.php?nm={query}
	searchURL := fmt.Sprintf("/tracker.php?nm=%s", url.QueryEscape(query))
	resp, err := client.PostForm(searchURL, formData)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check if redirected to login (not authenticated)
	if strings.Contains(resp.Request.URL.Path, "login.php") {
		return nil, fmt.Errorf("not authenticated: redirected to login page")
	}

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Find all torrent rows
	var torrents []RutrackerTorrent

	doc.Find("tr.hl-tr").Each(func(i int, s *goquery.Selection) {
		// Extract topic ID
		topicID, exists := s.Attr("data-topic_id")
		if !exists {
			return
		}

		// Extract title
		titleEl := s.Find(".t-title-col .tLink").First()
		title := strings.TrimSpace(titleEl.Text())
		if title == "" {
			return
		}

		// Extract category
		categoryEl := s.Find(".f-name-col a").First()
		category := strings.TrimSpace(categoryEl.Text())

		// Extract size
		sizeEl := s.Find(".tor-size a").First()
		size := strings.TrimSpace(sizeEl.Text())
		// Remove "↓" symbol
		size = strings.TrimSuffix(size, " ↓")

		// Extract seeds
		seedsEl := s.Find(".seedmed b").First()
		seedsStr := strings.TrimSpace(seedsEl.Text())
		seeds, _ := strconv.Atoi(seedsStr)

		// Extract leeches
		leechesEl := s.Find(".leechmed").First()
		leechesStr := strings.TrimSpace(leechesEl.Text())
		leeches, _ := strconv.Atoi(leechesStr)

		// Extract author
		authorEl := s.Find(".u-name-col a").First()
		author := strings.TrimSpace(authorEl.Text())

		// Extract date
		dateEl := s.Find("td.small.nowrap p").First()
		date := strings.TrimSpace(dateEl.Text())

		torrents = append(torrents, RutrackerTorrent{
			TopicID:  topicID,
			Title:    title,
			Category: category,
			Size:     size,
			Seeds:    seeds,
			Leeches:  leeches,
			Author:   author,
			Date:     date,
		})
	})

	runtime.LogInfo(a.ctx, fmt.Sprintf("Found %d torrents", len(torrents)))

	return torrents, nil
}

// GetRutrackerMagnetLink gets magnet link for a torrent by topic ID
func (a *App) GetRutrackerMagnetLink(topicID string) (string, error) {
	runtime.LogInfo(a.ctx, fmt.Sprintf("Getting magnet link for topic: %s", topicID))

	if topicID == "" {
		return "", fmt.Errorf("topic ID is empty")
	}

	// Load cookies
	cookies, err := a.LoadCookiesFromFile()
	if err != nil {
		return "", fmt.Errorf("not authenticated: %v", err)
	}

	// Create HTTP client and load cookies
	client, err := NewRutrackerClient()
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP client: %v", err)
	}

	err = client.LoadCookies(cookies)
	if err != nil {
		return "", fmt.Errorf("failed to load cookies: %v", err)
	}

	// GET viewtopic.php?t={topicID}
	resp, err := client.Get(fmt.Sprintf("/viewtopic.php?t=%s", topicID))
	if err != nil {
		return "", fmt.Errorf("failed to load topic page: %v", err)
	}
	defer resp.Body.Close()

	// Check if redirected to login
	if strings.Contains(resp.Request.URL.Path, "login.php") {
		return "", fmt.Errorf("not authenticated: redirected to login page")
	}

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Find magnet link
	magnetLink, exists := doc.Find("a.magnet-link").Attr("href")
	if !exists || magnetLink == "" {
		return "", fmt.Errorf("magnet link not found")
	}

	runtime.LogInfo(a.ctx, "Magnet link found successfully")

	return magnetLink, nil
}
