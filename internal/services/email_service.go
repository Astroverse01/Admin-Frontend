package services

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"time"

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

// SendCSVFilesByEmail sends multiple CSV files as email attachments
func (s *EmailService) SendCSVFilesByEmail(csvFilePaths []string, date string) error {
	if len(csvFilePaths) == 0 {
		log.Println("[EmailService] No CSV files to send")
		return nil
	}

	// Validate email configuration from .env
	if err := s.validateEmailConfig(); err != nil {
		return fmt.Errorf("email configuration error: %w. Please ensure EMAIL_USER and EMAIL_PASS are configured in .env file", err)
	}

	// Gmail SMTP settings (matching Nodemailer Gmail service)
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	emailFrom := s.cfg.EmailUser
	emailTo := s.cfg.EmailUser // Send to same address as requested

	// Setup authentication
	auth := smtp.PlainAuth("", s.cfg.EmailUser, s.cfg.EmailPass, smtpHost)

	// Prepare email
	subject := fmt.Sprintf("Daily Report - %s", date)
	body := fmt.Sprintf("Please find attached the daily reports for %s.\n\nTotal collections: %d\n\nThis is an automated email.", date, len(csvFilePaths))

	// Create message
	message := s.createMessage(emailFrom, emailTo, subject, body, csvFilePaths)

	// Send email
	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	if err := smtp.SendMail(addr, auth, emailFrom, []string{emailTo}, message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("[EmailService] Successfully sent email with %d CSV attachments", len(csvFilePaths))
	return nil
}

// createMessage creates the email message with attachments
func (s *EmailService) createMessage(from, to, subject, body string, attachments []string) []byte {
	boundary := "----=_NextPart_" + generateBoundary()

	message := fmt.Sprintf("From: %s\r\n", from)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", boundary)
	message += "\r\n"

	// Add body
	message += fmt.Sprintf("--%s\r\n", boundary)
	message += "Content-Type: text/plain; charset=UTF-8\r\n"
	message += "Content-Transfer-Encoding: 7bit\r\n"
	message += "\r\n"
	message += body + "\r\n"
	message += "\r\n"

	// Add attachments
	for _, filePath := range attachments {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			log.Printf("[EmailService] Warning: File does not exist: %s", filePath)
			continue
		}

		fileData, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("[EmailService] Warning: Failed to read file %s: %v", filePath, err)
			continue
		}

		fileName := filepath.Base(filePath)
		message += fmt.Sprintf("--%s\r\n", boundary)
		message += fmt.Sprintf("Content-Type: text/csv; name=\"%s\"\r\n", fileName)
		message += "Content-Transfer-Encoding: base64\r\n"
		message += fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", fileName)
		message += "\r\n"
		// Wrap base64 content at 76 characters per line (MIME standard)
		message += wrapBase64(encodeBase64(fileData)) + "\r\n"
		message += "\r\n"
	}

	message += fmt.Sprintf("--%s--\r\n", boundary)

	return []byte(message)
}

// generateBoundary generates a unique boundary string for MIME multipart messages
func generateBoundary() string {
	return fmt.Sprintf("%d_%d", os.Getpid(), time.Now().Unix())
}

// encodeBase64 encodes data to base64
func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// wrapBase64 wraps base64 string at 76 characters per line (MIME standard)
func wrapBase64(data string) string {
	const lineLength = 76
	var result string
	for i := 0; i < len(data); i += lineLength {
		end := i + lineLength
		if end > len(data) {
			end = len(data)
		}
		result += data[i:end] + "\r\n"
	}
	return result
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
