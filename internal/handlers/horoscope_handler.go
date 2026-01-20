package handlers

import (
	"admin-be/internal/dto"
	"admin-be/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type HoroscopeHandler struct {
	horoscopeService *services.HoroscopeService
	adminID          string
}

func NewHoroscopeHandler(horoscopeService *services.HoroscopeService, adminID string) *HoroscopeHandler {
	return &HoroscopeHandler{
		horoscopeService: horoscopeService,
		adminID:          adminID,
	}
}

func (h *HoroscopeHandler) ListHoroscopes(c *gin.Context) {
	// Parse query parameters
	date := c.Query("date")
	signName := c.Query("signName")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "12"))
	if err != nil || limit < 1 {
		limit = 12
	}

	// Call service
	response, err := h.horoscopeService.ListHoroscopes(c.Request.Context(), date, signName, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to fetch horoscopes: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *HoroscopeHandler) BulkCreateHoroscopes(c *gin.Context) {
	var payload dto.BulkHoroscopeRequestPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if len(payload) == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "No horoscope data provided",
		})
		return
	}

	response, err := h.horoscopeService.BulkCreateHoroscopes(c.Request.Context(), payload, h.adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to process horoscopes: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *HoroscopeHandler) UpdateHoroscope(c *gin.Context) {
	horoscopeID := c.Param("horoscopeId")
	if horoscopeID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Horoscope ID is required",
		})
		return
	}

	var req dto.UpdateHoroscopeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// Validate request
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Validation failed: " + err.Error(),
		})
		return
	}

	// Call service
	response, err := h.horoscopeService.UpdateHoroscope(c.Request.Context(), horoscopeID, &req, h.adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to update horoscope: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *HoroscopeHandler) DeleteHoroscope(c *gin.Context) {
	horoscopeID := c.Param("horoscopeId")
	if horoscopeID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Horoscope ID is required",
		})
		return
	}

	// Call service
	response, err := h.horoscopeService.DeleteHoroscope(c.Request.Context(), horoscopeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to delete horoscope: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
