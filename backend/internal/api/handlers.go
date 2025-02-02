package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func handleHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func handlePosts(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://dummyjson.com/posts?limit=10")
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// Device represents a single device from the API response
type Device struct {
	// Define fields based on the actual API response structure
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DevicesResponse represents the structure of the devices list API response
type DevicesResponse struct {
	Devices []Device `json:"devices"` // Adjust according to actual API structure
	Count   int      `json:"count"`
}

// BindRequest represents the request structure for device binding
type BindRequest struct {
	DeviceToken string `json:"device_token"`
	ProductID   int    `json:"product_id"`
	Timestamp   int64  `json:"timestamp"`
}

// fetchDevices handles the GET request to retrieve devices
func fetchDevices(tm *TokenManager, qpAPIBase string) ([]byte, int, error) {
	timestamp := time.Now().Unix()
	url := fmt.Sprintf("%s/devices?timestamp=%d", qpAPIBase, timestamp)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	token, err := tm.GetToken()
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}

// bindDevice handles the POST request to bind a device
func bindDevice() error {
	deviceToken := os.Getenv("QP_DEVICE_TOKEN")
	productIDStr := os.Getenv("QP_PRODUCT_ID")

	if deviceToken == "" || productIDStr == "" {
		return fmt.Errorf("QP_DEVICE_TOKEN or QP_PRODUCT_ID environment variables not set")
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		return fmt.Errorf("invalid QP_PRODUCT_ID: %v", err)
	}

	tm := NewTokenManager()
	token, err := tm.GetToken()
	if err != nil {
		return fmt.Errorf("failed to get token: %v", err)
	}

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	bindReq := BindRequest{
		DeviceToken: deviceToken,
		ProductID:   productID,
		Timestamp:   timestamp,
	}

	jsonBody, err := json.Marshal(bindReq)
	if err != nil {
		return fmt.Errorf("failed to marshal bind request: %v", err)
	}

	// TODO update
	url := "https://apis.cleargrass.com/v1/apis/devices"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create bind request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("bind request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bind failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
func handleSensors(w http.ResponseWriter, r *http.Request) {
	tm := NewTokenManager()
	qpAPIBase := os.Getenv("QP_API_BASE")
	if qpAPIBase == "" {
		http.Error(w, "QP_API_BASE not configured", http.StatusInternalServerError)
		return
	}

	// Initial fetch of devices
	devicesBody, statusCode, err := fetchDevices(tm, qpAPIBase)
	if err != nil {
		http.Error(w, "Failed to fetch devices: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the response to check device count
	var devicesResp DevicesResponse
	if err := json.Unmarshal(devicesBody, &devicesResp); err != nil {
		log.Printf("Failed to parse devices response: %v", err)
	} else {
		if devicesResp.Count == 0 || len(devicesResp.Devices) == 0 {
			log.Println("No devices found. Attempting to bind device...")
			if err := bindDevice(); err != nil {
				log.Printf("Device binding failed: %v", err)
			} else {
				// Re-fetch devices after successful binding
				devicesBody, statusCode, err = fetchDevices(tm, qpAPIBase)
				if err != nil {
					http.Error(w, "Failed to re-fetch devices: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(devicesBody)
}

// func handleSensors(w http.ResponseWriter, r *http.Request) {
// 	tm := NewTokenManager()
// 	token, err := tm.GetToken()

// 	if err != nil {
// 		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
// 		return
// 	}

// 	qpAPIBase := os.Getenv("QP_API_BASE")
// 	if qpAPIBase == "" {
// 		http.Error(w, "QP_API_BASE not configured", http.StatusInternalServerError)
// 		return
// 	}

// 	// Create timestamp and URL
// 	timestamp := time.Now().Unix()
// 	url := fmt.Sprintf("%s/devices?timestamp=%d", qpAPIBase, timestamp)

// 	// Create new request
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	req.Header.Set("Authorization", "Bearer "+token)

// 	// Send request
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	fmt.Printf("Response: %v\n", resp)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// Read the entire response body
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		http.Error(w, "Failed to read API response", http.StatusInternalServerError)
// 		return
// 	}

// 	// Debug: Log the response body
// 	log.Printf("API Response Body: %s", body)

// 	// Set headers and write response
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(resp.StatusCode)
// 	w.Write(body)
// }
