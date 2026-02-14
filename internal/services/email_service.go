package services

import (
	"fmt"
	"log"
	"net/smtp"

	"admin-be/internal/config"
)

type EmailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		cfg: cfg,
	}
}

// validateEmailConfig validates that EMAIL_USER and EMAIL_PASS are provided from .env
func (s *EmailService) validateEmailConfig() error {
	if s.cfg.EmailUser == "" {
		return fmt.Errorf("EMAIL_USER is required in .env file")
	}
	if s.cfg.EmailPass == "" {
		return fmt.Errorf("EMAIL_PASS is required in .env file")
	}
	return nil
}

// SendNotificationEmail sends a notification email to a recipient
func (s *EmailService) SendNotificationEmail(to, subject, body string) error {
	if to == "" {
		log.Printf("[EmailService] Skipping email notification: recipient email is empty")
		return nil
	}

	// Validate email configuration
	if err := s.validateEmailConfig(); err != nil {
		return fmt.Errorf("email configuration error: %w", err)
	}

	// Gmail SMTP settings
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	emailFrom := s.cfg.EmailUser

	// Setup authentication
	auth := smtp.PlainAuth("", s.cfg.EmailUser, s.cfg.EmailPass, smtpHost)

	// Create simple message without attachments
	message := fmt.Sprintf("From: %s\r\n", emailFrom)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/plain; charset=UTF-8\r\n"
	message += "\r\n"
	message += body + "\r\n"

	// Send email
	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	if err := smtp.SendMail(addr, auth, emailFrom, []string{to}, []byte(message)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("[EmailService] Successfully sent notification email to %s", to)
	return nil
}

// SendNotificationEmailFrom sends a notification email with a custom From address (e.g. support@astrosway.com).
// Auth still uses cfg.EmailUser/EmailPass; From is used in message headers and envelope.
func (s *EmailService) SendNotificationEmailFrom(from, to, subject, body string) error {
	if to == "" {
		log.Printf("[EmailService] Skipping email notification: recipient email is empty")
		return nil
	}
	if from == "" {
		from = s.cfg.EmailUser
	}
	if err := s.validateEmailConfig(); err != nil {
		return fmt.Errorf("email configuration error: %w", err)
	}
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	auth := smtp.PlainAuth("", s.cfg.EmailUser, s.cfg.EmailPass, smtpHost)
	message := fmt.Sprintf("From: %s\r\n", from)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/plain; charset=UTF-8\r\n"
	message += "\r\n"
	message += body + "\r\n"
	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	if err := smtp.SendMail(addr, auth, from, []string{to}, []byte(message)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	log.Printf("[EmailService] Successfully sent notification email from %s to %s", from, to)
	return nil
}
