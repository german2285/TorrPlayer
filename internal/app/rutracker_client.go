package app

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// RutrackerClient handles HTTP requests to RuTracker
type RutrackerClient struct {
	httpClient *http.Client
	cookieJar  *cookiejar.Jar
	baseURL    string
}

// NewRutrackerClient creates a new RuTracker HTTP client
func NewRutrackerClient() (*RutrackerClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %v", err)
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Follow up to 10 redirects
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		},
	}

	return &RutrackerClient{
		httpClient: client,
		cookieJar:  jar,
		baseURL:    "https://rutracker.net/forum",
	}, nil
}

// LoadCookies loads cookies into client
func (c *RutrackerClient) LoadCookies(cookies []*http.Cookie) error {
	if len(cookies) == 0 {
		return fmt.Errorf("no cookies to load")
	}

	// Parse base URL
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return fmt.Errorf("failed to parse base URL: %v", err)
	}

	// Set cookies for RuTracker domain
	c.cookieJar.SetCookies(u, cookies)

	return nil
}

// GetCookies returns current cookies from jar
func (c *RutrackerClient) GetCookies() []*http.Cookie {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil
	}

	return c.cookieJar.Cookies(u)
}

// ClearCookies clears all cookies
func (c *RutrackerClient) ClearCookies() error {
	// Create a new empty cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("failed to create cookie jar: %v", err)
	}

	c.cookieJar = jar
	c.httpClient.Jar = jar

	return nil
}

// Get performs a GET request
func (c *RutrackerClient) Get(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %v", err)
	}

	// Set headers
	c.setCommonHeaders(req)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET request failed: %v", err)
	}

	return resp, nil
}

// Post performs a POST request
func (c *RutrackerClient) Post(path string, formData url.Values) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}

	// Set form data
	req.URL.RawQuery = formData.Encode()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set common headers
	c.setCommonHeaders(req)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("POST request failed: %v", err)
	}

	return resp, nil
}

// PostForm performs a POST request with form body
func (c *RutrackerClient) PostForm(path string, formData url.Values) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}

	// Encode form data as body
	req.URL.RawQuery = formData.Encode()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set common headers
	c.setCommonHeaders(req)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("POST request failed: %v", err)
	}

	return resp, nil
}

// setCommonHeaders sets common headers for RuTracker requests
func (c *RutrackerClient) setCommonHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ru,en;q=0.9")
	req.Header.Set("Referer", c.baseURL+"/index.php")
	req.Header.Set("Connection", "keep-alive")
}
