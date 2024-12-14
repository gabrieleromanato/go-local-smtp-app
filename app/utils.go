package app

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"mime/quotedprintable"
	"os"
	"strings"

	"net/smtp"
	"net/textproto"
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

func SendEmailViaSMTP(email Email, attachments []*multipart.FileHeader) error {
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_SERVER_HOST")
	smtpPort := os.Getenv("SMTP_SERVER_PORT")

	if smtpUser == "" || smtpPassword == "" || smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("SMTP credentials or server information are missing")
	}

	// SMTP server address
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	// Create the email buffer
	var emailBuffer bytes.Buffer
	writer := multipart.NewWriter(&emailBuffer)

	// Add headers
	headers := map[string]string{
		"From":         email.From,
		"To":           strings.Join(email.To, ", "),
		"Subject":      email.Subject,
		"MIME-Version": "1.0",
		"Content-Type": fmt.Sprintf("multipart/mixed; boundary=%s", writer.Boundary()),
	}

	for k, v := range headers {
		emailBuffer.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	emailBuffer.WriteString("\r\n")

	// Add email body as a part
	bodyHeader := textproto.MIMEHeader{}
	bodyHeader.Set("Content-Type", "text/plain; charset=utf-8")
	bodyHeader.Set("Content-Transfer-Encoding", "quoted-printable")
	bodyWriter, err := writer.CreatePart(bodyHeader)
	if err != nil {
		return fmt.Errorf("failed to create email body part: %w", err)
	}
	qpWriter := quotedprintable.NewWriter(bodyWriter)
	if _, err := qpWriter.Write([]byte(email.Body)); err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}
	qpWriter.Close()

	// Add attachments
	for _, fileHeader := range attachments {
		file, err := fileHeader.Open()
		if err != nil {
			return fmt.Errorf("failed to open attachment %s: %w", fileHeader.Filename, err)
		}
		defer file.Close()

		attachmentHeader := textproto.MIMEHeader{}
		attachmentHeader.Set("Content-Type", fileHeader.Header.Get("Content-Type"))
		attachmentHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileHeader.Filename))
		attachmentWriter, err := writer.CreatePart(attachmentHeader)
		if err != nil {
			return fmt.Errorf("failed to create attachment part for %s: %w", fileHeader.Filename, err)
		}

		if _, err := io.Copy(attachmentWriter, file); err != nil {
			return fmt.Errorf("failed to write attachment %s: %w", fileHeader.Filename, err)
		}
	}

	// Close the writer to finalize the email
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Auth for plain login
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	// Send the email
	if err := smtp.SendMail(addr, auth, email.From, email.To, emailBuffer.Bytes()); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
