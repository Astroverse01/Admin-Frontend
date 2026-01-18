package handlers

import (
	"astro-admin/internal/dto"
	"astro-admin/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AstroHandler struct {
	astroService *services.AstroService
}

func NewAstroHandler(astroService *services.AstroService) *AstroHandler {
	return &AstroHandler{
		astroService: astroService,
	}
}

func (h *AstroHandler) ListAstros(c *gin.Context) {
	// Parse query parameters
	name := c.Query("name")
	sort := c.Query("sort")
	if sort == "" {
		sort = "asc"
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	// Call service
	response, err := h.astroService.ListAstros(c.Request.Context(), name, sort, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to fetch astrologers: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AstroHandler) UpdateAstroStatus(c *gin.Context) {
	astroID := c.Param("astroId")
	if astroID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Astrologer ID is required",
		})
		return
	}

	var req dto.StatusUpdateRequest
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
	err := h.astroService.UpdateAstroStatus(c.Request.Context(), astroID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to update astrologer status: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.StatusUpdateResponse{
		Success: true,
		Message: "Astrologer " + req.Status + "d successfully",
		ID:      astroID,
	})
}

func (h *AstroHandler) ToggleVisibility(c *gin.Context) {
	astroID := c.Param("astroId")
	if astroID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Astrologer ID is required",
		})
		return
	}

	var req dto.VisibilityUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// No validation needed for boolean field - it can only be true or false

	// Call service
	err := h.astroService.ToggleVisibility(c.Request.Context(), astroID, req.Visible)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to update astrologer visibility: " + err.Error(),
		})
		return
	}

	visibility := "hidden"
	if req.Visible {
		visibility = "visible"
	}

	c.JSON(http.StatusOK, dto.StatusUpdateResponse{
		Success: true,
		Message: "Astrologer is now " + visibility + " to users",
		ID:      astroID,
	})
}