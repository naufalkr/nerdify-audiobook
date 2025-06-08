package controller

import (
	"catalog-service/data_layer/dto"
	"catalog-service/domain_layer/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TrackController struct {
	trackService *service.TrackService
}

func NewTrackController(trackService *service.TrackService) *TrackController {
	return &TrackController{
		trackService: trackService,
	}
}

// CreateTrack creates a new track
func (tc *TrackController) CreateTrack(c *gin.Context) {
	var req dto.CreateTrackRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	track, err := tc.trackService.CreateTrack(req)
	if err != nil {
		if err.Error() == "audiobook not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, track)
}

// GetTrackByID retrieves a track by ID
func (tc *TrackController) GetTrackByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	track, err := tc.trackService.GetTrackByID(uint(id))
	if err != nil {
		if err.Error() == "track not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, track)
}

// GetAllTracks retrieves all tracks with pagination
func (tc *TrackController) GetAllTracks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	req := dto.PaginationRequest{
		Page:  page,
		Limit: limit,
	}

	tracks, err := tc.trackService.GetAllTracks(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

// GetTracksByAudiobook retrieves all tracks for a specific audiobook
func (tc *TrackController) GetTracksByAudiobook(c *gin.Context) {
	audiobookID, err := strconv.ParseUint(c.Param("audiobook_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audiobook ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	tracks, err := tc.trackService.GetTracksByAudiobook(uint(audiobookID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

// UpdateTrack updates an existing track
func (tc *TrackController) UpdateTrack(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	var req dto.UpdateTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	track, err := tc.trackService.UpdateTrack(uint(id), req)
	if err != nil {
		if err.Error() == "track not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "audiobook not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, track)
}

// DeleteTrack deletes a track
func (tc *TrackController) DeleteTrack(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	if err := tc.trackService.DeleteTrack(uint(id)); err != nil {
		if err.Error() == "track not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Track deleted successfully"})
}

// UpdateTrackOrder updates the order of tracks within an audiobook
func (tc *TrackController) UpdateTrackOrder(c *gin.Context) {
	audiobookID, err := strconv.ParseUint(c.Param("audiobook_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid audiobook ID"})
		return
	}

	var req struct {
		TrackOrders []struct {
			TrackID uint `json:"track_id" binding:"required"`
			Order   int  `json:"order" binding:"required"`
		} `json:"track_orders" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert request to service format
	var trackOrders []struct {
		TrackID uint `json:"track_id"`
		Order   int  `json:"order"`
	}

	for _, trackOrder := range req.TrackOrders {
		trackOrders = append(trackOrders, struct {
			TrackID uint `json:"track_id"`
			Order   int  `json:"order"`
		}{
			TrackID: trackOrder.TrackID,
			Order:   trackOrder.Order,
		})
	}

	if err := tc.trackService.UpdateTrackOrder(uint(audiobookID), trackOrders); err != nil {
		if err.Error() == "audiobook not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Track order updated successfully"})
}
