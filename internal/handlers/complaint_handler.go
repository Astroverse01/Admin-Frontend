package handlers

import (
	"astro-admin/internal/dto"
	"astro-admin/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ComplaintHandler struct {
	complaintService *services.ComplaintService
}

func NewComplaintHandler(complaintService *services.ComplaintService) *ComplaintHandler {
	return &ComplaintHandler{
		complaintService: complaintService,
	}
}

func (h *ComplaintHandler) ListUserServiceComplaints(c *gin.Context) {
	// Parse query parameters
	serviceType := c.Query("serviceType")
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
	response, err := h.complaintService.ListUserServiceComplaints(c.Request.Context(), serviceType, status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to fetch complaints: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ComplaintHandler) GetComplaintDetails(c *gin.Context) {
	serviceType := c.Param("serviceType")
	orderId := c.Param("orderId")

	if serviceType == "" || orderId == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Service type and order ID are required",
		})
		return
	}

	// Call service
	response, err := h.complaintService.GetComplaintDetails(c.Request.Context(), serviceType, orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to fetch complaint details: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ComplaintHandler) AcceptRejectComplaint(c *gin.Context) {
	orderId := c.Param("orderId")
	if orderId == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Order ID is required",
		})
		return
	}

	var req dto.AcceptRejectRequest
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
	response, err := h.complaintService.AcceptRejectComplaint(c.Request.Context(), orderId, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to process complaint: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
