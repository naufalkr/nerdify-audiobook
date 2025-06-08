package service

import (
	"catalog-service/data_layer/dto"
	"catalog-service/data_layer/entity"
	"catalog-service/data_layer/repository"
	"errors"
	"time"

	"gorm.io/gorm"
)

type AnalyticsService struct {
	analyticsRepo repository.AnalyticsRepositoryInterface
}

func NewAnalyticsService(analyticsRepo repository.AnalyticsRepositoryInterface) *AnalyticsService {
	return &AnalyticsService{analyticsRepo: analyticsRepo}
}

// CreateAnalytics creates a new analytics event
func (s *AnalyticsService) CreateAnalytics(userID string, req dto.CreateAnalyticsRequest) (*dto.AnalyticsResponse, error) {
	analytics := entity.Analytics{
		AudiobookID:    req.AudiobookID,
		UserID:         userID,
		EventType:      req.EventType,
		EventTimestamp: time.Now(),
	}

	if err := s.analyticsRepo.Create(&analytics); err != nil {
		return nil, err
	}

	return &dto.AnalyticsResponse{
		ID:             analytics.ID,
		AudiobookID:    analytics.AudiobookID,
		UserID:         analytics.UserID,
		EventType:      analytics.EventType,
		EventTimestamp: analytics.EventTimestamp,
	}, nil
}

// GetAnalyticsByID retrieves analytics by ID
func (s *AnalyticsService) GetAnalyticsByID(id uint) (*dto.AnalyticsResponse, error) {
	analytics, err := s.analyticsRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("analytics not found")
		}
		return nil, err
	}

	return &dto.AnalyticsResponse{
		ID:             analytics.ID,
		AudiobookID:    analytics.AudiobookID,
		UserID:         analytics.UserID,
		EventType:      analytics.EventType,
		EventTimestamp: analytics.EventTimestamp,
	}, nil
}

// GetAllAnalytics retrieves all analytics with pagination
func (s *AnalyticsService) GetAllAnalytics(req dto.PaginationRequest) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get paginated results
	analytics, total, err := s.analyticsRepo.GetAll(offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var analyticsResponses []dto.AnalyticsResponse
	for _, analytic := range analytics {
		analyticsResponses = append(analyticsResponses, dto.AnalyticsResponse{
			ID:             analytic.ID,
			AudiobookID:    analytic.AudiobookID,
			UserID:         analytic.UserID,
			EventType:      analytic.EventType,
			EventTimestamp: analytic.EventTimestamp,
		})
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: analyticsResponses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// DeleteAnalytics deletes analytics by ID
func (s *AnalyticsService) DeleteAnalytics(id uint) error {
	_, err := s.analyticsRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("analytics not found")
		}
		return err
	}

	return s.analyticsRepo.Delete(id)
}

// GetAnalyticsByUserID retrieves analytics by user ID
func (s *AnalyticsService) GetAnalyticsByUserID(userID string, req dto.PaginationRequest) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get paginated results
	analytics, total, err := s.analyticsRepo.GetByUserID(userID, offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var analyticsResponses []dto.AnalyticsResponse
	for _, analytic := range analytics {
		analyticsResponses = append(analyticsResponses, dto.AnalyticsResponse{
			ID:             analytic.ID,
			AudiobookID:    analytic.AudiobookID,
			UserID:         analytic.UserID,
			EventType:      analytic.EventType,
			EventTimestamp: analytic.EventTimestamp,
		})
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: analyticsResponses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetAnalyticsByAudiobookID retrieves analytics by audiobook ID
func (s *AnalyticsService) GetAnalyticsByAudiobookID(audiobookID uint, req dto.PaginationRequest) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get paginated results
	analytics, total, err := s.analyticsRepo.GetByAudiobookID(audiobookID, offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var analyticsResponses []dto.AnalyticsResponse
	for _, analytic := range analytics {
		analyticsResponses = append(analyticsResponses, dto.AnalyticsResponse{
			ID:             analytic.ID,
			AudiobookID:    analytic.AudiobookID,
			UserID:         analytic.UserID,
			EventType:      analytic.EventType,
			EventTimestamp: analytic.EventTimestamp,
		})
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: analyticsResponses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetAnalyticsSummary retrieves analytics summary for an audiobook
func (s *AnalyticsService) GetAnalyticsSummary(audiobookID uint) (map[string]int64, error) {
	return s.analyticsRepo.GetAnalyticsSummary(audiobookID)
}

// GetAnalyticsByDateRange retrieves analytics within a date range
func (s *AnalyticsService) GetAnalyticsByDateRange(startDate, endDate time.Time, req dto.PaginationRequest) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get paginated results
	analytics, total, err := s.analyticsRepo.GetByDateRange(startDate, endDate, offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var analyticsResponses []dto.AnalyticsResponse
	for _, analytic := range analytics {
		analyticsResponses = append(analyticsResponses, dto.AnalyticsResponse{
			ID:             analytic.ID,
			AudiobookID:    analytic.AudiobookID,
			UserID:         analytic.UserID,
			EventType:      analytic.EventType,
			EventTimestamp: analytic.EventTimestamp,
		})
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: analyticsResponses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// CreateAnalyticsEvent creates a new analytics event (alias for CreateAnalytics)
func (s *AnalyticsService) CreateAnalyticsEvent(userID string, req dto.CreateAnalyticsRequest) (*dto.AnalyticsResponse, error) {
	return s.CreateAnalytics(userID, req)
}

// GetAnalyticsByUser retrieves analytics by user ID (alias for GetAnalyticsByUserID)
func (s *AnalyticsService) GetAnalyticsByUser(userID string, req dto.PaginationRequest) (*dto.ListResponse, error) {
	return s.GetAnalyticsByUserID(userID, req)
}

// GetAnalyticsByAudiobook retrieves analytics by audiobook ID (alias for GetAnalyticsByAudiobookID)
func (s *AnalyticsService) GetAnalyticsByAudiobook(audiobookID uint, req dto.PaginationRequest) (*dto.ListResponse, error) {
	return s.GetAnalyticsByAudiobookID(audiobookID, req)
}

// GetAnalyticsByEventType retrieves analytics by event type
func (s *AnalyticsService) GetAnalyticsByEventType(eventType string, req dto.PaginationRequest) (*dto.ListResponse, error) {
	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get paginated results
	analytics, total, err := s.analyticsRepo.GetByEventType(eventType, offset, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var analyticsResponses []dto.AnalyticsResponse
	for _, analytic := range analytics {
		analyticsResponses = append(analyticsResponses, dto.AnalyticsResponse{
			ID:             analytic.ID,
			AudiobookID:    analytic.AudiobookID,
			UserID:         analytic.UserID,
			EventType:      analytic.EventType,
			EventTimestamp: analytic.EventTimestamp,
		})
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &dto.ListResponse{
		Items: analyticsResponses,
		Pagination: dto.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}
