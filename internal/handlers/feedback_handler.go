package handlers

import (
	"admin-be/internal/dto"
	"admin-be/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FeedbackHandler struct {
	feedbackService *services.FeedbackService
}

func NewFeedbackHandler(feedbackService *services.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{
		feedbackService: feedbackService,
	}
}

func (h *FeedbackHandler) BulkCreateFeedbacks(c *gin.Context) {
	var payload dto.BulkFeedbackRequestPayload
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
			Message: "No feedback data provided",
		})
		return
	}

	// Validate each item in the payload
	validate := validator.New()
	for i, item := range payload {
		if err := validate.Struct(item); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Message: "Validation failed for item " + strconv.Itoa(i+1) + ": " + err.Error(),
			})
			return
		}
	}

	response, err := h.feedbackService.BulkCreateFeedbacks(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to create feedbacks: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

