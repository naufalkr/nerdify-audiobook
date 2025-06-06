package seed

import (
	"content-management-service/data_layer/entity"
	"log"

	"gorm.io/gorm"
)

// GenreSeeder seeds the genres table
func GenreSeeder(db *gorm.DB) error {
	log.Println("Seeding genres...")

	genres := []entity.Genre{
		{Name: "General Fiction"},
		{Name: "Humorous Fiction"},
		{Name: "Romance"},
		{Name: "Tragedy"},
		{Name: "Horror & Supernatural Fiction"},
		{Name: "Plays"},
		{Name: "Comedy"},
		{Name: "Dramatic Readings"},
		{Name: "Action & Adventure Fiction"},
		{Name: "Self-Help"},
		{Name: "Nautical & Marine Fiction"},
		{Name: "Myths"},
		{Name: "Legends & Fairy Tales"},
		{Name: "Action & Adventure"},
		{Name: "Children's Fiction"},
		{Name: "Coming of Age"},
		{Name: "Family Life"},
		{Name: "Psychology"},
		{Name: "Philosophy"},
		{Name: "Historical Fiction"},
		{Name: "Science Fiction"},
		{Name: "Fantasy"},
		{Name: "Mystery & Detective Fiction"},
		{Name: "Thriller"},
		{Name: "Crime Fiction"},
		{Name: "War & Military"},
		{Name: "Western"},
		{Name: "Biography & Autobiography"},
		{Name: "Memoirs"},
		{Name: "History"},
		{Name: "Politics"},
		{Name: "Social Sciences"},
		{Name: "Religion & Spirituality"},
		{Name: "Christianity"},
		{Name: "Islam"},
		{Name: "Buddhism"},
		{Name: "Hinduism"},
		{Name: "Judaism"},
		{Name: "Mythology"},
		{Name: "Poetry"},
		{Name: "Literature"},
		{Name: "Classics"},
		{Name: "Drama"},
		{Name: "Epic"},
		{Name: "Satire"},
		{Name: "Travel & Geography"},
		{Name: "Nature"},
		{Name: "Science"},
		{Name: "Mathematics"},
		{Name: "Physics"},
		{Name: "Chemistry"},
		{Name: "Biology"},
		{Name: "Medicine"},
		{Name: "Technology"},
		{Name: "Engineering"},
		{Name: "Computer Science"},
		{Name: "Economics"},
		{Name: "Business"},
		{Name: "Finance"},
		{Name: "Marketing"},
		{Name: "Management"},
		{Name: "Leadership"},
		{Name: "Education"},
		{Name: "Art"},
		{Name: "Music"},
		{Name: "Architecture"},
		{Name: "Photography"},
		{Name: "Cooking"},
		{Name: "Health & Wellness"},
		{Name: "Fitness"},
		{Name: "Nutrition"},
		{Name: "Mental Health"},
		{Name: "Relationships"},
		{Name: "Parenting"},
		{Name: "Personal Development"},
		{Name: "Motivation"},
		{Name: "Success"},
		{Name: "Productivity"},
		{Name: "Time Management"},
		{Name: "Communication"},
		{Name: "Public Speaking"},
		{Name: "Negotiation"},
		{Name: "Sales"},
		{Name: "Entrepreneurship"},
		{Name: "Startup"},
		{Name: "Innovation"},
		{Name: "Creativity"},
		{Name: "Design"},
		{Name: "Writing"},
		{Name: "Journalism"},
		{Name: "Media"},
		{Name: "Entertainment"},
		{Name: "Sports"},
		{Name: "Games"},
		{Name: "Hobbies"},
		{Name: "Crafts"},
		{Name: "Gardening"},
		{Name: "Pets"},
		{Name: "Animals"},
		{Name: "Environment"},
		{Name: "Climate Change"},
		{Name: "Sustainability"},
		{Name: "Green Living"},
		{Name: "Alternative Energy"},
	}

	// Check if genres already exist to avoid duplicates
	var count int64
	db.Model(&entity.Genre{}).Count(&count)
	if count > 0 {
		log.Println("Genres already seeded, skipping...")
		return nil
	}

	// Create genres in batches
	batchSize := 25
	for i := 0; i < len(genres); i += batchSize {
		end := i + batchSize
		if end > len(genres) {
			end = len(genres)
		}

		if err := db.Create(genres[i:end]).Error; err != nil {
			log.Printf("Error seeding genres batch %d-%d: %v", i, end-1, err)
			return err
		}
	}

	log.Printf("Successfully seeded %d genres", len(genres))
	return nil
}
