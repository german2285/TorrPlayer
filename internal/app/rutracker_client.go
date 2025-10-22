package app

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// RutrackerClient handles HTTP requests to RuTracker
type RutrackerClient struct {
	httpClient *http.Client
	cookieJar  *cookiejar.Jar
	baseURL    string
	ctx        context.Context // Context for logging
}

// NewRutrackerClient creates a new RuTracker HTTP client
func NewRutrackerClient(ctx context.Context) (*RutrackerClient, error) {
	runtime.LogInfo(ctx, "Creating new RuTracker HTTP client")

	jar, err := cookiejar.New(nil)
	if err != nil {
		runtime.LogError(ctx, fmt.Sprintf("Failed to create cookie jar: %v", err))
		return nil, fmt.Errorf("failed to create cookie jar: %v", err)
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Follow up to 10 redirects
			if len(via) >= 10 {
				runtime.LogWarning(ctx, fmt.Sprintf("Stopped after 10 redirects. Last URL: %s", req.URL.String()))
				return fmt.Errorf("stopped after 10 redirects")
			}
			runtime.LogInfo(ctx, fmt.Sprintf("Redirect %d: %s -> %s", len(via), via[len(via)-1].URL.String(), req.URL.String()))
			return nil
		},
	}

	runtime.LogInfo(ctx, "RuTracker HTTP client created successfully")

	return &RutrackerClient{
		httpClient: client,
		cookieJar:  jar,
		baseURL:    "https://rutracker.net/forum",
		ctx:        ctx,
	}, nil
}

// LoadCookies loads cookies into client
func (c *RutrackerClient) LoadCookies(cookies []*http.Cookie) error {
	runtime.LogInfo(c.ctx, fmt.Sprintf("Loading %d cookies into client", len(cookies)))

	if len(cookies) == 0 {
		runtime.LogError(c.ctx, "No cookies to load")
		return fmt.Errorf("no cookies to load")
	}

	// Parse base URL
	u, err := url.Parse(c.baseURL)
	if err != nil {
		runtime.LogError(c.ctx, fmt.Sprintf("Failed to parse base URL: %v", err))
		return fmt.Errorf("failed to parse base URL: %v", err)
	}

	// Log each cookie
	for i, cookie := range cookies {
		runtime.LogInfo(c.ctx, fmt.Sprintf("Cookie %d: %s = %s (Domain: %s, Path: %s)", i+1, cookie.Name, cookie.Value[:min(20, len(cookie.Value))]+"...", cookie.Domain, cookie.Path))
	}

	// Set cookies for RuTracker domain
	c.cookieJar.SetCookies(u, cookies)

	runtime.LogInfo(c.ctx, "Cookies loaded successfully")
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
	fullURL := c.baseURL + path
	runtime.LogInfo(c.ctx, fmt.Sprintf("GET request: %s", fullURL))

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		runtime.LogError(c.ctx, fmt.Sprintf("Failed to create GET request: %v", err))
		return nil, fmt.Errorf("failed to create GET request: %v", err)
	}

	// Set headers
	c.setCommonHeaders(req)

	// Log request details
	c.logRequest(req)

	// Execute request
	runtime.LogInfo(c.ctx, "Sending GET request...")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		runtime.LogError(c.ctx, fmt.Sprintf("GET request failed: %v", err))
		return nil, fmt.Errorf("GET request failed: %v", err)
	}

	// Log response details
	c.logResponse(resp)

	return resp, nil
}

// Post performs a POST request
func (c *RutrackerClient) Post(path string, formData url.Values) (*http.Response, error) {
	fullURL := c.baseURL + path
	runtime.LogInfo(c.ctx, fmt.Sprintf("POST request: %s", fullURL))
	runtime.LogInfo(c.ctx, fmt.Sprintf("Form data: %v", formData))

	req, err := http.NewRequest("POST", fullURL, nil)
	if err != nil {
		runtime.LogError(c.ctx, fmt.Sprintf("Failed to create POST request: %v", err))
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}

	// Set form data
	req.URL.RawQuery = formData.Encode()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set common headers
	c.setCommonHeaders(req)

	// Log request details
	c.logRequest(req)

	// Execute request
	runtime.LogInfo(c.ctx, "Sending POST request...")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		runtime.LogError(c.ctx, fmt.Sprintf("POST request failed: %v", err))
		return nil, fmt.Errorf("POST request failed: %v", err)
	}

	// Log response details
	c.logResponse(resp)

	return resp, nil
}

// PostForm performs a POST request with form body
func (c *RutrackerClient) PostForm(path string, formData url.Values) (*http.Response, error) {
	fullURL := c.baseURL + path
	runtime.LogInfo(c.ctx, fmt.Sprintf("POST FORM request: %s", fullURL))
	runtime.LogInfo(c.ctx, fmt.Sprintf("Form data: %v", formData))

	req, err := http.NewRequest("POST", fullURL, nil)
	if err != nil {
		runtime.LogError(c.ctx, fmt.Sprintf("Failed to create POST request: %v", err))
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}

	// Encode form data as body
	req.URL.RawQuery = formData.Encode()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set common headers
	c.setCommonHeaders(req)

	// Log request details
	c.logRequest(req)

	// Execute request
	runtime.LogInfo(c.ctx, "Sending POST FORM request...")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		runtime.LogError(c.ctx, fmt.Sprintf("POST request failed: %v", err))
		return nil, fmt.Errorf("POST request failed: %v", err)
	}

	// Log response details
	c.logResponse(resp)

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

// logRequest logs detailed request information
func (c *RutrackerClient) logRequest(req *http.Request) {
	runtime.LogInfo(c.ctx, "=== REQUEST ===")
	runtime.LogInfo(c.ctx, fmt.Sprintf("Method: %s", req.Method))
	runtime.LogInfo(c.ctx, fmt.Sprintf("URL: %s", req.URL.String()))

	// Log headers
	runtime.LogInfo(c.ctx, "Headers:")
	for name, values := range req.Header {
		for _, value := range values {
			runtime.LogInfo(c.ctx, fmt.Sprintf("  %s: %s", name, value))
		}
	}

	// Log cookies
	cookies := req.Cookies()
	if len(cookies) > 0 {
		runtime.LogInfo(c.ctx, fmt.Sprintf("Cookies (%d):", len(cookies)))
		for _, cookie := range cookies {
			runtime.LogInfo(c.ctx, fmt.Sprintf("  %s = %s", cookie.Name, cookie.Value[:min(20, len(cookie.Value))]+"..."))
		}
	} else {
		runtime.LogInfo(c.ctx, "No cookies in request")
	}

	// Dump full request for debugging
	dump, err := httputil.DumpRequestOut(req, false)
	if err == nil {
		runtime.LogInfo(c.ctx, fmt.Sprintf("Request dump:\n%s", string(dump)))
	}

	runtime.LogInfo(c.ctx, "=== END REQUEST ===")
}

// logResponse logs detailed response information
func (c *RutrackerClient) logResponse(resp *http.Response) {
	runtime.LogInfo(c.ctx, "=== RESPONSE ===")
	runtime.LogInfo(c.ctx, fmt.Sprintf("Status: %s (%d)", resp.Status, resp.StatusCode))
	runtime.LogInfo(c.ctx, fmt.Sprintf("Proto: %s", resp.Proto))
	runtime.LogInfo(c.ctx, fmt.Sprintf("Content-Length: %d", resp.ContentLength))

	// Log headers
	runtime.LogInfo(c.ctx, "Headers:")
	for name, values := range resp.Header {
		for _, value := range values {
			runtime.LogInfo(c.ctx, fmt.Sprintf("  %s: %s", name, value))
		}
	}

	// Log cookies
	cookies := resp.Cookies()
	if len(cookies) > 0 {
		runtime.LogInfo(c.ctx, fmt.Sprintf("Set-Cookie (%d):", len(cookies)))
		for _, cookie := range cookies {
			runtime.LogInfo(c.ctx, fmt.Sprintf("  %s = %s (Domain: %s, Path: %s, Secure: %v, HttpOnly: %v)",
				cookie.Name,
				cookie.Value[:min(20, len(cookie.Value))]+"...",
				cookie.Domain,
				cookie.Path,
				cookie.Secure,
				cookie.HttpOnly))
		}
	} else {
		runtime.LogInfo(c.ctx, "No Set-Cookie headers")
	}

	// Read and log body preview (first 500 bytes)
	if resp.Body != nil {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err == nil {
			runtime.LogInfo(c.ctx, fmt.Sprintf("Body length: %d bytes", len(bodyBytes)))

			// Log preview of body
			preview := string(bodyBytes)
			if len(preview) > 500 {
				preview = preview[:500] + "... [truncated]"
			}
			runtime.LogInfo(c.ctx, fmt.Sprintf("Body preview:\n%s", preview))

			// Restore body for later reading
			resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
	}

	runtime.LogInfo(c.ctx, "=== END RESPONSE ===")
}
