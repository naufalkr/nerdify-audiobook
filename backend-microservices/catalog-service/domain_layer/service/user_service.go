package service

import (
	"catalog-service/data_layer/dto"
	"catalog-service/data_layer/entity"
	"catalog-service/data_layer/repository"
	"errors"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo repository.UserRepositoryInterface
}

func NewUserService(userRepo repository.UserRepositoryInterface) *UserService {
	return &UserService{userRepo: userRepo}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	user := entity.User{
		ID:   req.ID,
		Role: req.Role,
	}

	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:   user.ID,
		Role: user.Role,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &dto.UserResponse{
		ID:   user.ID,
		Role: user.Role,
	}, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*dto.UserResponse, error) {
	// This would require implementing GetByEmail in the repository
	// For now, return an error indicating it's not implemented
	return nil, errors.New("GetUserByEmail not implemented - email field not available in current schema")
}

// GetAllUsers retrieves all users with pagination
func (s *UserService) GetAllUsers(req dto.PaginationRequest) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get paginated results
	users, total, err := s.userRepo.GetAll(offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID:   user.ID,
			Role: user.Role,
		})
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: userResponses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(id string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.Role = req.Role
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:   user.ID,
		Role: user.Role,
	}, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id string) error {
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return s.userRepo.Delete(id)
}

// GetUsersByRole retrieves users by role with pagination
func (s *UserService) GetUsersByRole(role string, req dto.PaginationRequest) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get paginated results
	users, total, err := s.userRepo.GetByRole(role, offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID:   user.ID,
			Role: user.Role,
		})
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: userResponses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}
