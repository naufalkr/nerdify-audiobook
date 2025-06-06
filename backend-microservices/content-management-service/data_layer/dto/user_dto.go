package dto

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	ID   string `json:"id" binding:"required,min=1,max=255"`
	Role string `json:"role" binding:"required,min=1,max=50"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Role string `json:"role" binding:"required,min=1,max=50"`
}

// UserResponse represents the response for user data
type UserResponse struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}
