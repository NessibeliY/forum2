package session

import (
	"database/sql"

	"forum/internal/models"
)

type SessionRepo struct {
	db *sql.DB
}

func NewSessionRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{
		db: db,
	}
}

type ISessionRepo interface {
	Insert(sessiong models.Session) error
	DeleteSessionByUser(user_id string) error
	GetSessionById(user_id string) (*models.Session, error)
	GetSessionByToken(token string) (*models.Session, error)
}

// insert new session in DB
func (s *SessionRepo) Insert(session models.Session) error {
	stmt := `INSERT INTO session(session_id,user_id,token,expire_time)VALUES(?,?,?,?)`
	if _, err := s.db.Exec(stmt, session.SessionId, session.UserId, session.Token, session.ExpireTime); err != nil {
		return err
	}
	return nil
}

func (s *SessionRepo) DeleteSessionByUser(user_id string) error {
	stmt := `DELETE FROM session WHERE user_id = ?`
	if _, err := s.db.Exec(stmt, user_id); err != nil {
		return err
	}
	return nil
}

// get session by id
func (s *SessionRepo) GetSessionById(user_id string) (*models.Session, error) {
	var session models.Session
	stmt := `SELECT * FROM session WHERE user_id = ?`
	row := s.db.QueryRow(stmt, user_id)

	err := row.Scan(&session.SessionId, &session.UserId, &session.Token, &session.ExpireTime)
	if err == sql.ErrNoRows {
		return nil, models.ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}

	return &session, nil
}

// get session by token
func (u *SessionRepo) GetSessionByToken(token string) (*models.Session, error) {
	var session models.Session
	stmt := `SELECT * FROM session WHERE token = ?`
	row := u.db.QueryRow(stmt, token)

	err := row.Scan(&session.SessionId, &session.UserId, &session.Token, &session.ExpireTime)
	if err == sql.ErrNoRows {
		return nil, models.ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}

	return &session, nil
}
