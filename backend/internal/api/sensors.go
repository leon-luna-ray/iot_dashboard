package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type DeviceInfo struct {
	MAC       string `json:"mac"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
}

type Device struct {
	Info DeviceInfo `json:"info"`
}

type DevicesResponse struct {
	Total   int      `json:"total"`
	Devices []Device `json:"devices"`
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

	fmt.Println("Response Status:", resp)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	fmt.Println("Response Body:", string(body))
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

	// Fetch devices
	devicesBody, statusCode, err := fetchDevices(tm, qpAPIBase)
	if err != nil {
		http.Error(w, "Failed to fetch devices: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var apiResponse DevicesResponse
	if err := json.Unmarshal(devicesBody, &apiResponse); err != nil {
		log.Printf("Failed to parse devices response: %v", err)
		http.Error(w, "Failed to parse devices response", http.StatusInternalServerError)
		return
	}

	fmt.Println("Total devices:", apiResponse.Total)
	fmt.Println("Devices array length:", len(apiResponse.Devices))

	// If no devices are found, attempt to bind the device
	// TODO Move to frontend as an api endpoint
	if apiResponse.Total == 0 || len(apiResponse.Devices) == 0 {
		log.Println("❓ No devices found. Attempting to bind device...")
		if err := bindDevice(); err != nil {
			log.Printf("Device binding failed: %v", err)
			http.Error(w, "Device binding failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Re-fetch devices after successful binding
		devicesBody, statusCode, err = fetchDevices(tm, qpAPIBase)
		if err != nil {
			http.Error(w, "Failed to re-fetch devices: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Return the devices list
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(devicesBody)
}
