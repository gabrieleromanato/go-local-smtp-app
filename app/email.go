package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

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

func populateUsersTableFromAuthFile(db *sql.DB, authFile string) error {
	lines := strings.Split(authFile, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid auth file format")
		}
		email := parts[0]
		password := parts[1]
		encPassword := HashString(password)
		_, err := db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", email, encPassword)
		if err != nil {
			return err
		}
	}
	return nil
}

func applyMigrations(db *sql.DB, migrationsDir string) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}

	sort.Strings(migrationFiles)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INT AUTO_INCREMENT PRIMARY KEY,
			migration_name VARCHAR(255) UNIQUE,
			run_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	for _, fileName := range migrationFiles {
		migrationPath := filepath.Join(migrationsDir, fileName)

		var count int
		query := "SELECT COUNT(*) FROM migrations WHERE migration_name = ?"
		if err := db.QueryRow(query, fileName).Scan(&count); err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}
		if count > 0 {
			continue
		}

		migrationContent, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}

		_, err = db.Exec(string(migrationContent))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", fileName, err)
		}

		_, err = db.Exec("INSERT INTO migrations (migration_name) VALUES (?)", fileName)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", fileName, err)
		}
	}

	return nil
}

func NewEmailStore(dsn string) (*EmailStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = applyMigrations(db, "./migrations")
	if err != nil {
		return nil, err
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		authFile, _ := os.ReadFile("./authfile")
		if err := populateUsersTableFromAuthFile(db, string(authFile)); err != nil {
			return nil, err
		}
	}

	return &EmailStore{Db: db}, nil
}

func (store *EmailStore) SaveEmail(from string, to []string, subject, body string) (int, error) {
	toString := strings.Join(to, ", ")
	query := `INSERT INTO emails (from_email, to_email, subject, body) VALUES (?, ?, ?, ?)`
	result, err := store.Db.Exec(query, from, toString, subject, body)
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
