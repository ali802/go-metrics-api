package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// SystemStatus structures the JSON payload that our API will output
type SystemStatus struct {
	Status         string    `json:"status"`
	Timestamp      time.Time `json:"timestamp"`
	GoVersion      string    `json:"go_version"`
	NumCPU         int       `json:"num_cpu"`
	Goroutines     int       `json:"goroutines_active"`
	AllocMemMB     uint64    `json:"allocated_memory_mb"`
}

func main() {
	// Look for an environment variable named PORT, default to 8080 if not found
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()

	// --- Production Middleware Section ---
	r.Use(middleware.RequestID)   // Injects a unique ID into every request for tracing logs
	r.Use(middleware.RealIP)      // Captures the actual client IP, skipping reverse proxies
	r.Use(middleware.Logger)      // Automatically logs all incoming traffic to the console
	r.Use(middleware.Recoverer)   // Prevents the API from crashing completely if a bad error occurs
	r.Use(middleware.Timeout(60 * time.Second))

	// --- API Endpoints ---
	// 1. Basic Health Check for orchestrators like Kubernetes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "UP"})
	})

	// 2. Real-time Core Metrics Endpoint
	r.Get("/api/v1/metrics", handleMetrics)

	log.Printf("🚀 Production System Metrics API starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// Function to calculate and return internal runtime statistics
func handleMetrics(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m) // Queries Go's memory engine

	status := SystemStatus{
		Status:     "HEALTHY",
		Timestamp:  time.Now(),
		GoVersion:  runtime.Version(),
		NumCPU:     runtime.NumCPU(),          // Number of logical CPU cores allocated to container
		Goroutines: runtime.NumGoroutine(),    // Active internal execution threads running
		AllocMemMB: m.Alloc / 1024 / 1024,     // Converts raw bytes to megabytes
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
