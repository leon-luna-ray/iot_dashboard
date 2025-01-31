package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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
	mu        sync.Mutex
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
	client := &http.Client{Timeout: 10 * time.Second}

	// Auth
	auth := fmt.Sprintf("%s:%s", tm.appKey, tm.appSecret)
	basicAuth := base64.StdEncoding.EncodeToString([]byte(auth))

	// Create request
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("scope", "device_full_access")

	req, err := http.NewRequest("POST", tm.authURL, strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Println("Request err:", err)
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+basicAuth)

	// Send request
	resp, err := client.Do(req)
	fmt.Println("Response:", resp)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get token: status %d", resp.StatusCode)
	}

	// Parse response
	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.token = result.AccessToken
	tm.expiry = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)

	return nil
}

func (tm *TokenManager) GetToken() (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Refresh token if expired or about to expire (within 30 seconds)
	if time.Now().Add(30 * time.Second).After(tm.expiry) {
		if err := tm.fetchToken(); err != nil {
			return "", fmt.Errorf("failed to refresh token: %v", err)
		}
	}

	if tm.token == "" {
		return "", errors.New("no valid token available")
	}

	return tm.token, nil
}
