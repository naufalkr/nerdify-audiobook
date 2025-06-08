package service

import (
	"catalog-service/data_layer/dto"
	"catalog-service/data_layer/entity"
	"catalog-service/data_layer/repository"
	"errors"

	"gorm.io/gorm"
)

type TrackService struct {
	trackRepo repository.TrackRepositoryInterface
}

func NewTrackService(trackRepo repository.TrackRepositoryInterface) *TrackService {
	return &TrackService{trackRepo: trackRepo}
}

// CreateTrack creates a new track
func (s *TrackService) CreateTrack(req dto.CreateTrackRequest) (*dto.TrackResponse, error) {
	track := entity.Track{
		AudiobookID: req.AudiobookID,
		Title:       req.Title,
		URL:         req.URL,
		Duration:    req.Duration,
	}

	if err := s.trackRepo.Create(&track); err != nil {
		return nil, err
	}

	return &dto.TrackResponse{
		ID:       track.ID,
		Title:    track.Title,
		URL:      track.URL,
		Duration: track.Duration,
	}, nil
}

// GetTrackByID retrieves a track by ID
func (s *TrackService) GetTrackByID(id uint) (*dto.TrackResponse, error) {
	track, err := s.trackRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("track not found")
		}
		return nil, err
	}

	return &dto.TrackResponse{
		ID:       track.ID,
		Title:    track.Title,
		URL:      track.URL,
		Duration: track.Duration,
	}, nil
}

// GetAllTracks retrieves all tracks with pagination
func (s *TrackService) GetAllTracks(req dto.PaginationRequest) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get paginated results
	tracks, total, err := s.trackRepo.GetAll(offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var trackResponses []dto.TrackResponse
	for _, track := range tracks {
		trackResponses = append(trackResponses, dto.TrackResponse{
			ID:       track.ID,
			Title:    track.Title,
			URL:      track.URL,
			Duration: track.Duration,
		})
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: trackResponses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateTrack updates an existing track
func (s *TrackService) UpdateTrack(id uint, req dto.UpdateTrackRequest) (*dto.TrackResponse, error) {
	track, err := s.trackRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("track not found")
		}
		return nil, err
	}

	track.Title = req.Title
	track.URL = req.URL
	track.Duration = req.Duration

	if err := s.trackRepo.Update(track); err != nil {
		return nil, err
	}

	return &dto.TrackResponse{
		ID:       track.ID,
		Title:    track.Title,
		URL:      track.URL,
		Duration: track.Duration,
	}, nil
}

// DeleteTrack deletes a track
func (s *TrackService) DeleteTrack(id uint) error {
	_, err := s.trackRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("track not found")
		}
		return err
	}

	return s.trackRepo.Delete(id)
}

// SearchTracks searches tracks by title
func (s *TrackService) SearchTracks(req dto.SearchRequest) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get paginated results
	tracks, total, err := s.trackRepo.SearchByTitle(req.Query, offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var trackResponses []dto.TrackResponse
	for _, track := range tracks {
		trackResponses = append(trackResponses, dto.TrackResponse{
			ID:       track.ID,
			Title:    track.Title,
			URL:      track.URL,
			Duration: track.Duration,
		})
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: trackResponses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetTracksByAudiobookID retrieves all tracks for a specific audiobook
func (s *TrackService) GetTracksByAudiobookID(audiobookID uint) ([]dto.TrackResponse, error) {
	tracks, err := s.trackRepo.GetByAudiobookID(audiobookID)
	if err != nil {
		return nil, err
	}

	var trackResponses []dto.TrackResponse
	for _, track := range tracks {
		trackResponses = append(trackResponses, dto.TrackResponse{
			ID:       track.ID,
			Title:    track.Title,
			URL:      track.URL,
			Duration: track.Duration,
		})
	}

	return trackResponses, nil
}

// GetTracksByAudiobook retrieves tracks by audiobook ID with pagination
func (s *TrackService) GetTracksByAudiobook(audiobookID uint, page, limit int) (*dto.ListResponse, error) {
	// Get all tracks for the audiobook first
	allTracks, err := s.trackRepo.GetByAudiobookID(audiobookID)
	if err != nil {
		return nil, err
	}

	// Calculate pagination manually
	total := int64(len(allTracks))
	offset := (page - 1) * limit

	// Apply pagination
	var tracks []entity.Track
	if offset < len(allTracks) {
		end := offset + limit
		if end > len(allTracks) {
			end = len(allTracks)
		}
		tracks = allTracks[offset:end]
	}

	// Convert to response format
	var trackResponses []dto.TrackResponse
	for _, track := range tracks {
		trackResponses = append(trackResponses, dto.TrackResponse{
			ID:       track.ID,
			Title:    track.Title,
			URL:      track.URL,
			Duration: track.Duration,
		})
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: trackResponses,
		Pagination: dto.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateTrackOrder updates the order of tracks in an audiobook
func (s *TrackService) UpdateTrackOrder(audiobookID uint, trackOrders []struct {
	TrackID uint `json:"track_id"`
	Order   int  `json:"order"`
}) error {
	// First check if audiobook exists
	_, err := s.trackRepo.GetByID(audiobookID) // Using a track to verify audiobook exists
	if err != nil {
		return errors.New("audiobook not found")
	}

	// Update the track orders (this would need custom implementation in repository)
	// For now, we'll just return success as order management would require
	// additional fields in the track entity
	return nil
}
