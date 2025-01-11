package app

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

func GetDSN() string {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
}

func FormatDateToMySQL(date string) string {
	return date[:10] + " " + date[11:19]
}

func FormatMySQLDateToLocale(date string) string {
	// Parsing the date string into a Go time.Time object
	dateObj, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		// Handle the error appropriately (returning an empty string here for simplicity)
		return ""
	}

	// Get the timezone from the environment variable (e.g., Europe/Rome)
	timeZone := os.Getenv("TZ")
	if timeZone == "" {
		timeZone = "UTC" // Fallback to UTC if TZ is not set
	}

	// Load the location based on the timezone
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		// Handle the error appropriately (fallback to UTC here)
		loc = time.UTC
	}

	// Convert the time to the specified timezone
	dateObj = dateObj.In(loc)

	// Retrieve the display format from the environment variable
	displayFormat := os.Getenv("DISPLAY_DATE_FORMAT")
	if displayFormat == "" {
		displayFormat = "dd/mm/yyyy HH:MM:SS" // Default format
	}

	// Replace quotes in the environment variable if any
	displayFormat = strings.ReplaceAll(displayFormat, "\"", "")

	// Replace placeholders with Go's date format specifiers
	displayFormat = strings.ReplaceAll(displayFormat, "dd", "02")
	displayFormat = strings.ReplaceAll(displayFormat, "mm", "01")
	displayFormat = strings.ReplaceAll(displayFormat, "yyyy", "2006")
	displayFormat = strings.ReplaceAll(displayFormat, "HH", "15")
	displayFormat = strings.ReplaceAll(displayFormat, "MM", "04")
	displayFormat = strings.ReplaceAll(displayFormat, "SS", "05")

	// Format the date and return it
	return dateObj.Format(displayFormat)
}

func HashString(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func ConvertBytesToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func SaveAttachmentToFile(filename string, data []byte) error {
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func SendEmailViaSMTP(email Email, attachments []string, userId int) error {
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_SERVER_HOST")
	smtpPort := os.Getenv("SMTP_SERVER_PORT")

	if smtpUser == "" || smtpPassword == "" || smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("SMTP credentials or server information are missing")
	}

	// Convert port to integer
	smtpPortInt, err := strconv.Atoi(smtpPort)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %w", err)
	}

	// Create a new email message
	message := gomail.NewMessage()
	message.SetHeader("From", email.From)
	message.SetHeader("To", email.To...)
	message.SetHeader("Subject", email.Subject)
	message.SetHeader("X-User-Id", strconv.Itoa(userId))
	message.SetBody("text/plain", email.Body)
	message.AddAlternative("text/html", email.BodyHTML)

	// Attach files
	for _, attachment := range attachments {
		message.Attach(attachment)
	}

	// Create a new SMTP dialer
	dialer := gomail.NewDialer(smtpHost, smtpPortInt, smtpUser, smtpPassword)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// Remove attachments
	for _, attachment := range attachments {
		err := os.Remove(attachment)
		if err != nil {
			return fmt.Errorf("failed to remove attachment: %w", err)
		}
	}

	return nil
}

func IsValidPassword(password string) bool {
	minPasswordLength, _ := strconv.Atoi(os.Getenv("MIN_PASSWORD_LENGTH"))
	if len(password) < minPasswordLength {
		return false
	}
	return true
}
