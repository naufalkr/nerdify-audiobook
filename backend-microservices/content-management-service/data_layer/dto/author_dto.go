package dto

// CreateAuthorRequest represents the request to create a new author
type CreateAuthorRequest struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

// UpdateAuthorRequest represents the request to update an author
type UpdateAuthorRequest struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

// AuthorResponse represents the response for author data
type AuthorResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
