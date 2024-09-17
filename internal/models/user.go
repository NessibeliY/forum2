package models

import (
	"errors"

	"forum/internal/validator"
)

type contextKey string

const (
	// Define a key for the session
	SessionKey contextKey = "session"
)

var (
	ErrNotFound               error = errors.New("error: Not found")
	ErrUserNotFound           error = errors.New("user not found")
	ErrWrongPassword          error = errors.New("wrong password")
	ErrSessionNotFound        error = errors.New("session not found")
	ErrPostNotCreated         error = errors.New("post not created")
	ErrUUIDCreate             error = errors.New("id not generated")
	ErrNotCreated             error = errors.New("error: not created")
	ErrIncrementLikeInPost    error = errors.New("error: increment like in post")
	ErrDecrementLikeInPost    error = errors.New("error: decrement like in post")
	ErrDecrementDisLikeInPost error = errors.New("error: decrement dislike in post")
	ErrIncrementCommentInPost error = errors.New("error: increment comment in post")
	ErrDecrementCommentInpost error = errors.New("error: decrement comment in post")
	ErrDeletDisLikeInPost     error = errors.New("error: delete dislike in post")
	ErrDeleteLikeInPost       error = errors.New("error: delete like in post")
	ErrDeleteLikeInComment    error = errors.New("error: delete like in comment")
	ErrDeleteDisLikeInComment error = errors.New("error: delete dislike in comment")
	ErrDeletePost             error = errors.New("error: delete post")
	ErrDeleteComment          error = errors.New("error: delete comment")
)

type User struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UserName  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserRepository interface {
	Insert(user User) error
	GetUserByEmail(email string) (*User, error)
	GetUserUserID(id string) (*User, error)
}

type ErrorMessage struct {
	Email       string
	UserName    string
	Password    string
	Title       string
	Description string
	Tags        string
}

type UserData struct {
	UserName string
	Email    string
	Password string
	Errors   ErrorMessage
	IsAuth   bool
}

type AuthData struct {
	Email    string
	Password string
}

type LoginData struct {
	User   AuthData
	Errors ErrorMessage
	IsAuth bool
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(v.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(v.MinChars(password, 8), "password", "password contain at least 8 characters")
	v.Check(v.ValidPassword(password), "password", "password must be contain 1 upper character,1 lower character and 1 digit")
}

func ValidateUser(v *validator.Validator, u *User) {
	v.Check(u.UserName != "", "username", "must be provided")
	v.Check(len(u.UserName) > 3, "username", "username must be at least 3 characters long")

	ValidateEmail(v, u.Email)
	ValidatePassword(v, u.Password)
}
