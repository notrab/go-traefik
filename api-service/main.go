package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type APIResponse struct {
	Message   string            `json:"message"`
	UserID    string            `json:"user_id,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
	Headers   map[string]string `json:"headers,omitempty"`
}

func handleProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	rateLimit := r.Header.Get("X-Rate-Limit")

	response := APIResponse{
		Message:   "This is a protected endpoint",
		UserID:    userID,
		Timestamp: time.Now(),
		Headers: map[string]string{
			"X-User-ID":    userID,
			"X-Rate-Limit": rateLimit,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handlePublicEndpoint(w http.ResponseWriter, r *http.Request) {
	response := APIResponse{
		Message:   "This is a public endpoint (no auth required)",
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API service healthy"))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/protected", handleProtectedEndpoint)
	mux.HandleFunc("GET /api/data", handleProtectedEndpoint)
	mux.HandleFunc("GET /api/public", handlePublicEndpoint)
	mux.HandleFunc("GET /api/health", handleHealth)

	log.Println("API service starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
