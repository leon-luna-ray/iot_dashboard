package api

import (
	"io/ioutil"
	"net/http"
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

// func handleHomebridge(w http.ResponseWriter, r *http.Request, client *HomebridgeClient) {
// 	body, err := client.GetAccessories()
// 	if err != nil {
// 		http.Error(w, "Failed to fetch accessories: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(body)
// }
