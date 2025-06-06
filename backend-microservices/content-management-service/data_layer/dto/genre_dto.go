package dto

// CreateGenreRequest represents the request to create a new genre
type CreateGenreRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

// UpdateGenreRequest represents the request to update a genre
type UpdateGenreRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

// GenreResponse represents the response for genre data
type GenreResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
