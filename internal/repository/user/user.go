package user

import (
	"database/sql"
	"fmt"

	"forum/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

// insert new user
func (u *UserRepo) Insert(user models.User) error {
	stmt := `INSERT INTO users(id,created_at,username,email,password)VALUES(?,datetime('now','localtime'),?,?,?)`
	if _, err := u.db.Exec(stmt, user.ID, user.UserName, user.Email, user.Password); err != nil {
		return fmt.Errorf("INSERT INTO DB: %v", err)
	}
	return nil
}

// get user by email
func (u *UserRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	stmt := `SELECT * FROM users WHERE email = ?`
	row := u.db.QueryRow(stmt, email)

	err := row.Scan(&user.ID, &user.UserName, &user.CreatedAt, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, models.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) GetUserUserID(id string) (*models.User, error) {
	var user models.User
	stmt := "SELECT * FROM users WHERE ID = ?"
	row := u.db.QueryRow(stmt, id)
	err := row.Scan(&user.ID, &user.UserName, &user.CreatedAt, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, models.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
