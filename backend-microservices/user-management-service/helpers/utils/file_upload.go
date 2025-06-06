package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

const (
	UploadDir            = "uploads"
	MaxFileSize          = 15 * 1024 * 1024 // 15MB maximum upload size
	MaxResizedFileSize   = 1 * 1024 * 1024  // 1MB for resized images (reduced from 3MB)
	DefaultResizeQuality = 80               // JPEG quality (1-100)
)

// AllowedImageTypes contains the allowed MIME types for images
var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

// CopyFile copies a file from src to dst
func CopyFile(src io.Reader, dst io.Writer) error {
	_, err := io.Copy(dst, src)
	return err
}

// UploadFile handles file upload and returns the file path
func UploadFile(file *multipart.FileHeader, subDir string) (string, error) {
	// Check file size for absolute maximum limit before attempting any processing
	if file.Size > MaxFileSize {
		return "", fmt.Errorf("file size exceeds maximum limit of %d MB", MaxFileSize/(1024*1024))
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Check file type
	contentType := file.Header.Get("Content-Type")
	if !AllowedImageTypes[contentType] {
		return "", fmt.Errorf("file type not allowed")
	}

	// Create upload directory if it doesn't exist
	uploadPath := filepath.Join(UploadDir, subDir)
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", err
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		// If extension is missing, derive from content type
		switch contentType {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		}
	}

	filename := fmt.Sprintf("%s_%s%s", uuid.New().String(), time.Now().Format("20060102150405"), ext)
	filePath := filepath.Join(uploadPath, filename)
	var fileBytes []byte

	// Resize the image if it's larger than the target size
	if file.Size > MaxResizedFileSize {
		// Reset the file pointer to the beginning
		src.Seek(0, io.SeekStart)

		fmt.Printf("Resizing image from %d bytes to target size %d bytes\n", file.Size, MaxResizedFileSize)

		// Resize the image
		fileBytes, err = ResizeImage(src, contentType, MaxResizedFileSize)
		if err != nil {
			return "", fmt.Errorf("failed to resize image: %w", err)
		}

		fmt.Printf("Image resized successfully, new size: %d bytes\n", len(fileBytes))
	} else {
		// Just read the file if it's already smaller than target size
		fileBytes, err = io.ReadAll(src)
		if err != nil {
			return "", fmt.Errorf("failed to read file: %w", err)
		}

		fmt.Printf("Image already within size limits: %d bytes\n", len(fileBytes))
	}

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Write the file content
	if _, err := dst.Write(fileBytes); err != nil {
		// Clean up on error
		os.Remove(filePath)
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Return relative path
	return filepath.Join(subDir, filename), nil
}

// DeleteFile deletes a file from the uploads directory
func DeleteFile(filePath string) error {
	fullPath := filepath.Join(UploadDir, filePath)
	return os.Remove(fullPath)
}

// GetFileURL returns the full URL for a file
func GetFileURL(filePath string) string {
	if filePath == "" {
		return ""
	}
	return fmt.Sprintf("/uploads/%s", filePath)
}

// ResizeImage resizes an image to reduce its file size
func ResizeImage(src io.Reader, contentType string, targetSize int64) ([]byte, error) {
	// Decode the image
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %w", err)
	}

	// Start with original dimensions
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	quality := DefaultResizeQuality

	fmt.Printf("Original image dimensions: %dx%d, starting quality: %d\n", width, height, quality)

	var buf bytes.Buffer
	resizeAttempts := 0
	maxAttempts := 15

	for resizeAttempts < maxAttempts {
		buf.Reset() // Resize the image if not the first attempt
		if resizeAttempts > 0 {
			// Reduce dimensions more aggressively for larger attempts
			scaleFactor := 0.8 // More aggressive initial scaling (was 0.9)
			if resizeAttempts > 3 {
				scaleFactor = 0.6 // Even more aggressive scaling after 3 attempts (was 0.7)
			}

			// Reduce dimensions by scaleFactor each attempt
			width = int(float64(width) * scaleFactor)
			height = int(float64(height) * scaleFactor)

			// Also reduce quality for JPEG after first attempt
			if contentType == "image/jpeg" {
				if quality > 80 {
					quality -= 15 // More aggressive quality reduction
				} else if quality > 60 {
					quality -= 10
				} else if quality > 30 {
					quality -= 5
				} else if quality > 20 {
					quality = 20 // Set minimum quality
				}
			}
		}

		// Resize image if dimensions were changed
		var resizedImg image.Image
		if resizeAttempts > 0 {
			resizedImg = imaging.Resize(img, width, height, imaging.Lanczos)
		} else {
			resizedImg = img
		}

		// Encode based on content type
		switch contentType {
		case "image/jpeg":
			err = jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: quality})
		case "image/png":
			err = png.Encode(&buf, resizedImg)
		case "image/gif":
			err = gif.Encode(&buf, resizedImg, &gif.Options{NumColors: 256})
		default:
			return nil, fmt.Errorf("unsupported image format: %s", contentType)
		}

		if err != nil {
			return nil, fmt.Errorf("error encoding image: %w", err)
		}

		// Check if we've reached target size
		if buf.Len() <= int(targetSize) || resizeAttempts == maxAttempts-1 {
			fmt.Printf("Resize attempt %d: dimensions %dx%d, quality %d, size %d bytes\n",
				resizeAttempts+1, width, height, quality, buf.Len())
			break
		}

		fmt.Printf("Resize attempt %d: dimensions %dx%d, quality %d, size %d bytes, still too large\n",
			resizeAttempts+1, width, height, quality, buf.Len())

		resizeAttempts++
	}

	// Always return the resized image, even if still slightly larger than target
	fmt.Printf("Final image size after resize: %d bytes (target: %d bytes)\n", buf.Len(), targetSize)

	return buf.Bytes(), nil
}
