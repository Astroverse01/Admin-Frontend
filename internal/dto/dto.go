package dto

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

// UserListResponse represents the user list response
type UserListResponse struct {
	Success    bool       `json:"success"`
	Data       []UserData `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// UserData represents user data in list response
type UserData struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// AstroListResponse represents the astro list response
type AstroListResponse struct {
	Success    bool        `json:"success"`
	Data       []AstroData `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// AstroData represents astro data in list response
type AstroData struct {
	AstroID string `json:"astroId"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Visible string `json:"visible"`
}

// ComplaintListResponse represents complaint list response
type ComplaintListResponse struct {
	Success    bool            `json:"success"`
	Data       []ComplaintData `json:"data"`
	Pagination Pagination      `json:"pagination"`
}

// ComplaintData represents complaint data in list response
type ComplaintData struct {
	ServiceType string `json:"serviceType"`
	OrderID     string `json:"orderId"`
	AstroID     string `json:"astroId"`
	UserID      string `json:"userId"`
	CreatedOn   string `json:"createdOn"`
	Status      string `json:"status"`
	Comment     string `json:"comment"`
	ReportID    string `json:"reportId"`
}

// ServiceDetailsResponse represents service details response
type ServiceDetailsResponse struct {
	Success bool        `json:"success"`
	Data    ServiceData `json:"data"`
}

// ServiceData represents service data
type ServiceData struct {
	ServiceID     string        `json:"serviceId"`
	AstroID       string        `json:"astroId"`
	UserID        string        `json:"userId"`
	Conversation  []interface{} `json:"conversation,omitempty"`
	URL           string        `json:"url,omitempty"`
	RatePerMinute float64       `json:"ratePerMinute"`
	Type          string        `json:"type"`
	ReportID      string        `json:"reportId"`
}

// AcceptRejectRequest represents accept/reject complaint request
type AcceptRejectRequest struct {
	Action          string  `json:"action" validate:"required,oneof=accept reject"`
	Reason          string  `json:"reason" validate:"required"`
	AstroRefundMoney float64 `json:"astroRefundMoney"`
	UserRefundMoney  float64 `json:"userRefundMoney"`
	LegitimateTime  int     `json:"legitimateTime,omitempty"`
}

// AcceptRejectResponse represents accept/reject response
type AcceptRejectResponse struct {
	Success        bool    `json:"success"`
	Message        string  `json:"message"`
	ComplaintID    string  `json:"complaintId"`
	UserID         string  `json:"userId"`
	AstroID        string  `json:"astroId"`
	RefundedAmount float64 `json:"refundedAmount,omitempty"`
}

// ProblemListResponse represents problem list response
type ProblemListResponse struct {
	Success    bool          `json:"success"`
	Data       []ProblemData `json:"data"`
	Pagination Pagination    `json:"pagination"`
}

// ProblemData represents problem data
type ProblemData struct {
	ProblemID    string `json:"problemId"`
	UserID       string `json:"userId,omitempty"`
	AstroID      string `json:"astroId,omitempty"`
	Status       string `json:"status"`
	Comment      string `json:"comment"`
	ProblemTypes string `json:"problemTypes"`
}

// CloseComplaintRequest represents close complaint request
type CloseComplaintRequest struct {
	Reason string `json:"reason" validate:"required"`
}

// CloseComplaintResponse represents close complaint response
type CloseComplaintResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	ProblemID string `json:"problemId"`
	Status    string `json:"status"`
}

// HoroscopeListResponse represents horoscope list response
type HoroscopeListResponse struct {
	Success    bool            `json:"success"`
	Data       []HoroscopeData `json:"data"`
	Pagination Pagination      `json:"pagination"`
}

// HoroscopeData represents horoscope data
type HoroscopeData struct {
	HoroscopeID string `json:"horoscopeId"`
	SignName    string `json:"signName"`
	Description string `json:"description"`
	Date        string `json:"date"`
	IsActive    int    `json:"isAsctive"`
	CreatedOn   string `json:"CreatedOn"`
	UpdatedBy   string `json:"UpdatedBy"`
	AdminID     string `json:"adminId"`
}

// UpdateHoroscopeRequest represents update horoscope request
type UpdateHoroscopeRequest struct {
	SignName    string `json:"signName" validate:"required"`
	Description string `json:"description" validate:"required"`
	Date        string `json:"date" validate:"required"`
	IsActive    int    `json:"isAsctive" validate:"required"`
}

// HoroscopeResponse represents horoscope response
type HoroscopeResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	HoroscopeID string `json:"horoscopeId,omitempty"`
}

// BulkHoroscopeRequestItem represents a single horoscope item in bulk request
type BulkHoroscopeRequestItem struct {
	SignName    string `json:"signName" validate:"required"`
	Description string `json:"description" validate:"required"`
	Date        string `json:"date" validate:"required"`
	IsActive    int    `json:"isAsctive" validate:"required"`
}

// BulkHoroscopeRequestPayload represents the incoming array of horoscope items
type BulkHoroscopeRequestPayload []BulkHoroscopeRequestItem

// BulkHoroscopeResponse represents bulk create/update feedback
type BulkHoroscopeResponse struct {
	Success     bool     `json:"success"`
	Message     string   `json:"message"`
	HoroscopeID []string `json:"horoscopeId"`
}

// Pagination represents pagination info
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// StatusUpdateRequest represents status update request
type StatusUpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=active inactive"`
}

// VisibilityUpdateRequest represents visibility update request
type VisibilityUpdateRequest struct {
	Visible bool `json:"visible" validate:"required"`
}

// StatusUpdateResponse represents status update response
type StatusUpdateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      string `json:"id"`
}
