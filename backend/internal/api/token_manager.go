package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type TokenManager struct {
	token     string
	expiry    time.Time
	mu        sync.RWMutex
	appKey    string
	appSecret string
	authURL   string
}

var (
	tokenManager *TokenManager
	once         sync.Once
)

func NewTokenManager() *TokenManager {
	once.Do(func() {
		tokenManager = &TokenManager{
			appKey:    os.Getenv("QP_APP_KEY"),
			appSecret: os.Getenv("QP_APP_SECRET"),
			authURL:   os.Getenv("QP_AUTH_API_BASE") + "/token",
		}

		// Validate configuration
		if tokenManager.appKey == "" || tokenManager.appSecret == "" || tokenManager.authURL == "/token" {
			log.Fatal("Missing required environment variables: QP_APP_KEY, QP_APP_SECRET, QP_AUTH_API_BASE")
		}
	})
	return tokenManager
}

func (tm *TokenManager) fetchToken() error {
	client := &http.Client{}

	// Prepare form data
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("scope", "device_full_access")

	// Create request
	req, err := http.NewRequest("POST", tm.authURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(tm.appKey, tm.appSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response body:", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get token: status %d", resp.StatusCode)
	}

	// Parse response
	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	tm.token = result.AccessToken
	tm.expiry = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)

	fmt.Println("🪙 Token acquired:")

	return nil
}

func (tm *TokenManager) GetToken() (string, error) {
	tm.mu.RLock()
	if time.Now().Add(30 * time.Second).After(tm.expiry) {
		fmt.Println("🟡 Token expired, refreshing...")
		tm.mu.RUnlock() // Release the read lock before acquiring the write lock
		tm.mu.Lock()
		defer tm.mu.Unlock()
		if err := tm.fetchToken(); err != nil {
			return "", fmt.Errorf("failed to refresh token: %v", err)
		}
	} else {
		fmt.Println("🟢 Token is still valid")
		defer tm.mu.RUnlock()
	}

	if tm.token == "" {
		return "", errors.New("no valid token available")
	}

	return tm.token, nil
}
