package session

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"forum/internal/models"
)

type SessionService struct { //nolint:revive
	SessionRepo models.SessionRepository
}

func NewSessionService(sessionRepo models.SessionRepository) *SessionService {
	return &SessionService{
		SessionRepo: sessionRepo,
	}
}

func (s *SessionService) SetSession(userID string) (*models.Session, error) {
	err := s.SessionRepo.DeleteSessionByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("delete session: %w", err)
	}

	session, err := s.createSession(userID)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return session, nil
}

func (s *SessionService) createSession(userID string) (*models.Session, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("generate uuid: %w", err)
	}

	token, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("generate token uuid: %w", err)
	}

	expireTime := time.Now().Add(time.Hour * 2)
	session := models.Session{
		SessionID:  id.String(),
		UserID:     userID,
		Token:      token.String(),
		ExpireTime: expireTime,
	}

	if err = s.SessionRepo.AddSession(session); err != nil {
		return nil, fmt.Errorf("add session: %w", err)
	}

	newSession, err := s.SessionRepo.GetSessionByToken(token.String())
	if err != nil {
		return nil, fmt.Errorf("get session by token: %w", err)
	}
	return newSession, nil
}

func (s *SessionService) GetSessionByToken(token string) (*models.Session, error) {
	session, err := s.SessionRepo.GetSessionByToken(token)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	return session, nil
}

func (s *SessionService) GetSessionByUserID(userID string) (*models.Session, error) {
	session, err := s.SessionRepo.GetSessionByID(userID)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return session, nil
}

func (s *SessionService) DeleteSessionByToken(token string) error {
	return s.SessionRepo.DeleteSessionByUserID(token) //nolint:wrapcheck
}
