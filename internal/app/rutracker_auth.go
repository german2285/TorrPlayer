package app

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// CaptchaData represents CAPTCHA information for registration
type CaptchaData struct {
	ImageBase64 string `json:"imageBase64"`
	SID         string `json:"sid"`
	CodeField   string `json:"codeField"`
}

// LoginData represents login credentials
type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegistrationData represents registration form data
type RegistrationData struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	CaptchaCode string `json:"captchaCode"`
	CaptchaSID  string `json:"captchaSid"`
	CodeField   string `json:"codeField"`
}

// GetRegistrationCaptcha fetches CAPTCHA image and data for registration
func (a *App) GetRegistrationCaptcha() (*CaptchaData, error) {
	runtime.LogInfo(a.ctx, "Fetching registration CAPTCHA")

	// Create HTTP client
	client, err := NewRutrackerClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %v", err)
	}

	// GET profile.php?mode=register
	resp, err := client.Get("/profile.php?mode=register")
	if err != nil {
		return nil, fmt.Errorf("failed to load registration page: %v", err)
	}
	defer resp.Body.Close()

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Extract CAPTCHA image URL
	captchaImg := doc.Find("img[src*='captcha']").First()
	if captchaImg.Length() == 0 {
		return nil, fmt.Errorf("CAPTCHA image not found")
	}

	captchaURL, exists := captchaImg.Attr("src")
	if !exists {
		return nil, fmt.Errorf("CAPTCHA URL not found")
	}

	// Extract cap_sid
	capSIDInput := doc.Find("input[name='cap_sid']").First()
	if capSIDInput.Length() == 0 {
		return nil, fmt.Errorf("cap_sid not found")
	}

	capSID, exists := capSIDInput.Attr("value")
	if !exists {
		return nil, fmt.Errorf("cap_sid value not found")
	}

	// Extract cap_code field name (dynamic)
	capCodeInput := doc.Find("input[name^='cap_code']").First()
	if capCodeInput.Length() == 0 {
		return nil, fmt.Errorf("cap_code field not found")
	}

	capCodeField, exists := capCodeInput.Attr("name")
	if !exists {
		return nil, fmt.Errorf("cap_code field name not found")
	}

	// Build full CAPTCHA image URL
	if !strings.HasPrefix(captchaURL, "http") {
		captchaURL = "https://static.rutracker.cc/" + strings.TrimPrefix(captchaURL, "/")
	}

	runtime.LogInfo(a.ctx, fmt.Sprintf("CAPTCHA URL: %s", captchaURL))

	// Download CAPTCHA image
	imgResp, err := client.Get(strings.TrimPrefix(captchaURL, client.baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to download CAPTCHA image: %v", err)
	}
	defer imgResp.Body.Close()

	// Read image bytes
	imageBytes, err := io.ReadAll(imgResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read CAPTCHA image: %v", err)
	}

	// Convert to base64
	base64Image := base64.StdEncoding.EncodeToString(imageBytes)
	imageDataURL := "data:image/jpeg;base64," + base64Image

	runtime.LogInfo(a.ctx, "CAPTCHA fetched successfully")

	return &CaptchaData{
		ImageBase64: imageDataURL,
		SID:         capSID,
		CodeField:   capCodeField,
	}, nil
}

// RegisterOnRuTracker registers a new account on RuTracker
func (a *App) RegisterOnRuTracker(data *RegistrationData) error {
	runtime.LogInfo(a.ctx, fmt.Sprintf("Registering user: %s", data.Username))

	// Validate input
	if data.Username == "" || data.Password == "" || data.Email == "" {
		return fmt.Errorf("username, password, and email are required")
	}

	if len(data.Password) > 20 {
		return fmt.Errorf("password must be max 20 characters")
	}

	if data.CaptchaCode == "" || data.CaptchaSID == "" || data.CodeField == "" {
		return fmt.Errorf("CAPTCHA data is required")
	}

	// Create HTTP client
	client, err := NewRutrackerClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %v", err)
	}

	// Prepare form data
	formData := url.Values{}
	formData.Set("username", data.Username)
	formData.Set("new_pass", data.Password)
	formData.Set("user_email", data.Email)
	formData.Set("cap_sid", data.CaptchaSID)
	formData.Set(data.CodeField, data.CaptchaCode)
	formData.Set("reg_agreed", "1")
	formData.Set("user_timezone_x2", "6")  // GMT+3 Moscow time
	formData.Set("user_gender_id", "0")     // Gender: Hidden

	// POST registration
	resp, err := client.PostForm("/profile.php?mode=register", formData)
	if err != nil {
		return fmt.Errorf("registration request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	// Check for error messages
	errorMsg := doc.Find(".has-validation-error").First().Text()
	if errorMsg != "" {
		return fmt.Errorf("registration failed: %s", strings.TrimSpace(errorMsg))
	}

	// Check for success (cookies should be set)
	cookies := client.GetCookies()
	if len(cookies) == 0 {
		return fmt.Errorf("registration failed: no cookies received")
	}

	// Save cookies
	err = a.SaveCookiesToFile(cookies)
	if err != nil {
		return fmt.Errorf("failed to save cookies: %v", err)
	}

	runtime.LogInfo(a.ctx, "Registration successful")
	return nil
}

// LoginToRuTracker logs in with username and password
func (a *App) LoginToRuTracker(data *LoginData) error {
	runtime.LogInfo(a.ctx, fmt.Sprintf("Logging in user: %s", data.Username))

	// Validate input
	if data.Username == "" || data.Password == "" {
		return fmt.Errorf("username and password are required")
	}

	// Create HTTP client
	client, err := NewRutrackerClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %v", err)
	}

	// Prepare form data
	formData := url.Values{}
	formData.Set("login_username", data.Username)
	formData.Set("login_password", data.Password)
	formData.Set("login", "вход")

	// POST login
	resp, err := client.PostForm("/login.php", formData)
	if err != nil {
		return fmt.Errorf("login request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check if redirected to login page (failed auth)
	if strings.Contains(resp.Request.URL.Path, "login.php") {
		return fmt.Errorf("login failed: invalid credentials")
	}

	// Get cookies
	cookies := client.GetCookies()
	if len(cookies) == 0 {
		return fmt.Errorf("login failed: no cookies received")
	}

	// Check if bb_session exists
	hasSession := false
	for _, cookie := range cookies {
		if cookie.Name == "bb_session" {
			hasSession = true
			break
		}
	}

	if !hasSession {
		return fmt.Errorf("login failed: no session cookie received")
	}

	// Save cookies
	err = a.SaveCookiesToFile(cookies)
	if err != nil {
		return fmt.Errorf("failed to save cookies: %v", err)
	}

	runtime.LogInfo(a.ctx, "Login successful")
	return nil
}

// CheckAuthStatus checks if user is authenticated
func (a *App) CheckAuthStatus() (bool, error) {
	runtime.LogInfo(a.ctx, "Checking auth status")

	// Try to load cookies
	cookies, err := a.LoadCookiesFromFile()
	if err != nil {
		runtime.LogInfo(a.ctx, "No valid cookies found")
		return false, nil
	}

	// Create HTTP client and load cookies
	client, err := NewRutrackerClient()
	if err != nil {
		return false, fmt.Errorf("failed to create HTTP client: %v", err)
	}

	err = client.LoadCookies(cookies)
	if err != nil {
		return false, fmt.Errorf("failed to load cookies: %v", err)
	}

	// Try to access tracker.php (requires auth)
	resp, err := client.Get("/tracker.php")
	if err != nil {
		return false, fmt.Errorf("failed to check auth: %v", err)
	}
	defer resp.Body.Close()

	// If redirected to login.php, not authenticated
	if strings.Contains(resp.Request.URL.Path, "login.php") {
		runtime.LogInfo(a.ctx, "Not authenticated (redirected to login)")
		return false, nil
	}

	runtime.LogInfo(a.ctx, "Authenticated successfully")
	return true, nil
}

// LogoutFromRuTracker logs out and deletes cookies
func (a *App) LogoutFromRuTracker() error {
	runtime.LogInfo(a.ctx, "Logging out")

	// Delete cookies file
	err := a.DeleteCookiesFile()
	if err != nil {
		return fmt.Errorf("failed to delete cookies: %v", err)
	}

	runtime.LogInfo(a.ctx, "Logged out successfully")
	return nil
}
