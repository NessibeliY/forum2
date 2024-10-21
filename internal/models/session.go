package models

import "time"

type Session struct {
	SessionId  string
	UserID     string
	Token      string
	ExpireTime time.Time
}

type SessionRepository interface {
	AddSession(session Session) error
	DeleteSessionByUserID(user_id string) error
	GetSessionById(user_id string) (*Session, error)
	GetSessionByToken(token string) (*Session, error)
}

type SessionService interface {
	SetSession(UserID string) (*Session, error)
	GetSessionByToken(token string) (*Session, error)
	GetSessionByUserId(user_id string) (*Session, error)
}
