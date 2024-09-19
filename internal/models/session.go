package models

import "time"

type Session struct {
	SessionId  string
	UserId     string
	Token      string
	ExpireTime time.Time
}

type SessionRepository interface {
	Insert(sessiong Session) error
	DeleteSessionByUser(user_id string) error
	GetSessionById(user_id string) (*Session, error)
	GetSessionByToken(token string) (*Session, error)
}

type SessionService interface {
	CreateSession(userID, email string) (*Session, error)
	GetSessionByToken(token string) (*Session, error)
	IsSession(userID string) bool
	GetSessionByUserId(user_id string) (*Session, error)
	DeleteSessionByUser(user_id string) error
}
