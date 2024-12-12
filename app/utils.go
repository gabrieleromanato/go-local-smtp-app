package app

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"net/smtp"
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

func SendEmailViaSMTP(email Email) error {
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_SERVER_HOST")
	smtpPort := os.Getenv("SMTP_SERVER_PORT")

	if smtpUser == "" || smtpPassword == "" || smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("SMTP credentials or server information are missing")
	}

	// SMTP server address
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	// Join the recipients as a comma-separated string for the "To" header
	toHeader := strings.Join(email.To, ", ")

	// Create the email message
	message := fmt.Sprintf("From: %s\r\n", email.From) +
		fmt.Sprintf("To: %s\r\n", toHeader) +
		fmt.Sprintf("Subject: %s\r\n\r\n", email.Subject) +
		email.Body

	// Auth for plain login
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	// Send the email
	err := smtp.SendMail(addr, auth, email.From, email.To, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
