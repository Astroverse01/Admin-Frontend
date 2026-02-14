package handlers

import (
	"admin-be/internal/dto"
	"admin-be/internal/services"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	log.Println("[ListUsers Handler] Request received")

	// Parse query parameters
	name := c.Query("name")
	sort := c.Query("sort")
	if sort == "" {
		sort = "asc"
	}
	updatedOnFrom := c.Query("updatedOnFrom") // optional, YYYY-MM-DD
	updatedOnTo := c.Query("updatedOnTo")     // optional, YYYY-MM-DD

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		log.Printf("[ListUsers Handler] Invalid page parameter: %v, defaulting to 1", err)
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		log.Printf("[ListUsers Handler] Invalid limit parameter: %v, defaulting to 10", err)
		limit = 10
	}

	// Validate updatedOn date format if provided
	if updatedOnFrom != "" {
		if _, err := time.Parse("2006-01-02", updatedOnFrom); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Message: "Invalid updatedOnFrom format. Use YYYY-MM-DD (e.g., 2025-01-15)",
			})
			return
		}
	}
	if updatedOnTo != "" {
		if _, err := time.Parse("2006-01-02", updatedOnTo); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Message: "Invalid updatedOnTo format. Use YYYY-MM-DD (e.g., 2025-01-15)",
			})
			return
		}
	}

	log.Printf("[ListUsers Handler] Query parameters - name: %s, sort: %s, page: %d, limit: %d, updatedOnFrom: %s, updatedOnTo: %s", name, sort, page, limit, updatedOnFrom, updatedOnTo)

	// Call service
	log.Println("[ListUsers Handler] Calling userService.ListUsers")
	response, err := h.userService.ListUsers(c.Request.Context(), name, sort, page, limit, updatedOnFrom, updatedOnTo)
	if err != nil {
		log.Printf("[ListUsers Handler] Error from service: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to fetch users: " + err.Error(),
		})
		return
	}

	log.Printf("[ListUsers Handler] Successfully fetched users. Total: %d, Returned: %d", response.Pagination.Total, len(response.Data))
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) DeactivateUser(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "User ID is required",
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
	err := h.userService.DeactivateUser(c.Request.Context(), userID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to update user status: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.StatusUpdateResponse{
		Success: true,
		Message: "User " + req.Status + "d successfully",
		ID:      userID,
	})
}
