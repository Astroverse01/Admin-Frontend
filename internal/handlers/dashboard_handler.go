package handlers

import (
	"admin-be/internal/dto"
	"admin-be/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService *services.DashboardService
}

func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetDailyMetrics handles GET /admin/dashboard/metrics
// Query parameter: date (optional, format: YYYY-MM-DD, defaults to today)
func (h *DashboardHandler) GetDailyMetrics(c *gin.Context) {
	// Get date from query parameter, default to today
	dateStr := c.DefaultQuery("date", "")
	if dateStr == "" {
		// Default to today's date in YYYY-MM-DD format
		dateStr = time.Now().UTC().Format("2006-01-02")
	}

	// Validate date format
	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid date format. Use YYYY-MM-DD (e.g., 2024-01-15)",
		})
		return
	}

	// Call service
	response, err := h.dashboardService.GetDailyMetrics(c.Request.Context(), dateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to fetch dashboard metrics: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

