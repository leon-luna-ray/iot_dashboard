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
	homebridgeClient *HomebridgeClient
)

type Config struct {
	ENV             string
	HB_USERNAME     string
	HB_PASSWORD     string
	HB_BASE_API_URL string
}

var config Config

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config = Config{
		HB_USERNAME:     os.Getenv("HB_USERNAME"),
		HB_PASSWORD:     os.Getenv("HB_PASSWORD"),
		HB_BASE_API_URL: os.Getenv("HB_BASE_API_URL"),
	}

	// Initialize the Homebridge client
	homebridgeClient = NewHomebridgeClient(config.HB_BASE_API_URL, config.HB_USERNAME, config.HB_PASSWORD)
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
		body, err := homebridgeClient.GetAccessories()
		if err != nil {
			http.Error(w, "Failed to fetch accessories: "+err.Error(), http.StatusInternalServerError)
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

type HomebridgeClient struct {
	BaseURL    string
	Username   string
	Password   string
	Token      string
	TokenLock  sync.Mutex
	HTTPClient *http.Client
}

func NewHomebridgeClient(baseURL, username, password string) *HomebridgeClient {
	return &HomebridgeClient{
		BaseURL:    baseURL,
		Username:   username,
		Password:   password,
		HTTPClient: &http.Client{},
	}
}

func (c *HomebridgeClient) Login() error {
	data := map[string]string{
		"username": c.Username,
		"password": c.Password,
		"otp":      "",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to login, status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	c.TokenLock.Lock()
	c.Token = result["token"].(string)
	c.TokenLock.Unlock()

	return nil
}

func (c *HomebridgeClient) GetAccessories() ([]byte, error) {
	if c.Token == "" {
		if err := c.Login(); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest("GET", c.BaseURL+"/api/accessories", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
