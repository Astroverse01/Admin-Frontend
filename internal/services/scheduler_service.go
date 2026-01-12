package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"astro-admin/internal/config"
	"astro-admin/internal/repository"
	"astro-admin/internal/utils"

	"github.com/robfig/cron/v3"
)

type SchedulerService struct {
	repo         repository.DailyReportRepository
	emailService *EmailService
	cfg          *config.Config
	cron         *cron.Cron
	outputDir    string
	collections  []string
}

func NewSchedulerService(repo repository.DailyReportRepository, emailService *EmailService, cfg *config.Config) *SchedulerService {
	return &SchedulerService{
		repo:         repo,
		emailService: emailService,
		cfg:          cfg,
		cron:         cron.New(cron.WithSeconds()),
		outputDir:    filepath.Join(os.TempDir(), "daily_reports"),
		collections: []string{
			"appointments",
			"chat",
			"conversionHistory",
			"feedback",
			"ivrCall",
			"orderScore",
			"rewards",
			"serviceReports",
			"user",
			"userPayment",
			"userProblem",
			"videoCall",
		},
	}
}

// Start starts the cron scheduler
func (s *SchedulerService) Start() error {
	// Validate email configuration from .env before starting
	if err := s.emailService.validateEmailConfig(); err != nil {
		return fmt.Errorf("failed to start scheduler: %w", err)
	}

	// Ensure output directory exists
	if err := utils.EnsureDir(s.outputDir); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Schedule cron job to run daily at 11:55 PM
	// Cron format (with seconds): second minute hour day month weekday
	// "0 55 23 * * *" means 23:55:00 (11:55 PM) every day
	_, err := s.cron.AddFunc("0 55 23 * * *", func() {
		log.Println("[Scheduler] Daily report job triggered at 11:55 PM")
		if err := s.GenerateAndSendDailyReports(); err != nil {
			log.Printf("[Scheduler] Error generating daily reports: %v", err)
		}
	})

	if err != nil {
		return fmt.Errorf("failed to schedule cron job: %w", err)
	}

	s.cron.Start()
	log.Println("[Scheduler] Daily report scheduler started. Will run at 11:55 PM every day.")
	return nil
}

// Stop stops the cron scheduler
func (s *SchedulerService) Stop() {
	s.cron.Stop()
	log.Println("[Scheduler] Daily report scheduler stopped")
}

// GenerateAndSendDailyReports fetches today's records from all collections and sends them via email
func (s *SchedulerService) GenerateAndSendDailyReports() error {
	log.Println("[Scheduler] Starting daily report generation...")

	// Get today's date range in UTC (00:00:00 to 23:59:59.999999999)
	// Use UTC to match the createdOn field format in MongoDB
	now := time.Now().UTC()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endTime := startTime.Add(24 * time.Hour)

	log.Printf("[Scheduler] Fetching records for date: %s (from %s to %s)",
		startTime.Format("2006-01-02"),
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	var csvFilePaths []string

	// Process each collection
	for _, collectionName := range s.collections {
		log.Printf("[Scheduler] Processing collection: %s", collectionName)

		// Fetch records from MongoDB
		records, err := s.repo.GetTodayRecords(ctx, collectionName, startTime, endTime)
		if err != nil {
			log.Printf("[Scheduler] Error fetching records from %s: %v", collectionName, err)
			// Continue with other collections even if one fails
			continue
		}

		log.Printf("[Scheduler] Found %d records in collection: %s", len(records), collectionName)

		// Generate CSV file
		csvPath, err := utils.GenerateCSV(records, collectionName, s.outputDir)
		if err != nil {
			log.Printf("[Scheduler] Error generating CSV for %s: %v", collectionName, err)
			// Continue with other collections even if one fails
			continue
		}

		csvFilePaths = append(csvFilePaths, csvPath)
		log.Printf("[Scheduler] Generated CSV for %s: %s", collectionName, csvPath)
	}

	if len(csvFilePaths) == 0 {
		log.Println("[Scheduler] No CSV files generated. Skipping email send.")
		return nil
	}

	// Send email with all CSV files
	log.Printf("[Scheduler] Sending email with %d CSV files...", len(csvFilePaths))
	if err := s.emailService.SendCSVFilesByEmail(csvFilePaths, startTime.Format("2006-01-02")); err != nil {
		log.Printf("[Scheduler] Error sending email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Println("[Scheduler] Daily report generation completed successfully")

	// Clean up CSV files after sending (optional)
	go s.cleanupCSVFiles(csvFilePaths)

	return nil
}

// cleanupCSVFiles removes CSV files after sending email (runs in background)
func (s *SchedulerService) cleanupCSVFiles(filePaths []string) {
	time.Sleep(5 * time.Minute) // Wait 5 minutes before cleanup
	for _, filePath := range filePaths {
		if err := os.Remove(filePath); err != nil {
			log.Printf("[Scheduler] Error cleaning up file %s: %v", filePath, err)
		} else {
			log.Printf("[Scheduler] Cleaned up file: %s", filePath)
		}
	}
}

// RunManually allows manual execution of the daily report job (for testing)
func (s *SchedulerService) RunManually() error {
	log.Println("[Scheduler] Manual execution triggered")
	return s.GenerateAndSendDailyReports()
}
