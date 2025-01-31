package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func handleSensors(w http.ResponseWriter, r *http.Request) {
	tm := NewTokenManager()
	token, err := tm.GetToken()
	fmt.Println("Token:", token)

	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	qpAPIBase := os.Getenv("QP_API_BASE")
	if qpAPIBase == "" {
		http.Error(w, "QP_API_BASE not configured", http.StatusInternalServerError)
		return
	}

	// Create timestamp and URL
	timestamp := time.Now().Unix()
	url := fmt.Sprintf("%s/devices?timestamp=%d", qpAPIBase, timestamp)

	// Create new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Printf("Response: %v\n", resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the entire response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read API response", http.StatusInternalServerError)
		return
	}

	// Debug: Log the response body
	log.Printf("API Response Body: %s", body)

	// Set headers and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

// func handleHomebridge(w http.ResponseWriter, r *http.Request, client *HomebridgeClient) {
// 	body, err := client.GetAccessories()
// 	if err != nil {
// 		http.Error(w, "Failed to fetch accessories: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(body)
// }
