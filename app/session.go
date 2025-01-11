package app

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/jhillyerd/enmime"
)

type Session struct {
	store  *EmailStore
	from   string
	to     []string
	buffer strings.Builder
}

func (s *Session) CreateUser(email, password string) error {
	encPassword := HashString(password)
	_, err := s.store.Db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", email, encPassword)
	if err != nil {
		return fmt.Errorf("Error while creating the user: %w", err)
	}
	return nil
}

func (s *Session) UserExists(email, password string) bool {
	encPassword := HashString(password)
	var count int
	err := s.store.Db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? AND password = ?", email, encPassword).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func (s *Session) GetUserId(username, password string) (int, error) {
	encPassword := HashString(password)
	var id int
	err := s.store.Db.QueryRow("SELECT id FROM users WHERE email = ? AND password = ?", username, encPassword).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Mail accetta il mittente dell'email
func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.from = from
	return nil
}

// Rcpt accetta il destinatario dell'email
func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.to = append(s.to, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	maxAttachmentSize := os.Getenv("MAX_ATTACHMENT_SIZE")
	MAX_ATTACHMENT_SIZE, err := strconv.Atoi(maxAttachmentSize)
	env, err := enmime.ReadEnvelope(r)
	if err != nil {
		return fmt.Errorf("Error while parsing email: %w", err)
	}

	subject := env.GetHeader("Subject")
	body := env.Text
	bodyHTML := env.HTML
	userID := 1
	if env.GetHeader("X-User-Id") != "" {
		userID, err = strconv.Atoi(env.GetHeader("X-User-Id"))
		if err != nil {
			return fmt.Errorf("Error while parsing user ID: %w", err)
		}
	}
	emailId, err := s.store.SaveEmail(userID, s.from, s.to, subject, body, bodyHTML)
	if err != nil {
		return fmt.Errorf("Error while saving the email: %w", err)
	}

	for _, att := range env.Attachments {
		strContent := ConvertBytesToBase64(att.Content)
		if len(att.Content) > MAX_ATTACHMENT_SIZE {
			filename := fmt.Sprintf("attachments/%s", att.FileName)
			err := SaveAttachmentToFile(filename, att.Content)
			if err != nil {
				return fmt.Errorf("Error while saving the attachment: %w", err)
			}
			strContent = filename
		}
		err := s.store.SaveAttachment(emailId, att.ContentType, att.FileName, strContent)
		if err != nil {
			return fmt.Errorf("Error while saving the attachment: %w", err)
		}
	}

	log.Printf("Email saved: From=%s, To=%v, Subject=%s, Body=%s", s.from, s.to, subject, body)
	return nil
}

func (s *Session) Reset() {
	s.from = ""
	s.to = nil
	s.buffer.Reset()
}

func (s *Session) Logout() error {
	return nil
}

func (s *Session) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

func (s *Session) Auth(mech string) (sasl.Server, error) {
	return sasl.NewPlainServer(func(identity, username, password string) error {
		found := s.UserExists(username, password)
		if !found {
			return errors.New("Authentication failed")
		}
		return nil
	}), nil
}

func getUsernameAndPasswordFromFile() []User {
	data, err := os.ReadFile("authfile")
	if err != nil {
		log.Fatalf("Errore while reading the authfile: %v", err)
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) < 1 {
		log.Fatalf("Invalid authfile")
	}
	output := make([]User, 0)
	for _, line := range lines {
		str := strings.TrimSpace(line)
		parts := strings.Split(str, ":")
		if len(parts) != 2 {
			continue
		}
		output = append(output, User{
			Email:    parts[0],
			Password: parts[1],
		})
	}
	return output

}
