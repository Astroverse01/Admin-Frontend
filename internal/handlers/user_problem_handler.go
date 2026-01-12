package handlers

import (
	"astro-admin/internal/dto"
	"astro-admin/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserProblemHandler struct {
	userProblemService *services.UserProblemService
}

func NewUserProblemHandler(userProblemService *services.UserProblemService) *UserProblemHandler {
	return &UserProblemHandler{
		userProblemService: userProblemService,
	}
}

func (h *UserProblemHandler) ListUserGeneralComplaints(c *gin.Context) {
	// Parse query parameters
	problemTypes := c.Query("problemTypes")
	status := c.Query("status")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	// Call service
	response, err := h.userProblemService.ListUserGeneralComplaints(c.Request.Context(), problemTypes, status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to fetch complaints: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserProblemHandler) CloseComplaint(c *gin.Context) {
	problemID := c.Param("problemId")
	if problemID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Problem ID is required",
		})
		return
	}

	var req dto.CloseComplaintRequest
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
	response, err := h.userProblemService.CloseComplaint(c.Request.Context(), problemID, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to close complaint: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

