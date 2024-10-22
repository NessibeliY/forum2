package models

import "time"

type Session struct {
	SessionID  string
	UserID     string
	Token      string
	ExpireTime time.Time
}

type SessionRepository interface {
	AddSession(session Session) error
	DeleteSessionByUserID(userID string) error
	GetSessionByID(userID string) (*Session, error)
	GetSessionByToken(token string) (*Session, error)
}

type SessionService interface {
	SetSession(userID string) (*Session, error)
	GetSessionByToken(token string) (*Session, error)
	GetSessionByUserID(userID string) (*Session, error)
}
