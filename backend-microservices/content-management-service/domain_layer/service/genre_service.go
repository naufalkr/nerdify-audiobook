package service

import (
	"content-management-service/data_layer/dto"
	"content-management-service/data_layer/entity"
	"content-management-service/data_layer/repository"
	"errors"

	"gorm.io/gorm"
)

type GenreService struct {
	genreRepo repository.GenreRepositoryInterface
}

func NewGenreService(genreRepo repository.GenreRepositoryInterface) *GenreService {
	return &GenreService{genreRepo: genreRepo}
}

// CreateGenre creates a new genre
func (s *GenreService) CreateGenre(req dto.CreateGenreRequest) (*dto.GenreResponse, error) {
	genre := entity.Genre{
		Name: req.Name,
	}

	if err := s.genreRepo.Create(&genre); err != nil {
		return nil, err
	}

	return &dto.GenreResponse{
		ID:   genre.ID,
		Name: genre.Name,
	}, nil
}

// GetGenreByID retrieves a genre by ID
func (s *GenreService) GetGenreByID(id uint) (*dto.GenreResponse, error) {
	genre, err := s.genreRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("genre not found")
		}
		return nil, err
	}

	return &dto.GenreResponse{
		ID:   genre.ID,
		Name: genre.Name,
	}, nil
}

// GetAllGenres retrieves all genres with pagination
func (s *GenreService) GetAllGenres(req dto.PaginationRequest) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get paginated results
	genres, total, err := s.genreRepo.GetAll(offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var genreResponses []dto.GenreResponse
	for _, genre := range genres {
		genreResponses = append(genreResponses, dto.GenreResponse{
			ID:   genre.ID,
			Name: genre.Name,
		})
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: genreResponses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateGenre updates an existing genre
func (s *GenreService) UpdateGenre(id uint, req dto.UpdateGenreRequest) (*dto.GenreResponse, error) {
	genre, err := s.genreRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("genre not found")
		}
		return nil, err
	}

	genre.Name = req.Name
	if err := s.genreRepo.Update(genre); err != nil {
		return nil, err
	}

	return &dto.GenreResponse{
		ID:   genre.ID,
		Name: genre.Name,
	}, nil
}

// DeleteGenre deletes a genre
func (s *GenreService) DeleteGenre(id uint) error {
	_, err := s.genreRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("genre not found")
		}
		return err
	}

	return s.genreRepo.Delete(id)
}

// SearchGenres searches genres by name
func (s *GenreService) SearchGenres(req dto.SearchRequest) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get paginated results
	genres, total, err := s.genreRepo.SearchByName(req.Query, offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var genreResponses []dto.GenreResponse
	for _, genre := range genres {
		genreResponses = append(genreResponses, dto.GenreResponse{
			ID:   genre.ID,
			Name: genre.Name,
		})
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: genreResponses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetGenresByIDs retrieves multiple genres by their IDs
func (s *GenreService) GetGenresByIDs(ids []uint) ([]dto.GenreResponse, error) {
	genres, err := s.genreRepo.GetByIDs(ids)
	if err != nil {
		return nil, err
	}

	var genreResponses []dto.GenreResponse
	for _, genre := range genres {
		genreResponses = append(genreResponses, dto.GenreResponse{
			ID:   genre.ID,
			Name: genre.Name,
		})
	}

	return genreResponses, nil
}
