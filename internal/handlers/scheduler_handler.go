package handlers

import (
	"admin-be/internal/dto"
	"admin-be/internal/services"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type SchedulerHandler struct {
	schedulerService *services.SchedulerService
}

func NewSchedulerHandler(schedulerService *services.SchedulerService) *SchedulerHandler {
	return &SchedulerHandler{
		schedulerService: schedulerService,
	}
}

// TriggerManually manually triggers the daily report generation (for testing)
func (h *SchedulerHandler) TriggerManually(c *gin.Context) {
	if err := h.schedulerService.RunManually(); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to trigger scheduler: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Daily report generation triggered successfully",
	})
}

// GenerateReports generates reports for a specific date range
func (h *SchedulerHandler) GenerateReports(c *gin.Context) {
	var req dto.GenerateReportRequest
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

	// Create context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Minute)
	defer cancel()

	// Generate reports
	csvFiles, err := h.schedulerService.GenerateReportsByDateRange(ctx, req.StartDate, req.EndDate)
	if err != nil {
		log.Printf("[SchedulerHandler] Error generating reports: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "Failed to generate reports: " + err.Error(),
		})
		return
	}

	// Prepare response with file information - include ALL collections
	log.Printf("[SchedulerHandler] Received %d collections from service", len(csvFiles))
	var files []dto.CSVFileInfo
	for collection, result := range csvFiles {
		var fileName string
		if result.FilePath != "" {
			fileName = filepath.Base(result.FilePath)
		} else {
			// Generate filename even if file creation failed (for consistency)
			// This ensures the frontend knows the report was attempted
			fileName = fmt.Sprintf("%s_%s_to_%s.csv", collection, req.StartDate, req.EndDate)
		}
		files = append(files, dto.CSVFileInfo{
			Collection:  collection,
			FileName:    fileName,
			RecordCount: result.RecordCount,
		})
		log.Printf("[SchedulerHandler] Added collection to response: %s (records: %d)", collection, result.RecordCount)
	}

	log.Printf("[SchedulerHandler] Successfully generated %d reports", len(files))
	c.JSON(http.StatusOK, dto.GenerateReportResponse{
		Success: true,
		Message: "Reports generated successfully",
		Files:   files,
	})
}

// DownloadReport downloads a specific CSV report
func (h *SchedulerHandler) DownloadReport(c *gin.Context) {
	fileName := c.Param("fileName")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "File name is required",
		})
		return
	}

	// Security check: prevent directory traversal
	if strings.Contains(fileName, "..") || strings.Contains(fileName, "/") || strings.Contains(fileName, "\\") {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid file name",
		})
		return
	}

	// Get full file path
	filePath := h.schedulerService.GetCSVFilePath(fileName)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Message: "File not found",
		})
		return
	}

	// Serve the file for download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "text/csv")
	c.File(filePath)
}

