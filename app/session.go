package app

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
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
	// Usa enmime per analizzare l'email
	env, err := enmime.ReadEnvelope(r)
	if err != nil {
		return fmt.Errorf("errore durante il parsing dell'email: %w", err)
	}

	// Estrai i dati
	subject := env.GetHeader("Subject") // Ottieni il Subject dall'intestazione
	body := env.Text                    // Corpo in formato testo semplice
	if body == "" {
		body = env.HTML // Se il testo semplice è vuoto, prova il corpo HTML
	}

	// Salva l'email nel database
	emailId, err := s.store.SaveEmail(s.from, s.to, subject, body)
	if err != nil {
		return fmt.Errorf("errore durante il salvataggio dell'email: %w", err)
	}

	// Salva il contenuto dell'allegato
	for _, att := range env.Attachments {
		strContent := ConvertBytesToBase64(att.Content)
		if len(att.Content) > MAX_ATTACHMENT_SIZE {
			// Salva su file
			filename := fmt.Sprintf("attachments/%d_%s", emailId, att.FileName)
			err := SaveAttachmentToFile(filename, att.Content)
			if err != nil {
				return fmt.Errorf("errore durante il salvataggio dell'allegato: %w", err)
			}
			strContent = filename
		}
		err := s.store.SaveAttachment(emailId, att.ContentType, att.FileName, strContent)
		if err != nil {
			return fmt.Errorf("errore durante il salvataggio dell'allegato: %w", err)
		}
	}

	log.Printf("Email salvata nel database: From=%s, To=%v, Subject=%s, Body=%s", s.from, s.to, subject, body)
	return nil
}

// Reset resetta la sessione
func (s *Session) Reset() {
	s.from = ""
	s.to = nil
	s.buffer.Reset()
}

// Logout termina la sessione
func (s *Session) Logout() error {
	return nil
}

func (s *Session) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

func (s *Session) Auth(mech string) (sasl.Server, error) {
	users := getUsernameAndPasswordFromFile()
	return sasl.NewPlainServer(func(identity, username, password string) error {
		found := false
		for _, user := range users {
			if username == user.Email && Md5String(password) == user.Password {
				found = true
				break
			}
		}
		if !found {
			return errors.New("Autenticazione fallita")
		}
		return nil
	}), nil
}

func getUsernameAndPasswordFromFile() []User {
	data, err := os.ReadFile("authfile")
	if err != nil {
		log.Fatalf("Errore nella lettura del file: %v", err)
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) < 1 {
		log.Fatalf("File non valido")
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
