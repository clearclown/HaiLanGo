package image

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	// ErrConversionFailed is returned when image conversion fails
	ErrConversionFailed = errors.New("image conversion failed")
	// ErrUnsupportedFormat is returned when the input format is not supported
	ErrUnsupportedFormat = errors.New("unsupported image format")
)

// ConvertHEICToJPEG converts HEIC image data to JPEG format
// Uses sips on macOS or ImageMagick on other platforms
func ConvertHEICToJPEG(heicData []byte) ([]byte, error) {
	// Create a temporary file for the HEIC input
	tmpDir := os.TempDir()
	inputPath := filepath.Join(tmpDir, fmt.Sprintf("input_%d.heic", os.Getpid()))
	outputPath := filepath.Join(tmpDir, fmt.Sprintf("output_%d.jpg", os.Getpid()))

	// Clean up temp files on exit
	defer os.Remove(inputPath)
	defer os.Remove(outputPath)

	// Write HEIC data to temp file
	if err := os.WriteFile(inputPath, heicData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}

	// Try sips first (macOS)
	if err := convertWithSips(inputPath, outputPath); err == nil {
		return os.ReadFile(outputPath)
	}

	// Try ImageMagick as fallback
	if err := convertWithImageMagick(inputPath, outputPath); err == nil {
		return os.ReadFile(outputPath)
	}

	// Try libheif (heif-convert) as another fallback
	if err := convertWithHeifConvert(inputPath, outputPath); err == nil {
		return os.ReadFile(outputPath)
	}

	return nil, ErrConversionFailed
}

// convertWithSips uses macOS sips command to convert HEIC to JPEG
func convertWithSips(inputPath, outputPath string) error {
	cmd := exec.Command("sips", "-s", "format", "jpeg", inputPath, "--out", outputPath)
	return cmd.Run()
}

// convertWithImageMagick uses ImageMagick's convert command
func convertWithImageMagick(inputPath, outputPath string) error {
	// Try 'magick' first (ImageMagick 7+)
	cmd := exec.Command("magick", "convert", inputPath, outputPath)
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Fallback to 'convert' (ImageMagick 6)
	cmd = exec.Command("convert", inputPath, outputPath)
	return cmd.Run()
}

// convertWithHeifConvert uses libheif's heif-convert command
func convertWithHeifConvert(inputPath, outputPath string) error {
	cmd := exec.Command("heif-convert", inputPath, outputPath)
	return cmd.Run()
}

// IsHEIC checks if the file is a HEIC image by extension
func IsHEIC(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".heic" || ext == ".heif"
}

// ConvertIfNeeded converts HEIC to JPEG if needed, otherwise returns original data
func ConvertIfNeeded(filename string, data []byte) ([]byte, string, error) {
	if IsHEIC(filename) {
		jpegData, err := ConvertHEICToJPEG(data)
		if err != nil {
			return nil, "", fmt.Errorf("failed to convert HEIC to JPEG: %w", err)
		}
		// Change extension to .jpg
		newFilename := strings.TrimSuffix(filename, filepath.Ext(filename)) + ".jpg"
		return jpegData, newFilename, nil
	}
	return data, filename, nil
}

// ConvertHEICReaderToJPEG converts HEIC from a reader to JPEG bytes
func ConvertHEICReaderToJPEG(reader io.Reader) ([]byte, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read HEIC data: %w", err)
	}
	return ConvertHEICToJPEG(data)
}

// ConvertHEICToJPEGReader converts HEIC to JPEG and returns a reader
func ConvertHEICToJPEGReader(heicData []byte) (io.Reader, error) {
	jpegData, err := ConvertHEICToJPEG(heicData)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(jpegData), nil
}
