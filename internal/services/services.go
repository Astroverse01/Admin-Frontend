package services

import (
	"admin-be/internal/database"
	"admin-be/internal/dto"
	"admin-be/internal/models"
	"admin-be/internal/repository"
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/flexstack/uuid"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuthService handles authentication logic
type AuthService struct {
	jwtSecret string
}

func NewAuthService(jwtSecret string) *AuthService {
	return &AuthService{
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) ValidateCredentials(username, password string) bool {
	// Static validation for now - replace with database lookup later
	return username == "admin" && password == "admin123"
}

func (s *AuthService) GenerateToken(username string) (string, error) {
	// Create JWT token with 24-hour expiration
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// UserService handles user management logic
type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) ListUsers(ctx context.Context, name string, sort string, page, limit int) (*dto.UserListResponse, error) {
	log.Printf("[ListUsers Service] Called with params - name: %s, sort: %s, page: %d, limit: %d", name, sort, page, limit)

	// Build filter
	filter := make(map[string]interface{})
	if name != "" {
		// Search in both name and fullName fields
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": regexp.QuoteMeta(name), "$options": "i"}},
			{"fullName": bson.M{"$regex": regexp.QuoteMeta(name), "$options": "i"}},
		}
		log.Printf("[ListUsers Service] Built filter with name regex for both name and fullName fields")
	} else {
		log.Println("[ListUsers Service] No name filter, using empty filter")
	}

	// Calculate pagination
	skip := int64((page - 1) * limit)
	log.Printf("[ListUsers Service] Pagination - skip: %d, limit: %d", skip, limit)

	// Get users
	log.Println("[ListUsers Service] Calling userRepo.FindAll")
	users, err := s.userRepo.FindAll(ctx, filter, skip, int64(limit))
	if err != nil {
		log.Printf("[ListUsers Service] Error from userRepo.FindAll: %v", err)
		return nil, fmt.Errorf("failed to fetch users from repository: %w", err)
	}
	log.Printf("[ListUsers Service] Successfully fetched %d users from repository", len(users))

	// Get total count
	log.Println("[ListUsers Service] Calling userRepo.Count")
	total, err := s.userRepo.Count(ctx, filter)
	if err != nil {
		log.Printf("[ListUsers Service] Error from userRepo.Count: %v", err)
		return nil, fmt.Errorf("failed to count users: %w", err)
	}
	log.Printf("[ListUsers Service] Total users count: %d", total)

	// Convert to DTO
	var userData []dto.UserData
	for _, user := range users {
		status := "active"
		if user.IsDelete == 1 {
			status = "inactive"
		}

		// Use name field if available, otherwise fall back to fullName field
		name := user.Name
		if name == "" && user.FullName != "" {
			name = user.FullName
		}

		userData = append(userData, dto.UserData{
			UserID: user.UserID,
			Name:   name,
			Status: status,
		})
	}
	log.Printf("[ListUsers Service] Converted %d users to DTO", len(userData))

	// Apply sorting
	if sort == "desc" {
		log.Println("[ListUsers Service] Applying descending sort")
		// Reverse the slice for descending order
		for i, j := 0, len(userData)-1; i < j; i, j = i+1, j-1 {
			userData[i], userData[j] = userData[j], userData[i]
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	log.Printf("[ListUsers Service] Calculated total pages: %d", totalPages)

	log.Printf("[ListUsers Service] Returning response with %d users, total: %d, page: %d/%d", len(userData), total, page, totalPages)
	return &dto.UserListResponse{
		Success: true,
		Data:    userData,
		Pagination: dto.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (s *UserService) DeactivateUser(ctx context.Context, userID string, status string) error {
	var isDelete int
	if status == "inactive" {
		isDelete = 1
	} else {
		isDelete = 0
	}

	update := map[string]interface{}{
		"isDelete": isDelete,
	}

	return s.userRepo.UpdateByUserID(ctx, userID, update)
}

// AstroService handles astrologer management logic
type AstroService struct {
	astroRepo repository.AstroRepository
}

func NewAstroService(astroRepo repository.AstroRepository) *AstroService {
	return &AstroService{
		astroRepo: astroRepo,
	}
}

func (s *AstroService) ListAstros(ctx context.Context, name string, sort string, page, limit int) (*dto.AstroListResponse, error) {
	// Build filter
	filter := make(map[string]interface{})
	if name != "" {
		// Search in both name and fullName fields
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": regexp.QuoteMeta(name), "$options": "i"}},
			{"fullName": bson.M{"$regex": regexp.QuoteMeta(name), "$options": "i"}},
		}
	}

	// Calculate pagination
	skip := int64((page - 1) * limit)

	// Get astros
	astros, err := s.astroRepo.FindAll(ctx, filter, skip, int64(limit))
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.astroRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	var astroData []dto.AstroData
	for _, astro := range astros {
		status := "active"
		if astro.IsDelete == 1 {
			status = "inactive"
		}

		visible := "hidden"
		if astro.IsActive == 1 {
			visible = "visible"
		}

		// Use name field if available, otherwise fall back to fullName field
		name := astro.Name
		if name == "" && astro.FullName != "" {
			name = astro.FullName
		}

		astroData = append(astroData, dto.AstroData{
			AstroID: astro.AstroID,
			Name:    name,
			Status:  status,
			Visible: visible,
		})
	}

	// Apply sorting
	if sort == "desc" {
		// Reverse the slice for descending order
		for i, j := 0, len(astroData)-1; i < j; i, j = i+1, j-1 {
			astroData[i], astroData[j] = astroData[j], astroData[i]
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.AstroListResponse{
		Success: true,
		Data:    astroData,
		Pagination: dto.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (s *AstroService) UpdateAstroStatus(ctx context.Context, astroID string, status string) error {
	var isDelete, isActive int
	if status == "inactive" {
		// When deactivating: set isDelete=1 and isActive=0 (status=inactive, visibility=hidden)
		isDelete = 1
		isActive = 0
	} else {
		// When activating: set isDelete=0 and isActive=1 (status=active, visibility=visible)
		isDelete = 0
		isActive = 1
	}

	update := map[string]interface{}{
		"isDelete": isDelete,
		"isActive": isActive,
	}

	return s.astroRepo.UpdateByAstroID(ctx, astroID, update)
}

func (s *AstroService) ToggleVisibility(ctx context.Context, astroID string, visible bool) error {
	// Only update isActive field (visibility)
	// This allows toggling visibility independently of deactivation status
	var isActive int
	if visible {
		isActive = 1
	} else {
		isActive = 0
	}

	update := map[string]interface{}{
		"isActive": isActive,
	}

	return s.astroRepo.UpdateByAstroID(ctx, astroID, update)
}

// ComplaintService handles service complaint logic
type ComplaintService struct {
	serviceReportRepo repository.ServiceReportRepository
	userRepo          repository.UserRepository
	astroRepo         repository.AstroRepository
	serviceRepo       repository.ServiceRepository
	emailService      *EmailService
}

func NewComplaintService(serviceReportRepo repository.ServiceReportRepository, userRepo repository.UserRepository, astroRepo repository.AstroRepository, mongoDB *database.MongoDB, emailService *EmailService) *ComplaintService {
	return &ComplaintService{
		serviceReportRepo: serviceReportRepo,
		userRepo:          userRepo,
		astroRepo:         astroRepo,
		serviceRepo:       repository.NewServiceRepository(mongoDB),
		emailService:      emailService,
	}
}

func (s *ComplaintService) ListUserServiceComplaints(ctx context.Context, serviceType, status string, page, limit int) (*dto.ComplaintListResponse, error) {
	// Build filter
	filter := make(map[string]interface{})
	if serviceType != "" {
		filter["serviceType"] = serviceType
	}
	if status != "" {
		filter["status"] = status
	}

	// Calculate pagination
	skip := int64((page - 1) * limit)

	// Get complaints
	complaints, err := s.serviceReportRepo.FindAll(ctx, filter, skip, int64(limit))
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.serviceReportRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Convert to DTO with names
	var complaintData []dto.ComplaintData
	for _, complaint := range complaints {
		// Fetch user name
		userName := "Unknown"
		if user, err := s.userRepo.FindByUserID(ctx, complaint.UserID); err == nil {
			if user.Name != "" {
				userName = user.Name
			} else if user.FullName != "" {
				userName = user.FullName
			}
		}

		// Fetch astrologer name
		astroName := "Unknown"
		if astro, err := s.astroRepo.FindByAstroID(ctx, complaint.AstroID); err == nil {
			if astro.Name != "" {
				astroName = astro.Name
			} else if astro.FullName != "" {
				astroName = astro.FullName
			}
		}

		complaintData = append(complaintData, dto.ComplaintData{
			ServiceType: complaint.ServiceType,
			OrderID:     complaint.OrderID,
			AstroID:     complaint.AstroID,
			AstroName:   astroName,
			UserID:      complaint.UserID,
			UserName:    userName,
			CreatedOn:   complaint.CreatedOn.Format("2006-01-02 15:04:05"),
			Status:      complaint.Status,
			Comment:     complaint.Comment,
			ReportID:    complaint.ReportID,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.ComplaintListResponse{
		Success: true,
		Data:    complaintData,
		Pagination: dto.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (s *ComplaintService) GetComplaintDetails(ctx context.Context, serviceType, orderId string) (*dto.ServiceDetailsResponse, error) {
	// Trim whitespace from inputs
	serviceType = strings.TrimSpace(serviceType)
	orderId = strings.TrimSpace(orderId)

	// Normalize serviceType to match collection names (camelCase for ivrCall and videoCall)
	//serviceTypeLower := strings.ToLower(serviceType)

	// Map to correct camelCase format based on collection names
	switch serviceType {
	case "ivrCall":
		serviceType = "ivrCall" // Collection name is "ivrCall"
	case "videoCall":
		serviceType = "videoCall" // Collection name is "videoCall"
	case "chat":
		serviceType = "chat" // Collection name is "chat"
	default:
		// Preserve original if it doesn't match known patterns
		// No assignment needed - serviceType already has its value
	}

	// First, verify that a complaint exists for this orderId and serviceType
	// Repository now handles case-insensitive matching via regex
	log.Printf("Searching for complaint with orderId: %s, serviceType: %s", orderId, serviceType)
	complaint, err := s.serviceReportRepo.FindByOrderIDAndServiceType(ctx, orderId, serviceType)
	if err != nil {
		// Fallback: If searching for "ivrCall" fails, try "call" (legacy data might use "call")
		if serviceType == "ivrCall" {
			log.Printf("Initial search failed for 'ivrCall', trying fallback 'call'")
			complaint, err = s.serviceReportRepo.FindByOrderIDAndServiceType(ctx, orderId, "ivrCall")
			if err == nil {
				log.Printf("Found complaint using fallback 'call', updating serviceType to 'ivrCall'")
				serviceType = "ivrCall" // Keep the normalized serviceType for switch statement
			}
		}

		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				log.Printf("No complaint found for orderId: %s, serviceType: %s", orderId, serviceType)
				return nil, fmt.Errorf("complaint not found for orderId %s and serviceType %s", orderId, serviceType)
			}
			log.Printf("Error finding complaint: %v", err)
			return nil, fmt.Errorf("failed to find complaint: %w", err)
		}
	}
	log.Printf("Complaint found: reportID=%s, orderID=%s, serviceType=%s", complaint.ReportID, complaint.OrderID, complaint.ServiceType)

	var serviceData dto.ServiceData
	reportID := complaint.ReportID

	switch serviceType {
	case "chat":
		// orderId is a foreign key in serviceReports collection that references chatId in the chat collection
		log.Printf("Looking up chat with chatId (orderId): %s", orderId)
		chat, err := s.serviceRepo.FindChatByChatID(ctx, orderId)
		if err != nil {
			log.Printf("Error finding chat by chatId %s: %v (error type: %T)", orderId, err, err)
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, fmt.Errorf("chat service not found for orderId %s", orderId)
			}
			return nil, fmt.Errorf("failed to find chat service: %w", err)
		}
		log.Printf("Chat object found: %+v", chat)
		serviceData = dto.ServiceData{
			ServiceID:     chat.ChatId,
			AstroID:       chat.AstroID,
			UserID:        chat.UserID,
			Conversation:  chat.Conversation,
			RatePerMinute: chat.RatePerMinute,
			Type:          "chat",
			ReportID:      reportID,
		}
	case "ivrCall":
		ivr, err := s.serviceRepo.FindIvrByIvrID(ctx, orderId)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, fmt.Errorf("IVR service not found for orderId %s", orderId)
			}
			return nil, fmt.Errorf("failed to find IVR service: %w", err)
		}
		if ivr.URL == "" {
			return nil, fmt.Errorf("audio is not available for this orderId")
		}
		serviceData = dto.ServiceData{
			ServiceID:     ivr.IvrID,
			AstroID:       ivr.AstroID,
			UserID:        ivr.UserID,
			URL:           ivr.URL,
			RatePerMinute: ivr.RatePerMinute,
			Type:          "ivr",
			ReportID:      reportID,
		}
	case "videoCall":
		video, err := s.serviceRepo.FindVideoByVideoID(ctx, orderId)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, fmt.Errorf("video service not found for orderId %s", orderId)
			}
			return nil, fmt.Errorf("failed to find video service: %w", err)
		}
		serviceData = dto.ServiceData{
			ServiceID:     video.VideoID,
			AstroID:       video.AstroID,
			UserID:        video.UserID,
			URL:           video.URL,
			RatePerMinute: video.RatePerMinute,
			Type:          "video",
			ReportID:      reportID,
		}
	default:
		return nil, fmt.Errorf("invalid service type: %s", serviceType)
	}

	return &dto.ServiceDetailsResponse{
		Success: true,
		Data:    serviceData,
	}, nil
}

func (s *ComplaintService) AcceptRejectComplaint(ctx context.Context, orderId string, req *dto.AcceptRejectRequest) (*dto.AcceptRejectResponse, error) {
	// Get complaint details by orderId
	complaint, err := s.serviceReportRepo.FindByOrderID(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to find complaint: %w", err)
	}

	if req.Action == "accept" {
		// Check if status is "open"
		if complaint.Status != "open" {
			return nil, fmt.Errorf("complaint status is not open, cannot accept")
		}

		// Update serviceReports status to "closed"
		err = s.serviceReportRepo.UpdateByOrderID(ctx, orderId, map[string]interface{}{
			"status": "closed",
			"reason": req.Reason,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update complaint status: %w", err)
		}

		// Update user: totalAmount += userRefundMoney (only if userRefundMoney > 0)
		if req.UserRefundMoney > 0 {
			user, err := s.userRepo.FindByUserID(ctx, complaint.UserID)
			if err != nil {
				return nil, fmt.Errorf("failed to find user: %w", err)
			}

			err = s.userRepo.UpdateByUserID(ctx, complaint.UserID, map[string]interface{}{
				"totalAmount": user.TotalAmount + req.UserRefundMoney,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to update user totalAmount: %w", err)
			}
		}

		// Update astrologer: totalEarned += astroRefundMoney (only if astroRefundMoney > 0)
		if req.AstroRefundMoney > 0 {
			astro, err := s.astroRepo.FindByAstroID(ctx, complaint.AstroID)
			if err != nil {
				return nil, fmt.Errorf("failed to find astrologer: %w", err)
			}

			err = s.astroRepo.UpdateByAstroID(ctx, complaint.AstroID, map[string]interface{}{
				"totalEarned": astro.TotalEarned + req.AstroRefundMoney,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to update astrologer totalEarned: %w", err)
			}
		}

		return &dto.AcceptRejectResponse{
			Success:        true,
			Message:        "Complaint accepted successfully",
			ComplaintID:    complaint.OrderID,
			UserID:         complaint.UserID,
			AstroID:        complaint.AstroID,
			RefundedAmount: req.UserRefundMoney,
		}, nil
	} else {
		// Reject complaint - Update serviceReports → status = "rejected"
		err = s.serviceReportRepo.UpdateByOrderID(ctx, orderId, map[string]interface{}{
			"status": "rejected",
			"reason": req.Reason,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update complaint status: %w", err)
		}

		// Fetch user and astro for email notifications
		user, err := s.userRepo.FindByUserID(ctx, complaint.UserID)
		if err != nil {
			log.Printf("Warning: Failed to fetch user for email notification: %v", err)
		}

		astro, err := s.astroRepo.FindByAstroID(ctx, complaint.AstroID)
		if err != nil {
			log.Printf("Warning: Failed to fetch astro for email notification: %v", err)
		}

		// Send email notifications to both user and astro
		if user != nil && user.Email != "" {
			userSubject := "Complaint Rejected"
			userBody := fmt.Sprintf("Dear %s,\n\nYour complaint (Order ID: %s) has been rejected.\n\nReason: %s\n\nThank you for your understanding.", user.Name, complaint.OrderID, req.Reason)
			if err := s.emailService.SendNotificationEmail(user.Email, userSubject, userBody); err != nil {
				log.Printf("Warning: Failed to send email notification to user %s: %v", user.Email, err)
			}
		}

		if astro != nil && astro.Email != "" {
			astroSubject := "Complaint Rejected"
			astroBody := fmt.Sprintf("Dear %s,\n\nA complaint (Order ID: %s) has been rejected.\n\nReason: %s\n\nThank you.", astro.Name, complaint.OrderID, req.Reason)
			if err := s.emailService.SendNotificationEmail(astro.Email, astroSubject, astroBody); err != nil {
				log.Printf("Warning: Failed to send email notification to astro %s: %v", astro.Email, err)
			}
		}

		return &dto.AcceptRejectResponse{
			Success:     true,
			Message:     "Complaint rejected",
			ComplaintID: complaint.OrderID,
			UserID:      complaint.UserID,
			AstroID:     complaint.AstroID,
		}, nil
	}
}

// UserProblemService handles user general complaint logic
type UserProblemService struct {
	userProblemRepo repository.UserProblemRepository
	userRepo        repository.UserRepository
}

func NewUserProblemService(userProblemRepo repository.UserProblemRepository, userRepo repository.UserRepository) *UserProblemService {
	return &UserProblemService{
		userProblemRepo: userProblemRepo,
		userRepo:        userRepo,
	}
}

func (s *UserProblemService) ListUserGeneralComplaints(ctx context.Context, problemTypes, status string, page, limit int) (*dto.ProblemListResponse, error) {
	// Build filter
	filter := make(map[string]interface{})
	if problemTypes != "" {
		filter["problemTypes"] = problemTypes
	}
	if status != "" {
		filter["status"] = status
	}

	// Calculate pagination
	skip := int64((page - 1) * limit)

	// Get problems
	problems, err := s.userProblemRepo.FindAll(ctx, filter, skip, int64(limit))
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.userProblemRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Convert to DTO and fetch user fullName
	var problemData []dto.ProblemData
	for _, problem := range problems {
		fullName := ""
		if problem.UserID != "" {
			user, err := s.userRepo.FindByUserID(ctx, problem.UserID)
			if err == nil && user != nil {
				fullName = user.FullName
			}
		}
		problemData = append(problemData, dto.ProblemData{
			ProblemID:    problem.ProblemID,
			UserID:       problem.UserID,
			FullName:     fullName,
			Status:       problem.Status,
			Comment:      problem.Comment,
			ProblemTypes: problem.ProblemTypes,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.ProblemListResponse{
		Success: true,
		Data:    problemData,
		Pagination: dto.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (s *UserProblemService) CloseComplaint(ctx context.Context, problemID string, reason string) (*dto.CloseComplaintResponse, error) {
	// Update complaint status
	err := s.userProblemRepo.UpdateByProblemID(ctx, problemID, map[string]interface{}{
		"status": "closed",
		"reason": reason,
	})
	if err != nil {
		return nil, err
	}

	// Get user details for notification (if needed in future)
	// problem, err := s.userProblemRepo.FindByProblemID(ctx, problemID)
	// if err != nil {
	// 	return nil, err
	// }
	// user, err := s.userRepo.FindByUserID(ctx, problem.UserID)
	// if err != nil {
	// 	return nil, err
	// }

	// Send notification (simplified)
	// sendNotification(user.FCMToken, fmt.Sprintf("Your complaint has been closed by Admin. Reason: %s", reason))

	return &dto.CloseComplaintResponse{
		Success:   true,
		Message:   "Complaint closed successfully",
		ProblemID: problemID,
		Status:    "closed",
	}, nil
}

// AstroProblemService handles astro general complaint logic
type AstroProblemService struct {
	astroProblemRepo repository.AstroProblemRepository
	astroRepo        repository.AstroRepository
}

func NewAstroProblemService(astroProblemRepo repository.AstroProblemRepository, astroRepo repository.AstroRepository) *AstroProblemService {
	return &AstroProblemService{
		astroProblemRepo: astroProblemRepo,
		astroRepo:        astroRepo,
	}
}

func (s *AstroProblemService) ListAstroGeneralComplaints(ctx context.Context, problemTypes, status string, page, limit int) (*dto.ProblemListResponse, error) {
	// Build filter
	filter := make(map[string]interface{})
	if problemTypes != "" {
		filter["problemTypes"] = problemTypes
	}
	if status != "" {
		filter["status"] = status
	}

	// Calculate pagination
	skip := int64((page - 1) * limit)

	// Get problems
	problems, err := s.astroProblemRepo.FindAll(ctx, filter, skip, int64(limit))
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.astroProblemRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	var problemData []dto.ProblemData
	for _, problem := range problems {
		problemData = append(problemData, dto.ProblemData{
			ProblemID:    problem.ProblemID,
			AstroID:      problem.AstroID,
			Status:       problem.Status,
			Comment:      problem.Comment,
			ProblemTypes: problem.ProblemTypes,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.ProblemListResponse{
		Success: true,
		Data:    problemData,
		Pagination: dto.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (s *AstroProblemService) CloseComplaint(ctx context.Context, problemID string, reason string) (*dto.CloseComplaintResponse, error) {
	// Update complaint status
	err := s.astroProblemRepo.UpdateByProblemID(ctx, problemID, map[string]interface{}{
		"status": "closed",
		"reason": reason,
	})
	if err != nil {
		return nil, err
	}

	// Get astro details for notification (if needed in future)
	// problem, err := s.astroProblemRepo.FindByProblemID(ctx, problemID)
	// if err != nil {
	// 	return nil, err
	// }
	// astro, err := s.astroRepo.FindByAstroID(ctx, problem.AstroID)
	// if err != nil {
	// 	return nil, err
	// }

	// Send notification (simplified)
	// sendNotification(astro.FCMToken, fmt.Sprintf("Your complaint has been closed by Admin. Reason: %s", reason))

	return &dto.CloseComplaintResponse{
		Success:   true,
		Message:   "Complaint closed successfully",
		ProblemID: problemID,
		Status:    "closed",
	}, nil
}

// HoroscopeService handles horoscope management logic
type HoroscopeService struct {
	horoscopeRepo repository.HoroscopeRepository
}

func NewHoroscopeService(horoscopeRepo repository.HoroscopeRepository) *HoroscopeService {
	return &HoroscopeService{
		horoscopeRepo: horoscopeRepo,
	}
}

func normalizeSignName(sign string) string {
	sign = strings.TrimSpace(sign)
	if sign == "" {
		return ""
	}

	return strings.ToLower(sign)
}

func (s *HoroscopeService) ListHoroscopes(ctx context.Context, date, signName string, page, limit int) (*dto.HoroscopeListResponse, error) {
	// Build filter
	filter := make(map[string]interface{})
	if date != "" {
		filter["date"] = date
	}
	if signName != "" {
		filter["signName"] = signName
	}

	// Calculate pagination
	skip := int64((page - 1) * limit)

	// Get horoscopes
	horoscopes, err := s.horoscopeRepo.FindAll(ctx, filter, skip, int64(limit))
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.horoscopeRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	var horoscopeData []dto.HoroscopeData
	for _, horoscope := range horoscopes {
		horoscopeData = append(horoscopeData, dto.HoroscopeData{
			HoroscopeID: horoscope.HoroscopeID,
			SignName:    horoscope.SignName,
			Description: horoscope.Description,
			Date:        horoscope.Date,
			IsActive:    horoscope.IsActive,
			CreatedOn:   horoscope.CreatedOn.Format("2006-01-02 15:04:05"),
			UpdatedBy:   horoscope.UpdatedBy.Format("2006-01-02 15:04:05"),
			AdminID:     horoscope.AdminID,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.HoroscopeListResponse{
		Success: true,
		Data:    horoscopeData,
		Pagination: dto.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (s *HoroscopeService) BulkCreateHoroscopes(ctx context.Context, payload dto.BulkHoroscopeRequestPayload, adminID string) (*dto.BulkHoroscopeResponse, error) {
	if len(payload) == 0 {
		return nil, fmt.Errorf("no horoscope data provided")
	}

	var createdIDs, updatedIDs []string

	for _, item := range payload {
		sign := normalizeSignName(item.SignName)
		if sign == "" || item.Description == "" || item.Date == "" {
			continue
		}

		existing, err := s.horoscopeRepo.FindBySignNameAndDate(ctx, sign, item.Date)
		if err == nil && existing != nil {
			update := map[string]interface{}{
				"signName":    sign,
				"description": item.Description,
				"date":        item.Date,
				"isAsctive":   1,
				"adminId":     adminID,
				"UpdatedBy":   time.Now(),
			}

			if err := s.horoscopeRepo.UpdateByHoroscopeID(ctx, existing.HoroscopeID, update); err != nil {
				return nil, err
			}
			updatedIDs = append(updatedIDs, existing.HoroscopeID)
			continue
		}

		// Generate UUID v7 for horoscopeId
		horoscopeUUID, err := uuid.NewV7()
		if err != nil {
			return nil, fmt.Errorf("failed to generate UUID v7: %w", err)
		}

		now := time.Now()
		entry := &models.Horoscope{
			HoroscopeID:    horoscopeUUID.String(),
			SignName:       sign,
			Description:    item.Description,
			Date:           item.Date,
			IsActive:       1,
			CreatedOn:      now,
			CreatedOnLower: now,
			UpdatedBy:      now,
			AdminID:        adminID,
		}

		if err := s.horoscopeRepo.Create(ctx, entry); err != nil {
			return nil, err
		}
		createdIDs = append(createdIDs, entry.HoroscopeID)
	}

	// Combine created and updated IDs
	allHoroscopeIDs := append(createdIDs, updatedIDs...)

	totalProcessed := len(allHoroscopeIDs)
	message := "Horoscope created successfully"
	if totalProcessed == 0 {
		message = "No horoscopes created or updated"
	}

	return &dto.BulkHoroscopeResponse{
		Success:     totalProcessed > 0,
		Message:     message,
		HoroscopeID: allHoroscopeIDs,
	}, nil
}

func (s *HoroscopeService) UpdateHoroscope(ctx context.Context, horoscopeID string, req *dto.UpdateHoroscopeRequest, adminID string) (*dto.HoroscopeResponse, error) {
	update := map[string]interface{}{
		"signName":    strings.ToLower(req.SignName),
		"description": req.Description,
		"date":        req.Date,
		"isAsctive":   1,
		"adminId":     adminID,
		"UpdatedBy":   time.Now(),
	}

	err := s.horoscopeRepo.UpdateByHoroscopeID(ctx, horoscopeID, update)
	if err != nil {
		return nil, err
	}

	return &dto.HoroscopeResponse{
		Success:     true,
		Message:     "Horoscope updated successfully",
		HoroscopeID: horoscopeID,
	}, nil
}

func (s *HoroscopeService) DeleteHoroscope(ctx context.Context, horoscopeID string) (*dto.HoroscopeResponse, error) {
	err := s.horoscopeRepo.DeleteByHoroscopeID(ctx, horoscopeID)
	if err != nil {
		return nil, err
	}

	return &dto.HoroscopeResponse{
		Success:     true,
		Message:     "Horoscope deleted successfully",
		HoroscopeID: horoscopeID,
	}, nil
}

// DashboardService handles dashboard metrics logic
type DashboardService struct {
	serviceRepo repository.ServiceRepository
}

func NewDashboardService(serviceRepo repository.ServiceRepository) *DashboardService {
	return &DashboardService{
		serviceRepo: serviceRepo,
	}
}

// GetDailyMetrics calculates daily metrics for chat, ivrCall, and videoCall grouped by lastStatus
// date format: YYYY-MM-DD (e.g., "2024-01-15")
func (s *DashboardService) GetDailyMetrics(ctx context.Context, date string) (*dto.DashboardMetricsResponse, error) {
	// Parse date
	targetDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	// Set time range for the entire day (00:00:00 to 23:59:59.999999999)
	startTime := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, time.UTC)
	endTime := startTime.Add(24 * time.Hour)

	// Get metrics for each collection
	chatMetrics, err := s.serviceRepo.CountByLastStatus(ctx, "chat", startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat metrics: %w", err)
	}

	ivrCallMetrics, err := s.serviceRepo.CountByLastStatus(ctx, "ivrCall", startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get ivrCall metrics: %w", err)
	}

	videoCallMetrics, err := s.serviceRepo.CountByLastStatus(ctx, "videoCall", startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get videoCall metrics: %w", err)
	}

	// Convert to ServiceMetrics DTO
	chatServiceMetrics := s.convertToServiceMetrics(chatMetrics)
	ivrCallServiceMetrics := s.convertToServiceMetrics(ivrCallMetrics)
	videoCallServiceMetrics := s.convertToServiceMetrics(videoCallMetrics)

	return &dto.DashboardMetricsResponse{
		Success: true,
		Data: dto.DashboardMetrics{
			Date:      date,
			Chat:      chatServiceMetrics,
			IvrCall:   ivrCallServiceMetrics,
			VideoCall: videoCallServiceMetrics,
		},
	}, nil
}

// convertToServiceMetrics converts map[string]int64 to ServiceMetrics
func (s *DashboardService) convertToServiceMetrics(metrics map[string]int64) dto.ServiceMetrics {
	serviceMetrics := dto.ServiceMetrics{
		Failed:   metrics["failed"],
		Request:  metrics["request"],
		Complete: metrics["complete"],
		Issue:    metrics["issue"],
		Reject:   metrics["reject"],
	}

	// Calculate total
	serviceMetrics.Total = serviceMetrics.Failed + serviceMetrics.Request + serviceMetrics.Complete + serviceMetrics.Issue + serviceMetrics.Reject

	return serviceMetrics
}

// FeedbackService handles feedback management logic
type FeedbackService struct {
	feedbackRepo repository.FeedbackRepository
}

func NewFeedbackService(feedbackRepo repository.FeedbackRepository) *FeedbackService {
	return &FeedbackService{
		feedbackRepo: feedbackRepo,
	}
}

func (s *FeedbackService) BulkCreateFeedbacks(ctx context.Context, payload dto.BulkFeedbackRequestPayload) (*dto.BulkFeedbackResponse, error) {
	if len(payload) == 0 {
		return nil, fmt.Errorf("no feedback data provided")
	}

	var feedbacks []*models.Feedback
	var feedbackIDs []string

	for _, item := range payload {
		now := time.Now()
		feedback := &models.Feedback{
			AstroID:    item.AstroID,
			Comment:    item.Comment,
			Name:       item.Name,
			Rating:     item.Rating,
			ProfilePic: item.ProfilePic,
			CreatedOn:  now,
			FeedbackID: item.FeedbackID,
		}

		feedbacks = append(feedbacks, feedback)
		feedbackIDs = append(feedbackIDs, feedback.FeedbackID)
	}

	// Bulk insert all feedbacks
	if err := s.feedbackRepo.BulkCreate(ctx, feedbacks); err != nil {
		return nil, fmt.Errorf("failed to create feedbacks: %w", err)
	}

	return &dto.BulkFeedbackResponse{
		Success:    true,
		Message:    "Feedbacks created successfully",
		FeedbackID: feedbackIDs,
	}, nil
}
