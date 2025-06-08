package seed

import (
	"catalog-service/data_layer/entity"
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// AudiobookGenreSeeder seeds the audiobook_genres junction table
func AudiobookGenreSeeder(db *gorm.DB) error {
	log.Println("Seeding audiobook-genre relationships...")

	// Get audiobooks and genres from database
	var audiobooks []entity.Audiobook
	var genres []entity.Genre

	if err := db.Find(&audiobooks).Error; err != nil {
		log.Printf("Error fetching audiobooks: %v", err)
		return err
	}

	if err := db.Find(&genres).Error; err != nil {
		log.Printf("Error fetching genres: %v", err)
		return err
	}

	if len(audiobooks) == 0 || len(genres) == 0 {
		log.Println("Audiobooks and genres must be seeded first")
		return nil
	}

	// Seed random generator
	rand.Seed(time.Now().UnixNano())

	// Define some genre mappings based on audiobook titles
	genreMap := map[string][]string{
		"Pride and Prejudice":                   {"General Fiction", "Romance", "Humorous Fiction"},
		"Romeo and Juliet":                      {"Tragedy", "Romance", "Plays", "Drama"},
		"Treasure Island":                       {"Action & Adventure Fiction", "Children's Fiction"},
		"Jane Eyre":                             {"General Fiction", "Romance", "Coming of Age"},
		"Heidi":                                 {"Children's Fiction", "Family Life", "Coming of Age"},
		"Twenty Thousand Leagues Under the Sea": {"Action & Adventure Fiction", "Science Fiction", "Nautical & Marine Fiction"},
		"The Tell-Tale Heart":                   {"Horror & Supernatural Fiction", "Psychology"},
		"The Great Gatsby":                      {"General Fiction", "Tragedy", "Historical Fiction"},
		"The Jungle Book":                       {"Children's Fiction", "Action & Adventure Fiction", "Myths"},
		"A Christmas Carol":                     {"General Fiction", "Fantasy", "Holiday"},
	}

	// Process each audiobook
	processedCount := 0
	for _, audiobook := range audiobooks {
		// Check if this audiobook already has genre associations
		existingCount := db.Model(&audiobook).Association("Genres").Count()
		if existingCount > 0 {
			continue // Skip if already has genres
		}

		var selectedGenres []entity.Genre

		// Try to match by title first
		if genreNames, exists := genreMap[audiobook.Title]; exists {
			for _, genreName := range genreNames {
				for _, genre := range genres {
					if genre.Name == genreName {
						selectedGenres = append(selectedGenres, genre)
						break
					}
				}
			}
		}

		// If no specific mapping found, assign based on patterns or randomly
		if len(selectedGenres) == 0 {
			selectedGenres = assignGenresByPattern(audiobook, genres)
		}

		// Ensure at least 1-4 genres per audiobook
		if len(selectedGenres) == 0 {
			// Assign random genres
			numGenres := rand.Intn(3) + 1 // 1-3 genres
			selectedGenreIndices := make(map[int]bool)

			for len(selectedGenreIndices) < numGenres {
				idx := rand.Intn(len(genres))
				if !selectedGenreIndices[idx] {
					selectedGenreIndices[idx] = true
					selectedGenres = append(selectedGenres, genres[idx])
				}
			}
		}

		// Associate genres with audiobook
		if len(selectedGenres) > 0 {
			if err := db.Model(&audiobook).Association("Genres").Append(selectedGenres); err != nil {
				log.Printf("Error associating genres with audiobook %s: %v", audiobook.Title, err)
				continue
			}
			processedCount++
		}
	}

	log.Printf("Successfully created genre associations for %d audiobooks", processedCount)
	return nil
}

// assignGenresByPattern assigns genres based on audiobook title patterns
func assignGenresByPattern(audiobook entity.Audiobook, genres []entity.Genre) []entity.Genre {
	var selectedGenres []entity.Genre
	title := audiobook.Title

	// Helper function to find genre by name
	findGenre := func(name string) *entity.Genre {
		for _, genre := range genres {
			if genre.Name == name {
				return &genre
			}
		}
		return nil
	}

	// Pattern matching for genres
	patterns := map[string][]string{
		// Children's books
		"Little|Children|Kid|Young|Boy|Girl": {"Children's Fiction", "Family Life"},

		// Adventure books
		"Adventure|Quest|Journey|Island|Sea|Ocean|Wild|Frontier": {"Action & Adventure Fiction", "Action & Adventure"},

		// Romance patterns
		"Love|Heart|Wedding|Marriage|Bride": {"Romance", "General Fiction"},

		// Horror/Mystery patterns
		"Dark|Death|Murder|Mystery|Secret|Shadow|Ghost|Horror": {"Horror & Supernatural Fiction", "Mystery & Detective Fiction"},

		// Historical patterns
		"War|Battle|Empire|King|Queen|Prince|Princess|Castle|Medieval": {"Historical Fiction", "War & Military"},

		// Fantasy/Sci-Fi patterns
		"Magic|Wizard|Dragon|Space|Future|Time|Machine|Robot": {"Fantasy", "Science Fiction"},

		// Classical literature
		"Classic|Tale|Story|Fable": {"Classics", "Literature"},

		// Biography patterns
		"Life|Biography|Memoir|Story of": {"Biography & Autobiography", "Memoirs"},

		// Self-help patterns
		"How to|Guide|Success|Improve|Better": {"Self-Help", "Personal Development"},

		// Poetry patterns
		"Poem|Poetry|Verse|Sonnet": {"Poetry", "Literature"},
	}

	// Check title against patterns
	for pattern, genreNames := range patterns {
		if matchesPattern(title, pattern) {
			for _, genreName := range genreNames {
				if genre := findGenre(genreName); genre != nil {
					selectedGenres = append(selectedGenres, *genre)
				}
			}
			break // Use first matching pattern
		}
	}

	// If still no genres found, assign default fiction genre
	if len(selectedGenres) == 0 {
		if genre := findGenre("General Fiction"); genre != nil {
			selectedGenres = append(selectedGenres, *genre)
		}
		if genre := findGenre("Literature"); genre != nil {
			selectedGenres = append(selectedGenres, *genre)
		}
	}

	// Limit to maximum 4 genres and add some randomness
	if len(selectedGenres) > 4 {
		selectedGenres = selectedGenres[:4]
	}

	// Randomly add one more genre if we have less than 3
	if len(selectedGenres) < 3 && len(genres) > 0 {
		randomGenre := genres[rand.Intn(len(genres))]
		// Check if not already selected
		alreadySelected := false
		for _, selected := range selectedGenres {
			if selected.ID == randomGenre.ID {
				alreadySelected = true
				break
			}
		}
		if !alreadySelected {
			selectedGenres = append(selectedGenres, randomGenre)
		}
	}

	return selectedGenres
}

// matchesPattern checks if title matches any of the pipe-separated patterns
func matchesPattern(title, pattern string) bool {
	// Simple pattern matching - check if any keyword exists in title
	keywords := []string{}
	current := ""

	for _, char := range pattern {
		if char == '|' {
			if current != "" {
				keywords = append(keywords, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		keywords = append(keywords, current)
	}

	// Check if any keyword appears in title (case-insensitive)
	titleLower := toLower(title)
	for _, keyword := range keywords {
		keywordLower := toLower(keyword)
		if containsWord(titleLower, keywordLower) {
			return true
		}
	}

	return false
}

// toLower converts string to lowercase
func toLower(s string) string {
	result := ""
	for _, char := range s {
		if char >= 'A' && char <= 'Z' {
			result += string(char + 32)
		} else {
			result += string(char)
		}
	}
	return result
}

// containsWord checks if a word exists in text
func containsWord(text, word string) bool {
	return containsSubstring(text, word)
}

// containsSubstring helper function for substring search
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
