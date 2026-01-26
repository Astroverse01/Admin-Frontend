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

	"go.mongodb.org/mongo-driver/bson"
)

type SchedulerService struct {
	repo        repository.DailyReportRepository
	cfg         *config.Config
	outputDir   string
	collections []string
}

func NewSchedulerService(repo repository.DailyReportRepository, cfg *config.Config) *SchedulerService {
	return &SchedulerService{
		repo:      repo,
		cfg:       cfg,
		outputDir: filepath.Join(os.TempDir(), "daily_reports"),
		collections: []string{
			"appointments",
			"astrologer",
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

// GenerateAndSendDailyReports fetches the previous day's records from all collections and generates CSV files
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
		log.Println("[Scheduler] No CSV files generated.")
		return nil
	}

	log.Printf("[Scheduler] Daily report generation completed successfully. Generated %d CSV files", len(csvFilePaths))
	return nil
}

// CollectionReportResult represents the result of generating a report for a collection
type CollectionReportResult struct {
	FilePath    string
	RecordCount int
	Error       error
}

// GenerateReportsByDateRange generates reports for a specific date range
// startDate format: YYYY-MM-DD, endDate format: YYYY-MM-DD
// Collects data from 12:05 AM of startDate to 11:59 PM of endDate
func (s *SchedulerService) GenerateReportsByDateRange(ctx context.Context, startDate, endDate string) (map[string]CollectionReportResult, error) {
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
	// Set end time to 11:59 PM UTC (23:59)
	endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 0, 0, time.UTC)

	log.Printf("[Scheduler] Fetching records from %s to %s", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

	// Ensure output directory exists
	if err := utils.EnsureDir(s.outputDir); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	csvFiles := make(map[string]CollectionReportResult) // collection name -> result

	// Initialize all collections in the result map first to ensure they're all included
	for _, collectionName := range s.collections {
		csvFiles[collectionName] = CollectionReportResult{
			FilePath:    "",
			RecordCount: 0,
			Error:       nil,
		}
	}

	// Process each collection
	for _, collectionName := range s.collections {
		log.Printf("[Scheduler] Processing collection: %s", collectionName)

		// Fetch records from MongoDB
		records, err := s.repo.GetTodayRecords(ctx, collectionName, startTime, endTime)
		recordCount := 0
		if err != nil {
			log.Printf("[Scheduler] Error fetching records from %s: %v", collectionName, err)
			// Still try to generate empty CSV even on fetch error
			records = []bson.M{}
		} else {
			recordCount = len(records)
			log.Printf("[Scheduler] Found %d records in collection: %s", recordCount, collectionName)
		}

		// Generate CSV file even if no records (empty CSV) with date range in filename
		csvPath, err := utils.GenerateCSVWithDateRange(records, collectionName, s.outputDir, startTime, endTime)
		if err != nil {
			log.Printf("[Scheduler] Error generating CSV for %s: %v", collectionName, err)
			csvFiles[collectionName] = CollectionReportResult{
				FilePath:    "",
				RecordCount: recordCount,
				Error:       err,
			}
			continue
		}

		csvFiles[collectionName] = CollectionReportResult{
			FilePath:    csvPath,
			RecordCount: recordCount,
			Error:       nil,
		}
		log.Printf("[Scheduler] Generated CSV for %s: %s", collectionName, csvPath)
	}

	log.Printf("[Scheduler] Report generation completed. Processed %d collections", len(csvFiles))

	// Verify all collections are present in the result
	expectedCount := len(s.collections)
	actualCount := len(csvFiles)
	if actualCount != expectedCount {
		log.Printf("[Scheduler] WARNING: Expected %d collections but only %d in result. Missing collections:", expectedCount, actualCount)
		for _, collectionName := range s.collections {
			if _, exists := csvFiles[collectionName]; !exists {
				log.Printf("[Scheduler] WARNING: Missing collection: %s", collectionName)
				// Add it with empty result
				csvFiles[collectionName] = CollectionReportResult{
					FilePath:    "",
					RecordCount: 0,
					Error:       fmt.Errorf("collection was not processed"),
				}
			}
		}
	}

	log.Printf("[Scheduler] Final result contains %d collections", len(csvFiles))
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
