package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"justcompress/compressor"
)

// CompressResponse is the JSON response for the compression endpoint.
type CompressResponse struct {
	Success        bool    `json:"success"`
	OriginalSize   int64   `json:"originalSize"`
	CompressedSize int64   `json:"compressedSize"`
	Savings        float64 `json:"savings"`
	Format         string  `json:"format"`
	Filename       string  `json:"filename"`
	Duration       string  `json:"duration"`
	Error          string  `json:"error,omitempty"`
}

// CompressHandler handles image compression requests.
func CompressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	start := time.Now()

	// Parse multipart form (max 50MB)
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		sendError(w, "Failed to parse form data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, header, err := r.FormFile("image")
	if err != nil {
		sendError(w, "No image file provided: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read image data
	imageData, err := io.ReadAll(file)
	if err != nil {
		sendError(w, "Failed to read image data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	originalSize := int64(len(imageData))

	// Get compression parameters
	qualityStr := r.FormValue("quality")
	quality := 75 // default
	if qualityStr != "" {
		q, err := strconv.Atoi(qualityStr)
		if err == nil && q >= 1 && q <= 100 {
			quality = q
		}
	}

	outputFormat := r.FormValue("format")
	if outputFormat == "" {
		outputFormat = "jpeg"
	}

	lossless := r.FormValue("lossless") == "true"

	// Compress the image
	opts := compressor.Options{
		Quality:  quality,
		Format:   outputFormat,
		Lossless: lossless,
	}

	compressed, err := compressor.Compress(imageData, opts)
	if err != nil {
		sendError(w, "Compression failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	compressedSize := int64(len(compressed))
	savings := 0.0
	if originalSize > 0 {
		savings = (1.0 - float64(compressedSize)/float64(originalSize)) * 100.0
	}

	duration := time.Since(start)

	// Set content type for the compressed image
	contentType := "image/jpeg"
	ext := ".jpg"
	switch outputFormat {
	case "png":
		contentType = "image/png"
		ext = ".png"
	case "webp":
		contentType = "image/webp"
		ext = ".webp"
	}

	// Check if the client wants JSON metadata or the binary image
	if r.FormValue("metadata") == "true" {
		resp := CompressResponse{
			Success:        true,
			OriginalSize:   originalSize,
			CompressedSize: compressedSize,
			Savings:        savings,
			Format:         outputFormat,
			Filename:       header.Filename,
			Duration:       duration.String(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Return the compressed image binary
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"compressed_%s%s\"", header.Filename, ext))
	w.Header().Set("X-Original-Size", strconv.FormatInt(originalSize, 10))
	w.Header().Set("X-Compressed-Size", strconv.FormatInt(compressedSize, 10))
	w.Header().Set("X-Savings", fmt.Sprintf("%.1f%%", savings))
	w.Header().Set("X-Duration", duration.String())
	w.Header().Set("Content-Length", strconv.Itoa(len(compressed)))
	w.Write(compressed)
}

// HealthHandler returns the server health status.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"service": "JustCompress Go Server",
	})
}

func sendError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(CompressResponse{
		Success: false,
		Error:   msg,
	})
}
