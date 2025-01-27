// Main file for the server
package main

import (
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

	println(`
    Server started at http://localhost:9090 📡
    `)
	http.ListenAndServe(":9090", nil)
}
