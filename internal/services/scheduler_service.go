package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"admin-be/internal/config"
	"admin-be/internal/repository"
	"admin-be/internal/utils"

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

// Start initializes the scheduler service (cron job removed - reports can be generated manually via API)
func (s *SchedulerService) Start() error {
	// Ensure output directory exists
	if err := utils.EnsureDir(s.outputDir); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	log.Printf("[Scheduler] Scheduler service initialized successfully")
	log.Printf("[Scheduler] Daily cron job has been removed - reports can be generated manually via API")

	return nil
}

// Stop stops the scheduler service (no-op since cron job has been removed)
func (s *SchedulerService) Stop() {
	log.Println("[Scheduler] Scheduler service stopped")
}

// GenerateAndSendDailyReports fetches the previous day's records from all collections and sends them via email
func (s *SchedulerService) GenerateAndSendDailyReports() error {
	log.Println("[Scheduler] Starting daily report generation...")

	// Get the previous day's date range in UTC (00:00:00 to 23:59:59.999999999)
	// When running at 2:35 PM, we fetch data from the previous day (yesterday)
	// Use UTC to match the createdOn field format in MongoDB
	now := time.Now().UTC()
	// Subtract one day to get yesterday's date
	yesterday := now.AddDate(0, 0, -1)
	startTime := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)
	endTime := startTime.Add(24 * time.Hour)

	log.Printf("[Scheduler] Fetching records for previous day: %s (from %s to %s)",
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

		// Generate CSV file with date range in filename
		csvPath, err := utils.GenerateCSVWithDateRange(records, collectionName, s.outputDir, startTime, endTime)
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

// GenerateReportsByDateRange generates reports for a specific date range
// startDate format: YYYY-MM-DD, endDate format: YYYY-MM-DD
// Collects data from 5 AM of startDate to 11:55 PM of endDate
func (s *SchedulerService) GenerateReportsByDateRange(ctx context.Context, startDate, endDate string) (map[string]string, error) {
	log.Printf("[Scheduler] Starting report generation for date range: %s to %s", startDate, endDate)

	// Parse dates
	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}
	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	// Set start time to 12:05 AM UTC (00:05)
	startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 5, 0, 0, time.UTC)
	// Set end time to 11:55 PM UTC (23:55)
	endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 55, 0, 0, time.UTC)

	log.Printf("[Scheduler] Fetching records from %s to %s", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

	// Ensure output directory exists
	if err := utils.EnsureDir(s.outputDir); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	csvFiles := make(map[string]string) // collection name -> file path

	// Process each collection
	for _, collectionName := range s.collections {
		log.Printf("[Scheduler] Processing collection: %s", collectionName)

		// Fetch records from MongoDB
		records, err := s.repo.GetTodayRecords(ctx, collectionName, startTime, endTime)
		if err != nil {
			log.Printf("[Scheduler] Error fetching records from %s: %v", collectionName, err)
			continue
		}

		log.Printf("[Scheduler] Found %d records in collection: %s", len(records), collectionName)

		// Generate CSV file even if no records (empty CSV) with date range in filename
		csvPath, err := utils.GenerateCSVWithDateRange(records, collectionName, s.outputDir, startTime, endTime)
		if err != nil {
			log.Printf("[Scheduler] Error generating CSV for %s: %v", collectionName, err)
			continue
		}

		csvFiles[collectionName] = csvPath
		log.Printf("[Scheduler] Generated CSV for %s: %s", collectionName, csvPath)
	}

	log.Printf("[Scheduler] Report generation completed. Generated %d CSV files", len(csvFiles))
	return csvFiles, nil
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

// GetCSVFilePath returns the full path for a CSV file
func (s *SchedulerService) GetCSVFilePath(fileName string) string {
	return filepath.Join(s.outputDir, fileName)
}
