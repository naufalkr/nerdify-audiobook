package dto

// Common response structures

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// ListResponse represents a paginated list response
type ListResponse struct {
	Items      interface{}        `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=20" binding:"min=1,max=100"`
}

// SearchRequest represents search parameters
type SearchRequest struct {
	Query string `form:"q"`
	PaginationRequest
}
