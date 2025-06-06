package dto

// CreateTrackRequest represents the request to create a new track
type CreateTrackRequest struct {
	AudiobookID uint   `json:"audiobook_id" binding:"required"`
	Title       string `json:"title" binding:"required,min=1,max=255"`
	URL         string `json:"url" binding:"required,url,max=255"`
	Duration    string `json:"duration"`
}

// UpdateTrackRequest represents the request to update a track
type UpdateTrackRequest struct {
	Title    string `json:"title" binding:"required,min=1,max=255"`
	URL      string `json:"url" binding:"required,url,max=255"`
	Duration string `json:"duration"`
}

// TrackResponse represents the response for track data
type TrackResponse struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Duration string `json:"duration"`
}
