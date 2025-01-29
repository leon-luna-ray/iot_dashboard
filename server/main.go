package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

//go:embed internal/static/dist/*
var frontend embed.FS

var (
	token     string
	tokenLock sync.Mutex
)

type Config struct {
	ENV         string
	HB_USERNAME string
	HB_PASSWORD string
}

var config Config

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config = Config{
		HB_USERNAME: os.Getenv("HB_USERNAME"),
		HB_PASSWORD: os.Getenv("HB_PASSWORD"),
	}
}

func main() {

	staticFiles, err := fs.Sub(frontend, "internal/static/dist")
	if err != nil {
		panic(err)
	}

	fs := http.FileServer(http.FS(staticFiles))

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "internal/static/dist/index.html")
			return
		}
		fs.ServeHTTP(w, r)
	}))

	http.HandleFunc("/api/v1/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	http.HandleFunc("/api/v1/posts", func(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/api/v1/homebridge", func(w http.ResponseWriter, r *http.Request) {
		// TODO make separte function for this
		tokenLock.Lock()
		if token == "" {
			fmt.Println("Token is empty")
			tokenLock.Unlock()
			err := homebridgeLogin()
			if err != nil {
				http.Error(w, "Failed to login to Homebridge API", http.StatusInternalServerError)
				return
			}
			tokenLock.Lock()
		}
		tokenLock.Unlock()
		print(token)
		req, err := http.NewRequest("GET", "http://localhost:8581/api/accessories", nil)
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Failed to connect to Homebridge API", http.StatusInternalServerError)
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

func homebridgeLogin() error {
	// TODO: move this to a separate function
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	data := map[string]string{
		"username": config.HB_USERNAME,
		"password": config.HB_PASSWORD,
		"otp":      "",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://192.168.50.2:8581/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	fmt.Println("Response: ", resp)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to login, status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	tokenLock.Lock()
	token = result["token"].(string)
	tokenLock.Unlock()

	return nil
}
