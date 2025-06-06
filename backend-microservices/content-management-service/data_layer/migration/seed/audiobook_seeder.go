package seed

import (
	"content-management-service/data_layer/entity"
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// AudiobookSeeder seeds the audiobooks table
func AudiobookSeeder(db *gorm.DB) error {
	log.Println("Seeding audiobooks...")

	// Check if audiobooks already exist to avoid duplicates
	var count int64
	db.Model(&entity.Audiobook{}).Count(&count)
	if count > 0 {
		log.Println("Audiobooks already seeded, skipping...")
		return nil
	}

	// Get authors and readers from database
	var authors []entity.Author
	var readers []entity.Reader

	if err := db.Find(&authors).Error; err != nil {
		log.Printf("Error fetching authors: %v", err)
		return err
	}

	if err := db.Find(&readers).Error; err != nil {
		log.Printf("Error fetching readers: %v", err)
		return err
	}

	if len(authors) == 0 || len(readers) == 0 {
		return fmt.Errorf("authors and readers must be seeded first")
	}

	// Seed random generator
	rand.Seed(time.Now().UnixNano())

	audiobooks := []entity.Audiobook{
		{
			Title:            "Pride and Prejudice",
			AuthorID:         getAuthorIDByName(authors, "Jane Austen"),
			ReaderID:         getReaderIDByName(readers, "Karen Savage"),
			Description:      "A romantic novel about Elizabeth Bennet and Mr. Darcy, exploring themes of love, marriage, and social class in Georgian England.",
			ImageURL:         "https://example.com/images/pride-and-prejudice.jpg",
			Language:         "English",
			YearOfPublishing: 1813,
			TotalDuration:    "12:34:56",
		},
		{
			Title:            "Romeo and Juliet",
			AuthorID:         getAuthorIDByName(authors, "William Shakespeare"),
			ReaderID:         getReaderIDByName(readers, "Group"),
			Description:      "The tragic story of two young star-crossed lovers whose deaths ultimately reconcile their feuding families.",
			ImageURL:         "https://example.com/images/romeo-and-juliet.jpg",
			Language:         "English",
			YearOfPublishing: 1597,
			TotalDuration:    "3:45:23",
		},
		{
			Title:            "Treasure Island",
			AuthorID:         getAuthorIDByName(authors, "Robert Louis Stevenson"),
			ReaderID:         getReaderIDByName(readers, "Phil Chenevert"),
			Description:      "A coming-of-age story about young Jim Hawkins and his adventures in search of treasure.",
			ImageURL:         "https://example.com/images/treasure-island.jpg",
			Language:         "English",
			YearOfPublishing: 1883,
			TotalDuration:    "8:12:45",
		},
		{
			Title:            "Jane Eyre",
			AuthorID:         getAuthorIDByName(authors, "Charlotte Brontë"),
			ReaderID:         getReaderIDByName(readers, "Elizabeth Klett"),
			Description:      "The story of an orphaned girl who becomes a governess and falls in love with her brooding employer.",
			ImageURL:         "https://example.com/images/jane-eyre.jpg",
			Language:         "English",
			YearOfPublishing: 1847,
			TotalDuration:    "19:08:32",
		},
		{
			Title:            "Heidi",
			AuthorID:         getAuthorIDByName(authors, "Johanna Spyri"),
			ReaderID:         getReaderIDByName(readers, "Andrea Fiore"),
			Description:      "The story of a young orphan girl who lives with her grandfather in the Swiss Alps.",
			ImageURL:         "https://example.com/images/heidi.jpg",
			Language:         "English",
			YearOfPublishing: 1881,
			TotalDuration:    "6:45:18",
		},
		{
			Title:            "Twenty Thousand Leagues Under the Sea",
			AuthorID:         getAuthorIDByName(authors, "Jules Verne"),
			ReaderID:         getReaderIDByName(readers, "Bob Neufeld"),
			Description:      "The adventures of Captain Nemo and his submarine Nautilus as seen from the perspective of Professor Aronnax.",
			ImageURL:         "https://example.com/images/twenty-thousand-leagues.jpg",
			Language:         "English",
			YearOfPublishing: 1870,
			TotalDuration:    "14:23:07",
		},
		{
			Title:            "The Tell-Tale Heart",
			AuthorID:         getAuthorIDByName(authors, "Edgar Allan Poe"),
			ReaderID:         getReaderIDByName(readers, "John W. Michaels"),
			Description:      "A short story about an unnamed narrator who insists on his sanity after murdering an old man.",
			ImageURL:         "https://example.com/images/tell-tale-heart.jpg",
			Language:         "English",
			YearOfPublishing: 1843,
			TotalDuration:    "0:23:45",
		},
		{
			Title:            "The Great Gatsby",
			AuthorID:         getAuthorIDByName(authors, "F. Scott Fitzgerald"),
			ReaderID:         getReaderIDByName(readers, "Meredith Hughes"),
			Description:      "A critique of the American Dream set in the Jazz Age, following Nick Carraway and the mysterious Jay Gatsby.",
			ImageURL:         "https://example.com/images/great-gatsby.jpg",
			Language:         "English",
			YearOfPublishing: 1925,
			TotalDuration:    "5:32:18",
		},
		{
			Title:            "The Jungle Book",
			AuthorID:         getAuthorIDByName(authors, "Rudyard Kipling"),
			ReaderID:         getReaderIDByName(readers, "Sue Anderson"),
			Description:      "A collection of stories about Mowgli, a boy raised by wolves in the Indian jungle.",
			ImageURL:         "https://example.com/images/jungle-book.jpg",
			Language:         "English",
			YearOfPublishing: 1894,
			TotalDuration:    "7:15:42",
		},
		{
			Title:            "A Christmas Carol",
			AuthorID:         getAuthorIDByName(authors, "Charles Dickens"),
			ReaderID:         getReaderIDByName(readers, "Ruth Golding"),
			Description:      "The story of Ebenezer Scrooge's transformation on Christmas Eve through visits from three ghosts.",
			ImageURL:         "https://example.com/images/christmas-carol.jpg",
			Language:         "English",
			YearOfPublishing: 1843,
			TotalDuration:    "3:28:15",
		},
	}

	// Generate additional random audiobooks
	titles := []string{
		"The Adventures of Tom Sawyer", "Wuthering Heights", "Moby Dick", "The Picture of Dorian Gray",
		"The Time Machine", "Dracula", "Frankenstein", "Little Women", "The Scarlet Letter",
		"The Adventures of Huckleberry Finn", "Great Expectations", "Oliver Twist", "David Copperfield",
		"The Metamorphosis", "The Trial", "War and Peace", "Anna Karenina", "Crime and Punishment",
		"The Brothers Karamazov", "The Cherry Orchard", "Uncle Tom's Cabin", "The Call of the Wild",
		"White Fang", "The Sea-Wolf", "Martin Eden", "The Iron Heel", "The People of the Abyss",
		"The Wonderful Wizard of Oz", "The Wind in the Willows", "Peter Pan", "Alice's Adventures in Wonderland",
		"Through the Looking-Glass", "The Secret Garden", "A Little Princess", "The Railway Children",
		"Anne of Green Gables", "Rebecca of Sunnybrook Farm", "Pollyanna", "The Five Little Peppers",
		"Eight Cousins", "Rose in Bloom", "An Old-Fashioned Girl", "Jack and Jill", "Under the Lilacs",
	}

	descriptions := []string{
		"A timeless classic that has captivated readers for generations with its compelling characters and engaging plot.",
		"An extraordinary tale of adventure, romance, and human nature that explores the depths of the human condition.",
		"A masterpiece of literature that weaves together themes of love, loss, and redemption in a beautifully crafted narrative.",
		"A gripping story that takes listeners on an unforgettable journey through different worlds and experiences.",
		"A profound exploration of society, morality, and the human spirit that resonates with readers across cultures.",
		"An epic tale of courage, sacrifice, and triumph that showcases the best and worst of human nature.",
		"A thought-provoking work that challenges conventional wisdom and offers new perspectives on life and society.",
		"A beautiful story of growth, discovery, and self-realization that speaks to the universal human experience.",
		"A compelling narrative that combines rich character development with intricate plot twists and turns.",
		"An inspiring tale of perseverance, hope, and the indomitable human spirit in the face of adversity.",
	}

	languages := []string{"English", "French", "German", "Spanish", "Italian", "Russian", "Portuguese"}
	durations := []string{"2:15:30", "4:32:18", "6:45:22", "8:12:45", "10:33:12", "12:28:38", "15:42:55", "18:15:20"}

	for i, title := range titles {
		if i >= len(authors) || i >= len(readers) {
			break
		}

		audiobook := entity.Audiobook{
			Title:            title,
			AuthorID:         authors[i%len(authors)].ID,
			ReaderID:         readers[i%len(readers)].ID,
			Description:      descriptions[i%len(descriptions)],
			ImageURL:         fmt.Sprintf("https://example.com/images/%s.jpg", fmt.Sprintf("book-%d", i+1)),
			Language:         languages[i%len(languages)],
			YearOfPublishing: 1800 + (i % 200),
			TotalDuration:    durations[i%len(durations)],
		}
		audiobooks = append(audiobooks, audiobook)
	}

	// Create audiobooks in batches
	batchSize := 10
	for i := 0; i < len(audiobooks); i += batchSize {
		end := i + batchSize
		if end > len(audiobooks) {
			end = len(audiobooks)
		}

		if err := db.Create(audiobooks[i:end]).Error; err != nil {
			log.Printf("Error seeding audiobooks batch %d-%d: %v", i, end-1, err)
			return err
		}
	}

	log.Printf("Successfully seeded %d audiobooks", len(audiobooks))
	return nil
}

// Helper function to get author ID by name
func getAuthorIDByName(authors []entity.Author, name string) uint {
	for _, author := range authors {
		if author.Name == name {
			return author.ID
		}
	}
	// Return first author if not found
	if len(authors) > 0 {
		return authors[0].ID
	}
	return 1
}

// Helper function to get reader ID by name
func getReaderIDByName(readers []entity.Reader, name string) uint {
	for _, reader := range readers {
		if reader.Name == name {
			return reader.ID
		}
	}
	// Return first reader if not found
	if len(readers) > 0 {
		return readers[0].ID
	}
	return 1
}
