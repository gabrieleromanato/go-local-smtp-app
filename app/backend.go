package app

import "github.com/emersion/go-smtp"

type Backend struct {
	Store *EmailStore
}

func (b *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{
		store: b.Store,
	}, nil
}
