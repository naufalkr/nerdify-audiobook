package seed

import (
	"content-management-service/data_layer/entity"
	"log"

	"gorm.io/gorm"
)

// AuthorSeeder seeds the authors table
func AuthorSeeder(db *gorm.DB) error {
	log.Println("Seeding authors...")

	authors := []entity.Author{
		{Name: "Jane Austen"},
		{Name: "LibriVox Volunteers"},
		{Name: "William Shakespeare"},
		{Name: "Robert Louis Stevenson"},
		{Name: "Charlotte Brontë"},
		{Name: "Johanna Spyri"},
		{Name: "Jules Verne"},
		{Name: "William Walker Atkinson"},
		{Name: "Edgar Allan Poe"},
		{Name: "Jacob & Wilhelm Grimm"},
		{Name: "F. Scott Fitzgerald"},
		{Name: "Rudyard Kipling"},
		{Name: "Geronimo"},
		{Name: "Charles Dickens"},
		{Name: "Mark Twain"},
		{Name: "Arthur Conan Doyle"},
		{Name: "Lewis Carroll"},
		{Name: "Oscar Wilde"},
		{Name: "H.G. Wells"},
		{Name: "Bram Stoker"},
		{Name: "Mary Shelley"},
		{Name: "Emily Brontë"},
		{Name: "George Orwell"},
		{Name: "Virginia Woolf"},
		{Name: "James Joyce"},
		{Name: "Franz Kafka"},
		{Name: "Leo Tolstoy"},
		{Name: "Fyodor Dostoevsky"},
		{Name: "Anton Chekhov"},
		{Name: "Herman Melville"},
		{Name: "Walt Whitman"},
		{Name: "Emily Dickinson"},
		{Name: "Robert Frost"},
		{Name: "Edgar Rice Burroughs"},
		{Name: "L. Frank Baum"},
		{Name: "Kenneth Grahame"},
		{Name: "Louisa May Alcott"},
		{Name: "Nathaniel Hawthorne"},
		{Name: "Washington Irving"},
		{Name: "Henry James"},
		{Name: "Edith Wharton"},
		{Name: "Theodore Dreiser"},
		{Name: "Sinclair Lewis"},
		{Name: "Jack London"},
		{Name: "O. Henry"},
		{Name: "Willa Cather"},
		{Name: "Stephen Crane"},
		{Name: "Ambrose Bierce"},
		{Name: "Kate Chopin"},
		{Name: "Zane Grey"},
		{Name: "Joseph Conrad"},
		{Name: "Thomas Hardy"},
		{Name: "George Eliot"},
		{Name: "Anthony Trollope"},
		{Name: "William Makepeace Thackeray"},
		{Name: "Daniel Defoe"},
		{Name: "Jonathan Swift"},
		{Name: "Miguel de Cervantes"},
		{Name: "Alexandre Dumas"},
		{Name: "Victor Hugo"},
		{Name: "Gustave Flaubert"},
		{Name: "Émile Zola"},
		{Name: "Guy de Maupassant"},
		{Name: "Honoré de Balzac"},
		{Name: "Voltaire"},
		{Name: "Jean-Jacques Rousseau"},
		{Name: "Molière"},
		{Name: "Johann Wolfgang von Goethe"},
		{Name: "Friedrich Nietzsche"},
		{Name: "Thomas Mann"},
		{Name: "Hermann Hesse"},
		{Name: "Rainer Maria Rilke"},
		{Name: "Heinrich Heine"},
		{Name: "E.T.A. Hoffmann"},
		{Name: "Brothers Grimm"},
		{Name: "Hans Christian Andersen"},
		{Name: "Dante Alighieri"},
		{Name: "Giovanni Boccaccio"},
		{Name: "Niccolò Machiavelli"},
		{Name: "Petrarch"},
		{Name: "Homer"},
		{Name: "Sophocles"},
		{Name: "Euripides"},
		{Name: "Aeschylus"},
		{Name: "Aristotle"},
		{Name: "Plato"},
		{Name: "Socrates"},
		{Name: "Marcus Aurelius"},
		{Name: "Cicero"},
		{Name: "Virgil"},
		{Name: "Ovid"},
		{Name: "Horace"},
		{Name: "Sun Tzu"},
		{Name: "Confucius"},
		{Name: "Lao Tzu"},
		{Name: "Kahlil Gibran"},
		{Name: "Rumi"},
		{Name: "Omar Khayyam"},
		{Name: "Hafez"},
		{Name: "Saadi"},
		{Name: "Ferdowsi"},
		{Name: "Al-Ghazali"},
		{Name: "Ibn Khaldun"},
		{Name: "Averroes"},
		{Name: "Avicenna"},
	}

	// Check if authors already exist to avoid duplicates
	var count int64
	db.Model(&entity.Author{}).Count(&count)
	if count > 0 {
		log.Println("Authors already seeded, skipping...")
		return nil
	}

	// Create authors in batches
	batchSize := 20
	for i := 0; i < len(authors); i += batchSize {
		end := i + batchSize
		if end > len(authors) {
			end = len(authors)
		}

		if err := db.Create(authors[i:end]).Error; err != nil {
			log.Printf("Error seeding authors batch %d-%d: %v", i, end-1, err)
			return err
		}
	}

	log.Printf("Successfully seeded %d authors", len(authors))
	return nil
}
