package service

import (
	"context"
	"encoding/json"
	"microservice/user/data-layer/entity"
	"microservice/user/data-layer/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// LogActivityRequest represents a request to log an activity
type LogActivityRequest struct {
	UserID     string      `json:"user_id"`
	TenantID   string      `json:"tenant_id,omitempty"`
	EntityID   string      `json:"entity_id"`
	EntityType string      `json:"entity_type"`
	Action     string      `json:"action"`
	Detail     string      `json:"detail,omitempty"`
	OldValue   interface{} `json:"old_value,omitempty"`
	NewValue   interface{} `json:"new_value,omitempty"`
	IPAddress  string      `json:"ip_address,omitempty"`
	UserAgent  string      `json:"user_agent,omitempty"`
}

// AuditService is the interface for audit logging
type AuditService interface {
	LogActivity(ctx context.Context, req *LogActivityRequest) error
	GetLogs(ctx context.Context, page, limit int, entityType, entityID, userID, tenantID string) ([]entity.AuditLog, int, error)
	GetLogByID(ctx context.Context, id string) (entity.AuditLog, error)
	GetLogsWithDateRange(ctx context.Context, page, limit int, entityType, entityID, userID, tenantID string, startDate, endDate *time.Time) ([]entity.AuditLog, int, error)
	GetStatistics(ctx context.Context, period, entityType string) (map[string]interface{}, error)
}

type auditService struct {
	repo repository.AuditLogRepository
}

// NewAuditService creates a new audit service
func NewAuditService(repo repository.AuditLogRepository) AuditService {
	return &auditService{
		repo: repo,
	}
}

// LogActivity logs an activity
func (s *auditService) LogActivity(ctx context.Context, req *LogActivityRequest) error {
	// Convert old and new values to JSON if provided
	var oldValueJSON, newValueJSON datatypes.JSON
	var err error

	if req.OldValue != nil {
		tmpJSON, err := json.Marshal(req.OldValue)
		if err != nil {
			return err
		}
		oldValueJSON = datatypes.JSON(tmpJSON)
	}

	if req.NewValue != nil {
		tmpJSON, err := json.Marshal(req.NewValue)
		if err != nil {
			return err
		}
		newValueJSON = datatypes.JSON(tmpJSON)
	}

	// Parse UUIDs
	var userUUID *uuid.UUID
	var tenantUUID, entityUUID uuid.UUID

	if req.UserID != "" {
		parsed, err := uuid.Parse(req.UserID)
		if err != nil {
			return err
		}
		userUUID = &parsed
	}

	if req.TenantID != "" {
		tenantUUID, err = uuid.Parse(req.TenantID)
		if err != nil {
			return err
		}
	}

	if req.EntityID != "" {
		entityUUID, err = uuid.Parse(req.EntityID)
		if err != nil {
			return err
		}
	}

	// Create audit log entry
	log := entity.AuditLog{
		ID:         uuid.New(),
		UserID:     userUUID,
		TenantID:   tenantUUID,
		EntityID:   entityUUID,
		EntityType: req.EntityType,
		Action:     req.Action,
		OldValues:  oldValueJSON,
		NewValues:  newValueJSON,
		IPAddress:  req.IPAddress,
		UserAgent:  req.UserAgent,
		CreatedAt:  time.Now(),
	}

	return s.repo.Create(ctx, nil, log)
}

// GetLogs retrieves audit logs with filtering and pagination
func (s *auditService) GetLogs(ctx context.Context, page, limit int, entityType, entityID, userID, tenantID string) ([]entity.AuditLog, int, error) {
	// Use repository methods to implement the search logic
	var logs []entity.AuditLog
	var total int

	// Try to filter by different criteria
	if entityID != "" && entityType != "" {
		foundLogs, err := s.repo.FindByEntityTypeAndID(ctx, entityType, entityID)
		if err != nil {
			return nil, 0, err
		}
		logs = foundLogs
		total = len(logs)
	} else if userID != "" {
		foundLogs, err := s.repo.FindByUserID(ctx, userID)
		if err != nil {
			return nil, 0, err
		}
		logs = foundLogs
		total = len(logs)
	} else if tenantID != "" {
		foundLogs, err := s.repo.FindByTenantID(ctx, tenantID)
		if err != nil {
			return nil, 0, err
		}
		logs = foundLogs
		total = len(logs)
	} else {
		// If no filters provided, get all logs with pagination (for superadmin)
		offset := (page - 1) * limit
		foundLogs, totalCount, err := s.repo.FindAll(ctx, limit, offset)
		if err != nil {
			return nil, 0, err
		}
		return foundLogs, totalCount, nil
	}

	// Apply pagination
	startIndex := (page - 1) * limit
	endIndex := startIndex + limit
	if startIndex >= len(logs) {
		return []entity.AuditLog{}, total, nil
	}
	if endIndex > len(logs) {
		endIndex = len(logs)
	}

	return logs[startIndex:endIndex], total, nil
}

// GetLogByID retrieves a single audit log by ID
func (s *auditService) GetLogByID(ctx context.Context, id string) (entity.AuditLog, error) {
	return s.repo.FindByID(ctx, id)
}

// GetLogsWithDateRange retrieves audit logs with filtering by date range and pagination
func (s *auditService) GetLogsWithDateRange(ctx context.Context, page, limit int, entityType, entityID, userID, tenantID string, startDate, endDate *time.Time) ([]entity.AuditLog, int, error) {
	// Use repository methods to implement the search logic with date range
	var logs []entity.AuditLog
	var total int

	// In a real implementation, you would have a repository method that supports date filtering
	// For this example, we'll use the existing methods and then filter by date in memory
	if entityID != "" && entityType != "" {
		foundLogs, err := s.repo.FindByEntityTypeAndID(ctx, entityType, entityID)
		if err != nil {
			return nil, 0, err
		}
		logs = foundLogs
	} else if userID != "" {
		foundLogs, err := s.repo.FindByUserID(ctx, userID)
		if err != nil {
			return nil, 0, err
		}
		logs = foundLogs
	} else if tenantID != "" {
		foundLogs, err := s.repo.FindByTenantID(ctx, tenantID)
		if err != nil {
			return nil, 0, err
		}
		logs = foundLogs
	} else {
		// If no filters provided, get all logs (for superadmin)
		// We'll get a larger set and filter by date range, then paginate
		offset := 0
		limit := 10000 // Get more records for date filtering
		foundLogs, _, err := s.repo.FindAll(ctx, limit, offset)
		if err != nil {
			return nil, 0, err
		}
		logs = foundLogs
	}

	// Filter by date range in memory
	if startDate != nil || endDate != nil {
		filteredLogs := make([]entity.AuditLog, 0)
		for _, log := range logs {
			if startDate != nil && log.CreatedAt.Before(*startDate) {
				continue
			}
			if endDate != nil && log.CreatedAt.After(*endDate) {
				continue
			}
			filteredLogs = append(filteredLogs, log)
		}
		logs = filteredLogs
	}

	total = len(logs)

	// Apply pagination
	startIndex := (page - 1) * limit
	endIndex := startIndex + limit
	if startIndex >= len(logs) {
		return []entity.AuditLog{}, total, nil
	}
	if endIndex > len(logs) {
		endIndex = len(logs)
	}

	return logs[startIndex:endIndex], total, nil
}

// GetStatistics retrieves statistics about audit logs
func (s *auditService) GetStatistics(ctx context.Context, period, entityType string) (map[string]interface{}, error) {
	// In a real implementation, this would query the database for aggregate data
	// For this example, we'll implement a simple in-memory calculation

	if entityType != "" {
		// If entity type is specified, get logs for that entity type
		// This is an example - we would need additional repository methods for this
		return map[string]interface{}{
			"period":      period,
			"entity_type": entityType,
			"counts": map[string]int{
				"create": 10, // These would be actual counts from the database
				"update": 25,
				"delete": 5,
			},
			"users": map[string]int{
				"admin": 15,
				"user1": 10,
				"user2": 15,
			},
		}, nil
	}

	// Mock statistics based on period
	switch period {
	case "weekly":
		return map[string]interface{}{
			"period": "weekly",
			"days":   []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
			"counts": []int{12, 15, 8, 10, 5, 2, 3}, // Example counts per day
			"action_types": map[string]int{
				"create": 25,
				"update": 20,
				"delete": 10,
			},
		}, nil
	case "monthly":
		return map[string]interface{}{
			"period": "monthly",
			"weeks":  []string{"Week 1", "Week 2", "Week 3", "Week 4"},
			"counts": []int{42, 38, 56, 35}, // Example counts per week
			"action_types": map[string]int{
				"create": 70,
				"update": 65,
				"delete": 36,
			},
		}, nil
	default: // daily
		return map[string]interface{}{
			"period": "daily",
			"hours":  []string{"00:00", "04:00", "08:00", "12:00", "16:00", "20:00"},
			"counts": []int{5, 2, 15, 25, 18, 10}, // Example counts per time period
			"action_types": map[string]int{
				"create": 30,
				"update": 35,
				"delete": 10,
			},
		}, nil
	}
}
