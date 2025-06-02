package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
)

type AuthService struct {
	redis *redis.Client
}

func NewAuthService(redisURL string) *AuthService {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	return &AuthService{redis: rdb}
}

func (a *AuthService) handleAuth(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	// Extract API key (assuming format: "Bearer <api-key>")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	apiKey := parts[1]

	// Validate API key against Redis
	ctx := context.Background()
	userID, err := a.redis.Get(ctx, fmt.Sprintf("api_key:%s", apiKey)).Result()
	if err == redis.Nil {
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("Redis error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set response headers for Traefik to forward
	w.Header().Set("X-User-ID", userID)
	w.Header().Set("X-Rate-Limit", "1000") // Example rate limit

	// Return 200 OK to allow request through
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (a *AuthService) handleCreateKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.FormValue("user_id")
	apiKey := r.FormValue("api_key")

	if userID == "" || apiKey == "" {
		http.Error(w, "user_id and api_key are required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err := a.redis.Set(ctx, fmt.Sprintf("api_key:%s", apiKey), userID, 0).Err()
	if err != nil {
		log.Printf("Redis error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "API key created for user: %s\n", userID)
}

func (a *AuthService) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Auth service healthy"))
}

func main() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	authService := NewAuthService(redisURL)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /auth", authService.handleAuth)
	mux.HandleFunc("POST /auth", authService.handleAuth)
	mux.HandleFunc("POST /auth/create-key", authService.handleCreateKey)
	mux.HandleFunc("GET /auth/health", authService.handleHealth)

	log.Println("Auth service starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
