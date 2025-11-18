package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mukund/mediaconvert/internal/analytics"
	"github.com/mukund/mediaconvert/internal/auth"
)

type AnalyticsHandler struct {
	analytics *analytics.Client
}

func NewAnalyticsHandler(analyticsClient *analytics.Client) *AnalyticsHandler {
	return &AnalyticsHandler{
		analytics: analyticsClient,
	}
}

// GetJobStats returns aggregated job statistics
func (h *AnalyticsHandler) GetJobStats(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if h.analytics == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Analytics not available"})
		return
	}

	// Parse time range (default: last 7 days)
	days := c.DefaultQuery("days", "7")
	daysInt := 7
	if d, err := time.ParseDuration(days + "d"); err == nil {
		daysInt = int(d.Hours() / 24)
	}

	stats, err := h.analytics.GetJobStats(c.Request.Context(), uint64(userID), daysInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetJobTimeline returns job metrics over time
func (h *AnalyticsHandler) GetJobTimeline(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if h.analytics == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Analytics not available"})
		return
	}

	// Parse time range
	days := c.DefaultQuery("days", "7")
	daysInt := 7
	if d, err := time.ParseDuration(days + "d"); err == nil {
		daysInt = int(d.Hours() / 24)
	}

	// Parse interval (hour, day)
	interval := c.DefaultQuery("interval", "hour")

	timeline, err := h.analytics.GetJobTimeline(c.Request.Context(), uint64(userID), daysInt, interval)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch timeline"})
		return
	}

	c.JSON(http.StatusOK, timeline)
}

// GetPipelineStats returns statistics by pipeline
func (h *AnalyticsHandler) GetPipelineStats(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if h.analytics == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Analytics not available"})
		return
	}

	days := c.DefaultQuery("days", "30")
	daysInt := 30
	if d, err := time.ParseDuration(days + "d"); err == nil {
		daysInt = int(d.Hours() / 24)
	}

	stats, err := h.analytics.GetPipelineStats(c.Request.Context(), uint64(userID), daysInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pipeline statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
