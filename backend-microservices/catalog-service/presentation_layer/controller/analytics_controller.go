package controller

import (
	"catalog-service/data_layer/dto"
	"catalog-service/domain_layer/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AnalyticsController struct {
	analyticsService *service.AnalyticsService
}

func NewAnalyticsController(analyticsService *service.AnalyticsService) *AnalyticsController {
	return &AnalyticsController{
		analyticsService: analyticsService,
	}
}

// CreateAnalyticsEvent creates a new analytics event
func (ac *AnalyticsController) CreateAnalyticsEvent(c *gin.Context) {
	var req dto.CreateAnalyticsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract user ID from context (in a real app, this would come from JWT/auth middleware)
	userID := c.GetString("user_id")
	if userID == "" {
		userID = "anonymous" // Default for now
	}

	analytics, err := ac.analyticsService.CreateAnalyticsEvent(userID, req)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "audiobook not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, analytics)
}

// GetAnalyticsByID retrieves analytics by ID
func (ac *AnalyticsController) GetAnalyticsByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid analytics ID"})
		return
	}

	analytics, err := ac.analyticsService.GetAnalyticsByID(uint(id))
	if err != nil {
		if err.Error() == "analytics not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetAnalyticsByDateRange retrieves analytics within a date range
func (ac *AnalyticsController) GetAnalyticsByDateRange(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date parameters are required (format: YYYY-MM-DD)"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	analytics, err := ac.analyticsService.GetAnalyticsByDateRange(startDate, endDate, dto.PaginationRequest{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetAnalyticsByUser retrieves analytics for a specific user
func (ac *AnalyticsController) GetAnalyticsByUser(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	paginationReq := dto.PaginationRequest{
		Page:  page,
		Limit: limit,
	}

	analytics, err := ac.analyticsService.GetAnalyticsByUser(c.Param("user_id"), paginationReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetAnalyticsByAudiobook retrieves analytics for a specific audiobook
func (ac *AnalyticsController) GetAnalyticsByAudiobook(c *gin.Context) {
	audiobookID, err := strconv.ParseUint(c.Param("audiobook_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audiobook ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	paginationReq := dto.PaginationRequest{
		Page:  page,
		Limit: limit,
	}

	analytics, err := ac.analyticsService.GetAnalyticsByAudiobook(uint(audiobookID), paginationReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetAnalyticsByEventType retrieves analytics by event type
func (ac *AnalyticsController) GetAnalyticsByEventType(c *gin.Context) {
	eventType := c.Param("event_type")
	if eventType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event type parameter is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	paginationReq := dto.PaginationRequest{
		Page:  page,
		Limit: limit,
	}

	analytics, err := ac.analyticsService.GetAnalyticsByEventType(eventType, paginationReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetAnalyticsSummary retrieves analytics summary for an audiobook
func (ac *AnalyticsController) GetAnalyticsSummary(c *gin.Context) {
	audiobookID, err := strconv.ParseUint(c.Param("audiobook_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audiobook ID"})
		return
	}

	summary, err := ac.analyticsService.GetAnalyticsSummary(uint(audiobookID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// DeleteAnalytics deletes analytics record
func (ac *AnalyticsController) DeleteAnalytics(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid analytics ID"})
		return
	}

	if err := ac.analyticsService.DeleteAnalytics(uint(id)); err != nil {
		if err.Error() == "analytics not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Analytics deleted successfully"})
}
