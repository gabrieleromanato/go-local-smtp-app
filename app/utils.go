package app

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

func FormatDateToMySQL(date string) string {
	return date[:10] + " " + date[11:19]
}

func FormatMySQLDateToLocale(date string) string {
	parts := strings.Split(date, " ")
	datePart := parts[0]
	timePart := parts[1]
	dateParts := strings.Split(datePart, "-")
	year := dateParts[0]
	month := dateParts[1]
	day := dateParts[2]
	return fmt.Sprintf("%s/%s/%s %s", day, month, year, timePart)
}

func Md5String(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
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

func SendEmailViaSMTP(email Email, attachments []string) error {
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
	message.SetBody("text/plain", email.Body)

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
