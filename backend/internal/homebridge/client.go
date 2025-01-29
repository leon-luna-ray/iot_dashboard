package homebridge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

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
