package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"justcompress/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Determine the directory of the executable to resolve relative paths
	staticDir := "../"
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		staticDir = filepath.Join(exeDir, "..") + string(os.PathSeparator)
	}
	// On Windows, ensure correct path separator
	if runtime.GOOS == "windows" {
		staticDir = filepath.ToSlash(staticDir)
	}

	mux := http.NewServeMux()

	// Serve static files (frontend)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	// API endpoints
	mux.HandleFunc("/api/compress", handlers.CompressHandler)
	mux.HandleFunc("/api/health", handlers.HealthHandler)

	// Serve the frontend
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			// Try to serve as a static file (css, js)
			filePath := filepath.Join(staticDir, r.URL.Path)
			if _, err := os.Stat(filePath); err == nil {
				http.ServeFile(w, r, filePath)
				return
			}
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	})

	// CORS middleware
	handler := corsMiddleware(mux)

	fmt.Printf("🚀 JustCompress server running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
