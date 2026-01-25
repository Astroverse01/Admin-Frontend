package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string             `bson:"userId" json:"userId"`
	Name        string             `bson:"name,omitempty" json:"name"`
	FullName    string             `bson:"fullName,omitempty" json:"-"`
	Email       string             `bson:"email" json:"email"`
	IsDelete    int                `bson:"isDelete" json:"isDelete"`
	TotalAmount float64            `bson:"totalAmount" json:"totalAmount"`
	FCMToken    string             `bson:"fcmToken" json:"fcmToken"`
	Platform    string             `bson:"platform" json:"platform"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// Astrologer represents an astrologer in the system
type Astrologer struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AstroID     string             `bson:"astroId" json:"astroId"`
	Name        string             `bson:"name,omitempty" json:"name"`
	FullName    string             `bson:"fullName,omitempty" json:"-"`
	Email       string             `bson:"email" json:"email"`
	IsDelete    int                `bson:"isDelete" json:"isDelete"`
	IsActive    int                `bson:"isActive" json:"isActive"`
	TotalEarned float64            `bson:"totalEarned" json:"totalEarned"`
	FCMToken    string             `bson:"fcmToken" json:"fcmToken"`
	Platform    string             `bson:"platform" json:"platform"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// ServiceReport represents a service complaint
type ServiceReport struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ReportID      string             `bson:"reportId" json:"reportId"`
	ServiceType   string             `bson:"serviceType" json:"serviceType"`
	OrderID       string             `bson:"orderId" json:"orderId"`
	AstroID       string             `bson:"astroId" json:"astroId"`
	UserID        string             `bson:"userId" json:"userId"`
	Status        string             `bson:"status" json:"status"`
	Comment       string             `bson:"comment" json:"comment"`
	IssueResolved int                `bson:"issueResolved" json:"issueResolved"`
	Reason        string             `bson:"reason" json:"reason"`
	CreatedOn     time.Time          `bson:"createdOn" json:"createdOn"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// UserProblem represents a general user complaint
type UserProblem struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProblemID    string             `bson:"problemId" json:"problemId"`
	UserID       string             `bson:"userId" json:"userId"`
	Status       string             `bson:"status" json:"status"`
	Comment      string             `bson:"comment" json:"comment"`
	ProblemTypes string             `bson:"problemTypes" json:"problemTypes"`
	Reason       string             `bson:"reason" json:"reason"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// AstroProblem represents a general astrologer complaint
type AstroProblem struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProblemID    string             `bson:"problemId" json:"problemId"`
	AstroID      string             `bson:"astroId" json:"astroId"`
	Status       string             `bson:"status" json:"status"`
	Comment      string             `bson:"comment" json:"comment"`
	ProblemTypes string             `bson:"problemTypes" json:"problemTypes"`
	Reason       string             `bson:"reason" json:"reason"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// Horoscope represents a horoscope entry
type Horoscope struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	HoroscopeID    string             `bson:"horoscopeId" json:"horoscopeId"`
	CreatedOn      time.Time          `bson:"CreatedOn" json:"CreatedOn"`
	AdminID        string             `bson:"adminId" json:"adminId"`
	UpdatedBy      time.Time          `bson:"UpdatedBy" json:"UpdatedBy"`
	SignName       string             `bson:"signName" json:"signName"`
	Description    string             `bson:"description" json:"description"`
	Date           string             `bson:"date" json:"date"`
	IsActive       int                `bson:"isAsctive" json:"isAsctive"`
	CreatedOnLower time.Time          `bson:"createdOn" json:"createdOn"`
}

// Chat represents a chat service record
type Chat struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ChatId          string             `bson:"chatId" json:"chatId"`
	AstroID         string             `bson:"astroId" json:"astroId"`
	UserID          string             `bson:"userId" json:"userId"`
	Conversation    []interface{}      `bson:"conversation" json:"conversation"`
	RatePerMinute   float64            `bson:"ratePerMinute" json:"ratePerMinute"`
	SpendMoney      float64            `bson:"spendMoney" json:"spendMoney"`
	PaymentReceived int                `bson:"paymentReceived" json:"paymentReceived"`
	LastStatus      string             `bson:"lastStatus" json:"lastStatus"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
}

// IvrCall represents an IVR service record
type IvrCall struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	IvrID           string             `bson:"ivrId" json:"ivrId"`
	AstroID         string             `bson:"astroId" json:"astroId"`
	UserID          string             `bson:"userId" json:"userId"`
	URL             string             `bson:"url" json:"url"`
	UserUrl         string             `bson:"userUrl" json:"userUrl"`
	RatePerMinute   float64            `bson:"ratePerMinute" json:"ratePerMinute"`
	SpendMoney      float64            `bson:"spendMoney" json:"spendMoney"`
	PaymentReceived int                `bson:"paymentReceived" json:"paymentReceived"`
	LastStatus      string             `bson:"lastStatus" json:"lastStatus"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
}

// VideoCall represents a video service record
type VideoCall struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	VideoID         string             `bson:"videoId" json:"videoId"`
	AstroID         string             `bson:"astroId" json:"astroId"`
	UserID          string             `bson:"userId" json:"userId"`
	URL             string             `bson:"url" json:"url"`
	RatePerMinute   float64            `bson:"ratePerMinute" json:"ratePerMinute"`
	SpendMoney      float64            `bson:"spendMoney" json:"spendMoney"`
	PaymentReceived int                `bson:"paymentReceived" json:"paymentReceived"`
	LastStatus      string             `bson:"lastStatus" json:"lastStatus"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
}

// Feedback represents a feedback entry
type Feedback struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AstroID    string             `bson:"astroId" json:"astroId"`
	Comment    string             `bson:"comment" json:"comment"`
	Name       string             `bson:"name" json:"name"`
	Rating     int                `bson:"rating" json:"rating"`
	ProfilePic string             `bson:"profilePic" json:"profilePic"`
	CreatedOn  time.Time          `bson:"createdOn" json:"createdOn"`
	FeedbackID string             `bson:"feedbackId" json:"feedbackId"`
}
