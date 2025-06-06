package dto

// CreateAudiobookRequest represents the request to create a new audiobook
type CreateAudiobookRequest struct {
	Title            string `json:"title" binding:"required,min=1,max=255"`
	AuthorID         uint   `json:"author_id" binding:"required"`
	ReaderID         uint   `json:"reader_id" binding:"required"`
	Description      string `json:"description"`
	ImageURL         string `json:"image_url"`
	Language         string `json:"language"`
	YearOfPublishing int    `json:"year_of_publishing"`
	TotalDuration    string `json:"total_duration"`
	GenreIDs         []uint `json:"genre_ids"`
}

// UpdateAudiobookRequest represents the request to update an audiobook
type UpdateAudiobookRequest struct {
	Title            string `json:"title" binding:"required,min=1,max=255"`
	AuthorID         uint   `json:"author_id" binding:"required"`
	ReaderID         uint   `json:"reader_id" binding:"required"`
	Description      string `json:"description"`
	ImageURL         string `json:"image_url"`
	Language         string `json:"language"`
	YearOfPublishing int    `json:"year_of_publishing"`
	TotalDuration    string `json:"total_duration"`
	GenreIDs         []uint `json:"genre_ids"`
}

// AudiobookResponse represents the response for audiobook data
type AudiobookResponse struct {
	ID               uint            `json:"id"`
	Title            string          `json:"title"`
	Author           AuthorResponse  `json:"author"`
	Reader           ReaderResponse  `json:"reader"`
	Description      string          `json:"description"`
	ImageURL         string          `json:"image_url"`
	Language         string          `json:"language"`
	YearOfPublishing int             `json:"year_of_publishing"`
	TotalDuration    string          `json:"total_duration"`
	Genres           []GenreResponse `json:"genres"`
	Tracks           []TrackResponse `json:"tracks,omitempty"`
}

// AudiobookListResponse represents the response for audiobook list
type AudiobookListResponse struct {
	ID               uint            `json:"id"`
	Title            string          `json:"title"`
	Author           AuthorResponse  `json:"author"`
	Reader           ReaderResponse  `json:"reader"`
	ImageURL         string          `json:"image_url"`
	Language         string          `json:"language"`
	YearOfPublishing int             `json:"year_of_publishing"`
	TotalDuration    string          `json:"total_duration"`
	Genres           []GenreResponse `json:"genres"`
}

// AudiobookFilter represents filtering options for audiobooks
type AudiobookFilter struct {
	AuthorID uint `form:"author_id"`
	ReaderID uint `form:"reader_id"`
	GenreID  uint `form:"genre_id"`
}
