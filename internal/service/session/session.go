package session

import (
	"time"

	"github.com/gofrs/uuid"

	"forum/internal/models"
)

type SessionService struct {
	SessionRepo models.SessionRepository
}

func NewSessionService(sessionRepo models.SessionRepository) *SessionService {
	return &SessionService{
		SessionRepo: sessionRepo,
	}
}

func (s *SessionService) SetSession(UserID string) (*models.Session, error) {
	err := s.SessionRepo.DeleteSessionByUserID(UserID)
	if err != nil {
		return nil, err
	}

	session, err := s.createSession(UserID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *SessionService) createSession(UserID string) (*models.Session, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	token, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	expire_time := time.Now().Add(time.Hour * 2)
	session := models.Session{
		SessionId:  id.String(),
		UserID:     UserID,
		Token:      token.String(),
		ExpireTime: expire_time,
	}

	if err = s.SessionRepo.AddSession(session); err != nil {
		return nil, err
	}

	new_session, err := s.SessionRepo.GetSessionByToken(token.String())
	if err != nil {
		return nil, err
	}
	return new_session, err
}

func (s *SessionService) GetSessionByToken(token string) (*models.Session, error) {
	session, err := s.SessionRepo.GetSessionByToken(token)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionService) GetSessionByUserId(user_id string) (*models.Session, error) {
	session, err := s.SessionRepo.GetSessionById(user_id)
	if err != nil {
		return nil, err
	}

	return session, err
}
