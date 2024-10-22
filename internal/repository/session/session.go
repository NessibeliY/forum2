package session

import (
	"database/sql"
	"fmt"

	"forum/internal/models"
)

type SessionRepo struct { //nolint:revive
	db *sql.DB
}

func NewSessionRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{
		db: db,
	}
}

func (s *SessionRepo) AddSession(session models.Session) error {
	stmt := `INSERT INTO sessions(session_id,user_id,token,expire_time)VALUES(?,?,?,?)`
	_, err := s.db.Exec(stmt, session.SessionID, session.UserID, session.Token, session.ExpireTime)
	return err //nolint:wrapcheck
}

func (s *SessionRepo) DeleteSessionByUserID(userID string) error {
	stmt := `DELETE FROM sessions WHERE user_id = ?`
	_, err := s.db.Exec(stmt, userID)
	return err //nolint:wrapcheck
}

func (s *SessionRepo) GetSessionByID(userID string) (*models.Session, error) {
	var session models.Session
	stmt := `SELECT * FROM sessions WHERE user_id = ?`
	row := s.db.QueryRow(stmt, userID)

	err := row.Scan(&session.SessionID, &session.UserID, &session.Token, &session.ExpireTime)
	if err == sql.ErrNoRows {
		return nil, models.ErrSessionNotFound
	} else if err != nil {
		return nil, fmt.Errorf("row scan: %w", err)
	}

	return &session, nil
}

func (s *SessionRepo) GetSessionByToken(token string) (*models.Session, error) {
	var session models.Session
	stmt := `SELECT * FROM sessions WHERE token = ?`
	row := s.db.QueryRow(stmt, token)

	err := row.Scan(&session.SessionID, &session.UserID, &session.Token, &session.ExpireTime)
	if err == sql.ErrNoRows {
		return nil, models.ErrSessionNotFound
	} else if err != nil {
		return nil, fmt.Errorf("row scan: %w", err)
	}

	return &session, nil
}
