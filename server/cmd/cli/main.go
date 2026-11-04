package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"justcompress/compressor"
)

func main() {
	// CLI flags
	input := flag.String("input", "", "Input image path (required)")
	output := flag.String("output", "", "Output image path (auto-generated if empty)")
	quality := flag.Int("quality", 75, "Compression quality (1-100)")
	format := flag.String("format", "", "Output format: jpeg, png (auto-detect from extension if empty)")
	lossless := flag.Bool("lossless", false, "Use lossless compression")
	maxWidth := flag.Uint("max-width", 0, "Max width (0 = no resize)")
	maxHeight := flag.Uint("max-height", 0, "Max height (0 = no resize)")

	flag.Parse()

	if *input == "" {
		fmt.Println("❌ Error: -input flag is required")
		fmt.Println("\nUsage:")
		fmt.Println("  justcompress-cli -input photo.jpg -quality 60")
		fmt.Println("  justcompress-cli -input photo.png -format jpeg -quality 80")
		fmt.Println("  justcompress-cli -input photo.jpg -max-width 1920 -max-height 1080")
		fmt.Println()
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Read input file
	data, err := os.ReadFile(*input)
	if err != nil {
		fmt.Printf("❌ Failed to read input file: %v\n", err)
		os.Exit(1)
	}

	originalSize := int64(len(data))

	// Auto-detect output format
	outFormat := *format
	if outFormat == "" {
		ext := strings.ToLower(filepath.Ext(*input))
		switch ext {
		case ".png":
			outFormat = "png"
		case ".webp":
			outFormat = "jpeg" // no Go webp encoder, fallback
		default:
			outFormat = "jpeg"
		}
	}

	// Auto-generate output path
	outPath := *output
	if outPath == "" {
		ext := ".jpg"
		switch outFormat {
		case "png":
			ext = ".png"
		case "webp":
			ext = ".webp"
		}
		base := strings.TrimSuffix(filepath.Base(*input), filepath.Ext(*input))
		dir := filepath.Dir(*input)
		outPath = filepath.Join(dir, base+"_compressed"+ext)
	}

	// Show info
	fmt.Println("┌─────────────────────────────────────────┐")
	fmt.Println("│         🗜️  JustCompress CLI             │")
	fmt.Println("└─────────────────────────────────────────┘")
	fmt.Printf("  📄 Input:    %s\n", *input)
	fmt.Printf("  📦 Size:     %s\n", compressor.FormatFileSize(originalSize))
	fmt.Printf("  🎯 Quality:  %d\n", *quality)
	fmt.Printf("  📐 Format:   %s\n", outFormat)
	if *lossless {
		fmt.Println("  🔒 Mode:     Lossless")
	} else {
		fmt.Println("  ⚡ Mode:     Lossy")
	}
	fmt.Println()

	// Compress
	start := time.Now()

	opts := compressor.Options{
		Quality:   *quality,
		Format:    outFormat,
		Lossless:  *lossless,
		MaxWidth:  *maxWidth,
		MaxHeight: *maxHeight,
	}

	compressed, err := compressor.Compress(data, opts)
	if err != nil {
		fmt.Printf("❌ Compression failed: %v\n", err)
		os.Exit(1)
	}

	duration := time.Since(start)
	compressedSize := int64(len(compressed))
	savings := (1.0 - float64(compressedSize)/float64(originalSize)) * 100.0

	// Write output
	if err := os.WriteFile(outPath, compressed, 0644); err != nil {
		fmt.Printf("❌ Failed to write output file: %v\n", err)
		os.Exit(1)
	}

	// Results
	fmt.Println("  ✅ Compression complete!")
	fmt.Printf("  📄 Output:   %s\n", outPath)
	fmt.Printf("  📦 Original: %s\n", compressor.FormatFileSize(originalSize))
	fmt.Printf("  📦 Result:   %s\n", compressor.FormatFileSize(compressedSize))

	if savings > 0 {
		fmt.Printf("  💰 Saved:    %.1f%%\n", savings)
	} else {
		fmt.Printf("  ⚠️  Size increased by %.1f%%\n", -savings)
	}

	fmt.Printf("  ⏱️  Time:     %s\n", duration)
	fmt.Println()
}
