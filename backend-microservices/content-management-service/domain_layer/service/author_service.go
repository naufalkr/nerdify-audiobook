package service

import (
	"content-management-service/data_layer/dto"
	"content-management-service/data_layer/entity"
	"content-management-service/data_layer/repository"
	"errors"

	"gorm.io/gorm"
)

type AuthorService struct {
	authorRepo repository.AuthorRepositoryInterface
}

func NewAuthorService(authorRepo repository.AuthorRepositoryInterface) *AuthorService {
	return &AuthorService{authorRepo: authorRepo}
}

// CreateAuthor creates a new author
func (s *AuthorService) CreateAuthor(req dto.CreateAuthorRequest) (*dto.AuthorResponse, error) {
	author := entity.Author{
		Name: req.Name,
	}

	if err := s.authorRepo.Create(&author); err != nil {
		return nil, err
	}

	return &dto.AuthorResponse{
		ID:   author.ID,
		Name: author.Name,
	}, nil
}

// GetAuthorByID retrieves an author by ID
func (s *AuthorService) GetAuthorByID(id uint) (*dto.AuthorResponse, error) {
	author, err := s.authorRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("author not found")
		}
		return nil, err
	}

	return &dto.AuthorResponse{
		ID:   author.ID,
		Name: author.Name,
	}, nil
}

// GetAllAuthors retrieves all authors with pagination
func (s *AuthorService) GetAllAuthors(req dto.PaginationRequest) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get paginated results
	authors, total, err := s.authorRepo.GetAll(offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var authorResponses []dto.AuthorResponse
	for _, author := range authors {
		authorResponses = append(authorResponses, dto.AuthorResponse{
			ID:   author.ID,
			Name: author.Name,
		})
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: authorResponses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateAuthor updates an existing author
func (s *AuthorService) UpdateAuthor(id uint, req dto.UpdateAuthorRequest) (*dto.AuthorResponse, error) {
	author, err := s.authorRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("author not found")
		}
		return nil, err
	}

	author.Name = req.Name
	if err := s.authorRepo.Update(author); err != nil {
		return nil, err
	}

	return &dto.AuthorResponse{
		ID:   author.ID,
		Name: author.Name,
	}, nil
}

// DeleteAuthor deletes an author
func (s *AuthorService) DeleteAuthor(id uint) error {
	_, err := s.authorRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("author not found")
		}
		return err
	}

	return s.authorRepo.Delete(id)
}

// SearchAuthors searches authors by name
func (s *AuthorService) SearchAuthors(req dto.SearchRequest) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get paginated results
	authors, total, err := s.authorRepo.SearchByName(req.Query, offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var authorResponses []dto.AuthorResponse
	for _, author := range authors {
		authorResponses = append(authorResponses, dto.AuthorResponse{
			ID:   author.ID,
			Name: author.Name,
		})
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: authorResponses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}
