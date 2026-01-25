package repository

import (
	"admin-be/internal/models"
	"context"
	"time"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	FindAll(ctx context.Context, filter map[string]interface{}, skip, limit int64) ([]models.User, error)
	FindByUserID(ctx context.Context, userID string) (*models.User, error)
	UpdateByUserID(ctx context.Context, userID string, update map[string]interface{}) error
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)
}

// AstroRepository defines the interface for astrologer data operations
type AstroRepository interface {
	FindAll(ctx context.Context, filter map[string]interface{}, skip, limit int64) ([]models.Astrologer, error)
	FindByAstroID(ctx context.Context, astroID string) (*models.Astrologer, error)
	UpdateByAstroID(ctx context.Context, astroID string, update map[string]interface{}) error
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)
}

// ServiceReportRepository defines the interface for service report operations
type ServiceReportRepository interface {
	FindAll(ctx context.Context, filter map[string]interface{}, skip, limit int64) ([]models.ServiceReport, error)
	FindByReportID(ctx context.Context, reportID string) (*models.ServiceReport, error)
	FindByOrderID(ctx context.Context, orderId string) (*models.ServiceReport, error)
	FindByOrderIDAndServiceType(ctx context.Context, orderId, serviceType string) (*models.ServiceReport, error)
	UpdateByReportID(ctx context.Context, reportID string, update map[string]interface{}) error
	UpdateByOrderID(ctx context.Context, orderId string, update map[string]interface{}) error
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)
}

// UserProblemRepository defines the interface for user problem operations
type UserProblemRepository interface {
	FindAll(ctx context.Context, filter map[string]interface{}, skip, limit int64) ([]models.UserProblem, error)
	FindByProblemID(ctx context.Context, problemID string) (*models.UserProblem, error)
	UpdateByProblemID(ctx context.Context, problemID string, update map[string]interface{}) error
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)
}

// AstroProblemRepository defines the interface for astro problem operations
type AstroProblemRepository interface {
	FindAll(ctx context.Context, filter map[string]interface{}, skip, limit int64) ([]models.AstroProblem, error)
	FindByProblemID(ctx context.Context, problemID string) (*models.AstroProblem, error)
	UpdateByProblemID(ctx context.Context, problemID string, update map[string]interface{}) error
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)
}

// HoroscopeRepository defines the interface for horoscope operations
type HoroscopeRepository interface {
	FindAll(ctx context.Context, filter map[string]interface{}, skip, limit int64) ([]models.Horoscope, error)
	FindByHoroscopeID(ctx context.Context, horoscopeID string) (*models.Horoscope, error)
	FindBySignNameAndDate(ctx context.Context, signName, date string) (*models.Horoscope, error)
	Create(ctx context.Context, horoscope *models.Horoscope) error
	UpdateByHoroscopeID(ctx context.Context, horoscopeID string, update map[string]interface{}) error
	DeleteByHoroscopeID(ctx context.Context, horoscopeID string) error
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)
}

// ServiceRepository defines the interface for service operations (chat, ivr, video)
type ServiceRepository interface {
	FindChatByChatID(ctx context.Context, chatId string) (*models.Chat, error)
	FindIvrByIvrID(ctx context.Context, ivrID string) (*models.IvrCall, error)
	FindVideoByVideoID(ctx context.Context, videoID string) (*models.VideoCall, error)
	UpdateChatByChatID(ctx context.Context, chatID string, update map[string]interface{}) error
	UpdateIvrByIvrID(ctx context.Context, ivrID string, update map[string]interface{}) error
	UpdateVideoByVideoID(ctx context.Context, videoID string, update map[string]interface{}) error
	CountByLastStatus(ctx context.Context, collectionName string, startTime, endTime time.Time) (map[string]int64, error)
}

// FeedbackRepository defines the interface for feedback operations
type FeedbackRepository interface {
	Create(ctx context.Context, feedback *models.Feedback) error
	BulkCreate(ctx context.Context, feedbacks []*models.Feedback) error
}

