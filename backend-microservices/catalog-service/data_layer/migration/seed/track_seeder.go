package seed

import (
	"catalog-service/data_layer/entity"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"gorm.io/gorm"
)

// AudiobookData represents the structure from updateAudiobooks.json
type AudiobookData struct {
	Data []struct {
		ID      int `json:"id"`
		Details struct {
			Title            string   `json:"title"`
			AuthorID         int      `json:"authorId"`
			Author           string   `json:"author"`
			Genres           []string `json:"genres"`
			YearOfPublishing int      `json:"yearOfPublishing"`
			Language         string   `json:"language"`
			ReaderID         int      `json:"readerId"`
			Reader           string   `json:"reader"`
			YoutubeVideoURL  string   `json:"youtubeVideoUrl"`
			LibrivoxPageURL  string   `json:"librivoxPageUrl"`
			TotalDuration    string   `json:"totalDuration"`
			Description      string   `json:"description"`
		} `json:"details"`
		Tracks []struct {
			Title    string `json:"title"`
			URL      string `json:"url"`
			Duration string `json:"duration"`
		} `json:"tracks"`
	} `json:"data"`
}

// TrackSeeder seeds the tracks table with actual data from updateAudiobooks.json
func TrackSeeder(db *gorm.DB) error {
	log.Println("Seeding tracks...")

	// Check if tracks already exist to avoid duplicates
	var count int64
	db.Model(&entity.Track{}).Count(&count)
	if count > 0 {
		log.Println("Tracks already seeded, skipping...")
		return nil
	}

	// Get audiobooks from database to map with JSON data
	var audiobooks []entity.Audiobook
	if err := db.Find(&audiobooks).Error; err != nil {
		log.Printf("Error fetching audiobooks: %v", err)
		return err
	}

	if len(audiobooks) == 0 {
		return fmt.Errorf("audiobooks must be seeded first")
	}

	// Load JSON data
	jsonData, err := loadAudiobookData()
	if err != nil {
		log.Printf("Error loading audiobook data: %v", err)
		return err
	}

	// Create a map for faster lookup of JSON data by title
	jsonDataMap := make(map[string][]struct {
		Title    string `json:"title"`
		URL      string `json:"url"`
		Duration string `json:"duration"`
	})

	for _, item := range jsonData.Data {
		if item.Details.Title != "" && len(item.Tracks) > 0 {
			jsonDataMap[item.Details.Title] = item.Tracks
		}
	}

	var tracks []entity.Track

	// Generate tracks for each audiobook
	for _, audiobook := range audiobooks {
		// Try to find matching tracks in JSON data
		if tracksData, exists := jsonDataMap[audiobook.Title]; exists {
			// Use actual track data from JSON
			for _, trackData := range tracksData {
				track := entity.Track{
					AudiobookID: audiobook.ID,
					Title:       trackData.Title,
					URL:         trackData.URL,
					Duration:    trackData.Duration,
				}
				tracks = append(tracks, track)
			}
			log.Printf("Added %d tracks for audiobook: %s", len(tracksData), audiobook.Title)
		} else {
			// Fallback: generate default tracks if no JSON data found
			log.Printf("No track data found for audiobook: %s, using fallback", audiobook.Title)
			fallbackTracks := generateFallbackTracks(audiobook)
			tracks = append(tracks, fallbackTracks...)
		}
	}

	// Create tracks in batches
	batchSize := 50
	for i := 0; i < len(tracks); i += batchSize {
		end := i + batchSize
		if end > len(tracks) {
			end = len(tracks)
		}

		if err := db.Create(tracks[i:end]).Error; err != nil {
			log.Printf("Error seeding tracks batch %d-%d: %v", i, end-1, err)
			return err
		}
	}

	log.Printf("Successfully seeded %d tracks", len(tracks))
	return nil
}

// loadAudiobookData loads the audiobook data from updateAudiobooks.json
func loadAudiobookData() (*AudiobookData, error) {
	// Get the path to the JSON file
	jsonPath := filepath.Join("delete_later", "updateAudiobooks.json")

	// Read the JSON file
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read audiobook data file: %v", err)
	}

	// Parse JSON
	var audiobookData AudiobookData
	if err := json.Unmarshal(data, &audiobookData); err != nil {
		return nil, fmt.Errorf("failed to parse audiobook data: %v", err)
	}

	return &audiobookData, nil
}

// generateFallbackTracks creates fallback tracks if no JSON data is found
func generateFallbackTracks(audiobook entity.Audiobook) []entity.Track {
	var tracks []entity.Track

	// Generate 5-10 basic tracks as fallback
	numTracks := 5

	for i := 1; i <= numTracks; i++ {
		track := entity.Track{
			AudiobookID: audiobook.ID,
			Title:       fmt.Sprintf("Chapter %d", i),
			URL:         fmt.Sprintf("https://example.com/audio/audiobook_%d/track_%02d.mp3", audiobook.ID, i),
			Duration:    "00:15:00", // Default 15 minutes
		}
		tracks = append(tracks, track)
	}

	return tracks
}
