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
	runtime.LogInfo(a.ctx, "========== SearchRuTracker START ==========")
	runtime.LogInfo(a.ctx, fmt.Sprintf("Search query: '%s'", query))

	if query == "" {
		runtime.LogError(a.ctx, "Search query is empty")
		return nil, fmt.Errorf("search query is empty")
	}

	// Load cookies
	runtime.LogInfo(a.ctx, "Loading cookies for authentication...")
	cookies, err := a.LoadCookiesFromFile()
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Not authenticated: %v", err))
		return nil, fmt.Errorf("not authenticated: %v", err)
	}
	runtime.LogInfo(a.ctx, fmt.Sprintf("Loaded %d cookies from file", len(cookies)))

	// Create HTTP client and load cookies
	runtime.LogInfo(a.ctx, "Creating HTTP client...")
	client, err := NewRutrackerClient(a.ctx)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to create HTTP client: %v", err))
		return nil, fmt.Errorf("failed to create HTTP client: %v", err)
	}

	runtime.LogInfo(a.ctx, "Loading cookies into HTTP client...")
	err = client.LoadCookies(cookies)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to load cookies: %v", err))
		return nil, fmt.Errorf("failed to load cookies: %v", err)
	}

	// Encode query to Windows-1251
	runtime.LogInfo(a.ctx, "Encoding query to Windows-1251...")
	queryWin1251, err := charmap.Windows1251.NewEncoder().String(query)
	if err != nil {
		runtime.LogWarning(a.ctx, fmt.Sprintf("Failed to encode to Windows-1251, using UTF-8: %v", err))
		queryWin1251 = query
	}
	runtime.LogInfo(a.ctx, fmt.Sprintf("Encoded query: '%s'", queryWin1251))

	// Prepare form data
	runtime.LogInfo(a.ctx, "Preparing form data...")
	formData := url.Values{}
	formData.Set("max", "1")
	formData.Set("nm", queryWin1251)
	runtime.LogInfo(a.ctx, fmt.Sprintf("Form data: %v", formData))

	// POST tracker.php?nm={query}
	searchURL := fmt.Sprintf("/tracker.php?nm=%s", url.QueryEscape(query))
	runtime.LogInfo(a.ctx, fmt.Sprintf("Search URL: %s", searchURL))
	runtime.LogInfo(a.ctx, "Sending search request...")
	resp, err := client.PostForm(searchURL, formData)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Search request failed: %v", err))
		return nil, fmt.Errorf("search request failed: %v", err)
	}
	defer resp.Body.Close()

	runtime.LogInfo(a.ctx, fmt.Sprintf("Response status: %s", resp.Status))
	runtime.LogInfo(a.ctx, fmt.Sprintf("Final URL after redirects: %s", resp.Request.URL.String()))

	// Check if redirected to login (not authenticated)
	if strings.Contains(resp.Request.URL.Path, "login.php") {
		runtime.LogError(a.ctx, "Redirected to login page - not authenticated")
		return nil, fmt.Errorf("not authenticated: redirected to login page")
	}

	// Parse HTML
	runtime.LogInfo(a.ctx, "Parsing HTML response...")
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to parse HTML: %v", err))
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}
	runtime.LogInfo(a.ctx, "HTML parsed successfully")

	// Find all torrent rows
	var torrents []RutrackerTorrent
	runtime.LogInfo(a.ctx, "Searching for torrent rows (tr.hl-tr)...")

	doc.Find("tr.hl-tr").Each(func(i int, s *goquery.Selection) {
		runtime.LogInfo(a.ctx, fmt.Sprintf("Processing torrent row %d...", i+1))
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

	runtime.LogInfo(a.ctx, fmt.Sprintf("Search completed: found %d torrents", len(torrents)))

	// Log details of found torrents
	for i, t := range torrents {
		runtime.LogInfo(a.ctx, fmt.Sprintf("Torrent %d: ID=%s, Title='%s', Seeds=%d, Leeches=%d, Size=%s",
			i+1, t.TopicID, t.Title, t.Seeds, t.Leeches, t.Size))
	}

	runtime.LogInfo(a.ctx, "========== SearchRuTracker END ==========")
	return torrents, nil
}

// GetRutrackerMagnetLink gets magnet link for a torrent by topic ID
func (a *App) GetRutrackerMagnetLink(topicID string) (string, error) {
	runtime.LogInfo(a.ctx, "========== GetRutrackerMagnetLink START ==========")
	runtime.LogInfo(a.ctx, fmt.Sprintf("Topic ID: %s", topicID))

	if topicID == "" {
		runtime.LogError(a.ctx, "Topic ID is empty")
		return "", fmt.Errorf("topic ID is empty")
	}

	// Load cookies
	runtime.LogInfo(a.ctx, "Loading cookies for authentication...")
	cookies, err := a.LoadCookiesFromFile()
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Not authenticated: %v", err))
		return "", fmt.Errorf("not authenticated: %v", err)
	}
	runtime.LogInfo(a.ctx, fmt.Sprintf("Loaded %d cookies from file", len(cookies)))

	// Create HTTP client and load cookies
	runtime.LogInfo(a.ctx, "Creating HTTP client...")
	client, err := NewRutrackerClient(a.ctx)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to create HTTP client: %v", err))
		return "", fmt.Errorf("failed to create HTTP client: %v", err)
	}

	runtime.LogInfo(a.ctx, "Loading cookies into HTTP client...")
	err = client.LoadCookies(cookies)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to load cookies: %v", err))
		return "", fmt.Errorf("failed to load cookies: %v", err)
	}

	// GET viewtopic.php?t={topicID}
	topicURL := fmt.Sprintf("/viewtopic.php?t=%s", topicID)
	runtime.LogInfo(a.ctx, fmt.Sprintf("Loading topic page: %s", topicURL))
	resp, err := client.Get(topicURL)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to load topic page: %v", err))
		return "", fmt.Errorf("failed to load topic page: %v", err)
	}
	defer resp.Body.Close()

	runtime.LogInfo(a.ctx, fmt.Sprintf("Response status: %s", resp.Status))
	runtime.LogInfo(a.ctx, fmt.Sprintf("Final URL after redirects: %s", resp.Request.URL.String()))

	// Check if redirected to login
	if strings.Contains(resp.Request.URL.Path, "login.php") {
		runtime.LogError(a.ctx, "Redirected to login page - not authenticated")
		return "", fmt.Errorf("not authenticated: redirected to login page")
	}

	// Parse HTML
	runtime.LogInfo(a.ctx, "Parsing HTML response...")
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Failed to parse HTML: %v", err))
		return "", fmt.Errorf("failed to parse HTML: %v", err)
	}
	runtime.LogInfo(a.ctx, "HTML parsed successfully")

	// Find magnet link
	runtime.LogInfo(a.ctx, "Searching for magnet link (a.magnet-link)...")
	magnetLink, exists := doc.Find("a.magnet-link").Attr("href")
	if !exists || magnetLink == "" {
		runtime.LogError(a.ctx, "Magnet link not found in HTML")
		// Log HTML snippet for debugging
		html, _ := doc.Html()
		if len(html) > 1000 {
			runtime.LogInfo(a.ctx, fmt.Sprintf("HTML snippet:\n%s", html[:1000]+"... [truncated]"))
		} else {
			runtime.LogInfo(a.ctx, fmt.Sprintf("HTML:\n%s", html))
		}
		return "", fmt.Errorf("magnet link not found")
	}

	runtime.LogInfo(a.ctx, fmt.Sprintf("Magnet link found: %s", magnetLink[:50]+"..."))
	runtime.LogInfo(a.ctx, "========== GetRutrackerMagnetLink END ==========")

	return magnetLink, nil
}
