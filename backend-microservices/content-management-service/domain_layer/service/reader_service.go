package service

import (
	"content-management-service/data_layer/dto"
	"content-management-service/data_layer/entity"
	"content-management-service/data_layer/repository"
	"errors"

	"gorm.io/gorm"
)

type ReaderService struct {
	readerRepo repository.ReaderRepositoryInterface
}

func NewReaderService(readerRepo repository.ReaderRepositoryInterface) *ReaderService {
	return &ReaderService{readerRepo: readerRepo}
}

// CreateReader creates a new reader
func (s *ReaderService) CreateReader(req dto.CreateReaderRequest) (*dto.ReaderResponse, error) {
	reader := entity.Reader{
		Name: req.Name,
	}

	if err := s.readerRepo.Create(&reader); err != nil {
		return nil, err
	}

	return &dto.ReaderResponse{
		ID:   reader.ID,
		Name: reader.Name,
	}, nil
}

// GetReaderByID retrieves a reader by ID
func (s *ReaderService) GetReaderByID(id uint) (*dto.ReaderResponse, error) {
	reader, err := s.readerRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reader not found")
		}
		return nil, err
	}

	return &dto.ReaderResponse{
		ID:   reader.ID,
		Name: reader.Name,
	}, nil
}

// GetAllReaders retrieves all readers with pagination
func (s *ReaderService) GetAllReaders(req dto.PaginationRequest) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get paginated results
	readers, total, err := s.readerRepo.GetAll(offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var readerResponses []dto.ReaderResponse
	for _, reader := range readers {
		readerResponses = append(readerResponses, dto.ReaderResponse{
			ID:   reader.ID,
			Name: reader.Name,
		})
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: readerResponses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateReader updates an existing reader
func (s *ReaderService) UpdateReader(id uint, req dto.UpdateReaderRequest) (*dto.ReaderResponse, error) {
	reader, err := s.readerRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reader not found")
		}
		return nil, err
	}

	reader.Name = req.Name
	if err := s.readerRepo.Update(reader); err != nil {
		return nil, err
	}

	return &dto.ReaderResponse{
		ID:   reader.ID,
		Name: reader.Name,
	}, nil
}

// DeleteReader deletes a reader
func (s *ReaderService) DeleteReader(id uint) error {
	_, err := s.readerRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("reader not found")
		}
		return err
	}

	return s.readerRepo.Delete(id)
}

// SearchReaders searches readers by name
func (s *ReaderService) SearchReaders(req dto.SearchRequest) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get paginated results
	readers, total, err := s.readerRepo.SearchByName(req.Query, offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var readerResponses []dto.ReaderResponse
	for _, reader := range readers {
		readerResponses = append(readerResponses, dto.ReaderResponse{
			ID:   reader.ID,
			Name: reader.Name,
		})
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: readerResponses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// SearchReadersByName searches readers by name with pagination
func (s *ReaderService) SearchReadersByName(name string, page, limit int) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated results
	readers, total, err := s.readerRepo.SearchByName(name, offset, limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var readerResponses []dto.ReaderResponse
	for _, reader := range readers {
		readerResponses = append(readerResponses, dto.ReaderResponse{
			ID:   reader.ID,
			Name: reader.Name,
		})
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: readerResponses,
		Pagination: dto.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}
