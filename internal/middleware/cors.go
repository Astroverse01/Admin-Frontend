package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles CORS headers for all requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Allow specific origins for development
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:5173", // Vite default port
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
		}

		// Check if origin is in allowed list
		allowed := false
		for _, ao := range allowedOrigins {
			if origin == ao {
				allowed = true
				break
			}
		}

		// For development: be more permissive - allow any localhost origin
		// In production, restrict this to specific domains only
		if origin != "" {
			// Allow if it's in the allowed list OR if it's a localhost origin (for development)
			if allowed || strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1") {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}

		// Set CORS headers
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

