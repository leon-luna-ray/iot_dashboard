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

	// Set headers
	req.Header.Set("Authorization", "Bearer "+token)

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read error response", http.StatusInternalServerError)
			return
		}
		http.Error(w, string(body), resp.StatusCode)
		return
	}

	// Copy headers from response
	for k, v := range resp.Header {
		w.Header()[k] = v
	}

	// Copy status code and body
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
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
