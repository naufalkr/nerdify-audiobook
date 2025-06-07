package service

import (
	"content-management-service/data_layer/dto"
	"content-management-service/data_layer/entity"
	"content-management-service/data_layer/repository"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type AudiobookService struct {
	audiobookRepo repository.AudiobookRepositoryInterface
	authorRepo    repository.AuthorRepositoryInterface
	readerRepo    repository.ReaderRepositoryInterface
	genreRepo     repository.GenreRepositoryInterface
	trackRepo     repository.TrackRepositoryInterface
	analyticsRepo repository.AnalyticsRepositoryInterface
}

func NewAudiobookService(
	audiobookRepo repository.AudiobookRepositoryInterface,
	authorRepo repository.AuthorRepositoryInterface,
	readerRepo repository.ReaderRepositoryInterface,
	genreRepo repository.GenreRepositoryInterface,
	trackRepo repository.TrackRepositoryInterface,
	analyticsRepo repository.AnalyticsRepositoryInterface,
) *AudiobookService {
	return &AudiobookService{
		audiobookRepo: audiobookRepo,
		authorRepo:    authorRepo,
		readerRepo:    readerRepo,
		genreRepo:     genreRepo,
		trackRepo:     trackRepo,
		analyticsRepo: analyticsRepo,
	}
}

// CreateAudiobook creates a new audiobook
func (s *AudiobookService) CreateAudiobook(req dto.CreateAudiobookRequest) (*dto.AudiobookResponse, error) {
	// Validate author exists
	_, err := s.authorRepo.GetByID(req.AuthorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("author not found")
		}
		return nil, err
	}

	// Validate reader exists
	_, err = s.readerRepo.GetByID(req.ReaderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reader not found")
		}
		return nil, err
	}

	// Create audiobook
	audiobook := entity.Audiobook{
		Title:            req.Title,
		AuthorID:         req.AuthorID,
		ReaderID:         req.ReaderID,
		Description:      req.Description,
		ImageURL:         req.ImageURL,
		Language:         req.Language,
		YearOfPublishing: req.YearOfPublishing,
		TotalDuration:    req.TotalDuration,
	}

	if err := s.audiobookRepo.Create(&audiobook); err != nil {
		return nil, err
	}

	// Associate genres if provided
	if len(req.GenreIDs) > 0 {
		if err := s.audiobookRepo.AssignGenres(audiobook.ID, req.GenreIDs); err != nil {
			return nil, err
		}
	}

	// Get audiobook with relations for response
	audiobookWithRelations, err := s.audiobookRepo.GetByIDWithRelations(audiobook.ID)
	if err != nil {
		return nil, err
	}

	return s.convertToAudiobookResponse(audiobookWithRelations), nil
}

// GetAudiobookByID retrieves an audiobook by ID with all relationships
func (s *AudiobookService) GetAudiobookByID(id uint) (*dto.AudiobookResponse, error) {
	audiobook, err := s.audiobookRepo.GetByIDWithRelations(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("audiobook not found")
		}
		return nil, err
	}

	return s.convertToAudiobookResponse(audiobook), nil
}

// GetAllAudiobooks retrieves all audiobooks with pagination
func (s *AudiobookService) GetAllAudiobooks(req dto.PaginationRequest) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get paginated results with relations
	audiobooks, total, err := s.audiobookRepo.GetAllWithRelations(offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var audiobookResponses []dto.AudiobookListResponse
	for _, audiobook := range audiobooks {
		audiobookResponses = append(audiobookResponses, s.convertToAudiobookListResponse(&audiobook))
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: audiobookResponses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateAudiobook updates an existing audiobook
func (s *AudiobookService) UpdateAudiobook(id uint, req dto.UpdateAudiobookRequest) (*dto.AudiobookResponse, error) {
	// Get existing audiobook
	audiobook, err := s.audiobookRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("audiobook not found")
		}
		return nil, err
	}

	// Validate author exists
	_, err = s.authorRepo.GetByID(req.AuthorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("author not found")
		}
		return nil, err
	}

	// Validate reader exists
	_, err = s.readerRepo.GetByID(req.ReaderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reader not found")
		}
		return nil, err
	}

	// Update audiobook fields
	audiobook.Title = req.Title
	audiobook.AuthorID = req.AuthorID
	audiobook.ReaderID = req.ReaderID
	audiobook.Description = req.Description
	audiobook.ImageURL = req.ImageURL
	audiobook.Language = req.Language
	audiobook.YearOfPublishing = req.YearOfPublishing
	audiobook.TotalDuration = req.TotalDuration

	if err := s.audiobookRepo.Update(audiobook); err != nil {
		return nil, err
	}

	// Update genres if provided
	if len(req.GenreIDs) > 0 {
		// Remove all existing genres first
		if err := s.audiobookRepo.RemoveGenres(audiobook.ID, []uint{}); err != nil {
			return nil, err
		}
		// Add new genres
		if err := s.audiobookRepo.AssignGenres(audiobook.ID, req.GenreIDs); err != nil {
			return nil, err
		}
	}

	// Get updated audiobook with relations
	updatedAudiobook, err := s.audiobookRepo.GetByIDWithRelations(audiobook.ID)
	if err != nil {
		return nil, err
	}

	return s.convertToAudiobookResponse(updatedAudiobook), nil
}

// DeleteAudiobook deletes an audiobook
func (s *AudiobookService) DeleteAudiobook(id uint) error {
	// Check if audiobook exists
	_, err := s.audiobookRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("audiobook not found")
		}
		return err
	}

	// ✅ Clean up all related data before deleting audiobook

	// 1. Remove all genre associations
	if err := s.audiobookRepo.RemoveAllGenres(id); err != nil {
		return fmt.Errorf("failed to remove genre associations: %v", err)
	}

	// 2. Delete all tracks for this audiobook
	if err := s.trackRepo.DeleteByAudiobookID(id); err != nil {
		return fmt.Errorf("failed to delete tracks: %v", err)
	}

	// 3. Delete all analytics records for this audiobook
	if err := s.analyticsRepo.DeleteByAudiobookID(id); err != nil {
		return fmt.Errorf("failed to delete analytics: %v", err)
	}

	// 4. Finally delete the audiobook itself
	if err := s.audiobookRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete audiobook: %v", err)
	}

	return nil
}

// SearchAudiobooks searches audiobooks by title
func (s *AudiobookService) SearchAudiobooks(req dto.SearchRequest) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get paginated results
	audiobooks, total, err := s.audiobookRepo.SearchByTitle(req.Query, offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var audiobookResponses []dto.AudiobookListResponse
	for _, audiobook := range audiobooks {
		audiobookResponses = append(audiobookResponses, s.convertToAudiobookListResponse(&audiobook))
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: audiobookResponses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetAudiobooksByAuthorID retrieves audiobooks by author ID
func (s *AudiobookService) GetAudiobooksByAuthorID(authorID uint, req dto.PaginationRequest) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get paginated results
	audiobooks, total, err := s.audiobookRepo.GetByAuthorID(authorID, offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var audiobookResponses []dto.AudiobookListResponse
	for _, audiobook := range audiobooks {
		audiobookResponses = append(audiobookResponses, s.convertToAudiobookListResponse(&audiobook))
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: audiobookResponses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// Helper methods
func (s *AudiobookService) convertToAudiobookResponse(audiobook *entity.Audiobook) *dto.AudiobookResponse {
	response := &dto.AudiobookResponse{
		ID:               audiobook.ID,
		Title:            audiobook.Title,
		Description:      audiobook.Description,
		ImageURL:         audiobook.ImageURL,
		Language:         audiobook.Language,
		YearOfPublishing: audiobook.YearOfPublishing,
		TotalDuration:    audiobook.TotalDuration,
		Author: dto.AuthorResponse{
			ID:   audiobook.Author.ID,
			Name: audiobook.Author.Name,
		},
		Reader: dto.ReaderResponse{
			ID:   audiobook.Reader.ID,
			Name: audiobook.Reader.Name,
		},
	}

	// Convert genres
	for _, genre := range audiobook.Genres {
		response.Genres = append(response.Genres, dto.GenreResponse{
			ID:   genre.ID,
			Name: genre.Name,
		})
	}

	// Convert tracks
	for _, track := range audiobook.Tracks {
		response.Tracks = append(response.Tracks, dto.TrackResponse{
			ID:       track.ID,
			Title:    track.Title,
			URL:      track.URL,
			Duration: track.Duration,
		})
	}

	return response
}

func (s *AudiobookService) convertToAudiobookListResponse(audiobook *entity.Audiobook) dto.AudiobookListResponse {
	// Convert genres
	var genres []dto.GenreResponse
	for _, genre := range audiobook.Genres {
		genres = append(genres, dto.GenreResponse{
			ID:   genre.ID,
			Name: genre.Name,
		})
	}

	// Convert author
	var author *dto.AuthorResponse
	if audiobook.Author != nil {
		author = &dto.AuthorResponse{
			ID:   audiobook.Author.ID,
			Name: audiobook.Author.Name,
		}
	}

	// Convert reader
	var reader *dto.ReaderResponse
	if audiobook.Reader != nil {
		reader = &dto.ReaderResponse{
			ID:   audiobook.Reader.ID,
			Name: audiobook.Reader.Name,
		}
	}

	return dto.AudiobookListResponse{
		ID:                audiobook.ID,
		Title:             audiobook.Title,
		Author:            author,
		Reader:            reader,
		ImageURL:          audiobook.ImageURL,
		Language:          audiobook.Language,
		YearOfPublishing:  audiobook.YearOfPublishing,  // Tambahkan ini
		TotalDuration:     audiobook.TotalDuration,
		Genres:            genres,
	}
}

// GetAudiobooks retrieves audiobooks with filtering options
func (s *AudiobookService) GetAudiobooks(filter dto.AudiobookFilter, page, limit int) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (page - 1) * limit

	var audiobooks []entity.Audiobook
	var total int64
	var err error

	// Apply filters based on the provided filter criteria
	if filter.AuthorID > 0 {
		audiobooks, total, err = s.audiobookRepo.GetByAuthorID(filter.AuthorID, offset, limit)
	} else if filter.ReaderID > 0 {
		audiobooks, total, err = s.audiobookRepo.GetByReaderID(filter.ReaderID, offset, limit)
	} else if filter.GenreID > 0 {
		audiobooks, total, err = s.audiobookRepo.GetByGenreID(filter.GenreID, offset, limit)
	} else {
		audiobooks, total, err = s.audiobookRepo.GetAllWithRelations(offset, limit)
	}

	if err != nil {
		return nil, err
	}

	// Convert to response format
	var audiobookResponses []dto.AudiobookListResponse
	for _, audiobook := range audiobooks {
		audiobookResponses = append(audiobookResponses, s.convertToAudiobookListResponse(&audiobook))
	}

	// Calculate total pages
	totalPages := total / int64(limit)
	if total%int64(limit) > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: audiobookResponses,
		Pagination: dto.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: int(totalPages),
		},
	}, nil
}

// AddGenresToAudiobook adds genres to an audiobook
func (s *AudiobookService) AddGenresToAudiobook(audiobookID uint, genreIDs []uint) error {
	// Check if audiobook exists
	_, err := s.audiobookRepo.GetByID(audiobookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("audiobook not found")
		}
		return err
	}

	// Validate that all genres exist
	for _, genreID := range genreIDs {
		_, err := s.genreRepo.GetByID(genreID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("one or more genres not found")
			}
			return err
		}
	}

	return s.audiobookRepo.AssignGenres(audiobookID, genreIDs)
}

// RemoveGenresFromAudiobook removes genres from an audiobook
func (s *AudiobookService) RemoveGenresFromAudiobook(audiobookID uint, genreIDs []uint) error {
	// Check if audiobook exists
	_, err := s.audiobookRepo.GetByID(audiobookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("audiobook not found")
		}
		return err
	}

	return s.audiobookRepo.RemoveGenres(audiobookID, genreIDs)
}
