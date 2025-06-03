package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Audiobook struct {
	ID2 struct {
		Oid string `json:"oid"`
	} `json:"_id"`
	ID          int    `json:"id"`
	ImgUrl      string `json:"imgUrl"`
	Details     struct {
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
}

func connectDB() *sql.DB {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Connected to database successfully!")
	return db
}

func seedDatabase(db *sql.DB) {
	// Create tables dengan schema yang sesuai SQLC
	createTablesQuery := `
	-- Enable UUID extension
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	-- Books table (sesuai dengan yang diharapkan SQLC)
	CREATE TABLE IF NOT EXISTS books (
		book_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		title TEXT,
		author TEXT,
		reader TEXT,
		genre TEXT,
		summary TEXT,          -- Column yang dibutuhkan aplikasi
		image_url TEXT,
		librivox_url TEXT,
		language TEXT,
		total_duration TEXT,
		year_of_publishing INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Seeks table untuk menyimpan posisi audio
	CREATE TABLE IF NOT EXISTS seeks (
		id SERIAL PRIMARY KEY,
		user_id TEXT,
		book_chapter TEXT,
		seek_position INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id, book_chapter)
	);

	-- Authors table
	CREATE TABLE IF NOT EXISTS authors (
		id SERIAL PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Genres table
	CREATE TABLE IF NOT EXISTS genres (
		id SERIAL PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(createTablesQuery)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}
	fmt.Println("✓ Database tables created successfully!")
}

func seedAudiobooks(db *sql.DB) {
	// Read audiobooks.json
	data, err := ioutil.ReadFile("data/audiobooks.json")
	if err != nil {
		log.Fatal("Failed to read audiobooks.json:", err)
	}

	var response struct {
		Data []Audiobook `json:"data"`
	}

	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Fatal("Failed to parse JSON:", err)
	}

	// Insert query sesuai dengan struktur table
	insertQuery := `
	INSERT INTO books (book_id, title, author, reader, genre, summary, image_url, librivox_url, language, total_duration, year_of_publishing)
	VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	ON CONFLICT (book_id) DO NOTHING;`

	for _, book := range response.Data {
		// Convert genres array to string
		genreStr := ""
		if len(book.Details.Genres) > 0 {
			genreStr = strings.Join(book.Details.Genres, ", ")
		}

		// Use description as summary
		summary := book.Details.Description
		if summary == "" {
			summary = fmt.Sprintf("Audiobook: %s by %s", book.Details.Title, book.Details.Author)
		}

		_, err := db.Exec(insertQuery,
			book.Details.Title,                 // title
			book.Details.Author,                // author
			book.Details.Reader,                // reader
			genreStr,                          // genre
			summary,                           // summary (dari description)
			book.ImgUrl,                       // image_url
			book.Details.LibrivoxPageURL,      // librivox_url
			book.Details.Language,             // language
			book.Details.TotalDuration,        // total_duration
			book.Details.YearOfPublishing,     // year_of_publishing
		)

		if err != nil {
			log.Printf("Failed to insert book %s: %v", book.Details.Title, err)
		} else {
			fmt.Printf("✓ Inserted: %s by %s\n", book.Details.Title, book.Details.Author)
		}
	}

	fmt.Printf("Seeding completed! Inserted %d audiobooks.\n", len(response.Data))
}

func seedAuthors(db *sql.DB) {
	data, err := ioutil.ReadFile("data/authors.json")
	if err != nil {
		log.Println("authors.json not found, skipping...")
		return
	}

	var authors []map[string]interface{}
	json.Unmarshal(data, &authors)

	for _, author := range authors {
		if name, ok := author["name"].(string); ok {
			_, err := db.Exec("INSERT INTO authors (name) VALUES ($1) ON CONFLICT (name) DO NOTHING", name)
			if err != nil {
				log.Printf("Failed to insert author %s: %v", name, err)
			} else {
				fmt.Printf("✓ Inserted author: %s\n", name)
			}
		}
	}
}

func seedGenres(db *sql.DB) {
	data, err := ioutil.ReadFile("data/genres.json")
	if err != nil {
		log.Println("genres.json not found, skipping...")
		return
	}

	var genres []map[string]interface{}
	json.Unmarshal(data, &genres)

	for _, genre := range genres {
		if name, ok := genre["name"].(string); ok {
			_, err := db.Exec("INSERT INTO genres (name) VALUES ($1) ON CONFLICT (name) DO NOTHING", name)
			if err != nil {
				log.Printf("Failed to insert genre %s: %v", name, err)
			} else {
				fmt.Printf("✓ Inserted genre: %s\n", name)
			}
		}
	}
}

func main() {
	fmt.Println("Starting database seeding...")

	db := connectDB()
	defer db.Close()

	// Create database schema first
	seedDatabase(db)

	// Seed all data
	seedAudiobooks(db)
	seedAuthors(db)
	seedGenres(db)

	fmt.Println("Database seeding completed successfully!")
}