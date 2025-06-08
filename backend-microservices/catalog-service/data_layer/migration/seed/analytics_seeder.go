package seed

import (
	"catalog-service/data_layer/entity"
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// AnalyticsSeeder seeds the analytics table
func AnalyticsSeeder(db *gorm.DB) error {
	log.Println("Seeding analytics...")

	// Check if analytics already exist to avoid duplicates
	var count int64
	db.Model(&entity.Analytics{}).Count(&count)
	if count > 0 {
		log.Println("Analytics already seeded, skipping...")
		return nil
	}

	// Get users and audiobooks from database
	var users []entity.User
	var audiobooks []entity.Audiobook

	if err := db.Find(&users).Error; err != nil {
		log.Printf("Error fetching users: %v", err)
		return err
	}

	if err := db.Find(&audiobooks).Error; err != nil {
		log.Printf("Error fetching audiobooks: %v", err)
		return err
	}

	if len(users) == 0 || len(audiobooks) == 0 {
		return fmt.Errorf("users and audiobooks must be seeded first")
	}

	// Seed random generator
	rand.Seed(time.Now().UnixNano())

	var analytics []entity.Analytics

	// Event types for analytics
	eventTypes := []string{
		"play_started",
		"play_paused",
		"play_resumed",
		"play_stopped",
		"play_completed",
		"bookmark_added",
		"bookmark_removed",
		"volume_changed",
		"speed_changed",
		"chapter_skipped",
		"chapter_rewound",
		"audiobook_liked",
		"audiobook_unliked",
		"audiobook_shared",
		"audiobook_downloaded",
		"audiobook_rated",
		"playlist_created",
		"playlist_updated",
		"playlist_deleted",
		"search_performed",
		"filter_applied",
		"profile_updated",
		"subscription_started",
		"subscription_cancelled",
		"feedback_submitted",
	}

	// Generate analytics events for the past 6 months
	now := time.Now()
	startDate := now.AddDate(0, -6, 0)

	// Generate analytics for each user
	for _, user := range users {
		// Generate between 10-100 events per user
		numEvents := rand.Intn(91) + 10

		for i := 0; i < numEvents; i++ {
			// Random audiobook
			audiobook := audiobooks[rand.Intn(len(audiobooks))]

			// Random event type
			eventType := eventTypes[rand.Intn(len(eventTypes))]

			// Random timestamp within the past 6 months
			randomDays := rand.Intn(180)
			randomHours := rand.Intn(24)
			randomMinutes := rand.Intn(60)
			eventTime := startDate.AddDate(0, 0, randomDays).Add(
				time.Duration(randomHours)*time.Hour +
					time.Duration(randomMinutes)*time.Minute,
			)

			analytic := entity.Analytics{
				AudiobookID:    audiobook.ID,
				UserID:         user.ID,
				EventType:      eventType,
				EventTimestamp: eventTime,
			}

			analytics = append(analytics, analytic)
		}
	}

	// Generate some popular audiobook events (more frequent events for certain audiobooks)
	popularAudiobooks := audiobooks[:min(5, len(audiobooks))]

	for _, audiobook := range popularAudiobooks {
		// Generate additional events for popular audiobooks
		for i := 0; i < 200; i++ {
			// Random user
			user := users[rand.Intn(len(users))]

			// Bias towards play events for popular content
			playEvents := []string{"play_started", "play_completed", "audiobook_liked", "audiobook_shared"}
			eventType := playEvents[rand.Intn(len(playEvents))]

			// Random timestamp within the past 3 months (more recent for popular content)
			randomDays := rand.Intn(90)
			randomHours := rand.Intn(24)
			randomMinutes := rand.Intn(60)
			eventTime := now.AddDate(0, -3, 0).AddDate(0, 0, randomDays).Add(
				time.Duration(randomHours)*time.Hour +
					time.Duration(randomMinutes)*time.Minute,
			)

			analytic := entity.Analytics{
				AudiobookID:    audiobook.ID,
				UserID:         user.ID,
				EventType:      eventType,
				EventTimestamp: eventTime,
			}

			analytics = append(analytics, analytic)
		}
	}

	// Generate daily analytics patterns
	for days := 0; days < 30; days++ {
		currentDay := now.AddDate(0, 0, -days)

		// Generate 50-200 events per day
		dailyEvents := rand.Intn(151) + 50

		for i := 0; i < dailyEvents; i++ {
			user := users[rand.Intn(len(users))]
			audiobook := audiobooks[rand.Intn(len(audiobooks))]
			eventType := eventTypes[rand.Intn(len(eventTypes))]

			// Random time during the day
			randomHour := rand.Intn(24)
			randomMinute := rand.Intn(60)
			randomSecond := rand.Intn(60)

			eventTime := time.Date(
				currentDay.Year(), currentDay.Month(), currentDay.Day(),
				randomHour, randomMinute, randomSecond, 0, currentDay.Location(),
			)

			analytic := entity.Analytics{
				AudiobookID:    audiobook.ID,
				UserID:         user.ID,
				EventType:      eventType,
				EventTimestamp: eventTime,
			}

			analytics = append(analytics, analytic)
		}
	}

	log.Printf("Generated %d analytics events", len(analytics))

	// Create analytics in batches
	batchSize := 100
	for i := 0; i < len(analytics); i += batchSize {
		end := i + batchSize
		if end > len(analytics) {
			end = len(analytics)
		}

		if err := db.Create(analytics[i:end]).Error; err != nil {
			log.Printf("Error seeding analytics batch %d-%d: %v", i, end-1, err)
			return err
		}
	}

	log.Printf("Successfully seeded %d analytics events", len(analytics))
	return nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
