// package api

// import (
// 	"crypto/hmac"
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"time"
// )

// type WebhookRequest struct {
// 	Signature struct {
// 		Timestamp int64  `json:"timestamp"`
// 		Token     string `json:"token"`
// 		Signature string `json:"signature"`
// 	} `json:"signature"`
// 	Payload json.RawMessage `json:"payload"`
// }

// func handleWebhook(w http.ResponseWriter, r *http.Request) {
// 	// Validate request method
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Read request body
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Error reading body", http.StatusBadRequest)
// 		return
// 	}

// 	// Unmarshal JSON
// 	var req WebhookRequest
// 	if err := json.Unmarshal(body, &req); err != nil {
// 		http.Error(w, "Invalid JSON", http.StatusBadRequest)
// 		return
// 	}

// 	// Signature validation
// 	appSecret := os.Getenv("QP_APP_SECRET")
// 	if appSecret == "" {
// 		http.Error(w, "Server misconfigured", http.StatusInternalServerError)
// 		return
// 	}

// 	// Check timestamp (within 5 minutes)
// 	if time.Now().Unix()-req.Signature.Timestamp > 300 {
// 		http.Error(w, "Expired request", http.StatusBadRequest)
// 		return
// 	}

// 	// Generate signature
// 	mac := hmac.New(sha256.New, []byte(appSecret))
// 	mac.Write([]byte(strconv.FormatInt(req.Signature.Timestamp, 10) + req.Signature.Token))
// 	expectedSignature := hex.EncodeToString(mac.Sum(nil))

// 	// Compare signatures
// 	if !hmac.Equal([]byte(expectedSignature), []byte(req.Signature.Signature)) {
// 		http.Error(w, "Invalid signature", http.StatusUnauthorized)
// 		return
// 	}

// 	// Process payload
// 	fmt.Printf("Received valid payload: %s\n", req.Payload)

// 	// Respond quickly
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(`{"status": "ok"}`))
// }
