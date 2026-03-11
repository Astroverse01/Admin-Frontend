package handlers

import (
	"admin-be/internal/dto"
	"admin-be/internal/services"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/flexstack/uuid"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BlogHandler struct {
	blogService *services.BlogService
	s3Service   *services.S3Service
}

func NewBlogHandler(blogService *services.BlogService, s3Service *services.S3Service) *BlogHandler {
	return &BlogHandler{blogService: blogService, s3Service: s3Service}
}

type UploadCoverImageResponse struct {
	Success bool   `json:"success"`
	Key     string `json:"key"`
	URL     string `json:"url"`
}

func (h *BlogHandler) UploadCoverImage(c *gin.Context) {
	if h.s3Service == nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Success: false, Message: "S3 is not configured"})
		return
	}

	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "file is required (multipart field: file)"})
		return
	}

	// UI says up to 5MB
	const maxBytes = 5 * 1024 * 1024
	if fh.Size > maxBytes {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "file too large (max 5MB)"})
		return
	}

	f, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "failed to open file"})
		return
	}
	defer f.Close()

	b, err := services.ReadAllWithLimit(f, maxBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".webp", ".gif":
	default:
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "unsupported file type"})
		return
	}

	contentType := fh.Header.Get("Content-Type")
	if contentType == "" {
		// best-effort default
		contentType = "application/octet-stream"
	}

	// key: blogs/covers/<uuid>.<ext>
	u, err := uuid.NewV7()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Success: false, Message: "failed to generate file id"})
		return
	}
	key := fmt.Sprintf("blogs/covers/%s%s", u.String(), ext)

	keyOut, urlOut, err := h.s3Service.UploadPublic(c.Request.Context(), key, b, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, UploadCoverImageResponse{Success: true, Key: keyOut, URL: urlOut})
}

func (h *BlogHandler) CreateBlog(c *gin.Context) {
	var req dto.CreateBlogRequest
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("bind error:", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "Invalid request body"})
		return
	}

	// If multipart/form-data with an uploaded file, upload to S3 and set coverImage to the S3 URL.
	if strings.HasPrefix(c.GetHeader("Content-Type"), "multipart/form-data") && h.s3Service != nil {
		fh, err := c.FormFile("coverImage")
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "coverImage file is required (multipart field: coverImage)"})
			return
		}

		const maxBytes = 5 * 1024 * 1024
		if fh.Size > maxBytes {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "file too large (max 5MB)"})
			return
		}

		f, err := fh.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "failed to open file"})
			return
		}
		defer f.Close()

		b, err := services.ReadAllWithLimit(f, maxBytes)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: err.Error()})
			return
		}

		ext := strings.ToLower(filepath.Ext(fh.Filename))
		switch ext {
		case ".png", ".jpg", ".jpeg", ".webp", ".gif":
		default:
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "unsupported file type"})
			return
		}

		contentType := fh.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		u, err := uuid.NewV7()
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Success: false, Message: "failed to generate file id"})
			return
		}
		key := fmt.Sprintf("blogs/covers/%s%s", u.String(), ext)

		_, urlOut, err := h.s3Service.UploadPublic(c.Request.Context(), key, b, contentType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Success: false, Message: err.Error()})
			return
		}

		// Use uploaded S3 URL as coverImage in the blog document.
		req.CoverImage = urlOut
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "Validation failed: " + err.Error()})
		return
	}

	resp, err := h.blogService.CreateBlog(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *BlogHandler) ListBlogs(c *gin.Context) {
	q := c.Query("q")
	status := c.Query("status")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	resp, err := h.blogService.ListBlogs(c.Request.Context(), q, status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Success: false, Message: "Failed to fetch blogs: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *BlogHandler) GetBlog(c *gin.Context) {
	blogID := c.Param("blogId")
	if blogID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "blogId is required"})
		return
	}

	resp, err := h.blogService.GetBlog(c.Request.Context(), blogID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *BlogHandler) UpdateBlog(c *gin.Context) {
	blogID := c.Param("blogId")
	if blogID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "blogId is required"})
		return
	}

	var req dto.UpdateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "Invalid request body"})
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "Validation failed: " + err.Error()})
		return
	}

	resp, err := h.blogService.UpdateBlog(c.Request.Context(), blogID, &req)
	if err != nil {
		// service returns user-friendly errors; treat most as 400, but missing as 404
		if err.Error() == "blog not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Success: false, Message: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *BlogHandler) DeleteBlog(c *gin.Context) {
	blogID := c.Param("blogId")
	if blogID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "blogId is required"})
		return
	}

	resp, err := h.blogService.DeleteBlog(c.Request.Context(), blogID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Success: false, Message: "Failed to delete blog: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

