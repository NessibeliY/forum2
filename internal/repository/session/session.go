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

// insert new session in DB
func (s *SessionRepo) AddSession(session models.Session) error {
	stmt := `INSERT INTO session(session_id,user_id,token,expire_time)VALUES(?,?,?,?)`
	if _, err := s.db.Exec(stmt, session.SessionId, session.UserId, session.Token, session.ExpireTime); err != nil {
		return err
	}
	return nil
}

func (s *SessionRepo) DeleteSessionByUserID(user_id string) error {
	stmt := `DELETE FROM session WHERE user_id = ?`
	if _, err := s.db.Exec(stmt, user_id); err != nil {
		return err
	}
	return nil
}

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
