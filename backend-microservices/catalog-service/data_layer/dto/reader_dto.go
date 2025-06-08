package dto

// CreateReaderRequest represents the request to create a new reader
type CreateReaderRequest struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

// UpdateReaderRequest represents the request to update a reader
type UpdateReaderRequest struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

// ReaderResponse represents the response for reader data
type ReaderResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
