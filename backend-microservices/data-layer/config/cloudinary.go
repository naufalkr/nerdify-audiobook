package config

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"microservice/user/helpers/utils"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
}

func NewCloudinaryClient(config CloudinaryConfig) (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(config.CloudName, config.APIKey, config.APISecret)
	if err != nil {
		return nil, err
	}
	return cld, nil
}

type CloudinaryService struct {
	client *cloudinary.Cloudinary
}

func NewCloudinaryService() (*CloudinaryService, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("cloudinary credentials not found in environment variables")
	}

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary: %v", err)
	}

	return &CloudinaryService{
		client: cld,
	}, nil
}

type UploadResult struct {
	URL      string
	PublicID string
	Width    int
	Height   int
	Format   string
	Size     int
}

type ImageTransformation struct {
	Width   int
	Height  int
	Crop    string
	Quality int
	Format  string
}

// Constants for file naming and formatting
const (
	// TimestampFormat is the standard format for file timestamps
	TimestampFormat = "20060102150405"
	// RandomSuffixLength is the number of random bytes to add for uniqueness
	RandomSuffixLength = 4
	// MaxFilenameLength is the maximum length to use from original filename
	MaxFilenameLength = 40
	// FilenameCleanRegex matches characters that should be replaced in filenames
	FilenameCleanRegex = `[^a-zA-Z0-9_-]`
)

// generatePublicID creates a unique, SEO-friendly ID for Cloudinary assets
// Format: sanitized-original-name_timestamp_randomstring
func generatePublicID(originalFilename string) (string, error) {
	// Get timestamp component
	timestamp := time.Now().UTC().Format(TimestampFormat)

	// Generate random component for uniqueness
	randomBytes := make([]byte, RandomSuffixLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random suffix: %w", err)
	}
	randomSuffix := hex.EncodeToString(randomBytes)

	// Process original filename
	// 1. Remove extension
	filename := filepath.Base(originalFilename)
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))

	// 2. Clean unwanted characters, replace spaces with dashes
	reg := regexp.MustCompile(FilenameCleanRegex)
	filename = reg.ReplaceAllString(filename, "")
	filename = strings.ReplaceAll(filename, " ", "-")

	// 3. Truncate if too long
	if len(filename) > MaxFilenameLength {
		filename = filename[:MaxFilenameLength]
	}

	// 4. Ensure we have at least some name component
	if filename == "" {
		filename = "file"
	}

	// Combine all parts
	return fmt.Sprintf("%s_%s_%s", filename, timestamp, randomSuffix), nil
}

func (s *CloudinaryService) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string, transformations *ImageTransformation) (*UploadResult, error) {
	fmt.Printf("Starting Cloudinary upload process for file: %s\n", file.Filename)

	var fileBytes []byte
	var contentType string

	// Check file size to see if we need to resize
	if file.Size > 10*1024*1024 { // 10MB (Cloudinary free tier limit)
		fmt.Printf("File too large for Cloudinary free tier (%d bytes), resizing...\n", file.Size)

		// Open the file for resizing
		src, err := file.Open()
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return nil, err
		}

		// Get content type
		contentType = file.Header.Get("Content-Type")

		// Target size slightly under Cloudinary limit (9.5MB)
		targetSize := int64(9.5 * 1024 * 1024)

		// Use existing resize utility
		resizedBytes, err := utils.ResizeImage(src, contentType, targetSize)
		if err != nil {
			src.Close()
			return nil, fmt.Errorf("failed to resize image: %w", err)
		}

		fmt.Printf("Image resized successfully, new size: %d bytes\n", len(resizedBytes))
		fileBytes = resizedBytes
		src.Close()
	} else {
		// Just read the original file if it's already small enough
		src, err := file.Open()
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return nil, err
		}
		defer src.Close()

		// Read file into memory
		fileBytes, err = io.ReadAll(src)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
	}
	// Create a new reader for the file bytes
	reader := bytes.NewReader(fileBytes)

	// Generate filename as timestamp only, without extension
	// Cloudinary will handle the file extension automatically based on the content
	publicID, err := generatePublicID(file.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to generate public ID: %w", err)
	}

	uploadParams := uploader.UploadParams{
		Folder:   folder,
		PublicID: publicID,
	}

	// Apply transformations if provided
	if transformations != nil {
		uploadParams.Transformation = fmt.Sprintf("c_%s,w_%d,h_%d,q_%d,f_%s",
			transformations.Crop,
			transformations.Width,
			transformations.Height,
			transformations.Quality,
			transformations.Format)
		fmt.Printf("Applied transformations: %s\n", uploadParams.Transformation)
	}

	fmt.Printf("Sending to Cloudinary with folder: %s, publicID: %s\n", folder, publicID)
	fmt.Printf("Cloudinary credentials - Cloud Name: %s, API Key: %s...\n",
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY")[:4]+"***")

	result, err := s.client.Upload.Upload(ctx, reader, uploadParams)
	if err != nil {
		fmt.Printf("Cloudinary upload error: %v\n", err)
		return nil, err
	}

	return &UploadResult{
		URL:      result.SecureURL,
		PublicID: result.PublicID,
		Width:    result.Width,
		Height:   result.Height,
		Format:   result.Format,
		Size:     result.Bytes,
	}, nil
}

func (s *CloudinaryService) DeleteFile(ctx context.Context, publicID string) error {
	_, err := s.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}

func (s *CloudinaryService) GenerateImageURL(publicID string, transformations *ImageTransformation) string {
	url, _ := s.client.Image(publicID)

	if transformations != nil {
		url.Transformation = fmt.Sprintf("c_%s,w_%d,h_%d,q_%d,f_%s",
			transformations.Crop,
			transformations.Width,
			transformations.Height,
			transformations.Quality,
			transformations.Format)
	}

	urlStr, _ := url.String()
	return urlStr
}
