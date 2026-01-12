package main

import (
	"log"
	"os"

	"astro-admin/internal/config"
	"astro-admin/internal/database"
	"astro-admin/internal/handlers"
	"astro-admin/internal/middleware"
	"astro-admin/internal/repository"
	"astro-admin/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize configuration
	cfg := config.Load()

	// Initialize database connections
	mongoDB := database.NewMongoDB(cfg.MongoURI)
	redisDB, err := database.NewRedis(cfg.RedisURL)
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v. Continuing without Redis...", err)
		redisDB = nil
	} else {
		log.Println("Redis connected successfully")
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(mongoDB)
	astroRepo := repository.NewAstroRepository(mongoDB)
	serviceReportRepo := repository.NewServiceReportRepository(mongoDB)
	userProblemRepo := repository.NewUserProblemRepository(mongoDB)
	astroProblemRepo := repository.NewAstroProblemRepository(mongoDB)
	horoscopeRepo := repository.NewHoroscopeRepository(mongoDB)
	dailyReportRepo := repository.NewDailyReportRepository(mongoDB)

	// Initialize services
	authService := services.NewAuthService(cfg.JWTSecret)
	userService := services.NewUserService(userRepo)
	astroService := services.NewAstroService(astroRepo)
	emailService := services.NewEmailService(cfg)
	complaintService := services.NewComplaintService(serviceReportRepo, userRepo, astroRepo, mongoDB, redisDB, emailService)
	userProblemService := services.NewUserProblemService(userProblemRepo, userRepo, redisDB)
	astroProblemService := services.NewAstroProblemService(astroProblemRepo, astroRepo, redisDB)
	horoscopeService := services.NewHoroscopeService(horoscopeRepo)
	schedulerService := services.NewSchedulerService(dailyReportRepo, emailService, cfg)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	astroHandler := handlers.NewAstroHandler(astroService)
	complaintHandler := handlers.NewComplaintHandler(complaintService)
	userProblemHandler := handlers.NewUserProblemHandler(userProblemService)
	astroProblemHandler := handlers.NewAstroProblemHandler(astroProblemService)
	horoscopeHandler := handlers.NewHoroscopeHandler(horoscopeService, cfg.AdminID)

	// Setup router
	router := gin.Default()

	// Public routes
	router.POST("/admin/login", authHandler.Login)

	// Protected admin routes
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// User management
		admin.GET("/users", userHandler.ListUsers)
		admin.PATCH("/users/:userId/deactivate", userHandler.DeactivateUser)

		// Astro management
		admin.GET("/astros", astroHandler.ListAstros)
		admin.PATCH("/astros/:astroId/status", astroHandler.UpdateAstroStatus)
		admin.PATCH("/astros/:astroId/visibility", astroHandler.ToggleVisibility)

		// User service complaints
		admin.GET("/user-service-complaints", complaintHandler.ListUserServiceComplaints)
		admin.GET("/user-service-complaints/:serviceType/:orderId", complaintHandler.GetComplaintDetails)
		admin.PATCH("/user-service-complaints/:orderId", complaintHandler.AcceptRejectComplaint)

		// User general complaints
		admin.GET("/user-general-complaints", userProblemHandler.ListUserGeneralComplaints)
		admin.PATCH("/user-general-complaints/:problemId/close", userProblemHandler.CloseComplaint)

		// Astro general complaints
		admin.GET("/astro-general-complaints", astroProblemHandler.ListAstroGeneralComplaints)
		admin.PATCH("/astro-general-complaints/:problemId/close", astroProblemHandler.CloseComplaint)

		// Horoscope management
		admin.GET("/horoscopes", horoscopeHandler.ListHoroscopes)
		admin.POST("/horoscopes/bulk", horoscopeHandler.BulkCreateHoroscopes)
		admin.PATCH("/horoscopes/:horoscopeId", horoscopeHandler.UpdateHoroscope)
		admin.DELETE("/horoscopes/:horoscopeId", horoscopeHandler.DeleteHoroscope)
	}

	// Start daily report scheduler
	if err := schedulerService.Start(); err != nil {
		log.Printf("Warning: Failed to start scheduler: %v. Continuing without scheduler...", err)
	} else {
		log.Println("Daily report scheduler started successfully")
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}

	// Stop scheduler gracefully when server stops
	schedulerService.Stop()
}
