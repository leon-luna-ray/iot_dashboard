// Main file for the server
package main

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func main() {
	fs := http.FileServer(http.Dir("../internal/static/dist"))

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join("../internal/static/dist", "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	}))

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	http.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
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
	})

	println(`
    Server started at http://localhost:9090 📡
    `)
	http.ListenAndServe(":9090", nil)
}
