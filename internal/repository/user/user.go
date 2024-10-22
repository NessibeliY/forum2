package user

import (
	"database/sql"
	"fmt"

	"forum/internal/models"
)

type UserRepo struct { //nolint:revive
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) AddUser(user models.User) error {
	stmt := `INSERT INTO users(id,created_at,username,email,password)VALUES(?,datetime('now','localtime'),?,?,?)`
	if _, err := u.db.Exec(stmt, user.ID, user.Username, user.Email, user.Password); err != nil {
		return fmt.Errorf("INSERT INTO DB: %w", err)
	}
	return nil
}

func (u *UserRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	stmt := `SELECT * FROM users WHERE email = ?`
	row := u.db.QueryRow(stmt, email)

	err := row.Scan(&user.ID, &user.Username, &user.CreatedAt, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, models.ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("row scan: %w", err)
	}

	return &user, nil
}

func (u *UserRepo) GetUserByUserID(id string) (*models.User, error) {
	var user models.User
	stmt := "SELECT * FROM users WHERE id = ?"
	row := u.db.QueryRow(stmt, id)
	err := row.Scan(&user.ID, &user.Username, &user.CreatedAt, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, models.ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("row scan: %w", err)
	}

	return &user, nil
}
