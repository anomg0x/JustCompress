package compressor

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	// Register image decoders
	_ "image/gif"

	"golang.org/x/image/webp"

	nfntResize "github.com/nfnt/resize"
)

// Options defines compression parameters.
type Options struct {
	Quality   int    // 1-100, used for JPEG/WebP lossy
	Format    string // "jpeg", "png", "webp"
	Lossless  bool   // if true, lossless compression where possible
	MaxWidth  uint   // optional max width for resizing (0 = no resize)
	MaxHeight uint   // optional max height for resizing (0 = no resize)
}

// Compress takes raw image bytes and returns compressed image bytes.
func Compress(data []byte, opts Options) ([]byte, error) {
	// Decode the image
	reader := bytes.NewReader(data)
	img, format, err := image.Decode(reader)
	if err != nil {
		// Try WebP specifically if standard decode fails
		reader.Seek(0, io.SeekStart)
		img, err = webp.Decode(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to decode image: %w", err)
		}
		format = "webp"
	}

	_ = format // original format for logging

	// Resize if needed
	if opts.MaxWidth > 0 || opts.MaxHeight > 0 {
		img = resizeImage(img, opts.MaxWidth, opts.MaxHeight)
	}

	// Encode in the target format
	var buf bytes.Buffer
	switch opts.Format {
	case "jpeg", "jpg":
		err = encodeJPEG(&buf, img, opts)
	case "png":
		err = encodePNG(&buf, img, opts)
	case "webp":
		// Go stdlib doesn't have a WebP encoder, so we fall back to PNG or JPEG
		// For "webp" output, we encode as high-quality JPEG (since no native encoder)
		err = encodeJPEG(&buf, img, opts)
	default:
		err = encodeJPEG(&buf, img, opts)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return buf.Bytes(), nil
}

// encodeJPEG encodes an image as JPEG with the given quality.
func encodeJPEG(w io.Writer, img image.Image, opts Options) error {
	quality := opts.Quality
	if quality < 1 {
		quality = 1
	}
	if quality > 100 {
		quality = 100
	}

	return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
}

// encodePNG encodes an image as PNG.
// PNG is always lossless, but we can optimise the compression level.
func encodePNG(w io.Writer, img image.Image, opts Options) error {
	encoder := &png.Encoder{}

	if opts.Lossless || opts.Quality > 80 {
		encoder.CompressionLevel = png.BestCompression
	} else if opts.Quality > 50 {
		encoder.CompressionLevel = png.DefaultCompression
	} else {
		encoder.CompressionLevel = png.BestSpeed
	}

	return encoder.Encode(w, img)
}

// resizeImage resizes the image to fit within maxWidth×maxHeight while preserving aspect ratio.
func resizeImage(img image.Image, maxWidth, maxHeight uint) image.Image {
	bounds := img.Bounds()
	origW := uint(bounds.Dx())
	origH := uint(bounds.Dy())

	if maxWidth == 0 {
		maxWidth = origW
	}
	if maxHeight == 0 {
		maxHeight = origH
	}

	// Only downscale, never upscale
	if origW <= maxWidth && origH <= maxHeight {
		return img
	}

	return nfntResize.Thumbnail(maxWidth, maxHeight, img, nfntResize.Lanczos3)
}

// GetImageDimensions returns the dimensions of the image.
func GetImageDimensions(data []byte) (width, height int, err error) {
	reader := bytes.NewReader(data)
	cfg, _, err := image.DecodeConfig(reader)
	if err != nil {
		return 0, 0, err
	}
	return cfg.Width, cfg.Height, nil
}

// FormatFileSize formats bytes into a human-readable string.
func FormatFileSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
