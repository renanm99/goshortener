package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type HealthResponse struct {
	Status      string    `json:"status"`
	Environment string    `json:"environment"`
	Version     string    `json:"version"`
	Timestamp   time.Time `json:"timestamp"`
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortID     string `json:"short_id"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	version := os.Getenv("VERSION")
	if version == "" {
		version = "dev"
	}

	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response := HealthResponse{
			Status:      "healthy",
			Environment: environment,
			Version:     version,
			Timestamp:   time.Now(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Readiness check endpoint
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		response := HealthResponse{
			Status:      "ready",
			Environment: environment,
			Version:     version,
			Timestamp:   time.Now(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Root endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"service":     "GoShorter URL Shortener",
				"version":     version,
				"environment": environment,
				"status":      "running",
			})
			return
		}

		// Simulate redirect (would lookup in database)
		http.Error(w, "Short URL not found", http.StatusNotFound)
	})

	// Shorten endpoint
	mux.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ShortenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.URL == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		// Generate a simple short ID (in production, this would be more sophisticated)
		shortID := fmt.Sprintf("abc%d", time.Now().Unix()%10000)

		response := ShortenResponse{
			ShortID:     shortID,
			OriginalURL: req.URL,
			ShortURL:    fmt.Sprintf("http://localhost:%s/%s", port, shortID),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	})

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 GoShorter starting on port %s", port)
	log.Printf("📦 Environment: %s", environment)
	log.Printf("🔖 Version: %s", version)
	log.Printf("🌐 Endpoints:")
	log.Printf("   - GET  /         - Service info")
	log.Printf("   - GET  /health   - Health check")
	log.Printf("   - GET  /ready    - Readiness check")
	log.Printf("   - POST /shorten  - Shorten URL")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}
