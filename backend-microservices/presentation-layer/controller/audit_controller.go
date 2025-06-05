package controller

import (
	"net/http"
	"strconv"
	"time"

	"microservice/user/data-layer/entity"
	"microservice/user/domain-layer/service"

	"github.com/gin-gonic/gin"
)

// Fungsi helper untuk memparsing tanggal dengan berbagai format
func parseDateFlexibly(dateStr, primaryFormat string) (time.Time, error) {
	// Pertama, coba dengan format yang ditentukan pengguna
	if t, err := time.Parse(primaryFormat, dateStr); err == nil {
		return t, nil
	}

	// Coba format-format umum jika format utama gagal
	commonFormats := []string{
		"2006-01-02",
		"02/01/2006",
		"01/02/2006",
		"2006/01/02",
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"02-01-2006",
		"01-02-2006",
	}

	for _, format := range commonFormats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	// Jika semua format gagal, kembalikan error dari format utama
	return time.Parse(primaryFormat, dateStr)
}

type AuditController struct {
	auditService service.AuditService
}

func NewAuditController(auditService service.AuditService) *AuditController {
	return &AuditController{
		auditService: auditService,
	}
}

// GetAuditLogs handles requests to get audit logs with pagination
func (c *AuditController) GetAuditLogs(ctx *gin.Context) {
	// Parse pagination parameters
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Get filters from query parameters
	entityType := ctx.Query("entity_type")
	entityID := ctx.Query("entity_id")
	userID := ctx.Query("user_id")
	tenantID := ctx.Query("tenant_id")

	// Get date range filters
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")
	dateFormat := ctx.DefaultQuery("date_format", "2006-01-02") // Allow custom date format

	var startDate, endDate *time.Time
	if startDateStr != "" {
		parsedStartDate, err := parseDateFlexibly(startDateStr, dateFormat)
		if err == nil {
			// Set to start of day (00:00:00)
			year, month, day := parsedStartDate.Date()
			startOfDay := time.Date(year, month, day, 0, 0, 0, 0, parsedStartDate.Location())
			startDate = &startOfDay
		}
	}

	if endDateStr != "" {
		parsedEndDate, err := parseDateFlexibly(endDateStr, dateFormat)
		if err == nil {
			// Set to end of day (23:59:59)
			year, month, day := parsedEndDate.Date()
			endOfDay := time.Date(year, month, day, 23, 59, 59, 999999999, parsedEndDate.Location())
			endDate = &endOfDay
		}
	}

	// Get logs based on filters with date range
	var logs []entity.AuditLog
	var total int

	// Use date range if provided, otherwise use standard method
	if startDate != nil || endDate != nil {
		logs, total, err = c.auditService.GetLogsWithDateRange(ctx.Request.Context(), page, limit, entityType, entityID, userID, tenantID, startDate, endDate)
	} else {
		logs, total, err = c.auditService.GetLogs(ctx.Request.Context(), page, limit, entityType, entityID, userID, tenantID)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": logs,
		"meta": gin.H{
			"page":      page,
			"limit":     limit,
			"total":     total,
			"last_page": (total + limit - 1) / limit,
		},
	})
}

// GetAuditLogByID handles requests to get a specific audit log by ID
func (c *AuditController) GetAuditLogByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Audit log ID is required"})
		return
	}

	log, err := c.auditService.GetLogByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": log})
}

// ExportAuditLogs handles requests to export audit logs in different formats
func (c *AuditController) ExportAuditLogs(ctx *gin.Context) {
	// Parse pagination parameters
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "100"))
	if err != nil || limit < 1 || limit > 1000 {
		limit = 100
	}

	// Get filters from query parameters
	entityType := ctx.Query("entity_type")
	entityID := ctx.Query("entity_id")
	userID := ctx.Query("user_id")
	tenantID := ctx.Query("tenant_id")
	format := ctx.DefaultQuery("format", "json")

	// Get date range filters
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")
	dateFormat := ctx.DefaultQuery("date_format", "2006-01-02") // Allow custom date format

	var startDate, endDate *time.Time
	if startDateStr != "" {
		// Try different date formats if the specified format fails
		parsedStartDate, err := parseDateFlexibly(startDateStr, dateFormat)
		if err == nil {
			// Set to start of day (00:00:00)
			year, month, day := parsedStartDate.Date()
			startOfDay := time.Date(year, month, day, 0, 0, 0, 0, parsedStartDate.Location())
			startDate = &startOfDay
		}
	}

	if endDateStr != "" {
		// Try different date formats if the specified format fails
		parsedEndDate, err := parseDateFlexibly(endDateStr, dateFormat)
		if err == nil {
			// Set to end of day (23:59:59)
			year, month, day := parsedEndDate.Date()
			endOfDay := time.Date(year, month, day, 23, 59, 59, 999999999, parsedEndDate.Location())
			endDate = &endOfDay
		}
	}

	// Get logs based on filters with date range
	logs, total, err := c.auditService.GetLogsWithDateRange(ctx.Request.Context(), page, limit, entityType, entityID, userID, tenantID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Handle different export formats
	switch format {
	case "csv":
		ctx.Header("Content-Disposition", "attachment; filename=audit_logs.csv")
		ctx.Header("Content-Type", "text/csv")
		// In a real implementation, you would convert logs to CSV format here
		// For simplicity, we'll just return JSON for now
		ctx.JSON(http.StatusOK, logs)
	case "pdf":
		ctx.Header("Content-Disposition", "attachment; filename=audit_logs.pdf")
		ctx.Header("Content-Type", "application/pdf")
		// In a real implementation, you would convert logs to PDF format here
		// For simplicity, we'll just return JSON for now
		ctx.JSON(http.StatusOK, logs)
	default: // json
		ctx.JSON(http.StatusOK, gin.H{
			"data": logs,
			"meta": gin.H{
				"page":      page,
				"limit":     limit,
				"total":     total,
				"last_page": (total + limit - 1) / limit,
			},
		})
	}
}

// GetAuditStatistics handles requests to get statistics about audit logs
func (c *AuditController) GetAuditStatistics(ctx *gin.Context) {
	// Get time period from query parameters (daily, weekly, monthly)
	period := ctx.DefaultQuery("period", "daily")

	// Get entity type filter if provided
	entityType := ctx.Query("entity_type")

	// Call service method to get statistics
	stats, err := c.auditService.GetStatistics(ctx.Request.Context(), period, entityType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": stats})
}

// GetUserAuditLogs handles requests to get audit logs for the authenticated user
func (c *AuditController) GetUserAuditLogs(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Convert userID to string
	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Parse pagination parameters
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Get filters from query parameters
	entityType := ctx.Query("entity_type")
	entityID := ctx.Query("entity_id")
	tenantID := ctx.Query("tenant_id")

	// Get date range filters
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")
	dateFormat := ctx.DefaultQuery("date_format", "2006-01-02") // Allow custom date format

	var startDate, endDate *time.Time
	if startDateStr != "" {
		parsedStartDate, err := parseDateFlexibly(startDateStr, dateFormat)
		if err == nil {
			// Set to start of day (00:00:00)
			year, month, day := parsedStartDate.Date()
			startOfDay := time.Date(year, month, day, 0, 0, 0, 0, parsedStartDate.Location())
			startDate = &startOfDay
		}
	}

	if endDateStr != "" {
		parsedEndDate, err := parseDateFlexibly(endDateStr, dateFormat)
		if err == nil {
			// Set to end of day (23:59:59)
			year, month, day := parsedEndDate.Date()
			endOfDay := time.Date(year, month, day, 23, 59, 59, 999999999, parsedEndDate.Location())
			endDate = &endOfDay
		}
	}

	// Get logs based on filters with date range, force userID to be the authenticated user
	var logs []entity.AuditLog
	var total int

	// Use date range if provided, otherwise use standard method
	if startDate != nil || endDate != nil {
		logs, total, err = c.auditService.GetLogsWithDateRange(ctx.Request.Context(), page, limit, entityType, entityID, userIDStr, tenantID, startDate, endDate)
	} else {
		logs, total, err = c.auditService.GetLogs(ctx.Request.Context(), page, limit, entityType, entityID, userIDStr, tenantID)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": logs,
		"meta": gin.H{
			"page":      page,
			"limit":     limit,
			"total":     total,
			"last_page": (total + limit - 1) / limit,
		},
	})
}
