package app

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Attachment struct {
	ID       int    `json:"id"`
	EmailID  int    `json:"email_id"`
	Type     string `json:"type"`
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type EmailStore struct {
	Db *sql.DB
}

type Email struct {
	ID          int          `json:"id"`
	From        string       `json:"from"`
	To          []string     `json:"to"`
	Subject     string       `json:"subject"`
	Body        string       `json:"body"`
	SentAt      string       `json:"sent_at"`
	Attachments []Attachment `json:"attachments"`
}

type EmailResponse struct {
	Emails []Email `json:"emails"`
	Pages  int     `json:"pages"`
	Page   int     `json:"page"`
}

func NewEmailStore(dsn string) (*EmailStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS emails (
		id INT AUTO_INCREMENT PRIMARY KEY,
		from_email VARCHAR(255),
		to_email TEXT,
		subject VARCHAR(255),
		body TEXT,
		sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	queryAttachments := `
	CREATE TABLE IF NOT EXISTS attachments (
		id INT AUTO_INCREMENT PRIMARY KEY,
		email_id INT,
		type VARCHAR(255),
		filename VARCHAR(255),
		content TEXT,
		FOREIGN KEY (email_id) REFERENCES emails(id) ON DELETE CASCADE
	)`
	_, err = db.Exec(queryAttachments)
	if err != nil {
		return nil, err
	}

	return &EmailStore{Db: db}, nil
}

func (store *EmailStore) SaveEmail(from string, to []string, subject, body string) (int, error) {
	toString := strings.Join(to, ", ")
	sentAt := time.Now()
	query := `INSERT INTO emails (from_email, to_email, subject, body, sent_at) VALUES (?, ?, ?, ?, ?)`
	result, err := store.Db.Exec(query, from, toString, subject, body, sentAt)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return int(id), nil
}

func (store *EmailStore) SaveAttachment(emailId int, attachmentType, filename, content string) error {
	query := `INSERT INTO attachments (email_id, type, filename, content) VALUES (?, ?, ?, ?)`
	_, err := store.Db.Exec(query, emailId, attachmentType, filename, content)
	return err
}

func (store *EmailStore) DeleteEmail(id int) error {
	query := "DELETE FROM emails WHERE id = ?"
	_, err := store.Db.Exec(query, id)
	return err
}

func (store *EmailStore) GetEmailAttachments(id int) ([]Attachment, error) {
	query := "SELECT * FROM attachments WHERE email_id = ?"
	rows, err := store.Db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attachments := []Attachment{}
	for rows.Next() {
		var id, emailId int
		var attachmentType, filename, content string
		rows.Scan(&id, &emailId, &attachmentType, &filename, &content)
		attachment := Attachment{ID: id, EmailID: emailId, Type: attachmentType, Filename: filename, Content: content}
		attachments = append(attachments, attachment)
	}

	return attachments, nil
}

func (store *EmailStore) ListEmails(page int) (EmailResponse, error) {
	emailsPerPage := os.Getenv("EMAILS_PER_PAGE")
	if emailsPerPage == "" {
		emailsPerPage = "10"
	}
	perPage, _ := strconv.Atoi(emailsPerPage)
	offset := (page - 1) * perPage
	total := 0
	err := store.Db.QueryRow("SELECT COUNT(*) FROM emails").Scan(&total)
	if err != nil {
		return EmailResponse{}, err
	}
	pages := (total + perPage - 1) / perPage

	query := "SELECT * FROM emails ORDER BY sent_at DESC LIMIT ? OFFSET ?"
	rows, err := store.Db.Query(query, perPage, offset)
	if err != nil {
		return EmailResponse{}, err
	}
	defer rows.Close()

	emails := []Email{}
	for rows.Next() {
		var from, to, subject, body, sentAt string
		var id int
		rows.Scan(&id, &from, &to, &subject, &body, &sentAt)
		attachments, _ := store.GetEmailAttachments(id)
		email := Email{
			ID:          id,
			From:        from,
			To:          strings.Split(to, ", "),
			Subject:     subject,
			Body:        body,
			SentAt:      sentAt,
			Attachments: attachments}
		emails = append(emails, email)
	}

	resp := EmailResponse{Emails: emails, Pages: pages, Page: page}
	return resp, nil
}

func (store *EmailStore) SearchEmails(query string, page int) (EmailResponse, error) {
	emailsPerPage := os.Getenv("EMAILS_PER_PAGE")
	if emailsPerPage == "" {
		emailsPerPage = "10"
	}
	perPage, _ := strconv.Atoi(emailsPerPage)
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage
	total := 0
	err := store.Db.QueryRow("SELECT COUNT(*) FROM emails WHERE subject LIKE ? OR body LIKE ?", "%"+query+"%", "%"+query+"%").Scan(&total)
	if err != nil {
		return EmailResponse{}, err
	}

	pages := (total + perPage - 1) / perPage

	sqlQuery := "SELECT * FROM emails WHERE subject LIKE ? OR body LIKE ? ORDER BY sent_at DESC LIMIT ? OFFSET ?"
	rows, err := store.Db.Query(sqlQuery, "%"+query+"%", "%"+query+"%", perPage, offset)
	if err != nil {
		return EmailResponse{}, err
	}
	defer rows.Close()

	emails := []Email{}
	for rows.Next() {
		var from, to, subject, body, sentAt string
		var id int
		if err := rows.Scan(&id, &from, &to, &subject, &body, &sentAt); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		attachments, _ := store.GetEmailAttachments(id)
		email := Email{
			ID:          id,
			From:        from,
			To:          strings.Split(to, ", "),
			Subject:     subject,
			Body:        body,
			SentAt:      sentAt,
			Attachments: attachments}
		emails = append(emails, email)
	}

	resp := EmailResponse{Emails: emails, Pages: pages, Page: page}
	return resp, nil
}
