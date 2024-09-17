package session

import (
	"time"

	"forum/internal/models"
	"forum/internal/repository/session"

	"github.com/gofrs/uuid"
)

type SessionService struct {
	SessionRepo session.ISessionRepo
}

func NewSessionService(sessionRepo session.ISessionRepo) *SessionService {
	return &SessionService{
		SessionRepo: sessionRepo,
	}
}

type ISessionService interface {
	CreateSession(userID, email string) (*models.Session, error)
	GetSessionByToken(token string) (*models.Session, error)
	IsSession(userID string) bool
	GetSessionByUserId(user_id string) (*models.Session, error)
	DeleteSessionByUser(user_id string) error
}

func (s *SessionService) CreateSession(userID, email string) (*models.Session, error) {
	// generate unique id for session
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	// generate token for session
	token, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	// expire time for session
	expire_time := time.Now().Add(time.Hour * 2)
	session := models.Session{
		SessionId:  id.String(),
		UserId:     userID,
		Token:      token.String(),
		ExpireTime: expire_time,
	}

	if err = s.SessionRepo.Insert(session); err != nil {
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

func (s *SessionService) IsSession(userID string) bool {
	_, err := s.SessionRepo.GetSessionById(userID)
	if err != nil {
		return false
	}

	return true
}

func (s *SessionService) GetSessionByUserId(user_id string) (*models.Session, error) {
	session, err := s.SessionRepo.GetSessionById(user_id)
	if err != nil {
		return nil, err
	}

	return session, err
}

func (s *SessionService) DeleteSessionByUser(user_id string) error {
	err := s.SessionRepo.DeleteSessionByUser(user_id)
	if err != nil {
		return err
	}

	return nil
}