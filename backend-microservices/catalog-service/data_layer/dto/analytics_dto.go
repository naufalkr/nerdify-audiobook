package dto

import "time"

// CreateAnalyticsRequest represents the request to create analytics event
type CreateAnalyticsRequest struct {
	AudiobookID uint   `json:"audiobook_id" binding:"required"`
	EventType   string `json:"event_type" binding:"required,oneof=VIEW PLAY_START PLAY_FINISH DOWNLOAD"`
}

// AnalyticsResponse represents the response for analytics data
type AnalyticsResponse struct {
	ID             uint      `json:"id"`
	AudiobookID    uint      `json:"audiobook_id"`
	UserID         string    `json:"user_id"`
	EventType      string    `json:"event_type"`
	EventTimestamp time.Time `json:"event_timestamp"`
}

// AnalyticsStatsResponse represents analytics statistics
type AnalyticsStatsResponse struct {
	AudiobookID    uint                `json:"audiobook_id"`
	AudiobookTitle string              `json:"audiobook_title"`
	TotalViews     int64               `json:"total_views"`
	TotalPlays     int64               `json:"total_plays"`
	EventsByType   map[string]int64    `json:"events_by_type"`
	RecentEvents   []AnalyticsResponse `json:"recent_events"`
}
